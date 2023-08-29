package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/dal/mysql"
	relation "douyin/kitex_gen/relation"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/mw/kafka"
	"douyin/mw/redis"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
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

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, request *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	resp = new(relation.RelationActionResponse)
	zap.L().Info("RelationClient action start",
		zap.Int64("user_id", request.UserId),
		zap.Int32("action_type", request.ActionType),
		zap.Int64("ToUserId", request.ToUserId),
	)
	fromUserId := request.GetUserId()
	toUserId := request.GetToUserId()

	switch request.ActionType {
	case 1: // 关注
		// 判断用户是否已经关注过了
		res, err := redis.IsInMyFollowList(uint(fromUserId), uint(toUserId))
		if res {
			return &relation.RelationActionResponse{
				StatusCode: common.CodeSuccess,
				StatusMsg:  common.MapErrMsg(common.CodeSuccess),
			}, nil
		}
		// 延迟双删
		redis.DelKey(uint(fromUserId), redis.FollowList)
		redis.DelKey(uint(toUserId), redis.FollowerList)

		//err = mysql.AddFollow(uint(fromUserId), uint(toUserId))
		err = kafka.FollowMQInstance.ProduceAddFollowMsg(uint(fromUserId), uint(toUserId))
		if err != nil {
			return &relation.RelationActionResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
			}, err
		}
		go func() {
			time.Sleep(1 * time.Second)

			redis.DelKey(uint(fromUserId), redis.FollowList)
			redis.DelKey(uint(toUserId), redis.FollowerList)
		}()

		return &relation.RelationActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		}, nil
	case 2: // 取关
		// 延迟双删
		redis.DelKey(uint(fromUserId), redis.FollowList)
		redis.DelKey(uint(toUserId), redis.FollowerList)

		//err = mysql.DeleteFollowById(uint(fromUserId), uint(toUserId))
		err = kafka.FollowMQInstance.ProduceDelFollowMsg(uint(fromUserId), uint(toUserId))
		if err != nil {
			return &relation.RelationActionResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
			}, err
		}
		go func() {
			time.Sleep(1 * time.Second)

			redis.DelKey(uint(fromUserId), redis.FollowList)
			redis.DelKey(uint(toUserId), redis.FollowerList)
		}()

		return &relation.RelationActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		}, nil

	default:
		return &relation.RelationActionResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		}, errors.New(common.MapErrMsg(common.CodeInvalidParam))
	}
}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, request *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	actionId := request.GetUserId()
	toUserId := request.GetToUserId()

	res := CheckAndSetRedisRelationKey(uint(toUserId), redis.FollowList)
	if res == 2 {
		return &relation.FollowListResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFollowListById(uint(toUserId))
	if err != nil {
		return &relation.FollowListResponse{
			StatusCode: common.CodeDBError,
			StatusMsg:  common.MapErrMsg(common.CodeDBError),
			UserList:   nil,
		}, err
	}
	followList := make([]*user.User, 0, len(id))
	for _, com := range id {
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			ActorId: actionId,
			UserId:  int64(com),
		})
		if err != nil {
			zap.L().Error("查询Follow用户信息失败", zap.Error(err))
			return &relation.FollowListResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),

				UserList: nil,
			}, err
		}
		followList = append(followList, userResp.GetUser())
	}
	return &relation.FollowListResponse{
		StatusCode: common.CodeSuccess,
		StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		UserList:   followList,
	}, nil
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, request *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	actionId := request.GetUserId()
	toUserId := request.GetToUserId()
	res := CheckAndSetRedisRelationKey(uint(toUserId), redis.FollowerList)
	if res == 2 {
		return &relation.FollowerListResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFollowerListById(uint(toUserId))
	if err != nil {
		return &relation.FollowerListResponse{
			StatusCode: common.CodeDBError,
			StatusMsg:  common.MapErrMsg(common.CodeDBError),
			UserList:   nil,
		}, err
	}
	followerList := make([]*user.User, 0, len(id))
	for _, com := range id {
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			ActorId: actionId,
			UserId:  int64(com),
		})
		if err != nil {
			zap.L().Error("查询Follower用户信息失败", zap.Error(err))
			return &relation.FollowerListResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
				UserList:   nil,
			}, err
		}
		followerList = append(followerList, userResp.GetUser())
	}
	return &relation.FollowerListResponse{
		StatusCode: common.CodeSuccess,
		StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		UserList:   followerList,
	}, nil
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, request *relation.FriendListRequest) (resp *relation.FriendListResponse, err error) {
	actionId := request.GetUserId()
	toUserId := request.GetToUserId()
	res := CheckAndSetRedisRelationKey(uint(toUserId), redis.FollowList)
	if res == 2 {
		return &relation.FriendListResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
			UserList:   nil,
		}, nil
	}
	res = CheckAndSetRedisRelationKey(uint(toUserId), redis.FollowerList)
	if res == 2 {
		return &relation.FriendListResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFriendListById(uint(toUserId))
	if err != nil {
		return &relation.FriendListResponse{
			StatusCode: common.CodeDBError,
			StatusMsg:  common.MapErrMsg(common.CodeDBError),
			UserList:   nil,
		}, err
	}
	friendList := make([]*user.User, 0, len(id))
	for _, com := range id {
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			ActorId: actionId,
			UserId:  int64(com),
		})
		if err != nil {
			zap.L().Error("查询Friend用户信息失败", zap.Error(err))
			return &relation.FriendListResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
				UserList:   nil,
			}, err
		}
		friendList = append(friendList, userResp.User)
	}
	return &relation.FriendListResponse{
		StatusCode: common.CodeSuccess,
		StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		UserList:   friendList,
	}, nil
}

// GetFollowListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowListCount(ctx context.Context, userId int64) (resp int32, err error) {
	res := CheckAndSetRedisRelationKey(uint(userId), redis.FollowList)
	if res == 2 {
		return 0, nil
	}
	count, _ := redis.GetFollowCountById(uint(userId))
	return int32(count), nil
}

// GetFollowerListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerListCount(ctx context.Context, userId int64) (resp int32, err error) {
	res := CheckAndSetRedisRelationKey(uint(userId), redis.FollowerList)
	if res == 2 {
		return 0, nil
	}
	count, _ := redis.GetFollowerCountById(uint(userId))
	return int32(count), nil
}

// IsFollowing 判断actionID的关注表里是否有toUserID
func (s *RelationServiceImpl) IsFollowing(ctx context.Context, request *relation.IsFollowingRequest) (resp bool, err error) {
	actionId := request.GetActorId()
	toUserId := request.GetUserId()

	res := CheckAndSetRedisRelationKey(uint(actionId), redis.FollowList)
	if res == 2 {
		return false, nil
	}

	found, err := redis.IsInMyFollowList(uint(actionId), uint(toUserId))
	return found, err
}

// IsFriend 判断二者是不是friend
func (s *RelationServiceImpl) IsFriend(ctx context.Context, request *relation.IsFriendRequest) (resp bool, err error) {
	// 从数据库查询是否已关注
	// 先判断actorid的关注ser和粉丝set是否存在
	res := CheckAndSetRedisRelationKey(uint(request.ActorId), redis.FollowerList)
	if res == 2 {
		return false, nil
	}
	res = CheckAndSetRedisRelationKey(uint(request.ActorId), redis.FollowList)
	if res == 2 {
		return false, nil
	}

	// 如果发生err了，返回false吧
	result1, err := redis.IsInMyFollowerList(uint(request.ActorId), uint(request.UserId))
	if err != nil {
		return false, err
	}
	result2, err := redis.IsInMyFollowList(uint(request.ActorId), uint(request.UserId))
	if err != nil {
		return false, err
	}
	//result, err := mysql.IsFriend(uint(request.ActorId), uint(request.UserId))
	return result2 && result1, err
}

// CheckAndSetRedisRelationKey
// 返回0 KeyExistsAndNotSet 表示这个key存在，未设置
// 返回1 KeyUpdated	表示，这个key不存在,已更新
// 返回2 KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func CheckAndSetRedisRelationKey(userId uint, key string) int {
	if redis.IsExistUserSetField(userId, key) {
		return redis.KeyExistsAndNotSet
	}
	switch key {
	case redis.FollowList:
		if redis.AcquireRelationLock(userId, key) {
			defer redis.ReleaseRelationLock(userId, key)
			//double check
			if redis.IsExistUserSetField(userId, key) {
				return redis.KeyExistsAndNotSet
			}
			id, err := mysql.GetFollowList(userId)
			if err != nil {
				zap.L().Error("mysql获取FollowList失败", zap.Error(err))
				return redis.KeyNotExistsInBoth
			}
			if len(id) == 0 {
				zap.L().Info("mysql没有该FollowList")
				return redis.KeyNotExistsInBoth
			}
			err = redis.SetFollowListByUserId(userId, id)
			if err != nil {
				zap.L().Error("redis更新FollowList失败", zap.Error(err))
				return redis.KeyNotExistsInBoth
			}
			return redis.KeyUpdated
		}
	case redis.FollowerList:
		if redis.AcquireRelationLock(userId, key) {
			defer redis.ReleaseRelationLock(userId, key)
			//double check
			if redis.IsExistUserSetField(userId, key) {
				return redis.KeyExistsAndNotSet
			}
			id, err := mysql.GetFollowerList(userId)
			if err != nil {
				zap.L().Error("mysql获取FollowerList失败", zap.Error(err))
				return redis.KeyNotExistsInBoth
			}
			if len(id) == 0 {
				zap.L().Info("mysql没有该FollowerList")
				return redis.KeyNotExistsInBoth
			}
			err = redis.SetFollowerListByUserId(userId, id)
			if err != nil {
				zap.L().Error("redis更新FollowerList失败", zap.Error(err))
				return redis.KeyNotExistsInBoth
			}
			return redis.KeyUpdated
		}
	default:
		return redis.KeyNotExistsInBoth
	}
	fmt.Println("CheckAndSetRedisRelationKey")
	// 重试
	time.Sleep(redis.RetryTime)
	return CheckAndSetRedisRelationKey(userId, key)
}
