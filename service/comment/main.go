package main

import (
	"douyin/common"
	"douyin/config"
	"douyin/constant"
	"douyin/dal/mysql"
	comment "douyin/kitex_gen/comment/commentservice"
	"douyin/logger"
	"douyin/mw/redis"
	"douyin/mw/rocketMQ"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"net"
	"strconv"
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
		zap.L().Error("Init redis failed, err:%v\n", zap.Error(err))
		return
	}
	//if err := kafka.Init(config.Conf); err != nil {
	//	zap.L().Error("Init kafka failed, err:%v\n", zap.Error(err))
	//	return
	//}
	//// 初始化comment模块的kafka
	//kafka.InitCommentKafka()
	if err := rocketMQ.Init(config.Conf); err != nil {
		zap.L().Error("Init rocketMQ failed, err:%v\n", zap.Error(err))
		return
	}
	// 初始化comment模块的rocketMQ
	rocketMQ.InitCommentMQ()

	// 初始化敏感词过滤器
	if err := common.InitSensitiveFilter(); err != nil {
		zap.L().Error("Init sensitive filter error", zap.Error(err))
	}

	// 初始化Bloom Filter
	common.InitCommentBloomFilter()
	common.LoadCommentVideoIdToBloomFilter()

	// 初始化雪花算法ID
	nodeNum, err := strconv.ParseInt(config.Conf.Node, 10, 64)
	if err != nil {
		zap.L().Error("Snowflake node num failed, err:%v\n", zap.Error(err))
		return
	}

	err = common.InitSnowflake(nodeNum)
	if err != nil {
		zap.L().Error("Snowflake node failed, err:%v\n", zap.Error(err))
		return
	}

	addr, err := net.ResolveTCPAddr("tcp", constant.CommentServicePort)
	if err != nil {
		panic(err)
	}

	svr := comment.NewServer(
		new(CommentServiceImpl),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
		server.WithRegistry(r),
		server.WithMuxTransport(),
	)
	err = svr.Run()

	if err != nil {
		log.Fatal(err)
	}
}
