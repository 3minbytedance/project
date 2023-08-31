package video

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/kitex_gen/video"
	"douyin/kitex_gen/video/videoservice"
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

// 限制上传文件的最大大小 50MB 最小大小10KB
const maxFileSize = 50 * 1024 * 1024
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
		client.WithMuxConnection(2),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func FeedList(ctx context.Context, c *app.RequestContext) {
	actorId, _ := common.GetCurrentUserID(c)
	zap.L().Info("FeedList", zap.Uint("actorID", actorId))
	latestTime := c.Query("latest_time")
	t, err := strconv.Atoi(latestTime)
	if latestTime == "" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	} else if err != nil || t < 0 {
		zap.L().Info("Parse videoIdStr err:", zap.Error(err))
		c.JSON(http.StatusOK, video.VideoFeedResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
	req := &video.VideoFeedRequest{
		UserId:     int64(actorId),
		LatestTime: &latestTime,
	}

	resp, err := videoClient.VideoFeed(ctx, req)

	if err != nil {
		zap.L().Error("Get feed list from video client err.", zap.Error(err))
		c.JSON(http.StatusOK, video.VideoFeedResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  common.MapErrMsg(resp.StatusCode),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetPublishList(ctx context.Context, c *app.RequestContext) {
	zap.L().Info("GetPublishList")
	actionId, _ := common.GetCurrentUserID(c)
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, video.VideoFeedResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
	request := &video.PublishVideoListRequest{
		FromUserId: int64(actionId),
		ToUserId:   userId,
	}
	resp, err := videoClient.GetPublishVideoList(ctx, request)
	if err != nil {
		zap.L().Error("GetPublishVideoListerr.", zap.Error(err))
		c.JSON(http.StatusOK, video.PublishVideoListResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  common.MapErrMsg(resp.StatusCode),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Publish(ctx context.Context, c *app.RequestContext) {
	zap.L().Info("Publish video")
	actionId, _ := common.GetCurrentUserID(c)
	title := c.PostForm("title")
	file, err := c.FormFile("data")

	if err != nil || title == "" || file.Size == 0 {
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidFileType(ext) {
		zap.L().Error("file type error")
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: common.CodeInvalidFileType,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidFileType),
		})
		return
	}

	// 校验文件大小
	if file.Size > maxFileSize || file.Size < minFileSize {
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: common.CodeInvalidFileSize,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidFileSize),
		})
		return
	}

	// 打开上传的文件
	src, err := file.Open()
	defer src.Close()
	if err != nil {
		zap.L().Error("打开文件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, video.PublishVideoResponse{
			StatusCode: common.CodeServerBusy,
			StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
		})
		return
	}
	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, video.PublishVideoResponse{
			StatusCode: common.CodeServerBusy,
			StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
		})
		return
	}
	req := &video.PublishVideoRequest{
		UserId: int64(actionId),
		Title:  title,
		Data:   data,
	}
	resp, err := videoClient.PublishVideo(ctx, req)
	if err != nil {
		zap.L().Error("PublishVideo err.", zap.Error(err))
		c.JSON(http.StatusOK, video.PublishVideoResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  common.MapErrMsg(resp.StatusCode),
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
