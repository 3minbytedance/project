package mw

import (
	"context"
	"douyin/common"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// Auth 鉴权中间件
func Auth() []app.HandlerFunc {
	return []app.HandlerFunc{
		func(ctx context.Context, c *app.RequestContext) {
			token := c.Query("token")
			// 没携带token
			if len(token) == 0 {
				// 没有token, 阻止后面函数执行
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Unauthorized",
				})

			} else {
				claims, err := common.ParseToken(token)
				if err != nil {
					// token有误，阻止后面函数执行
					c.Abort()
					c.JSON(http.StatusUnauthorized, Response{
						StatusCode: -1,
						StatusMsg:  "Token Error",
					})
				}
				c.Set(common.ContextUserIDKey, claims.ID)
				c.Next(ctx)

			}
		}}
}

// AuthWithoutLogin 未登录情况，若携带token,解析用户id放入context;如果没有携带，则将用户id默认为0
func AuthWithoutLogin() []app.HandlerFunc {
	return []app.HandlerFunc{func(ctx context.Context, c *app.RequestContext) {
		token := c.Query("token")
		var userId uint
		if len(token) == 0 {
			// 没有token, 设置userId为0
			userId = 0
		} else {
			claims, err := common.ParseToken(token)
			if err != nil {
				// token有误，阻止后面函数执行
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				userId = claims.ID
			}
			c.Set(common.ContextUserIDKey, userId)
			c.Next(ctx)
		}
	}}
}

// AuthBody 若token在请求体里，解析token
func AuthBody() []app.HandlerFunc {
	return []app.HandlerFunc{func(ctx context.Context, c *app.RequestContext) {
		token := c.PostForm("token")
		// 没携带token
		if len(token) == 0 {
			// 没有token, 阻止后面函数执行
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		} else {
			claims, err := common.ParseToken(token)
			if err != nil {
				// token有误，阻止后面函数执行
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			}
			c.Set(common.ContextUserIDKey, claims.ID)
			c.Next(ctx)
		}
	}}
}
