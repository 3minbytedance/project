package redis

import (
	"fmt"
	"strconv"
)

const ()

func GetWorkCountByUserId(userId uint) (int64, error) {
	key := UserKey + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, WorkCountField).Result()
	commentCount, _ := strconv.ParseInt(count, 10, 64)
	return commentCount, err
}

func GetNameByUserId(userId uint) (string, error) {
	key := UserKey + fmt.Sprintf("%d", userId)
	name, err := Rdb.HGet(Ctx, key, NameField).Result()
	return name, err
}

func IncrementWorkCountByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, WorkCountField, 1).Result()
	return err
}

func SetWorkCountByUserId(userId uint, workCount int64) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, WorkCountField, workCount).Err()
	return err
}

func SetNameByUserId(userId uint, name string) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, NameField, name).Err()
	return err
}