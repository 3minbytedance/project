package mongo

import (
	"context"
	"douyin/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

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
