package mysql

import (
	"fmt"
	"math/rand"
	"project/models"
	"project/utils"
)

func FindUserByName(name string) (models.User, bool) {
	user := models.User{}
	return user, DB.Where("name = ?", name).First(&user).RowsAffected != 0
}

//func FindUserStateByName(name string) (models.user, bool) {
//	userState := models.user{}
//	return userState, DB.Where("name = ?", name).First(&userState).RowsAffected != 0
//}

func FindUserByID(id uint) (models.UserInfo, bool) {
	user := models.UserInfo{}
	return user, DB.Where("id = ?", id).First(&user).RowsAffected != 0
}

//func FindUserStateByID(id int) (models.user, bool) {
//	userState := models.user{}
//	return userState, DB.Where("id = ?", id).First(&userState).RowsAffected != 0
//}

//// FindUserByToken todo 废弃，jwt解析自带信息
//func FindUserByToken(token string) (models.User, bool) {
//	user := models.User{}
//	userState := models.user{}
//	row := DB.Where("token = ?", token).First(&userState).RowsAffected
//	if row == 0 {
//		return user, false
//	}
//	// 应该在user表里面加id，而不是name
//	return user, DB.Where("name = ?", userState.Name).First(&user).RowsAffected != 0
//}

func FindUserInfoByUserId(userId uint) (models.User, bool) {
	user := models.User{}

	row := DB.Where("Id = ?", userId).First(&user).RowsAffected
	if row == 0 {
		return user, false
	}
}

func CheckUserRegisterInfo(username string, password string) (int32, string) {

	if len(username) == 0 || len(username) > 32 {
		return 1, "用户名不合法"
	}

	if len(password) <= 6 || len(password) > 32 {
		return 2, "密码不合法"
	}

	if _, ok := FindUserByName(username); ok {
		return 3, "用户已注册"
	}

	return 0, "合法"
}

func RegisterUserInfo(username string, password string) (int32, string, uint) {

	user := models.User{}
	user.Name = username

	// id默认自增
	//user.Id = uuid.New()

	// 将信息存储到数据库中
	salt := fmt.Sprintf("%06d", rand.Int())
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt)

	// 数据入库
	DB.Create(&user)
	fmt.Println("<<<<<<<<<id: ", user.Id)
	return 0, "注册成功", user.Id
}
