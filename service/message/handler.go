package main

import (
	"context"
	"douyin/common"
	"douyin/constant"
	"douyin/dal/model"
	"douyin/dal/mongo"
	message "douyin/kitex_gen/message"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
	"douyin/mw/kafka"
	"douyin/service/message/pack"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"time"
)

var relationClient relationservice.Client

func init() {
	// Etcd 服务发现
	r, err := etcd.NewEtcdResolver([]string{constant.EtcdAddr})
	if err != nil {
		log.Fatal(err)
	}
	relationClient, err = relationservice.NewClient(
		constant.RelationServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.RelationServiceName}))
	if err != nil {
		log.Fatal(err)
	}
}

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageChat implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageChat(ctx context.Context, request *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	// 检查好友关系
	isFriend, err := relationClient.IsFriend(ctx, &relation.IsFriendRequest{
		ActorId: request.GetFromUserId(),
		UserId:  request.GetToUserId(),
	})
	if err != nil {
		zap.L().Error("relationClient.IsFriend error", zap.Error(err))
		return &message.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Internal server error"),
		}, err
	}
	// 不是好友关系
	if !isFriend {
		zap.L().Info("Not a friend, cannot see chat list")
		return &message.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("The user is not your friend"),
		}, nil
	}
	// 获取聊天记录
	msgList, err := mongo.GetMessageList(
		uint(request.GetFromUserId()),
		uint(request.GetToUserId()),
		request.GetPreMsgTime())
	if err != nil {
		zap.L().Error("mongo.GetMessageList error", zap.Error(err))
		return &message.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Internal server error"),
		}, err
	}
	// 封装数据
	packedMsgList := pack.Messages(msgList)
	return &message.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   thrift.StringPtr("success"),
		MessageList: packedMsgList,
	}, nil
}

// MessageAction implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageAction(ctx context.Context, request *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	// 检查好友关系
	isFriend, err := relationClient.IsFriend(ctx, &relation.IsFriendRequest{
		ActorId: request.GetFromUserId(),
		UserId:  request.GetToUserId(),
	})
	if err != nil {
		zap.L().Error("relationClient.IsFriend error", zap.Error(err))
		return &message.MessageActionResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("Internal server error"),
		}, err
	}
	// 不是好友关系
	if !isFriend {
		zap.L().Info("Not a friend, cannot see chat list")
		return &message.MessageActionResponse{
			StatusCode: 1,
			StatusMsg:  thrift.StringPtr("The user is not your friend"),
		}, nil
	}

	messageData := &model.Message{
		ID:         int64(common.GetUid()),
		FromUserId: uint(request.GetFromUserId()),
		ToUserId:   uint(request.GetToUserId()),
		Content:    request.GetContent(),
		CreateTime: time.Now().Unix(),
	}

	// 聊天记录发向kafka
	go kafka.MessageMQInstance.Produce(messageData)

	////往mongo发送聊天记录
	//err = mongo.SendMessage(messageData)
	//if err != nil {
	//	log.Println("mongo.SendMessage err:", err)
	//	return
	//}

	return &message.MessageActionResponse{
		StatusCode: 0,
		StatusMsg:  thrift.StringPtr("Success"),
	}, nil
}
