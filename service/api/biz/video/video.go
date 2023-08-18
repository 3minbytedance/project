package video

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/kitex_gen/video"
	"douyin/kitex_gen/video/videoservice"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

var videoClient videoservice.Client

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
	userId, err := common.GetCurrentUserID(c)
	// not logged in
	//todo 后面再改
	var id *int32
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
	}
	if userId > 0 {
		loginId := int32(userId)
		id = &loginId
	}
	latestTime := c.Query("latest_time")
	unixTime, err := strconv.Atoi(latestTime)
	if latestTime == "" || latestTime == "0" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	} else if err != nil || unixTime < 0 {
		zap.L().Error("Parse videoIdStr err:", zap.Error(err))
		c.JSON(http.StatusOK, err.Error())
		return
	}
	req := &video.VideoFeedRequest{
		UserId:     id,
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

}

func Publish(ctx context.Context, c *app.RequestContext) {

}