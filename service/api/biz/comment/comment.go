package comment

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/kitex_gen/comment"
	"douyin/kitex_gen/comment/commentservice"
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

var commentClient commentservice.Client

func init() {

	// Etcd 服务发现
	r, err := etcd.NewEtcdResolver([]string{constant.EtcdAddr})
	if err != nil {
		log.Fatal(err)
	}

	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
		client.WithMuxConnection(2),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(ctx context.Context, c *app.RequestContext) {
	userId, err := common.GetCurrentUserID(c)
	// not logged in
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
		c.JSON(http.StatusOK, comment.CommentActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
	videoIdStr, videoIdExists := c.GetQuery("video_id")
	actionTypeStr, actionTypeExists := c.GetQuery("action_type")
	commentText, commentTextExists := c.GetQuery("comment_text")
	commentIdStr, commentIdExists := c.GetQuery("comment_id")

	// miss param, return
	if !videoIdExists || !actionTypeExists {
		c.JSON(http.StatusOK, comment.CommentActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}

	// invalid param, return
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, comment.CommentActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}

	switch actionType {
	case 1: // create comment
		if !commentTextExists {
			c.JSON(http.StatusOK, comment.CommentActionResponse{
				StatusCode: common.CodeInvalidParam,
				StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
			})
			return
		}
		req := &comment.CommentActionRequest{
			UserId:      int64(userId),
			VideoId:     videoId,
			ActionType:  int32(actionType),
			CommentText: proto.String(commentText),
		}
		resp, err := commentClient.CommentAction(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, comment.CommentActionResponse{
				StatusCode: resp.StatusCode,
				StatusMsg:  common.MapErrMsg(resp.StatusCode),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	case 2: // delete comment
		if !commentIdExists {
			c.JSON(http.StatusOK, comment.CommentActionResponse{
				StatusCode: common.CodeInvalidParam,
				StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
			})
			return
		}
		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, comment.CommentActionResponse{
				StatusCode: common.CodeInvalidParam,
				StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
			})
			return
		}
		req := &comment.CommentActionRequest{
			UserId:     int64(userId),
			VideoId:    videoId,
			ActionType: int32(actionType),
			CommentId:  proto.Int32(int32(commentId)),
		}

		resp, err := commentClient.CommentAction(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, comment.CommentActionResponse{
				StatusCode: resp.StatusCode,
				StatusMsg:  common.MapErrMsg(resp.StatusCode),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	default: // wrong action_type
		c.JSON(http.StatusOK, comment.CommentActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}
}

func List(ctx context.Context, c *app.RequestContext) {
	userId, err := common.GetCurrentUserID(c)
	// not logged in
	if err != nil {
		zap.L().Error("Get user id from ctx", zap.Error(err))
	}
	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		zap.L().Error("Parse videoIdStr err:", zap.Error(err))
		c.JSON(http.StatusOK, comment.CommentListResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}

	req := &comment.CommentListRequest{
		UserId:  int64(userId),
		VideoId: videoId,
	}

	resp, err := commentClient.GetCommentList(ctx, req)
	if err != nil {
		zap.L().Error("Get comment list from comment client err.", zap.Error(err))
		c.JSON(http.StatusOK, comment.CommentListResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  common.MapErrMsg(resp.StatusCode),
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}
