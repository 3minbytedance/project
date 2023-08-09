package mysql

import (
	"project/models"
)

func FindVideoByVideoId(videoId int64) (models.Video, bool) {
	video := models.Video{}
	return video, DB.Where("id = ?", videoId).First(&video).RowsAffected != 0
}

func FindVideosByAuthor(authorId int) ([]models.VideoRes, bool) {
	videos := make([]models.Video, 0)
	num := DB.Where("author_id = ?", authorId).Find(&videos).RowsAffected
	if num == 0 {
		return nil, false
	}
	videosRes := make([]models.VideoRes, 0)
	for _, v := range videos {
		user, _ := FindUserByID(uint(v.AuthorId))
		temp := models.VideoRes{
			Id:            int64(v.ID),
			Author:        user,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		}
		videosRes = append(videosRes, temp)
	}
	return videosRes, true
}
