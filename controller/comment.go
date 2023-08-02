package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/utils"
	"strconv"
	"time"
)

type CommentListResponse struct {
	models.Response
	CommentList []models.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	models.Response
	Comment models.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	if user, exist := models.FindUserByToken(utils.DB, token); !exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
	} else {
		switch actionType {
		case "1":
			text := c.Query("comment_text")
			commData := models.Comments{
				VideoId:    int64(videoId),
				UserId:     int64(user.ID),
				Content:    text,
				CreateTime: time.Now().Unix(),
			}
			comment := models.Comment{
				User:       user,
				Content:    text,
				CreateDate: models.CommentTime(commData.CreateTime, time.Now().Unix()),
			}

			// 更新评论信息
			utils.DB.Model(models.Comments{}).Create(&commData)
			// 更新视频信息
			video, b := models.FindVideoByVideoId(utils.DB, videoId)
			if !b {
				fmt.Println("未找到对应的视频")
			} else {
				num := video.CommentCount + 1
				utils.DB.Model(&video).Update("comment_count", strconv.Itoa(int(num)))
			}
			c.JSON(http.StatusOK, CommentActionResponse{Response: models.Response{StatusCode: 0},
				Comment: comment})
			return
		case "2": // 删除评论的操作，需要拿到commentId才能删，还要改对应的数据库

		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "没操作"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdStr)
	comments := models.GetComments(videoId)
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    models.Response{StatusCode: 0},
		CommentList: comments,
	})
}
