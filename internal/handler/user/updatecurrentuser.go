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

type UpdateCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *UpdateCurrentUserLogic) UpdateCurrentUser(req *types.UpdateUserRequest) (resp *types.GetUserResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Get email from JWT context
	email, err := auth.GetEmailFromContext(l.ctx)
	if err != nil {
		l.Errorf("Failed to get email from context: %v", err)
		return nil, err
	}

	// Get customer ID first
	customer, err := l.svcCtx.Levee.Customers.GetCustomerByEmail(l.ctx, email)
	if err != nil {
		l.Errorf("Failed to get customer %s: %v", email, err)
		return nil, err
	}

	// Update customer via Levee SDK
	customerInfo, err := l.svcCtx.Levee.Customers.UpdateCustomer(l.ctx, customer.ID, &levee.SDKUpdateCustomerRequest{
		Name: req.Name,
	})
	if err != nil {
		l.Errorf("Failed to update customer %s: %v", email, err)
		return nil, err
	}

	return &types.GetUserResponse{
		User: types.User{
			Id:        customerInfo.ID,
			Email:     customerInfo.Email,
			Name:      customerInfo.Name,
			CreatedAt: customerInfo.CreatedAt,
		},
	}, nil
}

// UpdateCurrentUserHandler handles the update current user profile request.
func UpdateCurrentUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &UpdateCurrentUserLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.UpdateCurrentUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
