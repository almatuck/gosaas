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

type ListInvitesLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ListInvitesHandler lists pending invites for an organization.
func ListInvitesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListInvitesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &ListInvitesLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ListInvites(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *ListInvitesLogic) ListInvites(req *types.ListInvitesRequest) (resp *types.ListInvitesResponse, err error) {
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

	// Check if user has permission to view invites (owner or admin)
	member, err := l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: req.Id,
		UserID:         userID.String(),
	})
	if err != nil {
		slog.Error("Failed to get member", "error", err)
		return nil, fmt.Errorf("you are not a member of this organization")
	}
	if member.Role != "owner" && member.Role != "admin" {
		return nil, fmt.Errorf("you do not have permission to view invites")
	}

	// List pending invites
	invites, err := l.svcCtx.DB.Queries.ListOrganizationInvites(l.ctx, req.Id)
	if err != nil {
		slog.Error("Failed to list invites", "error", err)
		return nil, err
	}

	// Convert to response type
	result := make([]types.OrganizationInvite, len(invites))
	for i, inv := range invites {
		result[i] = types.OrganizationInvite{
			Id:           inv.ID,
			Email:        inv.Email,
			Role:         inv.Role,
			InviterName:  inv.InviterName,
			InviterEmail: inv.InviterEmail,
			ExpiresAt:    time.Unix(inv.ExpiresAt, 0).Format(time.RFC3339),
			CreatedAt:    time.Unix(inv.CreatedAt, 0).Format(time.RFC3339),
		}
	}

	return &types.ListInvitesResponse{
		Invites: result,
	}, nil
}
