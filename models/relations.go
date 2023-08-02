package models

import "gorm.io/gorm"

type Relations struct {
	gorm.Model
	UserId      int64
	FollowingId *int64
	FollowedId  *int64
}

func (*Relations) TableName() string {
	return "relations"
}

func Follow() {

}

func UnFollow() {

}

func GetFollowList() {

}

func GetFollowerList() {

}
