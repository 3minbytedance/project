package kafka

import (
	"context"
	"douyin/dal/mysql"
	"encoding/json"
	"fmt"
)

type FollowMessage struct {
	Type     int
	UserId   uint
	FollowId uint
}

type FollowMQ struct {
	MQ
}

var (
	FollowMQInstance *FollowMQ
)

func InitFollowKafka() {
	FollowMQInstance = &FollowMQ{
		MQ{
			Topic:   "follows",
			GroupId: "follow_group",
		},
	}

	// 创建 Follow 业务的生产者和消费者实例
	FollowMQInstance.Producer = kafkaManager.NewProducer(FollowMQInstance.Topic)
	FollowMQInstance.Consumer = kafkaManager.NewConsumer(FollowMQInstance.Topic, FollowMQInstance.GroupId)

	go FollowMQInstance.Consume()
}

// ProduceAddFollowMsg 发布添加关注的消息, 向mysql中添加关注时, 调用此方法
func (m *FollowMQ) ProduceAddFollowMsg(userId, followId uint) {
	message := &FollowMessage{
		Type:     0,
		UserId:   userId,
		FollowId: followId,
	}
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		fmt.Println("kafka发送添加关注的消息失败：", err)
		return
	}
}

// ProduceDelFollowMsg 发布删除关注的消息, mysql删除关注时, 调用此方法
func (m *FollowMQ) ProduceDelFollowMsg(userId, followId uint) {
	message := &FollowMessage{
		Type:     1,
		UserId:   userId,
		FollowId: followId,
	}
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		fmt.Println("kafka发送删除关注的消息失败：", err)
		return
	}
}

// Consume 消费者消费消息
func (m *FollowMQ) Consume() {
	for {
		msg, err := m.Consumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("[FollowMQ]从消息队列中读取消息失败:", err)
		}

		// 解析消息
		var message FollowMessage
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			fmt.Println("[FollowMQ]解析消息失败:", err)
			continue
		}

		// 根据消息类型, 执行不同的操作
		switch message.Type {
		case 0:
			// 添加点赞
			err = mysql.AddFollow(message.UserId, message.FollowId)
			if err != nil {
				fmt.Println("[FollowMQ]添加点赞失败:", err)
				continue
			}
		case 1:
			// 删除点赞
			err = mysql.DeleteFollowById(message.UserId, message.FollowId)
			if err != nil {
				fmt.Println("[FollowMQ]删除点赞失败:", err)
				continue
			}
		}
	}
}
