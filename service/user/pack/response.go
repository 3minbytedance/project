package pack

import (
	user "douyin/kitex_gen/user"
)

func User(userId int32) *user.User {
	return &user.User{
		Id:              userId,
		Name:            "",
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "",
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
}
