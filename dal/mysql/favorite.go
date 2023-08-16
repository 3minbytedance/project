package mysql

import (
	"gorm.io/gorm"
	"log"
	"project/models"
)

var (
	IdTypeVideo = 1
	IdTypeUser  = 2
)

// GetFavoritesByIdFromMysql 从数据库中根据Id类型获取对应的数据
func GetFavoritesByIdFromMysql(id uint, idType int) ([]models.Favorite, int, error) {
	var (
		res  []models.Favorite
		rows int64
		err  error
	)
	switch idType {
	case IdTypeVideo:
		dbStruct := DB.Where("video_id = ?", id).Find(&res)
		rows = dbStruct.RowsAffected
		err = DB.Error
	case IdTypeUser:
		dbStruct := DB.Where("user_id = ?", id).Find(&res)
		rows = dbStruct.RowsAffected
		err = DB.Error
	}
	return res, int(rows), err
}

// AddUserFavorite 添加喜欢关系
func AddUserFavorite(userId, videoId uint) bool {
	follow := models.Favorite{UserId: userId, VideoId: videoId}
	result := DB.Model(&models.Favorite{}).Create(&follow)
	return result.RowsAffected != 0
}

// DeleteUserFavorite 删除喜欢关系
func DeleteUserFavorite(userId, videoId uint) error {
	favorite := models.Favorite{UserId: userId, VideoId: videoId}
	result := DB.Delete(&models.Favorite{}, favorite)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		log.Println("未找到 Follow", userId, videoId)
		return result.Error
	}
	return nil
}

//
//// GetFavoritedUserCount 根据用户id，返回该用户的点赞的视频数（外部使用）
//func GetFavoritedUserCount(userId int64) (int, error) {
//	db := daoMysql.DB
//	rdb := daoRedis.UserFavoriteRDB
//	_, num, err := getFavoritesById(db, rdb, userId, idTypeUser)
//	return num, err
//}
