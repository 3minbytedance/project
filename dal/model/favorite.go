package model

type Favorite struct {
	ID      uint `gorm:"primaryKey"`
	UserId  uint `gorm:"index;not null"`
	VideoId uint `gorm:"index;not null"`
}

type FavoriteAction struct {
	UserId     uint
	VideoId    uint
	ActionType int
}

func (*Favorite) TableName() string {
	return "favorite"
}

type FavoriteListResponse struct {
	FavoriteRes   Response
	VideoResponse []VideoResponse `json:"video_list"`
}
