package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/service"
	"project/utils"
	"strconv"
)

// FavoriteAction 点赞取消赞的操作
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	userToken, _ := utils.ParseToken(token)
	userId := userToken.ID

	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdStr)
	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.Atoi(actionTypeStr)

	err := service.FavoriteActions(userId, uint(videoId), actionType)
	//count, _ := service.GetFavoritesVideoCount(int64(videoId))
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "操作失败，err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg: "操作成功，actionType：" + actionTypeStr +
			" video_id is: " + videoIdStr})

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	_ = c.Query("token") //TODO 视频流客户端传递这个参数，用处Token续签、未登录的情况下查询关注返回false
	userIdStr := c.Query("user_id")
	userId, _ := strconv.Atoi(userIdStr)
	videoList, err := service.GetFavoriteList(uint(userId))
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "操作失败，err: " + err.Error(),
		})
	}
	c.JSON(http.StatusOK, models.FavoriteListResponse{
		FavoriteRes: models.Response{
			StatusCode: 0,
		},
		VideoResponse: videoList,
	})

}
