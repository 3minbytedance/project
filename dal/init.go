package dal

import (
	"douyin/config"
	"douyin/dal/mongo"
	"douyin/dal/mysql"
)

func Init(appConfig *config.AppConfig) error {
	err := mongo.Init(appConfig)
	if err != nil {
		return err
	}
	err = mysql.Init(appConfig)
	if err != nil {
		return err
	}
	return nil
}
