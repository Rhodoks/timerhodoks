package jobstore

import (
	"errors"
	"timerhodoks/pkg/job"
)

// raft状态机
// 存储所有任务的状态
// 按照哈希存储在不同的哈希桶中
type JobStore struct {
	Jobs  [job.HASH_BUC_NUM]map[uint64]*job.JobEntry //job按照哈希值分到不同的桶中
	MaxId uint64                                     //当前所使用最大ID
}

func NewJobStore() *JobStore {
	jobstore := JobStore{
		MaxId: 0,
	}
	for i := 0; i < job.HASH_BUC_NUM; i++ {
		jobstore.Jobs[i] = make(map[uint64]*job.JobEntry)
	}
	return &jobstore
}

// 获取一个未被使用的ID
func (store *JobStore) GetNewId() uint64 {
	store.MaxId++
	return store.MaxId
}

// 插入一个JobEntry，如果已有相同id，则返回error
func (store *JobStore) Insert(e *job.JobEntry) error {
	bucId := e.Hash
	if _, ok := store.Jobs[bucId][e.Id]; ok {
		return errors.New("insert a existing entry")
	}
	store.Jobs[bucId][e.Id] = e
	return nil
}

// 按照key查找JobEntry
func (store *JobStore) Find(id uint64) *job.JobEntry {
	bucId := job.Hash(id)
	return store.Jobs[bucId][id]
}

// 按照key删除JobEntry
func (store *JobStore) Delete(id uint64) *job.JobEntry {
	bucId := job.Hash(id)
	res, ok := store.Jobs[bucId][id]
	if ok {
		delete(store.Jobs[bucId], id)
		return res
	} else {
		return nil
	}
}

// 获取任务条目数
func (store *JobStore) Size() int {
	sum := 0
	for _, hashBuc := range store.Jobs {
		sum += len(hashBuc)
	}
	return sum
}

// 获取任务列表
func (store *JobStore) GetJobList(start int, end int) []*job.JobEntry {
	res := make([]*job.JobEntry, 0)
	cnt := 0
	for _, hashBuc := range store.Jobs {
		if cnt+len(hashBuc) <= start || cnt > end {
			cnt += len(hashBuc)
			continue
		}
		for _, job := range hashBuc {
			if cnt >= start && cnt <= end {
				res = append(res, job)
			}
			cnt++
		}
	}
	return res
}

// 获取某个哈希桶中的所有任务
func (store *JobStore) GetJobMap(hashBuc uint) *map[uint64]*job.JobEntry {
	return &store.Jobs[hashBuc]
}
