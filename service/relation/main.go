package main

import (
	"douyin/config"
	"douyin/constant"
	"douyin/dal/mysql"
	relation "douyin/kitex_gen/relation/relationservice"
	"douyin/logger"
	"douyin/mw/redis"
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

	if err := mysql.Init(config.Conf); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
		return
	}

	// 初始化Redis
	if err := redis.Init(config.Conf); err != nil {
		fmt.Printf("Init redis failed, err:%v\n", err)
		return
	}

	addr, err := net.ResolveTCPAddr("tcp", constant.RelationServicePort)
	if err != nil {
		panic(err)
	}

	svr := relation.NewServer(
		new(RelationServiceImpl),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.RelationServiceName}),
		server.WithRegistry(r),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
