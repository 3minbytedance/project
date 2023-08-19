package main

import (
	"context"
	"douyin/constant"
	"douyin/dal/mysql"
	relation "douyin/kitex_gen/relation"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}))
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
		zap.Int32("user_id", request.UserId),
		zap.Int32("action_type", request.ActionType),
		zap.Int32("ToUserId", request.ToUserId),
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

		err = mysql.AddFollow(uint(userId), uint(toUserId))
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

		err = mysql.DeleteFollowById(uint(request.UserId), uint(request.ToUserId))
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
	CheckAndSetRedisRelationKey(uint(request.UserId), redis.FollowList)
	id, err := redis.GetFollowListById(uint(request.UserId))
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
			ActorId: request.UserId,
			UserId:  int32(com),
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
	CheckAndSetRedisRelationKey(uint(request.UserId), redis.FollowerList)
	id, err := redis.GetFollowerListById(uint(request.UserId))
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
			ActorId: request.UserId,
			UserId:  int32(com),
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
	CheckAndSetRedisRelationKey(uint(request.UserId), redis.FollowList)
	CheckAndSetRedisRelationKey(uint(request.UserId), redis.FollowerList)
	id, err := redis.GetFriendListById(uint(request.UserId))
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
			UserId:  int32(com),
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

// CheckAndSetRedisRelationKey 返回true表示不存在这个key，并设置key
// 返回false表示已存在这个key
func CheckAndSetRedisRelationKey(userId uint, key string) bool {
	if redis.IsExistUserSetField(userId, key) {
		return false
	}
	//key不存在
	if key == redis.FollowList {
		id, err := mysql.GetFollowList(userId)
		if err != nil {
			zap.L().Error("mysql获取FollowList失败", zap.Error(err))
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
		err = redis.SetFollowerListByUserId(userId, id)
		if err != nil {
			zap.L().Error("redis更新FollowerList失败", zap.Error(err))
		}
	}
	return true

}

// GetFollowListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowListCount(ctx context.Context, userId int32) (resp int32, err error) {
	CheckAndSetRedisRelationKey(uint(userId), redis.FollowList)
	count, err := redis.GetFollowCountById(uint(userId))
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// GetFollowerListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerListCount(ctx context.Context, userId int32) (resp int32, err error) {
	CheckAndSetRedisRelationKey(uint(userId), redis.FollowerList)
	count, err := redis.GetFollowerCountById(uint(userId))
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// IsFollowing implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) IsFollowing(ctx context.Context, request *relation.IsFollowingRequest) (resp bool, err error) {
	// redis存在key
	if redis.IsExistUserSetField(uint(request.ActorId), redis.FollowList) {
		found := redis.IsInMyFollowList(uint(request.ActorId), uint(request.UserId))
		return found, nil
	}
	// redis不存在，从数据库查询是否已关注
	found := mysql.IsFollowing(uint(request.ActorId), uint(request.UserId))
	// 获取所有关注列表id
	followListId, err := mysql.GetFollowList(uint(request.ActorId))
	if err != nil {
		zap.L().Error("GetFollowList error", zap.Error(err))
		return false, err
	}
	// 往redis赋值
	go func() {
		err = redis.SetFollowListByUserId(uint(request.ActorId), followListId)
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
