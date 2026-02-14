package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
)

type ResetPasswordLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ResetPasswordHandler - Reset password with token
func ResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetPasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &ResetPasswordLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ResetPassword(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordRequest) (resp *types.MessageResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Reset password via Levee SDK
	_, err = l.svcCtx.Levee.Auth.ResetPassword(l.ctx, &levee.SDKResetPasswordRequest{
		Token:           req.Token,
		Password:        req.NewPassword,
		ConfirmPassword: req.NewPassword,
	})
	if err != nil {
		slog.Error("reset password failed", "error", err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Password has been reset successfully.",
	}, nil
}
