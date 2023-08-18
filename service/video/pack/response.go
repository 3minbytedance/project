package pack

import (
	"douyin/constant/biz"
	"douyin/dal/model"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/video"
)

func Feed(videoModel *model.Video, userModel *user.User) *video.Video {
	if videoModel == nil {
		return nil
	}
	return &video.Video{
		Id:       int32(videoModel.ID),
		Author:   userModel,
		PlayUrl:  biz.OSS + videoModel.VideoUrl,
		CoverUrl: biz.OSS + videoModel.CoverUrl,
		//FavoriteCount: favoriteCount,
		//CommentCount:  commentCount,
		//IsFavorite:    IsUserFavorite(userId, video.ID),
		Title: videoModel.Title,
	}
}
