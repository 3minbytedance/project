package mysql

import (
	"gorm.io/gorm"
	"log"
	"project/models"
)

func AddComment(comment *models.Comment) (uint, error) {
	result := DB.Model(models.Comment{}).Create(comment)
	// 判断是否创建成功
	if result.Error != nil {
		log.Println("创建 Comment 失败:", result.Error)
		return 0, result.Error
	} else {
		log.Println("成功创建 Comment")
		return comment.ID, nil
	}
}

func FindCommentsByVideoId(videoId uint) ([]models.Comment, error) {
	comments := make([]models.Comment, 0)
	result := DB.Where("video_id = ?", videoId).Find(&comments)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return comments, result.Error
	}
	return comments, nil
}

func FindCommentById(commentId uint) (models.Comment, error) {
	comment := models.Comment{}
	result := DB.Find(&comment, commentId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return comment, result.Error
	}
	return comment, nil
}

func DeleteCommentById(commentId uint) error {
	result := DB.Delete(&models.Comment{}, commentId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		log.Println("未找到 Comment")
		return result.Error
	}
	return nil
}

func GetCommentCnt(videoId uint) (int64, error) {
	var cnt int64
	err := DB.Model(&models.Comment{}).Where("video_id = ?", videoId).Count(&cnt).Error
	// 返回评论数和是否查询成功
	return cnt, err
}
