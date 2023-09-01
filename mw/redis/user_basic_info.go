package redis

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetWorkCountByUserId(userId uint) (int64, error) {
	baseSlice := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)

	count, err := Rdb.HGet(Ctx, key, WorkCountField).Result()
	commentCount, _ := strconv.ParseInt(count, 10, 64)
	return commentCount, err
}

func GetNameByUserId(userId uint) (string, error) {
	baseSlice := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	name, err := Rdb.HGet(Ctx, key, NameField).Result()
	return name, err
}

func IncrementWorkCountByUserId(userId uint) error {
	baseSlice := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, WorkCountField, 1).Result()
	return err
}

func SetWorkCountByUserId(userId uint, workCount int64) error {
	baseSlice := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	err := Rdb.HSet(Ctx, key, WorkCountField, workCount).Err()
	randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}

func SetNameByUserId(userId uint, name string) error {
	baseSlice := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	err := Rdb.HSet(Ctx, key, NameField, name).Err()
	randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}
