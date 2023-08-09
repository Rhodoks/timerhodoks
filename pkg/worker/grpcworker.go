package worker

import (
	"context"
	"errors"
	"time"

	pb "timerhodoks/pkg/workerpb"

	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcWorker struct {
	client *pb.GrpcWorkerClient
	cancel context.CancelFunc
	ctx    context.Context
	server string
	seme   *semaphore.Weighted
}

func NewGrpcWorker(limit int64, server string) *GrpcWorker {
	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	c := pb.NewGrpcWorkerClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	return &GrpcWorker{
		server: server,
		cancel: cancel,
		ctx:    ctx,
		client: &c,
		seme:   semaphore.NewWeighted(limit),
	}
}

func (worker *GrpcWorker) TryAcquire() bool {
	return worker.seme.TryAcquire(1)
}

func (worker *GrpcWorker) Execute(JobId uint64, TriggerTime time.Time, RetryNum int, ExecutorType string, ExecutorInfo string) error {
	var err error

	for i := 0; i <= RetryNum; i++ {
		reply, err := (*worker.client).Execute(worker.ctx, &pb.ExecuteRequest{
			JobId:        uint32(JobId),
			TriggerTime:  uint32(TriggerTime.Unix()),
			ExecutorType: ExecutorType,
			ExecutorInfo: ExecutorInfo,
		})
		if err == nil && reply.Ok == 1 { // grpc调用成功且执行成功
			return nil
		}
	}
	if err != nil {
		return err
	} else {
		return errors.New("error in executing")
	}
}

func (worker *GrpcWorker) Release() {
	worker.seme.Release(1)
}

func (worker *GrpcWorker) Exit() {
	worker.cancel()
}
