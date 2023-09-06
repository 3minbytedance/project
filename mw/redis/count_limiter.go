package redis

import (
	"strings"
	"time"
)

const limiterTime = 2 * time.Hour // 2 小时
const loginLimit = "loginLimit"
const commentLimit = "commentLimit"
const uploadLimit = "uploadLimit"

func IncrementLoginLimiterComment(ip string) int64 {
	baseSlice := []string{loginLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	val := Rdb.Incr(Ctx, ip).Val()
	return val
}

func IncrementCommentLimiterComment(ip string) int64 {
	baseSlice := []string{commentLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	val := Rdb.Incr(Ctx, ip).Val()
	return val
}

func IncrementUploadLimiterComment(ip string) int64 {
	baseSlice := []string{uploadLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	val := Rdb.Incr(Ctx, ip).Val()
	return val
}
