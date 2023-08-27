package relation

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
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

var relationClient relationservice.Client

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

	relationClient, err = relationservice.NewClient(
		constant.RelationServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.RelationServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(ctx context.Context, c *app.RequestContext) {
	actionId, err := common.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
		c.JSON(http.StatusNonAuthoritativeInfo, common.Response{StatusCode: 1, StatusMsg: "未授权的行为"})
		return
	}
	actionIdStr := strconv.FormatUint(uint64(actionId), 10)
	toUserIdStr, toUserIdExists := c.GetQuery("to_user_id")
	actionTypeStr, actionTypeExists := c.GetQuery("action_type")

	// miss param, return
	if !toUserIdExists || !actionTypeExists {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "参数不合法"})
		return
	}
	if actionIdStr == toUserIdStr {
		c.JSON(http.StatusOK, &relation.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "不能对自己操作",
		})
		return
	}

	// invalid param, return
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "参数不合法"})
		return
	}
	req := &relation.RelationActionRequest{
		UserId:     int64(actionId),
		ToUserId:   toUserId,
		ActionType: int32(actionType),
	}
	// TODO: judge userId
	zap.L().Debug("ACTIONTYPE", zap.Int("AT", int(actionType)))

	switch actionType {
	case 1, 2: // 关注
		resp, err := relationClient.RelationAction(ctx, req)
		if err != nil {
			c.JSON(http.StatusOK, &relation.RelationActionResponse{
				StatusCode: 1,
				StatusMsg:  "Server internal error.",
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	default: // wrong action_type
		c.JSON(http.StatusOK, &relation.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "Invalid param.",
		})
		return
	}
}

func FollowList(ctx context.Context, c *app.RequestContext) {
	actionId, err := common.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: 1,
			StatusMsg:  "Unauthorized operation.",
		})
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "参数不合法"})
		return
	}
	req := &relation.FollowListRequest{
		UserId:   int64(actionId),
		ToUserId: toUserId,
	}

	resp, err := relationClient.GetFollowList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: 1,
			StatusMsg:  "Server internal error.",
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}

func FollowerList(ctx context.Context, c *app.RequestContext) {
	// 已经有鉴权中间件，鉴过token了
	actionId, err := common.GetCurrentUserID(c)
	// not logged in
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "未授权的行为"})
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "参数不合法"})
		return
	}

	req := &relation.FollowerListRequest{
		UserId:   int64(actionId),
		ToUserId: toUserId,
	}

	resp, err := relationClient.GetFollowerList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowerListResponse{
			StatusCode: 1,
			StatusMsg:  "Server internal error.",
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}

func FriendList(ctx context.Context, c *app.RequestContext) {
	// 已经有鉴权中间件，鉴过token了
	actionId, err := common.GetCurrentUserID(c)
	// not logged in
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "未授权的行为"})
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "参数不合法"})
		return
	}

	req := &relation.FriendListRequest{
		UserId:   int64(actionId),
		ToUserId: toUserId,
	}
	resp, err := relationClient.GetFriendList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FriendListResponse{
			StatusCode: 1,
			StatusMsg:  "Server internal error.",
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}
