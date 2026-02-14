package user

import (
	"context"
	"log/slog"
	"net/http"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type UpdatePreferencesLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *UpdatePreferencesLogic) UpdatePreferences(req *types.UpdatePreferencesRequest) (resp *types.GetPreferencesResponse, err error) {
	// Preferences would be stored in your database or Levee custom fields
	// For the boilerplate, return the updated values
	return &types.GetPreferencesResponse{
		Preferences: types.UserPreferences{
			EmailNotifications: req.EmailNotifications,
			MarketingEmails:    req.MarketingEmails,
			Timezone:           req.Timezone,
			Language:           req.Language,
			Theme:              req.Theme,
		},
	}, nil
}

// UpdatePreferencesHandler handles the update user preferences request.
func UpdatePreferencesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePreferencesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &UpdatePreferencesLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.UpdatePreferences(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
