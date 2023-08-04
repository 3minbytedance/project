package models

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

/* 这个是用来记录用户喜欢的视频的
UerFavoriteRDB：用来存储用户点赞的视频的信息（视频ID集合），以及点赞的视频数量
VideoFavoritedRDB：用来存储视频被点赞的信息（用户ID集合），以及被点赞的数量
Redis持久化时间为两倍的过期时间，用定时器管理，每次持久化过期时间及以下的数据
*/

type Favorite struct {
	gorm.Model
	UserId     int64
	VideoId    int64
	IsFavorite bool
}

func (*Favorite) TableName() string {
	return "favorite"
}

type FavoriteListResponse struct {
	FavoriteRes Response
	VideoList   []Video
}

func GetFavoritesByUserFromMysql(db *gorm.DB, userId int64) ([]Favorite, int, error) {
	var res []Favorite
	rows := db.Where("user_id = ? and is_favorite = ?", userId, 1).Find(&res).RowsAffected
	if rows != 0 {
		log.Printf("can't find %d in table favorite", userId)
		return nil, 0, errors.New("can't find userId in mysql")
	}
	return res, int(rows), nil
}

func GetFavoritesByVideoFromMysql(db *gorm.DB, videoId int64) ([]Favorite, int, error) {
	var res []Favorite
	rows := db.Where("video_id = ? and is_favorite = ?", videoId, 1).Find(&res).RowsAffected
	if rows != 0 {
		log.Printf("can't find %d in table favorite", videoId)
		return nil, 0, errors.New("can't find userId in mysql")
	}
	return res, int(rows), nil
}

// GetFavoritesByUser 获取当前id的点赞的视频id列表
func GetFavoritesByUser(db *gorm.DB, rdb *redis.Client, userId int64) ([]int64, error) {
	userKey := strconv.FormatInt(userId, 10)
	result, err := rdb.Exists(userKey).Result()
	if err != nil {
		log.Println(err.Error())
	}
	if result > 0 {
		userFavoriteVideos, err := rdb.SMembers(userKey).Result()
		if err != nil {
			log.Println(err.Error())
		}
		res, err := convertStrListToInt64List(userFavoriteVideos)
		if err != nil {
			log.Println(err.Error())
		}
		return res, err
	} else {
		// 从数据库中获取数据
		favorites, num, err := GetFavoritesByUserFromMysql(db, userId)
		if err != nil {
			log.Println(err.Error())
		}
		videos := getIdListFromFavoriteSlice(favorites, 1)
		// key 不存在需要同步到redis，如果这个时候发生点赞怎么办
		numKey := fmt.Sprintf("%d:count", userId)
		loadSetToRedis(userKey, videos, rdb)
		loadCountToRedis(numKey, num, rdb)
		return videos, err
	}
}

func GetFavoritesByVideo(db *gorm.DB, rdb *redis.Client, videoId int64) ([]int64, error) {
	userKey := strconv.FormatInt(videoId, 10)
	result, err := rdb.Exists(userKey).Result()
	if err != nil {
		log.Println(err.Error())
	}
	if result > 0 {
		videoFavoritedUsers, err := rdb.SMembers(userKey).Result()
		if err != nil {
			log.Println(err.Error())
		}
		res, err := convertStrListToInt64List(videoFavoritedUsers)
		if err != nil {
			log.Println(err.Error())
		}
		return res, err
	} else {
		// 从数据库中获取数据
		favorites, num, err := GetFavoritesByVideoFromMysql(db, videoId)
		if err != nil {
			log.Println(err.Error())
		}
		users := getIdListFromFavoriteSlice(favorites, 2)
		// key 不存在需要同步到redis，如果这个时候发生点赞怎么办
		numKey := fmt.Sprintf("%d:count", videoId)
		loadSetToRedis(userKey, users, rdb)
		loadCountToRedis(numKey, num, rdb)
		return users, err
	}
}

// FavoriteActions 点赞，取消赞的操作过程
func FavoriteActions(db *gorm.DB, userRDB *redis.Client, videoRDB *redis.Client,
	userId int64, videoId int64, actionType int) error {
	//userKey := strconv.FormatInt(videoId, 10)
	//res, err := userRDB.Exists(userKey).Result()

	_, err := GetFavoritesByUser(db, userRDB, userId)
	if err != nil {
		return err
	}
	_, err = GetFavoritesByVideo(db, videoRDB, userId)
	if err != nil {
		return err
	}
	switch actionType {
	case 1:
		// 点赞
		userIdStr := strconv.FormatInt(userId, 10)
		err = userRDB.SAdd(userIdStr, videoId).Err()
		if err != nil {
			log.Println(err)
		}
		err = userRDB.Incr(fmt.Sprintf("%d:count", userId)).Err()
		if err != nil {
			log.Println(err)
		}
	case 2:

	}
	return nil
}

func StoreByTimer() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		select {
		case <-ticker.C:
			StoreData()
		}
	}()
}

func StoreData() {
	processData()
	storeDataToMysql()
}

// 辅助函数
// convertStrListToInt64List 将字符串列表转化为int64列表
func convertStrListToInt64List(strs []string) ([]int64, error) {
	res := make([]int64, 0)
	for _, v := range strs {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		res = append(res, int64(vInt))
	}
	return res, nil
}

// getIdListFromFavoriteSlice 从Favorite的slice中获取id的列表
func getIdListFromFavoriteSlice(favorites []Favorite, idType int) []int64 {
	res := make([]int64, 0)
	for _, fav := range favorites {
		switch idType {
		case 1:
			res = append(res, fav.UserId)
		case 2:
			res = append(res, fav.VideoId)
		}
	}
	return res
}

// loadSetToRedis 根据列表和对应的id将其存在redis中
func loadSetToRedis(id string, value []int64, rdb *redis.Client) {
	for _, v := range value {
		err := rdb.SAdd(id, v).Err()
		if err != nil {
			log.Println(err)
		}
	}
	err := rdb.Expire(id, time.Hour*2).Err()
	if err != nil {
		log.Println(err)
	}
}

// loadCountToRedis 将数值存储在redis中
func loadCountToRedis(id string, count int, rdb *redis.Client) {
	err := rdb.Set(id, count, time.Hour*2).Err()
	if err != nil {
		fmt.Println(err)
	}

}

// 处理数据
func processData() {

}

// 存储数据
func storeDataToMysql() {

}

// GetFavoritesVideoCount 根据视频id，返回该视频的点赞数（外部使用）
func GetFavoritesVideoCount() {

}

// GetFavoritedUserCount 根据用户id，返回该用户的点赞的视频数（外部使用）
func GetFavoritedUserCount() {

}
