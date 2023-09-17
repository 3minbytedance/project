package rocketMQ

import (
	"context"
	"douyin/config"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"strconv"
)

type Manager struct {
	Brokers []string
}

var rocketMQManager *Manager

type MQ struct {
	Topic    string
	GroupId  string
	Producer rocketmq.Producer
	Consumer rocketmq.PushConsumer
}

func Init(appConfig *config.AppConfig) (err error) {
	var conf *config.RocketMQConfig
	conf = appConfig.Remote.RocketMQConfig

	brokerUrl := conf.Address + ":" + strconv.Itoa(conf.Port)
	// 初始化 Manager
	brokers := []string{brokerUrl}
	rocketMQManager = NewRocketMQManager(brokers)

	return nil
}

func NewRocketMQManager(brokers []string) *Manager {
	return &Manager{
		Brokers: brokers,
	}
}

func (m *Manager) NewProducer(groupName, topic string) rocketmq.Producer {
	c, _ := rocketmq.NewProducer(
		producer.WithNameServer(m.Brokers), // 接入点地址
		producer.WithRetry(2),              // 重试次数
		producer.WithCreateTopicKey(topic),
		producer.WithGroupName(groupName),
	)
	return c
}

func (m *Manager) NewConsumer(groupId string) rocketmq.PushConsumer {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer(m.Brokers), // 接入点地址
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(groupId), // 分组名称
	)
	return c
}

// ProduceMessage 向 RocketMQ 写入消息的公共函数, 由于不同业务的消息格式不同, 所以使用 interface{} 代替
func (m *Manager) ProduceMessage(p rocketmq.Producer, message interface{}, topic string) (*primitive.SendResult, error) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: topic,
		Body:  messageBytes,
	})
	return result, err
}
