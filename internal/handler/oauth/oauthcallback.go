package oauth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type OAuthCallbackLogic struct {
	logger *slog.Logger
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
			httpx.ErrorResponse(w, err)
			return
		}

		l := &OAuthCallbackLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.OAuthCallback(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
