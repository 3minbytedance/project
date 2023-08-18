package main

import (
	"context"
	"douyin/constant"
	"douyin/dal/model"
	dalMySQL "douyin/dal/mysql"
	comment "douyin/kitex_gen/comment"
	"douyin/kitex_gen/comment/commentservice"
	favorite "douyin/kitex_gen/favorite"
	user "douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"douyin/kitex_gen/video"
	mwRedis "douyin/mw/redis"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}))
	if err != nil {
		log.Fatal(err)
	}
	commentClient, err = commentservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}))
	if err != nil {
		log.Fatal(err)
	}
}

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, request *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	userId := uint(request.UserId)
	videoId := uint(request.VideoId)
	actionType := int(request.ActionType)

	err = favoriteActions(userId, uint(videoId), actionType)
	//count, _ := service.GetFavoritesVideoCount(int64(videoId))
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
	userId := request.UserId

	favoritesByUserId, err := getFavoritesByUserId(uint(userId))
	if err != nil {
		return nil, err
	}
	videos := make([]video.Video, 0)
	for _, id := range favoritesByUserId {
		videoByVideoId, _ := dalMySQL.FindVideoByVideoId(id)
		userResp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{ActorId: request.GetUserId(),
			UserId: int32(videoByVideoId.AuthorId)})
		if err != nil {
			log.Println(err.Error())
		}
		commentCountResp, err := commentClient.GetCommentCount(ctx, &comment.CommentCountRequest{
			VideoId: int32(id),
		})
		if err != nil {
			log.Println(err.Error())
		}
		favoriteCount, err := getFavoritesVideoCount(id)
		if err != nil {
			log.Println(err.Error())
		}
		vid := video.Video{
			Id:            int32(videoByVideoId.ID),
			Author:        userResp.User,
			PlayUrl:       videoByVideoId.VideoUrl,
			CoverUrl:      videoByVideoId.CoverUrl,
			FavoriteCount: int32(favoriteCount),
			CommentCount:  commentCountResp.CommentCount,
			IsFavorite:    isUserFavorite(uint(userId), id),
			Title:         videoByVideoId.Title,
		}
		videos = append(videos, vid)
	}
	return &favorite.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("get favorite video list done"),
	}, err
}

// GetVideoFavoriteCount implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetVideoFavoriteCount(ctx context.Context, request *favorite.VideoFavoriteCountRequest) (resp *favorite.VideoFavoriteCountResponse, err error) {
	videoId := request.VideoId
	count, err := getFavoritesVideoCount(uint(videoId))
	if err != nil {
		return &favorite.VideoFavoriteCountResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr(err.Error()),
			Count:      0,
		}, err
	}
	return &favorite.VideoFavoriteCountResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
		Count:      int32(count),
	}, nil
}

// GetUserFavoriteCount implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetUserFavoriteCount(ctx context.Context, request *favorite.UserFavoriteCountRequest) (resp *favorite.UserFavoriteCountResponse, err error) {
	userId := uint(request.UserId)
	favoritesByUserId, err := getFavoritesByUserId(userId)
	if err != nil {
		return &favorite.UserFavoriteCountResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr(err.Error()),
			Count:      0,
		}, err
	}
	return &favorite.UserFavoriteCountResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("success"),
		Count:      int32(len(favoritesByUserId)),
	}, nil
}

// favoriteActions 点赞，取消赞的操作过程
func favoriteActions(userId uint, videoId uint, actionType int) error {
	// 判断是否在redis中，如果没有的话，一起加载到redis中
	_, err := getFavoritesByUserId(userId)
	if err != nil {
		return err
	}
	_, err = getFavoritesVideoCount(videoId)
	if err != nil {
		return err
	}
	getUserTotalFavoritedCount(userId)
	video, _ := dalMySQL.FindVideoByVideoId(videoId)
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
			fmt.Println(err)
		}
		// 更新用户喜欢的视频数量，这个不用，直接从set中获取
		// 更新视频被喜欢的数量
		err = mwRedis.IncrementFavoritedCountByVideoId(videoId)
		if err != nil {
			fmt.Println(err)
		}
		// 更新视频作者的被点赞量
		getUserTotalFavoritedCount(video.AuthorId)
		err = mwRedis.IncrementTotalFavoritedByUserId(video.AuthorId)
		if err != nil {
			fmt.Println(err)
		}
		dalMySQL.AddUserFavorite(userId, videoId)
		return err
	case 2:
		// 取消赞
		if !mwRedis.IsInUserFavoriteList(userId, videoId) {
			return errors.New("该视频未点赞")
		}
		// 更新视频被喜欢的用户列表
		err = mwRedis.DeleteFavoriteVideoFromList(userId, videoId)
		if err != nil {
			fmt.Println(err)
		}
		// 更新视频被喜欢的数量
		err = mwRedis.DecrementFavoritedCountByVideoId(videoId)
		if err != nil {
			fmt.Println(err)
		}
		getUserTotalFavoritedCount(video.AuthorId)
		// 更新视频作者的被点赞量
		err = mwRedis.DecrementTotalFavoritedByUserId(video.AuthorId)
		if err != nil {
			fmt.Println(err)
		}
		err = dalMySQL.DeleteUserFavorite(userId, videoId)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// GetFavoritesVideoCount 根据视频id，返回该视频的点赞数
