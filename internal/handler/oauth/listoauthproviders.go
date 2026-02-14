package oauth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type ListOAuthProvidersLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *ListOAuthProvidersLogic) ListOAuthProviders() (resp *types.ListOAuthProvidersResponse, err error) {
	if !l.svcCtx.Config.IsOAuthEnabled() {
		return nil, fmt.Errorf("OAuth feature is not enabled")
	}

	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("OAuth not available in this mode")
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// Get user's OAuth connections
	connections, err := l.svcCtx.DB.Queries.ListUserOAuthConnections(l.ctx, userID.String())
	if err != nil {
		slog.Error("Failed to list OAuth connections", "error", err)
		return nil, err
	}

	// Build map of connected providers
	connectedProviders := make(map[string]string)
	for _, conn := range connections {
		connectedProviders[conn.Provider] = conn.Email.String
	}

	// Build provider list
	var providers []types.OAuthProvider

	// Google
	if l.svcCtx.Config.IsGoogleOAuthEnabled() {
		email := connectedProviders["google"]
		providers = append(providers, types.OAuthProvider{
			Name:      "google",
			Connected: email != "",
			Email:     email,
		})
	}

	// GitHub
	if l.svcCtx.Config.IsGitHubOAuthEnabled() {
		email := connectedProviders["github"]
		providers = append(providers, types.OAuthProvider{
			Name:      "github",
			Connected: email != "",
			Email:     email,
		})
	}

	return &types.ListOAuthProvidersResponse{
		Providers: providers,
	}, nil
}

func ListOAuthProvidersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &ListOAuthProvidersLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ListOAuthProviders()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
