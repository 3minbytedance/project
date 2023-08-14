package model

import (
	"gorm.io/gorm"
	"time"
)

// Comment 数据库Model
type Comment struct {
	ID        uint `gorm:"primaryKey"`
	VideoId   uint `gorm:"index"` // 非唯一索引
	UserId    uint `gorm:"index"` // 非唯一索引
	Content   string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (*Comment) TableName() string {
	return "comments"
}
