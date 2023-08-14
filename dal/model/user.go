package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uint   `gorm:"primarykey"`
	Username        string `gorm:"uniqueIndex;size:32"` // 用户名称
	Password        string // 用户密码
	Avatar          string // 用户头像
	BackgroundImage string // 用户个人页顶部大图
	Signature       string `default:"默认签名"` // 个人简介
	Salt            string // 加密盐
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
