package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/constant/biz"
	"douyin/dal/model"
	"douyin/dal/mysql"
	"douyin/kitex_gen/comment/commentservice"
	"douyin/kitex_gen/favorite"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/kitex_gen/video"
	"douyin/mw/localcache"
	"douyin/mw/redis"
	"douyin/mw/rocketMQ"
	"github.com/allegro/bigcache/v3"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type favoriteMap = map[uint]map[uint]int

var (
	userClient    userservice.Client
	commentClient commentservice.Client

	favoriteData = make(favoriteMap)
	flushMutex   = sync.RWMutex{}
	mutex        = sync.Mutex{}
	mapChan      = make(chan favoriteMap)

	cache *bigcache.BigCache
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

	cache = localcache.Init(localcache.FavoriteVideo)

	go startTimer(mapChan)
	go consumerFavoriteMap(mapChan)
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
	switch actionType {
	case 1, 2:
		err = rocketMQ.FavoriteMQInstance.ProduceFavoriteMsg(&model.FavoriteAction{
			UserId:     userId,
			VideoId:    videoId,
			ActionType: actionType,
		})
		if err != nil {
			return &favorite.FavoriteActionResponse{
				StatusCode: common.CodeServerBusy,
				StatusMsg:  common.MapErrMsg(common.CodeServerBusy),
			}, nil
		}

		return &favorite.FavoriteActionResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		}, nil
	}
	return &favorite.FavoriteActionResponse{
		StatusCode: common.CodeInvalidParam,
		StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
	}, nil
}

// GetFavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteList(ctx context.Context, request *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	resp = new(favorite.FavoriteListResponse)
	actionId := request.GetActionId()
	userId := request.GetUserId()

	favoritesByUserId := mysql.GetFavoritesById(uint(userId))
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
		videoModel, found := getVideoByVideoId(id)
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
			favoriteCount, _ := checkAndSetVideoFavoriteCountKey(id)
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
			Title:    videoModel.Title,
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
	count, _ := checkAndSetVideoFavoriteCountKey(uint(videoId))
	return int32(count), nil
}

// GetUserFavoriteCount implements the FavoriteServiceImpl interface.
// 获取用户喜欢的视频列表数
func (s *FavoriteServiceImpl) GetUserFavoriteCount(ctx context.Context, userId int64) (resp int32, err error) {
	favoriteCount, _ := checkAndSetUserFavoriteCountKey(uint(userId))
	return int32(favoriteCount), nil
}

// GetUserTotalFavoritedCount implements the FavoriteServiceImpl interface.
// 获取用户的被喜欢总数
func (s *FavoriteServiceImpl) GetUserTotalFavoritedCount(ctx context.Context, userId int64) (resp int32, err error) {
	totalFavoriteCount, _ := checkAndSetTotalFavoriteFieldKey(uint(userId))
	if err != nil {
		return 0, err
	}
	return int32(totalFavoriteCount), nil
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

func addFavoriteActionToRedis(userId uint, videoId uint, actionType int, authorId uint) (status int) {
	switch actionType {
	case 1:
		// 判断是否在redis中，如果不在则从MySQL取防止对空key操作
		checkAndSetVideoFavoriteCountKey(videoId)
		checkAndSetTotalFavoriteFieldKey(authorId)
		checkAndSetUserFavoriteCountKey(userId)
		err := redis.ActionLike(userId, videoId, authorId)
		if err != nil {
			return biz.FavoriteActionError
		}
		go func() {
			common.AddToIsFavoriteBloom(userId, videoId)
			common.AddToFavoriteVideoIdBloom(strconv.Itoa(int(videoId)))
		}()
		return biz.FavoriteActionSuccess
	case 2:
		// 判断是否在redis中，如果不在则从MySQL取防止对空key操作
		checkAndSetVideoFavoriteCountKey(videoId)
		checkAndSetTotalFavoriteFieldKey(authorId)
		checkAndSetUserFavoriteCountKey(userId)
		err := redis.ActionCancelLike(userId, videoId, authorId)
		if err != nil {
			return biz.FavoriteActionError
		}
		return biz.FavoriteActionSuccess
	default:
		return biz.FavoriteActionError
	}
}

// IsUserFavorite 判断是否点赞
func isUserFavorite(userId, videoId uint) bool {
	exist := common.TestIsFavoriteBloom(userId, videoId)
	if !exist {
		return false
	}
	return mysql.IsFavorite(userId, videoId)
}

// checkAndSetVideoFavoriteCountKey
// 返回mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetVideoFavoriteCountKey(videoId uint) (videoFavoriteCount int64, status int) {
	if count, err := redis.GetFavoritedCountByVideoId(videoId); err == nil {
		return count, redis.KeyExistsAndNotSet
	}
	//key不存在
	if redis.AcquireFavoriteLock(videoId, redis.VideoFavoritedCountField) {
		defer redis.ReleaseFavoriteLock(videoId, redis.VideoFavoritedCountField)
		//double check
		if count, err := redis.GetFavoritedCountByVideoId(videoId); err == nil {
			return count, redis.KeyExistsAndNotSet
		}
		// redis中不存在，从数据库中读取

		exist := common.TestFavoriteVideoIdBloom(strconv.Itoa(int(videoId)))

		// 不存在
		if !exist {
			redis.SetVideoFavoritedCountByVideoId(videoId, 0)
			return 0, redis.KeyNotExistsInBoth
		}
		num, err := mysql.GetVideoFavoriteCountByVideoId(videoId)
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
			return 0, redis.KeyNotExistsInBoth
		}
		redis.SetVideoFavoritedCountByVideoId(videoId, num) // 加载视频被点赞数量
		return num, redis.KeyUpdated
	}
	time.Sleep(redis.RetryTime)
	return checkAndSetVideoFavoriteCountKey(videoId)
}

