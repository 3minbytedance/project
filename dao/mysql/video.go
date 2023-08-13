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
	var videos []models.Video
	return videos, DB.Where("id = ?", authorId).Find(&videos).RowsAffected != 0
}

func FindWorkCountsByAuthorId(authorId uint) int64 {
	var videos []models.Video
	return DB.Where("id = ?", authorId).Find(&videos).RowsAffected
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
	if result.RowsAffected != 0 {
		return false
	}
	return true
}

func GetLatestVideos(latestTime string) []models.Video {
	videos := make([]models.Video, 0)

	result := DB.Model(&models.Video{}).Where("created_at < ?", latestTime).Order("created_at DESC").Limit(30).Find(&videos)
	if result.RowsAffected == 0 {
		return []models.Video{}
	}
	return videos
}
