package rocketMQ

import (
	"context"
	"douyin/dal/model"
	"douyin/dal/mysql"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
)

type CommentMQ struct {
	MQ
}

var (
	CommentMQInstance *CommentMQ
)

func InitCommentKafka() {
	CommentMQInstance = &CommentMQ{
		MQ{
			Topic:   "comments",
			GroupId: "comment_group",
		},
	}

	// 创建 Comment 业务的生产者和消费者实例
	CommentMQInstance.Producer = rocketMQManager.NewProducer(CommentMQInstance.GroupId)
	err := CommentMQInstance.Producer.Start()
	if err != nil {
		panic("启动comment producer 失败")
	}

	CommentMQInstance.Consumer = rocketMQManager.NewConsumer(CommentMQInstance.Topic, CommentMQInstance.GroupId)

	err = CommentMQInstance.Consumer.Subscribe(CommentMQInstance.Topic, consumer.MessageSelector{}, Consum)
	if err != nil {
		panic("comment consumer 订阅失败")
	}

	err = CommentMQInstance.Consumer.Start()
	if err != nil {
		panic("启动comment consumer 失败")
	}
}

// Consume 消费添加或者删除评论的消息
func Consum(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		fmt.Printf("subscribe callback : %v \n", msgs[i])
		msg := msgs[i].Body

		var result json.RawMessage
		err := json.Unmarshal(msg, &result)
		if err != nil {
			log.Println("[CommentMQ]解析消息失败:", err)
			continue
		}

		// 解析消息, 消息类型可能为model.Comment, 也可能为CommentId, 如果是前者, 则添加评论, 如果是后者, 则删除评论
		// 解析为model.Comment, 则向数据库中添加评论
		message := new(model.Comment)
		err = json.Unmarshal(result, message)
		if err == nil {
			_, err = mysql.AddComment(message)
			if err != nil {
				log.Println("[CommentMQ]向mysql中添加评论失败:", err)
			}
			log.Println("[CommentMQ]向mysql中添加评论成功")
			continue
		}

		// 解析为整型, 即CommentId, 则从数据库中删除评论
		commentId := new(uint)
		err = json.Unmarshal(result, commentId)
		if err == nil {
			err = mysql.DeleteCommentById(*commentId)
			if err != nil {
				log.Println("[CommentMQ]从mysql中删除评论失败:", err)
			}
			log.Println("[CommentMQ]从mysql中删除评论成功")
			continue
		}
		log.Println("[CommentMQ]解析消息失败:", result)

	}
	return consumer.ConsumeSuccess, nil
}

// ProduceAddCommentMsg 发布添加评论的消息, 向mysql中添加评论时, 调用此方法
func (m *CommentMQ) ProduceAddCommentMsg(message *model.Comment) error {
	_, err := rocketMQManager.ProduceMessage(m.Producer, message, m.Topic)

	if err != nil {
		log.Println("rocketMQ 发送添加评论的消息失败：", err)
		return err
	}

	return nil
}

// ProduceDelCommentMsg 发布删除评论的消息, mysql删除评论时, 调用此方法
func (m *CommentMQ) ProduceDelCommentMsg(commentId uint) error {
	_, err := rocketMQManager.ProduceMessage(m.Producer, commentId, m.Topic)
	if err != nil {
		log.Println("rocketMQ 发送删除评论的消息失败：", err)
		return err
	}
	return nil
}
