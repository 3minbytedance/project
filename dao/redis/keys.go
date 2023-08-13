package redis

import (
	"fmt"
	"log"
)

const (
	//VideoKey hash 类型 key:video_videoId
	VideoKey = "video_"

	CommentCountField = "comment_count"
	//点赞数Field todo
)

const (
	//UserKey hash 类型 key:user_userId
	UserKey = "user_"

	WorkCountField     = "work_count"      //作品数
	NameField          = "name"            //用户名
	TotalFavoriteField = "total_favorited" //发布视频的总获赞数量
	FavoriteCountField = "favorite_count"  //喜欢数
	FollowCountField   = "follow_count"    //关注数
	FollowerCountField = "follower_count"  //粉丝数
)

const VideoPage = "video_page_"

func IsExistUserField(userId uint, field string) bool {
	key := UserKey + fmt.Sprintf("%d", userId)
	exists, err := Rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		log.Println("redis isExistUser 连接失败")
		return false
	}
	if !exists {
		return false
	}
	return true
}

func IsExistVideoField(videoId uint, field string) bool {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	exists, err := Rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		log.Println("redis isExistVideo 连接失败")
		return false
	}
	if !exists {
		return false
	}
	return true
}
