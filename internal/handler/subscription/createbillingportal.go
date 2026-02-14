package subscription

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

type CreateBillingPortalLogic struct {
	logx.Logger
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
		l.Errorf("Failed to get email from context: %v", err)
		return nil, err
	}

	// Get customer to get Stripe customer ID
	customer, err := l.svcCtx.Levee.Customers.GetCustomerByEmail(l.ctx, email)
	if err != nil {
		l.Errorf("Failed to get customer %s: %v", email, err)
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
		l.Errorf("Failed to get billing portal for %s: %v", email, err)
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
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.CreateBillingPortal()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
