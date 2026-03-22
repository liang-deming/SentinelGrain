package quota

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ruleEntry struct {
	Threshold int64
	Period    int64
}

// RuleCache RPC 侧内存规则表：启动全量加载 + 定时 Refresh；版本号参与 L1 key，避免 Admin 更新后脏读
type RuleCache struct {
	rdb *redis.Redis
	mu  sync.RWMutex
	// key: appId:resource
	rules   map[string]*ruleEntry
	version int64
	// cmdTimeout 单次 Redis 操作上限（Refresh 内使用）
	cmdTimeout time.Duration
}

// NewRuleCache 由 rpc svc 注入；cmdTimeout 与 CommandTimeout 对齐
func NewRuleCache(rdb *redis.Redis, cmdTimeout time.Duration) *RuleCache {
	return &RuleCache{
		rdb:        rdb,
		rules:      make(map[string]*ruleEntry),
		cmdTimeout: cmdTimeout,
	}
}

// Get 返回 threshold、period；无规则时 ok=false
func (c *RuleCache) Get(appId, resource string) (threshold, period int64, ok bool) {
	k := fmt.Sprintf("%s:%s", appId, resource)
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.rules[k]
	if !ok || e == nil {
		return 0, 0, false
	}
	return e.Threshold, e.Period, true
}

// CacheVersion 当前已加载的配置版本（来自 Redis KeyVersion），用于拼 L1 key
func (c *RuleCache) CacheVersion() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.version
}

// Refresh 从 Redis 全量重载规则与版本号
func (c *RuleCache) Refresh(ctx context.Context) error {
	repo := NewRedisRepo(c.rdb)
	rules, _, err := repo.List(ctx, ListQuery{Page: 1, Size: 100000})
	if err != nil {
		return err
	}
	verStr, err := c.rdb.GetCtx(ctx, KeyVersion())
	if err != nil {
		return err
	}
	ver, _ := strconv.ParseInt(verStr, 10, 64)

	m := make(map[string]*ruleEntry, len(rules))
	for i := range rules {
		r := &rules[i]
		k := fmt.Sprintf("%s:%s", r.AppId, r.Resource)
		m[k] = &ruleEntry{Threshold: r.Threshold, Period: r.Period}
	}
	c.mu.Lock()
	c.rules = m
	c.version = ver
	c.mu.Unlock()
	return nil
}

// StartPeriodicRefresh 定时热更（混合模式：启动全量加载 + 周期刷新）。
// TODO: 可接入 Redis Pub/Sub 或 Admin 主动通知，使 Save 后立即 Refresh，降低传播延迟。
func (c *RuleCache) StartPeriodicRefresh(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), c.cmdTimeout*10)
			_ = c.Refresh(ctx)
			cancel()
		}
	}()
}
