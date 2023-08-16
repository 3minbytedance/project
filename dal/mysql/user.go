package mysql

import (
	"douyin/common"
	"douyin/dal/model"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
)

func FindUserByName(name string) (user model.User, exist bool) {
	user = model.User{}
	if err := DB.Where("username = ?", name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, false
		}
		// 处理其他查询错误
		zap.L().Error("Database err", zap.Error(err))
		return user, false
	}
	return user, true
}

func FindUserByUserID(id uint) (user model.User, exist bool) {
	user = model.User{}
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, false
		}
		// 处理其他查询错误
		zap.L().Error("Database err", zap.Error(err))
		return user, false
	}
	return user, true
}

// TODO 待改

func RegisterUserInfo(username string, password string) (int32, string, uint) {

	user := model.User{}
	user.Username = username

	// id默认自增
	//user.Id = uuid.New()

	// 将信息存储到数据库中
	salt := fmt.Sprintf("%06d", rand.Int())
	user.Salt = salt
	user.Password = common.MakePassword(password, salt)

	// 数据入库
	DB.Create(&user)
	fmt.Println("<<<<<<<<<id: ", user.ID)
	return 0, "注册成功", user.ID
}

func CreateUser(user *model.User) (id uint, err error) {
	// 数据入库
	err = DB.Create(&user).Error
	id = user.ID
	return
}
