package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/models"
	"project/service"
	"strconv"
)

// Publish check token then save upload file to public directory
//func Publish(c *gin.Context) {
//	//TODO
//
//	token := c.Query("token")
//
//	if user, exist := mysql.FindUserByToken(token); !exist {
//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
//		return
//	} else {
//		data, err := c.FormFile("data")
//		if err != nil {
//			c.JSON(http.StatusOK, models.Response{
//				StatusCode: 1,
//				StatusMsg:  err.Error(),
//			})
//			return
//		}
//
//		filename := filepath.Base(data.Filename)
//		finalName := fmt.Sprintf("%d_%s", user.ID, filename)
//		saveFile := filepath.Join("./public/", finalName)
//		if err := c.SaveUploadedFile(data, saveFile); err != nil {
//			c.JSON(http.StatusOK, models.Response{
//				StatusCode: 1,
//				StatusMsg:  err.Error(),
//			})
//			return
//		}
//		// 更新视频信息
//		video := models.Video{
//			AuthorId:    int64(user.ID),
//			PlayUrl:     saveFile,
//			PublishTime: time.Now().Unix(),
//		}
//		mysql.DB.Model(models.Video{}).Create(&video)
//		// 视频封面如何获取？用户上传（自定义）+默认生成
//
//		// 更新用户作品数量
//		mysql.DB.Model(&user).Update("work_count", user.WorkCount+1)
//		c.JSON(http.StatusOK, models.Response{
//			StatusCode: 0,
//			StatusMsg:  finalName + " uploaded successfully",
//		})
//	}
//}

// PublishList 每个用户的自己的发布列表
func PublishList(c *gin.Context) {
	//鉴权  TODO
	id := c.Query("user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: int32(CodeInvalidParam),
			StatusMsg:  codeMsgMap[CodeInvalidParam]})
		return
	}
	_, err = service.GetPublishList(uint(userID))
	if err != nil {
		return
	}

	//TODO
}
