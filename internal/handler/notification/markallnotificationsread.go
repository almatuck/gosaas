package notification

import (
	"context"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type MarkAllNotificationsReadLogic struct {
	logx.Logger
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
		l.Errorf("Failed to get user ID: %v", err)
		return nil, err
	}

	// Mark all as read
	err = l.svcCtx.DB.Queries.MarkAllNotificationsRead(l.ctx, userID.String())
	if err != nil {
		l.Errorf("Failed to mark all notifications as read: %v", err)
		return nil, err
	}

	return &types.MessageResponse{Message: "All notifications marked as read"}, nil
}

func MarkAllNotificationsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &MarkAllNotificationsReadLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.MarkAllNotificationsRead()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
