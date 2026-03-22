package cache

import (
	"github.com/zeromicro/go-zero/core/collection"
	"time"
)

// L1Cache 本地缓存封装
// 仅缓存最近一次Check结果，TTL应保持在~500ms量级以确保与Redis窗口计数的最终一致性
type L1Cache struct {
	cache *collection.Cache
}

// NewL1Cache 创建本地缓存实例
func NewL1Cache(ttl time.Duration) (*L1Cache, error) {
	cache, err := collection.NewCache(ttl)
	if err != nil {
		return nil, err
	}
	return &L1Cache{
		cache: cache,
	}, nil
}

// CheckResult 限流检查结果
type CheckResult struct {
	Allowed   bool
	Remaining int64
}

// GetResult 获取缓存的检查结果
func (c *L1Cache) GetResult(key string) (*CheckResult, bool) {
	if val, ok := c.cache.Get(key); ok {
		return val.(*CheckResult), true
	}
	return nil, false
}

// SetResult 缓存检查结果
func (c *L1Cache) SetResult(key string, result *CheckResult) {
	c.cache.Set(key, result)
}

// Del 删除指定 key 的缓存项（测试或需要绕过陈旧命中时使用）
func (c *L1Cache) Del(key string) {
	if c == nil || c.cache == nil {
		return
	}
	c.cache.Del(key)
}
