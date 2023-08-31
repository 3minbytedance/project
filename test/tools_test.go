package test

import (
	"context"
	"douyin/common"
	"douyin/config"
	"douyin/dal/model"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func TestBloom(t *testing.T) {
	common.InitUserBloomFilter()
	common.AddToUserBloom("user1")
	common.AddToUserBloom("user2")
	common.AddToUserBloom("user3")
	common.AddToUserBloom("user4")
	common.AddToUserBloom("use1")
	common.AddToUserBloom("use2")

	assert.True(t, common.TestUserBloom("usera1"))
	assert.False(t, common.TestUserBloom("user5"))
}

func TestSensitiveFilter(t *testing.T) {
	err := common.InitSensitiveFilter()
	if err != nil {
		t.Log(err)
		return
	}
	word := common.ReplaceWord("傻逼吧卧槽啊你妈的")
	assert.True(t, word == "sad")

}

var Mongo *mongo.Database
var Ctx context.Context

func Init(appConfig *config.AppConfig) (err error) {
	var conf *config.MongoConfig
	if appConfig.Mode == config.LocalMode {
		conf = appConfig.Local.MongoConfig
	} else {
		conf = appConfig.Remote.MongoConfig
	}
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoUrl := fmt.Sprintf("mongodb://%s:%d", conf.Address, conf.Port)
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(mongoUrl))
	//.SetAuth(options.Credential{
	//		Username: conf.Username,
	//		Password: conf.Password,
	//	}
	if err != nil {
		log.Println("Connection MongoDB Error:", err)
		return
	}
	Mongo = client.Database("3minbytedance")
	return
}

func SendMessage(message *model.Message) (err error) {
	collection := Mongo.Collection("messages")
	_, err = collection.InsertOne(Ctx, message)
	if err != nil {
		log.Println("消息插入到 MongoDB失败。")
		return err
	}
	return
}

func GetMessageList(fromUserId, toUserId uint, preMsgTime int64) ([]*model.Message, error) {
	collection := Mongo.Collection("messages")
	filter := bson.M{
		"$and": []bson.M{
			{"$or": []bson.M{
				{"from_userid": fromUserId, "to_userid": toUserId},
				{"from_userid": toUserId, "to_userid": fromUserId},
			}},
			{"create_time": bson.M{"$gt": preMsgTime}}, // 添加时间戳条件
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

	//log.Println("找到的聊天记录数量:", len(messages))

	return messages, nil
}

func TestMongo(t *testing.T) {
	// 加载配置
	if err := config.Init(); err != nil {
		zap.L().Error("Load config failed, err:%v\n", zap.Error(err))
		return
	}
	err := Init(config.Conf)
	if err != nil {
		log.Println(err)
	}
}
