package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
)

// GetFollowCountById 根据userId查找关注数
func GetFollowCountById(userId uint) (int, error) {
	key := fmt.Sprintf("%s_%d", FollowList, userId)
	size, err := Rdb.SCard(Ctx, key).Result()
	return int(size), err
}

// GetFollowerCountById 根据userId查找粉丝数
func GetFollowerCountById(userId uint) (int, error) {
	key := fmt.Sprintf("%s_%d", FollowerList, userId)
	size, err := Rdb.SCard(Ctx, key).Result()
	return int(size), err
}

// GetFollowListById 根据userId查找关注list
func GetFollowListById(userId uint) ([]uint, error) {
	key := fmt.Sprintf("%s_%d", FollowList, userId)
	list, err := Rdb.SMembers(Ctx, key).Result()
	result := make([]uint, 0)
	if err != nil {
		return result, err
	}
	for _, i := range list {
		k, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(k))
	}
	return result, err
}

// GetFollowerListById 根据userId查找粉丝list
func GetFollowerListById(userId uint) ([]uint, error) {
	key := fmt.Sprintf("%s_%d", FollowerList, userId)
	list, err := Rdb.SMembers(Ctx, key).Result()
	result := make([]uint, 0)
	if err != nil {
		return result, err
	}
	for _, i := range list {
		k, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(k))
	}
	return result, err
}

// GetFriendListById 根据userId查找好友list
func GetFriendListById(userId uint) ([]uint, error) {
	key1 := fmt.Sprintf("%s_%d", FollowerList, userId)
	key2 := fmt.Sprintf("%s_%d", FollowList, userId)
	friend, err := Rdb.SInter(Ctx, key2, key1).Result()
	result := make([]uint, 0, len(friend))
	if err != nil {
		return result, err
	}
	for _, i := range friend {
		k, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(k))
	}
	return result, err
}

// SetFollowListByUserId 设置关注列表
func SetFollowListByUserId(userId uint, ids []uint) error {
	key := fmt.Sprintf("%s_%d", FollowList, userId)
	pipe := Rdb.Pipeline()
	for _, value := range ids {
		err := pipe.SAdd(Ctx, key, value).Err()
		if err != nil {
			return err
		}
	}
	zap.L().Info("Follow_LIST", zap.Any("List", ids))
	_, err := pipe.Exec(Ctx)
	return err
}

// SetFollowerListByUserId 设置粉丝列表
func SetFollowerListByUserId(userId uint, ids []uint) error {
	key := fmt.Sprintf("%s_%d", FollowerList, userId)
	pipe := Rdb.Pipeline()
	// 转换为[]interface{}
	for _, value := range ids {
		err := pipe.SAdd(Ctx, key, value).Err()
		if err != nil {
			return err
		}
	}
	zap.L().Info("Follower_LIST", zap.Any("List", ids))
	_, err := pipe.Exec(Ctx)
	return err
}

// IncreaseFollowCountByUserId 给Id对应的关注set加上 id
func IncreaseFollowCountByUserId(userId uint, id uint) error {
	key := fmt.Sprintf("%s_%d", FollowList, userId)
	err := Rdb.SAdd(Ctx, key, id).Err()
	return err
}

// DecreaseFollowCountByUserId userId 关注列表取关 followId
func DecreaseFollowCountByUserId(userId uint, followId uint) error {
	key := fmt.Sprintf("%s_%d", FollowList, userId)
	err := Rdb.SRem(Ctx, key, followId).Err()
	return err
}

// IncreaseFollowerCountByUserId 给userId粉丝列表加上 followid
func IncreaseFollowerCountByUserId(userId uint, followId uint) error {
	key := fmt.Sprintf("%s_%d", FollowerList, userId)
	err := Rdb.SAdd(Ctx, key, followId).Err()
	return err
}

// DecreaseFollowerCountByUserId 给userId对应的粉丝列表减去id
func DecreaseFollowerCountByUserId(userId uint, id uint) error {
	key := fmt.Sprintf("%s_%d", FollowerList, userId)
	err := Rdb.SRem(Ctx, key, id).Err()
	return err
}

// IsInMyFollowList userid是否关注了id
func IsInMyFollowList(userId uint, id uint) (bool, error) {
	key := fmt.Sprintf("%s_%d", FollowList, userId)
	found, err := Rdb.SIsMember(Ctx, key, id).Result()
	return found, err
}

//-------------------------------丢弃
// 给Id对应的关注数加一
//func IncreaseFollowCountByUserId(userId uint) error {
//	key := UserKey + fmt.Sprintf("%d", userId)
//	_, err := Rdb.HIncrBy(Ctx, key, FollowCountField, 1).Result()
//	return err
//}
//
//// 给Id对应的关注数减一
//func DecreaseFollowCountByUserId(userId uint) error {
//	key := UserKey + fmt.Sprintf("%d", userId)
//	_, err := Rdb.HIncrBy(Ctx, key, FollowCountField, -1).Result()
//	return err
//}
//
//// 给videoId对应的粉丝数加一
//func IncreaseFollowerCountByUserId(userId uint) error {
//	key := UserKey + fmt.Sprintf("%d", userId)
//	_, err := Rdb.HIncrBy(Ctx, key, FollowerCountField, 1).Result()
//	return err
//}
//
//// 给videoId对应的粉丝数减一
//func DecreaseFollowerCountByUserId(userId uint) error {
//	key := UserKey + fmt.Sprintf("%d", userId)
//	_, err := Rdb.HIncrBy(Ctx, key, FollowerCountField, -1).Result()
//	return err
//}
//
//// 设置关注数
//func SetFollowCountByUserId(userid uint, id uint) error {
//	key := UserKey + fmt.Sprintf("%d", userid)
//	err := Rdb.HSet(Ctx, key, FollowCountField, count).Err()
//	return err
//	//key := fmt.Sprintf("%d_%s", userid, FollowList)
//	//err := Rdb.SAdd(Ctx, key, id).Err()
//    //return err
//}
//
//// 设置粉丝数
//func SetFollowerCountByUserId(userid uint, count int64) error {
//	key := UserKey + fmt.Sprintf("%d", userid)
//	err := Rdb.HSet(Ctx, key, FollowerCountField, count).Err()
//	return err
//}
