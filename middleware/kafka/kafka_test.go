package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"project/models"
	"testing"
	"time"
)

func TestKafka(t *testing.T) {
	InitMessageKafka()
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoUrl := fmt.Sprintf("mongodb://%s:%d", "127.0.0.1", 27017)
	_, _ = mongo.Connect(Ctx, options.Client().ApplyURI(mongoUrl))
	//.SetAuth(options.Credential{
	//		Username: conf.Username,
	//		Password: conf.Password,
	//	}
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
		//err = mongo.SendMessage(message)
		//if err != nil {
		//	log.Println("kafka往mongo存入message失败：", err)
		//	return
		//}
	}

}
