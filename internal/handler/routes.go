// Package handler registers all API routes.
package handler

import (
	"net/http"

	"gosaas/internal/handler/admin"
	"gosaas/internal/handler/auth"
	"gosaas/internal/handler/notification"
	"gosaas/internal/handler/oauth"
	"gosaas/internal/handler/organization"
	"gosaas/internal/handler/subscription"
	"gosaas/internal/handler/user"
	"gosaas/internal/middleware"
	"gosaas/internal/svc"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(r chi.Router, svcCtx *svc.ServiceContext) {
	// Health check (no prefix, no auth)
	r.Get("/health", HealthCheckHandler(svcCtx))

	r.Route("/api/v1", func(r chi.Router) {
		// ── Public routes (no auth) ──────────────────────────────

		// Auth
		r.Get("/auth/config", auth.GetAuthConfigHandler(svcCtx))
		r.Post("/auth/forgot-password", auth.ForgotPasswordHandler(svcCtx))
		r.Post("/auth/login", auth.LoginHandler(svcCtx))
		r.Post("/auth/refresh", auth.RefreshTokenHandler(svcCtx))
		r.Post("/auth/register", auth.RegisterHandler(svcCtx))
		r.Post("/auth/resend-verification", auth.ResendVerificationHandler(svcCtx))
		r.Post("/auth/reset-password", auth.ResetPasswordHandler(svcCtx))
		r.Post("/auth/verify-email", auth.VerifyEmailHandler(svcCtx))

		// OAuth (public - callback and URL generation)
		r.Post("/oauth/{provider}/callback", oauth.OAuthCallbackHandler(svcCtx))
		r.Get("/oauth/{provider}/url", oauth.GetOAuthUrlHandler(svcCtx))

		// Organization invite (public - token-based)
		r.Get("/organizations/invites/{token}", organization.GetInviteByTokenHandler(svcCtx))

		// Subscription plans (public)
		r.Get("/subscription/plans", subscription.ListPlansHandler(svcCtx))

		// ── JWT-protected routes ─────────────────────────────────

		r.Group(func(r chi.Router) {
			// Apply JWT auth based on mode
			if svcCtx.UseLevee() {
				r.Use(middleware.LeveeTokenTranslator(svcCtx.Config.Auth.AccessSecret))
				r.Use(middleware.JWTAuth(svcCtx.Config.Auth.AccessSecret))
			} else {
				r.Use(middleware.JWTAuth(svcCtx.Config.Auth.AccessSecret))
			}

			// Notifications
			r.Get("/notifications", notification.ListNotificationsHandler(svcCtx))
			r.Delete("/notifications/{id}", notification.DeleteNotificationHandler(svcCtx))
			r.Put("/notifications/{id}/read", notification.MarkNotificationReadHandler(svcCtx))
			r.Put("/notifications/read-all", notification.MarkAllNotificationsReadHandler(svcCtx))
			r.Get("/notifications/unread-count", notification.GetUnreadCountHandler(svcCtx))

			// OAuth (authenticated - manage connections)
			r.Delete("/oauth/{provider}", oauth.DisconnectOAuthHandler(svcCtx))
			r.Get("/oauth/providers", oauth.ListOAuthProvidersHandler(svcCtx))

			// Organizations
			r.Post("/organizations", organization.CreateOrganizationHandler(svcCtx))
			r.Get("/organizations", organization.ListOrganizationsHandler(svcCtx))
			r.Get("/organizations/{id}", organization.GetOrganizationHandler(svcCtx))
			r.Put("/organizations/{id}", organization.UpdateOrganizationHandler(svcCtx))
			r.Delete("/organizations/{id}", organization.DeleteOrganizationHandler(svcCtx))
			r.Post("/organizations/{id}/invites", organization.InviteMemberHandler(svcCtx))
			r.Get("/organizations/{id}/invites", organization.ListInvitesHandler(svcCtx))
			r.Post("/organizations/{id}/leave", organization.LeaveOrganizationHandler(svcCtx))
			r.Get("/organizations/{id}/members", organization.ListMembersHandler(svcCtx))
			r.Delete("/organizations/{orgId}/invites/{inviteId}", organization.RevokeInviteHandler(svcCtx))
			r.Put("/organizations/{orgId}/members/{userId}", organization.UpdateMemberRoleHandler(svcCtx))
			r.Delete("/organizations/{orgId}/members/{userId}", organization.RemoveMemberHandler(svcCtx))
			r.Post("/organizations/invites/accept", organization.AcceptInviteHandler(svcCtx))
			r.Post("/organizations/switch", organization.SwitchOrganizationHandler(svcCtx))

			// Subscription
			r.Get("/subscription", subscription.GetSubscriptionHandler(svcCtx))
			r.Get("/subscription/billing-history", subscription.ListBillingHistoryHandler(svcCtx))
			r.Post("/subscription/billing-portal", subscription.CreateBillingPortalHandler(svcCtx))
			r.Post("/subscription/cancel", subscription.CancelSubscriptionHandler(svcCtx))
			r.Post("/subscription/check-feature", subscription.CheckFeatureHandler(svcCtx))
			r.Post("/subscription/checkout", subscription.CreateCheckoutHandler(svcCtx))
			r.Get("/subscription/usage", subscription.GetUsageStatsHandler(svcCtx))

			// User
			r.Get("/user/me", user.GetCurrentUserHandler(svcCtx))
			r.Put("/user/me", user.UpdateCurrentUserHandler(svcCtx))
			r.Delete("/user/me", user.DeleteAccountHandler(svcCtx))
			r.Post("/user/me/change-password", user.ChangePasswordHandler(svcCtx))
			r.Get("/user/me/preferences", user.GetPreferencesHandler(svcCtx))
			r.Put("/user/me/preferences", user.UpdatePreferencesHandler(svcCtx))

			// Admin (JWT + admin auth)
			r.Group(func(r chi.Router) {
				r.Use(svcCtx.AdminAuth)
				r.Get("/admin/stats", admin.GetAdminStatsHandler(svcCtx))
				r.Get("/admin/subscriptions", admin.AdminListSubscriptionsHandler(svcCtx))
				r.Get("/admin/users", admin.AdminListUsersHandler(svcCtx))
			})
		})
	})

	// CSRF token endpoint (no auth, no prefix)
	r.Get("/api/v1/csrf-token", func(w http.ResponseWriter, r *http.Request) {
		svcCtx.SecurityMiddleware.GetCSRFTokenHandler().ServeHTTP(w, r)
	})
}
