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

package raftstore

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"timerhodoks/pkg/coordinator"
	"timerhodoks/pkg/job"
	"timerhodoks/pkg/jobstore"
	"timerhodoks/pkg/raftnode"
	"timerhodoks/pkg/scheduler"

	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.etcd.io/raft/v3/raftpb"
)

// 负责调度器, 协调器, api接口, 状态机和raft层的业务交互
type RaftStore struct {
	proposeC    chan<- string
	Mutex       sync.RWMutex
	Jobs        jobstore.JobStoreInterface
	snapshotter *snap.Snapshotter
	scheduler   *scheduler.Scheduler
	coordinator *coordinator.Coordinator
}

func NewRaftStore(snapshotter *snap.Snapshotter, scheduler *scheduler.Scheduler, coordinator *coordinator.Coordinator,
	proposeC chan<- string, commitC <-chan *raftnode.Commit, errorC <-chan error) *RaftStore {
	s := &RaftStore{
		proposeC:    proposeC,
		Jobs:        jobstore.NewJobStore(),
		Mutex:       sync.RWMutex{},
		snapshotter: snapshotter,
		scheduler:   scheduler,
		coordinator: coordinator,
	}
	s.scheduler.ProposeCommitTime = s.ProposeCommitTime
	snapshot, err := s.loadSnapshot()
	if err != nil {
		log.Panic(err)
	}
	if snapshot != nil {
		log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
		if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
			log.Panic(err)
		}
	}
	go s.loopReadCommits(commitC, errorC) //开始读取commit
	go s.loopCheckCoordinator()           //开始监视任务分配变化
	return s
}

// 将更新任务的请求解析并送入raft提交channel，以下类似
func (s *RaftStore) ProposeUpdate(v []byte) error {
	raftlog := RaftLogUpdateJob{}
	err := raftlog.parse(v)
	if err != nil {
		log.Println("Invalid job update")
		return err
	}
	var args RaftLog = &raftlog
	var buf strings.Builder
	if err := gob.NewEncoder(&buf).Encode(&args); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
	return nil
}

func (s *RaftStore) ProposeInsert(v []byte) error {
	raftlog := RaftLogCreateJob{CommitTime: time.Now()}
	err := raftlog.parse(v)
	if err != nil {
		log.Println("Invalid job insert")
		return err
	}
	var args RaftLog = &raftlog
	var buf strings.Builder
	if err := gob.NewEncoder(&buf).Encode(&args); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
	return nil
}

func (s *RaftStore) ProposeDelete(v []byte) error {
	raftlog := RaftLogDeleteJob{}
	err := raftlog.parse(v)
	if err != nil {
		log.Println("Invalid job delete")
		return err
	}
	var args RaftLog = &raftlog
	var buf strings.Builder
	if err := gob.NewEncoder(&buf).Encode(&args); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
	return nil
}

func (s *RaftStore) ProposeCommitTime(id uint64, commitTime int64, success bool) {
	raftlog := RaftLogUpdateTime{
		Id:         id,
		CommitTime: commitTime,
		Success:    success,
	}
	var args RaftLog = &raftlog
	var buf strings.Builder
	if err := gob.NewEncoder(&buf).Encode(&args); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
}

// 查找给定id对应的任务
func (s *RaftStore) Lookup(key uint64) (string, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	job := s.Jobs.Find(key)
	if job != nil {
		return "", false
	}
	buf, _ := json.Marshal(job)
	return string(buf), true
}

// 更新操作
// 首先对存储和调度器上锁
// 然后尝试更新到存储和调度器（如果任务分配给自己了）中
// 以下类似
func (s *RaftStore) Update(log *RaftLogUpdateJob) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	job := s.Jobs.Find(log.Id)
	if job == nil {
		return errors.New("do not find the job to update")
	}
	err := json.Unmarshal([]byte(log.Data), job)

	if err != nil {
		return err
	}
	s.coordinator.Mutex.RLock()
	defer s.coordinator.Mutex.RUnlock()
	if s.coordinator.Query(job.Id) {
		s.scheduler.Mutex.Lock()
		defer s.scheduler.Mutex.Unlock()
		s.scheduler.Remove(job.Id)
		s.scheduler.Insert(job)
	}
	return nil
}

