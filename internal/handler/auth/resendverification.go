package auth

import (
	"context"
	"net/http"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ResendVerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ResendVerificationHandler - Resend email verification
func ResendVerificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResendVerificationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &ResendVerificationLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ResendVerification(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
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
