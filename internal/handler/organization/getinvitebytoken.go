package organization

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type GetInviteByTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetInviteByTokenHandler returns invite details by token.
func GetInviteByTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetInviteByTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &GetInviteByTokenLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetInviteByToken(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func (l *GetInviteByTokenLogic) GetInviteByToken(req *types.GetInviteByTokenRequest) (resp *types.GetInviteByTokenResponse, err error) {
	if !l.svcCtx.Config.IsOrganizationsEnabled() {
		return nil, fmt.Errorf("organizations feature is not enabled")
	}

	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("organizations not available in this mode")
	}

	// Get invite by token (this is a public endpoint - no auth required)
	invite, err := l.svcCtx.DB.Queries.GetInviteByToken(l.ctx, req.Token)
	if err != nil {
		l.Errorf("Failed to get invite: %v", err)
		return nil, fmt.Errorf("invalid or expired invite")
	}

	return &types.GetInviteByTokenResponse{
		Invite: types.OrganizationInvite{
			Id:               invite.ID,
			Email:            invite.Email,
			Role:             invite.Role,
			OrganizationName: invite.OrganizationName,
			ExpiresAt:        time.Unix(invite.ExpiresAt, 0).Format(time.RFC3339),
			CreatedAt:        time.Unix(invite.CreatedAt, 0).Format(time.RFC3339),
		},
		OrganizationName: invite.OrganizationName,
		OrganizationSlug: invite.OrganizationSlug,
	}, nil
}
