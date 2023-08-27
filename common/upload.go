package common

import (
	"context"
	"douyin/constant/biz"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func CreateDirectoryIfNotExist() error {
	if _, err := os.Stat(biz.FileLocalPath); os.IsNotExist(err) {
		// 创建文件夹
		err = os.MkdirAll(biz.FileLocalPath, 0777)
		if err != nil {
			zap.Error(err)
			return err
		}
	}
	return nil
}

func UploadToOSS(localPath string, remotePath string) error {
	c := getClient()

	_, _, err := c.Object.Upload(context.Background(), remotePath, localPath, nil)
	if err != nil {
		return err
	}
	return nil
}

func getClient() *cos.Client {
	u, _ := url.Parse(biz.OSS)
	cu, _ := url.Parse(biz.CuOSS)
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     biz.SecretId,
			SecretKey:    biz.SecretKey,
			SessionToken: biz.SessionToken,
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    true,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	return c
}

func GetVideoCover(videoName string) (string, error) {
	// 生成图片 UUID
	imgId := uuid.New().String()
	// 修改文件名
	imgName := strings.Replace(imgId, "-", "", -1) + ".jpg"
	//调用oss 获取封面图
	err := postSnapShot(videoName, imgName)
	if err != nil {
		return "", err
	}
	return imgName, nil
}

func postSnapShot(videoName string, imgName string) error {
	c := getClient()
	PostSnapshotOpt := &cos.PostSnapshotOptions{
		Input: &cos.JobInput{
			Object: videoName,
		},
		Time:   "1",
		Width:  720,
		Height: 1280,
		Format: "jpg",
		Output: &cos.JobOutput{
			Region: "ap-nanjing",
			Bucket: "tiktok-1319971229",
			Object: imgName,
		},
	}
	_, _, err := c.CI.PostSnapshot(context.Background(), PostSnapshotOpt)
	if err != nil {
		return err
	}
	return nil
}
