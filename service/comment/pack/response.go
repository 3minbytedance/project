package pack

import (
	"douyin/dal/model"
	"douyin/kitex_gen/comment"
	user "douyin/kitex_gen/user"
	"time"
)

func Comment(commentModel *model.Comment, userModel *user.User) *comment.Comment {
	if commentModel == nil {
		return nil
	}
	return &comment.Comment{
		Id:         int64(commentModel.ID),
		User:       userModel,
		Content:    commentModel.Content,
		CreateDate: TranslateTime(commentModel.CreatedAt.Unix()),
	}
}

// TranslateTime 返回mm-dd格式
func TranslateTime(createTime int64) string {
	t := time.Unix(createTime, 0)
	return t.Format("01-02")
}
