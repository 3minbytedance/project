package mysql

import (
	"douyin/dal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

// AddFollow 添加关注关系
func AddFollow(userId, followId uint) error {
	follow := model.UserFollow{UserId: userId, FollowId: followId}
	result := DB.Model(&model.UserFollow{}).Create(follow)
	// 判断是否创建成功
	if result.Error != nil {
		log.Println("创建 follow 失败:", result.Error)
		return result.Error
	} else {
		log.Println("成功创建 follow")
		return nil
	}
}

// DeleteFollowById 删除关注关系
func DeleteFollowById(userId, followId uint) error {
	follow := model.UserFollow{UserId: userId, FollowId: followId}
	result := DB.Delete(&model.UserFollow{}, follow)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		log.Println("未找到 Follow", userId, followId)
		return result.Error
	}
	return nil
}

// GetFollowCnt 关注数
func GetFollowCnt(userId uint) (int64, error) {
	var cnt int64
	err := DB.Model(&model.UserFollow{}).Where("user_id = ?", userId).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}

// GetFollowerCnt 粉丝数
func GetFollowerCnt(userId uint) (int64, error) {
	var cnt int64
	err := DB.Model(&model.UserFollow{}).Where("follow_id = ?", userId).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}

func IsFollowing(userA uint, userB uint) bool {
	var count int64
	DB.Model(&model.UserFollow{}).
		Where("user_id = ? AND follow_id = ?", userA, userB).
		Count(&count)
	return count > 0
}

// IsMutualFollow 是否互关
func IsMutualFollow(userA uint, userB uint) bool {
	isAFollowB := IsFollowing(userA, userB)
	isBFollowA := IsFollowing(userB, userA)
	return isAFollowB && isBFollowA
}

// GetFollowList 获取关注列表
func GetFollowList(userId uint) ([]uint, error) {
	var result []uint
	err := DB.Model(&model.UserFollow{}).Select("follow_id").Where("user_id = ?", userId).Scan(&result).Error
	return result, err
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(userId uint) ([]uint, error) {
	var result []uint
	err := DB.Model(&model.UserFollow{}).Select("user_id").Where("follow_id = ?", userId).Scan(&result).Error
	return result, err
}

func IsFriend(actorId, userId uint) (result bool, err error) {
	// 检查用户A是否关注了用户B，以及用户B是否关注了用户A
	var count int64
	res := DB.Model(&model.UserFollow{}).
		Where("user_id = ? AND follow_id = ?", actorId, userId).
		Or("user_id = ? AND follow_id = ?", userId, actorId).
		Count(&count)
	if res.Error != nil {
		zap.L().Error("Error occurred during query:", zap.Error(res.Error))
	}
	return count == 2, res.Error
}
