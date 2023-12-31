package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/dal/model"
	"douyin/dal/mysql"
	"douyin/kitex_gen/favorite/favoriteservice"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
	user "douyin/kitex_gen/user"
	"douyin/kitex_gen/video/videoservice"
	"douyin/mw/localcache"
	"douyin/mw/redis"
	"douyin/service/user/pack"
	"github.com/allegro/bigcache/v3"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

var relationClient relationservice.Client
var favoriteClient favoriteservice.Client
var videoClient videoservice.Client
var cache *bigcache.BigCache

func init() {
	// Etcd 服务发现
	r, err := etcd.NewEtcdResolver([]string{constant.EtcdAddr})
	if err != nil {
		log.Fatal(err)
	}
	relationClient, err = relationservice.NewClient(
		constant.RelationServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.RelationServiceName}),
		client.WithMuxConnection(2),
	)
	favoriteClient, err = favoriteservice.NewClient(
		constant.FavoriteServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
		client.WithMuxConnection(2),
	)
	videoClient, err = videoservice.NewClient(
		constant.VideoServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.VideoServiceName}),
		client.WithMuxConnection(2),
	)
	if err != nil {
		log.Fatal(err)
	}

	cache = localcache.Init(localcache.User)
}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, request *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {

	// 检查注册信息是否合理
	resp = new(user.UserRegisterResponse)

	userData := model.User{}
	userData.Name = request.Username
	userData.ID = common.GetUid()

	// 将信息存储到数据库中
	userData.Password, _ = common.MakePassword(request.Password)

	// 数据入库
	err = mysql.CreateUser(&userData)
	if err != nil {
		zap.L().Info("Create user err:", zap.Error(err))
		resp.StatusCode = common.CodeUsernameAlreadyExists
		resp.StatusMsg = common.MapErrMsg(common.CodeUsernameAlreadyExists)
		err = nil
		return
	}

	resp.UserId = int64(userData.ID)
	resp.Token = common.GenerateToken(userData.ID, request.Username)

	// 将token存入redis
	redis.SetToken(userData.ID, resp.Token)
	go func() {
		// 用户名存入Bloom Filter
		common.AddToUserBloom(request.Username)
	}()
	resp.StatusCode = common.CodeSuccess
	resp.StatusMsg = common.MapErrMsg(common.CodeSuccess)
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, request *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	resp = new(user.UserLoginResponse)

	exist := common.TestUserBloom(request.Username)

	// 用户名不存在
	if !exist {
		zap.L().Info("Check user exists info:", zap.Bool("exist", exist))
		resp.StatusCode = common.CodeUsernameNotFound
		resp.StatusMsg = common.MapErrMsg(common.CodeUsernameNotFound)
		return
	}

	// 用户名存在
	userModel, _, err := mysql.FindUserByName(request.Username)
	if err != nil {
		zap.L().Info("Find user by name err:", zap.Error(err))
		resp.StatusCode = common.CodeServerBusy
		resp.StatusMsg = common.MapErrMsg(common.CodeServerBusy)
		err = nil
		return
	}
	// 检查密码
	match := common.CheckPassword(request.Password, userModel.Password)
	if !match {
		zap.L().Info("User password wrong.")
		resp.StatusCode = common.CodeWrongLoginCredentials
		resp.StatusMsg = common.MapErrMsg(common.CodeWrongLoginCredentials)
		return
	}
	token := common.GenerateToken(userModel.ID, userModel.Name)
	resp.StatusCode = common.CodeSuccess
	resp.StatusMsg = common.MapErrMsg(common.CodeSuccess)

	resp.Token = token
	resp.UserId = int64(userModel.ID)
	// 将token存入redis
	redis.SetToken(userModel.ID, token)
	return
}

