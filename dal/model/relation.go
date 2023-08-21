package model

type UserFollow struct {
	// 用户的关注信息
	ID       int  `gorm:"primaryKey"`
	UserId   uint `gorm:"index;not null"` // 用户id
	FollowId uint `gorm:"index;not null"` // 关注用户id
}

func (*UserFollow) TableName() string {
	return "user_follow"
}
