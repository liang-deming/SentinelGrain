package logic

import (
	"context"
	"fmt"
	"sync"
	"time"

	"SentinelGrain/common/algorithm"
	"SentinelGrain/common/cache"
	"SentinelGrain/common/errors"
	"SentinelGrain/rpc/internal/svc"
	"SentinelGrain/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

const (
	// MVP选项：无规则时默认拒绝
	DefaultReject = true
)

var (
	// 对象池，用于复用CheckResult
	resultPool = sync.Pool{
		New: func() interface{} {
			return &cache.CheckResult{}
		},
	}
)

type CheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLogic {
	return &CheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckLogic) Check(in *pb.CheckRequest) (*pb.CheckResponse, error) {
	startTime := time.Now()

	// 构建限流Key；拼入配额版本号，Admin 更新阈值后 L1 不会长期命中旧结果
	ver := l.svcCtx.QuotaRules.CacheVersion()
	key := fmt.Sprintf("%s:%s:%s:v%d", in.AppId, in.Resource, in.Dimension, ver)

	// 先查L1缓存
	if l.svcCtx.L1Cache != nil {
		if result, ok := l.svcCtx.L1Cache.GetResult(key); ok {
			// 记录L1命中日志
			l.Infow("L1 cache hit",
				logx.Field("app_id", in.AppId),
				logx.Field("resource", in.Resource),
				logx.Field("allowed", result.Allowed),
				logx.Field("latency_ms", time.Since(startTime).Milliseconds()),
			)

			// 异步记录指标
			threading.GoSafe(func() {
				l1HitTotal.Inc(in.AppId, in.Resource)
				requestTotal.Inc(in.AppId, in.Resource)
				if !result.Allowed {
					rejectTotal.Inc(in.AppId, in.Resource)
				}
				checkLatency.WithLabelValues(in.AppId, in.Resource).
					Observe(time.Since(startTime).Seconds())
			})

			return &pb.CheckResponse{
				Allowed:   result.Allowed,
				Remaining: result.Remaining,
				Reason:    map[bool]string{true: "", false: errors.RateLimitExceeded}[result.Allowed],
			}, nil
		}
	}

	// 设置Redis操作超时
	redisCmdTimeout := time.Duration(l.svcCtx.Config.CommandTimeout) * time.Millisecond
	redisCtx, cancel := context.WithTimeout(l.ctx, redisCmdTimeout)
	defer cancel()

	// 查找限流规则（与 Admin 写入 Redis 的规则一致，经 QuotaRules 全量/定时加载）
	var threshold, period int64
	if th, p, ok := l.svcCtx.QuotaRules.Get(in.AppId, in.Resource); ok {
		threshold, period = th, p
	} else {
		// MVP选项：无规则时默认拒绝
		allowed := !DefaultReject
		l.Infow("No rule found",
			logx.Field("app_id", in.AppId),
			logx.Field("resource", in.Resource),
			logx.Field("allowed", allowed),
			logx.Field("latency_ms", time.Since(startTime).Milliseconds()),
		)

		// 异步记录指标
		threading.GoSafe(func() {
			requestTotal.Inc(in.AppId, in.Resource)
			if !allowed {
				rejectTotal.Inc(in.AppId, in.Resource)
			}
			checkLatency.WithLabelValues(in.AppId, in.Resource).
				Observe(time.Since(startTime).Seconds())
		})

		return &pb.CheckResponse{
			Allowed: allowed,
			Reason:  map[bool]string{true: "", false: errors.RuleNotFound}[allowed],
		}, nil
	}

	// 执行Redis滑动窗口判定
	now := time.Now().UnixMilli()
	windowMs := period * 1000 // 转换为毫秒

	allowed, remaining, err := algorithm.EvalSlideWindow(
		redisCtx,
		l.svcCtx.BizRedis,
		key,
		now,
		int64(windowMs),
		threshold,
		in.Cost,
	)

	if err != nil {
		l.Errorw("Failed to check quota",
			logx.Field("error", err),
			logx.Field("app_id", in.AppId),
			logx.Field("resource", in.Resource),
		)
		return &pb.CheckResponse{
			Allowed: false,
			Reason:  errors.InternalError,
		}, nil
	}

	// 缓存结果到L1 (仅缓存成功的Redis响应)
	if l.svcCtx.L1Cache != nil {
		result := resultPool.Get().(*cache.CheckResult)
		result.Allowed = allowed
		result.Remaining = remaining
		l.svcCtx.L1Cache.SetResult(key, result)
		resultPool.Put(result)
	}

	// 记录结构化日志
	l.Infow("Quota check completed",
		logx.Field("app_id", in.AppId),
		logx.Field("resource", in.Resource),
		logx.Field("allowed", allowed),
		logx.Field("latency_ms", time.Since(startTime).Milliseconds()),
	)

	// 异步记录指标
	threading.GoSafe(func() {
		requestTotal.Inc(in.AppId, in.Resource)
		if !allowed {
			rejectTotal.Inc(in.AppId, in.Resource)
		}
		checkLatency.WithLabelValues(in.AppId, in.Resource).
			Observe(time.Since(startTime).Seconds())
	})

	resp := &pb.CheckResponse{
		Allowed:   allowed,
		Remaining: remaining,
	}
	if !allowed {
		resp.Reason = errors.RateLimitExceeded
	}
	return resp, nil
}
