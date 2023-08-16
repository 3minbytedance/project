package kafka

import (
	"context"
	"douyin/dal/mysql"
	"encoding/json"
	"fmt"
)

// FavoriteMessage 往kafka中发送的消息
type FavoriteMessage struct {
	Type    int
	UserId  uint
	VideoId uint
}

type FavoriteMQ struct {
	MQ
}

var (
	FavoriteMQInstance *FavoriteMQ
)

func InitFavoriteKafka() {
	FavoriteMQInstance = &FavoriteMQ{
		MQ{
			Topic:   "favorites",
			GroupId: "favorite_group",
		},
	}

	// 创建 Favorite 业务的生产者和消费者实例
	FavoriteMQInstance.Producer = kafkaManager.NewProducer(FavoriteMQInstance.Topic)
	FavoriteMQInstance.Consumer = kafkaManager.NewConsumer(FavoriteMQInstance.Topic, FavoriteMQInstance.GroupId)

	go FavoriteMQInstance.Consume()
}

// ProduceAddFavoriteMsg 发布添加点赞的消息, 向mysql中添加点赞时, 调用此方法
func (m *FavoriteMQ) ProduceAddFavoriteMsg(userId, videoId uint) {
	message := &FavoriteMessage{
		Type:    0,
		UserId:  userId,
		VideoId: videoId,
	}
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		fmt.Println("kafka发送添加点赞的消息失败：", err)
		return
	}
}

// ProduceDelFavoriteMsg 发布删除点赞的消息, mysql删除点赞时, 调用此方法
func (m *FavoriteMQ) ProduceDelFavoriteMsg(userId, videoId uint) {
	message := &FavoriteMessage{
		Type:    1,
		UserId:  userId,
		VideoId: videoId,
	}
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		fmt.Println("kafka发送删除点赞的消息失败：", err)
		return
	}
}

// Consume 消费添加或者删除点赞的消息
func (m *FavoriteMQ) Consume() {
	for {
		msg, err := m.Consumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("[FavoriteMQ]从消息队列中读取消息失败:", err)
		}

		// 解析消息
		var message FavoriteMessage
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			fmt.Println("[FavoriteMQ]解析消息失败:", err)
			continue
		}

		// 根据消息类型, 执行不同的操作
		switch message.Type {
		case 0:
			// 添加点赞
			err = mysql.AddUserFavorite(message.UserId, message.VideoId)
			if err != nil {
				fmt.Println("[FavoriteMQ]添加点赞失败:", err)
				continue
			}
		case 1:
			// 删除点赞
			err = mysql.DeleteUserFavorite(message.UserId, message.VideoId)
			if err != nil {
				fmt.Println("[FavoriteMQ]删除点赞失败:", err)
				continue
			}
		}
	}
}
