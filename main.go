package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"project/config"
	"project/controller"
	"project/dao/mongo"
	"project/dao/mysql"
	"project/dao/redis"
	"project/router"
	"project/service"
)

func main() {
	// 0. 开启协程负责消息模块
	go service.RunMessageServer()

	// 1. 加载配置
	if err := config.Init(); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
		return
	}

	// 2. 初始化MySQL
	if err := mysql.Init(config.Conf); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
		return
	}

	// 3. 初始化Redis
	if err := redis.Init(config.Conf); err != nil {
		fmt.Printf("Init redis failed, err:%v\n", err)
		return
	}

	// 4. 初始化Mongo
	if err := mongo.Init(config.Conf); err != nil {
		fmt.Printf("Init mongo failed, err:%v\n", err)
		return
	}

	// 准备数据
	controller.PrepareData()

	// 初始化gin引擎
	r := gin.Default()

	// 注册路由
	router.InitRouter(r)

	// 启动服务
	r.Run(fmt.Sprintf(":%d", config.Conf.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