// checkAndSetTotalFavoriteFieldKey
// 获取userId发布的视频的总的获赞数量
// 返回mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetTotalFavoriteFieldKey(userId uint) (totalFavoriteCount int64, status int) {
	if totalCount, err := redis.GetTotalFavoritedByUserId(userId); err == nil {
		return totalCount, redis.KeyExistsAndNotSet
	}
	//key不存在 double check
	if redis.AcquireFavoriteLock(userId, redis.TotalFavoriteField) {
		defer redis.ReleaseFavoriteLock(userId, redis.TotalFavoriteField)
		//double check
		if totalCount, err := redis.GetTotalFavoritedByUserId(userId); err == nil {
			return totalCount, redis.KeyExistsAndNotSet
		}

		var total int64
		// 获取用户发布的视频列表
		videosByAuthorId, exist := mysql.FindVideosByAuthorId(userId)
		if !exist {
			redis.SetTotalFavoritedByUserId(userId, 0)
			return 0, redis.KeyNotExistsInBoth
		}
		for _, videoModel := range videosByAuthorId {
			count, err := mysql.GetVideoFavoriteCountByVideoId(videoModel.ID)
			if err != nil {
				continue
			}
			atomic.AddInt64(&total, count)
		}
		redis.SetTotalFavoritedByUserId(userId, total)
		return total, redis.KeyUpdated
	}
	time.Sleep(redis.RetryTime)
	return checkAndSetTotalFavoriteFieldKey(userId)
}

// checkAndSetUserFavoriteCountKey
// 返回mwRedis.KeyExistsAndNotSet 表示这个key存在，未设置
// 返回mwRedis.KeyUpdated 表示，这个key不存在,已更新
// 返回mwRedis.KeyNotExistsInBoth 表示，这个key在数据库和redis中都不存在，即缓存穿透
func checkAndSetUserFavoriteCountKey(userId uint) (int64, int) {
	if count, err := redis.GetUserFavoriteVideoCountById(userId); err == nil {
		return count, redis.KeyExistsAndNotSet
	}
	//key不存在 double check
	if redis.AcquireFavoriteLock(userId, redis.FavoriteCountFiled) {
		defer redis.ReleaseFavoriteLock(userId, redis.FavoriteCountFiled)

		favoriteCount, err := mysql.GetUserFavoriteCount(userId)
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
			return 0, redis.KeyNotExistsInBoth
		}
		// 同步到redis中
		redis.SetUserFavoriteVideoCountById(userId, favoriteCount)
		if err != nil {
			zap.L().Error("SetFavoriteListByUserId", zap.Error(err))
		}
		return favoriteCount, redis.KeyUpdated
	}
	// 重试
	time.Sleep(redis.RetryTime)
	return checkAndSetUserFavoriteCountKey(userId)
}

func getVideoByVideoId(videoId uint) (model.Video, bool) {
	// 从localCache中取
	if val, err := cache.Get(strconv.Itoa(int(videoId))); err == nil {
		var v model.Video
		msgpack.Unmarshal(val, &v)
		return v, true
	}

	videoModel, found := mysql.FindVideoByVideoId(videoId)
	if found {
		go func(uint, model.Video) {
			data, _ := msgpack.Marshal(&videoModel)
			cache.Set(strconv.Itoa(int(videoId)), data)
		}(videoId, videoModel)
	}
	return videoModel, found
}

func startTimer(msgChan chan<- favoriteMap) {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		flushMutex.Lock()
		if favoriteData != nil && len(favoriteData) != 0 {
			msgChan <- favoriteData // 发送 map 到通道
			favoriteData = make(favoriteMap)
		}
		flushMutex.Unlock()
	}
}

func consumerFavoriteMap(ch <-chan favoriteMap) {
	for {
		data := <-ch // 从通道接收数据
		addFavoriteList := make([]model.Favorite, 0, len(data))
		deleteFavoriteIDList := make([]model.Favorite, 0, len(data))
		for userId, innerMap := range data {
			for videoId, actionType := range innerMap {
				isValid, authorId := isFavoriteActionValid(userId, videoId)
				if !isValid {
					continue
				}
				switch actionType {
				case 1:
					// 重复点赞
					if mysql.IsFavorite(userId, videoId) {
						continue
					}
					favoriteAction := model.Favorite{
						UserId:  userId,
						VideoId: videoId,
					}
					addFavoriteList = append(addFavoriteList, favoriteAction)
				case 2:
					// 重复取消点赞
					id, exist := mysql.FindFavoriteByVideoId(userId, videoId)
					if !exist {
						continue
					}
					cancelAction := model.Favorite{
						ID: id,
					}
					deleteFavoriteIDList = append(deleteFavoriteIDList, cancelAction)
				}
				addFavoriteActionToRedis(userId, videoId, actionType, authorId)
			}
		}
		if len(addFavoriteList) != 0 {
			mysql.BatchCreateUserFavorite(addFavoriteList)
		}
		if len(deleteFavoriteIDList) != 0 {
			mysql.BatchDeleteUserFavorite(deleteFavoriteIDList)
		}
	}
}

// isFavoriteActionValid
// 用于判断点赞行为是否合法，若合法返回视频作者ID
func isFavoriteActionValid(userId uint, videoId uint) (valid bool, authorId uint) {
	videoModel, found := getVideoByVideoId(videoId)
	_, exist, err := mysql.FindUserByUserID(userId)
	if !exist || err != nil {
		return false, 0
	}
	if !found {
		return false, 0
	}
	return true, videoModel.AuthorId
}
