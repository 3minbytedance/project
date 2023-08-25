package main

import (
	"context"
	"douyin/constant"
	"douyin/constant/biz"
	"douyin/dal/model"
	dalMySQL "douyin/dal/mysql"
	"douyin/kitex_gen/comment/commentservice"
	favorite "douyin/kitex_gen/favorite"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/kitex_gen/video"
	"douyin/mw/kafka"
	mwRedis "douyin/mw/redis"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}))
	if err != nil {
		log.Fatal(err)
	}
	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}))
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

	err = favoriteActions(userId, videoId, actionType)
	if err != nil {
		return &favorite.FavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("action failed，err: " + err.Error()),
		}, err
	}
	return &favorite.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("action down"),
	}, nil
}

// GetFavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteList(ctx context.Context, request *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	resp = new(favorite.FavoriteListResponse)
	userId := request.GetUserId()
	actionId := request.GetActionId()

	favoritesByUserId, err := getFavoritesByUserId(uint(userId))
	if err != nil {
		return nil, err
	}
	videos := make([]*video.Video, 0, len(favoritesByUserId))
	for _, id := range favoritesByUserId {
		videoModel, b := dalMySQL.FindVideoByVideoId(id)
		if b == false {
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
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("get favorite video list done"),
		VideoList:  videos,
	}, err
}

// GetVideoFavoriteCount implements the FavoriteServiceImpl interface.
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
	favoritesByUserId, err := getFavoritesByUserId(uint(userId))
	if err != nil {
		return 0, err
	}
	return int32(len(favoritesByUserId)), nil
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
func favoriteActions(userId uint, videoId uint, actionType int) error {
	_, b := dalMySQL.FindVideoByVideoId(videoId)
	if b == false {
		return errors.New("video id not exist")
	}
	// 判断是否在redis中，如果没有的话，一起加载到redis中
	_, err := getFavoritesByUserId(userId)
	if err != nil {
		return err
	}
	_, err = getFavoritesVideoCount(videoId)
	if err != nil {
		return err
	}
	_, err = getUserTotalFavoritedCount(userId)
	videoModel, _ := dalMySQL.FindVideoByVideoId(videoId)
	// 如果err不为空，那么一定存在数据库中了
	switch actionType {
	case 1:
		// 点赞
		// 判断重复点赞
		if mwRedis.IsInUserFavoriteList(userId, videoId) {
			return errors.New("该视频已点赞")
		}
		// 更新用户喜欢的视频列表
		err = mwRedis.AddFavoriteVideoToList(userId, videoId)
		if err != nil {
			zap.L().Error("更新用户喜欢的视频列表", zap.Error(err))
			return err
		}
		// 更新用户喜欢的视频数量，这个不用，直接从set中获取
		// 更新视频被喜欢的数量
		err = mwRedis.IncrementFavoritedCountByVideoId(videoId)
		if err != nil {
			zap.L().Error("更新视频被喜欢的数量", zap.Error(err))
			return err
		}
		// 更新视频作者的被点赞量
		//getUserTotalFavoritedCount(videoModel.AuthorId)
		err = mwRedis.IncrementTotalFavoritedByUserId(videoModel.AuthorId)
		if err != nil {
			zap.L().Error("更新视频作者的被点赞量", zap.Error(err))
			return err
		}
		// dalMySQL.AddUserFavorite(userId, videoId)
		err := kafka.FavoriteMQInstance.ProduceAddFavoriteMsg(userId, videoId)
		if err != nil {
			zap.L().Error("更新视频作者的被点赞量", zap.Error(err))
		}
		return err
	case 2:
		// 取消赞
		if !mwRedis.IsInUserFavoriteList(userId, videoId) {
			return errors.New("该视频未点赞")
		}
		// 更新视频被喜欢的用户列表
		err = mwRedis.DeleteFavoriteVideoFromList(userId, videoId)
		if err != nil {
			zap.L().Error("更新视频被喜欢的用户列表", zap.Error(err))
			return err
		}
		// 更新视频被喜欢的数量
		err = mwRedis.DecrementFavoritedCountByVideoId(videoId)
		if err != nil {
			zap.L().Error("更新视频被喜欢的数量", zap.Error(err))
			return err
		}
		_, err = getUserTotalFavoritedCount(videoModel.AuthorId)
		// 更新视频作者的被点赞量
		err = mwRedis.DecrementTotalFavoritedByUserId(videoModel.AuthorId)
		if err != nil {
			zap.L().Error("更新视频作者的被点赞量", zap.Error(err))
			return err
		}
		// err = dalMySQL.DeleteUserFavorite(userId, videoId)
		err = kafka.FavoriteMQInstance.ProduceDelFavoriteMsg(userId, videoId)
		if err != nil {
			zap.L().Error("更新视频作者的被点赞量", zap.Error(err))
			return err
		}
	}
	return nil
}

