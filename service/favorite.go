package service

import (
	"errors"
	"fmt"
	"log"
	daoMySQL "project/dao/mysql"
	daoRedis "project/dao/redis"
	"project/models"
	"strconv"
)

// FavoriteActions 点赞，取消赞的操作过程
func FavoriteActions(userId uint, videoId uint, actionType int) error {
	// 判断是否在redis中，如果没有的话，一起加载到redis中
	_, err := GetFavoritesByUserId(userId)
	if err != nil {
		return err
	}
	_, err = GetFavoritesVideoCount(videoId)
	if err != nil {
		return err
	}
	video, _ := daoMySQL.FindVideoByVideoId(videoId)
	// 如果err不为空，那么一定存在数据库中了
	switch actionType {
	case 1:
		// 点赞
		// 更新用户喜欢的视频列表
		if daoRedis.IsInUserFavoriteList(userId, videoId) {
			return errors.New("该视频已点赞")
		}
		// 更新视频被喜欢的用户列表
		err = daoRedis.AddFavoriteVideoToList(userId, videoId)
		if err != nil {
			fmt.Println(err)
		}
		// 更新用户喜欢的视频数量，这个不用，直接从set中获取
		// 更新视频被喜欢的数量
		err = daoRedis.IncrementFavoritedCountByVideoId(videoId)
		if err != nil {
			fmt.Println(err)
		}
		// 更新视频作者的被点赞量
		err = daoRedis.IncrementTotalFavoritedByUserId(video.AuthorId)
		if err != nil {
			fmt.Println(err)
		}
		// 新增到数据库 todo
		return err
	case 2:
		// 取消赞
		if !daoRedis.IsInUserFavoriteList(userId, videoId) {
			return errors.New("该视频未点赞")
		}
		// 更新视频被喜欢的用户列表
		err = daoRedis.DeleteFavoriteVideoFromList(userId, videoId)
		if err != nil {
			fmt.Println(err)
		}
		// 更新视频被喜欢的数量
		err = daoRedis.DecrementFavoritedCountByVideoId(videoId)
		if err != nil {
			fmt.Println(err)
		}
		// 更新视频作者的被点赞量
		err = daoRedis.DecrementTotalFavoritedByUserId(video.AuthorId)
		if err != nil {
			fmt.Println(err)
		}
		// 数据库删除数据 todo
	}
	return nil
}

// GetFavoriteList 根据用户id获取用户点赞的视频列表
func GetFavoriteList(userId uint) ([]models.VideoResponse, error) {
	favoritesByUserId, err := GetFavoritesByUserId(userId)
	if err != nil {
		return nil, err
	}
	videos := make([]models.Video, 0)
	for _, id := range favoritesByUserId {
		videoByVideoId, _ := daoMySQL.FindVideoByVideoId(id)
		videos = append(videos, videoByVideoId)
	}
	// 从video的阿斗videoResponse
	videoResponses := make([]models.VideoResponse, 0)
	for _, video := range videos {
		user, _ := GetUserInfoByUserId(userId)
		commentCount, _ := GetCommentCount(video.ID)
		favoriteCount, _ := GetFavoritesVideoCount(video.ID)
		videoResponse := models.VideoResponse{
			ID:            video.ID,
			Author:        user,
			PlayUrl:       oss + video.VideoUrl,
			CoverUrl:      oss + video.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    true,
			Title:         video.Title,
		}
		videoResponses = append(videoResponses, videoResponse)
	}
	return videoResponses, err
}

// GetFavoritesVideoCount 根据视频id，返回该视频的点赞数
func GetFavoritesVideoCount(videoId uint) (int64, error) {
	// 判断redis中是否存在对应的video数据
	exits := daoRedis.IsExistVideoField(videoId, daoRedis.VideoFavoritedCountField)
	if exits {
		// redis中存在对应的数据
		count, err := daoRedis.GetFavoritedCountByVideoId(videoId)
		if err != nil {
			fmt.Println(err)
		}
		return count, err
	} else {
		// redis中不存在，从数据库中读取
		_, num, err := daoMySQL.GetFavoritesByIdFromMysql(videoId, daoMySQL.IdTypeVideo)
		if err != nil {
			log.Println(err)
		}
		err = daoRedis.SetTotalFavoritedByVideoId(videoId, int64(num)) // 加载视频被点赞数量
		if err != nil {
			fmt.Println(err)
		}
		return int64(num), err
	}
}

// GetFavoritesByUserId 获取当前id的点赞的视频id列表
func GetFavoritesByUserId(userId uint) ([]uint, error) {
	// 查看redis是否存在对应的user数据
	exist := daoRedis.IsExistUserSetField(userId, daoRedis.FavoriteList)
	if exist {
		// redis存在
		favoritesVideoIdList, err := daoRedis.GetFavoriteListByUserId(userId)
		if err != nil {
			fmt.Println(err)
		}
		return favoritesVideoIdList, err
	} else {
		// redis中没有对应的数据，从MYSQL数据库中获取数据
		favorites, _, err := daoMySQL.GetFavoritesByIdFromMysql(userId, daoMySQL.IdTypeUser)
		if err != nil {
			log.Println(err)
		}
		idList := getIdListFromFavoriteSlice(favorites, daoMySQL.IdTypeUser)
		// key 不存在需要同步到redis
		err = daoRedis.SetFavoriteListByUserId(userId, idList) // 加载到set中
		if err != nil {
			fmt.Println(err)
		}
		err = daoRedis.SetTotalFavoritedByUserId(userId, getUserTotalFavoritedCount(userId)) // 加载用户发布视频被点赞的总数
		if err != nil {
			fmt.Println(err)
		}
		return idList, err
	}
}

// 辅助函数
// getIdListFromFavoriteSlice 从Favorite的slice中获取id的列表
func getIdListFromFavoriteSlice(favorites []models.Favorite, idType int) []uint {
	res := make([]uint, 0)
	for _, fav := range favorites {
		switch idType {
		case 1:
			res = append(res, fav.ID)
		case 2:
			res = append(res, fav.VideoId)
		}
	}
	return res
}

// getUserTotalFavoritedCount获取用户发布视频的总的被点赞数量
func getUserTotalFavoritedCount(userId uint) int64 {
	var total int64
	var users []models.User
	// 获取用户发布的视频列表
	videosByAuthorId, exist := daoMySQL.FindVideosByAuthorId(userId)
	if !exist {
		return 0
	}
	idList := make([]string, 0)
	for _, video := range videosByAuthorId {
		idList = append(idList, strconv.Itoa(int(video.ID)))
	}

	total = daoMySQL.DB.Where("video_id IN ?", idList).Find(&users).RowsAffected
	return total
}
