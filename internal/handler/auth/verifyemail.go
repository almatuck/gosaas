package auth

import (
	"context"
	"fmt"
	"net/http"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type VerifyEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// VerifyEmailHandler - Verify email address with token
func VerifyEmailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EmailVerificationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &VerifyEmailLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.VerifyEmail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func (l *VerifyEmailLogic) VerifyEmail(req *types.EmailVerificationRequest) (resp *types.MessageResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Verify email via Levee SDK
	_, err = l.svcCtx.Levee.Auth.VerifyEmail(l.ctx, &levee.SDKVerifyEmailRequest{
		Token: req.Token,
	})
	if err != nil {
		l.Errorf("Email verification failed: %v", err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Email verified successfully.",
	}, nil
}