func (s *RaftStore) Insert(job *job.JobEntry) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	err := s.Jobs.Insert(job)
	if err != nil {
		return err
	}
	s.coordinator.Mutex.RLock()
	defer s.coordinator.Mutex.RUnlock()
	if s.coordinator.Query(job.Id) {
		s.scheduler.Mutex.Lock()
		defer s.scheduler.Mutex.Unlock()
		s.scheduler.Insert(job)
	}
	return nil
}

func (s *RaftStore) Delete(id uint64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	job := s.Jobs.Delete(id)
	if job == nil {
		return errors.New("do not find the job need to delete")
	}
	s.coordinator.Mutex.RLock()
	defer s.coordinator.Mutex.RUnlock()
	if s.coordinator.Query(id) {
		s.scheduler.Mutex.Lock()
		defer s.scheduler.Mutex.Unlock()
		s.scheduler.Remove(job.Id)
	}
	return nil
}

// 不断读取raft传递而来需要提交的日志，并进行应用
func (s *RaftStore) loopReadCommits(commitC <-chan *raftnode.Commit, errorC <-chan error) {
	for commit := range commitC {
		if commit == nil {
			// signaled to load snapshot
			snapshot, err := s.loadSnapshot()
			if err != nil {
				log.Panic(err)
			}
			if snapshot != nil {
				log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
				if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		for _, data := range commit.Data {
			var rlog RaftLog
			dec := gob.NewDecoder(bytes.NewBufferString(data))
			if err := dec.Decode(&rlog); err != nil {
				log.Fatalf("raftexample: could not decode message (%v)", err)
			}
			rlog.apply(s)
		}
		close(commit.ApplyDoneC)
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

// 检查任务分配是否变更
// 如果有，变更调度器中的任务
func (s *RaftStore) checkCoordinator() bool {
	if s.coordinator.UpdateFlag == 0 {
		return false
	}
	s.Mutex.RLock()
	s.coordinator.Mutex.Lock()
	s.scheduler.Mutex.Lock()
	defer s.scheduler.Mutex.Unlock()
	defer s.coordinator.Mutex.Unlock()
	defer s.Mutex.RUnlock()
	for bucId := uint(0); bucId < job.HASH_BUC_NUM; bucId++ {
		if s.coordinator.PlanJobAllocation.Test(bucId) && !s.coordinator.ApplyJobAllocation.Test(bucId) { //被分配的哈希桶，但是并没有更新
			for _, job := range *s.Jobs.GetJobMap(bucId) {
				s.scheduler.Insert(job)
			}
		} else if !s.coordinator.PlanJobAllocation.Test(bucId) && s.coordinator.ApplyJobAllocation.Test(bucId) { //被取消分配的哈希桶，但是并没有更新
			for id := range *s.Jobs.GetJobMap(bucId) {
				s.scheduler.Remove(id)
			}
		}
	}
	s.coordinator.ApplyJobAllocation = s.coordinator.PlanJobAllocation
	s.coordinator.PlanJobAllocation = coordinator.NewJobAllocation()
	s.coordinator.ApplyJobAllocation.CopyFull(s.coordinator.PlanJobAllocation)
	s.coordinator.UpdateFlag = 0
	return true
}

const CHECK_ALLOC_INTERVAL_MS = 500

// 定期检查协调器
func (s *RaftStore) loopCheckCoordinator() bool {
	for {
		time.Sleep(CHECK_ALLOC_INTERVAL_MS * time.Millisecond)
		s.checkCoordinator()
	}
}

// 生成快照
func (s *RaftStore) GetSnapshot() ([]byte, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return json.Marshal(s.Jobs)
}

// 加载快照
func (s *RaftStore) loadSnapshot() (*raftpb.Snapshot, error) {
	snapshot, err := s.snapshotter.Load()
	if err == snap.ErrNoSnapshot {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

// 从快照中恢复状态机，并清空调度器和协调器
func (s *RaftStore) recoverFromSnapshot(snapshot []byte) error {
	var store jobstore.JobStore
	if err := json.Unmarshal(snapshot, &store); err != nil {
		return err
	}
	s.Mutex.Lock()
	s.coordinator.Mutex.Lock()
	defer s.coordinator.Mutex.Unlock()
	defer s.Mutex.Unlock()

	s.coordinator.UpdateFlag = 0
	s.coordinator.PlanJobAllocation.ClearAll() //清空分配给自己的任务，然后清空调度器里的任务

	s.Jobs = &store

	return nil
}
