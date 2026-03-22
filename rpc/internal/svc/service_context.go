package svc

import (
	"SentinelGrain/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

type ServiceContext struct {
	Config      config.Config
	BizRedis    *redis.Redis     // Redis客户端
	L1Cache     *cache.L1Cache   // 本地缓存
	LimitScript string          // Redis Lua脚本SHA
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建Redis客户端
	redisClient := redis.MustNewRedis(c.Redis)

	// 创建本地缓存(如果启用)
	var l1cache *cache.L1Cache
	if c.L1Cache.Enabled {
		l1cache = cache.NewL1Cache(time.Duration(c.L1Cache.TTL) * time.Millisecond)
	}

	// 加载Lua脚本
	script := algorithm.GetSlideWindowScript()
	scriptSha, err := redisClient.ScriptLoad(script)
	if err != nil {
		panic("Failed to load limit script: " + err.Error())
	}

	return &ServiceContext{
		Config:      c,
		BizRedis:    redisClient,
		L1Cache:     l1cache,
		LimitScript: scriptSha,
	}
}
