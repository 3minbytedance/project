package service

import (
	"fmt"
	"project/dao/mysql"
	"project/dao/redis"
	"project/models"
	"strconv"
)

func AddFollow(userId, followId uint) error {
	// 评论信息
	err := mysql.AddFollow(userId, followId)
	go func() {
		err := redis.IncreaseFollowCountByUserId(userId, followId)
		if err != nil {
			return
		}
	}()
	go func() {
		err := redis.IncreaseFollowerCountByUserId(followId, followId)
		if err != nil {
			return
		}
	}()
	return err
}

func DeleteFollow(userId, followId uint) error {
	err := mysql.DeleteFollowById(userId, followId)
	go func() {
		err := redis.DecreaseFollowCountByUserId(userId, followId)
		if err != nil {
			return
		}
	}()
	go func() {
		err := redis.DecreaseFollowerCountByUserId(followId, followId)
		if err != nil {
			return
		}
	}()
	return err
}

func GetFollowList(userId uint) ([]models.UserResponse, error) {
	// redis存在key
	if redis.IsExistUserSetField(userId, redis.FollowList) {
		followListId, err := redis.GetFollowListById(userId)
		if err != nil {
			return nil, err
		}
		return GetUserModelByList(followListId)
	} else {
		followListId, err := mysql.GetFollowList(userId)
		if err != nil {
			return nil, err
		}
		// 往redis赋值
		go func() {
			err = redis.SetFollowListByUserId(userId, followListId)
			if err != nil {
				return
			}
		}()
		return GetUserModelByList(followListId)
	}
}

func GetFollowerList(userId uint) ([]models.UserResponse, error) {
	// redis存在key
	if redis.IsExistUserSetField(userId, redis.FollowerList) {
		followListId, err := redis.GetFollowerListById(userId)
		if err != nil {
			return nil, err
		}
		return GetUserModelByList(followListId)
	} else {
		followListId, err := mysql.GetFollowerList(userId)
		if err != nil {
			return nil, err
		}
		// 往redis赋值
		go func() {
			err = redis.SetFollowerListByUserId(userId, followListId)
			if err != nil {
				return
			}
		}()
		return GetUserModelByList(followListId)
	}
}

func GetFriendList(userId uint) ([]models.UserResponse, error) {
	key1 := fmt.Sprintf("%d_%s", userId, redis.FollowerList)
	key2 := fmt.Sprintf("%d_%s", userId, redis.FollowList)
	friend, err := redis.Rdb.SUnion(redis.Ctx, key2, key1).Result()
	if err != nil {
		return nil, err
	}
	var friend_list []uint
	for _, value := range friend {
		k, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		friend_list = append(friend_list, uint(k))
	}
	results, err := GetUserModelByList(friend_list)
	return results, err
}

func GetFollowCount(userId uint) (int64, error) {
	// redis存在key
	if redis.IsExistUserSetField(userId, redis.FollowList) {
		num, err := redis.GetFollowCountById(userId)
		return int64(num), err
	} else {
		followListId, err := mysql.GetFollowList(userId)
		if err != nil {
			return -1, err
		}
		// 往redis赋值
		go func() {
			err = redis.SetFollowListByUserId(userId, followListId)
			if err != nil {
				return
			}
		}()
		return int64(len(followListId)), nil
	}
}

func GetFollowerCount(userId uint) (int64, error) {
	// redis存在key
	if redis.IsExistUserSetField(userId, redis.FollowerList) {
		num, err := redis.GetFollowerCountById(userId)
		return int64(num), err
	} else {
		followListId, err := mysql.GetFollowerList(userId)
		if err != nil {
			return -1, err
		}
		// 往redis赋值
		go func() {
			err = redis.SetFollowerListByUserId(userId, followListId)
			if err != nil {
				return
			}
		}()
		return int64(len(followListId)), nil
	}
}

// GetUserModelByList 根据id获取model
func GetUserModelByList(id []uint) ([]models.UserResponse, error) {
	results := make([]models.UserResponse, len(id))
	for _, value := range id {
		result, ok := GetUserInfoByUserId(value)
		if !ok {
			return nil, fmt.Errorf("please try again")
		}
		results = append(results, result)
	}
	return results, nil
}

// todo 会耗费go的堆栈内存
// 找出两个数组共有的元素
func intersection(a, b []uint) (c []uint) {
	m := make(map[uint]bool)
	for _, item := range a {
		// storing value to the map
		m[item] = true
	}
	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return c
}
