package logic

import (
	"context"
	"testing"

	"SentinelGrain/common/errors"
	"SentinelGrain/common/quota"
	"SentinelGrain/rpc/internal/config"
	"SentinelGrain/rpc/internal/svc"
	"SentinelGrain/rpc/pb"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestCheckLogic_Check(t *testing.T) {
	s := miniredis.RunT(t)
	defer s.Close()

	conf := config.Config{
		BizRedis: redis.RedisConf{
			Host: s.Addr(),
			Type: "node",
		},
		CommandTimeout: 100,
		L1Cache: config.L1Config{
			Enabled: true,
			TTL:     500,
		},
	}

	svcCtx := svc.NewServiceContext(conf)

	t.Cleanup(func() {
		s.FlushAll()
		_ = svcCtx.QuotaRules.Refresh(context.Background())
		if svcCtx.L1Cache != nil {
			v := svcCtx.QuotaRules.CacheVersion()
			svcCtx.L1Cache.Del(LimiterKey("test_app", "test_api", "test_user", v))
		}
	})

	t.Run("No rule should be rejected", func(t *testing.T) {
		logic := NewCheckLogic(context.Background(), svcCtx)
		resp, err := logic.Check(&pb.CheckRequest{
			AppId:     "no_rule_app",
			Resource:  "no_rule_api",
			Dimension: "test_user",
			Cost:      1,
		})

		assert.NoError(t, err)
		assert.False(t, resp.Allowed)
		assert.Equal(t, errors.RuleNotFound, resp.Reason)
	})

	repo := quota.NewRedisRepo(svcCtx.BizRedis)
	ctx := context.Background()
	assert.NoError(t, repo.Save(ctx, &quota.Rule{
		AppId: "test_app", Resource: "test_api", Threshold: 5, Period: 1,
	}))
	assert.NoError(t, svcCtx.QuotaRules.Refresh(ctx))

	tests := []struct {
		name        string
		req         *pb.CheckRequest
		wantAllowed bool
		wantReason  string
	}{
		{
			name: "First request should pass",
			req: &pb.CheckRequest{
				AppId:     "test_app",
				Resource:  "test_api",
				Dimension: "test_user",
				Cost:      1,
			},
			wantAllowed: true,
			wantReason:  "",
		},
		{
			name: "L1 cache hit should return quickly",
			req: &pb.CheckRequest{
				AppId:     "test_app",
				Resource:  "test_api",
				Dimension: "test_user",
				Cost:      1,
			},
			wantAllowed: true,
			wantReason:  "",
		},
		{
			name: "Exceed threshold should be rejected",
			req: &pb.CheckRequest{
				AppId:     "test_app",
				Resource:  "test_api",
				Dimension: "test_user",
				Cost:      6,
			},
			wantAllowed: false,
			wantReason:  errors.RateLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Exceed threshold should be rejected" && svcCtx.L1Cache != nil {
				svcCtx.L1Cache.Del(LimiterKey(tt.req.AppId, tt.req.Resource, tt.req.Dimension, svcCtx.QuotaRules.CacheVersion()))
			}
			logic := NewCheckLogic(context.Background(), svcCtx)
			resp, err := logic.Check(tt.req)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantAllowed, resp.Allowed)
			if !tt.wantAllowed {
				assert.Equal(t, tt.wantReason, resp.Reason)
			}
		})
	}
}

func TestCheckLogic_RedisTimeout(t *testing.T) {
	s := miniredis.RunT(t)
	defer s.Close()

	conf := config.Config{
		BizRedis: redis.RedisConf{
			Host: s.Addr(),
			Type: "node",
		},
		CommandTimeout: 3000,
	}

	svcCtx := svc.NewServiceContext(conf)
	repo := quota.NewRedisRepo(svcCtx.BizRedis)
	assert.NoError(t, repo.Save(context.Background(), &quota.Rule{
		AppId: "test_app", Resource: "test_api", Threshold: 5, Period: 1,
	}))
	assert.NoError(t, svcCtx.QuotaRules.Refresh(context.Background()))

	logic := NewCheckLogic(context.Background(), svcCtx)

	s.SetError("LOADING Redis is loading the dataset in memory")

	resp, err := logic.Check(&pb.CheckRequest{
		AppId:     "test_app",
		Resource:  "test_api",
		Dimension: "test_user",
		Cost:      1,
	})

	assert.NoError(t, err)
	assert.False(t, resp.Allowed)
	assert.Equal(t, errors.InternalError, resp.Reason)
}
