package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/constant/biz"
	"douyin/dal/mysql"
	"douyin/kitex_gen/comment/commentservice"
	"douyin/kitex_gen/favorite/favoriteservice"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	video "douyin/kitex_gen/video"
	"douyin/mw/redis"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/google/uuid"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var userClient userservice.Client
var commentClient commentservice.Client
var favoriteClient favoriteservice.Client
var relationClient relationservice.Client

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
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
	)
	relationClient, err = relationservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.RelationServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// VideoFeed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoFeed(ctx context.Context, request *video.VideoFeedRequest) (resp *video.VideoFeedResponse, err error) {
	videos := mysql.GetLatestVideos(request.GetLatestTime())
	if len(videos) == 0 {
		zap.L().Error("根据LatestTime取视频失败", zap.Error(err))
		return &video.VideoFeedResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取视频列表失败"),
			VideoList:  nil,
		}, err
	}
	currentId := request.GetUserId()
	isLogged := false
	if currentId != 0 {
		isLogged = true
	}

	// 将查询结果转换为VideoResponse类型
	videoList := make([]*video.Video, 0, len(videos))
	for _, v := range videos {
		userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			UserId: int32(v.AuthorId),
		})
		if isLogged {
			following, _ := relationClient.IsFollowing(ctx, &relation.IsFollowingRequest{
				ActorId: currentId,
				UserId:  int32(v.AuthorId),
			})
			userResp.GetUser().SetIsFollow(following)
		}
		commentCount, _ := commentClient.GetCommentCount(ctx, int32(v.ID))
		favoriteCount, _ := favoriteClient.GetVideoFavoriteCount(ctx, int32(v.ID))

		videoResponse := video.Video{
			Id:            int32(v.ID),
			Author:        userResp.GetUser(),
			PlayUrl:       biz.OSS + v.VideoUrl,
			CoverUrl:      biz.OSS + v.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			//IsFavorite:    IsUserFavorite(userID, video.ID), // todo
			Title: v.Title,
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
	// 生成 UUID

	fileId := strings.Replace(uuid.New().String(), "-", "", -1)
	// 生成新的文件名
	videoFileName := fileId + ".mp4"

	//todo 待改
	//err := SaveUploadedFile(file, "./public/"+videoFileName)
	//if err != nil {
	//	return &video.PublishVideoResponse{
	//		StatusCode: 1,
	//		StatusMsg:  nil,
	//	}, nil
	//}

	// MQ 异步解耦,解决返回json阻塞 TODO

	imgName := GetVideoCover(videoFileName)

	//视频存储到oss
	go func() {
		if err := UploadToOSS(biz.FileLocalPath+request.GetTitle(), request.GetTitle()); err != nil {
			log.Fatal(err)
			return
		}
	}()

	// 图片存储到oss
	go func() {
		if err := UploadToOSS(biz.FileLocalPath+imgName, imgName); err != nil {
			log.Fatal(err)
			return
		}
	}()

	go func() {
		mysql.InsertVideo(videoFileName, imgName, uint(request.GetUserId()), request.GetTitle())
		if !redis.IsExistUserField(uint(request.GetUserId()), redis.WorkCountField) {
			cnt := mysql.FindWorkCountsByAuthorId(uint(request.GetUserId()))
			err := redis.SetWorkCountByUserId(uint(request.GetUserId()), cnt)
			if err != nil {
				log.Println("redis更新评论数失败", err)
				return
			}
			return
		}
		err := redis.IncrementWorkCountByUserId(uint(request.GetUserId()))
		if err != nil {
			log.Println(err)
		}
	}()
	return
}

// GetPublishVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishVideoList(ctx context.Context, request *video.PublishVideoListRequest) (resp *video.PublishVideoListResponse, err error) {
	videos, found := mysql.FindVideosByAuthorId(uint(request.GetFromUserId()))
	if !found {
		return &video.PublishVideoListResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("获取视频列表失败"),
		}, nil
	}
	// 将查询结果转换为VideoResponse类型
	videoList := make([]*video.Video, 0, len(videos))
	for _, v := range videos {
		userResp, _ := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
			UserId:  request.GetToUserId(),
			ActorId: request.GetFromUserId(),
		})
		commentCount, _ := commentClient.GetCommentCount(ctx, int32(v.ID))
		favoriteCount, _ := favoriteClient.GetVideoFavoriteCount(ctx, int32(v.ID))

		videoResponse := video.Video{
			Id:            int32(v.ID),
			Author:        userResp.GetUser(),
			PlayUrl:       biz.OSS + v.VideoUrl,
			CoverUrl:      biz.OSS + v.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			//IsFavorite:    IsUserFavorite(userID, video.ID), // todo
			Title: v.Title,
		}
		videoList = append(videoList, &videoResponse)
	}

	return &video.PublishVideoListResponse{
		StatusCode: 0,
		VideoList:  videoList,
	}, nil
}

// UploadToOSS  上传至腾讯OSS
func UploadToOSS(localPath string, remotePath string) error {
	reqUrl := biz.OSS
	u, _ := url.Parse(reqUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     biz.SecretId,
			SecretKey:    biz.SecretKey,
			SessionToken: biz.SessionToken,
		},
	})

	_, _, err := c.Object.Upload(context.Background(), remotePath, localPath, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoCover(fileName string) string {
	// 生成图片 UUID
	imgId := uuid.New().String()
	// 修改文件名
	imgName := strings.Replace(imgId, "-", "", -1) + ".jpg"
	//调用ffmpeg 获取封面图
	common.GetVideoFrame(biz.FileLocalPath+fileName, biz.FileLocalPath+imgName)
	return imgName
}

// GetWorkCount 返回某个用户的作品数
func (s *VideoServiceImpl) GetWorkCount(ctx context.Context, request *video.GetWorkCountRequest) (resp int32, err error) {
	userId := uint(request.GetUserId())
	// 从redis中获取作品数
	// 1. 缓存中有数据, 直接返回
	if redis.IsExistUserField(userId, redis.WorkCountField) {
		workCount, err := redis.GetWorkCountByUserId(userId)
		if err != nil {
			log.Println("从redis中获取作品数失败：", err)
		}
		return int32(workCount), nil
	}

	// 2. 缓存中没有数据，从数据库中获取
	workCount := mysql.FindWorkCountsByAuthorId(userId)
	log.Println("从数据库中获取作品数成功：", workCount)
	// 将作品数写入redis
	go func() {
		err := redis.SetWorkCountByUserId(userId, workCount)
		if err != nil {
			log.Println("将作品数写入redis失败：", err)
		}
	}()
	return int32(workCount), nil
}
