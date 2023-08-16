package main

import (
	"context"
	favorite "douyin/kitex_gen/favorite"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, request *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteList(ctx context.Context, request *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetVideoFavoriteCount implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetVideoFavoriteCount(ctx context.Context, request *favorite.VideoFavoriteCountRequest) (resp *favorite.VideoFavoriteCountResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserFavoriteCount implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetUserFavoriteCount(ctx context.Context, request *favorite.UserFavoriteCountRequest) (resp *favorite.UserFavoriteCountResponse, err error) {
	// TODO: Your code here...
	return
}
