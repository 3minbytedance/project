package test

import (
	"douyin/config"
	"douyin/dal/mysql"
	"douyin/logger"
	"douyin/mw/redis"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"strconv"
	"strings"
	"testing"
	"time"
)

func init() {
	// 加载配置
	if err := config.Init(); err != nil {
		zap.L().Error("Load config failed, err:%v\n", zap.Error(err))
		return
	}
	// 加载日志
	if err := logger.Init(config.Conf.LogConfig, config.Conf.Mode); err != nil {
		zap.L().Error("Init logger failed, err:%v\n", zap.Error(err))
		return
	}

	// 初始化数据库: mysql
	if err := mysql.Init(config.Conf); err != nil {
		zap.L().Error("Init mysql failed, err:%v\n", zap.Error(err))
		return
	}

	// 初始化中间件: redis + kafka
	if err := redis.Init(config.Conf); err != nil {
		zap.L().Error("Init middleware failed, err:%v\n", zap.Error(err))
		return
	}
}

func Benchmark_Add(b *testing.B) {
	for i := 0; i < 1000; i++ {
		GetName(1698866670832979968)
	}
}

var g singleflight.Group

// GetName 根据userId获取用户名
func GetName(userId uint) (string, error) {

	baseSlice := []string{redis.UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, redis.Delimiter)
	// 查缓存
	if name, err := redis.GetNameByUserId(userId); err == nil {
		return name, nil
	}
	v, err, _ := g.Do(key, func() (interface{}, error) {
		userModel, exist, _ := mysql.FindUserByUserID(userId)
		if !exist {
			return "", nil
		}
		// 将用户名写入redis
		err := redis.SetNameByUserId(userId, userModel.Name)
		if err != nil {
			zap.L().Error("将用户名写入redis失败：", zap.Error(err))
			return "", err
		}
		return userModel.Name, nil
	})
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

// GetName 根据userId获取用户名
func GetName2(userId uint) (string, bool) {
	// 从redis中获取用户名
	// 1. 缓存中有数据, 直接返回
	if name, err := redis.GetNameByUserId(userId); err == nil {
		return name, true
	}
	//缓存不存在，尝试从数据库中取
	if redis.AcquireUserLock(userId, redis.NameField) {
		defer redis.ReleaseUserLock(userId, redis.NameField)
		// 2. 缓存中没有数据，从数据库中获取
		userModel, exist, _ := mysql.FindUserByUserID(userId)
		if !exist {
			return "", false
		}
		// 将用户名写入redis
		err := redis.SetNameByUserId(userId, userModel.Name)
		if err != nil {
			zap.L().Error("将用户名写入redis失败：", zap.Error(err))
		}
		return userModel.Name, true
	}
	// 重试
	time.Sleep(redis.RetryTime)
	return GetName2(userId)
}
