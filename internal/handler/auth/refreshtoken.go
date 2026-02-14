package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
)

type RefreshTokenLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// RefreshTokenHandler - Refresh authentication token
func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &RefreshTokenLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.RefreshToken(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest) (resp *types.RefreshTokenResponse, err error) {
	// Use local auth when Levee is disabled
	if l.svcCtx.UseLocal() {
		return l.refreshTokenLocal(req)
	}

	// Use Levee when enabled
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("auth service not configured")
	}

	// Refresh token via Levee SDK
	authResp, err := l.svcCtx.Levee.Auth.RefreshToken(l.ctx, &levee.SDKRefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		slog.Error("token refresh failed", "error", err)
		return nil, err
	}

	// Parse expiry time
	expiresAt, _ := time.Parse(time.RFC3339, authResp.ExpiresAt)

	return &types.RefreshTokenResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    expiresAt.UnixMilli(),
	}, nil
}

// refreshTokenLocal handles token refresh with local SQLite auth
func (l *RefreshTokenLogic) refreshTokenLocal(req *types.RefreshTokenRequest) (*types.RefreshTokenResponse, error) {
	if l.svcCtx.Auth == nil {
		return nil, fmt.Errorf("local auth service not configured")
	}

	authResp, err := l.svcCtx.Auth.RefreshToken(l.ctx, req.RefreshToken)
	if err != nil {
		slog.Error("token refresh failed", "error", err)
		return nil, err
	}

	return &types.RefreshTokenResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    authResp.ExpiresAt.UnixMilli(),
	}, nil
}
