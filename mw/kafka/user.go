package kafka

import (
	"context"
	"douyin/dal/model"
	"douyin/dal/mysql"
	"encoding/json"
	"fmt"
)

type UserMQ struct {
	MQ
}

var (
	UserMQInstance *UserMQ
)

func InitUserKafka() {
	UserMQInstance = &UserMQ{
		MQ{
			Topic:   "users",
			GroupId: "user_group",
		},
	}

	// 创建 User 业务的生产者和消费者实例
	UserMQInstance.Producer = kafkaManager.NewProducer(UserMQInstance.Topic)
	UserMQInstance.Consumer = kafkaManager.NewConsumer(UserMQInstance.Topic, UserMQInstance.GroupId)

	go UserMQInstance.Consume()
}

// ProduceCreateUserMsg 发布创建用户的消息, 向mysql中创建用户时, 调用此方法
func (m *UserMQ) ProduceCreateUserMsg(message *models.User) {
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		fmt.Println("kafka发送添加点赞的消息失败：", err)
		return
	}
}

// Consume 消费创建用户的消息
func (m *UserMQ) Consume() {
	for {
		msg, err := m.Consumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("[UserMQ]从消息队列中读取消息失败:", err)
		}

		// 解析消息
		var user models.User
		err = json.Unmarshal(msg.Value, &user)
		if err != nil {
			fmt.Println("[UserMQ]解析消息失败:", err)
			continue
		}

		// 创建用户
		_, err = mysql.CreateUser(&user)
		if err != nil {
			fmt.Println("[UserMQ]创建用户失败:", err)
		}
	}
}
