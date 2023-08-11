package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/service"
	"strconv"
	"time"
)

func Feed(c *gin.Context) {
	_ = c.Query("token") //TODO 视频流客户端传递这个参数做什么？
	latestTime := c.Query("latest_time")
	if latestTime == "" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	}

	videoList, err := service.GetFeedList(latestTime)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: int32(CodeInvalidParam),
			StatusMsg:  codeMsgMap[CodeInvalidParam],
		})
		return
	}
	c.JSON(http.StatusOK, models.VideoListResponse{
		Response: models.Response{
			StatusCode: int32(CodeSuccess),
			StatusMsg:  codeMsgMap[CodeSuccess],
		},
		VideoResponse: videoList,
	})
}
