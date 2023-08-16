package mw

import (
	"douyin/config"
	"douyin/mw/kafka"
	"douyin/mw/redis"
)

func Init(appConfig *config.AppConfig) error {
	err := kafka.Init(appConfig)
	if err != nil {
		return err
	}

	err = redis.Init(appConfig)
	if err != nil {
		return err
	}
	return nil
}
