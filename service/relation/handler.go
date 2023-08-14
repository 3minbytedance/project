package main

import (
	"context"
	relation "douyin/kitex_gen/relation"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, req *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowListCount(ctx context.Context, req *relation.FollowListCountRequest) (resp *relation.FollowListCountResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowerListCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerListCount(ctx context.Context, req *relation.FollowerListCountRequest) (resp *relation.FollowerListCountResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, req *relation.FriendListRequest) (resp *relation.FriendListResponse, err error) {
	// TODO: Your code here...
	return
}

// IsFollowed implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) IsFollowed(ctx context.Context, req *relation.IsFollowedRequest) (resp *relation.IsFollowedResponse, err error) {
	// TODO: Your code here...
	return
}
