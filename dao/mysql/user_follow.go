package mysql

import (
	"gorm.io/gorm"
	"log"
	"project/models"
)

// 添加
func AddFollow(userId, followId uint) error {
	follow := models.UserFollow{UserId: userId, FollowId: followId}
	result := DB.Model(&models.UserFollow{}).Create(follow)
	// 判断是否创建成功
	if result.Error != nil {
		log.Println("创建 follow 失败:", result.Error)
		return result.Error
	} else {
		log.Println("成功创建 follow")
		return nil
	}
}

// delete
func DeleteFollowById(userId, followId uint) error {
	follow := models.UserFollow{UserId: userId, FollowId: followId}
	result := DB.Delete(&models.Comment{}, follow)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		log.Println("未找到 Comment")
		return result.Error
	}
	return nil
}

// 关注数
func GetFollowCnt(userId uint) (int64, error) {
	var cnt int64
	err := DB.Model(&models.UserFollow{}).Where("user_id = ?", userId).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}

// 粉丝数
func GetFollowerCnt(userId uint) (int64, error) {
	var cnt int64
	err := DB.Model(&models.UserFollow{}).Where("follow_id = ?", userId).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}

// 获取关注列表
func GetFollowList(userId uint) ([]uint, error) {
	var result []uint
	err := DB.Model(&models.UserFollow{}).Select("follow_id").Where("user_id = ?", userId).Scan(&result).Error
	return result, err
}

// 获取粉丝列表
func GetFollowerList(userId uint) ([]uint, error) {
	var result []uint
	err := DB.Model(&models.UserFollow{}).Select("user_id").Where("follow_id = ?", userId).Scan(&result).Error
	return result, err
}
