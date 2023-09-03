package redis

import (
	"douyin/dal/model"
	"github.com/cloudwego/hertz/pkg/common/json"
	red "github.com/redis/go-redis/v9"
)

func AddVideos(videos []model.Video) {
	for _, video := range videos {
		marshal, _ := json.Marshal(&video)
		Rdb.ZAdd(Ctx, VideoList, red.Z{
			Score: float64(video.CreatedAt), Member: marshal,
		})
	}
}

func AddVideo(video *model.Video) {
	marshal, _ := json.Marshal(video)
	Rdb.ZAdd(Ctx, VideoList, red.Z{
		Score: float64(video.CreatedAt), Member: marshal,
	})
}

func GetVideos(time string) []string {
	videos, _ := Rdb.ZRangeArgs(Ctx, red.ZRangeArgs{
		Key:     VideoList,
		ByScore: true,
		Rev:     true,
		Start:   0,
		Stop:    "(" + time, //(0,time)
		Offset:  0,
		Count:   29,
	}).Result()
	return videos
}
