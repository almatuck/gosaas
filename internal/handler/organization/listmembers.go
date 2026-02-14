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

type ListMembersLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ListMembersHandler lists all members of an organization.
func ListMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListMembersRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &ListMembersLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ListMembers(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *ListMembersLogic) ListMembers(req *types.ListMembersRequest) (resp *types.ListMembersResponse, err error) {
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

	// Check if user is a member of the organization
	_, err = l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: req.Id,
		UserID:         userID.String(),
	})
	if err != nil {
		slog.Error("Failed to get member", "error", err)
		return nil, fmt.Errorf("you are not a member of this organization")
	}

	// List all members
	members, err := l.svcCtx.DB.Queries.ListOrganizationMembers(l.ctx, req.Id)
	if err != nil {
		slog.Error("Failed to list members", "error", err)
		return nil, err
	}

	// Convert to response type
	result := make([]types.OrganizationMember, len(members))
	for i, m := range members {
		result[i] = types.OrganizationMember{
			Id:        m.ID,
			UserId:    m.UserID,
			Email:     m.Email,
			Name:      m.UserName,
			AvatarUrl: m.AvatarUrl.String,
			Role:      m.Role,
			JoinedAt:  time.Unix(m.JoinedAt, 0).Format(time.RFC3339),
		}
	}

	return &types.ListMembersResponse{
		Members: result,
	}, nil
}
