package main

import (
	"flag"
	"timerhodoks/pkg/worker"
)

func main() {
	grpcPort := flag.Int("grpcPort", 9023, "grpc server port")
	server := worker.GrpcWorkerServer{}
	server.ListenAndServe(*grpcPort)
}
