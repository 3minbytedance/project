package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/constant/biz"
	"douyin/dal/model"
	dalMySQL "douyin/dal/mysql"
	"douyin/kitex_gen/comment/commentservice"
	"douyin/kitex_gen/favorite"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/kitex_gen/video"
	"douyin/mw/kafka"
	mwRedis "douyin/mw/redis"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"sync/atomic"
	"time"
)

var (
	userClient    userservice.Client
	commentClient commentservice.Client
)

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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}),
		client.WithMuxConnection(1),
	)
	if err != nil {
		log.Fatal(err)
	}
	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
		client.WithMuxConnection(1),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, request *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	zap.L().Info("RelationClient action start",
		zap.Int64("user_id", request.UserId),
		zap.Int32("action_type", request.ActionType),
		zap.Int64("video_id", request.VideoId),
	)
	userId := uint(request.UserId)
	videoId := uint(request.VideoId)
	actionType := int(request.ActionType)
	res := favoriteActions(userId, videoId, actionType)
	switch res {
	case 0:
		return &favorite.FavoriteActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		}, nil
	case 1:
		return &favorite.FavoriteActionResponse{
			StatusCode: common.CodeFavoriteRepeat,
			StatusMsg:  common.MapErrMsg(common.CodeFavoriteRepeat),
		}, nil
	default:
		return &favorite.FavoriteActionResponse{
			StatusCode: common.CodeServerBusy,
			StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
		}, nil
	}
}

// GetFavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteList(ctx context.Context, request *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	resp = new(favorite.FavoriteListResponse)
	userId := request.GetUserId()
	actionId := request.GetActionId()

	favoritesByUserId, err := getFavoritesByUserId(uint(userId))
	if err != nil {
		return &favorite.FavoriteListResponse{
			StatusCode: common.CodeServerBusy,
			StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
		}, err
	}
	videos := make([]*video.Video, 0, len(favoritesByUserId))
	for _, id := range favoritesByUserId {
		videoModel, found := dalMySQL.FindVideoByVideoId(id)
		if !found {
			continue
		}
		userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			ActorId: actionId,
			UserId:  int64(videoModel.AuthorId),
		})
		commentCount, _ := commentClient.GetCommentCount(ctx, int64(id))
		favoriteCount, _ := getFavoritesVideoCount(id)
		vid := video.Video{
			Id:            int64(videoModel.ID),
			Author:        userResp.GetUser(),
			PlayUrl:       biz.OSS + videoModel.VideoUrl,
			CoverUrl:      biz.OSS + videoModel.CoverUrl,
			FavoriteCount: int32(favoriteCount),
			CommentCount:  commentCount,
			//判断当前请求ID是否点赞该视频
			IsFavorite: isUserFavorite(uint(actionId), id),
			Title:      videoModel.Title,
		}
		videos = append(videos, &vid)
	}
	return &favorite.FavoriteListResponse{
		StatusCode: common.CodeSuccess,
		StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		VideoList:  videos,
	}, nil
}

// GetVideoFavoriteCount implements the FavoriteServiceImpl interface.
// 获取视频的点赞数
func (s *FavoriteServiceImpl) GetVideoFavoriteCount(ctx context.Context, videoId int64) (resp int32, err error) {
	count, err := getFavoritesVideoCount(uint(videoId))
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// GetUserFavoriteCount implements the FavoriteServiceImpl interface.
// 获取用户喜欢的视频列表数
func (s *FavoriteServiceImpl) GetUserFavoriteCount(ctx context.Context, userId int64) (resp int32, err error) {
	res := checkAndSetUserFavoriteListKey(uint(userId), mwRedis.FavoriteList)
	// redis和mysql中没有对应的数据
	if res == 2 {
		return 0, nil
	}
	count, err := mwRedis.GetUserFavoriteVideoCountById(uint(userId))
	if err != nil {
		return 0, nil
	}
	return int32(count), nil
}

// GetUserTotalFavoritedCount implements the FavoriteServiceImpl interface.
// 获取用户的被喜欢总数
func (s *FavoriteServiceImpl) GetUserTotalFavoritedCount(ctx context.Context, userId int64) (resp int32, err error) {
	favoritesByUserId, err := getUserTotalFavoritedCount(uint(userId))
	if err != nil {
		return 0, err
	}
	return int32(favoritesByUserId), nil
}

// IsUserFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) IsUserFavorite(ctx context.Context, request *favorite.IsUserFavoriteRequest) (resp bool, err error) {
	userId := request.GetUserId()
	videoId := request.GetVideoId()
	return isUserFavorite(uint(userId), uint(videoId)), nil
}

