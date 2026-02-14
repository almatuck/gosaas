package auth

import (
	"context"
	"net/http"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type GetAuthConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetAuthConfigHandler - Get auth configuration (OAuth providers enabled)
func GetAuthConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &GetAuthConfigLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetAuthConfig()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
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
