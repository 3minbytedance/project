package redis

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

const (
	//VideoKey hash 类型 key:video_videoId
	VideoKey = "video"

	CommentCountField = "commentCount"

	VideoFavoritedCountField = "favoritedCount" // 视频被点赞总数
)

const (
	//UserKey hash 类型 key:user_userId
	UserKey = "user"

	WorkCountField     = "workCount"      //作品数
	NameField          = "name"           //用户名
	TotalFavoriteField = "totalFavorited" //发布视频的总获赞数量

	// FavoriteList  set类型
	FavoriteList = "favoriteList" //喜欢视频列表

	// FollowList and FollowerList  set类型
	FollowList   = "followList"   //关注列表
	FollowerList = "followerList" //粉丝列表
)

const TokenKey = "token:"

const (
	Lock               = "lock"
	RetryTime          = 30 * time.Millisecond
	KeyExistsAndNotSet = 0
	KeyUpdated         = 1
	KeyNotExistsInBoth = 2
)

func IsExistUserField(userId uint, field string) bool {
	key := fmt.Sprintf("%s:%d", UserKey, userId)
	exists, err := Rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		zap.L().Error("redis isExistVideo 连接失败")
		return false
	}
	if exists {
		randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
		expiration := time.Duration(randomSeconds) * time.Second
		Rdb.Expire(Ctx, key, expiration)
	}
	return exists
}

func IsExistVideoField(videoId uint, field string) bool {
	key := fmt.Sprintf("%s:%d", VideoKey, videoId)
	exists, err := Rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		zap.L().Error("redis isExistVideo 连接失败")
		return false
	}
	if exists {
		randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
		expiration := time.Duration(randomSeconds) * time.Second
		Rdb.Expire(Ctx, key, expiration)
	}
	return exists
}

// IsExistUserSetField 判断set类型的是否存在
func IsExistUserSetField(userId uint, field string) bool {
	key := fmt.Sprintf("%s:%d", field, userId)
	exists, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		zap.L().Error("redis isExistVideo 连接失败")
		return false
	}
	if exists != 0 {
		randomSeconds := rand.Intn(600) + 30 // 600秒到630秒之间的随机数
		expiration := time.Duration(randomSeconds) * time.Second
		Rdb.Expire(Ctx, key, expiration)
	}
	return exists != 0
}

// DelKey 根据参数合成并删除key
func DelKey(userId uint, field string) {
	key := fmt.Sprintf("%s:%d", field, userId)
	Rdb.Del(Ctx, key)
}

func AcquireCommentLock(videoId uint) bool {
	key := fmt.Sprintf("%s:%d_%s", CommentCountField, videoId, Lock)
	result, err := Rdb.SetNX(Ctx, key, 1, 30*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseCommentLock(videoId uint) {
	key := fmt.Sprintf("%s:%d_%s", CommentCountField, videoId, Lock)
	Rdb.Del(Ctx, key)
}

func AcquireUserLock(userId uint, field string) bool {
	var key string
	switch field {
	case NameField:
		key = fmt.Sprintf("%s:%d_%s", NameField, userId, Lock)
	case WorkCountField:
		key = fmt.Sprintf("%s:%d_%s", WorkCountField, userId, Lock)
	default:
		return false
	}
	result, err := Rdb.SetNX(Ctx, key, 1, 3*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseUserLock(userId uint, field string) {
	var key string
	switch field {
	case NameField:
		key = fmt.Sprintf("%s:%d_%s", NameField, userId, Lock)
	case WorkCountField:
		key = fmt.Sprintf("%s:%d_%s", WorkCountField, userId, Lock)
	default:
		return
	}
	Rdb.Del(Ctx, key)
}

func AcquireFavoriteLock(id uint, field string) bool {
	var key string
	switch field {
	case FavoriteList:
		key = fmt.Sprintf("%s:%d_%s", FavoriteList, id, Lock)
	case VideoFavoritedCountField:
		key = fmt.Sprintf("%s:%d_%s", VideoFavoritedCountField, id, Lock)
	case TotalFavoriteField:
		key = fmt.Sprintf("%s:%d_%s", TotalFavoriteField, id, Lock)
	default:
		return false
	}
	result, err := Rdb.SetNX(Ctx, key, 1, 2*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseFavoriteLock(id uint, field string) {
	var key string
	switch field {
	case FavoriteList:
		key = fmt.Sprintf("%s:%d_%s", FavoriteList, id, Lock)
	case VideoFavoritedCountField:
		key = fmt.Sprintf("%s:%d_%s", VideoFavoritedCountField, id, Lock)
	case TotalFavoriteField:
		key = fmt.Sprintf("%s:%d_%s", TotalFavoriteField, id, Lock)
	default:
		return
	}
	Rdb.Del(Ctx, key)
}

func AcquireRelationLock(id uint, field string) bool {
	var key string
	switch field {
	case FollowList:
		key = fmt.Sprintf("%s:%d_%s", FollowList, id, Lock)
	case FollowerList:
		key = fmt.Sprintf("%s:%d_%s", FollowerList, id, Lock)
	default:
		return false
	}
	result, err := Rdb.SetNX(Ctx, key, 1, 2*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseRelationLock(id uint, field string) {
	var key string
	switch field {
	case FollowList:
		key = fmt.Sprintf("%s:%d_%s", FollowList, id, Lock)
	case FollowerList:
		key = fmt.Sprintf("%s:%d_%s", FollowerList, id, Lock)
	default:
		return
	}
	Rdb.Del(Ctx, key)
}
