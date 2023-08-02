package models

import (
	"fmt"
	"gorm.io/gorm"
	"project/utils"
	"strconv"
	"time"
)

type Comments struct {
	gorm.Model
	VideoId       int64
	UserId        int64
	ReplyToUserId int64
	Content       string
	CreateTime    int64
}

func (*Comments) TableName() string {
	return "comments"
}

// Comment 原来demo的接口保留了
type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

func FindCommentsByVideoId(db *gorm.DB, videoId int) ([]Comments, bool) {
	comments := make([]Comments, 0)
	return comments, db.Where("video_id = ?", videoId).Find(&comments).RowsAffected != 0
}

// CommentTime 用来显示评论时间的，比如刚刚，几分钟之前，几小时之前，日期
func CommentTime(createTime int64, currentTime int64) string {
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

func GetComments(videoId int) []Comment {
	commentsByVideoId, b := FindCommentsByVideoId(utils.DB, videoId)
	if !b {
		fmt.Println("根据视频ID取评论失败")
		return nil
	}
	comments := make([]Comment, 0)
	for i, comm := range commentsByVideoId {
		user, b := FindUserByID(utils.DB, int(comm.UserId))
		if !b {
			fmt.Println("根据评论中的user_id找用户失败")
		}
		comment := Comment{
			Id:         int64(i),
			User:       user,
			Content:    comm.Content,
			CreateDate: CommentTime(comm.CreateTime, time.Now().Unix()),
		}
		comments = append(comments, comment)
	}
	return comments
}
