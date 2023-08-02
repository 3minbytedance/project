package models

import "gorm.io/gorm"

type User struct {
	// 这是刷到的用户的信息，不是当前用户的信息
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

func FindUserByName(db *gorm.DB, name string) (User, bool) {
	user := User{}
	return user, db.Where("name = ?", name).First(&user).RowsAffected != 0
}

func FindUserStateByName(db *gorm.DB, name string) (UserStates, bool) {
	userState := UserStates{}
	return userState, db.Where("name = ?", name).First(&userState).RowsAffected != 0
}

func FindUserByID(db *gorm.DB, id int) (User, bool) {
	user := User{}
	return user, db.Where("id = ?", id).First(&user).RowsAffected != 0
}

func FindUserStateByID(db *gorm.DB, id int) (UserStates, bool) {
	userState := UserStates{}
	return userState, db.Where("id = ?", id).First(&userState).RowsAffected != 0
}

func FindUserByToken(db *gorm.DB, token string) (User, bool) {
	user := User{}
	userState := UserStates{}
	row := db.Where("token = ?", token).First(&userState).RowsAffected
	if row == 0 || userState.IsLogOut {
		return user, false
	}
	// 应该在userStates表里面加id，而不是name
	return user, db.Where("name = ?", userState.Name).First(&user).RowsAffected != 0
}
