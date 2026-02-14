package organization

import (
	"context"
	"fmt"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type RemoveMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// RemoveMemberHandler removes a member from an organization.
func RemoveMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveMemberRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &RemoveMemberLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.RemoveMember(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func (l *RemoveMemberLogic) RemoveMember(req *types.RemoveMemberRequest) (resp *types.MessageResponse, err error) {
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

	// Get organization to check ownership
	org, err := l.svcCtx.DB.Queries.GetOrganizationByID(l.ctx, req.OrgId)
	if err != nil {
		l.Errorf("Failed to get organization: %v", err)
		return nil, fmt.Errorf("organization not found")
	}

	// Cannot remove the owner
	if req.UserId == org.OwnerID {
		return nil, fmt.Errorf("cannot remove the organization owner")
	}

	// Check if user has permission to remove members (owner or admin)
	currentMember, err := l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
		OrganizationID: req.OrgId,
		UserID:         userID.String(),
	})
	if err != nil {
		l.Errorf("Failed to get member: %v", err)
		return nil, fmt.Errorf("you are not a member of this organization")
	}
	if currentMember.Role != "owner" && currentMember.Role != "admin" {
		return nil, fmt.Errorf("you do not have permission to remove members")
	}

	// Admins cannot remove other admins or the owner
	if currentMember.Role == "admin" {
		targetMember, err := l.svcCtx.DB.Queries.GetOrganizationMember(l.ctx, db.GetOrganizationMemberParams{
			OrganizationID: req.OrgId,
			UserID:         req.UserId,
		})
		if err != nil {
			return nil, fmt.Errorf("user is not a member of this organization")
		}
		if targetMember.Role == "owner" || targetMember.Role == "admin" {
			return nil, fmt.Errorf("admins cannot remove other admins or the owner")
		}
	}

	// Remove the member
	err = l.svcCtx.DB.Queries.RemoveOrganizationMember(l.ctx, db.RemoveOrganizationMemberParams{
		OrganizationID: req.OrgId,
		UserID:         req.UserId,
	})
	if err != nil {
		l.Errorf("Failed to remove member: %v", err)
		return nil, err
	}

	return &types.MessageResponse{
		Message: "Member removed successfully",
	}, nil
}
