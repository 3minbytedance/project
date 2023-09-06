package redis

import (
	"strconv"
	"strings"
	"time"
)

const limiterTime = 2 * time.Hour // 2 小时
const loginLimit = "loginLimit"
const commentLimit = "commentLimit"
const uploadLimit = "uploadLimit"

func IncrementLoginLimiterComment(ip string) bool {
	baseSlice := []string{loginLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result, _ := Rdb.Get(Ctx, key).Result()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < 5 {
		Rdb.Incr(Ctx, key).Val()
		return true
	}
	return false
}

func IncrementCommentLimiterComment(ip string) bool {
	baseSlice := []string{commentLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result, _ := Rdb.Get(Ctx, key).Result()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < 20 {
		Rdb.Incr(Ctx, key).Val()
		return true
	}
	return false
}

func IncrementUploadLimiterComment(ip string) bool {
	baseSlice := []string{uploadLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result, _ := Rdb.Get(Ctx, key).Result()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < 3 {
		Rdb.Incr(Ctx, key).Val()
		return true
	}
	return false
}
