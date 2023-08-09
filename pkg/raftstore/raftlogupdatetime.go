package raftstore

import (
	"errors"
)

// 更新任务的完成时间，需要保证提交的时间点前的所有任务都已经执行完毕
type RaftLogUpdateTime struct {
	Id         uint64
	Success    bool
	CommitTime int64
}

func (log *RaftLogUpdateTime) apply(store *RaftStore) error {
	store.Mutex.Lock()
	defer store.Mutex.Unlock()
	job := store.Jobs.Find(log.Id)
	if job == nil {
		return errors.New("do not find the job to update time")
	}
	if log.Success {
		if log.CommitTime > job.LastSuccessTime {
			job.LastSuccessTime = log.CommitTime
		}
	} else {
		if log.CommitTime > job.LastFailureTime {
			job.LastFailureTime = log.CommitTime
		}
	}
	return nil
}
