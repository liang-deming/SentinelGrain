package quota

import (
	"context"
	"time"

	"SentinelGrain/api/internal/svc"
	"SentinelGrain/api/internal/types"
	"SentinelGrain/common/quota"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveQuotaLogic {
	return &SaveQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveQuotaLogic) SaveQuota(req *types.SaveQuotaReq) (resp *types.SaveQuotaResp, err error) {
	if req == nil {
		return &types.SaveQuotaResp{Code: 400, Msg: "empty body"}, nil
	}
	rule := &quota.Rule{
		AppId:     req.Rule.AppId,
		Resource:  req.Rule.Resource,
		Threshold: req.Rule.Threshold,
		Period:    req.Rule.Period,
	}
	timeout := time.Duration(l.svcCtx.Config.RedisCommandTimeout) * time.Millisecond
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	if err := l.svcCtx.QuotaRepo.Save(ctx, rule); err != nil {
		l.Infow("save quota validation failed", logx.Field("err", err.Error()))
		return &types.SaveQuotaResp{Code: 400, Msg: err.Error()}, nil
	}
	return &types.SaveQuotaResp{Code: 0, Msg: "ok"}, nil
}
