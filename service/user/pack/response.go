package pack

import (
	"douyin/dal/model"
	user "douyin/kitex_gen/user"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

func User(userModel *model.User) *user.User {
	if userModel == nil {
		return nil
	}
	return &user.User{
		Id:              int32(userModel.ID),
		Name:            userModel.Username,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          thrift.StringPtr(userModel.Avatar),
		BackgroundImage: thrift.StringPtr(userModel.BackgroundImage),
		Signature:       thrift.StringPtr(userModel.Signature),
		TotalFavorited:  nil,
		WorkCount:       nil,
		FavoriteCount:   nil,
	}
}

// TranslateTime 返回mm-dd格式
func TranslateTime(createTime int64) string {
	t := time.Unix(createTime, 0)
	return t.Format("01-02")
}
