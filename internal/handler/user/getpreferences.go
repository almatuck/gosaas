package user

import (
	"context"
	"net/http"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type GetPreferencesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *GetPreferencesLogic) GetPreferences() (resp *types.GetPreferencesResponse, err error) {
	// Preferences are stored locally - for the boilerplate, return defaults
	// In a real app, you'd store these in your database or in Levee custom fields
	return &types.GetPreferencesResponse{
		Preferences: types.UserPreferences{
			EmailNotifications: true,
			MarketingEmails:    true,
			Timezone:           "UTC",
			Language:           "en",
			Theme:              "system",
		},
	}, nil
}

// GetPreferencesHandler handles the get user preferences request.
func GetPreferencesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &GetPreferencesLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.GetPreferences()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
