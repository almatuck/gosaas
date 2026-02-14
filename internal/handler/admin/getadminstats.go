package admin

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type GetAdminStatsLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *GetAdminStatsLogic) GetAdminStats() (resp *types.AdminStatsResponse, err error) {
	// Only support local/standalone mode for admin stats
	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("admin stats only available in standalone mode")
	}

	// Calculate time boundaries
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 0, 0, 0, 0, now.Location()).Unix()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Unix()

	// Get total users
	totalUsers, err := l.svcCtx.DB.CountUsers(l.ctx)
	if err != nil {
		slog.Error("Failed to count users", "error", err)
		totalUsers = 0
	}

	// Get new users today
	newUsersToday, err := l.svcCtx.DB.CountUsersCreatedAfter(l.ctx, startOfDay)
	if err != nil {
		slog.Error("Failed to count users today", "error", err)
		newUsersToday = 0
	}

	// Get new users this week
	newUsersThisWeek, err := l.svcCtx.DB.CountUsersCreatedAfter(l.ctx, startOfWeek)
	if err != nil {
		slog.Error("Failed to count users this week", "error", err)
		newUsersThisWeek = 0
	}

	// Get new users this month
	newUsersThisMonth, err := l.svcCtx.DB.CountUsersCreatedAfter(l.ctx, startOfMonth)
	if err != nil {
		slog.Error("Failed to count users this month", "error", err)
		newUsersThisMonth = 0
	}

	// Get active subscriptions (paid plans only)
	activeSubscriptions, err := l.svcCtx.DB.CountActiveSubscriptions(l.ctx)
	if err != nil {
		slog.Error("Failed to count active subscriptions", "error", err)
		activeSubscriptions = 0
	}

	// Get trial subscriptions
	trialSubscriptions, err := l.svcCtx.DB.CountTrialSubscriptions(l.ctx)
	if err != nil {
		slog.Error("Failed to count trial subscriptions", "error", err)
		trialSubscriptions = 0
	}

	// Monthly revenue would require Stripe API integration
	// For now, return 0 (can be enhanced later)
	monthlyRevenue := 0

	return &types.AdminStatsResponse{
		TotalUsers:          int(totalUsers),
		ActiveSubscriptions: int(activeSubscriptions),
		TrialSubscriptions:  int(trialSubscriptions),
		MonthlyRevenue:      monthlyRevenue,
		Currency:            "usd",
		NewUsersToday:       int(newUsersToday),
		NewUsersThisWeek:    int(newUsersThisWeek),
		NewUsersThisMonth:   int(newUsersThisMonth),
	}, nil
}

func GetAdminStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &GetAdminStatsLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetAdminStats()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
