package subscription

import (
	"context"
	"log/slog"
	"net/http"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type ListPlansLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *ListPlansLogic) ListPlans() (resp *types.ListPlansResponse, err error) {
	// Plans are typically configured in Levee's product catalog
	// For the boilerplate, return example plans
	// You should update this with your actual plans from Levee
	plans := []types.SubscriptionPlan{
		{
			Id:          "free",
			Name:        "free",
			DisplayName: "Free",
			Description: "Get started with basic features",
			Price:       0,
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"Basic features",
				"Community support",
				"1 project",
			},
		},
		{
			Id:          "pro",
			Name:        "pro",
			DisplayName: "Pro",
			Description: "For professionals and small teams",
			Price:       2900,
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"All Free features",
				"Priority support",
				"Unlimited projects",
				"Advanced analytics",
			},
		},
		{
			Id:          "team",
			Name:        "team",
			DisplayName: "Team",
			Description: "For growing teams",
			Price:       9900,
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"All Pro features",
				"Team collaboration",
				"Admin dashboard",
				"API access",
				"Dedicated support",
			},
		},
	}

	return &types.ListPlansResponse{
		Plans: plans,
	}, nil
}

// ListPlansHandler handles requests to list all available subscription plans.
func ListPlansHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &ListPlansLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ListPlans()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
