package jobstore

import (
	"fmt"
	"testing"
	"time"
	"timerhodoks/pkg/job"
)

const INSERT_NUM = 1500000

func TestJobStore(t *testing.T) {
	jobstore := NewJobStore()

	for i := 0; i < INSERT_NUM; i++ {
		entry, _ := job.NewJobEntry(jobstore.GetNewId(), "* * * * * * *", "test", 1, "Shell", "test", time.Now())
		_ = jobstore.Insert(entry)
	}
	sum := 0
	max := 0
	min := INSERT_NUM
	for i := 0; i < job.HASH_BUC_NUM; i++ {
		x := len(jobstore.Jobs[i])
		sum += x
		if x > max {
			max = x
		}
		if x < min {
			min = x
		}
	}
	if sum != INSERT_NUM {
		t.Errorf("incorresponding number of jobs")
	}
	if max-min > 1 {
		t.Errorf("hash algorithm go wrong")
	}
	// 因为哈希函数里PRIME是质数，所以理论上哈希值应该是回环遍历的
}

func TestJobStoreMemoryLeak(t *testing.T) {
	jobstore := NewJobStore()

	for i := 0; i < INSERT_NUM; i++ {
		entry, _ := job.NewJobEntry(jobstore.GetNewId(), "* * * * * * *", "test", 1, "Shell", "test", time.Now())
		_ = jobstore.Insert(entry)
	}
	for i := 1; i <= INSERT_NUM; i++ {
		jobstore.Delete(uint64(i))
	}
	fmt.Println(1)
	for i := 0; i < INSERT_NUM; i++ {
		entry, _ := job.NewJobEntry(jobstore.GetNewId(), "* * * * * * *", "test", 1, "Shell", "test", time.Now())
		_ = jobstore.Insert(entry)
	}
	for i := 1; i <= INSERT_NUM; i++ {
		jobstore.Delete(uint64(i + INSERT_NUM))
	}

	for i := 1; i <= INSERT_NUM; i++ {
		jobstore.Find(uint64(i + 10*INSERT_NUM))
	}
	time.Sleep(5 * time.Second)
	// 应该内存为空
	fmt.Println(len(jobstore.Jobs[0]))
}
