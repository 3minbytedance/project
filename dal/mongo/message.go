package mongo

import (
	"douyin/dal/model"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"log"
)

func SendMessage(data []byte) (err error) {
	message := new(model.Message)
	err = json.Unmarshal(data, message)
	if err != nil {
		log.Println("kafka获取message反序列化失败：", err)
	}
	collection := Mongo.Collection("messages")
	_, err = collection.InsertOne(Ctx, message)
	if err != nil {
		fmt.Println("消息插入到 MongoDB失败。")
		return err
	}
	fmt.Println("消息已插入到 MongoDB。")
	return
}

func GetMessageList(fromUserId, toUserId uint, preMsgTime int64) ([]*model.Message, error) {
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

	var messages []*model.Message
	for cursor.Next(Ctx) {
		var message model.Message
		if err := cursor.Decode(&message); err != nil {
			log.Println("解码错误:", err)
			continue
		}
		messages = append(messages, &message)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("找到的聊天记录数量:", len(messages))

	return messages, nil
}
