package models

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	VideoId   uint `gorm:"primaryKey"`
	AuthorId  uint `gorm:"index"`
	VideoUrl  string
	CoverUrl  string
	Title     string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type VideoResponse struct {
	Id            uint     `json:"id"`
	Author        UserInfo `json:"author"` //TODO USERINFO?
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount uint     `json:"favorite_count"` //点赞数
	CommentCount  uint     `json:"comment_count"`  //评论数
	IsFavorite    bool     `json:"is_favorite"`    //是否点赞
}

// VideoListResponse 用户所有投稿过的视频
type VideoListResponse struct {
	Response
	VideoResponse []VideoResponse `json:"video_list,omitempty"`
}

// FeedListResponse 投稿时间倒序的视频列表
type FeedListResponse struct {
	Response
	NextTime      uint64          `json:"next_time,omitempty"`
	VideoResponse []VideoResponse `json:"video_list,omitempty"`
}

func (*Video) TableName() string {
	return "video"
}
