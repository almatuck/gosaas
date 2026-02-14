package subscription

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type GetUsageStatsLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *GetUsageStatsLogic) GetUsageStats() (resp *types.GetUsageStatsResponse, err error) {
	// Usage tracking would typically come from your application's database
	// or from Levee's metered billing. For the boilerplate, return placeholder data.
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	return &types.GetUsageStatsResponse{
		Stats: types.UsageStats{
			PeriodStart: periodStart.Format(time.RFC3339),
			PeriodEnd:   periodEnd.Format(time.RFC3339),
			Meters: map[string]int{
				"api_calls": 0,
				"storage":   0,
			},
		},
	}, nil
}

// GetUsageStatsHandler handles requests to get the current usage statistics.
func GetUsageStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &GetUsageStatsLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetUsageStats()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
