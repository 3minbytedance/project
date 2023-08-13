package redis

import (
	"fmt"
	"strconv"
)

func GetTotalFavoriteByUserId(userId uint) (int64, error) {
	key := UserKey + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, TotalFavoriteField).Result()
	totalFavorite, _ := strconv.ParseInt(count, 10, 64)
	return totalFavorite, err
}

func GetFavoriteCountByUserId(userId uint) (int64, error) {
	key := UserKey + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, FavoriteCountField).Result()
	favoriteCount, _ := strconv.ParseInt(count, 10, 64)
	return favoriteCount, err
}

func IncrementTotalFavoriteByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavoriteField, 1).Result()
	return err
}

func IncrementFavoriteCountByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, FavoriteCountField, 1).Result()
	return err
}

func DecrementTotalFavoriteByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavoriteField, -1).Result()
	return err
}

func DecrementFavoriteCountByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, FavoriteCountField, -1).Result()
	return err
}

func SetTotalFavoriteByUserId(userId uint, totalFavorite int64) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, TotalFavoriteField, totalFavorite).Err()
	return err
}

func SetFavoriteCountByUserId(userId uint, favoriteCount int64) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, FavoriteCountField, favoriteCount).Err()
	return err
}
