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
	return videos, DB.Where("author_id = ?", authorId).Find(&videos).RowsAffected != 0
}

// InsertVideo return 插入视频的id，是否插入成功
func InsertVideo(videoUrl string, coverUrl string, authorID uint, title string) (uint, bool) {
	video := models.Video{
		AuthorId:  authorID,
		VideoUrl:  videoUrl,
		CoverUrl:  coverUrl,
		Title:     title,
		CreatedAt: time.Now(),
	}
	result := DB.Create(&video)
	if result.Error != nil {
		return uint(10), false
	}
	return video.VideoId, true
}

func GetLatestVideos(latestTime string) []models.Video {
	var videos []models.Video
	result := DB.Where("created_at < ?", latestTime).Order("created_at DESC").Limit(30).Find(&videos)
	if result.Error != nil {
		return []models.Video{}
	}
	return videos
}
