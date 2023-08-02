package controller

import (
	"fmt"
	"os"
	"project/models"
	"project/utils"
	"time"
)

var ServerUrl = "https://" + os.Getenv("paas_url")
var LocalUrl = "http://loaclhost:8080"

var DemoVideos = []models.VideoRes{
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

var DemoComments = []models.Comment{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = models.User{
	Name:          "yyf",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
	//Avatar:          ServerUrl + "/public/avatar3.jpg",
	//BackgroundImage: ServerUrl + "/public/tx.jpeg",
	Avatar:          LocalUrl + "/public/avatar3.jpg",
	BackgroundImage: LocalUrl + "/public/tx.jpeg",
	Signature:       "这是个大帅逼",
	TotalFavorited:  99999,
	WorkCount:       20,
	FavoriteCount:   100,
}

func PrepareData() {
	// 建表
	table := utils.DB.Migrator().HasTable(&models.User{})
	if !table {
		err := utils.DB.AutoMigrate(&models.User{})
		if err != nil {
			fmt.Println("create user table failed.")
		}
	}
	table = utils.DB.Migrator().HasTable(&models.UserStates{})
	if !table {
		err := utils.DB.AutoMigrate(&models.UserStates{})
		if err != nil {
			fmt.Println("create user_states table failed.")
		}
	}
	table = utils.DB.Migrator().HasTable(&models.Comments{})
	if !table {
		err := utils.DB.AutoMigrate(&models.Comments{})
		if err != nil {
			fmt.Println("create comments table failed.")
		}
	}
	table = utils.DB.Migrator().HasTable(&models.Video{})
	if !table {
		err := utils.DB.AutoMigrate(&models.Video{})
		if err != nil {
			fmt.Println("create video table failed.")
		}
	}

	table = utils.DB.Migrator().HasTable(&models.Relations{})
	if !table {
		err := utils.DB.AutoMigrate(&models.Relations{})
		if err != nil {
			fmt.Println("create relations table failed.")
		}
	}

	// 新建数据
	videoId := 1
	if _, b := models.FindVideoByVideoId(utils.DB, videoId); !b {
		// 没数据的时候
		videos := []models.Video{
			{
				AuthorId:      1,
				PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
				CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
				FavoriteCount: 100,
				CommentCount:  100,
				IsFavorite:    false,
				PublishTime:   time.Now().Unix(),
			},
		}
		utils.DB.Model(&models.Video{}).Create(&videos)
	}
	if _, b := models.FindCommentsByVideoId(utils.DB, videoId); !b {
		// 没数据的时候
		comments := []models.Comments{
			{
				VideoId:    1,
				UserId:     2,
				Content:    "真棒",
				CreateTime: time.Now().Unix(),
			},
			{
				VideoId:    1,
				UserId:     2,
				Content:    "厉害了厉害了",
				CreateTime: time.Now().Unix(),
			},
		}
		utils.DB.Model(&models.Comments{}).Create(&comments)
	}
}
