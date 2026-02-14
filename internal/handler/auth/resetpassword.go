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

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ResetPasswordHandler - Reset password with token
func ResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetPasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &ResetPasswordLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ResetPassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
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
		l.Errorf("Reset password failed: %v", err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Password has been reset successfully.",
	}, nil
}
