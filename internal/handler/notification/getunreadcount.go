package notification

import (
	"context"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type GetUnreadCountLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *GetUnreadCountLogic) GetUnreadCount() (resp *types.GetUnreadCountResponse, err error) {
	// Check if notifications are enabled
	if !l.svcCtx.Config.IsNotificationsEnabled() {
		return &types.GetUnreadCountResponse{Count: 0}, nil
	}

	if !l.svcCtx.UseLocal() {
		return &types.GetUnreadCountResponse{Count: 0}, nil
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// Get unread count
	count, err := l.svcCtx.DB.Queries.CountUnreadNotifications(l.ctx, userID.String())
	if err != nil {
		slog.Error("Failed to count unread notifications", "error", err)
		return nil, err
	}

	return &types.GetUnreadCountResponse{Count: int(count)}, nil
}

func GetUnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &GetUnreadCountLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetUnreadCount()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
