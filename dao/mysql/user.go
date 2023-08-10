package mysql

import (
	"project/models"
)

func FindUserByName(name string) (user models.User, exist bool) {
	user = models.User{}
	return user, DB.Where("name = ?", name).First(&user).RowsAffected != 0
}

//func FindUserStateByName(name string) (models.user, bool) {
//	userState := models.user{}
//	return userState, DB.Where("name = ?", name).First(&userState).RowsAffected != 0
//}

func FindUserByID(id uint) (models.User, bool) {
	user := models.User{}
	return user, DB.Where("id = ?", id).First(&user).RowsAffected != 0
}

//func FindUserStateByID(id int) (models.user, bool) {
//	userState := models.user{}
//	return userState, DB.Where("id = ?", id).First(&userState).RowsAffected != 0
//}

func FindUserInfoByUserId(userId uint) (models.User, bool) {
	user := models.User{}
	row := DB.Where("Id = ?", userId).First(&user).RowsAffected
	if row == 0 {
		return models.User{}, false
	}

	return user, true
}

func CreateUser(user *models.User) (id uint, err error) {
	// 数据入库
	err = DB.Create(&user).Error
	id = user.Id
	return

}
