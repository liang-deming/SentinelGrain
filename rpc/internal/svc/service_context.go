package svc

import (
	"context"
	"time"

	"SentinelGrain/common/algorithm"
	"SentinelGrain/common/cache"
	"SentinelGrain/common/quota"
	"SentinelGrain/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config      config.Config
	BizRedis    *redis.Redis   // Redis客户端
	L1Cache     *cache.L1Cache // 本地缓存
	LimitScript string       // Redis Lua脚本SHA
	QuotaRules  *quota.RuleCache
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient := redis.MustNewRedis(c.Redis)

	var l1cache *cache.L1Cache
	if c.L1Cache.Enabled {
		var err error
		l1cache, err = cache.NewL1Cache(time.Duration(c.L1Cache.TTL) * time.Millisecond)
		if err != nil {
			panic("Failed to create L1 cache: " + err.Error())
		}
	}

	script := algorithm.GetSlideWindowScript()
	scriptSha, err := redisClient.ScriptLoad(script)
	if err != nil {
		panic("Failed to load limit script: " + err.Error())
	}

	cmdTimeout := time.Duration(c.CommandTimeout) * time.Millisecond
	if cmdTimeout <= 0 {
		cmdTimeout = 3 * time.Second
	}
	rc := quota.NewRuleCache(redisClient, cmdTimeout)
	loadCtx, cancel := context.WithTimeout(context.Background(), cmdTimeout*10)
	defer cancel()
	if err := rc.Refresh(loadCtx); err != nil {
		panic("quota RuleCache initial load: " + err.Error())
	}

	return &ServiceContext{
		Config:      c,
		BizRedis:    redisClient,
		L1Cache:     l1cache,
		LimitScript: scriptSha,
		QuotaRules:  rc,
	}
}
