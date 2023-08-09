package mysql

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"log"
	"project/models"
	"strconv"
	"time"
)

var (
	IdTypeVideo     = 1
	IdTypeUser      = 2
	Expiration      = time.Hour * 2
	StoreExpiration = Expiration / 2
)

/* 这个是用来记录用户喜欢的视频的
UerFavoriteRDB：用来存储用户点赞的视频的信息（视频ID集合），以及点赞的视频数量
VideoFavoritedRDB：用来存储视频被点赞的信息（用户ID集合），以及被点赞的数量
Redis持久化时间为两倍的过期时间，用定时器管理，每次持久化过期时间及以下的数据
*/

func getFavoritesByIdFromMysql(db *gorm.DB, id int64, idType int) ([]models.Favorite, int, error) {
	var (
		res  []models.Favorite
		rows int64
		err  error
	)
	switch idType {
	case IdTypeVideo:
		dbStruct := db.Where("video_id = ?", id).Find(&res)
		rows = dbStruct.RowsAffected
		err = db.Error
	case IdTypeUser:
		dbStruct := db.Where("user_id = ?", id).Find(&res)
		rows = dbStruct.RowsAffected
		err = db.Error
	}
	return res, int(rows), err
}

// GetFavoritesByUserId 获取当前id的点赞的视频id列表
func GetFavoritesByUserId(db *gorm.DB, rdb *redis.Client, userId int64) ([]int64, error) {
	idList, _, err := GetFavoritesById(db, rdb, userId, IdTypeUser)
	return idList, err
}

func GetFavoritesByVideoId(db *gorm.DB, rdb *redis.Client, videoId int64) ([]int64, error) {
	idList, _, err := GetFavoritesById(db, rdb, videoId, IdTypeVideo)
	return idList, err
}

func GetFavoritesById(db *gorm.DB, rdb *redis.Client, id int64, idType int) ([]int64, int, error) {
	// 先从redis中取数据
	key := strconv.FormatInt(id, 10)
	numKey := fmt.Sprintf("%d:count", id)
	result, err := rdb.Exists(numKey).Result()
	if err != nil {
		log.Println(err.Error())
		return []int64{}, 0, err
	}
	if result > 0 {
		// redis里有对应数据的情况
		favoritesStrList, err := rdb.SMembers(key).Result()
		if err != nil {
			log.Println(err.Error())
		}
		res, err := convertStrListToInt64List(favoritesStrList)
		if err != nil {
			log.Println(err.Error())
		}
		return res, len(res), err
	} else {
		// redis中没有对应的数据，从MYSQL数据库中获取数据
		favorites, num, err := getFavoritesByIdFromMysql(db, id, idType)
		if err != nil {
			log.Println(err.Error())
		}
		idList := getIdListFromFavoriteSlice(favorites, idType)
		// key 不存在需要同步到redis，如果这个时候发生点赞怎么办
		loadSetToRedis(key, idList, rdb)
		loadCountToRedis(numKey, num, rdb)
		return idList, num, err
	}
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
func getIdListFromFavoriteSlice(favorites []models.Favorite, idType int) []int64 {
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
	if len(value) == 0 {

		err := rdb.SAdd(id).Err()
		if err != nil {
			log.Println(err)
		}
	} else {
		for _, v := range value {
			err := rdb.SAdd(id, v).Err()
			if err != nil {
				log.Println(err)
			}
		}
	}
	err := rdb.Expire(id, Expiration).Err()
	if err != nil {
		log.Println(err)
	}
}

// loadCountToRedis 将数值存储在redis中
func loadCountToRedis(id string, count int, rdb *redis.Client) {
	err := rdb.Set(id, count, Expiration).Err()
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

//
//// GetFavoritedUserCount 根据用户id，返回该用户的点赞的视频数（外部使用）
//func GetFavoritedUserCount(userId int64) (int, error) {
//	db := daoMysql.DB
//	rdb := daoRedis.UserFavoriteRDB
//	_, num, err := getFavoritesById(db, rdb, userId, idTypeUser)
//	return num, err
//}
