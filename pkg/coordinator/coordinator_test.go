package coordinator

import (
	"fmt"
	"testing"
	"timerhodoks/pkg/job"
)

func TestHttpExecutor(t *testing.T) {
	x := NewJobAllocation()
	y := NewJobAllocation()
	x.Set(1)
	x.Set(2)
	y.Set(2)
	buffer := make([]uint, job.HASH_BUC_NUM)
	j, buffer := x.NextSetMany(0, buffer)
	fmt.Println(j, buffer)
}
