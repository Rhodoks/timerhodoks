// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.etcd.io/raft/v3/raftpb"

	"timerhodoks/pkg/config"
	"timerhodoks/pkg/coordinator"
	"timerhodoks/pkg/httpapi"
	"timerhodoks/pkg/raftnode"
	"timerhodoks/pkg/raftstore"
	"timerhodoks/pkg/scheduler"
	"timerhodoks/pkg/web"
	"timerhodoks/pkg/worker"
)

func init() {
	raftstore.GobInit()
}

const PROPOSE_CHAN_BUFFER = 1000
const MAX_GOROUTINE_NUM_PER_WORKER = 1000

func main() {
	workers_conf := flag.String("workers", "", "url of workers' grpc server, seperated by comma")
	config_path := flag.String("config", `./config.json`, "file path to config json")
	id := flag.Int("id", 1, "node ID, 1-indexed")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	// 处理worker信息
	fmt.Println(*workers_conf)
	workers_urls := strings.SplitAfter(*workers_conf, ",")
	workers := []worker.Worker{worker.NewLocalWorker(MAX_GOROUTINE_NUM_PER_WORKER)}
	for _, url := range workers_urls {
		if url == "" {
			continue
		}
		workers = append(workers, worker.NewGrpcWorker(MAX_GOROUTINE_NUM_PER_WORKER, url))
	}

	// 处理兄弟集群信息
	fileData, err := os.ReadFile(*config_path)
	if err != nil {
		log.Fatal(err)
	}
	configs, err := config.ParseConfigsFromJson(fileData)
	if err != nil {
		log.Fatal(err)
	}
	raftCluster := make([]string, 0)
	grpcCluster := make([]string, 0)
	for _, conf := range configs.Cluster {
		raftCluster = append(raftCluster, `http://`+conf.Ip+":"+strconv.Itoa(conf.RaftPort))
		grpcCluster = append(grpcCluster, conf.Ip+":"+strconv.Itoa(conf.GrpcPort))
	}

	// 主服务
	proposeC := make(chan string, PROPOSE_CHAN_BUFFER)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	var store *raftstore.RaftStore
	scheduler := scheduler.NewScheduler(workers)
	getSnapshot := func() ([]byte, error) { return store.GetSnapshot() }
	commitC, errorC, snapshotterReady, raftNode := raftnode.NewRaftNode(*id, raftCluster, *join, getSnapshot, proposeC, confChangeC)

	coordinator := coordinator.NewCoordinator(scheduler, grpcCluster, raftNode, configs.Cluster[*id-1].GrpcPort)
	store = raftstore.NewRaftStore(<-snapshotterReady, scheduler, coordinator, proposeC, commitC, errorC)
	scheduler.ProposeCommitTime = store.ProposeCommitTime
	go scheduler.Running()

	// api和web服务
	http.Handle("/api/job", &httpapi.JobHttpAPI{Store: store})
	http.Handle("/api/jobs", &httpapi.JobsHttpAPI{Store: store})
	http.Handle("/", &web.IndexServer{RaftNode: raftNode})
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.ListenAndServe(":"+strconv.Itoa(configs.Cluster[*id-1].HttpPort), nil)
}
