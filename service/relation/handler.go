package main

import (
	"context"
	"douyin/constant"
	"douyin/dal/mysql"
	relation "douyin/kitex_gen/relation"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/mw/kafka"
	"douyin/mw/redis"
	"github.com/apache/thrift/lib/go/thrift"
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}))
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
				StatusCode: 0,
				StatusMsg:  thrift.StringPtr("success, 用户关注过了"),
			}, nil
		}
		// 延迟双删
		redis.DelKey(uint(fromUserId), redis.FollowList)
		redis.DelKey(uint(toUserId), redis.FollowerList)

		//err = mysql.AddFollow(uint(fromUserId), uint(toUserId))
		err = kafka.FollowMQInstance.ProduceAddFollowMsg(uint(fromUserId), uint(toUserId))
		if err != nil {
			return &relation.RelationActionResponse{
				StatusCode: 1,
				StatusMsg:  thrift.StringPtr("内部错误"),
			}, err
		}
		go func() {
			time.Sleep(1 * time.Second)

			redis.DelKey(uint(fromUserId), redis.FollowList)
			redis.DelKey(uint(toUserId), redis.FollowerList)
		}()

		return &relation.RelationActionResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success"),
		}, nil
	case 2: // 取关
		// 延迟双删
		redis.DelKey(uint(fromUserId), redis.FollowList)
		redis.DelKey(uint(toUserId), redis.FollowerList)

		//err = mysql.DeleteFollowById(uint(fromUserId), uint(toUserId))
		err = kafka.FollowMQInstance.ProduceDelFollowMsg(uint(fromUserId), uint(toUserId))
		if err != nil {
			return &relation.RelationActionResponse{
				StatusCode: 1,
				StatusMsg:  thrift.StringPtr("内部错误"),
			}, err
		}
		go func() {
			time.Sleep(1 * time.Second)

			redis.DelKey(uint(fromUserId), redis.FollowList)
			redis.DelKey(uint(toUserId), redis.FollowerList)
		}()

		return &relation.RelationActionResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success"),
		}, nil

	default:
		return &relation.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("参数不合法"),
		}, nil
	}
}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, request *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	actionId := request.GetUserId()
	toUserId := request.GetToUserId()

	res := CheckAndSetRedisRelationKey(uint(toUserId), redis.FollowList)
	if res == 2 {
		return &relation.FollowListResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有关注用户"),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFollowListById(uint(toUserId))
	if err != nil {
		return &relation.FollowListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取Follow list失败"),
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
				StatusCode: 1,
				StatusMsg:  thrift.StringPtr("查询Follow list用户信息失败"),
				UserList:   nil,
			}, err
		}
		followList = append(followList, userResp.GetUser())
	}
	return &relation.FollowListResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
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
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有粉丝用户"),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFollowerListById(uint(toUserId))
	if err != nil {
		return &relation.FollowerListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取Follow list失败"),
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
				StatusCode: 1,
				StatusMsg:  thrift.StringPtr("查询Follower list用户信息失败"),
				UserList:   nil,
			}, err
		}
		followerList = append(followerList, userResp.GetUser())
	}
	return &relation.FollowerListResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
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
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有关注用户，所以没有friend"),
			UserList:   nil,
		}, nil
	}
	res = CheckAndSetRedisRelationKey(uint(toUserId), redis.FollowerList)
	if res == 2 {
		return &relation.FriendListResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有粉丝用户，所以没有friend"),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFriendListById(uint(toUserId))
	if err != nil {
		return &relation.FriendListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取Friend list失败"),
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
				StatusCode: 1,
				StatusMsg:  thrift.StringPtr("查询Friend list用户信息失败"),
				UserList:   nil,
			}, err
		}
		friendList = append(friendList, userResp.User)
	}
	return &relation.FriendListResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
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
// 返回0表示这个key存在，返回1表示，这个key不存在，已更新。返回2表示，这个key不存在，这个用户没有所对应的值
func CheckAndSetRedisRelationKey(userId uint, key string) int {
	var m sync.RWMutex
	if redis.IsExistUserSetField(userId, key) {
		return 0
	}
	//key不存在 double check
	m.Lock()
	defer m.Unlock()
	if redis.IsExistUserSetField(userId, key) {
		return 0
	}
	switch key {
	case redis.FollowList:
		id, err := mysql.GetFollowList(userId)
		if err != nil || len(id) == 0 {
			zap.L().Error("mysql获取FollowList失败", zap.Error(err))
			return 2
		}
		err = redis.SetFollowListByUserId(userId, id)
		if err != nil {
			zap.L().Error("redis更新FollowList失败", zap.Error(err))
			return 2
		}
		return 1
	case redis.FollowerList:
		id, err := mysql.GetFollowerList(userId)
		if err != nil || len(id) == 0 {
			zap.L().Error("mysql获取FollowerList失败", zap.Error(err))
			return 2
		}
		err = redis.SetFollowerListByUserId(userId, id)
		if err != nil {
			zap.L().Error("redis更新FollowerList失败", zap.Error(err))
			return 2
		}
		return 1
	default:
		return 2
	}
}
