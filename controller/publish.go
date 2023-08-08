package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/service"
	"strconv"
)

// GetPublishList 每个用户的自己的发布列表
func GetPublishList(c *gin.Context) {
	//鉴权  TODO
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: int32(CodeInvalidParam),
			StatusMsg:  codeMsgMap[CodeInvalidParam]})
		return
	}
	videoList, err := service.GetPublishList(uint(userID))
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
