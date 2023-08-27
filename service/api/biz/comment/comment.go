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

	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
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
		c.JSON(http.StatusOK, "Unauthorized operation.")
		return
	}
	videoIdStr, videoIdExists := c.GetQuery("video_id")
	actionTypeStr, actionTypeExists := c.GetQuery("action_type")
	commentText, commentTextExists := c.GetQuery("comment_text")
	commentIdStr, commentIdExists := c.GetQuery("comment_id")

	// miss param, return
	if !videoIdExists || !actionTypeExists {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}

	// invalid param, return
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, "Invalid Params.")
		return
	}

	switch actionType {
	case 1: // create comment
		if !commentTextExists {
			c.JSON(http.StatusOK, "Invalid Params.")
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
			c.JSON(http.StatusOK, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	case 2: // delete comment
		if !commentIdExists {
			c.JSON(http.StatusOK, &comment.CommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid param.",
			})
			return
		}
		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, &comment.CommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid comment ID.",
			})
		}
		req := &comment.CommentActionRequest{
			UserId:     int64(userId),
			VideoId:    videoId,
			ActionType: int32(actionType),
			CommentId:  proto.Int32(int32(commentId)),
		}

		resp, err := commentClient.CommentAction(ctx, req)
		if err != nil {
			c.JSON(http.StatusOK, &comment.CommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "Server internal error.",
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	default: // wrong action_type
		c.JSON(http.StatusOK, &comment.CommentActionResponse{
			StatusCode: 1,
			StatusMsg:  "Invalid param.",
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
		c.JSON(http.StatusOK, err.Error())
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
			StatusCode:  1,
			StatusMsg:   "Server internal error.",
			CommentList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)

}
