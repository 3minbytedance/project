package relation

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"github.com/apache/thrift/lib/go/thrift"

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
	userId, err := common.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
		c.JSON(http.StatusOK, "Unauthorized operation.")
		return
	}
	toUserIdStr, toUserIdExists := c.GetQuery("to_user_id")
	actionTypeStr, actionTypeExists := c.GetQuery("action_type")

	// miss param, return
	if !toUserIdExists || !actionTypeExists {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}

	// invalid param, return
	toUserId, err := strconv.ParseUint(toUserIdStr, 10, 32)
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}
	userIdUint := int32(userId)
	req := &relation.RelationActionRequest{
		UserId:     userIdUint,
		ToUserId:   int32(toUserId),
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
				StatusMsg:  thrift.StringPtr("Server internal error."),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	default: // wrong action_type
		c.JSON(http.StatusOK, &relation.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Invalid param."),
		})
		return
	}
}

func FollowList(ctx context.Context, c *app.RequestContext) {
	userId, err := common.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Unauthorized operation."),
		})
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseUint(toUserIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}
	req := &relation.FollowListRequest{
		UserId: int32(userId),
		ToUserId: int32(toUserId),
	}

	resp, err := relationClient.GetFollowList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Server internal error."),
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}

func FollowerList(ctx context.Context, c *app.RequestContext) {
	// 已经有鉴权中间件，鉴过token了
	userId, err := common.GetCurrentUserID(c)
	// not logged in
	if err != nil {
		c.JSON(http.StatusOK, "Unauthorized operation.")
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseUint(toUserIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}

	req := &relation.FollowerListRequest{
		UserId:   int32(userId),
		ToUserId: int32(toUserId),
	}

	resp, err := relationClient.GetFollowerList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowerListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Server internal error."),
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}

func FriendList(ctx context.Context, c *app.RequestContext) {
	// 已经有鉴权中间件，鉴过token了
	userId, err := common.GetCurrentUserID(c)
	// not logged in
	if err != nil {
		c.JSON(http.StatusOK, "Unauthorized operation.")
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseUint(toUserIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}

	req := &relation.FriendListRequest{
		UserId: int32(userId),
		ToUserId: int32(toUserId),
	}
	resp, err := relationClient.GetFriendList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FriendListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Server internal error."),
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}
