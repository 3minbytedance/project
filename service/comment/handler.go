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
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"strconv"
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
		client.WithMuxConnection(2),
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
		zap.Int64("comment_id", request.GetCommentId()),
		zap.String("comment_text", request.GetCommentText()),
	)
	videoId := uint(request.GetVideoId())
	userId := request.GetUserId()
	// 查询user是否存在，并在新增评论后返回该用户信息
	userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
		ActorId: userId,
		UserId:  userId,
	})
	if userResp.GetUser() == nil || err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = err.Error()
		err = nil
		return resp, nil
	}

	switch request.GetActionType() {
	case 1: // 新增评论
		commentData := model.Comment{
			ID:        common.GetUid(),
			UserId:    uint(userId),
			VideoId:   videoId,
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

		go func() {
			common.AddToCommentBloom(strconv.Itoa(int(videoId)))
		}()
		// cache aside
		redis.DelVideoHashField(videoId, redis.CommentCountField)

		// 封装返回数据
		//comment := pack.Comment(&commentData, user.User)
		return &comment.CommentActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
			Comment:    pack.Comment(&commentData, userResp.GetUser()),
		}, nil

	case 2:

		// 查询评论id是否属于该用户
		belongsToUser, err := mysql.IsCommentBelongsToUser(request.CommentId, request.UserId)
		if err != nil {
			return &comment.CommentActionResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
			}, err
		}
		if !belongsToUser {
			return &comment.CommentActionResponse{
				StatusCode: common.CodeInvalidCommentAction,
				StatusMsg:  common.MapErrMsg(common.CodeInvalidCommentAction),
			}, nil
		}

		// err = mysql.DeleteCommentById(uint(request.GetCommentId()))
		err = kafka.CommentMQInstance.ProduceDelCommentMsg(uint(request.GetCommentId()))
		if err != nil {
			zap.L().Error("DeleteCommentById error", zap.Error(err))
			return &comment.CommentActionResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
			}, err
		}

		// cache aside
		redis.DelVideoHashField(videoId, redis.CommentCountField)
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
		if userResp == nil {
			continue
		}
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
	_, count := checkAndSetRedisCommentKey(uint(videoId))
	return int32(count), nil
}

// checkAndSetRedisCommentKey
// 返回true表示不存在这个key，并更新key，并返回评论数
// 返回false表示已存在这个key，未更新，并返回评论数
func checkAndSetRedisCommentKey(videoId uint) (isSet bool, commentCount int64) {
	//缓存中有数据
	if count, err := redis.GetCommentCountByVideoId(videoId); err == nil {
		return false, count
	}
	//缓存不存在，尝试从数据库中取
	if redis.AcquireCommentLock(videoId) {
		defer redis.ReleaseCommentLock(videoId)

		exist := common.TestCommentBloom(strconv.Itoa(int(videoId)))

		// 不存在
		if !exist {
			err := redis.SetCommentCountByVideoId(videoId, 0)
			if err != nil {
				zap.L().Error("redis更新评论数失败", zap.Error(err))
				return false, 0
			}
			return true, 0
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
	// 重试
	time.Sleep(redis.RetryTime)
	return checkAndSetRedisCommentKey(videoId)
}
