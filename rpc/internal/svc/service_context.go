package svc

import (
	"SentinelGrain/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

type ServiceContext struct {
	Config      config.Config
	BizRedis    *redis.Redis        // Redis客户端
	L1Cache     *collection.Cache   // 本地缓存
	LimitScript string             // Redis Lua脚本SHA
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建Redis客户端
	redisClient := redis.MustNewRedis(c.Redis)

	// 创建本地缓存(如果启用)
	var l1cache *collection.Cache
	if c.L1Cache.Enabled {
		l1cache = collection.NewCache(time.Duration(c.L1Cache.TTL) * time.Millisecond)
	}

	return &ServiceContext{
		Config:   c,
		BizRedis: redisClient,
		L1Cache:  l1cache,
	}
}
