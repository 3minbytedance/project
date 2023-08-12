package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
	"strconv"
)

const (
	KeyFollowCount   = "follow_count"
	KeyFollowerCount = "follower_count  "
)

// 判断是否存在此建
func Is_Exist(userId uint, s string) bool {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.Exists(Ctx, key, s).Err()
	if err != nil {
		return false
	}
	return true
}

// 根据userId查找关注数
func GetFollowCountById(userId uint) (int, error) {
	key := User + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, KeyFollowCount).Result()
	commentCount, _ := strconv.Atoi(count)
	return commentCount, err
}

// 根据userId查找粉丝数
func GetFollowerCountById(userId uint) (int, error) {
	key := User + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, KeyFollowerCount).Result()
	commentCount, _ := strconv.Atoi(count)
	return commentCount, err
}

// 给Id对应的关注数加一
func IncreaseFollowCountByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, KeyFollowCount, 1).Result()
	return err
}

// 给Id对应的关注数减一
func DecreaseFollowCountByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, KeyFollowCount, -1).Result()
	return err
}

// 给videoId对应的粉丝数加一
func IncreaseFollowerCountByUserId(userId uint) error {
	key := Video + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, KeyFollowerCount, 1).Result()
	return err
}

// 给videoId对应的粉丝数减一
func DecreaseFollowerCountByUserId(userId uint) error {
	key := Video + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, KeyFollowerCount, -1).Result()
	return err
}

// 设置关注数
func SetFollowCountByUserId(userid uint, count int64) error {
	key := Video + fmt.Sprintf("%d", userid)
	err := Rdb.HSet(Ctx, key, KeyFollowCount, count).Err()
	return err
}

// 设置粉丝数
func SetFollowerCountByUserId(userid uint, count int64) error {
	key := Video + fmt.Sprintf("%d", userid)
	err := Rdb.HSet(Ctx, key, KeyFollowerCount, count).Err()
	return err
}
