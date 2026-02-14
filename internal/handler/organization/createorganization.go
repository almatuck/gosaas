package organization

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"gosaas/internal/auth"
	"gosaas/internal/db"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/google/uuid"
)

type CreateOrganizationLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// CreateOrganizationHandler creates a new organization.
func CreateOrganizationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrganizationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &CreateOrganizationLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.CreateOrganization(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func (l *CreateOrganizationLogic) CreateOrganization(req *types.CreateOrganizationRequest) (resp *types.CreateOrganizationResponse, err error) {
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

	// Generate slug from name if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	// Check if slug already exists
	found, err := l.svcCtx.DB.Queries.CheckSlugExists(l.ctx, slug)
	if err != nil {
		slog.Error("Failed to check slug", "error", err)
		return nil, err
	}
	if found > 0 {
		return nil, fmt.Errorf("organization with this slug already exists")
	}

	// Create organization
	orgID := uuid.New().String()
	org, err := l.svcCtx.DB.Queries.CreateOrganization(l.ctx, db.CreateOrganizationParams{
		ID:      orgID,
		Name:    req.Name,
		Slug:    slug,
		LogoUrl: sql.NullString{String: req.LogoUrl, Valid: req.LogoUrl != ""},
		OwnerID: userID.String(),
	})
	if err != nil {
		slog.Error("Failed to create organization", "error", err)
		return nil, err
	}

	// Add owner as member with owner role
	memberID := uuid.New().String()
	_, err = l.svcCtx.DB.Queries.AddOrganizationMember(l.ctx, db.AddOrganizationMemberParams{
		ID:             memberID,
		OrganizationID: orgID,
		UserID:         userID.String(),
		Role:           "owner",
	})
	if err != nil {
		slog.Error("Failed to add owner as member", "error", err)
		return nil, err
	}

	// Set as current organization for the user
	err = l.svcCtx.DB.Queries.SetCurrentOrganization(l.ctx, db.SetCurrentOrganizationParams{
		OrganizationID: sql.NullString{String: orgID, Valid: true},
		UserID:         userID.String(),
	})
	if err != nil {
		slog.Error("Failed to set current organization", "error", err)
		// Non-fatal, continue
	}

	return &types.CreateOrganizationResponse{
		Organization: types.Organization{
			Id:        org.ID,
			Name:      org.Name,
			Slug:      org.Slug,
			LogoUrl:   org.LogoUrl.String,
			OwnerId:   org.OwnerID,
			CreatedAt: time.Unix(org.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt: time.Unix(org.UpdatedAt, 0).Format(time.RFC3339),
		},
	}, nil
}

func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	// Trim hyphens from ends
	slug = strings.Trim(slug, "-")
	return slug
}
