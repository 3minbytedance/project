package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/dal/model"
	"douyin/dal/mysql"
	comment "douyin/kitex_gen/comment"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/mw/kafka"
	"douyin/mw/redis"
	"douyin/service/comment/pack"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}),
		client.WithMuxConnection(1),
	)
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
		zap.Int64("user_id", request.GetUserId()),
		zap.Int64("video_id", request.GetVideoId()),
		zap.Int32("action_type", request.GetActionType()),
		zap.Int32("comment_id", request.GetCommentId()),
		zap.String("comment_text", request.GetCommentText()),
	)
	videoId := request.GetVideoId()
	// 查询user是否存在，并在新增评论后返回该用户信息
	userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
		ActorId: request.GetUserId(),
		UserId:  request.GetUserId(),
	})
	if userResp.GetUser() == nil || err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = err.Error()
		return resp, err
	}

	switch request.GetActionType() {
	case 1: // 新增评论
		commentData := model.Comment{
			UserId:    uint(request.GetUserId()),
			VideoId:   uint(request.GetVideoId()),
			Content:   common.ReplaceWord(request.GetCommentText()),
			CreatedAt: time.Now(),
		}
		// _, err = mysql.AddComment(&commentData)
		err = kafka.CommentMQInstance.ProduceAddCommentMsg(&commentData)
		if err != nil {
			resp.StatusCode = common.CodeDBError
			resp.StatusMsg = common.MapErrMsg(common.CodeDBError)
			return
		}
		// 增加redis
		go func() {
			// 如果video不存在于redis，查询数据库并插入redis评论数
			isSetKey, _ := checkAndSetRedisCommentKey(uint(videoId))
			if isSetKey {
				return
			}
			// 如果video存在于redis，更新commentCount
			err = redis.IncrementCommentCountByVideoId(uint(videoId))
			if err != nil {
				zap.L().Error("更新videoId的评论数失败", zap.Error(err))
			}
		}()
		// 封装返回数据
		//comment := pack.Comment(&commentData, user.User)
		return &comment.CommentActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
			Comment:    pack.Comment(&commentData, userResp.GetUser()),
		}, nil

	case 2:
		// 设置redis
		go func() {
			// todo
			isSetKey, _ := checkAndSetRedisCommentKey(uint(videoId))
			if isSetKey {
				return
			}
			err = redis.DecrementCommentCountByVideoId(uint(videoId))
			if err != nil {
				zap.L().Error("DecrementCommentCountByVideoId error", zap.Error(err))
				return
			}
		}()

		// err = mysql.DeleteCommentById(uint(request.GetCommentId()))
		err = kafka.CommentMQInstance.ProduceDelCommentMsg(uint(request.GetCommentId()))
		if err != nil {
			zap.L().Error("DeleteCommentById error", zap.Error(err))
			return &comment.CommentActionResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
			}, err
		}
		return &comment.CommentActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		}, nil
	default:
		return &comment.CommentActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		}, errors.New(common.MapErrMsg(common.CodeInvalidParam))
	}
}

// GetCommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentList(ctx context.Context, request *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	comments, err := mysql.FindCommentsByVideoId(uint(request.VideoId))
	if err != nil {
		zap.L().Error("根据视频ID取评论失败", zap.Error(err))
		return &comment.CommentListResponse{
			StatusCode:  common.CodeDBError,
			StatusMsg:   common.MapErrMsg(common.CodeDBError),
			CommentList: nil,
		}, err
	}
	commentList := make([]*comment.Comment, 0, len(comments))
	userResponses := make([]*user.UserInfoByIdResponse, len(comments))
	actionId := request.GetUserId()
	var wg sync.WaitGroup
	wg.Add(len(comments))

	for i, com := range comments {
		go func(index int, c model.Comment) {
			defer wg.Done()
			userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
				ActorId: actionId,
				UserId:  int64(c.UserId),
			})
			if err != nil {
				err = nil
				zap.L().Error("查询评论用户信息失败", zap.Error(err))
				return
			}
			userResponses[index] = userResp
		}(i, com)
	}
	wg.Wait()

	// 处理 userResponses
	for i, com := range comments {
		userResp := userResponses[i]
		commentList = append(commentList, pack.Comment(&com, userResp.GetUser()))
	}
	return &comment.CommentListResponse{
		StatusCode:  common.CodeSuccess,
		StatusMsg:   common.MapErrMsg(common.CodeSuccess),
		CommentList: commentList,
	}, nil
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, videoId int64) (resp int32, err error) {
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

// checkAndSetRedisCommentKey
// 返回true表示不存在这个key，并设置key
// 返回false表示已存在这个key，cnt数返回0
func checkAndSetRedisCommentKey(videoId uint) (isSet bool, count int64) {
	if redis.IsExistVideoField(videoId, redis.CommentCountField) {
		return false, 0
	}
	//缓存不存在，尝试从数据库中取
	if redis.AcquireCommentLock(videoId) {
		defer redis.ReleaseCommentLock(videoId)
		//double check
		if redis.IsExistVideoField(videoId, redis.CommentCountField) {
			return false, 0
		}
		// 获取最新commentCount
		cnt, err := mysql.GetCommentCnt(videoId)
		if err != nil {
			zap.L().Error("mysql获取评论数失败", zap.Error(err))
			return false, 0
		}
		// 设置最新commentCount
		err = redis.SetCommentCountByVideoId(videoId, cnt)
		if err != nil {
			zap.L().Error("redis更新评论数失败", zap.Error(err))
			return false, 0
		}
		return true, cnt
	}
	fmt.Println("重试checkAndSetRedisCommentKey")
	// 重试
	time.Sleep(redis.RetryTime)
	return checkAndSetRedisCommentKey(videoId)
}
