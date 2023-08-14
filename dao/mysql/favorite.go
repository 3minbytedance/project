package mysql

import (
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

//
//// GetFavoritedUserCount 根据用户id，返回该用户的点赞的视频数（外部使用）
//func GetFavoritedUserCount(userId int64) (int, error) {
//	db := daoMysql.DB
//	rdb := daoRedis.UserFavoriteRDB
//	_, num, err := getFavoritesById(db, rdb, userId, idTypeUser)
//	return num, err
//}
