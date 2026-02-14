package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type DeleteAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *DeleteAccountLogic) DeleteAccount(req *types.DeleteAccountRequest) (resp *types.MessageResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Get email from JWT context
	email, err := auth.GetEmailFromContext(l.ctx)
	if err != nil {
		l.Errorf("Failed to get email from context: %v", err)
		return nil, err
	}

	// Verify password by attempting login via Levee SDK
	_, err = l.svcCtx.Levee.Auth.Login(l.ctx, &levee.SDKLoginRequest{
		Email:    email,
		Password: req.Password,
	})
	if err != nil {
		l.Errorf("Password verification failed for delete account: %v", err)
		return nil, errors.New("invalid password")
	}

	// Get customer ID
	customer, err := l.svcCtx.Levee.Customers.GetCustomerByEmail(l.ctx, email)
	if err != nil {
		l.Errorf("Failed to get customer %s: %v", email, err)
		return nil, err
	}

	// Delete customer via Levee SDK
	_, err = l.svcCtx.Levee.Customers.DeleteCustomer(l.ctx, customer.ID)
	if err != nil {
		l.Errorf("Failed to delete customer %s: %v", email, err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Account deleted successfully.",
	}, nil
}

// DeleteAccountHandler handles the delete account request for the current user.
func DeleteAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAccountRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &DeleteAccountLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.DeleteAccount(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
