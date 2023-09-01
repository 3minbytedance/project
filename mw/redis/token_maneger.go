package redis

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

const expireTime = 7 * 24 * time.Hour // 7天

// SetToken 设置token
func SetToken(userId uint, token string) {
	// userId作为key
	baseSlice := []string{TokenKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	err := Rdb.Set(Ctx, key, token, expireTime).Err()
	if err != nil {
		zap.L().Error("SetToken failed", zap.Error(err))
	}
}

// TokenIsExisted 判断用户对应的token是否存在
func TokenIsExisted(userId uint) bool {
	baseSlice := []string{TokenKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	// 判断key是否存在
	exists, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		return false
	}
	return err == nil && exists == 1
}
