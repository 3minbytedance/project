package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/models"
	"project/service"
	"project/utils"
	"strconv"
	"time"
)

func Feed(c *gin.Context) {
	token := c.Query("token") //TODO 视频流客户端传递这个参数，用处Token续签、未登录的情况下查询关注返回false
	latestTime := c.Query("latest_time")
	unixTime, err := strconv.Atoi(latestTime)

	// 先判断是否为null，如果不是再判断参数是否合法
	if latestTime == "" || latestTime == "0" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	} else if err != nil || unixTime < 0 {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: int32(CodeInvalidParam),
			StatusMsg:  CodeInvalidParam.Msg(),
		})
		return
	}
	var userId uint
	isLogged := false
	//todo 改为如果token在redis中查到
	if token != "" {
		userToken, _ := utils.ParseToken(token)
		userId = userToken.ID
		isLogged = true
	}
	videoList, nextTime, err := service.GetFeedList(latestTime, isLogged, userId)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, models.Response{
			StatusCode: int32(CodeSuccess),
			StatusMsg:  CodeInvalidParam.Msg(),
		})
		return
	}
	c.JSON(http.StatusOK, models.FeedListResponse{

		Response: models.Response{
			StatusCode: int32(CodeSuccess),
			StatusMsg:  CodeSuccess.Msg(),
		},
		NextTime:      nextTime,
		VideoResponse: videoList,
	})
}
