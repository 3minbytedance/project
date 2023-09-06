package mw

import (
	"context"
	"douyin/common"
	"douyin/mw/redis"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

// RequestLoginLimiter 限制登录请求次数
func RequestLoginLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		success := redis.IncrementLoginLimiterCount(ip)
		if !success {
			c.Abort()
			c.JSON(http.StatusOK, Response{
				StatusCode: common.CodeLimiterCount,
				StatusMsg:  common.MapErrMsg(common.CodeLimiterCount),
			})
			return
		}
		c.Next(ctx)
	}
}

// RequestCommentLimiter 限制评论请求次数
func RequestCommentLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		success := redis.IncrementCommentLimiterCount(ip)
		if !success {
			c.Abort()
			c.JSON(http.StatusOK, Response{
				StatusCode: common.CodeLimiterCount,
				StatusMsg:  common.MapErrMsg(common.CodeLimiterCount),
			})
			return
		}
		c.Next(ctx)
	}
}

// RequestUploadLimiter 限制上传请求次数
func RequestUploadLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		success := redis.IncrementUploadLimiterCount(ip)
		if !success {
			c.Abort()
			c.JSON(http.StatusOK, Response{
				StatusCode: common.CodeLimiterCount,
				StatusMsg:  common.MapErrMsg(common.CodeLimiterCount),
			})
			return
		}
		c.Next(ctx)
	}
}

// RequestMessageLimiter 限制聊天请求次数
func RequestMessageLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		success := redis.IncrementMessageLimiterCount(ip)
		if !success {
			c.Abort()
			c.JSON(http.StatusOK, Response{
				StatusCode: common.CodeLimiterCount,
				StatusMsg:  common.MapErrMsg(common.CodeLimiterCount),
			})
			return
		}
		c.Next(ctx)
	}
}
