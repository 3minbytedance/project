package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/service"
	"strconv"
)

// FavoriteAction 点赞取消赞的操作
func FavoriteAction(c *gin.Context) {
	// 鉴权
	//token := c.Query("token")
	userIdStr := c.Query("user_id")
	userId, _ := strconv.Atoi(userIdStr)
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdStr)
	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.Atoi(actionTypeStr)
	err := service.FavoriteActions(int64(userId), int64(videoId), actionType)
	count, _ := service.GetFavoritesVideoCount(int64(videoId))
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "操作失败，err: " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 0,
			StatusMsg: "操作成功，actionType：" + actionTypeStr +
				" video_id is: " + videoIdStr + " count is: " + strconv.Itoa(count)})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, _ := strconv.Atoi(userIdStr)
	list, err := service.GetFavoriteList(int64(userId))
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "操作失败，err: " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, service.FavoriteListResponse{
			FavoriteRes: models.Response{
				StatusCode: 0,
			},
			VideoList: list,
		})
	}
}
