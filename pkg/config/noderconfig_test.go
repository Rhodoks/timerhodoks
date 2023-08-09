package config

import (
	"testing"
)

var Json string = `{"Cluster":[{"Id":1,"Ip":"127.0.0.1","RaftPort":18000,"HttpPort":18001,"GrpcPort":18002},{"Id":2,"Ip":"127.0.0.1","RaftPort":28000,"HttpPort":28001,"GrpcPort":28002},{"Id":3,"Ip":"127.0.0.1","RaftPort":38000,"HttpPort":38001,"GrpcPort":38002}]}`

func TestNodeConfig(t *testing.T) {
	config, err := ParseConfigsFromJson([]byte(Json))
	if err != nil {
		t.Errorf("fail to parse: %v", err)
	}
	if len(config.Cluster) != 3 {
		t.Errorf("wrong config length")
	}
	if config.Cluster[2].Id != 3 || config.Cluster[2].Ip != "127.0.0.1" || config.Cluster[2].GrpcPort != 38002 {
		t.Errorf("wrong config data")
	}
}
