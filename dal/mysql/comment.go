package mysql

import (
	"douyin/dal/model"
	"go.uber.org/zap"
)

func AddComment(comment *model.Comment) (uint, error) {
	result := DB.Model(model.Comment{}).Create(comment)
	// 判断是否创建成功
	if result.Error != nil {
		zap.L().Error("创建 Comment 失败:", zap.Error(result.Error))
		return 0, result.Error
	} else {
		return comment.ID, nil
	}
}
