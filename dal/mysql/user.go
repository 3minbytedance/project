package mysql

import (
	"douyin/dal/model"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

func FindUserByName(name string) (user model.User, exist bool, err error) {
	user = model.User{}
	if err = DB.Where("name = ?", name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, false, nil
		}
		// 处理其他查询错误
		zap.L().Error("Database err", zap.Error(err))
		return user, false, err
	}
	return user, true, nil
}

func FindUserByUserID(id uint) (user model.User, exist bool, err error) {
	user = model.User{}
	if err = DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, false, err
		}
		// 处理其他查询错误
		zap.L().Error("Database err", zap.Error(err))
		return user, false, err
	}
	return user, true, err

}

func CreateUser(user *model.User) (id uint, err error) {
	// 数据入库
	err = DB.Create(&user).Error
	id = user.ID
	return
}
