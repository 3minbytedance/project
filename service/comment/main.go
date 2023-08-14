package main

import (
	"douyin/constant"
	comment "douyin/kitex_gen/comment/commentservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
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

	svr := comment.NewServer(
		new(CommentServiceImpl),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.CommentServiceName}),
		server.WithRegistry(r),
	)
	err = svr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
