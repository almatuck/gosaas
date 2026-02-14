package user

import (
	"context"
	"fmt"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ChangePasswordLogic struct {
	logx.Logger
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
		l.Errorf("Failed to get email from context: %v", err)
		return nil, err
	}

	// Change password via Levee SDK
	_, err = l.svcCtx.Levee.Auth.ChangePassword(l.ctx, &levee.SDKChangePasswordRequest{
		Email:           email,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		l.Errorf("Failed to change password for %s: %v", email, err)
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
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &ChangePasswordLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ChangePassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
