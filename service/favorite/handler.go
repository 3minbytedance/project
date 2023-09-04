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
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"strconv"
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
		client.WithMuxConnection(2),
	)
	if err != nil {
		log.Fatal(err)
	}
	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
		client.WithMuxConnection(2),
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
	case biz.FavoriteActionSuccess:
		return &favorite.FavoriteActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		}, nil
	case biz.FavoriteActionRepeat:
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
	userRespCh := make(chan *user.UserInfoByIdResponse)
	commentCountCh := make(chan int32)
	favoriteCountCh := make(chan int32)
	isFavoriteCh := make(chan bool)
	defer func() {
		close(commentCountCh)
		close(favoriteCountCh)
		close(isFavoriteCh)
		close(userRespCh)
	}()
	for _, id := range favoritesByUserId {
		videoModel, found := dalMySQL.FindVideoByVideoId(id)
		if !found {
			continue
		}
		go func() {
			userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
				ActorId: actionId,
				UserId:  int64(videoModel.AuthorId),
			})
			userRespCh <- userResp
		}()

		go func() {
			commentCount, _ := commentClient.GetCommentCount(ctx, int64(id))
			commentCountCh <- commentCount
		}()

		go func() {
			favoriteCount, _ := getFavoritesVideoCount(id)
			favoriteCountCh <- int32(favoriteCount)
		}()

		go func() {
			//判断当前请求ID是否点赞该视频
			isFavorite := isUserFavorite(uint(actionId), id)
			isFavoriteCh <- isFavorite
		}()

		vid := video.Video{
			Id:       int64(videoModel.ID),
			PlayUrl:  biz.OSS + videoModel.VideoUrl,
			CoverUrl: biz.OSS + videoModel.CoverUrl,
			Title: videoModel.Title,
		}

		for receivedCount := 0; receivedCount < 4; receivedCount++ {
			select {
			case userResp := <-userRespCh:
				vid.SetAuthor(userResp.GetUser())
			case favoriteCount := <-favoriteCountCh:
				vid.SetFavoriteCount(favoriteCount)
			case isFavorite := <-isFavoriteCh:
				vid.SetIsFavorite(isFavorite)
			case commentCount := <-commentCountCh:
				vid.SetCommentCount(commentCount)
			case <-time.After(3 * time.Second):
				zap.L().Error("3s overtime.")
				break
			}
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
	if res == mwRedis.KeyNotExistsInBoth {
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
	if userId == 0 {
		return false, nil
	}
	videoId := request.GetVideoId()
	return isUserFavorite(uint(userId), uint(videoId)), nil
}

// favoriteActions 点赞，取消赞的操作过程
// FavoriteActionSuccess 表示点赞/取消点赞成功
// FavoriteActionRepeat 表示重复点赞/取消点赞
// FavoriteActionError 表示其他错误
func favoriteActions(userId uint, videoId uint, actionType int) (status int) {
	videoModel, found := dalMySQL.FindVideoByVideoId(videoId)
	if !found {
		return biz.FavoriteActionError
	}
	switch actionType {
	case 1:
		// 点赞
		if mwRedis.AcquireFavoriteLock(userId, mwRedis.FavoriteAction) {
			defer mwRedis.ReleaseFavoriteLock(userId, mwRedis.FavoriteAction)
			// 判断是否在redis中，防止对空key操作
			checkAndSetUserFavoriteListKey(userId, mwRedis.FavoriteList)
			checkAndSetVideoFavoriteCountKey(videoId, mwRedis.VideoFavoritedCountField)
			// 判断重复点赞
			if isUserFavorite(userId, videoId) {
				return biz.FavoriteActionRepeat
			}
			err := mwRedis.ActionLike(userId, videoId, videoModel.AuthorId)
			if err != nil {
				return biz.FavoriteActionError
			}
			go func() {
				err := kafka.FavoriteMQInstance.ProduceAddFavoriteMsg(userId, videoId)
				if err != nil {
					zap.L().Error("更新MySQL点赞表err", zap.Error(err))
					return
				}
			}()
			go func() {
				common.AddToFavoriteUserIdBloom(strconv.Itoa(int(userId)))
				common.AddToFavoriteVideoIdBloom(strconv.Itoa(int(videoId)))
			}()
			return biz.FavoriteActionSuccess
		}
		time.Sleep(mwRedis.RetryTime)
		return favoriteActions(userId, videoId, actionType)
	case 2:
		if mwRedis.AcquireFavoriteLock(userId, mwRedis.FavoriteAction) {
			defer mwRedis.ReleaseFavoriteLock(userId, mwRedis.FavoriteAction)
			// 判断是否在redis中，防止对空key操作
			checkAndSetUserFavoriteListKey(userId, mwRedis.FavoriteList)
			checkAndSetVideoFavoriteCountKey(videoId, mwRedis.VideoFavoritedCountField)
			if !isUserFavorite(userId, videoId) {
				return biz.FavoriteActionRepeat
			}
			err := mwRedis.ActionCancelLike(userId, videoId, videoModel.AuthorId)
			if err != nil {
				return biz.FavoriteActionError
			}
			go func() {
				err = kafka.FavoriteMQInstance.ProduceDelFavoriteMsg(userId, videoId)
				if err != nil {
					zap.L().Error("更新MySQL点赞表err", zap.Error(err))
					return
				}
			}()
			return biz.FavoriteActionSuccess
		}
		time.Sleep(mwRedis.RetryTime)
		return favoriteActions(userId, videoId, actionType)
	default:
		return biz.FavoriteActionError
	}
}

// GetUserTotalFavoritedCount 获取用户发布视频的总的被点赞数量
func getUserTotalFavoritedCount(userId uint) (int64, error) {
	totalFavoriteCount, _ := checkAndSetTotalFavoriteFieldKey(userId, mwRedis.TotalFavoriteField)
	return totalFavoriteCount, nil
}

// GetFavoritesVideoCount 根据视频id，返回该视频的点赞数
func getFavoritesVideoCount(videoId uint) (int64, error) {
	// 判断redis中是否存在对应的video数据
	count, _ := checkAndSetVideoFavoriteCountKey(videoId, mwRedis.VideoFavoritedCountField)
	return count, nil
}

// getFavoritesByUserId
// 获取当前user_id的点赞的视频id列表
func getFavoritesByUserId(userId uint) ([]uint, error) {
	res := checkAndSetUserFavoriteListKey(userId, mwRedis.FavoriteList)
	// redis和mysql中没有对应的数据
	if res == mwRedis.KeyNotExistsInBoth {
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
	if res == mwRedis.KeyNotExistsInBoth {
		return false
	}
	return mwRedis.IsInUserFavoriteList(userId, videoId)
}

// checkAndSetUserFavoriteListKey
// 返回mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
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

		exist := common.TestFavoriteUserIdBloom(strconv.Itoa(int(userId)))

		// 不存在
		if !exist {
			return mwRedis.KeyNotExistsInBoth
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
		// key 不存在需要同步到redis set中
		err = mwRedis.SetFavoriteListByUserId(userId, idList)
		if err != nil {
			zap.L().Error("SetFavoriteListByUserId", zap.Error(err))
		}
		return mwRedis.KeyUpdated
	}
	// 重试
	time.Sleep(mwRedis.RetryTime)
	return checkAndSetUserFavoriteListKey(userId, key)
}

// checkAndSetVideoFavoriteCountKey
// 返回mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetVideoFavoriteCountKey(videoId uint, key string) (videoFavoriteCount int64, status int) {
	if count, err := mwRedis.GetFavoritedCountByVideoId(videoId); err == nil {
		return count, mwRedis.KeyExistsAndNotSet
	}
	//key不存在
	if mwRedis.AcquireFavoriteLock(videoId, mwRedis.VideoFavoritedCountField) {
		defer mwRedis.ReleaseFavoriteLock(videoId, mwRedis.VideoFavoritedCountField)
		//double check
		if count, err := mwRedis.GetFavoritedCountByVideoId(videoId); err == nil {
			return count, mwRedis.KeyExistsAndNotSet
		}
		// redis中不存在，从数据库中读取

		exist := common.TestFavoriteVideoIdBloom(strconv.Itoa(int(videoId)))

		// 不存在
		if !exist {
			mwRedis.SetVideoFavoritedCountByVideoId(videoId, 0)
			return 0, mwRedis.KeyNotExistsInBoth
		}

		_, num, err := dalMySQL.GetFavoritesByIdFromMysql(videoId, dalMySQL.IdTypeVideo)
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
			return 0, mwRedis.KeyNotExistsInBoth
		}
		err = mwRedis.SetVideoFavoritedCountByVideoId(videoId, int64(num)) // 加载视频被点赞数量
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
		}
		return int64(num), mwRedis.KeyUpdated
	}
	time.Sleep(mwRedis.RetryTime)
	return checkAndSetVideoFavoriteCountKey(videoId, key)
}

