package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"project/config"
)

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

	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password:     conf.Password, // 密码
		DB:           conf.DB,       // 数据库
		PoolSize:     conf.PoolSize, // 连接池大小
		MinIdleConns: conf.MinIdleConns,
	})

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
