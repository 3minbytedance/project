package redis

import (
	"fmt"
	"strconv"
)

const (
	TotalFavorite = "total_favorited" //发布视频的总获赞数量
	WorkCount     = "work_count"      //作品数
	FavoriteCount = "favorite_count"  //喜欢数
	Name          = "name"            //用户名
)

func GetWorkCountByUserId(userId uint) (int64, error) {
	key := User + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, WorkCount).Result()
	commentCount, _ := strconv.ParseInt(count, 10, 64)
	return commentCount, err
}

func GetNameByUserId(userId uint) (string, error) {
	key := User + fmt.Sprintf("%d", userId)
	name, err := Rdb.HGet(Ctx, key, Name).Result()
	return name, err
}

func GetTotalFavoriteByUserId(userId uint) (int64, error) {
	key := User + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, TotalFavorite).Result()
	totalFavorite, _ := strconv.ParseInt(count, 10, 64)
	return totalFavorite, err
}

func GetFavoriteCountByUserId(userId uint) (int64, error) {
	key := User + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, FavoriteCount).Result()
	favoriteCount, _ := strconv.ParseInt(count, 10, 64)
	return favoriteCount, err
}

func IncrementWorkCountByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, WorkCount, 1).Result()
	return err
}

func IncrementTotalFavoriteByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavorite, 1).Result()
	return err
}

func IncrementFavoriteCountByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, FavoriteCount, 1).Result()
	return err
}

func DecrementTotalFavoriteByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavorite, -1).Result()
	return err
}

func DecrementFavoriteCountByUserId(userId uint) error {
	key := User + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, FavoriteCount, -1).Result()
	return err
}

func SetNameByUserId(userId uint, name string) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, Name, name).Err()
	return err
}

func SetTotalFavoriteByUserId(userId uint, totalFavorite int64) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, TotalFavorite, totalFavorite).Err()
	return err
}

func SetWorkCountByUserId(userId uint, workCount int64) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, WorkCount, workCount).Err()
	return err
}

func SetFavoriteCountByUserId(userId uint, favoriteCount int64) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, FavoriteCount, favoriteCount).Err()
	return err
}
