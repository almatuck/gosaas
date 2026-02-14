package oauth

import (
	"context"
	"fmt"
	"net/http"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type OAuthCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// OAuthCallback is deprecated - OAuth callbacks are handled directly at /oauth/{provider}/callback
// This endpoint exists for API compatibility but should not be called directly.
func (l *OAuthCallbackLogic) OAuthCallback(req *types.OAuthLoginRequest) (resp *types.OAuthLoginResponse, err error) {
	return nil, fmt.Errorf("OAuth callbacks should use /oauth/%s/callback (browser redirect), not the API endpoint", req.Provider)
}

func OAuthCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OAuthLoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &OAuthCallbackLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.OAuthCallback(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
