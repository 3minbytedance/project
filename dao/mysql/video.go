package mysql

import (
	"project/models"
	"time"
)

func FindVideoByVideoId(videoId uint) (models.Video, bool) {
	video := models.Video{}
	return video, DB.Where("id = ?", videoId).First(&video).RowsAffected != 0
}

// FindVideosByAuthorId 返回查询到的列表及是否出错
// 若未找到，返回空列表
func FindVideosByAuthorId(authorId uint) ([]models.Video, bool) {
	videos := make([]models.Video, 0)
	return videos, DB.Where(" author_id = ?", authorId).Find(&videos).RowsAffected != 0
}

func FindWorkCountsByAuthorId(authorId uint) int64 {
	var count int64
	DB.Model(&models.Video{}).Where("author_id = ?", authorId).Count(&count)
	return count
}

// InsertVideo return 是否插入成功
func InsertVideo(videoUrl string, coverUrl string, authorID uint, title string) bool {
	video := models.Video{
		AuthorId:  authorID,
		VideoUrl:  videoUrl,
		CoverUrl:  coverUrl,
		Title:     title,
		CreatedAt: time.Now().Unix(),
	}
	result := DB.Model(models.Video{}).Create(&video)
	return result.RowsAffected != 0
}

func GetLatestVideos(latestTime string) []models.Video {
	videos := make([]models.Video, 0)

	DB.Model(&models.Video{}).Where("created_at < ?", latestTime).Order("created_at DESC").Limit(30).Find(&videos)

	return videos
}