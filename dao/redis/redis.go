package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project/config"
	"time"
)

var Ctx = context.Background()

// RdbComment Comment模块Rdb
var RdbComment *redis.Client

// RdbExpireTime key的过期时间
var RdbExpireTime time.Duration
var (
	RDB               *redis.Client
	UserFavoriteRDB   *redis.Client
	VideoFavoritedRDB *redis.Client
)

func Init(appConfig *config.AppConfig) (err error) {
	var conf *config.RedisConfig
	if appConfig.Mode == config.LocalMode {
		conf = appConfig.Local.RedisConfig
	} else {
		conf = appConfig.Remote.RedisConfig
	}
	// 获取conf中的过期时间, 单位为s
	RdbExpireTime = time.Duration(conf.ExpireTime) * time.Second

	RdbComment = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password:     conf.Password,  // 密码
		DB:           conf.CommentDB, // 数据库
		PoolSize:     conf.PoolSize,  // 连接池大小
		MinIdleConns: conf.MinIdleConns,
	})
	if err = RdbComment.Ping(Ctx).Err(); err != nil {
		return nil
	}

	UserFavoriteRDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password:     conf.Password,       // 密码
		DB:           conf.UerFavoriteRDB, // 数据库
		PoolSize:     conf.PoolSize,       // 连接池大小
		MinIdleConns: conf.MinIdleConns,
	})

	VideoFavoritedRDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password:     conf.Password,          // 密码
		DB:           conf.VideoFavoritedRDB, // 数据库
		PoolSize:     conf.PoolSize,          // 连接池大小
		MinIdleConns: conf.MinIdleConns,
	})

	_, err = RDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = RDB.Close()
}
