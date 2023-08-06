package models

import (
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"project/utils"
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

// todo 废弃，jwt解析自带信息
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

func CheckUserRegisterInfo(username string, password string) (int32, string) {

	if len(username) == 0 || len(username) > 32 {
		return 1, "用户名不合法"
	}

	if len(password) <= 6 || len(password) > 32 {
		return 2, "密码不合法"
	}

	if _, ok := FindUserByName(utils.DB, username); ok {
		return 3, "用户已注册"
	}

	return 0, "合法"
}

func RegisterUserInfo(username string, password string) (int32, string, int64, string) {

	// todo 对密码加密
	user := User{}
	user.Name = username

	// 生成token，id
	//user.ID = uuid.New()
	// 将信息存储到数据库中

	// salt密码加密
	userStates := UserStates{}
	userStates.Name = username
	salt := fmt.Sprintf("%06d", rand.Int())
	userStates.Salt = salt
	userStates.Password = utils.MakePassword(password, salt)
	userStates.Token = utils.GenerateToken(int64(user.ID), username)

	// 数据入库
	utils.DB.Create(&userStates)
	utils.DB.Create(&user)
	fmt.Println("<<<<<<<<<id: ", user.ID)
	return 0, "注册成功", int64(user.ID), userStates.Token
}
