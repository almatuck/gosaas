package organization

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type GetOrganizationLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetOrganizationHandler returns organization details by ID.
func GetOrganizationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOrganizationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &GetOrganizationLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetOrganization(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *GetOrganizationLogic) GetOrganization(req *types.GetOrganizationRequest) (resp *types.GetOrganizationResponse, err error) {
	if !l.svcCtx.Config.IsOrganizationsEnabled() {
		return nil, fmt.Errorf("organizations feature is not enabled")
	}

	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("organizations not available in this mode")
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// Get organization
	org, err := l.svcCtx.DB.Queries.GetOrganizationByID(l.ctx, req.Id)
	if err != nil {
		slog.Error("Failed to get organization", "error", err)
		return nil, fmt.Errorf("organization not found")
	}

	// Check if user is a member
	member, err := l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: req.Id,
		UserID:         userID.String(),
	})
	if err != nil {
		slog.Error("Failed to get member", "error", err)
		return nil, fmt.Errorf("you are not a member of this organization")
	}

	return &types.GetOrganizationResponse{
		Organization: types.Organization{
			Id:        org.ID,
			Name:      org.Name,
			Slug:      org.Slug,
			LogoUrl:   org.LogoUrl.String,
			OwnerId:   org.OwnerID,
			CreatedAt: time.Unix(org.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt: time.Unix(org.UpdatedAt, 0).Format(time.RFC3339),
		},
		Role: member.Role,
	}, nil
}
