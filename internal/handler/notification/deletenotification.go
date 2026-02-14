package notification

import (
	"context"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type DeleteNotificationLogic struct {
	logx.Logger
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
		l.Errorf("Failed to get user ID: %v", err)
		return nil, err
	}

	// Delete notification
	err = l.svcCtx.DB.Queries.DeleteNotification(l.ctx, db.DeleteNotificationParams{
		ID:     req.Id,
		UserID: userID.String(),
	})
	if err != nil {
		l.Errorf("Failed to delete notification: %v", err)
		return nil, err
	}

	return &types.MessageResponse{Message: "Notification deleted"}, nil
}

func DeleteNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteNotificationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &DeleteNotificationLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.DeleteNotification(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
