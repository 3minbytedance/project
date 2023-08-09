package service

import (
	"errors"
	"fmt"
	"log"
	daoMysql "project/dao/mysql"
	daoRedis "project/dao/redis"
	"project/models"
	"strconv"
)

type FavoriteListResponse struct {
	FavoriteRes models.Response
	//VideoList   []Video
	VideoList []int64
}

// FavoriteActions 点赞，取消赞的操作过程
func FavoriteActions(userId int64, videoId int64, actionType int) error {
	var (
		db       = daoMysql.DB
		userRDB  = daoRedis.UserFavoriteRDB
		videoRDB = daoRedis.VideoFavoritedRDB
	)
	_, err := daoMysql.GetFavoritesByUserId(db, userRDB, userId)
	if err != nil {
		return err
	}
	_, err = daoMysql.GetFavoritesByVideoId(db, videoRDB, userId)
	if err != nil {
		return err
	}
	userIdStr := strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	isMember, _ := userRDB.SIsMember(userIdStr, videoId).Result()
	switch actionType {
	case 1:
		// 点赞
		// 更新用户喜欢的视频列表
		if isMember {
			return errors.New("该视频已点赞")
		}
		err = userRDB.SAdd(userIdStr, videoId).Err()
		if err != nil {
			log.Println(err)
		}
		// 更新用户喜欢的视频数量
		err = userRDB.Incr(fmt.Sprintf("%d:count", userId)).Err()
		if err != nil {
			log.Println(err)
		}
		// 更新视频被喜欢的用户列表
		err = videoRDB.SAdd(videoIdStr, userId).Err()
		if err != nil {
			log.Println(err)
		}
		// 更新视频被喜欢的数量
		err = videoRDB.Incr(fmt.Sprintf("%d:count", videoId)).Err()
		if err != nil {
			log.Println(err)
		}
		// 新增到数据库
		db.Create(&models.Favorite{UserId: userId, VideoId: videoId})
	case 2:
		// 取消赞
		// 更新用户喜欢的视频列表
		if !isMember {
			return errors.New("该视频未点赞")
		}
		err = userRDB.SRem(userIdStr, videoId).Err()
		if err != nil {
			log.Println(err)
		}
		// 更新用户喜欢的视频数量
		err = userRDB.Decr(fmt.Sprintf("%d:count", userId)).Err()
		if err != nil {
			log.Println(err)
		}
		// 更新视频被喜欢的用户列表
		err = videoRDB.SRem(videoIdStr, userId).Err()
		if err != nil {
			log.Println(err)
		}
		// 更新视频被喜欢的数量
		err = videoRDB.Decr(fmt.Sprintf("%d:count", videoId)).Err()
		if err != nil {
			log.Println(err)
		}
		db.Where("user_id = ? and video_id = ?", userId, videoId).Delete(&models.Favorite{})
	}
	// 更新过期时间
	userRDB.Expire(fmt.Sprintf("%d:count", userId), daoMysql.Expiration)
	userRDB.Expire(userIdStr, daoMysql.Expiration)
	videoRDB.Expire(fmt.Sprintf("%d:count", videoId), daoMysql.Expiration)
	videoRDB.Expire(videoIdStr, daoMysql.Expiration)

	return nil
}

func GetFavoriteList(userId int64) ([]int64, error) {
	var (
		db      = daoMysql.DB
		userRDB = daoRedis.UserFavoriteRDB
	)
	favoritesByUserId, err := daoMysql.GetFavoritesByUserId(db, userRDB, userId)
	if err != nil {
		return nil, err
	}
	return favoritesByUserId, err
}

// GetFavoritesVideoCount 根据视频id，返回该视频的点赞数（外部使用）
func GetFavoritesVideoCount(videoId int64) (int, error) {
	db := daoMysql.DB
	rdb := daoRedis.VideoFavoritedRDB
	_, num, err := daoMysql.GetFavoritesById(db, rdb, videoId, daoMysql.IdTypeVideo)
	return num, err
}
