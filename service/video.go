package service

import (
	"context"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"project/dao/mysql"
	"project/models"
)

func uplaod() {

	// MQ 异步解耦 TODO

	// 存储到oss

	//写入到本地，创建临时文件

	//将封面图上传到oss

}

// UploadVideo
func UploadVideo(localPath string, remotePath string) error {
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

func GetPublishList(userID uint) ([]models.VideoResponse, error) {
	videos, err := mysql.FindVideosByAuthorId(userID)
	if err != nil {
		return nil, err
	}
	// 将查询结果转换为VideoResponse类型
	var videoResponses []models.VideoResponse
	for _, video := range videos {
		videoResponse := models.VideoResponse{
			Id:            video.VideoId,
			Author:        models.User{}, //TODO
			PlayUrl:       video.VideoUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: getFavoriteCount(video.VideoId),    // TODO
			CommentCount:  getCommentCount(video.VideoId),     // TODO
			IsFavorite:    isUserFavorite(111, video.VideoId), // TODO  userId,videoID
		}
		videoResponses = append(videoResponses, videoResponse)
	}

	return videoResponses, nil
}

func getFavoriteCount(uint) uint { return 1 }

func getCommentCount(uint) uint { return 1 }

func isUserFavorite(uint, uint) bool { return true }
