package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/dal/model"
	"douyin/dal/mysql"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
	user "douyin/kitex_gen/user"
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
	// 用户名不存在
	if !exist {
		zap.L().Info("Check user exists info:", zap.Bool("exist", exist))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Username not exist")
		return
	}

	// 用户名存在
	user, _, err := mysql.FindUserByName(request.Username)
	if err != nil {
		zap.L().Info("Find user by name err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server Internal error")
		return
	}
	// 检查密码
	match := common.CheckPassword(request.Password, user.Salt, user.Password)
	if !match {
		zap.L().Info("User password wrong.")
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Wrong password.")
		return
	}
	token := common.GenerateToken(user.ID, user.Name)
	resp.StatusCode = 0
	resp.StatusMsg = thrift.StringPtr("success")
	resp.Token = token
	resp.UserId = int32(user.ID)

	err = redis.SetToken(token, user.ID)
	if err != nil {
		zap.L().Error("Set token err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server internal error.")
		return
	}
	return
}

// GetUserInfoById implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfoById(ctx context.Context, request *user.UserInfoByIdRequest) (resp *user.UserInfoByIdResponse, err error) {
	resp = new(user.UserInfoByIdResponse)
	// userId不为0 -> 查询userId的用户信息，顺便查是不是actorId的关注，然后设置isFavorite
	// userId为0 -> 单纯查询actorId信息
	queryId := request.GetActorId()
	if request.GetUserId() != 0 {
		queryId = request.GetUserId()
	}
	user, exist, err := mysql.FindUserByUserID(uint(queryId))

	if err != nil {
		zap.L().Error("Check user exists err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server Internal error")
		return
	}
	// 用户名不存在
	if !exist {
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("User ID not exist")
		return
	}
	resp.StatusCode = 0
	resp.StatusMsg = thrift.StringPtr("success")
	resp.User = pack.User(&user)

	// 检查是否已关注
	zap.L().Info("IDS", zap.Any("actorId", request.ActorId), zap.Any("userId", request.UserId))
	relationResp, err := relationClient.IsFollowing(ctx, &relation.IsFollowingRequest{
		ActorId: request.GetActorId(),
		UserId:  request.GetUserId(),
	})
	if err != nil {
		zap.L().Error("Check user exists err:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server Internal error")
		return
	}
	resp.User.IsFollow = relationResp.GetResult_()
	return
}

// GetUserInfoByName implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfoByName(ctx context.Context, request *user.UserInfoByNameRequest) (resp *user.UserInfoByNameResponse, err error) {
	// TODO: Your code here...
	return
}

// CheckUserExists implements the UserServiceImpl interface.
func (s *UserServiceImpl) CheckUserExists(ctx context.Context, request *user.UserExistsRequest) (resp *user.UserExistsResponse, err error) {
	// TODO: Your code here...
	return
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
