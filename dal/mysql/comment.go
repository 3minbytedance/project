package mysql

import (
	"douyin/dal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

func AddComment(comment *model.Comment) (uint, error) {
	result := DB.Model(model.Comment{}).Create(comment)
	// 判断是否创建成功
	if result.Error != nil {
		zap.L().Error("创建 Comment 失败:", zap.Error(result.Error))
		return 0, result.Error
	} else {
		return comment.ID, nil
	}
}

func FindCommentsByVideoId(videoId uint) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	result := DB.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	log.Println(comments)
	return comments, nil
}

func FindCommentById(commentId uint) (model.Comment, error) {
	comment := model.Comment{}
	result := DB.Find(&comment, commentId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return comment, result.Error
	}
	return comment, nil
}

func DeleteCommentById(commentId uint) error {
	result := DB.Delete(&model.Comment{}, commentId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		log.Println("未找到 Comment")
		return result.Error
	}
	return nil
}

func GetCommentCnt(videoId uint) (int64, error) {
	var cnt int64
	err := DB.Model(&model.Comment{}).Where("video_id = ?", videoId).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}