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
func FavoriteActions(userId int64, videoId int64, actionType int) error {
	userIdStr := strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	switch actionType {
	case 1:
		// 点赞
		// 更新用户喜欢的视频列表

		//if isMember {
		//	return errors.New("该视频已点赞")
		//}
		//err = userRDB.SAdd(daoRedis.Ctx, userIdStr, videoId).Err()
		//if err != nil {
		//	log.Println(err)
		//}
		// 更新用户喜欢的视频数量

		// 更新视频被喜欢的用户列表

		// 更新视频被喜欢的数量

		// 新增到数据库
		//db.Create(&models.Favorite{UserId: userId, VideoId: videoId})
	case 2:
		// 取消赞
		// 更新用户喜欢的视频列表
		//if !isMember {
		//	return errors.New("该视频未点赞")
		//}
		//err = userRDB.SRem(daoRedis.Ctx, userIdStr, videoId).Err()
		//if err != nil {
		//	log.Println(err)
		//}
		// 更新用户喜欢的视频数量

		// 更新视频被喜欢的用户列表

		// 更新视频被喜欢的数量

	}

	return nil
}

func GetFavoriteList(userId int64) ([]models.VideoResponse, error) {
	var (
		db      = daoMySQL.DB
		userRDB = daoRedis.UserFavoriteRDB
	)
	favoritesByUserId, err := daoMySQL.GetFavoritesByUserId(db, userRDB, userId)
	if err != nil {
		return nil, err
	}
	videos := make([]models.Video, 0)
	for _, id := range favoritesByUserId {
		videoByVideoId, _ := daoMySQL.FindVideoByVideoId(id)
		videos = append(videos, videoByVideoId)
	}
	return videos, err
}

//GetFavoritesVideoCount 根据视频id，返回该视频的点赞数（外部使用）
//func GetFavoritesVideoCount(videoId int64) (int, error) {
//	db := daoMySQL.DB
//	rdb := daoRedis.VideoFavoritedRDB
//	_, num, err := daoMySQL.GetFavoritesById(db, rdb, videoId, daoMySQL.IdTypeVideo)
//	return num, err
//}
