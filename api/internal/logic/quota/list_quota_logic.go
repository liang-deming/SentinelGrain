package quota

import (
	"context"
	"time"

	"SentinelGrain/api/internal/svc"
	"SentinelGrain/api/internal/types"
	"SentinelGrain/common/quota"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuotaLogic {
	return &ListQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListQuotaLogic) ListQuota(req *types.ListQuotaReq) (resp *types.ListQuotaResp, err error) {
	timeout := time.Duration(l.svcCtx.Config.RedisCommandTimeout) * time.Millisecond
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	rules, total, err := l.svcCtx.QuotaRepo.List(ctx, quota.ListQuery{
		AppId:    req.AppId,
		Resource: req.Resource,
		Page:     page,
		Size:     size,
	})
	if err != nil {
		return nil, err
	}
	out := make([]types.QuotaRule, 0, len(rules))
	for i := range rules {
		r := &rules[i]
		out = append(out, types.QuotaRule{
			AppId:      r.AppId,
			Resource:   r.Resource,
			Threshold:  r.Threshold,
			Period:     r.Period,
			UpdateTime: r.UpdateTime,
		})
	}
	return &types.ListQuotaResp{Total: int64(total), Rules: out}, nil
}
