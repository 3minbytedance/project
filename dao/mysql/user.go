package mysql

import (
	"fmt"
	"math/rand"
	"project/models"
	"project/utils"
)

func FindUserByName(name string) (user models.User, exist bool) {
	user = models.User{}
	return user, DB.Where("name = ?", name).First(&user).RowsAffected != 0
}

func FindUserByUserID(id uint) (models.User, bool) {
	user := models.User{}
	return user, DB.Where("id = ?", id).First(&user).RowsAffected != 0
}

// TODO 待改

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
	fmt.Println("<<<<<<<<<id: ", user.ID)
	return 0, "注册成功", user.ID
}

func CreateUser(user *models.User) (id uint, err error) {
	// 数据入库
	err = DB.Create(&user).Error
	id = user.ID
	return

}
