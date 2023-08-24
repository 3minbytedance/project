package main

import (
	"context"
	"douyin/constant"
	"douyin/constant/biz"
	"douyin/dal/mysql"
	"douyin/kitex_gen/comment/commentservice"
	"douyin/kitex_gen/favorite"
	"douyin/kitex_gen/favorite/favoriteservice"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	video "douyin/kitex_gen/video"
	"douyin/mw/kafka"
	"douyin/mw/redis"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/google/uuid"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var userClient userservice.Client
var commentClient commentservice.Client
var favoriteClient favoriteservice.Client

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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}))
	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
	)
	favoriteClient, err = favoriteservice.NewClient(
		constant.FavoriteServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// VideoFeed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoFeed(ctx context.Context, request *video.VideoFeedRequest) (resp *video.VideoFeedResponse, err error) {
	zap.L().Info("VideoFeed", zap.Any("request", request))
	latestTime := request.GetLatestTime()
	videos := mysql.GetLatestVideos(latestTime)
	if len(videos) == 0 {
		zap.L().Info("根据LatestTime取视频失败")
		return &video.VideoFeedResponse{
			StatusCode: 0,
			StatusMsg:  thrift.StringPtr("获取视频列表失败"),
			VideoList:  nil,
		}, err
	}
	currentId := request.GetUserId()

	// 将查询结果转换为VideoResponse类型
	videoList := make([]*video.Video, 0, len(videos))
	for _, v := range videos {
		userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			ActorId: currentId,
			UserId:  int64(v.AuthorId),
		})
		commentCount, _ := commentClient.GetCommentCount(ctx, int64(v.ID))
		favoriteCount, _ := favoriteClient.GetVideoFavoriteCount(ctx, int64(v.ID))
		isFavorite, _ := favoriteClient.IsUserFavorite(ctx, &favorite.IsUserFavoriteRequest{
			UserId:  currentId,
			VideoId: int64(v.ID),
		})
		videoResponse := video.Video{
			Id:            int64(v.ID),
			Author:        userResp.GetUser(),
			PlayUrl:       biz.OSS + v.VideoUrl,
			CoverUrl:      biz.OSS + v.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
		}
		videoList = append(videoList, &videoResponse)
	}
	nextTime := videos[len(videos)-1].CreatedAt

	return &video.VideoFeedResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		VideoList:  videoList,
		NextTime:   &nextTime,
	}, nil
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, request *video.PublishVideoRequest) (resp *video.PublishVideoResponse, err error) {
	// 根据UUID生成新的文件名
	videoFileName := strings.Replace(uuid.New().String(), "-", "", -1) + ".mp4"
	videoPath := biz.FileLocalPath + videoFileName
	err = os.WriteFile(videoPath, request.GetData(), 0644)
	if err != nil {
		return &video.PublishVideoResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("文件上传失败"),
		}, err
	}

	go func() {
		// 通过MQ异步处理视频的上传操作, 包括上传到OSS, 保存到MySQL, 更新redis
		zap.L().Info("上传视频发送到消息队列", zap.String("videoPath", videoPath))
		kafka.VideoMQInstance.Produce(&kafka.VideoMessage{
			VideoPath:     videoPath,
			VideoFileName: videoFileName,
			UserID:        uint(request.GetUserId()),
			Title:         request.GetTitle(),
		})
	}()

	return &video.PublishVideoResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("Success"),
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
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取视频列表失败"),
		}, nil
	}
	// 将查询结果转换为VideoResponse类型
	videoList := make([]*video.Video, 0, len(videos))
	//toUserId的用户信息
	userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
		ActorId: fromUserId,
		UserId:  toUserId,
	})
	//toUserId 发布的视频
	for _, v := range videos {
		commentCount, _ := commentClient.GetCommentCount(ctx, int64(v.ID))
		favoriteCount, _ := favoriteClient.GetVideoFavoriteCount(ctx, int64(v.ID))
		//判断当前请求用户是否点赞该视频
		isFavorite, _ := favoriteClient.IsUserFavorite(ctx, &favorite.IsUserFavoriteRequest{
			UserId:  fromUserId,
			VideoId: int64(v.ID),
		})
		videoResponse := video.Video{
			Id:            int64(v.ID),
			Author:        userResp.GetUser(),
			PlayUrl:       biz.OSS + v.VideoUrl,
			CoverUrl:      biz.OSS + v.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
		}
		videoList = append(videoList, &videoResponse)
	}

	return &video.PublishVideoListResponse{
		StatusCode: 0,
		VideoList:  videoList,
	}, nil
}

// GetWorkCount implements the VideoServiceImpl interface.
// 获取userId的作品数
func (s *VideoServiceImpl) GetWorkCount(ctx context.Context, userId int64) (resp int32, err error) {
	// 从redis中获取作品数
	// 1. 缓存中有数据, 直接返回
	if redis.IsExistUserField(uint(userId), redis.WorkCountField) {
		workCount, err := redis.GetWorkCountByUserId(uint(userId))
		if err != nil {
			zap.L().Error("从redis中获取作品数失败", zap.Error(err))
		}
		return int32(workCount), nil
	}

	// 2. 缓存中没有数据，从数据库中获取
	workCount := mysql.FindWorkCountsByAuthorId(uint(userId))
	// 将作品数写入redis
	go func() {
		err := redis.SetWorkCountByUserId(uint(userId), workCount)
		if err != nil {
			zap.L().Error("将作品数写入redis失败", zap.Error(err))
		}
	}()
	return int32(workCount), nil
}
