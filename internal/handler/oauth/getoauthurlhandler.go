// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package oauth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gosaas/internal/logic/oauth"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

// Get OAuth authorization URL
func GetOAuthUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOAuthUrlRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := oauth.NewGetOAuthUrlLogic(r.Context(), svcCtx)
		resp, err := l.GetOAuthUrl(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
