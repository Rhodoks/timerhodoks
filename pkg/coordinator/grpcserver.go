package coordinator

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "timerhodoks/pkg/coordinatorpb"

	"github.com/bits-and-blooms/bitset"
	"google.golang.org/grpc"
)

// 返回任务分配信息
func (s *Coordinator) GetCoordinatingInfo(ctx context.Context, in *pb.EmptyMessage) (*pb.CoordinatingInfoReply, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	res := pb.CoordinatingInfoReply{Load: 1, JobAllocation: uint64ToByte(s.PlanJobAllocation.Bytes())}
	return &res, nil
}

// 接受Leader的任务分配
func (s *Coordinator) UpdateAllocation(ctx context.Context, in *pb.AllocationRequest) (*pb.EmptyMessage, error) {
	uint64slice, err := byteToUint64(in.JobAllocation)
	if err != nil {
		return nil, err
	}
	plan := bitset.From(uint64slice)
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.PlanJobAllocation = plan
	s.UpdateFlag += 1
	return &pb.EmptyMessage{}, nil
}

func (s *Coordinator) ListenAndServe(grpcPort int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("grpc failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterCoordinatorServer(server, s)
	log.Printf("grpc server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
