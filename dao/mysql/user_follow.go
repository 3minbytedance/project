package mysql

import (
	"gorm.io/gorm"
	"log"
	"project/models"
)

// 添加
func AddFollow(user_id, follow_id uint) (bool, error) {
	follow := models.User_follow{user_id, follow_id}
	result := DB.Model(models.User_follow{}).Create(follow)
	// 判断是否创建成功
	if result.Error != nil {
		log.Println("创建 follow 失败:", result.Error)
		return false, result.Error
	} else {
		log.Println("成功创建 follow")
		return true, nil
	}
}

// delete
func DeleteFollowById(user_id, follow_id uint) error {
	follow := models.User_follow{user_id, follow_id}
	result := DB.Delete(&models.Comment{}, follow)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		log.Println("未找到 Comment")
		return result.Error
	}
	return nil
}

// 关注数
func GetFollowCnt(user_id uint) (int64, error) {
	var cnt int64
	err := DB.Model(&models.User_follow{}).Where("user_id = ?", user_id).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}

// 粉丝数
func GetFollowerCnt(user_id uint) (int64, error) {
	var cnt int64
	err := DB.Model(&models.User_follow{}).Where("follow_Id = ?", user_id).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}
