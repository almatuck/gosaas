package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	App struct {
		BaseURL        string `json:",optional"`
		Domain         string `json:",optional"`
		ProductionMode bool   `json:",default=false"`
		AdminEmail     string `json:",optional"` // For Let's Encrypt notifications
	}
	Admin struct {
		Username string `json:",optional"` // Backoffice admin username
		Password string `json:",optional"` // Backoffice admin password
	}
	Auth struct {
		AccessSecret       string
		AccessExpire       int64
		RefreshTokenExpire int64 `json:",default=604800"` // 7 days in seconds
	}
	Database struct {
		// SQLite (used when Levee is disabled)
		SQLitePath string `json:",default=./data/gosaas.db"`

		// PostgreSQL (optional, for scaling later)
		Host            string `json:",default=localhost"`
		Port            int    `json:",default=5432"`
		User            string `json:",default=postgres"`
		Password        string `json:",optional"`
		DBName          string `json:",default=gosaas"`
		SSLMode         string `json:",default=disable"`
		MaxOpenConns    int    `json:",default=25"`
		MaxIdleConns    int    `json:",default=5"`
		ConnMaxLifetime int    `json:",default=300"` // seconds
	}
	Stripe struct {
		SecretKey      string `json:",optional"` // sk_test_xxx or sk_live_xxx
		PublishableKey string `json:",optional"` // pk_test_xxx or pk_live_xxx
		WebhookSecret  string `json:",optional"` // whsec_xxx
		SuccessURL     string `json:",default=/app/account"`
		CancelURL      string `json:",default=/pricing"`
	}
	// Products defines subscription products for standalone mode (without Levee)
	// These are synced to Stripe on startup. See docs/stripe.md for configuration details.
	Products []Product `json:",optional"`
	Security struct {
		// CSRF protection settings
		CSRFEnabled      bool   `json:",default=true"`
		CSRFSecret       string `json:",optional"`      // If empty, uses Auth.AccessSecret
		CSRFTokenExpiry  int64  `json:",default=43200"` // 12 hours in seconds
		CSRFSecureCookie bool   `json:",default=true"`

		// Rate limiting settings
		RateLimitEnabled      bool `json:",default=true"`
		RateLimitRequests     int  `json:",default=100"` // requests per interval
		RateLimitInterval     int  `json:",default=60"`  // interval in seconds
		RateLimitBurst        int  `json:",default=20"`  // burst size
		AuthRateLimitRequests int  `json:",default=5"`   // auth endpoints rate limit
		AuthRateLimitInterval int  `json:",default=60"`  // auth rate limit interval

		// Security headers settings
		EnableSecurityHeaders bool   `json:",default=true"`
		ContentSecurityPolicy string `json:",optional"` // Override default CSP
		AllowedOrigins        string `json:",optional"` // CORS allowed origins (comma-separated)

		// HTTPS enforcement
		ForceHTTPS bool `json:",default=false"` // Set to true in production

		// Input validation
		MaxRequestBodySize int64 `json:",default=10485760"` // 10MB default
		MaxURLLength       int   `json:",default=2048"`
	}
	Email struct {
		SMTPHost    string `json:",optional"`
		SMTPPort    int    `json:",default=587"`
		SMTPUser    string `json:",optional"`
		SMTPPass    string `json:",optional"`
		FromAddress string `json:",optional"`
		FromName    string `json:",default=gosaas"`
		ReplyTo     string `json:",optional"`
		BaseURL     string `json:",default=http://localhost:5173"` // For links in emails
	}
	Analytics struct {
		Enabled       bool   `json:",default=true"`    // Enable/disable analytics
		Provider      string `json:",default=console"` // Provider: segment, mixpanel, posthog, console
		APIKey        string `json:",optional"`        // Provider API key
		Endpoint      string `json:",optional"`        // Custom endpoint URL
		BatchSize     int    `json:",default=50"`      // Events per batch
		FlushInterval int    `json:",default=30"`      // Flush interval in seconds
		Debug         bool   `json:",default=false"`   // Debug mode for verbose logging
	}
	Subscription struct {
		Enabled             bool `json:",default=true"` // Enable subscription features
		EnforceQuotas       bool `json:",default=true"` // Enforce usage quotas
		FreeTierAnalyses    int  `json:",default=5"`    // Free tier monthly analysis limit
		FreeTierHistoryDays int  `json:",default=7"`    // Free tier history retention days
		ProTierHistoryDays  int  `json:",default=30"`   // Pro tier history retention days
		TeamTierHistoryDays int  `json:",default=-1"`   // Team tier history retention (-1 = unlimited)
	}
	AI struct {
		Enabled     bool   `json:",default=true"`                     // Enable AI features
		APIKey      string `json:",optional"`                         // Anthropic API key
		Model       string `json:",default=claude-sonnet-4-20250514"` // Claude model to use
		MaxTokens   int    `json:",default=4096"`                     // Max tokens per response
		TimeoutSecs int    `json:",default=60"`                       // Request timeout in seconds
	}
	CostAlerts struct {
		Enabled            bool    `json:",default=true"` // Enable cost alert emails
		AdminEmail         string  `json:",optional"`     // Email to receive cost alerts
		DailyCostThreshold float64 `json:",default=10.0"` // Alert if daily costs exceed this USD amount
		UserCostThreshold  float64 `json:",default=1.0"`  // Alert if single user's daily costs exceed this
	}
	Outlet struct {
		BaseURL string `json:",optional"` // Outlet.sh instance URL (e.g., https://email.yourdomain.com)
		APIKey  string `json:",optional"` // Outlet.sh API key

		// List slugs for subscriber management
		WaitlistListSlug   string `json:",default=waitlist"`
		NewsletterListSlug string `json:",default=newsletter"`
		UsersListSlug      string `json:",default=users"`

		// Sequence slugs for email automation
		OnboardingSequence      string `json:",default=onboarding"`
		TrialConversionSequence string `json:",default=trial-conversion"`

		// Template slugs for transactional emails
		EmailVerificationTemplate string `json:",default=email-verification"`
		PasswordResetTemplate     string `json:",default=password-reset"`
		WelcomeTemplate           string `json:",default=welcome"`
	}
	Levee struct {
		APIKey      string `json:",optional"`     // Levee API key
		BaseURL     string `json:",optional"`     // Custom API endpoint (optional)
		GRPCAddress string `json:",optional"`     // gRPC endpoint for LLM streaming (e.g., levee.localrivet.com:9889)
		Enabled     bool   `json:",default=true"` // Enable Levee integration

		// Checkout redirect URLs
		CheckoutSuccessURL string `json:",default=/app/billing/success"` // Redirect after successful checkout
		CheckoutCancelURL  string `json:",default=/pricing"`             // Redirect after cancelled checkout

		// List slugs
		WaitlistListSlug    string `json:",default=waitlist"`
		NewsletterListSlug  string `json:",default=newsletter-subscribers"`
		UsersListSlug       string `json:",default=users"`
		ProductHuntListSlug string `json:",default=product-hunt-launchers"`

		// Sequence slugs
		WaitlistNurtureSequence string `json:",default=waitlist-nurture"`
		OnboardingSequence      string `json:",default=onboarding-sequence"`
		TrialConversionSequence string `json:",default=free-to-pro-upgrade"`

		// Product slugs for checkout (these map to Levee price nicknames)
		FreeProductSlug        string `json:",default=free"`
		ProMonthlyProductSlug  string `json:",default=pro"`
		ProYearlyProductSlug   string `json:",default=pro-yearly"`
		TeamMonthlyProductSlug string `json:",default=team"`
		TeamYearlyProductSlug  string `json:",default=team-yearly"`

		// Template slugs for transactional emails
		EmailVerificationTemplate     string `json:",default=email-verification"`
		PasswordResetTemplate         string `json:",default=password-reset"`
		PasswordChangedTemplate       string `json:",default=password-changed"`
		AccountDeletedTemplate        string `json:",default=account-deleted"`
		WaitlistConfirmTemplate       string `json:",default=waitlist-confirm"`
		SubscriptionWelcomeTemplate   string `json:",default=subscription-welcome"`
		SubscriptionRenewalTemplate   string `json:",default=subscription-renewal"`
		SubscriptionCancelledTemplate string `json:",default=subscription-cancelled"`
		PaymentFailedTemplate         string `json:",default=payment-failed"`
		TrialEndingTemplate           string `json:",default=trial-ending"`
		CostAlertTemplate             string `json:",default=cost-alert"`
	}
}

