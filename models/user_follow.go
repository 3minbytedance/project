package models

type User_follow struct {
	// 用户的关注信息
	User_Id   uint `gorm:"primaryKey;not null"` // 用户id
	Follow_Id uint `gorm:"not null"`            // 关注用户id
}

func (*User_follow) TableName() string {
	return "user_follow"
}
