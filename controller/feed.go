package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Feed(c *gin.Context) {
	//TODO 鉴权
	_ = c.Query("token")
	latestTime := c.Query("latest_time")
	if latestTime == "" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	}
}
