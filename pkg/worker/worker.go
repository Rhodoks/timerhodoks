package worker

import "time"

// worker接口
// 本地goroutine集群和远程grpc调用都可以作为worker
// 可以限制同时执行的工作，防止突发高压任务潮拉起过多的go routine
type Worker interface {
	TryAcquire() bool                                                                                          //尝试获取信号量
	Execute(JobId uint64, TriggerTime time.Time, RetryNum int, ExecutorType string, ExecutorInfo string) error //执行
	Release()                                                                                                  //释放信号量
}
