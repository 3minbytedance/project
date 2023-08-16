package service

import (
	"fmt"
	"log"
	"math/rand"
	"project/dao/mysql"
	"project/dao/redis"
	"project/models"
	"project/utils"
)

func RegisterUser(username string, password string) (id uint, err error) {

	user := models.User{}
	user.Name = username

	// 将信息存储到数据库中
	salt := fmt.Sprintf("%06d", rand.Int())
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt)

	// 数据入库
	userId, err := mysql.CreateUser(&user)
	return userId, err
}

func GetUserInfoByUserId(userId uint) (models.UserResponse, bool) {
	name, exist := GetName(userId)
	if !exist {
		return models.UserResponse{}, false
	}
	followCount, err := GetFollowCount(userId)
	if err != nil {
		return models.UserResponse{}, false
	}
	followerCount, err := GetFollowerCount(userId)
	if err != nil {
		return models.UserResponse{}, false
	}
	workCount, err := GetWorkCount(userId)
	if err != nil {
		return models.UserResponse{}, false
	}
	favoriteCount := GetUserFavoriteCount(userId)
	totalFavoritedCount := GetUserTotalFavoritedCount(userId)
	userResponse := models.UserResponse{
		ID:              userId,
		Name:            name,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        false,
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "",
		TotalFavorited:  totalFavoritedCount,
		WorkCount:       workCount,
		FavoriteCount:   favoriteCount,
	}
	return userResponse, true
}

// GetName 获得作品数
func GetName(userId uint) (string, bool) {
	// 从redis中获取用户名
	// 1. 缓存中有数据, 直接返回
	if redis.IsExistUserField(userId, redis.NameField) {
		name, err := redis.GetNameByUserId(userId)
		if err != nil {
			log.Println("从redis中获取用户名失败：", err)
		}
		return name, true
	}

	// 2. 缓存中没有数据，从数据库中获取
	user, exist := mysql.FindUserByUserID(userId)
	if !exist {
		return "", false
	}
	// 将用户名写入redis
	go func() {
		err := redis.SetNameByUserId(userId,user.Name)
		if err != nil {
			log.Println("将用户名写入redis失败：", err)
		}
	}()
	return user.Name, true
}

func GetUserByName(username string) (user models.User, b bool) {
	return mysql.FindUserByName(username)
}

func CheckUserRegisterInfo(username string, password string) (int32, string) {

	if len(username) == 0 || len(username) > 32 {
		return 1, "用户名不合法"
	}

	if len(password) <= 6 || len(password) > 32 {
		return 2, "密码不合法"
	}

	if _, exist := GetUserByName(username); exist {
		return 3, "用户已注册"
	}

	return 0, "注册成功"
}
