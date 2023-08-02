package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorId      int64
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"` // 这是demo中的字段，感觉在视频信息里面存不妥，应该放在favorite表里面
	// 下面是一些保留字段， 部分是用来推荐的
	PublishTime int64  // 发布时间
	Title       string // 视频标题
	Topic       string // 视频主题类型
	IsLong      int    // 视频长度是否大于1分钟 0 否， 1是
}

// VideoRes demo中的struct，暂且保留了
type VideoRes struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

func (*Video) TableName() string {
	return "video"
}

func FindVideoByVideoId(db *gorm.DB, videoId int) (Video, bool) {
	video := Video{}
	return video, db.Where("id = ?", videoId).First(&video).RowsAffected != 0
}

func FindVideosByAuthor(db *gorm.DB, authorId int) ([]VideoRes, bool) {
	videos := make([]Video, 0)
	num := db.Where("author_id = ?", authorId).Find(&videos).RowsAffected
	if num == 0 {
		return nil, false
	}
	videosRes := make([]VideoRes, 0)
	for _, v := range videos {
		user, _ := FindUserByID(db, int(v.AuthorId))
		temp := VideoRes{
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
