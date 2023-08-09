package scheduler

import (
	"container/heap"
	"testing"
	"time"
	"timerhodoks/pkg/job"
)

func TestSchedulerHeap(t *testing.T) {
	job0, _ := job.NewJobEntry(1, "0 0 0 * * * 2099", "begin_of_minute", 1, "test", "test", time.Now())
	job1, _ := job.NewJobEntry(1, "0 * * 2 * * 2099", "begin_of_minute", 1, "test", "test", time.Now())
	job2, _ := job.NewJobEntry(1, "30 * * * * * 2099", "middle_of_minute", 1, "test", "test", time.Now())
	scheduler := NewScheduler(nil)
	heap.Push(scheduler, NewSchedulerNode(job0))
	heap.Push(scheduler, NewSchedulerNode(job1))
	heap.Push(scheduler, NewSchedulerNode(job2))
	heap.Pop(scheduler)
	if scheduler.Heap[0].nextTriggerTime.After(scheduler.Heap[1].nextTriggerTime) {
		t.Errorf("Heap goes wrong")
	}
	t.Log(scheduler.Heap[0].nextTriggerTime)
	t.Log(scheduler.Heap[1].nextTriggerTime)
}
