package raftstore

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"time"
	"timerhodoks/pkg/job"
)

// 初始化gob
func GobInit() {
	//所有日志通过gob序列化存储到raft里
	//因为gob支持序列化接口类型
	gob.Register(&RaftLogCreateJob{})
	gob.Register(&RaftLogUpdateJob{})
	gob.Register(&RaftLogDeleteJob{})
	gob.Register(&RaftLogUpdateTime{})
}

type RaftLog interface {
	apply(*RaftStore) error
}

// 插入任务日志
type RaftLogCreateJob struct {
	Name         string
	CronLine     string
	ExecutorType string
	ExecutorInfo string
	RetryNum     int
	CommitTime   time.Time
	Cnt          int
}

func (log *RaftLogCreateJob) parse(buf []byte) error {
	log.Cnt = 1
	err := json.Unmarshal(buf, log)
	if err != nil {
		return err
	}

	// 检查是否合法：尝试构造JobEntry
	_, err = job.NewJobEntry(0, log.CronLine, log.Name, log.RetryNum, log.ExecutorType, log.ExecutorInfo, log.CommitTime)

	return err
}

func (log *RaftLogCreateJob) apply(store *RaftStore) error {
	for i := 0; i < log.Cnt; i++ {
		job, err := job.NewJobEntry(store.Jobs.GetNewId(), log.CronLine, log.Name, log.RetryNum, log.ExecutorType, log.ExecutorInfo, log.CommitTime)
		if err != nil {
			return err
		}
		store.Insert(job)
	}
	return nil
}

// 修改一条已有的任务
type RaftLogUpdateJob struct {
	Id   uint64
	Data string
}

func (log *RaftLogUpdateJob) parse(buf []byte) error {
	err := json.Unmarshal(buf, log)
	if err != nil {
		return err
	}
	if log.Id == 0 {
		return errors.New("do not have field Id or Id = 0")
	}
	return nil
}

func (log *RaftLogUpdateJob) apply(store *RaftStore) error {
	err := store.Update(log)
	return err
}

// 删除一条已有的任务
type RaftLogDeleteJob struct {
	Id uint64
}

func (log *RaftLogDeleteJob) parse(buf []byte) error {
	err := json.Unmarshal(buf, log)
	return err
}

func (log *RaftLogDeleteJob) apply(store *RaftStore) error {
	err := store.Delete(log.Id)
	if err != nil {
		return errors.New("do not find the job to delete")
	}
	return nil
}
