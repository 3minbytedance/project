package redis

import (
	"fmt"
	"strconv"
)

// GetTotalFavoritedByUserId 获取用户发布视频的所有获赞量
func GetTotalFavoritedByUserId(userId uint) (int64, error) {
	key := UserKey + fmt.Sprintf("%d", userId)
	count, err := Rdb.HGet(Ctx, key, TotalFavoriteField).Result()
	totalFavorited, _ := strconv.ParseInt(count, 10, 64)
	return totalFavorited, err
}

// GetFavoritedCountByVideoId 获取视频被点赞的数量
func GetFavoritedCountByVideoId(videoId uint) (int64, error) {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	count, err := Rdb.HGet(Ctx, key, VideoFavoritedCountField).Result()
	favoritedCount, _ := strconv.ParseInt(count, 10, 64)
	return favoritedCount, err
}

// IncrementTotalFavoritedByUserId 用户视频总被点赞量加1
func IncrementTotalFavoritedByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavoriteField, 1).Result()
	return err
}

// IncrementFavoritedCountByVideoId 视频被点赞数量加1
func IncrementFavoritedCountByVideoId(videoId uint) error {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, VideoFavoritedCountField, 1).Result()
	return err
}

// DecrementTotalFavoritedByUserId 用户视频总被点击量减1
func DecrementTotalFavoritedByUserId(userId uint) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavoriteField, -1).Result()
	return err
}

// DecrementFavoritedCountByVideoId 视频被点赞数量减1
func DecrementFavoritedCountByVideoId(videoId uint) error {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	//减少并返回
	_, err := Rdb.HIncrBy(Ctx, key, VideoFavoritedCountField, -1).Result()
	return err
}

// SetTotalFavoritedByUserId 设置用户发布视频的总的被点赞数
func SetTotalFavoritedByUserId(userId uint, totalFavorite int64) error {
	key := UserKey + fmt.Sprintf("%d", userId)
	err := Rdb.HSet(Ctx, key, TotalFavoriteField, totalFavorite).Err()
	return err
}

// SetVideoFavoritedCountByVideoId 设置该videoId下被点赞总量
func SetVideoFavoritedCountByVideoId(videoId uint, totalFavorited int64) error {
	key := VideoKey + fmt.Sprintf("%d", videoId)
	err := Rdb.HSet(Ctx, key, VideoFavoritedCountField, totalFavorited).Err()
	return err
}

// GetFavoriteListByUserId 根据userId查找喜欢的视频list
func GetFavoriteListByUserId(userId uint) ([]uint, error) {
	key := fmt.Sprintf("%d_%s", userId, FavoriteList)
	list, err := Rdb.SMembers(Ctx, key).Result()
	var result []uint
	for _, i := range list {
		k, _ := strconv.Atoi(i)
		result = append(result, uint(k))
	}
	return result, err
}

// SetFavoriteListByUserId 设置用户的点赞视频列表
func SetFavoriteListByUserId(userid uint, id []uint) error {
	key := fmt.Sprintf("%d_%s", userid, FavoriteList)

	b := make([]interface{}, len(id))
	for i, v := range id {
		b[i] = v
	}
	err := Rdb.SAdd(Ctx, key, b...).Err()
	return err
}

// AddFavoriteVideoToList 给用户的点赞视频列表加一个video
func AddFavoriteVideoToList(userId uint, videoId uint) error {
	key := fmt.Sprintf("%d_%s", userId, FavoriteList)
	err := Rdb.SAdd(Ctx, key, videoId).Err()
	return err
}

// DeleteFavoriteVideoFromList 给用户的点赞视频列表删除一个video
func DeleteFavoriteVideoFromList(userId uint, videoId uint) error {
	key := fmt.Sprintf("%d_%s", userId, FavoriteList)
	err := Rdb.SRem(Ctx, key, videoId).Err()
	return err
}

// IsInUserFavoriteList 判断用户点赞视频列表中是否有对应的video
func IsInUserFavoriteList(userId uint, videoId uint) bool {
	key := fmt.Sprintf("%d_%s", userId, FavoriteList)
	found, _ := Rdb.SIsMember(Ctx, key, videoId).Result()
	return found
}

func IsUserFavoriteNil(userId uint) bool {
	key := fmt.Sprintf("%d_%s", userId, FavoriteList)
	found, _ := Rdb.SIsMember(Ctx, key, 0).Result()
	return found
}

func DeleteUserFavoriteNil(userId uint) error {
	key := fmt.Sprintf("%d_%s", userId, FavoriteList)
	err := Rdb.SRem(Ctx, key, 0).Err()
	return err
}

// GetUserFavoriteVideoCountById 根据userId查找喜欢的视频数量
func GetUserFavoriteVideoCountById(userId uint) (int64, error) {
	key := fmt.Sprintf("%d_%s", userId, FavoriteList)
	size, err := Rdb.SCard(Ctx, key).Result()
	return size, err
}
