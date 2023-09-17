package rocketMQ

import (
	"douyin/dal/model"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"log"
)

type FavoriteMQ struct {
	MQ
}

var (
	FavoriteMQInstance *FavoriteMQ
)

func InitFavoriteMQ() rocketmq.PushConsumer {
	rlog.SetLogLevel("error")
	FavoriteMQInstance = &FavoriteMQ{
		MQ{
			Topic:   "favorite",
			GroupId: "favorite_group",
		},
	}

	// 创建 Comment 业务的生产者和消费者实例
	FavoriteMQInstance.Producer = rocketMQManager.NewProducer(FavoriteMQInstance.GroupId, FavoriteMQInstance.Topic)
	err := FavoriteMQInstance.Producer.Start()
	if err != nil {
		panic("启动favorite producer 失败")
	}

	FavoriteMQInstance.Consumer = rocketMQManager.NewConsumer(FavoriteMQInstance.GroupId)

	return FavoriteMQInstance.Consumer
}

// ProduceFavoriteMsg 生产点赞的消息
func (m *FavoriteMQ) ProduceFavoriteMsg(message *model.FavoriteAction) error {
	_, err := rocketMQManager.ProduceMessage(m.Producer, message, m.Topic)

	if err != nil {
		log.Println("rocketMQ 发送favorite的消息失败：", err)
		return err
	}
	return nil
}
