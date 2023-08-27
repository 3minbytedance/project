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
		c.JSON(http.StatusOK, relation.RelationActionResponse{
			StatusCode: common.CodeInvalidToken,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidToken),
		})
		return
	}
	actionIdStr := strconv.FormatUint(uint64(actionId), 10)
	toUserIdStr, toUserIdExists := c.GetQuery("to_user_id")
	actionTypeStr, actionTypeExists := c.GetQuery("action_type")

	// miss param, return
	if !toUserIdExists || !actionTypeExists {
		c.JSON(http.StatusOK, relation.RelationActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
	if actionIdStr == toUserIdStr {
		c.JSON(http.StatusOK, relation.RelationActionResponse{
			StatusCode: common.CodeFollowMyself,
			StatusMsg:  common.MapErrMsg(common.CodeFollowMyself),
		})
		return
	}

	// invalid param, return
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, relation.RelationActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
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
			c.JSON(http.StatusOK, relation.RelationActionResponse{
				StatusCode: resp.StatusCode,
				StatusMsg:  common.MapErrMsg(resp.StatusCode),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	default: // wrong action_type
		c.JSON(http.StatusOK, relation.RelationActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
}

func FollowList(ctx context.Context, c *app.RequestContext) {
	actionId, err := common.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: common.CodeInvalidToken,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidToken),
		})
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
	req := &relation.FollowListRequest{
		UserId:   int64(actionId),
		ToUserId: toUserId,
	}

	resp, err := relationClient.GetFollowList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusInternalServerError, relation.FollowListResponse{
			StatusCode: common.CodeServerBusy,
			StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
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
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: common.CodeInvalidToken,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidToken),
		})
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, relation.FollowListResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}

	req := &relation.FollowerListRequest{
		UserId:   int64(actionId),
		ToUserId: toUserId,
	}

	resp, err := relationClient.GetFollowerList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusInternalServerError, relation.FollowListResponse{
			StatusCode: common.CodeServerBusy,
			StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
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
		c.JSON(http.StatusOK, relation.FriendListResponse{
			StatusCode: common.CodeInvalidToken,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidToken),
		})
		return
	}
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, relation.FriendListResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}

	req := &relation.FriendListRequest{
		UserId:   int64(actionId),
		ToUserId: toUserId,
	}
	resp, err := relationClient.GetFriendList(ctx, req)
	if err != nil {
		zap.L().Error("Get follow list from relation client err.", zap.Error(err))
		c.JSON(http.StatusInternalServerError, relation.FriendListResponse{
			StatusCode: common.CodeServerBusy,
			StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}
