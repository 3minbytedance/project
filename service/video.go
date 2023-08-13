package service

import (
	"bufio"
	"context"
	"errors"
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

const (
	fileLocalPath = "./public/"
	oss           = "https://tiktok-1319971229.cos.ap-nanjing.myqcloud.com/"
	SecretId      = "AKIDFKMQPakpcN6tkV9oJg6PanzAGC0hGkCZ"
	SecretKey     = "MWXXLzQlutgMtLl5HH9pPp5CB0cfcMxR"
	SessionToken  = "SECRETTOKEN"
)

// UploadToOSS  上传至腾讯OSS
func UploadToOSS(localPath string, remotePath string) error {
	reqUrl := oss
	u, _ := url.Parse(reqUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     SecretId,
			SecretKey:    SecretKey,
			SessionToken: SessionToken,
		},
	})

	_, _, err := c.Object.Upload(context.Background(), remotePath, localPath, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoCover(fileName string) string {
	// 生成图片 UUID
	imgId := uuid.New().String()
	// 修改文件名
	imgName := strings.Replace(imgId, "-", "", -1) + ".jpg"
	//调用ffmpeg 获取封面图
	utils.GetVideoFrame(fileLocalPath+fileName, fileLocalPath+imgName)
	return imgName
}


func StoreVideoAndImg(videoName string, imgName string, authorId uint, title string) {
	//视频存储到oss
	go func() {
		if err := UploadToOSS(fileLocalPath+videoName, videoName); err != nil {
			log.Fatal(err)
			return
		}
	}()

	// 图片存储到oss
	go func() {
		if err := UploadToOSS(fileLocalPath+imgName, imgName); err != nil {
			log.Fatal(err)
			return
		}
	}()

	go func() {
		mysql.InsertVideo(videoName, imgName, authorId, title)
		if !redis.IsExistUserField(authorId, redis.WorkCountField){
			cnt := mysql.FindWorkCountsByAuthorId(authorId)
			err := redis.SetWorkCountByUserId(authorId, cnt)
			if err != nil {
				log.Println("redis更新评论数失败", err)
				return
			}
			return
		}
		err := redis.IncrementWorkCountByUserId(authorId)
		if err != nil {
			log.Println(err)
		}
	}()
}

// GetWorkCount 获得作品数
func GetWorkCount(userId uint) (int64, error) {
	// 从redis中获取作品数
	// 1. 缓存中有数据, 直接返回
	if redis.IsExistUserField(userId, redis.WorkCountField) {
		workCount, err := redis.GetWorkCountByUserId(userId)
		if err != nil {
			log.Println("从redis中获取作品数失败：", err)
		}
		return workCount, nil
	}

	// 2. 缓存中没有数据，从数据库中获取
	workCount := mysql.FindWorkCountsByAuthorId(userId)
	log.Println("从数据库中获取作品数成功：", workCount)
	// 将作品数写入redis
	go func() {
		err := redis.SetWorkCountByUserId(userId, workCount)
		if err != nil {
			log.Println("将作品数写入redis失败：", err)
		}
	}()
	return workCount, nil
}

func GetPublishList(userID uint) (videoResponses []models.VideoResponse) {
	videos, found := mysql.FindVideosByAuthorId(userID)
	if !found {
		return []models.VideoResponse{}
	}
	// 将查询结果转换为VideoResponse类型
	videoResponses = make([]models.VideoResponse, 0, len(videos))
	for _, video := range videos {
		user, _ := GetUserInfoByUserId(userID)
		commentCount, _ := GetCommentCount(video.ID)
		videoResponse := models.VideoResponse{
			ID:            video.ID,
			Author:        user,
			PlayUrl:       oss + video.VideoUrl,
			CoverUrl:      oss + video.CoverUrl,
			FavoriteCount: 0, // TODO
			CommentCount:  commentCount,
			IsFavorite:    isUserFavorite(111, video.ID), // TODO  userId,videoID
			Title:         video.Title,
		}
		videoResponses = append(videoResponses, videoResponse)
	}

	return videoResponses
}

func GetFeedList(latestTime string) ([]models.VideoResponse, int64, error) {
	videos := mysql.GetLatestVideos(latestTime)
	if len(videos) == 0 {
		return []models.VideoResponse{}, 0, errors.New("no videos")
	}
	// 将查询结果转换为VideoResponse类型
	videoResponses := make([]models.VideoResponse, 0, len(videos))
	for _, video := range videos {
		user, _ := GetUserInfoByUserId(video.AuthorId)
		commentCount, _ := GetCommentCount(video.ID)
		videoResponse := models.VideoResponse{
			ID:            video.ID,
			Author:        user,
			PlayUrl:       oss + video.VideoUrl,
			CoverUrl:      oss + video.CoverUrl,
			FavoriteCount: 0, // TODO
			CommentCount:  commentCount,
			IsFavorite:    isUserFavorite(111, video.ID), // TODO  userId,videoID
			Title:         video.Title,
		}

		videoResponses = append(videoResponses, videoResponse)
	}
	//todo 客户端刷新也是用这个时间？
	//本次返回的视频中，发布最早的时间
	nextTime := videos[len(videos)-1].CreatedAt
	return videoResponses, nextTime, nil
}

func getFavoriteCount(uint) uint { return 1 }

func isUserFavorite(uint, uint) bool { return false }

// todo io优化，待测
func UploadIOVideo(file *multipart.FileHeader) (string, error) {
	// 生成 UUID
	fileId := strings.Replace(uuid.New().String(), "-", "", -1)

	// 修改文件名
	fileName := fileId + ".mp4"

	// 创建临时文件
	tempFile, err := createTempFile(fileName)
	if err != nil {
		return "", err
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 创建缓冲写入器
	dest := bufio.NewWriter(tempFile)

	// 将上传的文件内容写入临时文件
	_, err = io.Copy(dest, src)
	if err != nil {
		return "", err
	}

	// 清空缓冲区并确保文件已写入磁盘
	if err = dest.Flush(); err != nil {
		return "", err
	}
	return fileName, nil
}

func createTempFile(fileName string) (*os.File, error) {
	tempDir := fileLocalPath // 临时文件夹路径

	// 创建临时文件夹（如果不存在）
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err := os.Mkdir(tempDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// 在临时文件夹中创建临时文件
	tmpfile, err := os.Create(tempDir + fileName)
	if err != nil {
		return nil, err
	}

	return tmpfile, nil
}
