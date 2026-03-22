package logic

import (
	"context"
	"testing"
	"time"

	"SentinelGrain/common/errors"
	"SentinelGrain/rpc/internal/config"
	"SentinelGrain/rpc/internal/svc"
	"SentinelGrain/rpc/pb"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestCheckLogic_Check(t *testing.T) {
	// 创建mock Redis服务器
	s := miniredis.RunT(t)
	defer s.Close()

	// 创建测试配置
	conf := config.Config{
		CommandTimeout: 100,
		L1Cache: config.L1Config{
			Enabled: true,
			TTL:     500,
		},
		Redis: redis.RedisConf{
			Host: s.Addr(),
			Type: "node",
		},
	}

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(conf)

	// 创建测试用例
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
				Resource: "test_api",
				Dimension: "test_user",
				Cost:     1,
			},
			wantAllowed: true,
			wantReason:  "",
		},
		{
			name: "L1 cache hit should return quickly",
			req: &pb.CheckRequest{
				AppId:     "test_app",
				Resource: "test_api",
				Dimension: "test_user",
				Cost:     1,
			},
			wantAllowed: true,
			wantReason:  "",
		},
		{
			name: "Exceed threshold should be rejected",
			req: &pb.CheckRequest{
				AppId:     "test_app",
				Resource: "test_api",
				Dimension: "test_user",
				Cost:     DefaultThreshold + 1,
			},
			wantAllowed: false,
			wantReason:  errors.RateLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	// 创建mock Redis服务器
	s := miniredis.RunT(t)
	defer s.Close()

	// 配置极短的超时时间
	conf := config.Config{
		CommandTimeout: 1, // 1ms
		Redis: redis.RedisConf{
			Host: s.Addr(),
			Type: "node",
		},
	}

	svcCtx := svc.NewServiceContext(conf)
	logic := NewCheckLogic(context.Background(), svcCtx)

	// 模拟Redis延迟
	s.FastForward(time.Millisecond * 10)

	resp, err := logic.Check(&pb.CheckRequest{
		AppId:     "test_app",
		Resource: "test_api",
		Dimension: "test_user",
		Cost:     1,
	})

	// 验证超时情况下的行为
	assert.NoError(t, err)
	assert.False(t, resp.Allowed)
	assert.Equal(t, errors.InternalError, resp.Reason)
}
