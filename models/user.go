package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	// 用户的信息社交平台信息， todo 作为常用个人信息，如果获取过于复杂可以考虑在redis中储存
	Id              uint   `gorm:"primaryKey"`          // 用户id
	Name            string `gorm:"uniqueIndex;size:32"` // 用户名称
	Password        string // 用户密码
	Avatar          string // 用户头像
	BackgroundImage string // 用户个人页顶部大图
	Signature       string // 个人简介
	Salt            string //加密盐
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

type UserResponse struct {
	// 用户的信息社交平台信息， todo 作为常用个人信息，如果获取过于复杂可以考虑在redis中储存
	Id              uint   `json:"id"`               // 用户id
	Name            string `json:"name"`             // 用户名称
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数量
	FavoriteCount   int64  `json:"favorite_count"`   // 点赞数量
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type UserDetailResponse struct {
	Response
	User UserResponse `json:"user,omitempty"`
}

func (*User) TableName() string {
	return "user"
}
