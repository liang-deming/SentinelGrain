// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package quota

import (
	"context"

	"SentinelGrain/api/internal/svc"
	"SentinelGrain/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存配额规则
func NewSaveQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveQuotaLogic {
	return &SaveQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveQuotaLogic) SaveQuota(req *types.SaveQuotaReq) (resp *types.SaveQuotaResp, err error) {
	// todo: add your logic here and delete this line

	return
}
