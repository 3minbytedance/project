package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/dal/model"
	"douyin/dal/mysql"
	"douyin/kitex_gen/comment"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/mw/redis"
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}))
	if err != nil {
		log.Fatal(err)
	}
}

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, request *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	resp = new(comment.CommentActionResponse)
	zap.L().Info("CommentClient action start",
		zap.Int32("user_id", request.GetUserId()),
		zap.Int32("video_id", request.GetVideoId()),
		zap.Int32("action_type", request.GetActionType()),
		zap.Int32("comment_id", request.GetCommentId()),
		zap.String("comment_text", request.GetCommentText()),
	)

	switch request.GetActionType() {
	case 1: // 新增评论
		commentData := model.Comment{
			UserId:  uint(request.GetUserId()),
			VideoId: uint(request.GetVideoId()),
			Content: common.ReplaceWord(request.GetCommentText()),
		}
		_, err = mysql.AddComment(&commentData)
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = thrift.StringPtr(err.Error())
			return
		}

		// 设置redis
		go func() {
			// 如果video不存在于redis，查询数据库并插入redis评论数
			isSetKey, _ := checkAndSetRedisCommentKey(uint(request.GetVideoId()))
			if isSetKey {
				return
			}
			// 如果video存在于redis，更新commentCount
			err = redis.IncrementCommentCountByVideoId(uint(request.GetVideoId()))
			if err != nil {
				log.Printf("更新videoId为%v的评论数失败  %v\n", request.GetVideoId(), err.Error())
			}
		}()

		// 查询user
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{ActorId: request.GetUserId()})
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
		// 设置redis
		go func() {
			err = redis.DecrementCommentCountByVideoId(uint(request.GetVideoId()))
			if err != nil {
				zap.L().Error("DecrementCommentCountByVideoId error", zap.Error(err))
				return
			}
		}()

		err = mysql.DeleteCommentById(uint(request.GetCommentId()))
		if err != nil {
			zap.L().Error("DeleteCommentById error", zap.Error(err))
			return &comment.CommentActionResponse{
				StatusCode: 1,
				StatusMsg:  thrift.StringPtr("Internal server error"),
			}, err
		}
		return &comment.CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success"),
		}, nil
	default:
		return &comment.CommentActionResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Invalid Param"),
		}, nil
	}

}

// GetCommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentList(ctx context.Context, request *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	comments, err := mysql.FindCommentsByVideoId(uint(request.VideoId))
	if err != nil {
		zap.L().Error("根据视频ID取评论失败", zap.Error(err))
		return &comment.CommentListResponse{
			StatusCode:  1,
			StatusMsg:   thrift.StringPtr("获取评论列表失败"),
			CommentList: nil,
		}, err
	}
	commentList := make([]*comment.Comment, 0, len(comments))
	for _, com := range comments {
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			ActorId: request.GetUserId(),
			UserId:  int32(com.UserId),
		})
		if err != nil {
			zap.L().Error("查询评论用户信息失败", zap.Error(err))
			return &comment.CommentListResponse{
				StatusCode:  1,
				StatusMsg:   thrift.StringPtr("查询评论用户信息失败"),
				CommentList: nil,
			}, err
		}
		commentList = append(commentList, pack.Comment(&com, userResp.User))
	}
	return &comment.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   thrift.StringPtr("success"),
		CommentList: commentList,
	}, nil
}

// checkAndSetRedisCommentKey
// 返回true表示不存在这个key，并设置key
// 返回false表示已存在这个key，cnt数返回0
func checkAndSetRedisCommentKey(videoId uint) (isSet bool, count int64) {
	var cnt int64
	if !redis.IsExistVideoField(videoId, redis.CommentCountField) {
		// 获取最新commentCount
		cnt, err := mysql.GetCommentCnt(videoId)
		if err != nil {
			zap.L().Error("mysql获取评论数失败", zap.Error(err))
		}
		// 设置最新commentCount
		err = redis.SetCommentCountByVideoId(videoId, cnt)
		if err != nil {
			zap.L().Error("redis更新评论数失败", zap.Error(err))
		}
		return true, cnt
	}
	return false, cnt
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, videoId int32) (resp int32, err error) {
	isSetKey, count := checkAndSetRedisCommentKey(uint(videoId))
	if isSetKey {
		return int32(count), nil
	}
	// 从redis中获取评论数
	count, err = redis.GetCommentCountByVideoId(uint(videoId))
	if err != nil {
		zap.L().Error("redis获取评论数失败", zap.Error(err))
		return 0, err
	}
	return int32(count), nil
}