func getFavoritesVideoCount(videoId uint) (int64, error) {
	// 判断redis中是否存在对应的video数据
	exits := mwRedis.IsExistVideoField(videoId, mwRedis.VideoFavoritedCountField)
	if exits {
		// redis中存在对应的数据
		count, err := mwRedis.GetFavoritedCountByVideoId(videoId)
		if err != nil {
			fmt.Println(err)
		}
		return count, err
	} else {
		// redis中不存在，从数据库中读取
		_, num, err := dalMySQL.GetFavoritesByIdFromMysql(videoId, dalMySQL.IdTypeVideo)
		if err != nil {
			log.Println(err)
		}
		err = mwRedis.SetVideoFavoritedCountByVideoId(videoId, int64(num)) // 加载视频被点赞数量
		if err != nil {
			fmt.Println(err)
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
			log.Println(err)
		}
		return favoritesVideoIdList, err
	}

	// redis中没有对应的数据，从MYSQL数据库中获取数据
	favorites, _, err := dalMySQL.GetFavoritesByIdFromMysql(userId, dalMySQL.IdTypeUser)
	if err != nil {
		log.Println(err)
	}
	if len(favorites) == 0 {
		//点赞数为0，设置为0
		if err = mwRedis.SetFavoriteListByUserId(userId, []uint{0}); err != nil {
			log.Println(err)
		}
		return []uint{}, err
	}
	//removeNilValue(userId)
	idList := getIdListFromFavoriteSlice(favorites, dalMySQL.IdTypeUser)
	// key 不存在需要同步到redis
	err = mwRedis.SetFavoriteListByUserId(userId, idList) // 加载到set中
	if err != nil {
		log.Println(err)
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

// GetUserTotalFavoritedCount 获取用户发布视频的总的被点赞数量
func getUserTotalFavoritedCount(userId uint) int64 {
	exits := mwRedis.IsExistUserField(userId, mwRedis.TotalFavoriteField)
	if exits {
		// redis中存在对应的数据
		count, err := mwRedis.GetTotalFavoritedByUserId(userId)
		if err != nil {
			fmt.Println(err)
		}
		return count
	}
	//redis 不存在
	var total int64
	// 获取用户发布的视频列表
	videosByAuthorId, exist := dalMySQL.FindVideosByAuthorId(userId)
	if !exist {
		return 0
	}
	for _, video := range videosByAuthorId {
		count, _ := getFavoritesVideoCount(video.ID)
		total += count
	}
	err := mwRedis.SetTotalFavoritedByUserId(userId, total)
	if err != nil {
		log.Println(err)
	}
	return total
}

// IsUserFavorite 判断是否点赞
func isUserFavorite(userId, videoId uint) bool {
	exist := mwRedis.IsExistUserSetField(userId, mwRedis.FavoriteList)
	if !exist {
		// redis中没有对应的数据，从MYSQL数据库中获取数据
		favorites, _, err := dalMySQL.GetFavoritesByIdFromMysql(userId, dalMySQL.IdTypeUser)
		if err != nil {
			log.Println(err)
		}
		if len(favorites) == 0 {
			//点赞数为0，设置0
			if err = mwRedis.SetFavoriteListByUserId(userId, []uint{0}); err != nil {
				log.Println(err)
			}
		} else {
			//removeNilValue(userId)
			idList := getIdListFromFavoriteSlice(favorites, dalMySQL.IdTypeUser)
			// key 不存在需要同步到redis
			err = mwRedis.SetFavoriteListByUserId(userId, idList) // 加载到set中
			if err != nil {
				log.Println(err)
			}
		}
	}
	return mwRedis.IsInUserFavoriteList(userId, videoId)
}

//func removeNilValue(userId uint) {
//	if mwRedis.IsUserFavoriteNil(userId) {
//		err := mwRedis.DeleteUserFavoriteNil(userId)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}
