package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
	"strconv"
)

const ()

// 根据userId查找关注数
func GetFollowCountById(userId uint) (int, error) {
	key := UserKey + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, FollowCountField).Result()
	commentCount, _ := strconv.Atoi(count)
	return commentCount, err
}

// 根据userId查找粉丝数
func GetFollowerCountById(userId uint) (int, error) {
	key := UserKey + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, FollowerCountField).Result()
	commentCount, _ := strconv.Atoi(count)
	return commentCount, err
}

// 给Id对应的关注数加一
func IncreaseFollowCountByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, FollowCountField, 1).Result()
	return err
}

// 给Id对应的关注数减一
func DecreaseFollowCountByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, FollowCountField, -1).Result()
	return err
}

// 给videoId对应的粉丝数加一
func IncreaseFollowerCountByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, FollowerCountField, 1).Result()
	return err
}

// 给videoId对应的粉丝数减一
func DecreaseFollowerCountByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, FollowerCountField, -1).Result()
	return err
}

// 设置关注数
func SetFollowCountByUserId(userid uint, count int64) error {
	key := UserKey + fmt.Sprintf("%d", userid)
	err := Rdb.HSet(Ctx, key, FollowCountField, count).Err()
	return err
}

// 设置粉丝数
func SetFollowerCountByUserId(userid uint, count int64) error {
	key := UserKey + fmt.Sprintf("%d", userid)
	err := Rdb.HSet(Ctx, key, FollowerCountField, count).Err()
	return err
}
