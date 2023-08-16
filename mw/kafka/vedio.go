package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type VideoMessage struct {
	VideoName string
	AuthorId  uint
	Title     string
}

type VideoMQ struct {
	MQ
}

var (
	VideoMQInstance *VideoMQ
)

func InitVideoKafka() {
	VideoMQInstance = &VideoMQ{
		MQ{
			Topic:   "videos",
			GroupId: "video_group",
		},
	}

	// 创建 Video 业务的生产者和消费者实例
	VideoMQInstance.Producer = kafkaManager.NewProducer(VideoMQInstance.Topic)
	VideoMQInstance.Consumer = kafkaManager.NewConsumer(VideoMQInstance.Topic, VideoMQInstance.GroupId)

	go VideoMQInstance.Consume()
}

// Produce 发布将本地视频上传到OSS的消息
func (m *VideoMQ) Produce(message *VideoMessage) {
	err := kafkaManager.ProduceMessage(m.Producer, message)
	if err != nil {
		fmt.Println("kafka发送添加视频的消息失败：", err)
		return
	}
}

// Consume 消费将本地视频上传到OSS的消息
func (m *VideoMQ) Consume() {
	for {
		msg, err := m.Consumer.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("[VideoMQ]从消息队列中读取消息失败:", err)
		}
		fmt.Printf("[VideoMQ]收到消息：%s\n", string(msg.Value))

		videoMsg := new(VideoMessage)
		err = json.Unmarshal(msg.Value, videoMsg)
		if err != nil {
			fmt.Println("[VideoMQ]解析消息失败:", err)
		}
		fmt.Printf("[VideoMQ]解析消息成功：%v\n", videoMsg)

		// FIXME: 下面的代码存在问题, 假如此时宕机了, 消息已经被消费了, 但是视频没有上传成功, 那么就会丢失视频
		go func() {
			//imgName := service.GetVideoCover(videoMsg.VideoName)
			//service.StoreVideoAndImg(videoMsg.VideoName, imgName, videoMsg.AuthorId, videoMsg.Title)
		}()
	}
}
