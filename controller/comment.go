package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/service"
	"strconv"
)

func CommentAction(c *gin.Context) {

	// TODO 在鉴权中间件完成, 并保存user信息到ctx中, 不在这里重新实现
	//token := c.Query("token")
	//if user, exist := models.FindUserByToken(utils.DB, token); !exist {
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
	//} else {

	actionType := c.Query("action_type")
	commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	// FIXME JWT完成后取消评论
	// userId, _ := utils.GetCurrentUserID(c)
	userId := int64(1)

	switch actionType {
	// 新增评论
	case "1":
		content := c.Query("comment_text")
		data, err := service.AddComment(videoId, userId, content)
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
		data, err := service.DeleteComment(videoId, userId, commentId)
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
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)
	commentList, err := service.GetCommentList(videoId)
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
