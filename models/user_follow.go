package models

type UserFollow struct {
	// 用户的关注信息
	UserId   uint `gorm:"primaryKey;index;not null"` // 用户id
	FollowId uint `gorm:"primaryKey;index;not null"` // 关注用户id
}

func (*UserFollow) TableName() string {
	return "user_follow"
}
