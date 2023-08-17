package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
	"strconv"
)

// GetFollowCountById 根据userId查找关注数
func GetFollowCountById(userId uint) (int, error) {
	key := fmt.Sprintf("%d_%s", userId, FollowList)
	size, err := Rdb.SCard(Ctx, key).Result()
	return int(size), err

}

// 根据userId查找粉丝数
func GetFollowerCountById(userId uint) (int, error) {
	key := fmt.Sprintf("%d_%s", userId, FollowerList)
	size, err := Rdb.SCard(Ctx, key).Result()
	return int(size), err
}

// 根据userId查找关注list
func GetFollowListById(userId uint) ([]uint, error) {
	key := fmt.Sprintf("%d_%s", userId, FollowList)
	list, err := Rdb.SMembers(Ctx, key).Result()
	var result []uint
	for _, i := range list {
		k, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(k))
	}
	return result, err
}

// 根据userId查找粉丝list
func GetFollowerListById(userId uint) ([]uint, error) {
	key := fmt.Sprintf("%d_%s", userId, FollowerList)
	list, err := Rdb.SMembers(Ctx, key).Result()
	var result []uint
	for _, i := range list {
		k, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(k))
	}
	return result, err
}

// 根据userId查找好友list
func GetFriendListById(userId uint) ([]uint, error) {
	key1 := fmt.Sprintf("%d_%s", userId, FollowerList)
	key2 := fmt.Sprintf("%d_%s", userId, FollowList)
	friend, err := Rdb.SUnion(Ctx, key2, key1).Result()
	var result []uint
	for _, i := range friend {
		k, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(k))
	}
	return result, err
}

// 设置关注列表
func SetFollowListByUserId(userid uint, id []uint) error {
	key := fmt.Sprintf("%d_%s", userid, FollowList)
	// 转换为[]interface{}
	b := make([]interface{}, len(id))
	for i := range id {
		b[i] = id[i]
	}
	err := Rdb.SAdd(Ctx, key, b...).Err()
	return err
}

// 设置粉丝列表
func SetFollowerListByUserId(userid uint, id []uint) error {
	key := fmt.Sprintf("%d_%s", userid, FollowerList)
	// 转换为[]interface{}
	b := make([]interface{}, len(id))
	for i := range id {
		b[i] = id[i]
	}
	err := Rdb.SAdd(Ctx, key, b...).Err()
	return err
}

// 给Id对应的关注set加上 id
func IncreaseFollowCountByUserId(userId uint, id uint) error {
	key := fmt.Sprintf("%d_%s", userId, FollowList)
	err := Rdb.SAdd(Ctx, key, id).Err()
	return err
}

// DecreaseFollowCountByUserId userId 关注列表取关 followId
func DecreaseFollowCountByUserId(userId uint, followId uint) error {
	key := fmt.Sprintf("%d_%s", userId, FollowList)
	err := Rdb.SRem(Ctx, key, followId).Err()
	return err
}

// IncreaseFollowerCountByUserId 给userId粉丝列表加上 followid
func IncreaseFollowerCountByUserId(userId uint, followId uint) error {
	key := fmt.Sprintf("%d_%s", userId, FollowerList)
	err := Rdb.SAdd(Ctx, key, followId).Err()
	return err
}

// 给userId对应的粉丝列表减去id
func DecreaseFollowerCountByUserId(userId uint, id uint) error {
	key := fmt.Sprintf("%d_%s", userId, FollowerList)
	err := Rdb.SRem(Ctx, key, id).Err()
	return err
}

// id是不是userid 的好友
func IsInMyFollowList(userId uint, id uint) bool {
	key := fmt.Sprintf("%d_%s", userId, FollowerList)
	found, _ := Rdb.SIsMember(Ctx, key, id).Result()
	return found
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