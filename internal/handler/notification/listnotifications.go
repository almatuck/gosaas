package notification

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type ListNotificationsLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *ListNotificationsLogic) ListNotifications(req *types.ListNotificationsRequest) (resp *types.ListNotificationsResponse, err error) {
	// Check if notifications are enabled
	if !l.svcCtx.Config.IsNotificationsEnabled() {
		return &types.ListNotificationsResponse{
			Notifications: []types.Notification{},
			UnreadCount:   0,
			TotalCount:    0,
		}, nil
	}

	if !l.svcCtx.UseLocal() {
		return &types.ListNotificationsResponse{
			Notifications: []types.Notification{},
			UnreadCount:   0,
			TotalCount:    0,
		}, nil
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// Set defaults
	pageSize := int64(20)
	if req.PageSize > 0 && req.PageSize <= 100 {
		pageSize = int64(req.PageSize)
	}
	page := 1
	if req.Page > 0 {
		page = req.Page
	}
	offset := int64((page - 1)) * pageSize

	var notifications []db.Notification
	if req.Unread {
		notifications, err = l.svcCtx.DB.Queries.ListUnreadNotifications(l.ctx, db.ListUnreadNotificationsParams{
			UserID:   userID.String(),
			PageSize: pageSize,
		})
	} else {
		notifications, err = l.svcCtx.DB.Queries.ListUserNotifications(l.ctx, db.ListUserNotificationsParams{
			UserID:     userID.String(),
			PageOffset: offset,
			PageSize:   pageSize,
		})
	}
	if err != nil {
		slog.Error("Failed to list notifications", "error", err)
		return nil, err
	}

	// Get unread count
	unreadCount, err := l.svcCtx.DB.Queries.CountUnreadNotifications(l.ctx, userID.String())
	if err != nil {
		slog.Error("Failed to count unread notifications", "error", err)
		return nil, err
	}

	// Convert to response type
	result := make([]types.Notification, len(notifications))
	for i, n := range notifications {
		result[i] = types.Notification{
			Id:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body.String,
			ActionUrl: n.ActionUrl.String,
			Icon:      n.Icon.String,
			CreatedAt: time.Unix(n.CreatedAt, 0).Format(time.RFC3339),
		}
		if n.ReadAt.Valid {
			result[i].ReadAt = time.Unix(n.ReadAt.Int64, 0).Format(time.RFC3339)
		}
	}

	return &types.ListNotificationsResponse{
		Notifications: result,
		UnreadCount:   int(unreadCount),
		TotalCount:    len(result),
	}, nil
}

func ListNotificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListNotificationsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &ListNotificationsLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ListNotifications(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
