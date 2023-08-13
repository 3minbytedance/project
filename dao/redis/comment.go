package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
	"strconv"
)

const ()

// GetCommentCountByVideoId 根据videoId查找评论数
func GetCommentCountByVideoId(videoId uint) (int64, error) {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	count, err := Rdb.HGet(Ctx, key, CommentCountField).Result()
	commentCount, _ := strconv.ParseInt(count, 10, 64)
	return commentCount, err
}

// IncrementCommentCountByVideoId 给videoId对应的评论数加一
func IncrementCommentCountByVideoId(videoId uint) error {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	//增加并返回评论数
	_, err := Rdb.HIncrBy(Ctx, key, CommentCountField, 1).Result()
	return err
}

// DecrementCommentCountByVideoId 给videoId对应的评论数减一
func DecrementCommentCountByVideoId(videoId uint) error {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	//减少并返回评论数
	_, err := Rdb.HIncrBy(Ctx, key, CommentCountField, -1).Result()
	return err
}

func SetCommentCountByVideoId(videoId uint, commentCount int64) error {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	err := Rdb.HSet(Ctx, key, CommentCountField, commentCount).Err()
	return err
}
