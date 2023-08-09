package job

const HASH_BIT_LEN = 16
const HASH_BUC_NUM = (1 << HASH_BIT_LEN)
const HASH_MASK = HASH_BUC_NUM - 1
const PRIME uint64 = 19260817

// 对任务Id进行哈希
// 哈希函数是x * p mod len，由于p是质数，所以哈希取值是回环的
func Hash(x uint64) uint64 {
	return (x * PRIME) & HASH_MASK
}
