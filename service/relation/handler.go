package main

import (
	"context"
	"douyin/dal/mysql"
	relation "douyin/kitex_gen/relation"
	"douyin/mw/redis"
	"go.uber.org/zap"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, request *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, request *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, request *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowListCount(ctx context.Context, request *relation.FollowListCountRequest) (resp *relation.FollowListCountResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowerListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerListCount(ctx context.Context, request *relation.FollowerListCountRequest) (resp *relation.FollowerListCountResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, request *relation.FriendListRequest) (resp *relation.FriendListResponse, err error) {
	// TODO: Your code here...
	return
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
