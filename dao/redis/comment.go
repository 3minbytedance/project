package redis

import (
	"github.com/redis/go-redis/v9"
	"strconv"
)

// AddMappingVideoIdToCommentId 向videoId对应的zset中添加commentId
func AddMappingVideoIdToCommentId(videoId, commentId int64, score int64) error {
	// 将videoId转为string
	videoIdStr := strconv.FormatInt(videoId, 10)
	// 向videoId对应的zset中添加commentId
	err := RdbVCId.ZAdd(RedisCtx, videoIdStr, redis.Z{
		Score:  float64(score),
		Member: commentId,
	}).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	err = RdbVCId.Expire(RedisCtx, videoIdStr, RdbExpireTime).Err()
	return err
}

// DeleteMappingVideoIdToCommentId 从videoId对应的zset中删除commentId
func DeleteMappingVideoIdToCommentId(videoId, commentId int64) error {
	// 将videoId转为string
	videoIdStr := strconv.FormatInt(videoId, 10)
	// 从videoId对应的zset中删除commentId
	err := RdbVCId.ZRem(RedisCtx, videoIdStr, commentId).Err()
	return err
}

// GetCommentIdListByVideoId 从videoId对应的zset中获取所有的commentId
func GetCommentIdListByVideoId(videoId int64) ([]string, error) {
	// 将videoId转为string
	videoIdStr := strconv.FormatInt(videoId, 10)
	commentIdStrList, err := RdbVCId.SMembers(RedisCtx, videoIdStr).Result()
	return commentIdStrList, err
}

// GetCommentCountByViedoId 根据videoId获取对应视频的评论数
func GetCommentCountByViedoId(videoId int64) (int64, error) {
	// 将videoId转为string
	videoIdStr := strconv.FormatInt(videoId, 10)
	count, err := RdbVCId.ZCard(RedisCtx, videoIdStr).Result()
	return count, err
}

// AddMappingCommentIdToVideoId 添加commentId到videoId的一对一映射
func AddMappingCommentIdToVideoId(commentId, videoId int64) error {
	// 将commentId转为string
	commentIdStr := strconv.FormatInt(commentId, 10)
	// 将videoId转为string
	videoIdStr := strconv.FormatInt(videoId, 10)
	// 添加commentId到videoId的一对一映射
	err := RdbCVId.Set(RedisCtx, commentIdStr, videoIdStr, RdbExpireTime).Err()
	return err
}

// DeleteMappingCommentIdToVideoId 删除commentId到videoId的一对一映射
func DeleteMappingCommentIdToVideoId(commentId int64) error {
	// 将commentId转为string
	commentIdStr := strconv.FormatInt(commentId, 10)
	// 删除commentId到videoId的一对一映射
	err := RdbCVId.Del(RedisCtx, commentIdStr).Err()
	return err
}

// GetVideoIdByCommentId 根据commentId获取对应视频的videoId
func GetVideoIdByCommentId(commentId int64) (string, error) {
	// 将commentId转为string
	commentIdStr := strconv.FormatInt(commentId, 10)
	videoIdStr, err := RdbCVId.Get(RedisCtx, commentIdStr).Result()
	if err != nil {
		return "0", err
	}
	return videoIdStr, err
}

// AddCommentByCommentId 添加commentId到comment的一对一映射
func AddCommentByCommentId(commentId int64, comment string) error {
	// 将commentId转为string
	commentIdStr := strconv.FormatInt(commentId, 10)
	// 添加commentId到comment的一对一映射
	err := RdbCIdComment.Set(RedisCtx, commentIdStr, comment, RdbExpireTime).Err()
	return err
}

// DeleteCommentByCommentId 删除commentId到comment的一对一映射
func DeleteCommentByCommentId(commentId int64) error {
	// 将commentId转为string
	commentIdStr := strconv.FormatInt(commentId, 10)
	// 删除commentId到comment的一对一映射
	err := RdbCIdComment.Del(RedisCtx, commentIdStr).Err()
	return err
}

// GetCommentByCommentId 添加commentId到comment内容
func GetCommentByCommentId(commentId int64) (string, error) {
	// 将commentId转为string
	commentIdStr := strconv.FormatInt(commentId, 10)
	comment, err := RdbCIdComment.Get(RedisCtx, commentIdStr).Result()
	if err != nil {
		return "", err
	}
	return comment, err
}
