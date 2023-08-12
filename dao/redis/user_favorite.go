package redis

import (
	"fmt"
	"strconv"
)

const (
	TotalFavorite = "total_favorited" //发布视频的总获赞数量
	FavoriteCount = "favorite_count"  //喜欢数
)

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

func SetTotalFavoriteByUserId(userId uint, totalFavorite int64) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, TotalFavorite, totalFavorite).Err()
	return err
}

func SetFavoriteCountByUserId(userId uint, favoriteCount int64) error {
	key := User + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, FavoriteCount, favoriteCount).Err()
	return err
}
