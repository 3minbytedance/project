package common

import (
	"douyin/dal/model"
	"douyin/dal/mysql"
	"github.com/bits-and-blooms/bloom/v3"
	"go.uber.org/zap"
	"log"
)

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

func LoadUsernamesToBloomFilter() {
	var usernames []string
	err := mysql.DB.Model(&model.User{}).Pluck("name", &usernames).Error
	if err != nil {
		log.Fatal("Failed to retrieve usernames from database:", err)
	}

	for _, username := range usernames {
		AddToBloom(username)
	}

	zap.L().Info("Loaded %d usernames to the bloom filter.\n", zap.Int("size", len(usernames)))
}
