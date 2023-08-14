package comment

import (
	"context"
	"douyin/constant"
	"douyin/kitex_gen/comment"
	"douyin/kitex_gen/comment/commentservice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
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
	userId, userIdExists := c.Get("userId")
	// not logged in
	if !userIdExists {
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
	videoId, err := strconv.ParseUint(videoIdStr, 10, 32)
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 32)
	commentId, err := strconv.ParseUint(commentIdStr, 10, 32)
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
			UserId:      userId.(uint32),
			VideoId:     uint32(videoId),
			ActionType:  uint32(actionType),
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
			c.JSON(http.StatusOK, "Invalid Params.")
			return
		}
		req := &comment.CommentActionRequest{
			UserId:     userId.(uint32),
			VideoId:    uint32(videoId),
			ActionType: uint32(actionType),
			CommentId:  proto.Uint32(uint32(commentId)),
		}

		resp, err := commentClient.CommentAction(ctx, req)
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	default: // wrong action_type
		c.JSON(http.StatusOK, "Invalid action type.")
		return
	}
}

func List(ctx context.Context, c *app.RequestContext) {

}
