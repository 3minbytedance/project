package main

import (
	"context"
	"douyin/constant"
	"douyin/dal/model"
	"douyin/dal/mysql"
	"douyin/kitex_gen/comment"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/service/comment/pack"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
)

var userClient userservice.Client

func init() {
	// Etcd 服务发现
	r, err := etcd.NewEtcdResolver([]string{constant.EtcdAddr})
	if err != nil {
		log.Fatal(err)
	}
	userClient, err = userservice.NewClient(
		constant.UserServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}))
	if err != nil {
		log.Fatal(err)
	}
}

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, request *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	// TODO: Your code here...
	resp = new(comment.CommentActionResponse)
	zap.L().Info("CommentClient action start",
		zap.Int32("user_id", request.UserId),
		zap.Int32("video_id", request.VideoId),
		zap.Int32("action_type", request.ActionType),
		zap.Int32("comment_id", request.GetCommentId()),
		zap.String("comment_text", request.GetCommentText()),
	)

	switch request.ActionType {
	case 1: // 新增评论
		commentData := model.Comment{
			UserId:  uint(request.UserId),
			VideoId: uint(request.VideoId),
			Content: request.GetCommentText(),
		}
		_, err = mysql.AddComment(&commentData)
		if err != nil {
			resp.StatusCode = int32(1)
			resp.StatusMsg = thrift.StringPtr(err.Error())
			return
		}
		// todo: redis

		// 查询user
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{UserId: request.UserId})
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = thrift.StringPtr(err.Error())
			return resp, err
		}

		// 封装返回数据
		//comment := pack.Comment(&commentData, user.User)
		return &comment.CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success"),
			Comment:    pack.Comment(&commentData, userResp.User),
		}, nil

	case 2:
		return
	default:
		return
	}

	return
}

// GetCommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentList(ctx context.Context, request *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, request *comment.CommentCountRequest) (resp *comment.CommentCountResponse, err error) {
	// TODO: Your code here...
	return
}
