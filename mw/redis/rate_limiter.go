package redis

import "github.com/redis/go-redis/v9"

func Ping() error {
	if _, err := Rdb.Ping(Ctx).Result(); err != nil {

		return err
	}
	return nil
}

func RunScript(script string, keys []string, args ...interface{}) (interface{}, error) {
	val, err := redis.NewScript(script).Run(Ctx, Rdb, keys, args).Result()
	if err != nil && err != redis.Nil {

	}
	return val, nil
}

func ClearAll() error {
	if _, err := Rdb.FlushAll(Ctx).Result(); err != nil {

		return err
	}
	return nil
}
