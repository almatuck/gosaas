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

type GetUnreadCountLogic struct {
	logx.Logger
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
		l.Errorf("Failed to get user ID: %v", err)
		return nil, err
	}

	// Get unread count
	count, err := l.svcCtx.DB.Queries.CountUnreadNotifications(l.ctx, userID.String())
	if err != nil {
		l.Errorf("Failed to count unread notifications: %v", err)
		return nil, err
	}

	return &types.GetUnreadCountResponse{Count: int(count)}, nil
}

func GetUnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &GetUnreadCountLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetUnreadCount()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
