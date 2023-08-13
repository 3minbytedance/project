package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"project/models"
	"project/service"
	"project/utils"
	"strconv"
	"strings"
)

// 限制上传文件的最大大小 200MB
const maxFileSize = 200 * 1024 * 1024
const minFileSize = 1 * 1024

// TODO
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	file, err := c.FormFile("data")
	if token == "" || title == "" || err != nil || file.Size == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: 400,
			StatusMsg:  "参数错"})
		return
	}
	userToken, _ := utils.ParseToken(token)
	userId := userToken.ID
	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidFileType(ext) {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: 400,
			StatusMsg:  "无效的文件类型"})
		return
	}

	// 校验文件大小
	if file.Size > maxFileSize || file.Size < minFileSize {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: 400,
			StatusMsg:  "文件过大或过小"})
		return
	}

	// 生成 UUID
	fileId := strings.Replace(uuid.New().String(), "-", "", -1)

	// 修改文件名
	videoFileName := fileId + ".mp4"

	//todo IO流优化待测，先用gin内置的
	err = c.SaveUploadedFile(file, "./public/"+videoFileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: 400,
			StatusMsg:  "上传失败"})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: int32(CodeSuccess),
		StatusMsg:  codeMsgMap[CodeSuccess]})

	// MQ 异步解耦,解决返回json阻塞 TODO

	imgName := service.GetVideoCover(videoFileName)
	service.StoreVideoAndImg(videoFileName, imgName, userId, title)
}

// GetPublishList 每个用户的自己的发布列表
func GetPublishList(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {

		c.JSON(http.StatusOK, models.Response{
			StatusCode: int32(CodeInvalidParam),
			StatusMsg:  codeMsgMap[CodeInvalidParam]})
		return
	}
	videoList := service.GetPublishList(uint(userID))
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
