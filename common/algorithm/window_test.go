package algorithm

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"testing"
	"time"
)

func TestEvalSlideWindow(t *testing.T) {
	// 创建mock Redis服务器
	s := miniredis.RunT(t)
	defer s.Close()

	// 创建Redis客户端
	client := redis.New(s.Addr())

	// 测试参数
	ctx := context.Background()
	key := "test:limiter"
	now := time.Now().UnixMilli()
	window := int64(1000) // 1秒窗口
	threshold := int64(5) // 限制5个请求
	cost := int64(1)

	tests := []struct {
		name           string
		key           string
		nowMs         int64
		windowMs      int64
		threshold     int64
		cost          int64
		wantAllowed   bool
		wantRemaining int64
		wantErr       bool
	}{
		{
			name:           "First request should pass",
			key:           key,
			nowMs:         now,
			windowMs:      window,
			threshold:     threshold,
			cost:          cost,
			wantAllowed:   true,
			wantRemaining: threshold - cost,
			wantErr:       false,
		},
		{
			name:           "Invalid window size",
			key:           key,
			nowMs:         now,
			windowMs:      0,
			threshold:     threshold,
			cost:          cost,
			wantAllowed:   false,
			wantRemaining: 0,
			wantErr:       true,
		},
		{
			name:           "Invalid threshold",
			key:           key,
			nowMs:         now,
			windowMs:      window,
			threshold:     0,
			cost:          cost,
			wantAllowed:   false,
			wantRemaining: 0,
			wantErr:       true,
		},
		{
			name:           "Invalid cost",
			key:           key,
			nowMs:         now,
			windowMs:      window,
			threshold:     threshold,
			cost:          0,
			wantAllowed:   false,
			wantRemaining: 0,
			wantErr:       true,
		},
		{
			name:           "Empty key",
			key:           "",
			nowMs:         now,
			windowMs:      window,
			threshold:     threshold,
			cost:          cost,
			wantAllowed:   false,
			wantRemaining: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, remaining, err := EvalSlideWindow(ctx, client, tt.key, tt.nowMs, tt.windowMs, tt.threshold, tt.cost)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantAllowed, allowed)
			assert.Equal(t, tt.wantRemaining, remaining)
		})
	}
}

func TestWindowSliding(t *testing.T) {
	// 创建mock Redis服务器
	s := miniredis.RunT(t)
	defer s.Close()

	// 创建Redis客户端
	client := redis.New(s.Addr())

	// 测试参数
	ctx := context.Background()
	key := "test:limiter"
	window := int64(1000) // 1秒窗口
	threshold := int64(2) // 限制2个请求
	cost := int64(1)

	// 第一个时间点：t0
	t0 := time.Now().UnixMilli()
	allowed, remaining, err := EvalSlideWindow(ctx, client, key, t0, window, threshold, cost)
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, threshold-cost, remaining)

	// 同一个窗口内：t0 + 100ms
	allowed, remaining, err = EvalSlideWindow(ctx, client, key, t0+100, window, threshold, cost)
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, threshold-cost*2, remaining)

	// 同一个窗口内，超过阈值：t0 + 200ms
	allowed, remaining, err = EvalSlideWindow(ctx, client, key, t0+200, window, threshold, cost)
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, int64(0), remaining)

	// 下一个窗口：t0 + window + 100ms
	allowed, remaining, err = EvalSlideWindow(ctx, client, key, t0+window+100, window, threshold, cost)
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, threshold-cost, remaining)
}

func TestPreloadScript(t *testing.T) {
	// 创建mock Redis服务器
	s := miniredis.RunT(t)
	defer s.Close()

	// 创建Redis客户端
	client := redis.New(s.Addr())

	// 测试脚本预加载
	sha, err := PreloadScript(client)
	assert.NoError(t, err)
	assert.NotEmpty(t, sha)

	// 测试nil客户端
	_, err = PreloadScript(nil)
	assert.Error(t, err)
}
