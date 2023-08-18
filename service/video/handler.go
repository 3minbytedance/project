package main

import (
	"context"
	"douyin/common"
	"douyin/constant/biz"
	"douyin/dal/mysql"
	video "douyin/kitex_gen/video"
	"douyin/mw/redis"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/google/uuid"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// VideoFeed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoFeed(ctx context.Context, request *video.VideoFeedRequest) (resp *video.VideoFeedResponse, err error) {
	videos := mysql.GetLatestVideos(request.GetLatestTime())
	if len(videos) == 0 {
		zap.L().Error("根据LatestTime取视频失败", zap.Error(err))
		return &video.VideoFeedResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取视频列表失败"),
			VideoList:  nil,
		}, err
	}

	return
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, request *video.PublishVideoRequest) (resp *video.PublishVideoResponse, err error) {
	// TODO: Your code here...
	return
}

// GetPublishVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishVideoList(ctx context.Context, request *video.PublishVideoListRequest) (resp *video.PublishVideoListResponse, err error) {
	// TODO: Your code here...
	return
}

func StoreVideoAndImg(videoName string, imgName string, authorId uint, title string) {
	//视频存储到oss
	go func() {
		if err := UploadToOSS(biz.FileLocalPath+videoName, videoName); err != nil {
			log.Fatal(err)
			return
		}
	}()

	// 图片存储到oss
	go func() {
		if err := UploadToOSS(biz.FileLocalPath+imgName, imgName); err != nil {
			log.Fatal(err)
			return
		}
	}()

	go func() {
		mysql.InsertVideo(videoName, imgName, authorId, title)
		if !redis.IsExistUserField(authorId, redis.WorkCountField) {
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

// UploadToOSS  上传至腾讯OSS
func UploadToOSS(localPath string, remotePath string) error {
	reqUrl := biz.OSS
	u, _ := url.Parse(reqUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     biz.SecretId,
			SecretKey:    biz.SecretKey,
			SessionToken: biz.SessionToken,
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
	common.GetVideoFrame(biz.FileLocalPath+fileName, biz.FileLocalPath+imgName)
	return imgName
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
