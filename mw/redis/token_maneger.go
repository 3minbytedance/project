package redis

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const expireTime = 7 * 24 * time.Hour // 7天

// SetToken 设置token
func SetToken(userId uint, token string) {
	// token字段作为key, userId作为value
	key := TokenKey + strconv.Itoa(int(userId))
	err := Rdb.Set(Ctx, key, token, expireTime).Err()
	if err != nil {
		zap.L().Error("SetToken failed", zap.Error(err))
	}
}

// TokenIsExisted 判断用户对应的token是否存在
func TokenIsExisted(userId uint) bool {
	key := TokenKey + strconv.Itoa(int(userId))
	// 判断key是否存在
	exists, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return err == nil && exists == 1
}
