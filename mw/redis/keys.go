package redis

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	//VideoKey hash 类型 key:video_videoId
	VideoKey = "video_"

	CommentCountField = "comment_count"

	VideoFavoritedCountField = "favorited_count" // 视频被点赞总数
)

const (
	//UserKey hash 类型 key:user_userId
	UserKey = "user_"

	WorkCountField     = "work_count"      //作品数
	NameField          = "name"            //用户名
	TotalFavoriteField = "total_favorited" //发布视频的总获赞数量

	// FavoriteList  set类型
	FavoriteList = "favorite_list" //喜欢视频列表

	// FollowList and FollowerList  set类型
	FollowList   = "follow_list"   //关注列表
	FollowerList = "follower_list" //粉丝列表
)

const TokenKey = "token_"

func IsExistUserField(userId uint, field string) bool {
	key := UserKey + fmt.Sprintf("%d", userId)
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
	key := VideoKey + fmt.Sprintf("%d", videoId)
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
	key := fmt.Sprintf("%s_%d", field, userId)
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

// 根据参数合成并删除key
func DelKey(userId uint, field string) {
	key := fmt.Sprintf("%s_%d", field, userId)
	Rdb.Del(Ctx, key)
}
