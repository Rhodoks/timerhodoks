package job

import (
	"errors"
	"time"

	"github.com/gorhill/cronexpr"
)

// raft协议上存储的任务条目
type JobEntry struct {
	Id              uint64 //任务ID
	Hash            uint64 //任务哈希（放在哪个哈希桶）
	CronLine        string //任务的Cron表达式
	Name            string //任务名
	RetryNum        int    //任务重试次数
	ExecutorType    string //任务执行器
	ExecutorInfo    string //任务执行参数
	CreateTime      int64  //任务创建时间
	LastSuccessTime int64  //上一次成功时间
	LastFailureTime int64  //上一次失败时间
}

func NewJobEntry(id uint64, cronLine string, name string, retryNum int, executerType string, executorInfo string, commitTime time.Time) (*JobEntry, error) {
	_, err := cronexpr.Parse(cronLine) //校验cron表达式
	if err != nil {
		return nil, err
	}

	if retryNum < 0 {
		return nil, errors.New("ReteyNum must be non-negative")
	}

	return &JobEntry{
		Id:              id,
		Hash:            Hash(id),
		CronLine:        cronLine,
		Name:            name,
		RetryNum:        retryNum,
		ExecutorType:    executerType,
		ExecutorInfo:    executorInfo,
		CreateTime:      commitTime.Unix(),
		LastSuccessTime: 0,
		LastFailureTime: 0,
	}, nil
}
