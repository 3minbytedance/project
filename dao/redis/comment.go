package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
	"strconv"
)

const (
	KeyCommentCount = "comment_count"
)

// GetCommentCountByVideoId 根据videoId查找评论数
func GetCommentCountByVideoId(videoId uint) (int, error) {
	key := Video + fmt.Sprintf("%d", videoId)
	count, err := Rdb.HGet(Ctx, key, KeyCommentCount).Result()
	commentCount, _ := strconv.Atoi(count)
	return commentCount, err
}

// IncrementCommentCountByVideoId 给videoId对应的评论数加一
func IncrementCommentCountByVideoId(videoId uint) error {
	key := Video + fmt.Sprintf("%d", videoId)
	_, err := Rdb.HIncrBy(Ctx, key, KeyCommentCount, 1).Result()
	//TODO 增加并返回评论数
	return err
}

// DecrementCommentCountByVideoId 给videoId对应的评论数减一
func DecrementCommentCountByVideoId(videoId uint) error {
	key := Video + fmt.Sprintf("%d", videoId)
	_, err := Rdb.HIncrBy(Ctx, key, KeyCommentCount, -1).Result()
	//TODO 减少并返回评论数
	return err
}

func SetCommentCountByVideoId(videoId uint, commentCount int64) error {
	key := Video + fmt.Sprintf("%d", videoId)
	err := Rdb.HSet(Ctx, key, KeyCommentCount, commentCount).Err()
	return err
}

// AddMappingVideoIdToCommentId 向videoId对应的zset中添加commentId
//func AddMappingVideoIdToCommentId(videoId, commentId int64, score int64) error {
//	// 将videoId转为string，封装成key：video_comments:12345 => [10001, 10002, 10003]
//	videoIdStr := strconv.FormatInt(videoId, 10)
//	key := KeyVideoToComments + videoIdStr
//	// 使用pipeline一次发送多条命令减少rtt
//	pipeline := RdbComment.Pipeline()
//	// 向videoId对应的ZSet中添加commentId
//	pipeline.ZAdd(Ctx, key, redis.Z{
//		Member: commentId,
//		Score:  float64(score),
//	})
//	// 设置过期时间
//	RdbComment.Expire(Ctx, key, RdbExpireTime)
//	_, err := pipeline.Exec(Ctx)
//	return err
//}
