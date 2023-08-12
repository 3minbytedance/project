package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"project/dao/mysql"
	"project/models"
	"project/utils"
	"strings"
)

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

func GetVideoCover(fileName string) string {
	// 生成图片 UUID
	imgId := uuid.New().String()
	// 修改文件名
	imgName := strings.Replace(imgId, "-", "", -1) + ".jpg"
	//调用ffmpeg 获取封面图
	utils.GetVideoFrame("./public/"+fileName, "./public/"+imgName)
	return imgName
}

// todo
func StoreVideoAndImg(videoName string, imgName string, authorId uint, title string) {
	// 视频存储到oss
	//if err := UploadToOSS("/dumpfile/"+fileName, fileName); err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//// 图片存储到oss
	//if err := UploadToOSS("/dumpfile/"+imgName, imgName); err != nil {
	//	log.Fatal(err)
	//	return
	//}
	fmt.Println(videoName, imgName, authorId, title)

	mysql.InsertVideo(videoName, imgName, authorId, title)
}

func GetPublishList(userID uint) ([]models.VideoResponse, bool) {
	videos, found := mysql.FindVideosByAuthorId(userID)
	if !found {
		return []models.VideoResponse{}, false
	}
	// 将查询结果转换为VideoResponse类型
	var videoResponses []models.VideoResponse
	for _, video := range videos {
		user, _ := GetUserInfoByUserId(userID)
		commentCount, _ := GetCommentCount(video.VideoId)
		videoResponse := models.VideoResponse{
			Id:            video.VideoId,
			Author:        user,
			PlayUrl:       video.VideoUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: 0, // TODO
			CommentCount:  commentCount,
			IsFavorite:    isUserFavorite(111, video.VideoId), // TODO  userId,videoID
			Title:         video.Title,
		}
		videoResponses = append(videoResponses, videoResponse)
	}

	return videoResponses, true
}

func GetFeedList(latestTime string) ([]models.VideoResponse, int64, error) {
	videos := mysql.GetLatestVideos(latestTime)
	if len(videos) == 0 {
		return []models.VideoResponse{}, 0, errors.New("no videos")
	}
	// 将查询结果转换为VideoResponse类型
	var videoResponses []models.VideoResponse
	for _, video := range videos {
		user, _ := GetUserInfoByUserId(video.AuthorId)
		commentCount, _ := GetCommentCount(video.VideoId)
		videoResponse := models.VideoResponse{
			Id:            video.VideoId,
			Author:        user,
			PlayUrl:       video.VideoUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: 0, // TODO
			CommentCount:  commentCount,
			IsFavorite:    isUserFavorite(111, video.VideoId), // TODO  userId,videoID
			Title:         video.Title,
		}

		videoResponses = append(videoResponses, videoResponse)
	}
	//本次返回的视频中，发布最早的时间
	nextTime := videos[len(videos)-1].CreatedAt.Unix()

	return videoResponses, nextTime, nil
}

func getFavoriteCount(uint) uint { return 1 }

func isUserFavorite(uint, uint) bool { return true }
