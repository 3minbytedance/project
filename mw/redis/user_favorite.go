package redis

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GetTotalFavoritedByUserId 获取用户发布视频的所有获赞量
func GetTotalFavoritedByUserId(userId uint) (int64, error) {
	baseSliceUser := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceUser, Delimiter)
	count, err := Rdb.HGet(Ctx, key, TotalFavoriteField).Result()
	totalFavorited, _ := strconv.ParseInt(count, 10, 64)
	return totalFavorited, err
}

// GetFavoritedCountByVideoId 获取视频被点赞的数量
func GetFavoritedCountByVideoId(videoId uint) (int64, error) {
	baseSliceVideo := []string{VideoKey, strconv.Itoa(int(videoId))}
	key := strings.Join(baseSliceVideo, Delimiter)

	count, err := Rdb.HGet(Ctx, key, VideoFavoritedCountField).Result()
	favoritedCount, _ := strconv.ParseInt(count, 10, 64)
	return favoritedCount, err
}

// IncrementTotalFavoritedByUserId 用户视频总被点赞量加1
func IncrementTotalFavoritedByUserId(userId uint) error {
	baseSliceUser := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceUser, Delimiter)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavoriteField, 1).Result()
	return err
}

// IncrementFavoritedCountByVideoId 视频被点赞数量加1
func IncrementFavoritedCountByVideoId(videoId uint) error {
	baseSliceVideo := []string{VideoKey, strconv.Itoa(int(videoId))}
	key := strings.Join(baseSliceVideo, Delimiter)
	//增加并返回
	_, err := Rdb.HIncrBy(Ctx, key, VideoFavoritedCountField, 1).Result()
	return err
}

// DecrementTotalFavoritedByUserId 用户视频总被点击量减1
func DecrementTotalFavoritedByUserId(userId uint) error {
	baseSliceUser := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceUser, Delimiter)
	_, err := Rdb.HIncrBy(Ctx, key, TotalFavoriteField, -1).Result()
	return err
}

// DecrementFavoritedCountByVideoId 视频被点赞数量减1
func DecrementFavoritedCountByVideoId(videoId uint) error {
	baseSliceVideo := []string{VideoKey, strconv.Itoa(int(videoId))}
	key := strings.Join(baseSliceVideo, Delimiter)
	//减少并返回
	_, err := Rdb.HIncrBy(Ctx, key, VideoFavoritedCountField, -1).Result()
	return err
}

