package models

import (
	"time"
)

type Video struct {
	VideoId   uint `gorm:"primaryKey"`
	AuthorId  uint `gorm:"index"`
	VideoUrl  string
	CoverUrl  string
	Title     string
	CreatedAt time.Time
	DeletedAt time.Time
}

type VideoResponse struct {
	Id            uint64 `json:"id"`
	User          User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount uint64 `json:"favorite_count"`
	CommentCount  uint64 `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
}

type VideoListResponse struct {
	Response
	NextTime      uint64          `json:"next_time,omitempty"`
	VideoResponse []VideoResponse `json:"video_list,omitempty"`
}

func (*Video) TableName() string {
	return "video"
}
