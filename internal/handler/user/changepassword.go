package user

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
)

type ChangePasswordLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) (resp *types.MessageResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Get email from JWT context
	email, err := auth.GetEmailFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get email from context", "error", err)
		return nil, err
	}

	// Change password via Levee SDK
	_, err = l.svcCtx.Levee.Auth.ChangePassword(l.ctx, &levee.SDKChangePasswordRequest{
		Email:           email,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		slog.Error("Failed to change password", "email", email, "error", err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Password changed successfully.",
	}, nil
}

// ChangePasswordHandler handles the change password request for authenticated users.
func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &ChangePasswordLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ChangePassword(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
