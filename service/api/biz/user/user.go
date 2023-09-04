package user

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

var userClient userservice.Client
var illegalCharRe *regexp.Regexp

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
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}),
		client.WithMuxConnection(2),
	)
	if err != nil {
		log.Fatal(err)
	}
	illegalCharRe = regexp.MustCompile(`[?;=]`)
}

func Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	statusCode := checkUserRegisterInfo(username, password)

	if statusCode != common.CodeSuccess {
		c.JSON(http.StatusOK, user.UserRegisterResponse{
			StatusCode: statusCode,
			StatusMsg:  common.MapErrMsg(statusCode),
		})
		return
	}

	resp, err := userClient.Register(ctx, &user.UserRegisterRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		zap.L().Error("Invoke userClient Register err:", zap.Error(err))
		c.JSON(http.StatusOK, user.UserRegisterResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  common.MapErrMsg(resp.StatusCode),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func Login(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	resp, err := userClient.Login(ctx, &user.UserLoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		zap.L().Error("Invoke userClient Login err:", zap.Error(err))
		c.JSON(http.StatusOK, user.UserLoginResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  common.MapErrMsg(resp.StatusCode),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func Info(ctx context.Context, c *app.RequestContext) {
	actorId, _ := c.Get(common.ContextUserIDKey)
	zap.L().Info("Info", zap.Uint("actorID", actorId.(uint)))
	userId := c.Query("user_id")
	userIdUint, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		zap.L().Error("Parse userId error", zap.Error(err))
		c.JSON(http.StatusOK, user.UserInfoByIdResponse{
			StatusCode: common.CodeInvalidParam,
			StatusMsg:  common.MapErrMsg(common.CodeInvalidParam),
		})
		return
	}

	resp, err := userClient.GetUserInfoById(ctx, &user.UserInfoByIdRequest{
		ActorId: int64(actorId.(uint)),
		UserId:  int64(userIdUint),
	})

	if err != nil {
		zap.L().Error("Invoke userClient getUserInfoById err:", zap.Error(err))
		c.JSON(http.StatusOK, user.UserInfoByIdResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  common.MapErrMsg(resp.StatusCode),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func checkUserRegisterInfo(username string, password string) int32 {
	if illegalCharRe.MatchString(username) {
		return common.CodeInvalidRegisterUsername
	}

	if len(username) == 0 || len(username) > 32 {
		return common.CodeInvalidRegisterUsername
	}

	if len(password) < 6 || len(password) > 32 {
		return common.CodeInvalidRegisterPassword
	}

	return common.CodeSuccess
}