// checkAndSetTotalFavoriteFieldKey
// 获取userId发布的视频的总的获赞数量
// 返回mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetTotalFavoriteFieldKey(userId uint, key string) (totalFavoriteCount int64, status int) {
	if totalCount, err := mwRedis.GetTotalFavoritedByUserId(userId); err == nil {
		return totalCount, mwRedis.KeyExistsAndNotSet
	}
	//key不存在 double check
	if mwRedis.AcquireFavoriteLock(userId, mwRedis.TotalFavoriteField) {
		defer mwRedis.ReleaseFavoriteLock(userId, mwRedis.TotalFavoriteField)
		//double check
		if totalCount, err := mwRedis.GetTotalFavoritedByUserId(userId); err == nil {
			return totalCount, mwRedis.KeyExistsAndNotSet
		}

		var total int64
		// 获取用户发布的视频列表
		videosByAuthorId, exist := dalMySQL.FindVideosByAuthorId(userId)
		if !exist {
			mwRedis.SetTotalFavoritedByUserId(userId, 0)
			return 0, mwRedis.KeyNotExistsInBoth
		}
		for _, videoModel := range videosByAuthorId {
			count, _ := getFavoritesVideoCount(videoModel.ID)
			atomic.AddInt64(&total, count)
		}
		err := mwRedis.SetTotalFavoritedByUserId(userId, total)
		if err != nil {
			zap.L().Error("SetTotalFavoriteByUserId失败", zap.Error(err))
		}
		return total, mwRedis.KeyUpdated
	}
	time.Sleep(mwRedis.RetryTime)
	return checkAndSetTotalFavoriteFieldKey(userId, key)
}
