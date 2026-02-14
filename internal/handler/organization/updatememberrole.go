package organization

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type UpdateMemberRoleLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// UpdateMemberRoleHandler updates a member's role in an organization.
func UpdateMemberRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMemberRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &UpdateMemberRoleLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.UpdateMemberRole(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *UpdateMemberRoleLogic) UpdateMemberRole(req *types.UpdateMemberRoleRequest) (resp *types.MessageResponse, err error) {
	if !l.svcCtx.Config.IsOrganizationsEnabled() {
		return nil, fmt.Errorf("organizations feature is not enabled")
	}

	if !l.svcCtx.UseLocal() {
		return nil, fmt.Errorf("organizations not available in this mode")
	}

	// Validate role
	validRoles := map[string]bool{"owner": true, "admin": true, "member": true}
	if !validRoles[req.Role] {
		return nil, fmt.Errorf("invalid role: must be owner, admin, or member")
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// Get organization to check ownership
	org, err := l.svcCtx.DB.Queries.GetOrganizationByID(l.ctx, req.OrgId)
	if err != nil {
		slog.Error("Failed to get organization", "error", err)
		return nil, fmt.Errorf("organization not found")
	}

	// Check if user is owner or admin
	currentMember, err := l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: req.OrgId,
		UserID:         userID.String(),
	})
	if err != nil {
		slog.Error("Failed to get member", "error", err)
		return nil, fmt.Errorf("you are not a member of this organization")
	}

	// Only owner can change roles
	if currentMember.Role != "owner" {
		return nil, fmt.Errorf("only the owner can change member roles")
	}

	// Cannot change owner's role (transfer ownership is a separate operation)
	if req.UserId == org.OwnerID && req.Role != "owner" {
		return nil, fmt.Errorf("cannot change the owner's role")
	}

	// Cannot make someone else an owner through this endpoint
	if req.Role == "owner" && req.UserId != org.OwnerID {
		return nil, fmt.Errorf("cannot assign owner role; use transfer ownership instead")
	}

	// Verify target user is a member
	_, err = l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: req.OrgId,
		UserID:         req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("user is not a member of this organization")
	}

	// Update the role
	err = l.svcCtx.DB.Queries.UpdateMemberRole(l.ctx, db.UpdateMemberRoleParams{
		Role:           req.Role,
		OrganizationID: req.OrgId,
		UserID:         req.UserId,
	})
	if err != nil {
		slog.Error("Failed to update member role", "error", err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Member role updated successfully",
	}, nil
}
