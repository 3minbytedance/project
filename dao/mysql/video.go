package mysql

import (
	"project/models"
	"time"
)

func FindVideoByVideoId(videoId int64) (models.Video, bool) {
	video := models.Video{}
	result := DB.First(&video, videoId)
	if result.Error != nil {
		return video, false
	}
	return video, true
}

func FindVideosByAuthor(authorId uint) ([]models.Video, bool) {
	var videos []models.Video
	result := DB.Where("author_id = ?", authorId).Find(&videos)
	if result.Error != nil {
		return []models.Video{}, false
	}
	return videos, true
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
		return -1, false
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