// Product defines a subscription product with its prices
// Products are synced to Stripe on startup (standalone mode only)
type Product struct {
	// Slug is the unique identifier for this product (e.g., "free", "pro", "team")
	// Used in checkout URLs and API calls
	Slug string `json:"slug"`

	// Name is the display name shown to customers (e.g., "Pro Plan")
	Name string `json:"name"`

	// Description is shown on pricing page and Stripe checkout
	Description string `json:"description,optional"`

	// Features is a list of features included in this plan
	// Displayed on pricing pages and used for feature gating
	Features []string `json:"features,optional"`

	// Prices defines the pricing options for this product
	// A product can have multiple prices (e.g., monthly and yearly)
	Prices []Price `json:"prices,optional"`

	// Default marks this as the default plan for new users (only one should be true)
	Default bool `json:"default,optional"`

	// StripeProductID is auto-populated after syncing to Stripe
	StripeProductID string `json:"stripeProductId,optional"`
}

// Price defines a single pricing option for a product
type Price struct {
	// Slug identifies this price (e.g., "monthly", "yearly")
	Slug string `json:"slug"`

	// Amount in cents (e.g., 2900 for $29.00, 0 for free)
	Amount int64 `json:"amount"`

	// Currency code (default: "usd")
	Currency string `json:"currency,default=usd"`

	// Interval: "month", "year", or "one_time"
	Interval string `json:"interval,default=month"`

	// IntervalCount: billing frequency (default: 1)
	// e.g., interval=month + intervalCount=3 = quarterly billing
	IntervalCount int64 `json:"intervalCount,default=1"`

	// TrialDays: number of trial days (0 = no trial)
	TrialDays int64 `json:"trialDays,optional"`

	// StripePriceID is auto-populated after syncing to Stripe
	StripePriceID string `json:"stripePriceId,optional"`
}
