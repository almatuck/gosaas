package organization

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type SwitchOrganizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// SwitchOrganizationHandler switches the current user's active organization.
func SwitchOrganizationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SwitchOrganizationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &SwitchOrganizationLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.SwitchOrganization(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(req *types.SwitchOrganizationRequest) (resp *types.MessageResponse, err error) {
	if !l.svcCtx.Config.IsOrganizationsEnabled() {
		return nil, fmt.Errorf("organizations feature is not enabled")
	}

	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("organizations not available in this mode")
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		l.Errorf("Failed to get user ID: %v", err)
		return nil, err
	}

	// Verify user is a member of the target organization
	_, err = l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: req.OrganizationId,
		UserID:         userID.String(),
	})
	if err != nil {
		l.Errorf("Failed to get member: %v", err)
		return nil, fmt.Errorf("you are not a member of this organization")
	}

	// Set as current organization
	err = l.svcCtx.DB.Queries.SetCurrentOrganization(l.ctx, db.SetCurrentOrganizationParams{
		OrganizationID: sql.NullString{String: req.OrganizationId, Valid: true},
		UserID:         userID.String(),
	})
	if err != nil {
		l.Errorf("Failed to switch organization: %v", err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Switched to organization",
	}, nil
}
