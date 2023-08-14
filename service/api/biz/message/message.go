package message

import (
	"context"
	"douyin/constant"
	"douyin/kitex_gen/message/messageservice"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
)

var messageClient messageservice.Client

func init() {
	//r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(config.CommentServiceName),
	//	provider.WithExportEndpoint(config.EnvConfig.EXPORT_ENDPOINT),
	//	provider.WithInsecure(),
	//)

	var err error
	messageClient, err = messageservice.NewClient(
		constant.CommentServiceName,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(ctx context.Context, c *app.RequestContext) {

}

func Chat(ctx context.Context, c *app.RequestContext) {

}
