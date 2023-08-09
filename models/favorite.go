package models

import "gorm.io/gorm"

type Favorite struct {
	Id       uint
	UserId   int64
	VideoId  int64
	DeleteAt gorm.DeletedAt
}

func (*Favorite) TableName() string {
	return "favorite"
}
