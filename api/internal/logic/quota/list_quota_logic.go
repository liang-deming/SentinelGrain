// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package quota

import (
	"context"

	"SentinelGrain/api/internal/svc"
	"SentinelGrain/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询配额规则列表
func NewListQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuotaLogic {
	return &ListQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListQuotaLogic) ListQuota(req *types.ListQuotaReq) (resp *types.ListQuotaResp, err error) {
	// todo: add your logic here and delete this line

	return
}
