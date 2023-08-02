package models

import "gorm.io/gorm"

// UserActions 记录用户行为的，后面用来进行视频的推荐，字段暂定
type UserActions struct {
	gorm.Model
	UserId        int64 // 操作者ID
	VideoId       int64 // 视频ID
	ActionType    int   // 0 点赞 1 评论 2 分享
	WhetherFinish int   // 0 没看完 1 看完
	WatchingCount int   // 观看数量
	ActionTime    int64 // 观看时间
}

func (*UserActions) TableName() string {
	return "user_action"
}