// GetUserInfoById implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfoById(ctx context.Context, request *user.UserInfoByIdRequest) (resp *user.UserInfoByIdResponse, err error) {
	resp = new(user.UserInfoByIdResponse)
	actionId := request.GetActorId()
	isLogged := false
	if actionId != 0 {
		isLogged = true
	}
	userId := request.GetUserId()
	if userId == 0 {
		resp.StatusCode = common.CodeInvalidParam
		resp.StatusMsg = common.MapErrMsg(common.CodeInvalidParam)
		return
	}
	name, exist := GetName(uint(userId))
	// 用户名不存在
	if !exist {
		resp.StatusCode = common.CodeUserNotFound
		resp.StatusMsg = common.MapErrMsg(common.CodeUserNotFound)
		return
	}

	followCountCh := make(chan int32)
	followerCountCh := make(chan int32)
	workCountCh := make(chan int32)
	favoriteCountCh := make(chan int32)
	totalFavoriteCountCh := make(chan int32)
	isFollowCh := make(chan bool)

	defer func() {
		close(followCountCh)
		close(followerCountCh)
		close(workCountCh)
		close(favoriteCountCh)
		close(totalFavoriteCountCh)
		close(isFollowCh)
	}()

	// 关注数
	go func() {
		followCount, _ := relationClient.GetFollowListCount(ctx, userId)
		followCountCh <- followCount
	}()

	// 粉丝数
	go func() {
		followerCount, _ := relationClient.GetFollowerListCount(ctx, userId)
		followerCountCh <- followerCount
	}()

	// 作品数
	go func() {
		workCount, _ := videoClient.GetWorkCount(ctx, userId)
		workCountCh <- workCount
	}()

	// 喜欢数
	go func() {
		favoriteCount, _ := favoriteClient.GetUserFavoriteCount(ctx, userId)
		favoriteCountCh <- favoriteCount
	}()

	// 总的被点赞数
	go func() {
		totalFavoriteCount, _ := favoriteClient.GetUserTotalFavoritedCount(ctx, userId)
		totalFavoriteCountCh <- totalFavoriteCount
	}()

	go func(isLog bool) {
		if isLog {
			isFollow, _ := relationClient.IsFollowing(ctx, &relation.IsFollowingRequest{
				ActorId: actionId,
				UserId:  userId,
			})
			isFollowCh <- isFollow
			return
		}
		isFollowCh <- false
	}(isLogged)

	var followCount, followerCount, workCount, favoriteCount, totalFavoriteCount int32
	var isFollow bool

	resp.SetUser(pack.User(userId))
	resp.User.SetName(name)
	// 从通道接收结果
	for receivedCount := 0; receivedCount < 6; receivedCount++ {
		select {
		case followCount = <-followCountCh:
			resp.User.SetFollowCount(followCount)
		case followerCount = <-followerCountCh:
			resp.User.SetFollowerCount(followerCount)
		case workCount = <-workCountCh:
			resp.User.SetWorkCount(workCount)
		case favoriteCount = <-favoriteCountCh:
			resp.User.SetFavoriteCount(favoriteCount)
		case totalFavoriteCount = <-totalFavoriteCountCh:
			resp.User.SetTotalFavorited(strconv.Itoa(int(totalFavoriteCount)))
		case isFollow = <-isFollowCh:
			resp.User.SetIsFollow(isFollow)
		case <-time.After(3 * time.Second):
			zap.L().Error("3s overtime.")
			break
		}
	}

	resp.StatusCode = common.CodeSuccess
	resp.StatusMsg = common.MapErrMsg(common.CodeSuccess)
	return
}

// GetName 根据userId获取用户名
func GetName(userId uint) (string, bool) {
	if val, err := cache.Get(strconv.Itoa(int(userId))); err == nil {
		return string(val), true
	}
	// 从redis中获取用户名
	// 1. 缓存中有数据, 直接返回
	if name, err := redis.GetNameByUserId(userId); err == nil {
		go func(uint, string) {
			cache.Set(strconv.Itoa(int(userId)), []byte(name))
		}(userId, name)
		return name, true
	}
	//缓存不存在，尝试从数据库中取
	if redis.AcquireUserLock(userId, redis.NameField) {
		defer redis.ReleaseUserLock(userId, redis.NameField)
		// 2. 缓存中没有数据，从数据库中获取
		userModel, exist, _ := mysql.FindUserByUserID(userId)
		if !exist {
			return "", false
		}
		// 将用户名写入redis
		err := redis.SetNameByUserId(userId, userModel.Name)
		if err != nil {
			zap.L().Error("将用户名写入redis失败：", zap.Error(err))
		}
		go func(uint, string) {
			cache.Set(strconv.Itoa(int(userId)), []byte(userModel.Name))
		}(userId, userModel.Name)
		return userModel.Name, true
	}
	// 重试
	time.Sleep(redis.RetryTime)
	return GetName(userId)
}
