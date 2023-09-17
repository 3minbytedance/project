package main

import (
	"context"
	"douyin/dal/model"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
)


// Consume 消费点赞信息
func Consume(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		result := msgs[i].Body

		message := new(model.FavoriteAction)
		if err := json.Unmarshal(result, message); err == nil {
			flushMutex.RLock()
			mutex.Lock()
			if favoriteData[message.UserId] == nil {
				favoriteData[message.UserId] = make(map[uint]int)
			}
			switch message.ActionType {
			case 1, 2:
				favoriteData[message.UserId][message.VideoId] = message.ActionType
			}
			mutex.Unlock()
			flushMutex.RUnlock()
			continue
		}
		zap.L().Error("[FavoriteMQ]解析消息失败:", zap.Binary("favorite", result))
	}
	return consumer.ConsumeSuccess, nil
}
