package kafka

import (
	"context"
	"douyin/dal/mysql"
	"douyin/mw/redis"
	"douyin/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
)

type VideoMessage struct {
	VideoPath     string
	VideoFileName string
	UserID        uint
	Title         string
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
		videoMsg := new(VideoMessage)
		err = json.Unmarshal(msg.Value, videoMsg)
		if err != nil {
			fmt.Println("[VideoMQ]解析消息失败:", err)
		}
		go func() {
			zap.L().Info("开始处理视频消息", zap.Any("videoMsg", videoMsg))
			//视频存储到oss
			if err = utils.UploadToOSS(videoMsg.VideoPath, videoMsg.VideoFileName); err != nil {
				zap.L().Error("上传视频到OSS失败", zap.Error(err))
			}

			//利用oss功能获取封面图
			imgName, err := utils.GetVideoCover(videoMsg.VideoFileName)
			if err != nil {
				zap.L().Error("图片截帧失败", zap.Error(err))
			}

			// 视频信息存储到MySQL
			mysql.InsertVideo(videoMsg.VideoFileName, imgName, videoMsg.UserID, videoMsg.Title)

			// 更新redis中的用户作品数
			if !redis.IsExistUserField(videoMsg.UserID, redis.WorkCountField) {
				cnt := mysql.FindWorkCountsByAuthorId(videoMsg.UserID)
				err := redis.SetWorkCountByUserId(videoMsg.UserID, cnt)
				if err != nil {
					zap.L().Error("redis更新作品数失败", zap.Error(err))
					return
				}
			}
			err = redis.IncrementWorkCountByUserId(videoMsg.UserID)
			if err != nil {
				zap.L().Error("redis增加其作品数失败", zap.Error(err))
				return
			}
			zap.L().Info("视频消息处理成功", zap.Any("videoMsg", videoMsg))
		}()
	}
}
