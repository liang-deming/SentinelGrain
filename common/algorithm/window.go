package algorithm

import (
	"context"
	_ "embed"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

//go:embed window.lua
var slideWindowScript string

// GetSlideWindowScript 获取滑动窗口Lua脚本
func GetSlideWindowScript() string {
	return slideWindowScript
}

// EvalSlideWindow 执行滑动窗口限流判定
func EvalSlideWindow(ctx context.Context, redis *redis.Redis, key string, nowMs, windowMs, threshold, cost int64) (allowed bool, remaining int64, err error) {
	// 执行Lua脚本
	resp, err := redis.EvalCtx(ctx, slideWindowScript, []string{key}, []interface{}{nowMs, windowMs, threshold, cost})
	if err != nil {
		return false, 0, err
	}

	// 解析返回值
	result := resp.([]interface{})
	allowed = result[0].(int64) == 1
	remaining = result[1].(int64)

	return allowed, remaining, nil
}
