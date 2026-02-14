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

type MarkAllNotificationsReadLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *MarkAllNotificationsReadLogic) MarkAllNotificationsRead() (resp *types.MessageResponse, err error) {
	// Check if notifications are enabled
	if !l.svcCtx.Config.IsNotificationsEnabled() {
		return &types.MessageResponse{Message: "Notifications not enabled"}, nil
	}

	if !l.svcCtx.UseLocal() {
		return &types.MessageResponse{Message: "All notifications marked as read"}, nil
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// Mark all as read
	err = l.svcCtx.DB.Queries.MarkAllNotificationsRead(l.ctx, userID.String())
	if err != nil {
		slog.Error("Failed to mark all notifications as read", "error", err)
		return nil, err
	}

	return &types.MessageResponse{Message: "All notifications marked as read"}, nil
}

func MarkAllNotificationsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &MarkAllNotificationsReadLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.MarkAllNotificationsRead()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
