package redis

import (
	_ "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GetFollowCountById 根据userId查找关注数
func GetFollowCountById(userId uint) (int, error) {
	baseSlice := []string{FollowList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	size, err := Rdb.SCard(Ctx, key).Result()
	return int(size), err
}

// GetFollowerCountById 根据userId查找粉丝数
func GetFollowerCountById(userId uint) (int, error) {
	baseSlice := []string{FollowerList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	size, err := Rdb.SCard(Ctx, key).Result()
	return int(size), err
}

// GetFollowListById 根据userId查找关注list
func GetFollowListById(userId uint) ([]uint, error) {
	baseSlice := []string{FollowList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
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
	baseSlice := []string{FollowerList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
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
	baseSliceFollower := []string{FollowerList, strconv.Itoa(int(userId))}
	keyFollower := strings.Join(baseSliceFollower, Delimiter)
	baseSliceFollow := []string{FollowList, strconv.Itoa(int(userId))}
	keyFollow := strings.Join(baseSliceFollow, Delimiter)

	friend, err := Rdb.SInter(Ctx, keyFollow, keyFollower).Result()
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
	baseSlice := []string{FollowList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	pipe := Rdb.Pipeline()
	for _, value := range ids {
		err := pipe.SAdd(Ctx, key, value).Err()
		if err != nil {
			return err
		}
	}
	zap.L().Info("Follow_LIST", zap.Any("List", ids))
	_, err := pipe.Exec(Ctx)
	randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}

// SetFollowerListByUserId 设置粉丝列表
func SetFollowerListByUserId(userId uint, ids []uint) error {
	baseSlice := []string{FollowerList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
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
	randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}

// IncreaseFollowCountByUserId 给Id对应的关注set加上 id
func IncreaseFollowCountByUserId(userId uint, id uint) error {
	baseSlice := []string{FollowList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	err := Rdb.SAdd(Ctx, key, id).Err()
	return err
}

// DecreaseFollowCountByUserId userId 关注列表取关 followId
func DecreaseFollowCountByUserId(userId uint, followId uint) error {
	baseSlice := []string{FollowList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	err := Rdb.SRem(Ctx, key, followId).Err()
	return err
}

// IncreaseFollowerCountByUserId 给userId粉丝列表加上 followid
func IncreaseFollowerCountByUserId(userId uint, followId uint) error {
	baseSlice := []string{FollowerList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	err := Rdb.SAdd(Ctx, key, followId).Err()
	return err
}

// DecreaseFollowerCountByUserId 给userId对应的粉丝列表减去id
func DecreaseFollowerCountByUserId(userId uint, id uint) error {
	baseSlice := []string{FollowerList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	err := Rdb.SRem(Ctx, key, id).Err()
	return err
}

// IsInMyFollowList userid的follow list是否存在id
func IsInMyFollowList(userId uint, id uint) (bool, error) {
	baseSlice := []string{FollowList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	found, err := Rdb.SIsMember(Ctx, key, id).Result()
	return found, err
}

// IsInMyFollowerList userid的follower list是否存在id
func IsInMyFollowerList(userId uint, id uint) (bool, error) {
	baseSlice := []string{FollowerList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	found, err := Rdb.SIsMember(Ctx, key, id).Result()
	return found, err
}

// ActionFollow
// 更新fromUserId关注和toUserId粉丝
func ActionFollow(fromUserId, toUserId uint) error {
	baseSliceFollow := []string{FollowList, strconv.Itoa(int(fromUserId))}
	keyFollow := strings.Join(baseSliceFollow, Delimiter)
	baseSliceFollower := []string{FollowerList, strconv.Itoa(int(toUserId))}
	keyFollower := strings.Join(baseSliceFollower, Delimiter)

	pipe := Rdb.TxPipeline()
	pipe.SAdd(Ctx, keyFollow, toUserId)
	pipe.SAdd(Ctx, keyFollower, fromUserId)
	_, err := pipe.Exec(Ctx)
	return err
}

// ActionCancelFollow
// 更新fromUserId关注和toUserId粉丝
func ActionCancelFollow(fromUserId, toUserId uint) error {
	baseSliceFollow := []string{FollowList, strconv.Itoa(int(fromUserId))}
	keyFollow := strings.Join(baseSliceFollow, Delimiter)
	baseSliceFollower := []string{FollowerList, strconv.Itoa(int(toUserId))}
	keyFollower := strings.Join(baseSliceFollower, Delimiter)

	pipe := Rdb.TxPipeline()
	pipe.SRem(Ctx, keyFollow, toUserId)
	pipe.SRem(Ctx, keyFollower, fromUserId)
	_, err := pipe.Exec(Ctx)
	return err
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
