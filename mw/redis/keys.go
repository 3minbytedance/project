package redis

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

const Delimiter = ":"

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

const VideoList = "videos"

const TokenKey = "token"

const (
	Lock               = "lock"
	FavoriteAction     = "favoriteAction"
	FollowAction       = "followAction"
	RetryTime          = 30 * time.Millisecond
	KeyExistsAndNotSet = 0
	KeyUpdated         = 1
	KeyNotExistsInBoth = 2
)

func IsExistUserField(userId uint, field string) bool {
	baseSlice := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	exists, err := Rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		zap.L().Error("redis isExistVideo 连接失败")
		return false
	}
	return exists
}

func IsExistVideoField(videoId uint, field string) bool {
	baseSlice := []string{VideoKey, strconv.Itoa(int(videoId))}
	key := strings.Join(baseSlice, Delimiter)
	exists, err := Rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		zap.L().Error("redis isExistVideo 连接失败")
		return false
	}
	return exists
}

// IsExistUserSetField 判断set类型的是否存在
func IsExistUserSetField(userId uint, field string) bool {
	baseSlice := []string{field, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	exists, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		zap.L().Error("redis isExistVideo 连接失败")
		return false
	}
	return exists != 0
}

// DelVideoHashField 根据参数删除video Hash field
func DelVideoHashField(videoId uint, field string) {
	baseSlice := []string{VideoKey, strconv.Itoa(int(videoId))}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.HDel(Ctx, key, field)
}

// DelUserHashField 根据参数删除user Hash field
func DelUserHashField(userId uint, field string) {
	baseSlice := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.HDel(Ctx, key, field)
}

// DelVideoKey 删除videos ZSet key
func DelVideoKey() {
	key := VideoList
	Rdb.Del(Ctx, key)
}

func AcquireCommentLock(videoId uint) bool {
	baseSlice := []string{CommentCountField, strconv.Itoa(int(videoId)), Lock}
	key := strings.Join(baseSlice, Delimiter)
	result, err := Rdb.SetNX(Ctx, key, 1, 2*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseCommentLock(videoId uint) {
	baseSlice := []string{CommentCountField, strconv.Itoa(int(videoId)), Lock}
	key := strings.Join(baseSlice, Delimiter)
	Rdb.Del(Ctx, key)
}

func AcquireUserLock(userId uint, field string) bool {
	var key string
	var baseSlice []string
	switch field {
	case NameField:
		baseSlice = []string{NameField, strconv.Itoa(int(userId)), Lock}
	case WorkCountField:
		baseSlice = []string{WorkCountField, strconv.Itoa(int(userId)), Lock}
	default:
		return false
	}
	key = strings.Join(baseSlice, Delimiter)
	result, err := Rdb.SetNX(Ctx, key, 1, 2*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseUserLock(userId uint, field string) {
	var key string
	var baseSlice []string
	switch field {
	case NameField:
		baseSlice = []string{NameField, strconv.Itoa(int(userId)), Lock}
	case WorkCountField:
		baseSlice = []string{WorkCountField, strconv.Itoa(int(userId)), Lock}
	default:
		return
	}
	key = strings.Join(baseSlice, Delimiter)
	Rdb.Del(Ctx, key)
}

func AcquireFavoriteLock(id uint, field string) bool {
	var key string
	var baseSlice []string
	switch field {
	case FavoriteList:
		baseSlice = []string{FavoriteList, strconv.Itoa(int(id)), Lock}
	case VideoFavoritedCountField:
		baseSlice = []string{VideoFavoritedCountField, strconv.Itoa(int(id)), Lock}
	case TotalFavoriteField:
		baseSlice = []string{TotalFavoriteField, strconv.Itoa(int(id)), Lock}
	case FavoriteAction:
		baseSlice = []string{FavoriteAction, strconv.Itoa(int(id)), Lock}
	default:
		return false
	}
	key = strings.Join(baseSlice, Delimiter)
	result, err := Rdb.SetNX(Ctx, key, 1, 2*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseFavoriteLock(id uint, field string) {
	var key string
	var baseSlice []string
	switch field {
	case FavoriteList:
		baseSlice = []string{FavoriteList, strconv.Itoa(int(id)), Lock}
	case VideoFavoritedCountField:
		baseSlice = []string{VideoFavoritedCountField, strconv.Itoa(int(id)), Lock}
	case TotalFavoriteField:
		baseSlice = []string{TotalFavoriteField, strconv.Itoa(int(id)), Lock}
	case FavoriteAction:
		baseSlice = []string{FavoriteAction, strconv.Itoa(int(id)), Lock}
	default:
		return
	}
	key = strings.Join(baseSlice, Delimiter)
	Rdb.Del(Ctx, key)
}

func AcquireRelationLock(id uint, field string) bool {
	var key string
	var baseSlice []string
	switch field {
	case FollowList:
		baseSlice = []string{FollowList, strconv.Itoa(int(id)), Lock}
	case FollowerList:
		baseSlice = []string{FollowerList, strconv.Itoa(int(id)), Lock}
	case FollowAction:
		baseSlice = []string{FollowAction, strconv.Itoa(int(id)), Lock}
	default:
		return false
	}
	key = strings.Join(baseSlice, Delimiter)
	result, err := Rdb.SetNX(Ctx, key, 1, 2*time.Second).Result()
	if err != nil {
		zap.L().Error("获取锁失败", zap.Error(err))
		return false
	}
	return result
}

func ReleaseRelationLock(id uint, field string) {
	var key string
	var baseSlice []string
	switch field {
	case FollowList:
		baseSlice = []string{FollowList, strconv.Itoa(int(id)), Lock}
	case FollowerList:
		baseSlice = []string{FollowerList, strconv.Itoa(int(id)), Lock}
	case FollowAction:
		baseSlice = []string{FollowAction, strconv.Itoa(int(id)), Lock}
	default:
		return
	}
	key = strings.Join(baseSlice, Delimiter)
	Rdb.Del(Ctx, key)
}
