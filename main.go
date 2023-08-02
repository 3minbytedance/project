package main

import (
	"github.com/gin-gonic/gin"
	"project/controller"
	"project/router"
	"project/service"
	"project/utils"
)

func main() {
	go service.RunMessageServer()
	utils.InitConfig()
	utils.InitMysql()
	controller.PrepareData()

	r := gin.Default()
	router.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
