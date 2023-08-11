package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"log"
	"project/models"
)

func SendMessage(message *models.Message) (err error) {
	collection := Mongo.Collection("messages")
	_, err = collection.InsertOne(Ctx, message)
	if err != nil {
		fmt.Println("消息插入到 MongoDB失败。")
		return err
	}
	fmt.Println("消息已插入到 MongoDB。")
	return
}

func GetMessageList(fromUserId, toUserId uint, preMsgTime int64) ([]models.Message, error) {
	collection := Mongo.Collection("messages")
	filter := bson.M{
		"$and": []bson.M{
			{"$or": []bson.M{
				{"fromuserid": fromUserId, "touserid": toUserId},
				{"fromuserid": toUserId, "touserid": fromUserId},
			}},
			{"createtime": bson.M{"$gt": preMsgTime}}, // 添加时间戳条件
		},
	}

	cursor, err := collection.Find(Ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(Ctx)

	var messages []models.Message
	for cursor.Next(Ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			log.Println("解码错误:", err)
			continue
		}
		messages = append(messages, message)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("找到的聊天记录数量:", len(messages))

	return messages, nil
}
