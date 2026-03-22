// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package quota

import (
	"net/http"

	"SentinelGrain/api/internal/logic/quota"
	"SentinelGrain/api/internal/svc"
	"SentinelGrain/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询配额规则列表
func ListQuotaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListQuotaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := quota.NewListQuotaLogic(r.Context(), svcCtx)
		resp, err := l.ListQuota(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
