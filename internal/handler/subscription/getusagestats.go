package subscription

import (
	"context"
	"net/http"
	"time"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type GetUsageStatsLogic struct {
	logx.Logger
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
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetUsageStats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
