package redis

import "time"

const expireTime = 7 * 24 * time.Hour // 7天

// SetToken 设置token
func SetToken(token string, userId uint) error {
	// token字段作为key, userId作为value
	key := TokenKey + token
	err := Rdb.Set(Ctx, key, userId, expireTime).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetToken 获取token
func GetToken(token string) uint {
	key := TokenKey + token
	// 判断key是否存在
	exists, err := Rdb.Exists(Ctx, key).Result()
	if exists == 0 {
		return 0
	}
	// 获取key对应的value
	userId, err := Rdb.Get(Ctx, key).Uint64()
	if err != nil {
		return 0
	}
	// 重置过期时间
	err = Rdb.Expire(Ctx, key, expireTime).Err()
	if err != nil {
		return 0
	}
	return uint(userId)
}
