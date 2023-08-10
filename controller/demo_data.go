package controller

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"project/dao/mysql"
	"project/models"
	"project/service"
	"time"
)

var ServerUrl = "https://" + os.Getenv("paas_url")
var LocalUrl = "http://loaclhost:8080"

var DemoVideos = []models.VideoResponse{
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoComments = []models.CommentResponse{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

// TODO
var DemoUser = models.UserInfo{
	Name:            "yyf",
	FollowCount:     0,
	FollowerCount:   0,
	IsFollow:        false,
	Avatar:          LocalUrl + "/public/avatar3.jpg",
	BackgroundImage: LocalUrl + "/public/tx.jpeg",
	Signature:       "这是个大帅逼",
	TotalFavorited:  99999,
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
	table = mysql.DB.Migrator().HasTable(&models.UserStates{})
	if !table {
		err := mysql.DB.AutoMigrate(&models.UserStates{})
		if err != nil {
			fmt.Println("create user_states table failed.")
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

	// 新建数据
	videoId := int64(1)
	if _, b := mysql.FindVideoByVideoId(videoId); !b {
		// 没数据的时候
		videos := []models.Video{
			{
				AuthorId:  1,
				VideoUrl:  "https://www.w3schools.com/html/movie.mp4",
				CoverUrl:  "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
				Title:     "hello world",
				CreatedAt: time.Now(),
			},
		}
		mysql.DB.Model(&models.Video{}).Create(&videos)
	}
	if _, err := service.GetCommentList(videoId); err == nil {
		// 没数据的时候
		comments := []models.Comment{
			{
				VideoId: 1,
				UserId:  2,
				Content: "真棒",
				Model:   gorm.Model{CreatedAt: time.Now()},
			},
			{
				VideoId: 1,
				UserId:  2,
				Content: "厉害了厉害了",
				Model:   gorm.Model{CreatedAt: time.Now()},
			},
		}

		mysql.DB.Model(&models.Comment{}).Create(&comments)
	}
}
