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

func GetVideos(time string) []model.Video {
	videos, _ := Rdb.ZRevRangeByScore(Ctx, "videos",
		&redisClient.ZRangeBy{
			Min:    "0",
			Max:    time, // 根据需要的时间格式进行转换
			Offset: 0,
			Count:  30,
		}).Result()
	video := make([]model.Video, 0, 30)
	var v model.Video
	fmt.Println("redis get video")
	for _, val := range videos {
		fmt.Println(val)
		err := json.Unmarshal([]byte(val), &v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		video = append(video, v)
	}
	return video
}
