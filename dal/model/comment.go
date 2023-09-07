package model

import (
	"gorm.io/gorm"
	"time"
)

// Comment 数据库Model
type Comment struct {
	ID        uint `gorm:"primaryKey"`
	VideoId   uint `gorm:"index"` // 非唯一索引
	UserId    uint
	Content   string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (*Comment) TableName() string {
	return "comments"
}

// CommentResponse 返回数据的Model
type CommentResponse struct {
	Id         int64        `json:"id"`
	User       UserResponse `json:"user"`
	Content    string       `json:"content"`
	CreateDate string       `json:"create_date"`
}

type CommentListResponse struct {
	Response
	CommentList []CommentResponse `json:"comment_list"`
}

type CommentActionResponse struct {
	Response
	Comment CommentResponse `json:"comment"`
}
