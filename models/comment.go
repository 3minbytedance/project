package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Comment 数据库Model
type Comment struct {
	CommentId uint `gorm:"primaryKey"`
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

// TranslateTime 用来显示评论时间的，比如刚刚，几分钟之前，几小时之前，日期
func TranslateTime(createTime int64, currentTime int64) string {
	diff := currentTime - createTime
	if diff < 0 {
		fmt.Println("时间差小于0，数据存储可能出错")
		return "未知"
	}
	var res string
	if diff < 60 {
		res = "刚刚"
	} else if diff < 60*60 {
		res = strconv.Itoa(int(diff/60)) + "分钟前"
	} else if diff < 60*60*24 {
		res = strconv.Itoa(int(diff/3600)) + "小时前"
	} else {
		year, month, day := time.Unix(createTime, 0).Date()
		res = strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	}
	return res
}
