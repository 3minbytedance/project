package dal

import (
	"douyin/dal/mongo"
	"douyin/dal/mysql"
)

func Init() {
	mongo.Init()
	mysql.Init()
}
