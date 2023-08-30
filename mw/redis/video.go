package redis

import (
	"douyin/dal/model"
	"encoding/json"
	"fmt"
	redisClient "github.com/redis/go-redis/v9"
)

func AddVideo(video model.Video) {
	marshal, _ := json.Marshal(&video)
	Rdb.ZAdd(Ctx, "videos", redisClient.Z{
		Score: float64(video.CreatedAt), Member: marshal,
	})
}

func GetVideos(time string) {
	videos, _ := Rdb.ZRevRangeByScoreWithScores(Ctx, "videos",
		&redisClient.ZRangeBy{
			Min:    "0",
			Max:    time, // 根据需要的时间格式进行转换
			Offset: 0,
			Count:  30,
		}).Result()

	for _, val := range videos {
		fmt.Println(val.Score)
		fmt.Println(val.Member)
	}
}