// SetTotalFavoritedByUserId 设置用户发布视频的总的被点赞数
func SetTotalFavoritedByUserId(userId uint, totalFavorite int64) error {
	baseSliceUser := []string{UserKey, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceUser, Delimiter)
	err := Rdb.HSet(Ctx, key, TotalFavoriteField, totalFavorite).Err()
	randomSeconds := 600 + rand.Intn(31) // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}

// SetVideoFavoritedCountByVideoId 设置该videoId下被点赞总量
func SetVideoFavoritedCountByVideoId(videoId uint, totalFavorited int64) error {
	baseSliceVideo := []string{VideoKey, strconv.Itoa(int(videoId))}
	key := strings.Join(baseSliceVideo, Delimiter)

	err := Rdb.HSet(Ctx, key, VideoFavoritedCountField, totalFavorited).Err()
	randomSeconds := 600 + rand.Intn(31) // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}

// GetFavoriteListByUserId 根据userId查找喜欢的视频list
func GetFavoriteListByUserId(userId uint) ([]uint, error) {
	baseSliceFavorite := []string{FavoriteList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceFavorite, Delimiter)
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
	baseSliceFavorite := []string{FavoriteList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceFavorite, Delimiter)

	pipe := Rdb.Pipeline()
	for _, value := range id {
		err := pipe.SAdd(Ctx, key, value).Err()
		if err != nil {
			return err
		}
	}
	_, err := pipe.Exec(Ctx)
	randomSeconds := 600 + rand.Intn(31) // 600秒到630秒之间的随机数
	expiration := time.Duration(randomSeconds) * time.Second
	Rdb.Expire(Ctx, key, expiration)
	return err
}

// AddFavoriteVideoToList 给用户的点赞视频列表加一个video
func AddFavoriteVideoToList(userId uint, videoId uint) error {
	baseSliceFavorite := []string{FavoriteList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceFavorite, Delimiter)

	err := Rdb.SAdd(Ctx, key, videoId).Err()
	return err
}

// DeleteFavoriteVideoFromList 给用户的点赞视频列表删除一个video
func DeleteFavoriteVideoFromList(userId uint, videoId uint) error {
	baseSliceFavorite := []string{FavoriteList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceFavorite, Delimiter)
	err := Rdb.SRem(Ctx, key, videoId).Err()
	return err
}

// IsInUserFavoriteList 判断用户点赞视频列表中是否有对应的video
func IsInUserFavoriteList(userId uint, videoId uint) bool {
	baseSliceFavorite := []string{FavoriteList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceFavorite, Delimiter)
	found, err := Rdb.SIsMember(Ctx, key, videoId).Result()
	if err != nil {
		return false
	}
	return found
}

// GetUserFavoriteVideoCountById 根据userId查找喜欢的视频数量
func GetUserFavoriteVideoCountById(userId uint) (int64, error) {
	baseSliceFavorite := []string{FavoriteList, strconv.Itoa(int(userId))}
	key := strings.Join(baseSliceFavorite, Delimiter)
	size, err := Rdb.SCard(Ctx, key).Result()
	return size, err
}

// ActionLike
// 更新用户喜欢的视频列表,更新视频被喜欢的数量,更新视频作者的被点赞量
func ActionLike(userId, videoId, authorId uint) error {
	// 用户喜欢的视频列表
	baseSliceFavorite := []string{FavoriteList, strconv.Itoa(int(userId))}
	favoriteListKey := strings.Join(baseSliceFavorite, Delimiter)
	// 视频被喜欢的数量
	baseSliceVideo := []string{VideoKey, strconv.Itoa(int(videoId))}
	videoKey := strings.Join(baseSliceVideo, Delimiter)
	// 视频作者的被点赞量
	baseSliceUser := []string{UserKey, strconv.Itoa(int(authorId))}
	userKey := strings.Join(baseSliceUser, Delimiter)

	pipe := Rdb.TxPipeline()
	pipe.SAdd(Ctx, favoriteListKey, videoId)
	pipe.Del(Ctx, videoKey, VideoFavoritedCountField)
	pipe.Del(Ctx, userKey, TotalFavoriteField)
	_, err := pipe.Exec(Ctx)
	return err
}

// ActionCancelLike
// 用户取消点赞，减少用户喜欢的视频列表,减少视频被喜欢的数量,减少视频作者的被点赞量
func ActionCancelLike(userId, videoId, authorId uint) error {
	// 用户喜欢的视频列表
	baseSliceFavorite := []string{FavoriteList,  strconv.Itoa(int(userId))}
	favoriteListKey := strings.Join(baseSliceFavorite, Delimiter)
	// 视频被喜欢的数量
	baseSliceVideo := []string{VideoKey, strconv.Itoa(int(videoId))}
	videoKey := strings.Join(baseSliceVideo, Delimiter)
	// 视频作者的被点赞量
	baseSliceUser := []string{UserKey,  strconv.Itoa(int(authorId))}
	userKey := strings.Join(baseSliceUser, Delimiter)
	pipe := Rdb.TxPipeline()
	pipe.SRem(Ctx, favoriteListKey, videoId)
	pipe.Del(Ctx, videoKey, VideoFavoritedCountField)
	pipe.Del(Ctx, userKey, TotalFavoriteField)
	_, err := pipe.Exec(Ctx)
	return err
}
