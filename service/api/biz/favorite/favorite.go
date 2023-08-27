package favorite

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/kitex_gen/favorite"
	"douyin/kitex_gen/favorite/favoriteservice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

var favoriteClient favoriteservice.Client

func init() {
	// OpenTelemetry 链路跟踪 还没配置好，先注释
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(config.CommentServiceName),
	//	provider.WithExportEndpoint("localhost:4317"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())

	// Etcd 服务发现
	r, err := etcd.NewEtcdResolver([]string{constant.EtcdAddr})
	if err != nil {
		log.Fatal(err)
	}

	favoriteClient, err = favoriteservice.NewClient(
		constant.FavoriteServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// Action 点赞取消赞的操作
func Action(ctx context.Context, c *app.RequestContext) {
	fromUserId, err := common.GetCurrentUserID(c)

	videoIdStr := c.Query("video_id")
	videoId, err := strconv.Atoi(videoIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, favorite.FavoriteActionResponse{
			StatusCode: 1,
		})
		return
	}
	actionTypeStr := c.Query("action_type")
	actionType, err := strconv.Atoi(actionTypeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, favorite.FavoriteActionResponse{
			StatusCode: 1,
		})
		return
	}
	req := &favorite.FavoriteActionRequest{
		UserId:     int64(fromUserId),
		VideoId:    int64(videoId),
		ActionType: int32(actionType),
	}

	resp, err := favoriteClient.FavoriteAction(ctx, req)
	if err != nil {
		zap.L().Error("FavoriteAction err.", zap.Error(err))
		c.JSON(http.StatusOK, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// List all users have same favorite video list
func List(ctx context.Context, c *app.RequestContext) {
	fromUserId, err := common.GetCurrentUserID(c)
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.Atoi(toUserIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, favorite.FavoriteListResponse{
			StatusCode: 1,
			StatusMsg:  "参数 error.",
		})
		return
	}

	req := &favorite.FavoriteListRequest{
		ActionId: int64(fromUserId),
		UserId:   int64(toUserId),
	}

	resp, err := favoriteClient.GetFavoriteList(ctx, req)
	if err != nil {
		zap.L().Error("GetFavoriteList err.", zap.Error(err))
		c.JSON(http.StatusOK, favorite.FavoriteListResponse{
			StatusCode: 1,
			StatusMsg:  "GetFavoriteList error.",
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}
