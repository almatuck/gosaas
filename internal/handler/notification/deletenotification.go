package notification

import (
	"context"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type DeleteNotificationLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *DeleteNotificationLogic) DeleteNotification(req *types.DeleteNotificationRequest) (resp *types.MessageResponse, err error) {
	if !l.svcCtx.Config.IsNotificationsEnabled() {
		return &types.MessageResponse{Message: "Notifications not enabled"}, nil
	}

	if !l.svcCtx.UseLocal() {
		return &types.MessageResponse{Message: "Notification deleted"}, nil
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// Delete notification
	err = l.svcCtx.DB.Queries.DeleteNotification(l.ctx, db.DeleteNotificationParams{
		ID:     req.Id,
		UserID: userID.String(),
	})
	if err != nil {
		slog.Error("Failed to delete notification", "error", err)
		return nil, err
	}

	return &types.MessageResponse{Message: "Notification deleted"}, nil
}

func DeleteNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteNotificationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &DeleteNotificationLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.DeleteNotification(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
