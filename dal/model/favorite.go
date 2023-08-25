package model

type Favorite struct {
	ID      uint `gorm:"primaryKey"`
	UserId  uint `gorm:"index;not null"`
	VideoId uint `gorm:"index;not null"`
}

func (*Favorite) TableName() string {
	return "favorite"
}

type FavoriteListResponse struct {
	FavoriteRes   Response
	VideoResponse []VideoResponse `json:"video_list"`
}
