package service

import (
	"bufio"
	"context"
	"github.com/google/uuid"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"project/dao/mysql"
	"project/dao/redis"
	"project/models"
	"project/utils"
	"strings"
)

func UploadVideo(file *multipart.FileHeader) error {
	// 生成 UUID
	id := uuid.New().String()

	// 修改文件名
	fileName := strings.Replace(id, "-", "", -1) + ".mp4"

	// 创建临时文件
	tmpfile, err := createTempFile(fileName)
	if err != nil {
		return err
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建缓冲写入器
	writer := bufio.NewWriter(tmpfile)

	// 将上传的文件内容写入临时文件
	_, err = io.Copy(writer, src)
	if err != nil {
		return err
	}

	// 清空缓冲区并确保文件已写入磁盘
	if err = writer.Flush(); err != nil {
		return err
	}
	return nil
}

func createTempFile(fileName string) (*os.File, error) {
	tmpDir := "/dumpfile" // 临时文件夹路径

	// 创建临时文件夹（如果不存在）
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err := os.Mkdir(tmpDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// 在临时文件夹中创建临时文件
	tmpfile, err := os.CreateTemp(tmpDir, fileName)
	if err != nil {
		return nil, err
	}

	return tmpfile, nil
}

// UploadToOSS  上传至腾讯OSS
func UploadToOSS(localPath string, remotePath string) error {
	req_url := "https://tiktok-1319971229.cos.ap-nanjing.myqcloud.com"
	u, _ := url.Parse(req_url)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     "AKIDFKMQPakpcN6tkV9oJg6PanzAGC0hGkCZ",
			SecretKey:    "MWXXLzQlutgMtLl5HH9pPp5CB0cfcMxR",
			SessionToken: "SECRETTOKEN",
		},
	})

	// 通过文件流上传对象
	fd, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = c.Object.Put(context.Background(), remotePath, fd, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoCover(fileName string) {
	// 生成图片 UUID
	imgId := uuid.New().String()
	// 修改文件名
	imgName := strings.Replace(imgId, "-", "", -1) + ".jpg"
	//调用ffmpeg 获取封面图
	utils.GetVideoFrame("/dumpfile/"+fileName, "/dumpfile/"+imgName)
}

func StoreVideoAndImg(videoUrl string, coverUrl string, authorID uint, title string) {
	// 视频存储到oss
	if err := UploadToOSS("/dumpfile/"+fileName, fileName); err != nil {
		log.Fatal(err)
		return
	}

	// 图片存储到oss
	if err := UploadToOSS("/dumpfile/"+imgName, imgName); err != nil {
		log.Fatal(err)
		return
	}

	mysql.InsertVideo(videoUrl, coverUrl, authorID, title)
}

func GetPublishList(userID uint) ([]models.VideoResponse, error) {
	videos, err := mysql.FindVideosByAuthorId(userID)
	if err != nil {
		return nil, err
	}
	// 将查询结果转换为VideoResponse类型
	var videoResponses []models.VideoResponse
	for _, video := range videos {
		user, _ := mysql.FindUserInfoByUserId(userID)
		commentCount, _ := redis.GetCommentCountByVideoId(video.VideoId)
		videoResponse := models.VideoResponse{
			Id:            video.VideoId,
			Author:        user,
			PlayUrl:       video.VideoUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: 0,                                  // TODO
			CommentCount:  commentCount,                       // TODO
			IsFavorite:    isUserFavorite(111, video.VideoId), // TODO  userId,videoID
		}
		videoResponses = append(videoResponses, videoResponse)
	}

	return videoResponses, nil
}

func getFavoriteCount(uint) uint { return 1 }

func getCommentCount(uint) uint { return 1 }

func isUserFavorite(uint, uint) bool { return true }
