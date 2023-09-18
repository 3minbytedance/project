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
	"douyin/kitex_gen/favorite/favoriteservice"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	video "douyin/kitex_gen/video"
	"douyin/mw/kafka"
	"douyin/mw/localcache"
	"douyin/mw/redis"
	"github.com/allegro/bigcache/v3"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/google/uuid"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var userClient userservice.Client
var commentClient commentservice.Client
var favoriteClient favoriteservice.Client
var cache *bigcache.BigCache

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

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
	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
		client.WithMuxConnection(2),
	)
	favoriteClient, err = favoriteservice.NewClient(
		constant.FavoriteServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
		client.WithMuxConnection(2),
	)
	if err != nil {
		log.Fatal(err)
	}

	cache = localcache.Init(localcache.WorkCount)
}

// VideoFeed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoFeed(ctx context.Context, request *video.VideoFeedRequest) (resp *video.VideoFeedResponse, err error) {
	zap.L().Info("VideoFeed", zap.Any("request", request))
	latestTime := request.GetLatestTime()
	//videos := mysql.GetLatestVideos(latestTime)
	videos := redis.GetVideos(latestTime)

	if videos == nil || len(videos) == 0 {
		//如果视频都看完，重置时间戳
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
		videos = redis.GetVideos(latestTime)
	}
	currentId := request.GetUserId()
	videoList := make([]*video.Video, 0, len(videos))
	userRespCh := make(chan *user.UserInfoByIdResponse)
	commentCountCh := make(chan int32)
	favoriteCountCh := make(chan int32)
	isFavoriteCh := make(chan bool)
	defer func() {
		close(userRespCh)
		close(commentCountCh)
		close(favoriteCountCh)
		close(isFavoriteCh)
	}()
	var videoModel model.Video
	for _, v := range videos {
		err = msgpack.Unmarshal([]byte(v), &videoModel)
		if err != nil {
			continue
		}
		go func() {
			userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
				ActorId: currentId,
				UserId:  int64(videoModel.AuthorId),
			})
			userRespCh <- userResp
		}()
		go func() {
			commentCount, _ := commentClient.GetCommentCount(ctx, int64(videoModel.ID))
			commentCountCh <- commentCount
		}()
		go func() {
			favoriteCount, _ := favoriteClient.GetVideoFavoriteCount(ctx, int64(videoModel.ID))
			favoriteCountCh <- favoriteCount
		}()
		//判断当前请求用户是否点赞该视频
		go func(id int64) {
			if id != 0 {
				isFavorite, _ := favoriteClient.IsUserFavorite(ctx, &favorite.IsUserFavoriteRequest{
					UserId:  id,
					VideoId: int64(videoModel.ID),
				})
				isFavoriteCh <- isFavorite
				return
			}
			isFavoriteCh <- false
		}(currentId)

		// 将查询结果转换为VideoResponse类型
		videoResponse := video.Video{
			Id:       int64(videoModel.ID),
			PlayUrl:  biz.OSS + videoModel.VideoUrl,
			CoverUrl: biz.OSS + videoModel.CoverUrl,
			Title:    videoModel.Title,
		}
		for receivedCount := 0; receivedCount < 4; receivedCount++ {
			select {
			case userResp := <-userRespCh:
				videoResponse.SetAuthor(userResp.GetUser())
			case favoriteCount := <-favoriteCountCh:
				videoResponse.SetFavoriteCount(favoriteCount)
			case isFavorite := <-isFavoriteCh:
				videoResponse.SetIsFavorite(isFavorite)
			case commentCount := <-commentCountCh:
				videoResponse.SetCommentCount(commentCount)
			case <-time.After(3 * time.Second):
				zap.L().Error("3s overtime.")
				break
			}
		}

		videoList = append(videoList, &videoResponse)
	}
	nextTime := videoModel.CreatedAt

	return &video.VideoFeedResponse{
		StatusCode: common.CodeSuccess,
		StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		VideoList:  videoList,
		NextTime:   nextTime,
	}, nil
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, request *video.PublishVideoRequest) (resp *video.PublishVideoResponse, err error) {
	// 根据UUID生成新的文件名
	videoFileName := strings.Replace(uuid.New().String(), "-", "", -1) + ".mp4"
	videoPath := biz.FileLocalPath + videoFileName
	err = os.WriteFile(videoPath, request.GetData(), biz.FileMode)
	if err != nil {
		err = nil
		return &video.PublishVideoResponse{
			StatusCode: common.CodeUploadFileError,
			StatusMsg:  common.MapErrMsg(common.CodeUploadFileError),
		}, nil
	}

	// 通过MQ异步处理视频的上传操作, 包括上传到OSS，截帧, 保存到MySQL, 更新redis
	zap.L().Info("上传视频发送到消息队列", zap.String("videoPath", videoPath))
	kafka.VideoMQInstance.Produce(&kafka.VideoMessage{
		VideoPath:     videoPath,
		VideoFileName: videoFileName,
		UserID:        uint(request.GetUserId()),
		Title:         request.GetTitle(),
	})

	return &video.PublishVideoResponse{
		StatusCode: common.CodeSuccess,
		StatusMsg:  common.MapErrMsg(common.CodeSuccess),
	}, nil
}

// GetPublishVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishVideoList(ctx context.Context, request *video.PublishVideoListRequest) (resp *video.PublishVideoListResponse, err error) {
	zap.L().Info("GetPublishVideoList", zap.Any("request", request))
	videos, found := mysql.FindVideosByAuthorId(uint(request.GetToUserId()))
	fromUserId := request.GetFromUserId()
	toUserId := request.GetToUserId()
	if !found {
		return &video.PublishVideoListResponse{
			StatusCode: common.CodeSuccess,
			StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		}, nil
	}
	// 将查询结果转换为VideoResponse类型
	videoList := make([]*video.Video, 0, len(videos))
	//toUserId的用户信息
	userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
		ActorId: fromUserId,
		UserId:  toUserId,
	})
	if userResp == nil {
		return &video.PublishVideoListResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		}, nil
	}
	//toUserId 发布的视频
	commentCountCh := make(chan int32)
	favoriteCountCh := make(chan int32)
	isFavoriteCh := make(chan bool)
	defer func() {
		close(commentCountCh)
		close(favoriteCountCh)
		close(isFavoriteCh)
	}()
	for _, v := range videos {
		go func() {
			commentCount, _ := commentClient.GetCommentCount(ctx, int64(v.ID))
			commentCountCh <- commentCount
		}()
		go func() {
			favoriteCount, _ := favoriteClient.GetVideoFavoriteCount(ctx, int64(v.ID))
			favoriteCountCh <- favoriteCount
		}()
		//判断当前请求用户是否点赞该视频
		go func(id int64) {
			if id != 0 {
				isFavorite, _ := favoriteClient.IsUserFavorite(ctx, &favorite.IsUserFavoriteRequest{
					UserId:  id,
					VideoId: int64(v.ID),
				})
				isFavoriteCh <- isFavorite
				return
			}
			isFavoriteCh <- false
		}(fromUserId)
		videoResponse := &video.Video{
			Id:       int64(v.ID),
			Author:   userResp.GetUser(),
			PlayUrl:  biz.OSS + v.VideoUrl,
			CoverUrl: biz.OSS + v.CoverUrl,
			Title:    v.Title,
		}
		for receivedCount := 0; receivedCount < 3; receivedCount++ {
			select {
			case favoriteCount := <-favoriteCountCh:
				videoResponse.SetFavoriteCount(favoriteCount)
			case isFavorite := <-isFavoriteCh:
				videoResponse.SetIsFavorite(isFavorite)
			case commentCount := <-commentCountCh:
				videoResponse.SetCommentCount(commentCount)
			case <-time.After(3 * time.Second):
				zap.L().Error("3s overtime.")
				break
			}
		}
		videoList = append(videoList, videoResponse)
	}

	return &video.PublishVideoListResponse{
		StatusCode: common.CodeSuccess,
		StatusMsg:  common.MapErrMsg(common.CodeSuccess),
		VideoList:  videoList,
	}, nil
}

// GetWorkCount implements the VideoServiceImpl interface.
// 获取userId的作品数
func (s *VideoServiceImpl) GetWorkCount(ctx context.Context, userId int64) (resp int32, err error) {
	return getWorkCount(uint(userId))
}

func getWorkCount(userId uint) (int32, error) {
	// 从localCache中取
	if val, err := cache.Get(strconv.Itoa(int(userId))); err == nil {
		count, _ := strconv.Atoi(string(val))
		return int32(count), nil
	}

	// 从redis中获取作品数
	// 1. 缓存中有数据, 直接返回
	if workCount, err := redis.GetWorkCountByUserId(userId); err == nil {
		go func(uint, int64) {
			cache.Set(strconv.Itoa(int(userId)), []byte(strconv.FormatInt(workCount, 10)))
		}(userId, workCount)
		return int32(workCount), nil
	}

	//缓存不存在，尝试从数据库中取
	if redis.AcquireUserLock(userId, redis.WorkCountField) {
		defer redis.ReleaseUserLock(userId, redis.WorkCountField)

		exist := common.TestWorkCountBloom(strconv.Itoa(int(userId)))

		// 不存在
		if !exist {
			err := redis.SetWorkCountByUserId(userId, 0)
			if err != nil {
				zap.L().Error("redis更新作品数失败", zap.Error(err))
			}
			go func(uint) {
				cache.Set(strconv.Itoa(int(userId)), []byte("0"))
			}(userId)
			return 0, nil
		}
		// 2. 从数据库中获取
		workCount := mysql.FindWorkCountsByAuthorId(userId)
		// 将作品数写入redis
		err := redis.SetWorkCountByUserId(userId, workCount)
		if err != nil {
			zap.L().Error("将作品数写入redis失败", zap.Error(err))
		}
		go func(uint, int64) {
			cache.Set(strconv.Itoa(int(userId)), []byte(strconv.FormatInt(workCount, 10)))
		}(userId, workCount)

		return int32(workCount), nil
	}
	//重试
	time.Sleep(redis.RetryTime)
	return getWorkCount(userId)
}

func InitVideoListToRedis() {
	latestTime := strconv.FormatInt(time.Now().Unix(), 10)
	videos := mysql.GetAllVideos(latestTime)
	redis.DelVideoKey()
	redis.AddVideos(videos)
}

func AddWorkCount(userId uint) {
	// 更新redis中的用户作品数
	if !redis.IsExistUserField(userId, redis.WorkCountField) {
		cnt := mysql.FindWorkCountsByAuthorId(userId)
		err := redis.SetWorkCountByUserId(userId, cnt)
		if err != nil {
			zap.L().Error("redis更新作品数失败", zap.Error(err))
			return
		}
		return
	}
	err := redis.IncrementWorkCountByUserId(userId)
	if err != nil {
		zap.L().Error("redis增加其作品数失败", zap.Error(err))
		return
	}
}
