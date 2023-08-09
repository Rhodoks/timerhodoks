package scheduler

import (
	"container/heap"
	"sync"
	"time"
	"timerhodoks/pkg/job"
	"timerhodoks/pkg/worker"

	"github.com/gorhill/cronexpr"
)

// 堆节点，维护了每个节点所在的下标，通过哈希表从ID索引节点。
// 可以做到O(log)的插入删除修改和O(1)的查询
type SchedulerNode struct {
	HeapId          int    //在堆中的下标
	id              uint64 //任务Id
	cron            *cronexpr.Expression
	executorInfo    string
	executorType    string
	retryNum        int
	finishTime      time.Time //一个时间点，其之前的所有任务均已完成
	lastTriggerTime time.Time //上一次触发的时间
	nextTriggerTime time.Time //下一次触发的时间
	mutex           sync.Mutex
}

func int64Max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func NewSchedulerNode(e *job.JobEntry) *SchedulerNode {
	cron := cronexpr.MustParse(e.CronLine)

	StartTime := time.Unix(int64Max(e.CreateTime, int64Max(e.LastFailureTime, e.LastSuccessTime)), 0) //重新插入的任务，其开始时间是（创建时间，上次成功时间，上次失败时间）中的最大值
	return &SchedulerNode{
		HeapId:          -1,
		id:              e.Id,
		cron:            cron,
		executorType:    e.ExecutorType,
		executorInfo:    e.ExecutorInfo,
		retryNum:        e.RetryNum,
		finishTime:      StartTime,
		lastTriggerTime: StartTime,
		nextTriggerTime: cron.Next(StartTime),
		mutex:           sync.Mutex{},
	}
}

type Scheduler struct {
	Heap              []*SchedulerNode
	IdMap             map[uint64]*SchedulerNode
	Mutex             sync.RWMutex
	ProposeCommitTime func(uint64, int64, bool) //传入函数指针避免循环依赖，更新任务完成时间。
	workers           []worker.Worker
	worker_ptr        int
}

func NewScheduler(workers []worker.Worker) *Scheduler {
	h := &Scheduler{
		Heap:    []*SchedulerNode{},
		IdMap:   make(map[uint64]*SchedulerNode),
		workers: workers,
	}
	heap.Init(h)
	return h
}

// 定义堆的比较函数
func (s *Scheduler) Less(i, j int) bool {
	t1 := s.Heap[i].nextTriggerTime
	t2 := s.Heap[j].nextTriggerTime
	if t1.Equal(t2) {
		return s.Heap[i].id < s.Heap[j].id
	}
	return t1.Before(t2)
}

// 获取堆大小
func (s *Scheduler) Len() int {
	return len(s.Heap)
}

// 定义堆的交换函数，维护任务结构体中的堆下标
func (s *Scheduler) Swap(i, j int) {
	s.Heap[i], s.Heap[j] = s.Heap[j], s.Heap[i]
	s.Heap[i].HeapId, s.Heap[j].HeapId = s.Heap[j].HeapId, s.Heap[i].HeapId
}

// 定义堆的Push操作
func (s *Scheduler) Push(x interface{}) {
	n := len(s.Heap)
	node := x.(*SchedulerNode)
	node.HeapId = n
	s.Heap = append(s.Heap, node)
	s.IdMap[node.id] = node
}

// 定义堆的Pop操作
func (s *Scheduler) Pop() interface{} {
	old := s.Heap
	n := len(old)
	node := old[n-1]
	node.HeapId = -1 // for safety
	s.Heap = old[0 : n-1]
	delete(s.IdMap, node.id)
	return node
}

// 移除给定id的任务
func (s *Scheduler) Remove(id uint64) {
	node := s.IdMap[id]
	if node.HeapId >= 0 {
		heap.Remove(s, node.HeapId)
	}
	delete(s.IdMap, id)
}

// 插入给定任务
func (s *Scheduler) Insert(job *job.JobEntry) {
	node := NewSchedulerNode(job)
	s.IdMap[job.Id] = node
	if (node.nextTriggerTime == time.Time{}) { //没有下一次触发的任务不必放入调度器
		return
	}
	heap.Push(s, node)
}

// 执行任务，并且等待直至该任务之前的触发均已处理完毕，提交并且释放worker
func (s *Scheduler) ExecuteAndUpdate(worker worker.Worker, lastTriggerTime time.Time, curTriggerTime time.Time, node *SchedulerNode) {
	err := worker.Execute(node.id, curTriggerTime, node.retryNum, node.executorType, node.executorInfo)
	for {
		node.mutex.Lock()
		if node.finishTime.Equal(lastTriggerTime) || node.finishTime.After(lastTriggerTime) {
			node.finishTime = curTriggerTime
			node.mutex.Unlock()
			break
			//完成的任务会自动在此等待，只有提交时间推进到上次提交时间才会写入
			//删除或更新时，node直接被丢弃
		}
		node.mutex.Unlock()
		time.Sleep(100 * time.Microsecond)
	}
	if err == nil {
		s.ProposeCommitTime(node.id, curTriggerTime.Unix(), true)
	} else {
		s.ProposeCommitTime(node.id, curTriggerTime.Unix(), false)
	}
	worker.Release()
}

// 检查堆顶的任务和当前时间，触发到时间的任务
func (s *Scheduler) Schedule() int {
	cnt := 0
	now := time.Now()
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	for len(s.Heap) > 0 && s.Heap[0].nextTriggerTime.Before(now) {
		node := s.Heap[0]
		lastTriggerTime := node.lastTriggerTime
		curTriggerTime := node.nextTriggerTime
		node.lastTriggerTime = node.nextTriggerTime
		node.nextTriggerTime = node.cron.Next(node.nextTriggerTime)
		if (node.nextTriggerTime == time.Time{}) { //没有下一次触发时间了
			node.HeapId = -1
			s.Pop() //从堆中删除该任务
		} else {
			heap.Fix(s, 0) //更新任务
		}
		for s.worker_ptr++; ; s.worker_ptr++ {
			id := s.worker_ptr % len(s.workers)
			if s.workers[id].TryAcquire() {
				go s.ExecuteAndUpdate(s.workers[id], lastTriggerTime, curTriggerTime, node)
				break
			}
		}
	}
	return cnt
}

func (s *Scheduler) Running() {
	for {
		s.Schedule()
	}
}
