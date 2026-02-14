package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
)

type DeleteAccountLogic struct {
	logger *slog.Logger
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
		slog.Error("Failed to get email from context", "error", err)
		return nil, err
	}

	// Verify password by attempting login via Levee SDK
	_, err = l.svcCtx.Levee.Auth.Login(l.ctx, &levee.SDKLoginRequest{
		Email:    email,
		Password: req.Password,
	})
	if err != nil {
		slog.Error("Password verification failed for delete account", "error", err)
		return nil, errors.New("invalid password")
	}

	// Get customer ID
	customer, err := l.svcCtx.Levee.Customers.GetCustomerByEmail(l.ctx, email)
	if err != nil {
		slog.Error("Failed to get customer", "email", email, "error", err)
		return nil, err
	}

	// Delete customer via Levee SDK
	_, err = l.svcCtx.Levee.Customers.DeleteCustomer(l.ctx, customer.ID)
	if err != nil {
		slog.Error("Failed to delete customer", "email", email, "error", err)
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
			httpx.ErrorResponse(w, err)
			return
		}

		l := &DeleteAccountLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.DeleteAccount(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
