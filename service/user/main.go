package main

import (
	"douyin/common"
	"douyin/config"
	"douyin/constant"
	"douyin/dal/mysql"
	user "douyin/kitex_gen/user/userservice"
	"douyin/logger"
	"douyin/mw"
	"fmt"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {
	// OpenTelemetry 链路跟踪
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(config.CommentServiceName),
	//	provider.WithExportEndpoint("localhost:4317"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())

	// Etcd 服务发现
	r, err := etcd.NewEtcdRegistry([]string{constant.EtcdAddr}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	// 加载配置
	if err := config.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	// 加载日志
	if err := logger.Init(config.Conf.LogConfig, config.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 初始化数据库: mysql
	if err := mysql.Init(config.Conf); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
		return
	}

	// 初始化中间件: redis + kafka
	if err := mw.Init(config.Conf); err != nil {
		fmt.Printf("Init middleware failed, err:%v\n", err)
		return
	}

	// 初始化Bloom Filter
	common.InitBloomFilter()
	common.LoadUsernamesToBloomFilter()

	addr, err := net.ResolveTCPAddr("tcp", constant.UserServicePort)
	if err != nil {
		panic(err)
	}

	svr := user.NewServer(
		new(UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.UserServiceName}),
		server.WithRegistry(r),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
