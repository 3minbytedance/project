package mw

import (
	"context"
	"douyin/common"
	"douyin/mw/redis"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// Auth 鉴权中间件
func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Query("token")
		// 解析token
		claims, err := common.ParseToken(token)
		if err != nil {
			// token有误（含len(token)=0的情况），阻止后面函数执行
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
			return
		}
		// 查看token是否在redis中, 若在，给token续期, 若不在，则阻止后面函数执行
		exist := redis.TokenIsExisted(uint(claims.ID))
		if !exist {
			// token有误，阻止后面函数执行
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
			return
		}
		// 给token续期
		redis.SetToken(uint(claims.ID), token)

		zap.L().Debug("CLAIM-ID", zap.Int("ID", int(claims.ID)))
		c.Set(common.ContextUserIDKey, claims.ID)
		c.Next(ctx)
	}
}

// AuthWithoutLogin 未登录情况，若携带token,解析用户id放入context;如果没有携带，则将用户id默认为0
func AuthWithoutLogin() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Query("token")
		var userId uint
		var tokenValid bool
		claims, err := common.ParseToken(token)
		if err != nil {
			// token有误或token为空，未登录
			tokenValid = false
			userId = 0
		} else {
			// 查看token是否在redis中, 若在，则返回用户id, 并且给token续期, 若不在，则将userID设为0
			exist := redis.TokenIsExisted(claims.ID)
			if !exist {
				zap.L().Debug("Token is not existed, user id ", zap.Int("ID", int(claims.ID)))
				// token有误，设置userId为0,tokenValid为false
				userId = 0
				tokenValid = false
			} else {
				userId = claims.ID
				// 给token续期
				redis.SetToken(uint(claims.ID), token)
				tokenValid = true
			}
		}

		zap.L().Debug("to")
		zap.L().Debug("USER-ID", zap.Int("ID", int(userId)))
		c.Set(common.TokenValid, tokenValid)
		c.Set(common.ContextUserIDKey, userId)
		c.Next(ctx)
	}
}

// AuthBody 若token在请求体里，解析token
func AuthBody() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
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
			// 查看token是否在redis中, 若在，则返回用户id, 并且给token续期, 若不在，则返回0
			exist := redis.TokenIsExisted(uint(claims.ID))
			if !exist {
				// token有误，阻止后面函数执行
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			}
			// 给token续期
			redis.SetToken(uint(claims.ID), token)
			c.Set(common.ContextUserIDKey, claims.ID)
			c.Next(ctx)
		}
	}
}
