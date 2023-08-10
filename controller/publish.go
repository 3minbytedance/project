package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"project/models"
	"project/service"
	"strconv"
	"project/utils"
)

// 限制上传文件的最大大小 200MB
const maxFileSize = 200 * 1024 * 1024

func Publish(c *gin.Context) {
	userId, err := utils.GetCurrentUserID(c)
	if err != nil {
	}
	title := c.Query("title")

	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 校验文件类型
	ext := filepath.Ext(file.Filename)
	if !isValidFileType(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// 校验文件大小
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the limit"})
		return
	}

	if err = service.UploadVideo(file); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: 400,
			StatusMsg:  "上传失败"})
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: int32(CodeSuccess),
		StatusMsg:  codeMsgMap[CodeSuccess]})

	// MQ 异步解耦,解决返回json阻塞 TODO
	service.GetVideoCover()
	service.StoreVideoAndImg()
}

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

// 校验文件类型是否为视频类型
func isValidFileType(fileExt string) bool {
	validExts := []string{".mp4", ".avi", ".mov"}
	for _, ext := range validExts {
		if fileExt == ext {
			return true
		}
	}
	return false
}