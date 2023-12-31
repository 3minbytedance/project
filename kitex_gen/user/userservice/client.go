// Code generated by Kitex v0.6.2. DO NOT EDIT.

package userservice

import (
	"context"
	user "douyin/kitex_gen/user"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Register(ctx context.Context, request *user.UserRegisterRequest, callOptions ...callopt.Option) (r *user.UserRegisterResponse, err error)
	Login(ctx context.Context, request *user.UserLoginRequest, callOptions ...callopt.Option) (r *user.UserLoginResponse, err error)
	GetUserInfoById(ctx context.Context, request *user.UserInfoByIdRequest, callOptions ...callopt.Option) (r *user.UserInfoByIdResponse, err error)
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
	return &kUserServiceClient{
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

type kUserServiceClient struct {
	*kClient
}

func (p *kUserServiceClient) Register(ctx context.Context, request *user.UserRegisterRequest, callOptions ...callopt.Option) (r *user.UserRegisterResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Register(ctx, request)
}

func (p *kUserServiceClient) Login(ctx context.Context, request *user.UserLoginRequest, callOptions ...callopt.Option) (r *user.UserLoginResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Login(ctx, request)
}

func (p *kUserServiceClient) GetUserInfoById(ctx context.Context, request *user.UserInfoByIdRequest, callOptions ...callopt.Option) (r *user.UserInfoByIdResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserInfoById(ctx, request)
}
