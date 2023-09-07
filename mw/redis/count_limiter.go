package redis

import (
	"strconv"
	"strings"
	"time"
)

const (
	limiterTime = 2 * time.Hour // 2 小时

	loginLimit    = "loginLimit"
	registerLimit = "registerLimit"
	commentLimit  = "commentLimit"
	uploadLimit   = "uploadLimit"
	messageLimit  = "messageLimit"
)

const (
	loginMaxCount    = 5
	registerMaxCount = 5
	commentMaxCount  = 30
	uploadMaxCount   = 3
	messageMaxCount  = 50
)

func IncrementLoginLimiterCount(ip string) bool {
	baseSlice := []string{loginLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result := Rdb.Get(Ctx, key).Val()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < loginMaxCount {
		Rdb.Incr(Ctx, key)
		return true
	}
	return false
}

func IncrementRegisterLimiterCount(ip string) bool {
	baseSlice := []string{registerLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result := Rdb.Get(Ctx, key).Val()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < registerMaxCount {
		Rdb.Incr(Ctx, key)
		return true
	}
	return false
}

func IncrementCommentLimiterCount(ip string) bool {
	baseSlice := []string{commentLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result := Rdb.Get(Ctx, key).Val()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < commentMaxCount {
		Rdb.Incr(Ctx, key)
		return true
	}
	return false
}

func IncrementUploadLimiterCount(ip string) bool {
	baseSlice := []string{uploadLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result := Rdb.Get(Ctx, key).Val()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < uploadMaxCount {
		Rdb.Incr(Ctx, key)
		return true
	}
	return false
}

func IncrementMessageLimiterCount(ip string) bool {
	baseSlice := []string{messageLimit, ip}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.SetNX(Ctx, key, 0, limiterTime)
	result := Rdb.Get(Ctx, key).Val()
	count, err := strconv.Atoi(result)
	if err != nil {
		return false
	}
	if count < messageMaxCount {
		Rdb.Incr(Ctx, key)
		return true
	}
	return false
}
