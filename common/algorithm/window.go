package algorithm

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

//go:embed window.lua
var slideWindowScript string

// GetSlideWindowScript 获取滑动窗口Lua脚本
func GetSlideWindowScript() string {
	return slideWindowScript
}

// EvalSlideWindow 执行滑动窗口限流判定
// key: Redis限流计数器key
// nowMs: 当前时间戳(毫秒)
// windowMs: 窗口大小(毫秒)
// threshold: 阈值
// cost: 本次请求消耗的配额数
func EvalSlideWindow(ctx context.Context, redis *redis.Redis, key string, nowMs, windowMs, threshold, cost int64) (allowed bool, remaining int64, err error) {
	if redis == nil {
		return false, 0, fmt.Errorf("redis client is nil")
	}
	if key == "" {
		return false, 0, fmt.Errorf("key is empty")
	}
	if windowMs <= 0 {
		return false, 0, fmt.Errorf("invalid window size: %d", windowMs)
	}
	if threshold <= 0 {
		return false, 0, fmt.Errorf("invalid threshold: %d", threshold)
	}
	if cost <= 0 {
		return false, 0, fmt.Errorf("invalid cost: %d", cost)
	}

	// 执行Lua脚本
	resp, err := redis.EvalCtx(ctx, slideWindowScript, []string{key}, []interface{}{nowMs, windowMs, threshold, cost})
	if err != nil {
		return false, 0, fmt.Errorf("failed to eval script: %v", err)
	}

	// 解析返回值
	result, ok := resp.([]interface{})
	if !ok || len(result) != 2 {
		return false, 0, fmt.Errorf("invalid script response format")
	}

	// 转换返回值
	allowedInt, ok := result[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("invalid allowed value type")
	}
	remaining, ok = result[1].(int64)
	if !ok {
		return false, 0, fmt.Errorf("invalid remaining value type")
	}

	return allowedInt == 1, remaining, nil
}

// PreloadScript 预加载Lua脚本到Redis并返回SHA
func PreloadScript(redis *redis.Redis) (string, error) {
	if redis == nil {
		return "", fmt.Errorf("redis client is nil")
	}
	return redis.ScriptLoad(slideWindowScript)
}
