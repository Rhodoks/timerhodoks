package executor

import (
	"errors"
	"time"
)

// 执行器接口
type Executor interface {
	Execute(jobId uint64, triggerTime time.Time) error
	Parse(buf []byte) error
}

// 解析执行器
func ParseExecutor(executorType string, executorinfo string) (Executor, error) {
	if executorType == "Shell" {
		return ParseExecutorShell([]byte(executorinfo))
	}
	if executorType == "Http" {
		return ParseExecutorHttp([]byte(executorinfo))
	}
	return nil, errors.New("can not find the corresponding executor type")
}
