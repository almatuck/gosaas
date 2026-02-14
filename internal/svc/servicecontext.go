package svc

import (
	"log/slog"
	"net/http"

	"gosaas/internal/config"
	"gosaas/internal/db"
	"gosaas/internal/local"
	"gosaas/internal/middleware"

	leveeSDK "github.com/almatuck/levee-go"
)

type ServiceContext struct {
	Config             config.Config
	SecurityMiddleware *middleware.SecurityMiddleware
	AdminAuth          func(http.Handler) http.Handler // Admin backoffice auth

	// Levee SDK (used when Levee.Enabled=true)
	Levee *leveeSDK.Client

	// Local services (used when Levee.Enabled=false)
	DB      *db.Store
	Auth    *local.AuthService
	Billing *local.BillingService
	Email   *local.EmailService // SMTP or Outlet.sh (available in both modes)
}

// NewServiceContext creates a new service context
func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize security middleware
	securityMw := middleware.NewSecurityMiddleware(c)
	slog.Info("Security middleware initialized")

	// Initialize admin auth middleware
	adminAuth := middleware.NewAdminAuthMiddleware(c.Admin.Username)

	svc := &ServiceContext{
		Config:             c,
		SecurityMiddleware: securityMw,
		AdminAuth:          adminAuth.Handle,
	}

	// Initialize email service (SMTP or Outlet.sh - available in both modes)
	emailService := local.NewEmailService(c)
	if emailService.IsConfigured() {
		svc.Email = emailService
		slog.Info("Email service initialized")
	} else {
		slog.Info("Email not configured - transactional emails disabled")
	}

	if c.IsLeveeEnabled() && c.Levee.APIKey != "" {
		// Mode 1: Use Levee for auth, billing, email
		baseURL := c.Levee.BaseURL
		if baseURL == "" {
			baseURL = "https://api.levee.sh"
		}
		client, err := leveeSDK.NewClient(c.Levee.APIKey, baseURL)
		if err != nil {
			slog.Error("Failed to create Levee client", "error", err)
		} else {
			svc.Levee = client
			slog.Info("Levee SDK client initialized (using Levee for auth/billing)")
		}
	} else {
		// Mode 2: Use local SQLite + direct Stripe
		slog.Info("Levee disabled - using local SQLite + direct Stripe")

		database, err := db.NewSQLite(c.Database.SQLitePath)
		if err != nil {
			slog.Error("Failed to initialize SQLite database", "error", err)
		} else {
			svc.DB = database
			slog.Info("SQLite database initialized", "path", c.Database.SQLitePath)

			svc.Auth = local.NewAuthService(database, c)
			svc.Billing = local.NewBillingService(database, c)
			slog.Info("Local auth and billing services initialized")
		}
	}

	return svc
}

// Close closes any open connections
func (svc *ServiceContext) Close() {
	if svc.DB != nil {
		svc.DB.Close()
		slog.Info("SQLite database connection closed")
	}
	slog.Info("Service context closed")
}

// UseLevee returns true if Levee is enabled and configured
func (svc *ServiceContext) UseLevee() bool {
	return svc.Levee != nil
}

// UseLocal returns true if using local SQLite + Stripe
func (svc *ServiceContext) UseLocal() bool {
	return svc.DB != nil
}