// GetUserTotalFavoritedCount 获取用户发布视频的总的被点赞数量
func getUserTotalFavoritedCount(userId uint) (int64, error) {
	exits := mwRedis.IsExistUserField(userId, mwRedis.TotalFavoriteField)
	if exits {
		// redis中存在对应的数据
		count, err := mwRedis.GetTotalFavoritedByUserId(userId)
		if err != nil {
			fmt.Println(err)
		}
		return count, nil
	}
	//redis 不存在
	var total int64
	// 获取用户发布的视频列表
	videosByAuthorId, exist := dalMySQL.FindVideosByAuthorId(userId)
	if !exist {
		return 0, nil
	}
	for _, videoModel := range videosByAuthorId {
		count, _ := getFavoritesVideoCount(videoModel.ID)
		total += count
	}
	err := mwRedis.SetTotalFavoritedByUserId(userId, total)
	if err != nil {
		zap.L().Error("SetTotalFavoriteByUserId失败", zap.Error(err))
		return 0, err
	}
	return total, err
}

// GetFavoritesVideoCount 根据视频id，返回该视频的点赞数
func getFavoritesVideoCount(videoId uint) (int64, error) {
	// 判断redis中是否存在对应的video数据
	exits := mwRedis.IsExistVideoField(videoId, mwRedis.VideoFavoritedCountField)
	if exits {
		// redis中存在对应的数据
		count, err := mwRedis.GetFavoritedCountByVideoId(videoId)
		if err != nil {
			zap.L().Error("GetFavoritedCountByVideoId失败", zap.Error(err))
		}
		return count, err
	} else {
		// redis中不存在，从数据库中读取
		_, num, err := dalMySQL.GetFavoritesByIdFromMysql(videoId, dalMySQL.IdTypeVideo)
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
		}
		if num == 0 {
			return 0, nil
		}
		err = mwRedis.SetVideoFavoritedCountByVideoId(videoId, int64(num)) // 加载视频被点赞数量
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
		}
		return int64(num), err
	}
}

// getFavoritesByUserId 获取当前user_id的点赞的视频id列表
func getFavoritesByUserId(userId uint) ([]uint, error) {
	// 查看redis是否存在对应的user数据
	exist := mwRedis.IsExistUserSetField(userId, mwRedis.FavoriteList)
	if exist {
		// redis存在
		favoritesVideoIdList, err := mwRedis.GetFavoriteListByUserId(userId)
		if err != nil {
			zap.L().Error("GetFavoriteListByUserId", zap.Error(err))
		}
		return favoritesVideoIdList, err
	}

	// redis中没有对应的数据，从MYSQL数据库中获取数据
	favorites, favoriteLength, err := dalMySQL.GetFavoritesByIdFromMysql(userId, dalMySQL.IdTypeUser)
	if err != nil {
		zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
	}
	if favoriteLength == 0 {
		return []uint{}, nil
	}
	//removeNilValue(userId)
	idList := getIdListFromFavoriteSlice(favorites, dalMySQL.IdTypeUser)
	// key 不存在需要同步到redis
	err = mwRedis.SetFavoriteListByUserId(userId, idList) // 加载到set中
	if err != nil {
		zap.L().Error("SetFavoriteListByUserId", zap.Error(err))
	}
	return idList, nil
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
	exist := mwRedis.IsExistUserSetField(userId, mwRedis.FavoriteList)
	if !exist {
		// redis中没有对应的数据，从MYSQL数据库中获取数据
		favorites, favoriteLength, err := dalMySQL.GetFavoritesByIdFromMysql(userId, dalMySQL.IdTypeUser)
		if err != nil {
			zap.L().Error("GetFavoritesByIdFromMysql", zap.Error(err))
			return false
		}
		//点赞数为0
		if favoriteLength == 0 {
			return false
		}
		idList := getIdListFromFavoriteSlice(favorites, dalMySQL.IdTypeUser)
		// key 不存在需要同步到redis
		err = mwRedis.SetFavoriteListByUserId(userId, idList) // 加载到set中
		if err != nil {
			zap.L().Error("SetFavoriteListByUserId", zap.Error(err))
		}
	}
	return mwRedis.IsInUserFavoriteList(userId, videoId)
}
