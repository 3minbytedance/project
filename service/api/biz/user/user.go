package user

import (
	"context"
	"douyin/constant"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"net/http"
)

var userClient userservice.Client

func init() {
	// OpenTelemetry 链路跟踪 还没配置好，先注释
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(config.CommentServiceName),
	//	provider.WithExportEndpoint("localhost:4317"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())

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
	)
	if err != nil {
		log.Fatal(err)
	}
}

func Register(ctx context.Context, c *app.RequestContext) {

	username := c.Query("username")
	password := c.Query("password")

	resp, err := userClient.Register(ctx, &user.UserRegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		zap.L().Error("Invoke userClient Register err:", zap.Error(err))
		c.JSON(http.StatusOK, &user.UserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Server internal error"),
			UserId:     0,
			Token:      "",
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
		c.JSON(http.StatusOK, &user.UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Server Internal error"),
			UserId:     0,
			Token:      "",
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}
