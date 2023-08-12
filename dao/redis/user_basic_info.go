package redis

import (
	"fmt"
	"strconv"
)

const (
	WorkCount = "work_count" //作品数
	Name      = "name"       //用户名
	Avator    = "avator"     //头像url
	BgImg     = "bg_img"     //背景大图url
)

func GetWorkCountByUserId(userId uint) (int64, error) {
	key := User + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, WorkCount).Result()
	commentCount, _ := strconv.ParseInt(count, 10, 64)
	return commentCount, err
}

func GetNameByUserId(userId uint) (string, error) {
	key := User + fmt.Sprintf("%d", userId)
	name, err := Rdb.HGet(Ctx, key, Name).Result()
	return name, err
}

func IncrementWorkCountByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, WorkCount, 1).Result()
	return err
}

func SetWorkCountByUserId(userId uint, workCount int64) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, WorkCount, workCount).Err()
	return err
}

func SetNameByUserId(userId uint, name string) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, Name, name).Err()
	return err
}
