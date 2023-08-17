package mysql

import (
	"douyin/dal/model"
	"time"
)

func FindVideoByVideoId(videoId uint) (model.Video, bool) {
	video := model.Video{}
	return video, DB.Where("id = ?", videoId).First(&video).RowsAffected != 0
}

// FindVideosByAuthorId 返回查询到的列表及是否出错
// 若未找到，返回空列表
func FindVideosByAuthorId(authorId uint) ([]model.Video, bool) {
	videos := make([]model.Video, 0)
	return videos, DB.Where(" author_id = ?", authorId).Find(&videos).RowsAffected != 0
}

func FindWorkCountsByAuthorId(authorId uint) int64 {
	var count int64
	DB.Model(&model.Video{}).Where("author_id = ?", authorId).Count(&count)
	return count
}

// InsertVideo return 是否插入成功
func InsertVideo(videoUrl string, coverUrl string, authorID uint, title string) bool {
	video := model.Video{
		AuthorId:  authorID,
		VideoUrl:  videoUrl,
		CoverUrl:  coverUrl,
		Title:     title,
		CreatedAt: time.Now().Unix(),
	}
	result := DB.Model(model.Video{}).Create(&video)
	return result.RowsAffected != 0
}

func GetLatestVideos(latestTime string) []model.Video {
	videos := make([]model.Video, 0)

	DB.Model(&model.Video{}).Where("created_at < ?", latestTime).Order("created_at DESC").Limit(30).Find(&videos)

	return videos
}