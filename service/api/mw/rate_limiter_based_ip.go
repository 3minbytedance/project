package mw

import (
	"context"
	"douyin/common"
	"douyin/mw/redis"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strings"
)

// RateLimiter 限流器
func RateLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ipAddr := strings.Split(c.ClientIP(), ":")[0]
		permit, remain, err := redis.AcquireBucket(ipAddr)
		if err != nil || !permit {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: common.CodeLimiterCount,
				StatusMsg:  common.MapErrMsg(common.CodeLimiterCount),
			})
			c.Abort()
			return
		}
		c.Set("reqCount", remain)
		c.Next(ctx)
	}
}
