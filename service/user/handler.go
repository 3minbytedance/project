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
	"douyin/mw/redis"
	"douyin/service/user/pack"
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.RelationServiceName}))
	favoriteClient, err = favoriteservice.NewClient(
		constant.FavoriteServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
	)
	videoClient, err = videoservice.NewClient(
		constant.VideoServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.VideoServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
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
	statusCode, statusMsg := CheckUserRegisterInfo(request.Username, request.Password)
	resp.StatusCode = statusCode
	resp.StatusMsg = statusMsg

	if statusCode != common.CodeSuccess {
		return
	}

	userData := model.User{}
	userData.Name = request.Username
	userData.ID = common.GetUid()

	// 将信息存储到数据库中
	userData.Password, _ = common.MakePassword(request.Password)

	// 数据入库
	userId, err := mysql.CreateUser(&userData)
	if err != nil {
		zap.L().Error("Create user err:", zap.Error(err))
		resp.StatusCode = common.CodeServerBusy
		resp.StatusMsg = common.MapErrMsg(common.CodeServerBusy)
		return
	}

	resp.UserId = int64(userId)
	resp.Token = common.GenerateToken(userId, request.Username)

	// 将token存入redis
	redis.SetToken(userId, resp.Token)

	// 用户名存入Bloom Filter
	common.AddToBloom(request.Username)
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, request *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	resp = new(user.UserLoginResponse)

	exist := common.TestBloom(request.Username)

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

	// 关注数 粉丝数
	followCount, _ := relationClient.GetFollowListCount(ctx, userId)
	followerCount, _ := relationClient.GetFollowerListCount(ctx, userId)

	resp.StatusCode = common.CodeSuccess
	resp.StatusMsg = common.MapErrMsg(common.CodeSuccess)
	// 作品数
	workCount, _ := videoClient.GetWorkCount(ctx, userId)
	// 喜欢数
	favoriteCount, _ := favoriteClient.GetUserFavoriteCount(ctx, userId)
	// 总的被点赞数
	totalFavoriteCount, _ := favoriteClient.GetUserTotalFavoritedCount(ctx, userId)
	// 检查是否已关注

	//zap.L().Info("IDS", zap.Any("actorId", actionId), zap.Any("userId", userId))
	isFollow := false
	//已登录
	if isLogged {
		isFollow, err = relationClient.IsFollowing(ctx, &relation.IsFollowingRequest{
			ActorId: actionId,
			UserId:  userId,
		})
		if err != nil {
			zap.L().Error("relationClient err:", zap.Error(err))
			resp.StatusCode = common.CodeServerBusy
			resp.StatusMsg = common.MapErrMsg(common.CodeServerBusy)
			return
		}
	}

	resp.SetUser(pack.User(userId))
	resp.User.SetName(name)
	resp.User.SetFollowCount(followCount)
	resp.User.SetFollowerCount(followerCount)
	resp.User.SetIsFollow(isFollow)
	resp.User.SetWorkCount(workCount)
	resp.User.SetFavoriteCount(favoriteCount)
	resp.User.SetTotalFavorited(strconv.Itoa(int(totalFavoriteCount)))
	return
}

// GetName 根据userId获取用户名
func GetName(userId uint) (string, bool) {
	// 从redis中获取用户名
	// 1. 缓存中有数据, 直接返回
	if redis.IsExistUserField(userId, redis.NameField) {
		name, err := redis.GetNameByUserId(userId)
		if err != nil {
			zap.L().Error("从redis中获取用户名失败：", zap.Error(err))
			return "", false
		}
		return name, true
	}
	//缓存不存在，尝试从数据库中取
	if redis.AcquireUserLock(userId, redis.NameField) {
		defer redis.ReleaseUserLock(userId, redis.NameField)
		//double check
		if redis.IsExistUserField(userId, redis.NameField) {
			name, err := redis.GetNameByUserId(userId)
			if err != nil {
				zap.L().Error("从redis中获取用户名失败：", zap.Error(err))
				return "", false
			}
			return name, true
		}
		// 2. 缓存中没有数据，从数据库中获取
		userModel, exist, _ := mysql.FindUserByUserID(userId)
		if !exist {
			return "", false
		}
		// 将用户名写入redis
		err := redis.SetNameByUserId(userId, userModel.Name)
		if err != nil {
			log.Println("将用户名写入redis失败：", err)
		}
		return userModel.Name, true
	}
	// 重试
	time.Sleep(100 * time.Millisecond)
	return GetName(userId)
}

func CheckUserRegisterInfo(username string, password string) (int32, string) {

	if len(username) == 0 || len(username) > 32 {
		return common.CodeInvalidRegisterUsername, common.MapErrMsg(common.CodeInvalidRegisterUsername)
	}

	if len(password) < 6 || len(password) > 32 {
		return common.CodeInvalidRegisterPassword, common.MapErrMsg(common.CodeInvalidRegisterPassword)
	}

	return common.CodeSuccess, common.MapErrMsg(common.CodeSuccess)
}
