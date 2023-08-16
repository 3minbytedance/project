package common

import "github.com/bits-and-blooms/bloom/v3"

var bloomFilter *bloom.BloomFilter

func InitBloomFilter() {
	// 初始化布隆过滤器
	bloomFilter = bloom.NewWithEstimates(10000000, 0.05) // 假设预期元素数量为 10000000，误判率为 0.05
}

func AddToBloom(data string) {
	bloomFilter.Add([]byte(data))
}

func TestBloom(data string) bool {
	return bloomFilter.Test([]byte(data))
}
