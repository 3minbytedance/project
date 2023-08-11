package service

import (
	"fmt"
	"log"
	"project/dao/mysql"
	"project/dao/redis"
	"project/models"
)

func AddFollow(userId, followId uint) error {
	// 评论信息
	err := mysql.AddFollow(userId, followId)
	return err
}

func DeleteFollow(userId, followId uint) error {
	err := mysql.DeleteFollowById(userId, followId)
	return err
}

func GetFollowList(userId uint) ([]models.UserInfo, error) {
	follow, err := mysql.GetFollowList(userId)
	if err != nil {
		return nil, err
	}
	results, err := GetUserModelByList(follow)

	return results, err
}

func GetFollowerList(userId uint) ([]models.UserInfo, error) {
	follower, err := mysql.GetFollowerList(userId)
	if err != nil {
		return nil, err
	}
	results, err := GetUserModelByList(follower)

	return results, err
}

func GetFriendList(userId uint) ([]models.UserInfo, error) {
	follow, err := mysql.GetFollowList(userId)
	if err != nil {
		return nil, err
	}
	follower, err := mysql.GetFollowerList(userId)
	if err != nil {
		return nil, err
	}
	friend := intersection(follow, follower)

	results, err := GetUserModelByList(friend)
	return results, err
}

func GetFollowCount(userID uint) (int64, error) {
	// 从redis中获取评论数
	count, err := redis.GetFollowCountById(userID)
	if err != nil {
		log.Println("从redis中获取失败：", err)
		return 0, err
	}
	// 1. 缓存中有数据, 直接返回

	if count >= 0 {
		log.Println("从redis中获取评论数成功：", count)
		return int64(count), nil
	}

	// 2. 缓存中没有数据，从数据库中获取
	num, err := mysql.GetFollowCnt(userID)
	if err != nil {
		log.Println("从数据库中获取关注数失败：", err.Error())
		return 0, nil
	}
	log.Println("从数据库中获取关注数成功：", num)
	// 将评论数写入redis
	go func() {
		err = redis.SetFollowCountByUserId(userID, num)
		if err != nil {
			log.Println("将评论数写入redis失败：", err.Error())
		}
	}()
	return num, nil
}

func GetFollowerCount(userID uint) (int64, error) {
	// 从redis中获取评论数
	count, err := redis.GetFollowerCountById(userID)
	if err != nil {
		log.Println("从redis中获取失败：", err)
		return 0, err
	}
	// 1. 缓存中有数据, 直接返回
	if count >= 0 {
		log.Println("从redis中获取评论数成功：", count)
		return int64(count), nil
	}

	// 2. 缓存中没有数据，从数据库中获取
	num, err := mysql.GetFollowerCnt(userID)
	if err != nil {
		log.Println("从数据库中获取关注数失败：", err.Error())
		return 0, nil
	}
	log.Println("从数据库中获取关注数成功：", num)
	// 将评论数写入redis
	go func() {
		err = redis.SetFollowerCountByUserId(userID, num)
		if err != nil {
			log.Println("将评论数写入redis失败：", err.Error())
		}
	}()
	return num, nil
}

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

// 根据id获取model
func GetUserModelByList(id []uint) ([]models.UserInfo, error) {
	var results []models.UserInfo
	for _, value := range id {
		result, ok := GetUserInfoByUserId(value)
		if !ok {
			return nil, fmt.Errorf("please try again")
		}
		results = append(results, result)
	}
	return results, nil
}
