// Code generated by Kitex v0.6.2. DO NOT EDIT.

package relationservice

import (
	"context"
	relation "douyin/kitex_gen/relation"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	RelationAction(ctx context.Context, request *relation.RelationActionRequest, callOptions ...callopt.Option) (r *relation.RelationActionResponse, err error)
	GetFollowList(ctx context.Context, request *relation.FollowListRequest, callOptions ...callopt.Option) (r *relation.FollowListResponse, err error)
	GetFollowerList(ctx context.Context, request *relation.FollowerListRequest, callOptions ...callopt.Option) (r *relation.FollowerListResponse, err error)
	GetFriendList(ctx context.Context, request *relation.FriendListRequest, callOptions ...callopt.Option) (r *relation.FriendListResponse, err error)
	GetFollowListCount(ctx context.Context, userId int64, callOptions ...callopt.Option) (r int32, err error)
	GetFollowerListCount(ctx context.Context, userId int64, callOptions ...callopt.Option) (r int32, err error)
	IsFollowing(ctx context.Context, request *relation.IsFollowingRequest, callOptions ...callopt.Option) (r bool, err error)
	IsFriend(ctx context.Context, request *relation.IsFriendRequest, callOptions ...callopt.Option) (r bool, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kRelationServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kRelationServiceClient struct {
	*kClient
}

func (p *kRelationServiceClient) RelationAction(ctx context.Context, request *relation.RelationActionRequest, callOptions ...callopt.Option) (r *relation.RelationActionResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.RelationAction(ctx, request)
}

func (p *kRelationServiceClient) GetFollowList(ctx context.Context, request *relation.FollowListRequest, callOptions ...callopt.Option) (r *relation.FollowListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFollowList(ctx, request)
}

func (p *kRelationServiceClient) GetFollowerList(ctx context.Context, request *relation.FollowerListRequest, callOptions ...callopt.Option) (r *relation.FollowerListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFollowerList(ctx, request)
}

func (p *kRelationServiceClient) GetFriendList(ctx context.Context, request *relation.FriendListRequest, callOptions ...callopt.Option) (r *relation.FriendListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFriendList(ctx, request)
}

func (p *kRelationServiceClient) GetFollowListCount(ctx context.Context, userId int64, callOptions ...callopt.Option) (r int32, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFollowListCount(ctx, userId)
}

func (p *kRelationServiceClient) GetFollowerListCount(ctx context.Context, userId int64, callOptions ...callopt.Option) (r int32, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFollowerListCount(ctx, userId)
}

func (p *kRelationServiceClient) IsFollowing(ctx context.Context, request *relation.IsFollowingRequest, callOptions ...callopt.Option) (r bool, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.IsFollowing(ctx, request)
}

func (p *kRelationServiceClient) IsFriend(ctx context.Context, request *relation.IsFriendRequest, callOptions ...callopt.Option) (r bool, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.IsFriend(ctx, request)
}
