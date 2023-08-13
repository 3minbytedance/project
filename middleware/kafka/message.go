package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"project/dao/mongo"
	"project/models"
)

type MessageMQ struct {
	Topic    string
	GroupId  string
	Producer *kafka.Writer
	Consumer *kafka.Reader
}

var MessageMQInstance *MessageMQ

func InitMessageKafka() {
	MessageMQInstance = &MessageMQ{
		Topic:   "messages",
		GroupId: "message_group",
	}

	// 创建 Message 业务的生产者和消费者实例
	MessageMQInstance.Producer = kafkaManager.NewProducer(MessageMQInstance.Topic)
	MessageMQInstance.Consumer = kafkaManager.NewConsumer(MessageMQInstance.Topic, MessageMQInstance.GroupId)

	go MessageMQInstance.Consume()
}

func (m *MessageMQ) Produce(message *models.Message) {
	// 在 Message 业务中使用 Kafka Manager
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		log.Println("kafka发送message失败：", err)
		return
	}
}

func (m *MessageMQ) Consume() {
	// 在 Message 业务中使用 Kafka Manager
	//message := new(models.Message)
	for {
		msg, err := m.Consumer.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Failed to read message:", err)
		}

		fmt.Printf("Received message: %s\n", msg.Value)

		// 发送确认
		err = m.Consumer.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Println("Failed to commit message:", err)
		}

		message := new(models.Message)
		err = json.Unmarshal(msg.Value, message)
		if err != nil {
			log.Println("kafka获取message反序列化失败：", err)
		}
		// 消息入库mongodb
		err = mongo.SendMessage(message)
		if err != nil {
			log.Println("kafka往mongo存入message失败：", err)
			return
		}
	}
}
