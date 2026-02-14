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

type AdminListUsersLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *AdminListUsersLogic) AdminListUsers(req *types.AdminListUsersRequest) (resp *types.AdminListUsersResponse, err error) {
	// Only support local/standalone mode
	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("admin users list only available in standalone mode")
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

	search := req.Search

	// Get total count
	totalCount, err := l.svcCtx.DB.CountUsers(l.ctx)
	if err != nil {
		slog.Error("Failed to count users", "error", err)
		return nil, fmt.Errorf("failed to count users")
	}

	// Get users
	offset := int64((page - 1) * pageSize)
	users, err := l.svcCtx.DB.ListUsersPaginated(l.ctx, db.ListUsersPaginatedParams{
		Search:     search,
		PageOffset: offset,
		PageSize:   int64(pageSize),
	})
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		return nil, fmt.Errorf("failed to list users")
	}

	// Convert to response type
	adminUsers := make([]types.AdminUser, 0, len(users))
	for _, u := range users {
		// Get user's subscription to determine plan
		plan := "free"
		status := "active"
		sub, err := l.svcCtx.DB.GetSubscriptionByUserID(l.ctx, u.ID)
		if err == nil {
			plan = sub.PlanID
			status = sub.Status
		}

		adminUsers = append(adminUsers, types.AdminUser{
			Id:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Plan:      plan,
			Status:    status,
			CreatedAt: time.Unix(u.CreatedAt, 0).Format(time.RFC3339),
		})
	}

	return &types.AdminListUsersResponse{
		Users:      adminUsers,
		TotalCount: int(totalCount),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func AdminListUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminListUsersRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &AdminListUsersLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.AdminListUsers(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
