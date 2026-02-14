package oauth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type DisconnectOAuthLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *DisconnectOAuthLogic) DisconnectOAuth(req *types.DisconnectOAuthRequest) (resp *types.MessageResponse, err error) {
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

	// Verify user has this OAuth connection
	_, err = l.svcCtx.DB.Queries.GetOAuthConnectionByUserAndProvider(l.ctx, db.GetOAuthConnectionByUserAndProviderParams{
		UserID:   userID.String(),
		Provider: req.Provider,
	})
	if err != nil {
		return nil, fmt.Errorf("OAuth provider %s is not connected", req.Provider)
	}

	// Get user to check if they have a password set
	user, err := l.svcCtx.Auth.GetUserByEmail(l.ctx, "")
	if err == nil && user != nil {
		// Check if user has other login methods
		connections, _ := l.svcCtx.DB.Queries.ListUserOAuthConnections(l.ctx, userID.String())
		hasPassword := user.PasswordHash != ""

		if !hasPassword && len(connections) <= 1 {
			return nil, fmt.Errorf("cannot disconnect your only login method; please set a password first")
		}
	}

	// Delete the OAuth connection
	err = l.svcCtx.DB.Queries.DeleteOAuthConnectionByProvider(l.ctx, db.DeleteOAuthConnectionByProviderParams{
		UserID:   userID.String(),
		Provider: req.Provider,
	})
	if err != nil {
		slog.Error("Failed to disconnect OAuth", "error", err)
		return nil, err
	}

	slog.Info("User disconnected OAuth provider", "userId", userID.String(), "provider", req.Provider)

	return &types.MessageResponse{
		Message: fmt.Sprintf("Successfully disconnected %s", req.Provider),
	}, nil
}

func DisconnectOAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DisconnectOAuthRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &DisconnectOAuthLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.DisconnectOAuth(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
