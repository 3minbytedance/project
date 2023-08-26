package redis

import (
	"fmt"
	"go.uber.org/zap"
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
	key := fmt.Sprintf("%s_%d", FavoriteList, userId)
	list, err := Rdb.SMembers(Ctx, key).Result()
	var result []uint
	for _, i := range list {
		k, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(k))
	}
	return result, err
}

// SetFavoriteListByUserId 设置用户的点赞视频列表
func SetFavoriteListByUserId(userId uint, id []uint) error {
	key := fmt.Sprintf("%s_%d", FavoriteList, userId)
	pipe := Rdb.Pipeline()
	for _, value := range id {
		err := pipe.SAdd(Ctx, key, value).Err()
		if err != nil {
			return err
		}
	}
	zap.L().Info("Favorite_LIST", zap.Any("List", id))
	_, err := pipe.Exec(Ctx)
	return err
}

// AddFavoriteVideoToList 给用户的点赞视频列表加一个video
func AddFavoriteVideoToList(userId uint, videoId uint) error {
	key := fmt.Sprintf("%s_%d", FavoriteList, userId)
	err := Rdb.SAdd(Ctx, key, videoId).Err()
	return err
}

// DeleteFavoriteVideoFromList 给用户的点赞视频列表删除一个video
func DeleteFavoriteVideoFromList(userId uint, videoId uint) error {
	key := fmt.Sprintf("%s_%d", FavoriteList, userId)
	err := Rdb.SRem(Ctx, key, videoId).Err()
	return err
}

// IsInUserFavoriteList 判断用户点赞视频列表中是否有对应的video
func IsInUserFavoriteList(userId uint, videoId uint) bool {
	key := fmt.Sprintf("%s_%d", FavoriteList, userId)
	found, err := Rdb.SIsMember(Ctx, key, videoId).Result()
	if err != nil {
		return false
	}
	return found
}

// GetUserFavoriteVideoCountById 根据userId查找喜欢的视频数量
func GetUserFavoriteVideoCountById(userId uint) (int64, error) {
	key := fmt.Sprintf("%s_%d", FavoriteList, userId)
	size, err := Rdb.SCard(Ctx, key).Result()
	return size, err
}
