package jobstore

import "timerhodoks/pkg/job"

// raft状态机接口
type JobStoreInterface interface {
	GetNewId() uint64
	Size() int
	Insert(e *job.JobEntry) error
	Find(id uint64) *job.JobEntry
	Delete(id uint64) *job.JobEntry
	GetJobList(start int, end int) []*job.JobEntry
	GetJobMap(id uint) *map[uint64]*job.JobEntry
}
