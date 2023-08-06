package models

import (
	"gorm.io/gorm"
)

type User struct {
	// 用户的信息社交平台信息， todo 作为常用个人信息，如果获取过于复杂可以考虑在redis中储存
	gorm.Model
	//Id              int64  `json:"id,omitempty"`               // 用户id
	Name            string `json:"name,omitempty"`             // 用户名称
	FollowCount     int64  `json:"follow_count,omitempty"`     // 关注总数
	FollowerCount   int64  `json:"follower_count,omitempty"`   // 粉丝总数
	IsFollow        bool   `json:"is_follow,omitempty"`        // true-已关注，false-未关注
	Avatar          string `json:"avatar,omitempty"`           // 用户头像
	BackgroundImage string `json:"background_image,omitempty"` // 用户个人页顶部大图
	Signature       string `json:"signature,omitempty"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`            // 获赞数量
	WorkCount       int64  `json:"work_count"`                 // 作品数量
	FavoriteCount   int64  `json:"favorite_count"`             // 点赞数量
}

func (*User) TableName() string {
	return "user"
}

type UserStates struct {
	gorm.Model
	Name      string
	Password  string
	Salt      string
	Token     string
	LoginTime int64
	IsLogOut  bool
}

func (*UserStates) TableName() string {
	return "user_states"
}
