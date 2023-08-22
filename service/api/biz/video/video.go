package video

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/constant/biz"
	"douyin/kitex_gen/video"
	"douyin/kitex_gen/video/videoservice"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var videoClient videoservice.Client

// 限制上传文件的最大大小 200MB 最小大小10KB
const maxFileSize = 200 * 1024 * 1024
const minFileSize = 10 * 1024

func init() {

	// Etcd 服务发现
	r, err := etcd.NewEtcdResolver([]string{constant.EtcdAddr})
	if err != nil {
		log.Fatal(err)
	}

	videoClient, err = videoservice.NewClient(
		constant.VideoServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.VideoServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func FeedList(ctx context.Context, c *app.RequestContext) {
	actorId, _ := c.Get(common.ContextUserIDKey)

	latestTime := c.Query("latest_time")
	unixTime, err := strconv.Atoi(latestTime)
	if latestTime == "" || latestTime < biz.MinTime {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	} else if err != nil || unixTime < 0 {
		zap.L().Error("Parse videoIdStr err:", zap.Error(err))
		c.JSON(http.StatusOK, err.Error())
		return
	}
	req := &video.VideoFeedRequest{
		UserId:     int32(actorId.(uint)),
		LatestTime: &latestTime,
	}

	resp, err := videoClient.VideoFeed(ctx, req)

	if err != nil {
		zap.L().Error("Get feed list from video client err.", zap.Error(err))
		c.JSON(http.StatusOK, video.VideoFeedResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Server internal error."),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetPublishList(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, video.PublishVideoListResponse{
			StatusCode: 0,
		})
		return
	}
	request := &video.PublishVideoListRequest{
		FromUserId: 0,
		ToUserId:   int32(userId),
	}
	result, err := videoClient.GetPublishVideoList(ctx, request)

	c.JSON(http.StatusOK, result)
}

func Publish(ctx context.Context, c *app.RequestContext) {
	zap.L().Info("Publish video")
	token := c.PostForm("token")
	title := c.PostForm("title")
	file, err := c.FormFile("data")

	if err != nil || token == "" || title == "" || file.Size == 0 {
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("参数不合法"),
		})
		return
	}
	userToken, err := common.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("token参数不合法"),
		})
		return
	}
	userId := userToken.ID
	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidFileType(ext) {
		msg := "无效的文件类型"
		zap.L().Error(msg)
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: 1,
			StatusMsg:  &msg})
		return
	}

	// 校验文件大小
	if file.Size > maxFileSize || file.Size < minFileSize {
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("文件过大或过小"),
		})
		return
	}

	// 打开上传的文件
	src, err := file.Open()
	defer src.Close()
	if err != nil {
		zap.L().Error("打开文件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, video.PublishVideoResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("打开文件失败"),
		})
		return
	}
	data, err := io.ReadAll(src)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("读取文件内容失败: %s", err.Error()))
		return
	}

	req := &video.PublishVideoRequest{
		UserId: int32(userId),
		Title:  title,
		Data:   data,
	}

	resp, err := videoClient.PublishVideo(ctx, req)

	if err != nil {
		zap.L().Error("PublishVideo err.", zap.Error(err))
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("PublishVideo error."),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
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
