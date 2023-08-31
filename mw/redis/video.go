package redis

import (
	"douyin/dal/model"
	"encoding/json"
	red "github.com/redis/go-redis/v9"
)

func AddVideos(videos []model.Video) {
	for _, video := range videos {
		marshal, _ := json.Marshal(&video)
		Rdb.ZAdd(Ctx, "videos", red.Z{
			Score: float64(video.CreatedAt), Member: marshal,
		})
	}
}

func AddVideo(video model.Video) {
	marshal, _ := json.Marshal(&video)
	Rdb.ZAdd(Ctx, "videos", red.Z{
		Score: float64(video.CreatedAt), Member: marshal,
	})
}

func GetVideos(time string) []model.Video {
	videos, _ := Rdb.ZRangeArgs(Ctx, red.ZRangeArgs{
		Key:     VideoList,
		ByScore: true,
		Rev:     true,
		Start:   0,
		Stop:    "(" + time, //(0,time)
		Offset:  0,
		Count:   30,
	}).Result()

	videoList := make([]model.Video, 0, 30)
	var v model.Video
	for _, val := range videos {
		err := json.Unmarshal([]byte(val), &v)
		if err != nil {
			continue
		}
		videoList = append(videoList, v)
	}
	return videoList
}
