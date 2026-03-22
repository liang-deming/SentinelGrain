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

// 保存配额规则
func SaveQuotaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveQuotaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := quota.NewSaveQuotaLogic(r.Context(), svcCtx)
		resp, err := l.SaveQuota(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
