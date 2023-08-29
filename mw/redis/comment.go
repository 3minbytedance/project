package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

// GetCommentCountByVideoId 根据videoId查找评论数
func GetCommentCountByVideoId(videoId uint) (int64, error) {
	key := fmt.Sprintf("%s:%d", VideoKey, videoId)
	count, err := Rdb.HGet(Ctx, key, CommentCountField).Result()
	commentCount, _ := strconv.ParseInt(count, 10, 64)
	return commentCount, err
}

// IncrementCommentCountByVideoId 给videoId对应的评论数加一
func IncrementCommentCountByVideoId(videoId uint) error {
	key := fmt.Sprintf("%s:%d", VideoKey, videoId)
	//增加并返回评论数
	_, err := Rdb.HIncrBy(Ctx, key, CommentCountField, 1).Result()
	return err
}

// DecrementCommentCountByVideoId 给videoId对应的评论数减一
func DecrementCommentCountByVideoId(videoId uint) error {
	key := fmt.Sprintf("%s:%d", VideoKey, videoId)
	//减少并返回评论数
	_, err := Rdb.HIncrBy(Ctx, key, CommentCountField, -1).Result()
	return err
}

func SetCommentCountByVideoId(videoId uint, commentCount int64) error {
	key := fmt.Sprintf("%s:%d", VideoKey, videoId)
	err := Rdb.HSet(Ctx, key, CommentCountField, commentCount).Err()
	randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}
