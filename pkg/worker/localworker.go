package worker

import (
	"time"
	"timerhodoks/pkg/executor"

	"golang.org/x/sync/semaphore"
)

// 本地go routine工作者
type LocalWorker struct {
	seme *semaphore.Weighted
}

func NewLocalWorker(limit int64) *LocalWorker {
	return &LocalWorker{
		seme: semaphore.NewWeighted(limit),
	}
}

func (worker *LocalWorker) TryAcquire() bool {
	return worker.seme.TryAcquire(1)
}

func (worker *LocalWorker) Execute(JobId uint64, TriggerTime time.Time, RetryNum int, ExecutorType string, ExecutorInfo string) error {
	executor, err := executor.ParseExecutor(ExecutorType, ExecutorInfo)
	if err != nil {
		return err
	}
	for i := 0; i <= RetryNum; i++ {
		err = executor.Execute(uint64(JobId), TriggerTime)
		if err == nil {
			return nil
		}
	}
	return err
}

func (worker *LocalWorker) Release() {
	worker.seme.Release(1)
}

func (worker *LocalWorker) Exit() {

}
