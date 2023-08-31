package kafka

import (
	"context"
	"douyin/dal/model"
	"douyin/dal/mongo"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"time"
)

const (
	maxRetries     = 3     // 最大重试次数
	initialBackoff = 100   // 初始退避时间（毫秒）
	maxBackoff     = 10000 // 最大退避时间（毫秒）
)

type MessageMQ struct {
	MQ
}

var MessageMQInstance *MessageMQ

func InitMessageKafka() {
	MessageMQInstance = &MessageMQ{
		MQ{
			Topic:   "messages",
			GroupId: "message_group",
		},
	}

	// 创建 Message 业务的生产者和消费者实例
	MessageMQInstance.Producer = kafkaManager.NewProducer(MessageMQInstance.Topic)
	MessageMQInstance.Consumer = kafkaManager.NewConsumer(MessageMQInstance.Topic, MessageMQInstance.GroupId)

	go MessageMQInstance.Consume()
}

func (m *MessageMQ) Produce(message *model.Message) {
	// 在 Message 业务中使用 Kafka Manager
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		log.Println("kafka发送message失败：", err)
		return
	}
}

func (m *MessageMQ) Consume() {
	// 在 Message 业务中使用 Kafka Manager
	for {
		msg, err := m.Consumer.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Failed to read message:", err)
		}

		//fmt.Printf("Received message: %s\n", msg.Value)

		// 发送确认
		err = m.Consumer.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Println("Failed to commit message:", err)
		}

		message := new(model.Message)
		err = json.Unmarshal(msg.Value, message)
		if err != nil {
			log.Println("kafka获取message反序列化失败：", err)
		}

		// 重试插入数据库
		err = performDatabaseInsertWithRetry(msg.Value)
		if err != nil {
			zap.L().Error("Error inserting into database after retry:", zap.Error(err))
			return
			// 可以进行错误处理，如记录错误日志
		} else {
			// 确认消息已成功处理
			err := m.Consumer.CommitMessages(context.Background(), msg)
			if err != nil {
				zap.L().Error("Commit to kafka error", zap.Error(err))
				return
			}
		}
	}
}

func performDatabaseInsertWithRetry(data []byte) error {
	retries := 0
	backoff := initialBackoff

	for retries < maxRetries {
		err := performDatabaseInsert(data)
		if err == nil {
			return nil
		}

		// 记录错误日志
		fmt.Printf("Error inserting into database (retry %d): %v\n", retries+1, err)

		// 使用指数退避策略
		time.Sleep(time.Duration(backoff) * time.Millisecond)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}

		retries++
	}

	return fmt.Errorf("maximum retries reached")
}

func performDatabaseInsert(data []byte) error {
	// 假设这里是将 data 插入数据库的代码
	// 在插入操作时，可以设置适当的超时时间
	// 消息入库mongodb
	err := mongo.SendMessage(data)
	if err != nil {
		log.Println("kafka往mongo存入message失败：", err)
		return err
	}
	return err
}
