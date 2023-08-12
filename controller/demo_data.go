package controller

import (
	"fmt"
	"os"
	"project/dao/mysql"
	"project/models"
)

var ServerUrl = "https://" + os.Getenv("paas_url")
var LocalUrl = "http://loaclhost:8080"


var DemoComments = []models.CommentResponse{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

// TODO
var DemoUser = models.UserResponse{
	Name:            "yyf",
	FollowCount:     0,
	FollowerCount:   0,
	IsFollow:        false,
	Avatar:          LocalUrl + "/public/avatar3.jpg",
	BackgroundImage: LocalUrl + "/public/tx.jpeg",
	Signature:       "这是个大帅逼",
	TotalFavorited:  "99999",
	WorkCount:       20,
	FavoriteCount:   100,
}

func PrepareData() {
	// 建表
	table := mysql.DB.Migrator().HasTable(&models.User{})
	if !table {
		err := mysql.DB.AutoMigrate(&models.User{})
		if err != nil {
			fmt.Println("create user table failed.")
		}
	}
	table = mysql.DB.Migrator().HasTable(&models.Comment{})
	if !table {
		err := mysql.DB.AutoMigrate(&models.Comment{})
		if err != nil {
			fmt.Println("create comments table failed.")
		}
	}
	table = mysql.DB.Migrator().HasTable(&models.Video{})
	if !table {
		err := mysql.DB.AutoMigrate(&models.Video{})
		if err != nil {
			fmt.Println("create video table failed.")
		}
	}

	table = mysql.DB.Migrator().HasTable(&models.Relations{})
	if !table {
		err := mysql.DB.AutoMigrate(&models.Relations{})
		if err != nil {
			fmt.Println("create relations table failed.")
		}
	}
}
