package service

import (
	"context"
	"github.com/spf13/viper"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"project/dao/mysql"
	"project/models"
)

func Upload_video(name string, path string) error {
	req_url := viper.GetString("oss.tencent")
	u, _ := url.Parse(req_url)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     "SECRETID",
			SecretKey:    "SECRETKEY",
			SessionToken: "SECRETTOKEN",
		},
	})

	// 通过文件流上传对象
	fd, err := os.Open("./test")
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = c.Object.Put(context.Background(), name, fd, nil)
	if err != nil {
		return err
	}
	return nil
}

func Download_video(name string) (interface{}, error) {
	req_url := viper.GetString("oss.tencent")
	u, _ := url.Parse(req_url)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     "SECRETID",
			SecretKey:    "SECRETKEY",
			SessionToken: "SECRETTOKEN",
		},
	})

	resp, err := c.Object.Get(context.Background(), name, nil)
	if err != nil {
		return nil, err
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bs, nil
}

func GetPublishList(userID uint64) ([]models.VideoResponse, error) {
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
			FavoriteCount: getFavoriteCount(video.VideoId),       // TODO
			CommentCount:  getCommentCount(video.VideoId),        // TODO
			IsFavorite:    isUserFavorite(userId, video.VideoId), // TODO
		}
		videoResponses = append(videoResponses, videoResponse)
	}

	return videoResponses, nil
}

func getFavoriteCount(uint) uint { return 1 }

func getCommentCount(uint) uint { return 1 }

func isUserFavorite(uint, uint) bool { return true }
