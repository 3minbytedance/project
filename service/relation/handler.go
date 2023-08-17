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

	switch request.ActionType {
	case 1: // 关注
		// 延迟双删
		redis.DelKey(uint(request.UserId), redis.FollowList)
		redis.DelKey(uint(request.ToUserId), redis.FollowerList)

		err = mysql.AddFollow(uint(request.UserId), uint(request.ToUserId))
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
		return
	}
}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, request *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	CheckAndSetRedisRealationKey(uint(request.UserId), redis.FollowList)
	id, err := redis.GetFollowListById(uint(request.UserId))
	if err != nil {
		return &relation.FollowListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取Follow list失败"),
			UserList:   nil,
		}, err
	}
	followlist := make([]*user.User, 0)
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
		followlist = append(followlist, userResp.User)
	}
	return &relation.FollowListResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
		UserList:   followlist,
	}, nil
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, request *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	CheckAndSetRedisRealationKey(uint(request.UserId), redis.FollowerList)
	id, err := redis.GetFollowerListById(uint(request.UserId))
	if err != nil {
		return &relation.FollowerListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取Follow list失败"),
			UserList:   nil,
		}, err
	}
	followerlist := make([]*user.User, 0)
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
		followerlist = append(followerlist, userResp.User)
	}
	return &relation.FollowerListResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
		UserList:   followerlist,
	}, nil
}

// GetFollowListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowListCount(ctx context.Context, request *relation.FollowListCountRequest) (resp *relation.FollowListCountResponse, err error) {
	CheckAndSetRedisRealationKey(uint(request.UserId), redis.FollowList)
	count, err := redis.GetFollowCountById(uint(request.UserId))
	if err != nil {
		return &relation.FollowListCountResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取follow list count失败"),
			Count:      0,
		}, err
	}
	return &relation.FollowListCountResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
		Count:      int32(count),
	}, nil
}

// GetFollowerListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerListCount(ctx context.Context, request *relation.FollowerListCountRequest) (resp *relation.FollowerListCountResponse, err error) {
	CheckAndSetRedisRealationKey(uint(request.UserId), redis.FollowerList)
	count, err := redis.GetFollowerCountById(uint(request.UserId))
	if err != nil {
		return &relation.FollowerListCountResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取follow list count失败"),
			Count:      0,
		}, err
	}
	return &relation.FollowerListCountResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
		Count:      int32(count),
	}, nil
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, request *relation.FriendListRequest) (resp *relation.FriendListResponse, err error) {
	CheckAndSetRedisRealationKey(uint(request.UserId), redis.FollowList)
	CheckAndSetRedisRealationKey(uint(request.UserId), redis.FollowerList)
	id, err := redis.GetFriendListById(uint(request.UserId))
	if err != nil {
		return &relation.FriendListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取Friend list失败"),
			UserList:   nil,
		}, err
	}
	friendlist := make([]*user.User, 0)
	for _, com := range id {
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			ActorId: request.UserId,
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
		friendlist = append(friendlist, userResp.User)
	}
	return &relation.FriendListResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
		UserList:   friendlist,
	}, nil
}

// IsFollowing implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) IsFollowing(ctx context.Context, request *relation.IsFollowingRequest) (resp *relation.IsFollowingResponse, err error) {
	// redis存在key
	if redis.IsExistUserSetField(uint(request.ActorId), redis.FollowList) {
		found := redis.IsInMyFollowList(uint(request.ActorId), uint(request.UserId))
		return &relation.IsFollowingResponse{Result_: found}, nil
	}
	// 从数据库查询是否已关注
	found := mysql.IsFollowing(uint(request.ActorId), uint(request.UserId))
	// 获取所有关注列表id
	followListId, err := mysql.GetFollowList(uint(request.ActorId))
	if err != nil {
		zap.L().Error("GetFollowList error", zap.Error(err))
		return &relation.IsFollowingResponse{Result_: found}, err
	}
	// 往redis赋值
	go func() {
		err = redis.SetFollowListByUserId(uint(request.ActorId), followListId)
		if err != nil {
			zap.L().Error("SetFollowListByUserId error", zap.Error(err))
			return
		}
	}()
	return &relation.IsFollowingResponse{Result_: found}, err
}

// 返回true表示不存在这个key，并设置key
// 返回false表示已存在这个key
func CheckAndSetRedisRealationKey(userId uint, key string) bool {
	if redis.IsExistUserSetField(userId, key) {
		return false
	} else {
		if key == redis.FollowList {
			id, err := mysql.GetFollowList(userId)
			if err != nil {
				log.Println("mysql获取FollowList失败", err)
			}
			err = redis.SetFollowListByUserId(userId, id)
			if err != nil {
				log.Println("redis更新FollowList失败", err)
			}
		} else {
			id, err := mysql.GetFollowerList(userId)
			if err != nil {
				log.Println("mysql获取FollowerList失败", err)
			}
			err = redis.SetFollowerListByUserId(userId, id)
			if err != nil {
				log.Println("redis更新FollowerList失败", err)
			}
		}
		return true
	}
}
