package models

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

// CommentResponse 返回数据的Model
type CommentResponse struct {
	Id         int64        `json:"id,omitempty"`
	User       UserResponse `json:"user"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
}

type CommentListResponse struct {
	Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment CommentResponse `json:"comment,omitempty"`
}

// TranslateTime 返回mm-dd格式
func TranslateTime(createTime int64) string {
	t := time.Unix(createTime, 0)
	return t.Format("01-02")
}
