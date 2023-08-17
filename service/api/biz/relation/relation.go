package relation

import (
	"context"
	"douyin/constant"
	"douyin/kitex_gen/comment"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
	userId, userIdExists := c.Get("userId")
	// not logged in
	if !userIdExists {
		c.JSON(http.StatusOK, "Unauthorized operation.")
		return
	}
	to_user_id_str, to_user_id_Exists := c.GetQuery("to_user_id")
	actionTypeStr, actionTypeExists := c.GetQuery("action_type")

	// miss param, return
	if !to_user_id_Exists || !actionTypeExists {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}

	// invalid param, return
	to_user_id, err := strconv.ParseUint(to_user_id_str, 10, 32)
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}
	userIdUint := int32(userId.(uint))
	req := &relation.RelationActionRequest{
		UserId:     userIdUint,
		ToUserId:   int32(to_user_id),
		ActionType: int32(actionType),
	}
	// TODO: judge userId

	switch actionType {
	case 1: // 关注
	case 2: // 取关
		resp, err := relationClient.RelationAction(ctx, req)
		if err != nil {
			c.JSON(http.StatusOK, &comment.CommentActionResponse{
				StatusCode: 1,
				StatusMsg:  proto.String("Server internal error."),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	default: // wrong action_type
		c.JSON(http.StatusOK, &comment.CommentActionResponse{
			StatusCode: 1,
			StatusMsg:  proto.String("Invalid param."),
		})
		return
	}
}

func Follow_List(ctx context.Context, c *app.RequestContext) {
	// 已经有鉴权中间件，鉴过token了
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		zap.L().Error("Parse videoIdStr err:", zap.Error(err))
		c.JSON(http.StatusOK, err.Error())
		return
	}

	req := &relation.FollowListRequest{
		UserId:  int32(userId),
		ActorId: 1,
	}

	resp, err := relationClient.GetFollowList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: 1,
			StatusMsg:  proto.String("Server internal error."),
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}

func Follower_List(ctx context.Context, c *app.RequestContext) {
	// 已经有鉴权中间件，鉴过token了
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		zap.L().Error("Parse videoIdStr err:", zap.Error(err))
		c.JSON(http.StatusOK, err.Error())
		return
	}

	req := &relation.FollowerListRequest{
		UserId:  int32(userId),
		ActorId: 1,
	}

	resp, err := relationClient.GetFollowerList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowerListResponse{
			StatusCode: 1,
			StatusMsg:  proto.String("Server internal error."),
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}

func Friend_List(ctx context.Context, c *app.RequestContext) {
	// 已经有鉴权中间件，鉴过token了
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		zap.L().Error("Parse videoIdStr err:", zap.Error(err))
		c.JSON(http.StatusOK, err.Error())
		return
	}

	req := &relation.FriendListRequest{
		UserId: int32(userId),
	}

	resp, err := relationClient.GetFriendList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusOK, relation.FriendListResponse{
			StatusCode: 1,
			StatusMsg:  proto.String("Server internal error."),
			UserList:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}
