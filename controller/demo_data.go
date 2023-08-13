package controller

import (
	"fmt"
	"project/dao/mysql"
	"project/models"
)

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

	table = mysql.DB.Migrator().HasTable(&models.UserFollow{})
	if !table {
		err := mysql.DB.AutoMigrate(&models.UserFollow{})
		if err != nil {
			fmt.Println("create relations table failed.")
		}
	}
	table = mysql.DB.Migrator().HasTable(&models.Favorite{})
	if !table {
		err := mysql.DB.AutoMigrate(&models.Favorite{})
		if err != nil {
			fmt.Println("create favorite table failed.")
		}
	}

}
