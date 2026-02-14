package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/auth"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type CancelSubscriptionLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *CancelSubscriptionLogic) CancelSubscription(req *types.CancelSubscriptionRequest) (resp *types.CancelSubscriptionResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Get email from JWT context
	email, err := auth.GetEmailFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get email from context", "error", err)
		return nil, err
	}

	// Get current subscriptions from Levee SDK
	subsResp, err := l.svcCtx.Levee.Customers.ListCustomerSubscriptions(l.ctx, email)
	if err != nil {
		slog.Error("Failed to get subscriptions", "email", email, "error", err)
		return nil, err
	}

	if len(subsResp.Subscriptions) == 0 {
		return nil, errors.New("no active subscription found")
	}

	// Cancel the first active subscription via Levee SDK
	sub := subsResp.Subscriptions[0]
	_, err = l.svcCtx.Levee.Billing.CancelSubscription(l.ctx, sub.ID)
	if err != nil {
		slog.Error("Failed to cancel subscription", "email", email, "error", err)
		return nil, err
	}

	now := time.Now().Format(time.RFC3339)
	return &types.CancelSubscriptionResponse{
		Message:     "Subscription cancelled successfully.",
		CancelledAt: now,
		EffectiveAt: sub.CurrentPeriodEnd,
	}, nil
}

// CancelSubscriptionHandler handles requests to cancel the current subscription.
func CancelSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CancelSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &CancelSubscriptionLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.CancelSubscription(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
