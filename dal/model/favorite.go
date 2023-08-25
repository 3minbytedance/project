package model

type Favorite struct {
	ID      uint `gorm:"primaryKey"`
	UserId  uint `gorm:"index:idx_user_video,uniqueIndex:idx_user_video"`
	VideoId uint `gorm:"index:idx_user_video,uniqueIndex:idx_user_video"`
}

func (*Favorite) TableName() string {
	return "favorite"
}

type FavoriteListResponse struct {
	FavoriteRes   Response
	VideoResponse []VideoResponse `json:"video_list"`
}
