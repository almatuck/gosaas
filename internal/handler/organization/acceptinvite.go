package organization

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/google/uuid"
)

type AcceptInviteLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// AcceptInviteHandler accepts an organization invite.
func AcceptInviteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AcceptInviteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &AcceptInviteLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.AcceptInvite(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *AcceptInviteLogic) AcceptInvite(req *types.AcceptInviteRequest) (resp *types.AcceptInviteResponse, err error) {
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

	// Get invite by token
	invite, err := l.svcCtx.DB.Queries.GetInviteByToken(l.ctx, req.Token)
	if err != nil {
		slog.Error("Failed to get invite", "error", err)
		return nil, fmt.Errorf("invalid or expired invite")
	}

	// Check if already a member
	_, err = l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: invite.OrganizationID,
		UserID:         userID.String(),
	})
	if err == nil {
		// Already a member, delete invite and return success
		_ = l.svcCtx.DB.Queries.DeleteInviteByToken(l.ctx, req.Token)
		org, _ := l.svcCtx.DB.Queries.GetOrganizationByID(l.ctx, invite.OrganizationID)
		return &types.AcceptInviteResponse{
			Organization: types.Organization{
				Id:        org.ID,
				Name:      org.Name,
				Slug:      org.Slug,
				LogoUrl:   org.LogoUrl.String,
				OwnerId:   org.OwnerID,
				CreatedAt: time.Unix(org.CreatedAt, 0).Format(time.RFC3339),
				UpdatedAt: time.Unix(org.UpdatedAt, 0).Format(time.RFC3339),
			},
			Message: "You are already a member of this organization",
		}, nil
	}

	// Add user as member
	memberID := uuid.New().String()
	_, err = l.svcCtx.DB.Queries.AddOrganizationMember(l.ctx, db.AddOrganizationMemberParams{
		ID:             memberID,
		OrganizationID: invite.OrganizationID,
		UserID:         userID.String(),
		Role:           invite.Role,
	})
	if err != nil {
		slog.Error("Failed to add member", "error", err)
		return nil, err
	}

	// Delete the invite
	err = l.svcCtx.DB.Queries.DeleteInviteByToken(l.ctx, req.Token)
	if err != nil {
		slog.Error("Failed to delete invite", "error", err)
		// Non-fatal
	}

	// Set as current organization
	err = l.svcCtx.DB.Queries.SetCurrentOrganization(l.ctx, db.SetCurrentOrganizationParams{
		OrganizationID: sql.NullString{String: invite.OrganizationID, Valid: true},
		UserID:         userID.String(),
	})
	if err != nil {
		slog.Error("Failed to set current organization", "error", err)
		// Non-fatal
	}

	// Get organization details
	org, err := l.svcCtx.DB.Queries.GetOrganizationByID(l.ctx, invite.OrganizationID)
	if err != nil {
		slog.Error("Failed to get organization", "error", err)
		return nil, err
	}

	return &types.AcceptInviteResponse{
		Organization: types.Organization{
			Id:        org.ID,
			Name:      org.Name,
			Slug:      org.Slug,
			LogoUrl:   org.LogoUrl.String,
			OwnerId:   org.OwnerID,
			CreatedAt: time.Unix(org.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt: time.Unix(org.UpdatedAt, 0).Format(time.RFC3339),
		},
		Message: "Successfully joined " + org.Name,
	}, nil
}
