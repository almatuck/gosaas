package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// AdminAuthMiddleware checks if the authenticated user is an admin.
// This runs AFTER JWT middleware has validated the token.
type AdminAuthMiddleware struct {
	adminUsername string
}

func NewAdminAuthMiddleware(adminUsername string) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		adminUsername: adminUsername,
	}
}

func (m *AdminAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.adminUsername == "" {
			slog.Error("[AdminAuth] Admin username not configured")
			adminForbidden(w, "Admin access not configured")
			return
		}

		email := r.Context().Value("email")
		if email == nil {
			slog.Error("[AdminAuth] No email claim in JWT")
			adminForbidden(w, "Admin access required")
			return
		}

		emailStr, ok := email.(string)
		if !ok {
			slog.Error("[AdminAuth] Invalid email claim type in JWT")
			adminForbidden(w, "Admin access required")
			return
		}

		if emailStr != m.adminUsername {
			slog.Info("[AdminAuth] Non-admin user attempted admin access", "email", emailStr)
			adminForbidden(w, "Admin access required")
			return
		}

		slog.Info("[AdminAuth] Admin access granted", "email", emailStr)
		next.ServeHTTP(w, r)
	})
}

// adminForbidden sends a 403 response for non-admin users
func adminForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