// favoriteActions 点赞，取消赞的操作过程
// magic number 待改 todo
// 0 表示点赞/取消点赞成功
// 1 表示重复点赞/取消点赞
// 2 表示其他错误
func favoriteActions(userId uint, videoId uint, actionType int) int {
	videoModel, found := dalMySQL.FindVideoByVideoId(videoId)
	if !found {
		return 2
	}
	// 判断是否在redis中，防止对空key操作
	checkAndSetUserFavoriteListKey(userId, mwRedis.FavoriteList)
	checkAndSetVideoFavoriteCountKey(videoId, mwRedis.VideoFavoritedCountField)
	checkAndSetTotalFavoriteFieldKey(videoModel.AuthorId, mwRedis.TotalFavoriteField)

	switch actionType {
	case 1:
		// 点赞
		// 判断重复点赞
		if isUserFavorite(userId, videoId) {
			return 1
		}
		if mwRedis.AcquireFavoriteLock(userId, mwRedis.FavoriteAction) {
			defer mwRedis.ReleaseFavoriteLock(userId, mwRedis.FavoriteAction)
			if isUserFavorite(userId, videoId) {
				return 1
			}
			err := mwRedis.ActionLike(userId, videoId, videoModel.AuthorId)
			if err != nil {
				return 2
			}
			go func() {
				err := kafka.FavoriteMQInstance.ProduceAddFavoriteMsg(userId, videoId)
				if err != nil {
					zap.L().Error("更新MySQL点赞表err", zap.Error(err))
					return
				}
			}()
			return 0
		}
		time.Sleep(mwRedis.RetryTime)
		return favoriteActions(userId, videoId, actionType)
	case 2:
		if !isUserFavorite(userId, videoId) {
			return 1
		}
		if mwRedis.AcquireFavoriteLock(userId, mwRedis.FavoriteAction) {
			defer mwRedis.ReleaseFavoriteLock(userId, mwRedis.FavoriteAction)
			if !isUserFavorite(userId, videoId) {
				return 1
			}
			err := mwRedis.ActionCancelLike(userId, videoId, videoModel.AuthorId)
			if err != nil {
				return 2
			}
			go func() {
				err = kafka.FavoriteMQInstance.ProduceDelFavoriteMsg(userId, videoId)
				if err != nil {
					zap.L().Error("更新MySQL点赞表err", zap.Error(err))
					return
				}
			}()
			return 0
		}
		time.Sleep(mwRedis.RetryTime)
		return favoriteActions(userId, videoId, actionType)
	default:
		return 2
	}
}

// GetUserTotalFavoritedCount 获取用户发布视频的总的被点赞数量
func getUserTotalFavoritedCount(userId uint) (int64, error) {
	res := checkAndSetTotalFavoriteFieldKey(userId, mwRedis.TotalFavoriteField)
	if res == 2 {
		return 0, nil
	}
	count, err := mwRedis.GetTotalFavoritedByUserId(userId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetFavoritesVideoCount 根据视频id，返回该视频的点赞数
func getFavoritesVideoCount(videoId uint) (int64, error) {
	// 判断redis中是否存在对应的video数据
	res := checkAndSetVideoFavoriteCountKey(videoId, mwRedis.VideoFavoritedCountField)
	// redis和mysql中没有对应的数据
	if res == 2 {
		return 0, nil
	}
	count, err := mwRedis.GetFavoritedCountByVideoId(videoId)
	if err != nil {
		zap.L().Error("GetFavoritedCountByVideoId失败", zap.Error(err))
		return 0, err
	}
	return count, nil
}

// getFavoritesByUserId
// 获取当前user_id的点赞的视频id列表
func getFavoritesByUserId(userId uint) ([]uint, error) {
	res := checkAndSetUserFavoriteListKey(userId, mwRedis.FavoriteList)
	// redis和mysql中没有对应的数据
	if res == 2 {
		return []uint{}, nil
	}

	// redis存在
	favoritesVideoIdList, err := mwRedis.GetFavoriteListByUserId(userId)
	if err != nil {
		zap.L().Error("GetFavoriteListByUserId", zap.Error(err))
		return []uint{}, err
	}
	return favoritesVideoIdList, nil
}

// getIdListFromFavoriteSlice 从Favorite的slice中获取id的列表
func getIdListFromFavoriteSlice(favorites []model.Favorite, idType int) []uint {
	res := make([]uint, 0, len(favorites))
	for _, fav := range favorites {
		switch idType {
		case dalMySQL.IdTypeUser:
			res = append(res, fav.VideoId)
		case dalMySQL.IdTypeVideo:
			res = append(res, fav.UserId)
		}
	}
	return res
}

// IsUserFavorite 判断是否点赞
func isUserFavorite(userId, videoId uint) bool {
	res := checkAndSetUserFavoriteListKey(userId, mwRedis.FavoriteList)
	// redis和mysql中没有对应的数据
	if res == 2 {
		return false
	}
	return mwRedis.IsInUserFavoriteList(userId, videoId)
}

// checkAndSetUserFavoriteListKey
// 返回0 mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回1 mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回2 mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetUserFavoriteListKey(userId uint, key string) int {
	if mwRedis.IsExistUserSetField(userId, key) {
		return mwRedis.KeyExistsAndNotSet
	}
	//key不存在 double check
	if mwRedis.AcquireFavoriteLock(userId, mwRedis.FavoriteList) {
		defer mwRedis.ReleaseFavoriteLock(userId, mwRedis.FavoriteList)
		//double check
		if mwRedis.IsExistUserSetField(userId, key) {
			return mwRedis.KeyExistsAndNotSet
		}
		favorites, favoriteLength, err := dalMySQL.GetFavoritesByIdFromMysql(userId, dalMySQL.IdTypeUser)
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
			return mwRedis.KeyNotExistsInBoth
		}
		//点赞数为0
		if favoriteLength == 0 {
			return mwRedis.KeyNotExistsInBoth
		}
		idList := getIdListFromFavoriteSlice(favorites, dalMySQL.IdTypeUser)
		// key 不存在需要同步到redis
		err = mwRedis.SetFavoriteListByUserId(userId, idList) // 加载到set中
		if err != nil {
			zap.L().Error("SetFavoriteListByUserId", zap.Error(err))
			return mwRedis.KeyNotExistsInBoth
		}
		return mwRedis.KeyUpdated
	}
	fmt.Println("checkAndSetUserFavoriteListKey")
	// 重试
	time.Sleep(mwRedis.RetryTime)
	return checkAndSetTotalFavoriteFieldKey(userId, key)
}

