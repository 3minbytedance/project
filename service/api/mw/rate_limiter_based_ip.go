package mw

import (
	"context"
	"douyin/common"
	"douyin/mw/redis"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strings"
	"time"
)

var (
	timeNow                 = time.Now
	bucketSize      int     = 1000 // 令牌桶的容量
	refillPerSecond float64 = 60   // 每隔多长时间添加令牌
	refillToken     int     = 10   // 每次令牌桶填充操作时添加的令牌数量
)

// RateLimiter 限流器
func RateLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

		ipAddr := strings.Split(c.ClientIP(), ":")[0]
		permit, remain, err := Acquire(ipAddr)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: -1,
				StatusMsg:  "Failed to acquire IP",
			})
			c.Abort()
			return
		}

		if !permit {
			c.JSON(http.StatusTooManyRequests, Response{
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

func Acquire(key string) (bool, int, error) {
	now := timeNow()
	cacheKey := fmt.Sprintf("tokenbucket:%s", key)

	remain, err := redis.RunScript(
		script,
		[]string{cacheKey},
		now.Unix(),
		refillToken,
		refillPerSecond,
		bucketSize,
	)
	if err != nil {
		//context.WithField("err", err).Error("redis.RunScript failed")
		return false, 0, err
	}
	if remain.(int64) < 0 {
		return false, bucketSize, nil
	}
	return true, int(remain.(int64)), nil
}

const (
	script string = `
	local capacity = tonumber(ARGV[4])
	local refill = tonumber(ARGV[3])
	local refillToken = tonumber(ARGV[2])
	local ts = tonumber(ARGV[1])
	local lastUpdate = ts
	local remainToken = capacity

	local last = redis.call('HMGET', KEYS[1], 'ts', 'tokens')
	if last[1] then
		local lastTs = tonumber(last[1])
		local lastTokens = tonumber(last[2])
		local refillCount = math.floor((ts - lastTs) / refill)


		remainToken = math.min(capacity, lastTokens + (refillCount * refillToken))
		lastUpdate = math.min(ts, lastTs + (refillCount * refill))
	end

	if remainToken >= 0 then
			remainToken = remainToken - 1
	end
	redis.call('HMSET', KEYS[1], 'ts', ts, 'tokens', remainToken)
	redis.call('EXPIRE', KEYS[1], math.ceil(capacity / refill))
	return remainToken
	`
)
