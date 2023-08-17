package main

import (
	"context"
	"douyin/common"
	"douyin/dal/model"
	"douyin/dal/mysql"
	user "douyin/kitex_gen/user"
	"douyin/mw/redis"
	"douyin/service/user/pack"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"go.uber.org/zap"
	"math/rand"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, request *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	resp = new(user.UserRegisterResponse)
	_, exist, err := mysql.FindUserByName(request.Username)
	if err != nil {
		zap.L().Error("Find user by name:", zap.Error(err))
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("Server internal error.")
		return
	}
	// 检查用户名是否存在
	if exist {
		zap.L().Info("User already exists")
		resp.StatusCode = 1
		resp.StatusMsg = thrift.StringPtr("User already exists.")
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

	resp.StatusCode = 0
	resp.StatusMsg = thrift.StringPtr("success")
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
	user, exist, err := mysql.FindUserByUserID(uint(request.UserId))
	if err != nil {
		zap.L().Info("Check user exists err:", zap.Error(err))
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
	// todo: 封装好友关系等信息
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
