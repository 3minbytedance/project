package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project/config"
)

var RedisCtx = context.Background()

// RdbVCId key: videoId value: commentId(zset, score: createTime)
var RdbVCId *redis.Client

// RdbCVId  key: commentId value: videoId
var RdbCVId *redis.Client

// RdbCIdComment key: commentId value: comment
var RdbCIdComment *redis.Client

func Init(appConfig *config.AppConfig) (err error) {
	var conf *config.RedisConfig
	if appConfig.Mode == config.LocalMode {
		conf = appConfig.Local.RedisConfig
	} else {
		conf = appConfig.Remote.RedisConfig
	}

	RdbVCId = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password:     conf.Password, // 密码
		DB:           conf.VCIdDB,   // 数据库
		PoolSize:     conf.PoolSize, // 连接池大小
		MinIdleConns: conf.MinIdleConns,
	})
	// 判断是否连接成功, 不成功则返回错误
	_, err = RdbVCId.Ping(RedisCtx).Result()
	if err != nil {
		return err
	}
	RdbCVId = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password:     conf.Password, // 密码
		DB:           conf.CVIdDB,   // 数据库
		PoolSize:     conf.PoolSize, // 连接池大小
		MinIdleConns: conf.MinIdleConns,
	})
	// 判断是否连接成功, 不成功则返回错误
	_, err = RdbCVId.Ping(RedisCtx).Result()
	if err != nil {
		return err
	}
	RdbCIdComment = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password:     conf.Password,     // 密码
		DB:           conf.CIdCommentDB, // 数据库
		PoolSize:     conf.PoolSize,     // 连接池大小
		MinIdleConns: conf.MinIdleConns,
	})
	// 判断是否连接成功, 不成功则返回错误
	_, err = RdbCIdComment.Ping(RedisCtx).Result()
	if err != nil {
		return err
	}
	return nil
}
