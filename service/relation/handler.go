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
	userId := request.GetUserId()
	toUserId := request.GetToUserId()
	if userId == toUserId {
		return &relation.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("不能对自己操作"),
		}, nil
	}

	switch request.ActionType {
	case 1: // 关注
		// 延迟双删
		redis.DelKey(uint(userId), redis.FollowList)
		redis.DelKey(uint(toUserId), redis.FollowerList)

		// err = mysql.AddFollow(uint(userId), uint(toUserId))
		err = kafka.FollowMQInstance.ProduceAddFollowMsg(uint(userId), uint(toUserId))
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = thrift.StringPtr(err.Error())
			return
		}

		time.Sleep(1 * time.Second)

		redis.DelKey(uint(request.UserId), redis.FollowList)
		redis.DelKey(uint(request.ToUserId), redis.FollowerList)

		return &relation.RelationActionResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success"),
		}, nil
	case 2: // 取关
		// 延迟双删
		redis.DelKey(uint(request.UserId), redis.FollowList)
		redis.DelKey(uint(request.ToUserId), redis.FollowerList)

		// err = mysql.DeleteFollowById(uint(request.UserId), uint(request.ToUserId))
		err = kafka.FollowMQInstance.ProduceDelFollowMsg(uint(userId), uint(toUserId))
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = thrift.StringPtr(err.Error())
			return
		}

		time.Sleep(1 * time.Second)

		redis.DelKey(uint(request.UserId), redis.FollowList)
		redis.DelKey(uint(request.ToUserId), redis.FollowerList)

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
	res := CheckAndSetRedisRelationKey(uint(request.GetToUserId()), redis.FollowList)
	if res == 2 {
		return &relation.FollowListResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有关注用户"),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFollowListById(uint(request.GetToUserId()))
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
			ActorId: request.GetUserId(),
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
	res := CheckAndSetRedisRelationKey(uint(request.GetToUserId()), redis.FollowerList)
	if res == 2 {
		return &relation.FollowerListResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有粉丝用户"),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFollowerListById(uint(request.GetToUserId()))
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
			ActorId: request.GetUserId(),
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
	res := CheckAndSetRedisRelationKey(uint(request.GetToUserId()), redis.FollowList)
	if res == 2 {
		return &relation.FriendListResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有关注用户，所以没有friend"),
			UserList:   nil,
		}, nil
	}
	res = CheckAndSetRedisRelationKey(uint(request.GetToUserId()), redis.FollowerList)
	if res == 2 {
		return &relation.FriendListResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("success, 没有粉丝用户，所以没有friend"),
			UserList:   nil,
		}, nil
	}
	id, err := redis.GetFriendListById(uint(request.GetToUserId()))
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
			ActorId: request.GetUserId(),
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
	count, err := redis.GetFollowCountById(uint(userId))
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// GetFollowerListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerListCount(ctx context.Context, userId int64) (resp int32, err error) {
	res := CheckAndSetRedisRelationKey(uint(userId), redis.FollowerList)
	if res == 2 {
		return 0, nil
	}
	count, err := redis.GetFollowerCountById(uint(userId))
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// IsFollowing implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) IsFollowing(ctx context.Context, request *relation.IsFollowingRequest) (resp bool, err error) {
	actionId := request.GetActorId()
	toUserId := request.GetUserId()
	if actionId == toUserId {
		return true, nil
	}
	// redis存在key
	if redis.IsExistUserSetField(uint(actionId), redis.FollowList) {
		found := redis.IsInMyFollowList(uint(actionId), uint(toUserId))
		return found, nil
	}
	// redis不存在，从数据库查询是否已关注
	found := mysql.IsFollowing(uint(actionId), uint(toUserId))
	// 获取所有关注列表id
	followListId, err := mysql.GetFollowList(uint(actionId))
	if err != nil {
		zap.L().Error("GetFollowList error", zap.Error(err))
		return false, err
	}
	// 往redis赋值
	go func() {
		err = redis.SetFollowListByUserId(uint(actionId), followListId)
		if err != nil {
			zap.L().Error("SetFollowListByUserId error", zap.Error(err))
			return
		}
	}()
	return found, err
}

// IsFriend implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) IsFriend(ctx context.Context, request *relation.IsFriendRequest) (resp bool, err error) {
	// 从数据库查询是否已关注
	// todo 改为从redis取
	result, err := mysql.IsFriend(uint(request.ActorId), uint(request.UserId))
	return result, err
}

// CheckAndSetRedisRelationKey 返回0表示这个key存在，返回1表示，这个key不存在，已更新。返回2表示，这个key不存在，这个用户没有所对应的值
func CheckAndSetRedisRelationKey(userId uint, key string) int {
	if redis.IsExistUserSetField(userId, key) {
		return 0
	}
	//key不存在
	if key == redis.FollowList {
		id, err := mysql.GetFollowList(userId)
		if err != nil {
			zap.L().Error("mysql获取FollowList失败", zap.Error(err))
		}
		if len(id) == 0 {
			zap.L().Info("mysql获取FollowList为空，不更新redis")
			return 2
		}
		err = redis.SetFollowListByUserId(userId, id)
		if err != nil {
			zap.L().Error("redis更新FollowList失败", zap.Error(err))
		}
	} else {
		id, err := mysql.GetFollowerList(userId)
		if err != nil {
			zap.L().Error("mysql获取FollowerList失败", zap.Error(err))
		}
		if len(id) == 0 {
			zap.L().Info("mysql获取FollowerList为空，不更新redis")
			return 2
		}
		err = redis.SetFollowerListByUserId(userId, id)
		if err != nil {
			zap.L().Error("redis更新FollowerList失败", zap.Error(err))
		}
	}
	return 1
}
