package main

import (
	"douyin/common"
	"douyin/config"
	"douyin/constant"
	"douyin/dal/mysql"
	favorite "douyin/kitex_gen/favorite/favoriteservice"
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
	//_ "net/http/pprof"
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
	// 初始化点赞模块的kafka
	kafka.InitFavoriteKafka()

	// 初始化Bloom Filter
	common.InitIsFavoriteFilter()
	common.LoadIsFavoriteToBloomFilter()
	common.InitFavoriteVideoIdFilter()
	common.LoadFavoriteVideoIdToBloomFilter()

	addr, err := net.ResolveTCPAddr("tcp", constant.FavoriteServicePort)
	if err != nil {
		panic(err)
	}

	//go func() {
	//	ip := "0.0.0.0:8899"
	//	if err := http.ListenAndServe(ip, nil); err != nil {
	//		log.Printf("start pprof failed on %s\n", ip)
	//	}
	//}()

	svr := favorite.NewServer(
		new(FavoriteServiceImpl),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.FavoriteServiceName}),
		server.WithRegistry(r),
		server.WithMuxTransport(),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