// checkAndSetVideoFavoriteCountKey
// 返回0 mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回1 mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回2 mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetVideoFavoriteCountKey(videoId uint, key string) int {
	if mwRedis.IsExistVideoField(videoId, key) {
		return mwRedis.KeyExistsAndNotSet
	}
	//key不存在
	if mwRedis.AcquireFavoriteLock(videoId, mwRedis.VideoFavoritedCountField) {
		defer mwRedis.ReleaseFavoriteLock(videoId, mwRedis.VideoFavoritedCountField)
		//double check
		if mwRedis.IsExistVideoField(videoId, key) {
			return mwRedis.KeyExistsAndNotSet
		}
		// redis中不存在，从数据库中读取
		_, num, err := dalMySQL.GetFavoritesByIdFromMysql(videoId, dalMySQL.IdTypeVideo)
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
			return mwRedis.KeyNotExistsInBoth
		}
		if num == 0 {
			return mwRedis.KeyNotExistsInBoth
		}
		err = mwRedis.SetVideoFavoritedCountByVideoId(videoId, int64(num)) // 加载视频被点赞数量
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
			return mwRedis.KeyNotExistsInBoth
		}
		return mwRedis.KeyUpdated
	}
	fmt.Println("checkAndSetVideoFavoriteCountKey")
	time.Sleep(mwRedis.RetryTime)
	return checkAndSetTotalFavoriteFieldKey(videoId, key)
}

// checkAndSetTotalFavoriteFieldKey
// 返回0 mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回1 mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回2 mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetTotalFavoriteFieldKey(userId uint, key string) int {
	if mwRedis.IsExistUserField(userId, key) {
		return mwRedis.KeyExistsAndNotSet
	}
	//key不存在 double check
	if mwRedis.AcquireFavoriteLock(userId, mwRedis.TotalFavoriteField) {
		defer mwRedis.ReleaseFavoriteLock(userId, mwRedis.TotalFavoriteField)
		//double check
		if mwRedis.IsExistUserField(userId, key) {
			return mwRedis.KeyExistsAndNotSet
		}
		//redis 不存在
		var total int64
		// 获取用户发布的视频列表
		videosByAuthorId, exist := dalMySQL.FindVideosByAuthorId(userId)
		if !exist {
			return mwRedis.KeyNotExistsInBoth
		}
		for _, videoModel := range videosByAuthorId {
			count, _ := getFavoritesVideoCount(videoModel.ID)
			atomic.AddInt64(&total, count)
		}
		err := mwRedis.SetTotalFavoritedByUserId(userId, total)
		if err != nil {
			zap.L().Error("SetTotalFavoriteByUserId失败", zap.Error(err))
			return mwRedis.KeyNotExistsInBoth
		}
		return mwRedis.KeyUpdated
	}
	fmt.Println("checkAndSetTotalFavoriteFieldKey")
	time.Sleep(mwRedis.RetryTime)
	return checkAndSetTotalFavoriteFieldKey(userId, key)
}
