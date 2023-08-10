package redis

import (
	"fmt"
	_ "github.com/redis/go-redis/v9"
)

const (
	KeyCommentCount = "comment_count"
)

// GetCommentCountByVideoId 根据videoId查找评论数
func GetCommentCountByVideoId(videoId uint) (string, error) {
	key := Video + fmt.Sprintf("%d", videoId)
	count, err := Rdb.HGet(Ctx, key, KeyCommentCount).Result()
	return count, err
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
//	pipeline := Rdb.Pipeline()
//	// 向videoId对应的ZSet中添加commentId
//	pipeline.ZAdd(Ctx, key, redis.Z{
//		Member: commentId,
//		Score:  float64(score),
//	})
//	// 设置过期时间
//	Rdb.Expire(Ctx, key, RdbExpireTime)
//	_, err := pipeline.Exec(Ctx)
//	return err
//}

// DeleteMappingVideoIdToCommentId 从videoId对应的zset中删除commentId
//func DeleteMappingVideoIdToCommentId(videoId, commentId int64) error {
//	// 将videoId转为string，封装成key：video_comments:12345 => [10001, 10002, 10003]
//	key := KeyVideoToComments + strconv.FormatInt(videoId, 10)
//	// 从videoId对应的zset中删除commentId
//	err := Rdb.ZRem(Ctx, key, commentId).Err()
//	return err
//}

// GetCommentIdListByVideoId 从videoId对应的zset中获取所有的commentId
//func GetCommentIdListByVideoId(videoId int64) ([]string, error) {
//	// 将videoId转为string，封装成key：video_comments:12345 => [10001, 10002, 10003]
//	key := KeyVideoToComments + strconv.FormatInt(videoId, 10)
//	commentIdStrList, err := Rdb.SMembers(Ctx, key).Result()
//	return commentIdStrList, err
//}

// GetCommentCountByVideoId 根据videoId获取对应视频的评论数
//func GetCommentCountByVideoId(videoId int64) (int64, error) {
//	key := KeyCommentCount + strconv.FormatInt(videoId, 10)
//	commentCount, err := Rdb.Get(Ctx, key).Int64()
//	return commentCount, err
//}

// AddMappingCommentIdToVideoId 添加commentId到videoId的一对一映射
//func AddMappingCommentIdToVideoId(commentId, videoId int64) error {
//	// 封装key: comment_video:10001 => 12345
//	commentIdStr := strconv.FormatInt(commentId, 10)
//	key := KeyCommentToVideo + commentIdStr
//	// 添加commentId到videoId的一对一映射
//	err := Rdb.Set(Ctx, key, videoId, RdbExpireTime).Err()
//	return err
//}

// DeleteMappingCommentIdToVideoId 删除commentId到videoId的一对一映射
//func DeleteMappingCommentIdToVideoId(commentId int64) error {
//	// 封装key: comment_video:10001 => 12345
//	key := KeyCommentToVideo + strconv.FormatInt(commentId, 10)
//	// 删除commentId到videoId的一对一映射
//	err := Rdb.Del(Ctx, key).Err()
//	return err
//}

// GetVideoIdByCommentId 根据commentId获取对应视频的videoId
//func GetVideoIdByCommentId(commentId int64) (string, error) {
//	// 封装key：comment_data:10001 => {"id": "123", "author": "user123", "timestamp": "1679921230" }
//	key := KeyCommentData + strconv.FormatInt(commentId, 10)
//	videoIdStr, err := Rdb.Get(Ctx, key).Result()
//	if err != nil {
//		return "0", err
//	}
//	return videoIdStr, err
//}

// AddCommentByCommentId 添加commentId到comment的一对一映射
//func AddCommentByCommentId(commentId int64, comment string) error {
//	// 封装key：comment_data:10001 => {"id": "123", "author": "user123", "timestamp": "1679921230" }
//	key := KeyCommentData + strconv.FormatInt(commentId, 10)
//	// 添加commentId到comment的一对一映射
//	err := Rdb.Set(Ctx, key, comment, RdbExpireTime).Err()
//	return err
//}

// DeleteCommentByCommentId 删除commentId到comment的一对一映射
//func DeleteCommentByCommentId(commentId int64) error {
//	// 封装key：comment_data:10001 => {"id": "123", "author": "user123", "timestamp": "1679921230" }
//	key := KeyCommentData + strconv.FormatInt(commentId, 10)
//	// 删除commentId到comment的一对一映射
//	err := Rdb.Del(Ctx, key).Err()
//	return err
//}

// GetCommentByCommentId 添加commentId到comment内容
//func GetCommentByCommentId(commentId int64) (string, error) {
//	// 封装key：comment_data:10001 => {"id": "123", "author": "user123", "timestamp": "1679921230" }
//	key := KeyCommentData + strconv.FormatInt(commentId, 10)
//	comment, err := Rdb.Get(Ctx, key).Result()
//	if err != nil {
//		return "", err
//	}
//	return comment, err
//}
