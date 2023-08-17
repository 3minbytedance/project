// Code generated by Kitex v0.6.2. DO NOT EDIT.

package favoriteservice

import (
	"context"
	favorite "douyin/kitex_gen/favorite"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	FavoriteAction(ctx context.Context, request *favorite.FavoriteActionRequest, callOptions ...callopt.Option) (r *favorite.FavoriteActionResponse, err error)
	GetFavoriteList(ctx context.Context, request *favorite.FavoriteListRequest, callOptions ...callopt.Option) (r *favorite.FavoriteListResponse, err error)
	GetVideoFavoriteCount(ctx context.Context, request *favorite.VideoFavoriteCountRequest, callOptions ...callopt.Option) (r *favorite.VideoFavoriteCountResponse, err error)
	GetUserFavoriteCount(ctx context.Context, request *favorite.UserFavoriteCountRequest, callOptions ...callopt.Option) (r *favorite.UserFavoriteCountResponse, err error)
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
	return &kFavoriteServiceClient{
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

type kFavoriteServiceClient struct {
	*kClient
}

func (p *kFavoriteServiceClient) FavoriteAction(ctx context.Context, request *favorite.FavoriteActionRequest, callOptions ...callopt.Option) (r *favorite.FavoriteActionResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.FavoriteAction(ctx, request)
}

func (p *kFavoriteServiceClient) GetFavoriteList(ctx context.Context, request *favorite.FavoriteListRequest, callOptions ...callopt.Option) (r *favorite.FavoriteListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFavoriteList(ctx, request)
}

func (p *kFavoriteServiceClient) GetVideoFavoriteCount(ctx context.Context, request *favorite.VideoFavoriteCountRequest, callOptions ...callopt.Option) (r *favorite.VideoFavoriteCountResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideoFavoriteCount(ctx, request)
}

func (p *kFavoriteServiceClient) GetUserFavoriteCount(ctx context.Context, request *favorite.UserFavoriteCountRequest, callOptions ...callopt.Option) (r *favorite.UserFavoriteCountResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserFavoriteCount(ctx, request)
}