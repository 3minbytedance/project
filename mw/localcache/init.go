package localcache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"log"
	"time"
)

const (
	WorkCount int32 = iota
	User
	FavoriteVideo
)

func Init(rpcName int32) *bigcache.BigCache {

	var lifeWindow time.Duration
	switch rpcName {
	case WorkCount:
		lifeWindow = 10 * time.Second
	case User:
		lifeWindow = 30 * time.Minute
	case FavoriteVideo:
		lifeWindow = 30 * time.Minute
	default:
		lifeWindow = 10 * time.Second
	}

	config := bigcache.Config{
		// 缓存条目数量
		Shards: 1024,

		// 单个缓存条目过期时间
		LifeWindow: lifeWindow,

		// 最大缓存值大小
		MaxEntrySize: 50,

		// 最大缓存总值大小
		MaxEntriesInWindow: 1000 * 10 * 60,
	}
	cache, initErr := bigcache.New(context.Background(), config)
	if initErr != nil {
		log.Fatal(initErr)
	}
	return cache
}
