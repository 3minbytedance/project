package redis

import (
	"douyin/dal/model"
	"fmt"
	redisClient "github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func AddVideo(video model.Video) error {
	err := Rdb.ZAdd(Ctx, "videos", redisClient.Z{
		Score: float64(video.CreatedAt), Member: video.VideoUrl + video.CoverUrl,
	}).Err()
	return err
}

func GetVideos() {
	videos, _ := Rdb.ZRevRangeByScoreWithScores(Ctx, "videos",
		&redisClient.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.FormatInt(time.Now().Unix(), 10), // 根据需要的时间格式进行转换
			Offset: 0,
			Count:  30,
		}).Result()

	for _, val := range videos {
		fmt.Println(val.Score)
		fmt.Println(val.Member)
	}
}
