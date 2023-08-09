package worker

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"timerhodoks/pkg/executor"
	pb "timerhodoks/pkg/workerpb"

	"google.golang.org/grpc"
)

type GrpcWorkerServer struct {
	pb.UnimplementedGrpcWorkerServer
}

// grpc调用执行任务
func (s *GrpcWorkerServer) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteReply, error) {
	exector, err := executor.ParseExecutor(req.ExecutorType, req.ExecutorInfo)
	if err != nil {
		return &pb.ExecuteReply{Ok: 0}, nil
	}
	err = exector.Execute(uint64(req.JobId), time.Unix(int64(req.TriggerTime), 0))
	if err != nil {
		return &pb.ExecuteReply{Ok: 0}, nil
	}
	return &pb.ExecuteReply{Ok: 1}, nil
}

func (s *GrpcWorkerServer) ListenAndServe(grpcPort int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("grpc failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterGrpcWorkerServer(server, s)
	log.Printf("grpc server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
