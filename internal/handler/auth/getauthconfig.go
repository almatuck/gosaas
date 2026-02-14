package auth

import (
	"context"
	"log/slog"
	"net/http"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"gosaas/internal/httpx"
)

type GetAuthConfigLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetAuthConfigHandler - Get auth configuration (OAuth providers enabled)
func GetAuthConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &GetAuthConfigLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetAuthConfig()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *GetAuthConfigLogic) GetAuthConfig() (resp *types.AuthConfigResponse, err error) {
	// Return OAuth provider configuration
	// Only return enabled status if OAuth feature is enabled and in local mode
	googleEnabled := false
	githubEnabled := false

	if l.svcCtx.UseLocal() && l.svcCtx.Config.IsOAuthEnabled() {
		googleEnabled = l.svcCtx.Config.IsGoogleOAuthEnabled()
		githubEnabled = l.svcCtx.Config.IsGitHubOAuthEnabled()
	}

	return &types.AuthConfigResponse{
		GoogleEnabled: googleEnabled,
		GitHubEnabled: githubEnabled,
	}, nil
}
