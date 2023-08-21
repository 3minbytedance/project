package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

// 判断是否越权 中间件
func Is_authorized() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		//token := c.Query("token")
		//id := redis.GetToken(token)
		//userId, err := common.GetCurrentUserID(c)
		//if err != nil {
		//	zap.L().Error("Get user id from ctx", zap.Error(err))
		//	c.JSON(http.StatusUnauthorized, Response{
		//		StatusCode: -1,
		//		StatusMsg:  "Unauthorized operation",
		//	})
		//}
		//
		//if id != userId {
		//	c.Abort()
		//	c.JSON(http.StatusUnauthorized, Response{
		//		StatusCode: -1,
		//		StatusMsg:  "Unauthorized operation",
		//	})
		//}

	}
}
