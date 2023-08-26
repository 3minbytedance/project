package redis

import (
	"fmt"
	"log"
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

func IsExistUserField(userId uint, field string) bool {
	key := fmt.Sprintf("%s:%d", UserKey, userId)
	exists, err := Rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		log.Println("redis isExistUser 连接失败")
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
		log.Println("redis isExistVideo 连接失败")
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
		log.Println("redis isExistUser 连接失败")
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
