package organization

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/auth"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type ListOrganizationsLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ListOrganizationsHandler lists all organizations for the current user.
func ListOrganizationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &ListOrganizationsLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ListOrganizations()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *ListOrganizationsLogic) ListOrganizations() (resp *types.ListOrganizationsResponse, err error) {
	if !l.svcCtx.Config.IsOrganizationsEnabled() {
		return &types.ListOrganizationsResponse{Organizations: []types.Organization{}}, nil
	}

	if !l.svcCtx.UseLocal() {
		return &types.ListOrganizationsResponse{Organizations: []types.Organization{}}, nil
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	// List organizations for user
	orgs, err := l.svcCtx.DB.Queries.ListUserOrganizations(l.ctx, userID.String())
	if err != nil {
		slog.Error("Failed to list organizations", "error", err)
		return nil, err
	}

	// Get current organization
	currentOrgID, _ := l.svcCtx.DB.Queries.GetCurrentOrganization(l.ctx, userID.String())

	// Convert to response type
	result := make([]types.Organization, len(orgs))
	for i, org := range orgs {
		result[i] = types.Organization{
			Id:        org.ID,
			Name:      org.Name,
			Slug:      org.Slug,
			LogoUrl:   org.LogoUrl.String,
			OwnerId:   org.OwnerID,
			CreatedAt: time.Unix(org.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt: time.Unix(org.UpdatedAt, 0).Format(time.RFC3339),
		}
	}

	// Note: currentOrgID is available but not returned in the response type
	_ = currentOrgID

	return &types.ListOrganizationsResponse{
		Organizations: result,
	}, nil
}
