package organization

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type DeleteOrganizationLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// DeleteOrganizationHandler deletes an organization.
func DeleteOrganizationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteOrganizationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &DeleteOrganizationLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.DeleteOrganization(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *DeleteOrganizationLogic) DeleteOrganization(req *types.DeleteOrganizationRequest) (resp *types.MessageResponse, err error) {
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

	// Get organization to check ownership
	org, err := l.svcCtx.DB.Queries.GetOrganizationByID(l.ctx, req.Id)
	if err != nil {
		slog.Error("Failed to get organization", "error", err)
		return nil, fmt.Errorf("organization not found")
	}

	// Only the owner can delete an organization
	if org.OwnerID != userID.String() {
		return nil, fmt.Errorf("only the owner can delete this organization")
	}

	// Delete the organization (cascades to members and invites via foreign keys)
	err = l.svcCtx.DB.Queries.DeleteOrganization(l.ctx, req.Id)
	if err != nil {
		slog.Error("Failed to delete organization", "error", err)
		return nil, err
	}

	// Clear current organization for the user if this was their current org
	_ = l.svcCtx.DB.Queries.SetCurrentOrganization(l.ctx, db.SetCurrentOrganizationParams{
		OrganizationID: sql.NullString{Valid: false},
		UserID:         userID.String(),
	})

	return &types.MessageResponse{
		Message: "Organization deleted successfully",
	}, nil
}
