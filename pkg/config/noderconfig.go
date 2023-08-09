package config

import "encoding/json"

type ClusterConfig struct {
	Id       int
	Ip       string
	RaftPort int
	HttpPort int
	GrpcPort int
}

type NodeConfig struct {
	Cluster []ClusterConfig
}

// 从json中解析兄弟节点信息
func ParseConfigsFromJson(buf []byte) (*NodeConfig, error) {
	res := NodeConfig{
		Cluster: make([]ClusterConfig, 0),
	}
	err := json.Unmarshal(buf, &res)
	return &res, err
}
