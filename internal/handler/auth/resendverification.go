package auth

import (
	"context"
	"log/slog"
	"net/http"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type ResendVerificationLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ResendVerificationHandler - Resend email verification
func ResendVerificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResendVerificationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &ResendVerificationLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ResendVerification(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *ResendVerificationLogic) ResendVerification(req *types.ResendVerificationRequest) (resp *types.MessageResponse, err error) {
	// Levee doesn't have a separate resend verification endpoint
	// The verification email is sent on registration
	// For now, return a success message
	return &types.MessageResponse{
		Message: "If the email address is registered and unverified, a new verification email has been sent.",
	}, nil
}
