// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package subscription

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gosaas/internal/logic/subscription"
	"gosaas/internal/svc"
)

// Create billing portal session for subscription management
func CreateBillingPortalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := subscription.NewCreateBillingPortalLogic(r.Context(), svcCtx)
		resp, err := l.CreateBillingPortal()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
