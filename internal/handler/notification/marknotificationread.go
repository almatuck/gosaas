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

type MarkNotificationReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *MarkNotificationReadLogic) MarkNotificationRead(req *types.MarkNotificationReadRequest) (resp *types.MessageResponse, err error) {
	if !l.svcCtx.Config.IsNotificationsEnabled() {
		return &types.MessageResponse{Message: "Notifications not enabled"}, nil
	}

	if !l.svcCtx.UseLocal() {
		return &types.MessageResponse{Message: "Notification marked as read"}, nil
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		l.Errorf("Failed to get user ID: %v", err)
		return nil, err
	}

	// Mark as read
	err = l.svcCtx.DB.Queries.MarkNotificationRead(l.ctx, db.MarkNotificationReadParams{
		ID:     req.Id,
		UserID: userID.String(),
	})
	if err != nil {
		l.Errorf("Failed to mark notification as read: %v", err)
		return nil, err
	}

	return &types.MessageResponse{Message: "Notification marked as read"}, nil
}

func MarkNotificationReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarkNotificationReadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &MarkNotificationReadLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.MarkNotificationRead(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
