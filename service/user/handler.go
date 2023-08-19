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
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/video/videoservice"
	"douyin/mw/redis"
	"douyin/service/user/pack"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"math/rand"
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
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
	)
	videoClient, err = videoservice.NewClient(
		constant.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.VideoServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, request *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {

	// 检查注册信息是否合理
	resp = new(user.UserRegisterResponse)
	statusCode, statusMsg := CheckUserRegisterInfo(request.Username, request.Password)
	resp.StatusCode = statusCode
	resp.StatusMsg = thrift.StringPtr(statusMsg)

	if statusCode != 0 {
		return
	}

	userData := model.User{}
	userData.Name = request.Username

	// 用户名存入Bloom Filter
	common.AddToBloom(request.Username)

	// 将信息存储到数据库中
	salt := fmt.Sprintf("%06d", rand.Int())
	userData.Salt = salt
	userData.Password = common.MakePassword(request.Password, salt)

	// 数据入库
	userId, err := mysql.CreateUser(&userData)
	if err != nil {
		zap.L().Error("Create user err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server internal error.")
		return
	}

	// todo int32？
	resp.UserId = int32(userId)
	resp.Token = common.GenerateToken(userId, request.Username)

	err = redis.SetToken(resp.Token, userId)
	if err != nil {
		zap.L().Error("Set token err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server internal error.")
		return
	}
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, request *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	resp = new(user.UserLoginResponse)

	exist := common.TestBloom(request.Username)
	// todo 判断存在但不存在 判断不存在但存在？
	// 用户名不存在
	if !exist {
		zap.L().Info("Check user exists info:", zap.Bool("exist", exist))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Username not exist")
		return
	}

	// 用户名存在
	userModel, _, err := mysql.FindUserByName(request.Username)
	if err != nil {
		zap.L().Info("Find user by name err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server Internal error")
		return
	}
	// 检查密码
	match := common.CheckPassword(request.Password, userModel.Salt, userModel.Password)
	if !match {
		zap.L().Info("User password wrong.")
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Wrong password.")
		return
	}
	token := common.GenerateToken(userModel.ID, userModel.Name)
	resp.StatusCode = 0
	resp.StatusMsg = thrift.StringPtr("success")
	resp.Token = token
	resp.UserId = int32(userModel.ID)

	err = redis.SetToken(token, userModel.ID)
	if err != nil {
		zap.L().Error("Set token err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server internal error.")
		return
	}
	return
}

// GetUserInfoById
// 查询userId 的信息，并判断当前actionId是否和userId关注
func (s *UserServiceImpl) GetUserInfoById(ctx context.Context, request *user.UserInfoByIdRequest) (resp *user.UserInfoByIdResponse, err error) {

	// todo redis

	resp = new(user.UserInfoByIdResponse)
	actionId := request.GetActorId()
	isLogged := false
	if actionId != 0 {
		isLogged = true
	}
	userId := request.GetUserId()
	if userId == 0 {
		userId = actionId
	}
	name, exist := GetName(uint(userId))
	// 用户名不存在
	if !exist {
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("User ID not exist")
		return
	}

	// 关注数 粉丝数
	followCount, _ := relationClient.GetFollowListCount(ctx, userId)
	followerCount, _ := relationClient.GetFollowerListCount(ctx, userId)

	resp.StatusCode = 0
	resp.StatusMsg = thrift.StringPtr("success")
	// 作品数
	workCount, _ := videoClient.GetWorkCount(ctx, userId)
	// 喜欢数
	favoriteCount, _ := favoriteClient.GetUserTotalFavoritedCount(ctx, userId)
	// 总的被点赞数
	totalFavoriteCount, _ := favoriteClient.GetUserTotalFavoritedCount(ctx, userId)
	// 检查是否已关注

	zap.L().Info("IDS", zap.Any("actorId", actionId), zap.Any("userId", userId))
	isFollow := false
	//已登录
	if isLogged {
		isFollow, err = relationClient.IsFollowing(ctx, &relation.IsFollowingRequest{
			ActorId: actionId,
			UserId:  userId,
		})
		if err != nil {
			zap.L().Error("relationClient err:", zap.Error(err))
			resp.StatusCode = 1
			resp.StatusMsg = thrift.StringPtr("Server Internal error")
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
	resp.User.SetTotalFavorited(totalFavoriteCount)
	return
}

// GetName 根据userId获取用户名
func GetName(userId uint) (string, bool) {
	// 从redis中获取用户名
	// 1. 缓存中有数据, 直接返回
	if redis.IsExistUserField(userId, redis.NameField) {
		name, err := redis.GetNameByUserId(userId)
		if err != nil {
			log.Println("从redis中获取用户名失败：", err)
		}
		return name, true
	}

	// 2. 缓存中没有数据，从数据库中获取
	userModel, exist, _ := mysql.FindUserByUserID(userId)
	if !exist {
		return "", false
	}
	// 将用户名写入redis
	go func() {
		err := redis.SetNameByUserId(userId, userModel.Name)
		if err != nil {
			log.Println("将用户名写入redis失败：", err)
		}
	}()
	return userModel.Name, true
}

func CheckUserRegisterInfo(username string, password string) (int32, string) {

	if len(username) == 0 || len(username) > 32 {
		return 1, "username is illegal."
	}

	if len(password) <= 6 || len(password) > 32 {
		return 2, "password is illegal"
	}

	_, exist, err := mysql.FindUserByName(username)
	if err != nil {
		zap.L().Error("Find user by name:", zap.Error(err))
		return 1, "Server internal error."
	}
	// 检查用户名是否存在
	if exist {
		zap.L().Info("User already exists")
		return 1, "User already exists."
	}

	return 0, "success"
}
