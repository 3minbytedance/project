package models

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	ID        uint `gorm:"primaryKey"`
	AuthorId  uint `gorm:"index"`
	VideoUrl  string
	CoverUrl  string
	Title     string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type VideoResponse struct {
	ID            uint         `json:"id"`
	Author        UserResponse `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount int64        `json:"favorite_count"` //点赞数
	CommentCount  int64        `json:"comment_count"`  //评论数
	IsFavorite    bool         `json:"is_favorite"`    //是否点赞
	Title         string       `json:"title"`          //视频标题
}

// VideoListResponse 用户所有投稿过的视频
type VideoListResponse struct {
	Response
	VideoResponse []VideoResponse `json:"video_list,omitempty"`
}

// FeedListResponse 投稿时间倒序的视频列表
type FeedListResponse struct {
	Response
	NextTime      int64           `json:"next_time,omitempty"`
	VideoResponse []VideoResponse `json:"video_list,omitempty"`
}

func (*Video) TableName() string {
	return "video"
}
