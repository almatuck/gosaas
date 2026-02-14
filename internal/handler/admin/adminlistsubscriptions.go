package admin

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type AdminListSubscriptionsLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *AdminListSubscriptionsLogic) AdminListSubscriptions(req *types.AdminListSubscriptionsRequest) (resp *types.AdminListSubscriptionsResponse, err error) {
	// Only support local/standalone mode
	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("admin subscriptions list only available in standalone mode")
	}

	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	statusFilter := req.Status

	// Get total count (with filter)
	totalCount, err := l.svcCtx.DB.CountSubscriptionsFiltered(l.ctx, statusFilter)
	if err != nil {
		slog.Error("Failed to count subscriptions", "error", err)
		return nil, fmt.Errorf("failed to count subscriptions")
	}

	// Get subscriptions
	offset := int64((page - 1) * pageSize)
	subscriptions, err := l.svcCtx.DB.ListSubscriptionsPaginated(l.ctx, db.ListSubscriptionsPaginatedParams{
		StatusFilter: statusFilter,
		PageOffset:   offset,
		PageSize:     int64(pageSize),
	})
	if err != nil {
		slog.Error("Failed to list subscriptions", "error", err)
		return nil, fmt.Errorf("failed to list subscriptions")
	}

	// Convert to response type
	adminSubs := make([]types.AdminSubscription, 0, len(subscriptions))
	for _, s := range subscriptions {
		adminSubs = append(adminSubs, types.AdminSubscription{
			Id:        s.ID,
			UserId:    s.UserID,
			UserEmail: s.UserEmail,
			PlanName:  s.PlanID,
			Status:    s.Status,
			Amount:    0,   // Would need to look up from config/Stripe
			Currency:  "usd",
			Interval:  "month", // Default, would need to look up
			CreatedAt: time.Unix(s.CreatedAt, 0).Format(time.RFC3339),
		})
	}

	return &types.AdminListSubscriptionsResponse{
		Subscriptions: adminSubs,
		TotalCount:    int(totalCount),
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

func AdminListSubscriptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminListSubscriptionsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &AdminListSubscriptionsLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.AdminListSubscriptions(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
