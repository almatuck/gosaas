package subscription

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gosaas/internal/logic/subscription"
	"gosaas/internal/svc"
)

// Get current usage statistics
func GetUsageStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := subscription.NewGetUsageStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetUsageStats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
