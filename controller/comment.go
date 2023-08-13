package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/service"
	"project/utils"
	"strconv"
)

func CommentAction(c *gin.Context) {
	// TODO 在鉴权中间件完成, 并保存userId信息到ctx中, 不在这里重新实现
	//token := c.Query("token")
	//if user, exist := models.FindUserByToken(utils.DB, token); !exist {
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
	//} else {

	actionType := c.Query("action_type")
	commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	userId, _ := utils.GetCurrentUserID(c)

	switch actionType {
	// 新增评论
	case "1":
		content := c.Query("comment_text")
		data, err := service.AddComment(uint(videoId), uint(userId), content)

		if err != nil {
			c.JSON(http.StatusOK, models.CommentActionResponse{
				Response: models.Response{
					StatusCode: int32(CodeServerBusy),
					StatusMsg:  CodeServerBusy.Msg() + err.Error(),
				},
				Comment: data,
			})
		}
		c.JSON(http.StatusOK, models.CommentActionResponse{
			Response: models.Response{
				StatusCode: int32(CodeSuccess),
				StatusMsg:  CodeSuccess.Msg(),
			},
			Comment: data,
		})
		return

	// 删除评论
	case "2":
		data, err := service.DeleteComment(uint(videoId), uint(userId), uint(commentId))
		if err != nil {
			c.JSON(http.StatusOK, models.CommentActionResponse{
				Response: models.Response{
					StatusCode: int32(CodeServerBusy),
					StatusMsg:  CodeServerBusy.Msg() + err.Error(),
				},
				Comment: data,
			})
		}
		c.JSON(http.StatusOK, models.CommentActionResponse{
			Response: models.Response{
				StatusCode: int32(CodeSuccess),
				StatusMsg:  CodeSuccess.Msg(),
			},
			Comment: data})
		return
	default:
		c.JSON(http.StatusOK, models.Response{
			StatusCode: int32(CodeInvalidParam),
			StatusMsg:  CodeInvalidParam.Msg(),
		})
	}
}

func CommentList(c *gin.Context) {
	token := c.Query("token") //TODO 视频流客户端传递这个参数，用处Token续签、未登录的情况下查询关注返回false
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)
	userToken, _ := utils.ParseToken(token)
	userId := userToken.ID
	isLogged := false
	//todo 改为如果token在redis中查到
	if token != "" {
		isLogged = true
	}
	commentList, err := service.GetCommentList(uint(videoId), isLogged, userId)
	if err != nil {
		c.JSON(http.StatusOK, models.CommentListResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  "Found comments failed:" + err.Error(),
			},
			CommentList: nil,
		})
	}
	c.JSON(http.StatusOK, models.CommentListResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "Found comments successfully.",
		},
		CommentList: commentList,
	})
}
