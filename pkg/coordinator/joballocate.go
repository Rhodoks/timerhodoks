package coordinator

import (
	"timerhodoks/pkg/job"

	"github.com/bits-and-blooms/bitset"
)

var FULL_JOB_ALLOCATION = bitset.New(job.HASH_BUC_NUM)

func JobAllocationInit() {
	for i := uint(0); i < job.HASH_BUC_NUM; i++ {
		FULL_JOB_ALLOCATION.Set(i)
	}
}

func NewJobAllocation() *bitset.BitSet {
	return bitset.New(job.HASH_BUC_NUM)
}
