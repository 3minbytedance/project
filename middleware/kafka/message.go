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

var (
	messageProducer *kafka.Writer
	messageConsumer *kafka.Reader
	topic           = "messages"
)

func InitMessageKafka() {
	// 初始化 Kafka Manager
	brokers := []string{"localhost:9092"}
	kafkaManager := NewKafkaManager(brokers)

	// 创建 Message 业务的生产者和消费者实例
	messageProducer = kafkaManager.NewProducer(topic)
	messageConsumer = kafkaManager.NewConsumer(topic)

	go Consume()
}

func Produce(message *models.Message) {
	// 在 Message 业务中使用 Kafka Manager
	_ = kafkaManager.ProduceMessage(messageProducer, message)

}

func Consume() {
	// 在 Message 业务中使用 Kafka Manager
	//message := new(models.Message)
	for {
		msg, err := messageConsumer.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Failed to read message:", err)
		}

		fmt.Printf("Received message: %s\n", msg.Value)

		// 发送确认
		err = messageConsumer.CommitMessages(context.Background(), msg)
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
