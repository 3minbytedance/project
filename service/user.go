package service

import (
	"fmt"
	"math/rand"
	"project/dao/mysql"
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

	user, exist := mysql.FindUserByUserID(userId)
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
	userResponse := models.UserResponse{
		Id:              user.UserId,
		Name:            user.Name,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        false,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  "0",
		WorkCount:       0, //todo
		FavoriteCount:   0,
	}
	return userResponse, true
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

	return 0, "合法"
}
