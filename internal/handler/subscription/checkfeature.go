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
)

type CheckFeatureLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Feature access rules - customize based on your plans
var planFeatures = map[string][]string{
	"free": {"basic", "community_support"},
	"pro":  {"basic", "community_support", "priority_support", "analytics", "unlimited_projects"},
	"team": {"basic", "community_support", "priority_support", "analytics", "unlimited_projects", "team", "api", "admin"},
}

func (l *CheckFeatureLogic) CheckFeature(req *types.CheckFeatureRequest) (resp *types.CheckFeatureResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Get email from JWT context
	email, err := auth.GetEmailFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get email from context", "error", err)
		return nil, err
	}

	// Get current subscription from Levee SDK
	subsResp, err := l.svcCtx.Levee.Customers.ListCustomerSubscriptions(l.ctx, email)
	if err != nil {
		slog.Error("Failed to get subscriptions", "email", email, "error", err)
		return nil, err
	}

	subs := subsResp.Subscriptions

	// Default to free plan
	planName := "free"
	if len(subs) > 0 {
		planName = subs[0].ProductName
	}

	// Check if feature is available for this plan
	features, ok := planFeatures[planName]
	if !ok {
		features = planFeatures["free"]
	}

	hasAccess := false
	for _, f := range features {
		if f == req.Feature {
			hasAccess = true
			break
		}
	}

	message := ""
	if !hasAccess {
		message = "Upgrade your plan to access this feature."
	}

	return &types.CheckFeatureResponse{
		HasAccess: hasAccess,
		Feature:   req.Feature,
		PlanName:  planName,
		Message:   message,
	}, nil
}

// CheckFeatureHandler handles requests to check if a user has access to a feature.
func CheckFeatureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckFeatureRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &CheckFeatureLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.CheckFeature(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
