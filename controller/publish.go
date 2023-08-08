package controller

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"project/models"
	"project/service"
	"strconv"
	"strings"
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

func Publish(c *gin.Context) {
	//TODO 鉴权

	//TODO 检测文件类型

	file, err := c.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	// 生成 UUID
	id := uuid.New().String()

	// 修改文件名
	fileName := strings.Replace(id, "-", "", -1) + ".mp4"

	done := make(chan struct{})

	// 处理文件上传的并发函数
	go func(file *multipart.FileHeader) {

		src, err := file.Open()
		if err != nil {
			log.Println("Failed to open source file:", err)
			return
		}
		defer src.Close()

		dest, err := os.Create(fileName)
		if err != nil {
			log.Println("Failed to create destination file:", err)
			return
		}
		defer dest.Close()

		reader := bufio.NewReader(src)
		_, err = io.Copy(dest, reader)
		if err != nil {
			log.Println("Failed to copy file:", err)
			return
		}
		done <- struct{}{}
	}(file)

	<-done

	c.JSON(http.StatusOK, models.Response{
		StatusCode: int32(CodeSuccess),
		StatusMsg:  codeMsgMap[CodeSuccess]})

	// MQ 异步解耦 TODO

	//调用ffmpeg 获取封面图

	// 存储到oss
	service.UploadVideo(fileName, fileName)

	// 存储到oss
	service.UploadVideo(fileName, fileName)

	// 将视频及图片uuid写入数据库
}
