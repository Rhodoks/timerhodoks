package coordinator

import (
	"context"
	"math/rand"
	"sync"
	"time"
	pb "timerhodoks/pkg/coordinatorpb"
	"timerhodoks/pkg/job"
	"timerhodoks/pkg/raftnode"
	"timerhodoks/pkg/scheduler"

	"github.com/bits-and-blooms/bitset"
	"go.etcd.io/raft/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 协调器
// 提供grpc服务，告知Leader自己的任务分配情况
// 对于Leader节点，负责定期检查和更新全局的任务分配
type Coordinator struct {
	Scheduler          *scheduler.Scheduler
	UpdateFlag         uint32 //标志本地任务分配是否被修改。RaftStore会定期检查。
	Mutex              sync.RWMutex
	ApplyJobAllocation *bitset.BitSet //被RaftStore接受的任务分配
	PlanJobAllocation  *bitset.BitSet //最新被分配的任务
	peerGrpcServer     []string       //兄弟节点的grpc服务器
	raftNode           *raftnode.RaftNode
	pb.UnimplementedCoordinatorServer
}

func NewCoordinator(scheduler *scheduler.Scheduler, peerGrpcServer []string, raftNode *raftnode.RaftNode, grpcPort int) *Coordinator {
	res := &Coordinator{
		Scheduler:          scheduler,
		UpdateFlag:         0,
		Mutex:              sync.RWMutex{},
		ApplyJobAllocation: NewJobAllocation(),
		PlanJobAllocation:  NewJobAllocation(),
		peerGrpcServer:     peerGrpcServer,
		raftNode:           raftNode,
	}
	go res.Run(grpcPort)
	return res
}

const COORDINATE_INTERVAL_MS = 1000
const GRPC_TIMEOUT_MS = 1000
const GRPC_TRY_NUM = 3

// 通过grpc获得兄弟节点的任务分配和负载情况
func getCoordinatingInfo(peer string, allocation **bitset.BitSet, load *int, id int, chanF chan int) error {
	for i := 1; i <= GRPC_TRY_NUM; i++ {
		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()
		if err != nil {
			if i == GRPC_TRY_NUM { //到达重试次数上限
				chanF <- -1
				return err
			}
			continue
		}
		c := pb.NewCoordinatorClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT_MS*time.Millisecond)
		defer cancel()
		res, err := c.GetCoordinatingInfo(ctx, &pb.EmptyMessage{})
		if err != nil {
			if i == GRPC_TRY_NUM { //到达重试次数上限
				chanF <- -1
				return err
			}
			continue
		}
		uint64slice, err := byteToUint64(res.JobAllocation)
		if err != nil {
			if i == GRPC_TRY_NUM { //到达重试次数上限
				chanF <- -1
				return err
			}
			continue
		}
		*allocation = bitset.From(uint64slice)
		*load = int(res.Load)
		chanF <- id
		return nil
	}
	panic("Unreachable!")
}

// 更新兄弟节点的任务分配
func updateAllocation(peer string, allocation bitset.BitSet, id int, chanF chan int) error {
	for i := 1; i <= GRPC_TRY_NUM; i++ {
		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()
		if err != nil {
			if i == GRPC_TRY_NUM { //到达重试次数上限
				chanF <- -1
				return err
			}
			continue
		}
		c := pb.NewCoordinatorClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT_MS*time.Millisecond)
		defer cancel()
		buf := uint64ToByte(allocation.Bytes())
		_, err = c.UpdateAllocation(ctx, &pb.AllocationRequest{JobAllocation: buf})
		if err != nil {
			if i == GRPC_TRY_NUM { //到达重试次数上限
				chanF <- -1
				return err
			}
			continue
		}
		chanF <- id
		return nil
	}
	panic("Unreachable!")
}

func (c *Coordinator) coordinate() {
	//通过grpc获取兄弟节点的任务分配
	num := len(c.peerGrpcServer)
	allocations := make([]*bitset.BitSet, num)
	loads := make([]int, num)
	chanF := make(chan int, num)

	idleNodes := make([]int, 0)   //未分配任务的节点
	activeNodes := make([]int, 0) //存活的节点

	for id, peer := range c.peerGrpcServer {
		go getCoordinatingInfo(peer, &allocations[id], &loads[id], id, chanF)
	}
	for i := 0; i < num; i++ {
		id := <-chanF
		if id >= 0 {
			activeNodes = append(activeNodes, id)
			if allocations[id].Count() == 0 {
				idleNodes = append(idleNodes, id)
			}
		}
	} // 等待任务完成

	if len(idleNodes) > 0 { // 如果有空闲节点的话
		idVec := make([]uint, 0)
		for i := 0; i < job.HASH_BUC_NUM; i++ {
			idVec = append(idVec, uint(i))
		}
		rand.Shuffle(len(idVec), func(i, j int) {
			idVec[i], idVec[j] = idVec[j], idVec[i]
		}) //随机打乱

		for _, idleId := range idleNodes { //按照比例随机选取一些桶分配给空闲节点
			aver := job.HASH_BUC_NUM / len(activeNodes) //平均每个节点的任务数量
			for j := idleId * aver; j < aver+idleId*aver; j++ {
				for _, id := range activeNodes {
					allocations[id].Clear(idVec[j])
				}
				allocations[idleId].Set(idVec[j])
			}
		}
	}

	occupy := make([][]uint, job.HASH_BUC_NUM) //记录每个哈希桶被分配给了哪些节点
	buffer := make([]uint, job.HASH_BUC_NUM)
	for nodeId, alloc := range allocations {
		if alloc == nil {
			continue
		}
		_, buffer = alloc.NextSetMany(0, buffer)
		for _, occupyId := range buffer {
			occupy[occupyId] = append(occupy[occupyId], uint(nodeId))
		}
	}
	for i := uint(0); i < job.HASH_BUC_NUM; i++ {
		if len(occupy[i]) == 0 { //未被分配的哈希桶，则
			for {
				x := uint(rand.Int31()) % uint(num)
				if allocations[x] != nil { // x存活
					occupy[i] = append(occupy[i], x)
					allocations[x].Set(i)
					break
				}
			}
		} else if len(occupy[i]) > 1 { //被多分配的哈希桶，只保留一个
			arr := occupy[i]
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			}) //随机打乱
			for len(arr) > 1 {
				x := arr[len(arr)-1]
				allocations[x].Clear(x)
				arr = arr[:len(arr)-1]
			} //只保留第一个节点作为分配者
		}
	}

	for nodeId, peer := range c.peerGrpcServer {
		if allocations[nodeId] == nil {
			num--
			continue
		}
		updateAllocation(peer, *allocations[nodeId], nodeId, chanF)
	}
	for i := 0; i < num; i++ {
		<-chanF
	} // 等待任务完成
}

func (c *Coordinator) Run(grpcPort int) {
	go c.ListenAndServe(grpcPort)
	for {
		time.Sleep(COORDINATE_INTERVAL_MS * time.Millisecond) //每隔COORDINATE_INTERVAL_MS毫秒进行一次调度
		if c.raftNode.GetState() == raft.StateLeader {
			c.coordinate()
		}
	}
}

// 查询第id个哈希桶是否分配给了自己
func (c *Coordinator) Query(id uint64) bool {
	hash := job.Hash(id)
	return c.ApplyJobAllocation.Test(uint(hash))
}
