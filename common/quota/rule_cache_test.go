package quota

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 验收 P3：启动全量加载 + 定时 Refresh 语义等价（此处用连续 Refresh 模拟 Admin 写入后下一轮加载）
func TestRuleCache_StartupAndHotReload(t *testing.T) {
	s := miniredis.RunT(t)
	defer s.Close()

	rdb := redis.New(s.Addr())
	rc := NewRuleCache(rdb, 2*time.Second)
	ctx := context.Background()

	assert.NoError(t, rc.Refresh(ctx))
	_, _, ok := rc.Get("app1", "api1")
	assert.False(t, ok)

	repo := NewRedisRepo(rdb)
	assert.NoError(t, repo.Save(ctx, &Rule{AppId: "app1", Resource: "api1", Threshold: 10, Period: 1}))
	assert.NoError(t, rc.Refresh(ctx))

	th, p, ok := rc.Get("app1", "api1")
	assert.True(t, ok)
	assert.Equal(t, int64(10), th)
	assert.Equal(t, int64(1), p)
	assert.GreaterOrEqual(t, rc.CacheVersion(), int64(1))

	assert.NoError(t, repo.Save(ctx, &Rule{AppId: "app1", Resource: "api1", Threshold: 99, Period: 2}))
	assert.NoError(t, rc.Refresh(ctx))

	th, p, ok = rc.Get("app1", "api1")
	assert.True(t, ok)
	assert.Equal(t, int64(99), th)
	assert.Equal(t, int64(2), p)
}
