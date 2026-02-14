package subscription

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

type CreateBillingPortalLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *CreateBillingPortalLogic) CreateBillingPortal() (resp *types.CreateBillingPortalResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Get email from JWT context
	email, err := auth.GetEmailFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get email from context", "error", err)
		return nil, err
	}

	// Get customer to get Stripe customer ID
	customer, err := l.svcCtx.Levee.Customers.GetCustomerByEmail(l.ctx, email)
	if err != nil {
		slog.Error("Failed to get customer", "email", email, "error", err)
		return nil, err
	}

	// Build return URL
	returnURL := l.svcCtx.Config.App.BaseURL + "/account"

	// Get billing portal URL from Levee SDK
	portalResp, err := l.svcCtx.Levee.Billing.GetCustomerPortal(l.ctx, &levee.PortalRequest{
		CustomerID: customer.ID,
		ReturnUrl:  returnURL,
	})
	if err != nil {
		slog.Error("Failed to get billing portal", "email", email, "error", err)
		return nil, err
	}

	return &types.CreateBillingPortalResponse{
		PortalUrl: portalResp.PortalUrl,
	}, nil
}

// CreateBillingPortalHandler handles requests to create a billing portal session.
func CreateBillingPortalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &CreateBillingPortalLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.CreateBillingPortal()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
