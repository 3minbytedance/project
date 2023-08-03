package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"project/dao/mysql"
	"project/models"
	"time"
)

type VideoListResponse struct {
	models.Response
	VideoList []models.VideoRes `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.Query("token")

	if user, exist := models.FindUserByToken(mysql.DB, token); !exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
		return
	} else {
		data, err := c.FormFile("data")
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}

		filename := filepath.Base(data.Filename)
		finalName := fmt.Sprintf("%d_%s", user.ID, filename)
		saveFile := filepath.Join("./public/", finalName)
		if err := c.SaveUploadedFile(data, saveFile); err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		// 更新视频信息
		video := models.Video{
			AuthorId:    int64(user.ID),
			PlayUrl:     saveFile,
			PublishTime: time.Now().Unix(),
		}
		mysql.DB.Model(models.Video{}).Create(&video)
		// 视频封面如何获取？用户上传（自定义）+默认生成

		// 更新用户作品数量
		mysql.DB.Model(&user).Update("work_count", user.WorkCount+1)
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully",
		})
	}
}

// PublishList 每个用户的自己的发布列表
func PublishList(c *gin.Context) {
	token := c.Query("token")

	if user, exist := models.FindUserByToken(mysql.DB, token); !exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
		return
	} else {
		videos, b := models.FindVideosByAuthor(mysql.DB, int(user.ID))
		if !b {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户未发布过作品"})
			return
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: models.Response{
				StatusCode: 0,
			},
			VideoList: videos,
		})
	}
}
