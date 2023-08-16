package model

import "gorm.io/gorm"

type Favorite struct {
	ID       uint `gorm:"primaryKey"`
	UserId   uint `gorm:"index"`
	VideoId  uint `gorm:"index"`
	DeleteAt gorm.DeletedAt
}

func (*Favorite) TableName() string {
	return "favorite"
}

type FavoriteListResponse struct {
	FavoriteRes   Response
	VideoResponse []VideoResponse `json:"video_list"`
}
