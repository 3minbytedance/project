package redis

import "github.com/redis/go-redis/v9"

func Ping() error {
	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		return err
	}
	return nil
}

func RunScript(script string, keys []string, args ...interface{}) (interface{}, error) {
	val, err := redis.NewScript(script).Run(Ctx, Rdb, keys, args).Result()
	if err != nil && err != redis.Nil {
		return nil, nil
	}
	return val, nil
}

func ClearAll() error {
	if _, err := Rdb.FlushAll(Ctx).Result(); err != nil {
		return err
	}
	return nil
}

const (
	rateScript string = `
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

