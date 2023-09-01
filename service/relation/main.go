package main

import (
	"douyin/common"
	"douyin/config"
	"douyin/constant"
	"douyin/dal/mysql"
	relation "douyin/kitex_gen/relation/relationservice"
	"douyin/logger"
	"douyin/mw/kafka"
	"douyin/mw/redis"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func main() {
	// Etcd 服务发现
	r, err := etcd.NewEtcdRegistry([]string{constant.EtcdAddr}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	// 加载配置
	if err := config.Init(); err != nil {
		zap.L().Error("Load config failed, err:%v\n", zap.Error(err))
		return
	}
	// 加载日志
	if err := logger.Init(config.Conf.LogConfig, config.Conf.Mode); err != nil {
		zap.L().Error("Init logger failed, err:%v\n", zap.Error(err))
		return
	}

	// 初始化数据库: mysql
	if err := mysql.Init(config.Conf); err != nil {
		zap.L().Error("Init mysql failed, err:%v\n", zap.Error(err))
		return
	}

	// 初始化中间件: redis + kafka
	if err := redis.Init(config.Conf); err != nil {
		zap.L().Error("Init middleware failed, err:%v\n", zap.Error(err))
		return
	}
	if err := kafka.Init(config.Conf); err != nil {
		zap.L().Error("Init kafka failed, err:%v\n", zap.Error(err))
		return
	}
	// 初始化关注模块的kafka
	kafka.InitFollowKafka()

	// 初始化Bloom Filter
	common.InitRelationFollowIdFilter()
	common.LoadRelationFollowIdToBloomFilter()
	common.InitRelationFollowerIdFilter()
	common.LoadRelationFollowerIdToBloomFilter()

	addr, err := net.ResolveTCPAddr("tcp", constant.RelationServicePort)
	if err != nil {
		panic(err)
	}

	//pprof 监听
	// 对阻塞超过1纳秒的 goroutine 进行数据采集
	runtime.SetBlockProfileRate(1)
	go func() {
		http.ListenAndServe(":8002", nil)
	}()

	svr := relation.NewServer(
		new(RelationServiceImpl),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.RelationServiceName}),
		server.WithRegistry(r),
		server.WithMuxTransport(),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
