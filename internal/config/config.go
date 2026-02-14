package config

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// parseBool parses a string as boolean with a default value.
// Accepts: "true", "1", "yes" as true; empty or other values return default.
func parseBool(s string, defaultVal bool) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return defaultVal
	}
	return s == "true" || s == "1" || s == "yes"
}

// LoadConfig unmarshals YAML data into a Config with defaults pre-applied.
func LoadConfig(data []byte) (*Config, error) {
	c := DefaultConfig()
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// DefaultConfig returns a Config with all default values set.
func DefaultConfig() Config {
	return Config{
		Name: "gosaas",
		Host: "0.0.0.0",
		Port: 8888,
		App: AppConfig{
			ProductionMode: "false",
		},
		Auth: AuthConfig{
			RefreshTokenExpire: 604800,
		},
		Database: DatabaseConfig{
			SQLitePath:      "./data/gosaas.db",
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			DBName:          "gosaas",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300,
		},
		Stripe: StripeConfig{
			SuccessURL: "/app/account",
			CancelURL:  "/pricing",
		},
		Security: SecurityConfig{
			CSRFEnabled:           "true",
			CSRFTokenExpiry:       43200,
			CSRFSecureCookie:      "true",
			RateLimitEnabled:      "true",
			RateLimitRequests:     100,
			RateLimitInterval:     60,
			RateLimitBurst:        20,
			AuthRateLimitRequests: 5,
			AuthRateLimitInterval: 60,
			EnableSecurityHeaders: "true",
			ForceHTTPS:            "false",
			MaxRequestBodySize:    10485760,
			MaxURLLength:          2048,
		},
		Email: EmailConfig{
			SMTPPort: 587,
			FromName: "gosaas",
			BaseURL:  "http://localhost:5173",
		},
		Analytics: AnalyticsConfig{
			Enabled:       "true",
			Provider:      "console",
			BatchSize:     50,
			FlushInterval: 30,
			Debug:         "false",
		},
		Subscription: SubscriptionConfig{
			Enabled:             "true",
			EnforceQuotas:       "true",
			FreeTierAnalyses:    5,
			FreeTierHistoryDays: 7,
			ProTierHistoryDays:  30,
			TeamTierHistoryDays: -1,
		},
		AI: AIConfig{
			Enabled:     "true",
			Model:       "claude-sonnet-4-5-20250929",
			MaxTokens:   4096,
			TimeoutSecs: 60,
		},
		CostAlerts: CostAlertsConfig{
			Enabled:            "true",
			DailyCostThreshold: 10.0,
			UserCostThreshold:  1.0,
		},
		Outlet: OutletConfig{
			WaitlistListSlug:          "waitlist",
			NewsletterListSlug:        "newsletter",
			UsersListSlug:             "users",
			OnboardingSequence:        "onboarding",
			TrialConversionSequence:   "trial-conversion",
			EmailVerificationTemplate: "email-verification",
			PasswordResetTemplate:     "password-reset",
			WelcomeTemplate:           "welcome",
		},
		Levee: LeveeConfig{
			Enabled:                           "true",
			CheckoutSuccessURL:                "/app/billing/success",
			CheckoutCancelURL:                 "/pricing",
			WaitlistListSlug:                  "waitlist",
			NewsletterListSlug:                "newsletter-subscribers",
			UsersListSlug:                     "users",
			ProductHuntListSlug:               "product-hunt-launchers",
			WaitlistNurtureSequence:           "waitlist-nurture",
			OnboardingSequence:                "onboarding-sequence",
			TrialConversionSequence:           "free-to-pro-upgrade",
			FreeProductSlug:                   "free",
			ProMonthlyProductSlug:             "pro",
			ProYearlyProductSlug:              "pro-yearly",
			TeamMonthlyProductSlug:            "team",
			TeamYearlyProductSlug:             "team-yearly",
			EmailVerificationTemplate:         "email-verification",
			PasswordResetTemplate:             "password-reset",
			PasswordChangedTemplate:           "password-changed",
			AccountDeletedTemplate:            "account-deleted",
			WaitlistConfirmTemplate:           "waitlist-confirm",
			SubscriptionWelcomeTemplate:       "subscription-welcome",
			SubscriptionRenewalTemplate:       "subscription-renewal",
			SubscriptionCancelledTemplate:     "subscription-cancelled",
			PaymentFailedTemplate:             "payment-failed",
			TrialEndingTemplate:               "trial-ending",
			CostAlertTemplate:                 "cost-alert",
		},
		OAuth: OAuthConfig{
			GoogleEnabled: "false",
			GitHubEnabled: "false",
		},
		Features: FeaturesConfig{
			OrganizationsEnabled: "true",
			NotificationsEnabled: "true",
			OAuthEnabled:         "false",
		},
	}
}

type Config struct {
	Name     string         `yaml:"Name"`
	Host     string         `yaml:"Host"`
	Port     int            `yaml:"Port"`
	App      AppConfig      `yaml:"App"`
	Admin    AdminConfig    `yaml:"Admin"`
	Auth     AuthConfig     `yaml:"Auth"`
	Database DatabaseConfig `yaml:"Database"`
	Stripe   StripeConfig   `yaml:"Stripe"`
	Products []Product      `yaml:"Products"`
	Security SecurityConfig `yaml:"Security"`
	Email    EmailConfig    `yaml:"Email"`
	Analytics    AnalyticsConfig    `yaml:"Analytics"`
	Subscription SubscriptionConfig `yaml:"Subscription"`
	AI           AIConfig           `yaml:"AI"`
	CostAlerts   CostAlertsConfig   `yaml:"CostAlerts"`
	Outlet       OutletConfig       `yaml:"Outlet"`
	Levee        LeveeConfig        `yaml:"Levee"`
	OAuth        OAuthConfig        `yaml:"OAuth"`
	Features     FeaturesConfig     `yaml:"Features"`
}

type AppConfig struct {
	BaseURL        string `yaml:"BaseURL"`
	Domain         string `yaml:"Domain"`
	ProductionMode string `yaml:"ProductionMode"`
	AdminEmail     string `yaml:"AdminEmail"`
}

type AdminConfig struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

type AuthConfig struct {
	AccessSecret       string `yaml:"AccessSecret"`
	AccessExpire       int64  `yaml:"AccessExpire"`
	RefreshTokenExpire int64  `yaml:"RefreshTokenExpire"`
}

type DatabaseConfig struct {
	SQLitePath      string `yaml:"SQLitePath"`
	Host            string `yaml:"Host"`
	Port            int    `yaml:"Port"`
	User            string `yaml:"User"`
	Password        string `yaml:"Password"`
	DBName          string `yaml:"DBName"`
	SSLMode         string `yaml:"SSLMode"`
	MaxOpenConns    int    `yaml:"MaxOpenConns"`
	MaxIdleConns    int    `yaml:"MaxIdleConns"`
	ConnMaxLifetime int    `yaml:"ConnMaxLifetime"`
}

type StripeConfig struct {
	SecretKey      string `yaml:"SecretKey"`
	PublishableKey string `yaml:"PublishableKey"`
	WebhookSecret  string `yaml:"WebhookSecret"`
	SuccessURL     string `yaml:"SuccessURL"`
	CancelURL      string `yaml:"CancelURL"`
}

type SecurityConfig struct {
	CSRFEnabled           string `yaml:"CSRFEnabled"`
	CSRFSecret            string `yaml:"CSRFSecret"`
	CSRFTokenExpiry       int64  `yaml:"CSRFTokenExpiry"`
	CSRFSecureCookie      string `yaml:"CSRFSecureCookie"`
	RateLimitEnabled      string `yaml:"RateLimitEnabled"`
	RateLimitRequests     int    `yaml:"RateLimitRequests"`
	RateLimitInterval     int    `yaml:"RateLimitInterval"`
	RateLimitBurst        int    `yaml:"RateLimitBurst"`
	AuthRateLimitRequests int    `yaml:"AuthRateLimitRequests"`
	AuthRateLimitInterval int    `yaml:"AuthRateLimitInterval"`
	EnableSecurityHeaders string `yaml:"EnableSecurityHeaders"`
	ContentSecurityPolicy string `yaml:"ContentSecurityPolicy"`
	AllowedOrigins        string `yaml:"AllowedOrigins"`
	ForceHTTPS            string `yaml:"ForceHTTPS"`
	MaxRequestBodySize    int64  `yaml:"MaxRequestBodySize"`
	MaxURLLength          int    `yaml:"MaxURLLength"`
}

type EmailConfig struct {
	SMTPHost    string `yaml:"SMTPHost"`
	SMTPPort    int    `yaml:"SMTPPort"`
	SMTPUser    string `yaml:"SMTPUser"`
	SMTPPass    string `yaml:"SMTPPass"`
	FromAddress string `yaml:"FromAddress"`
	FromName    string `yaml:"FromName"`
	ReplyTo     string `yaml:"ReplyTo"`
	BaseURL     string `yaml:"BaseURL"`
}

type AnalyticsConfig struct {
	Enabled       string `yaml:"Enabled"`
	Provider      string `yaml:"Provider"`
	APIKey        string `yaml:"APIKey"`
	Endpoint      string `yaml:"Endpoint"`
	BatchSize     int    `yaml:"BatchSize"`
	FlushInterval int    `yaml:"FlushInterval"`
	Debug         string `yaml:"Debug"`
}

type SubscriptionConfig struct {
	Enabled             string `yaml:"Enabled"`
	EnforceQuotas       string `yaml:"EnforceQuotas"`
	FreeTierAnalyses    int    `yaml:"FreeTierAnalyses"`
	FreeTierHistoryDays int    `yaml:"FreeTierHistoryDays"`
	ProTierHistoryDays  int    `yaml:"ProTierHistoryDays"`
	TeamTierHistoryDays int    `yaml:"TeamTierHistoryDays"`
}

type AIConfig struct {
	Enabled     string `yaml:"Enabled"`
	APIKey      string `yaml:"APIKey"`
	Model       string `yaml:"Model"`
	MaxTokens   int    `yaml:"MaxTokens"`
	TimeoutSecs int    `yaml:"TimeoutSecs"`
}

type CostAlertsConfig struct {
	Enabled            string  `yaml:"Enabled"`
	AdminEmail         string  `yaml:"AdminEmail"`
	DailyCostThreshold float64 `yaml:"DailyCostThreshold"`
	UserCostThreshold  float64 `yaml:"UserCostThreshold"`
}

type OutletConfig struct {
	BaseURL                   string `yaml:"BaseURL"`
	APIKey                    string `yaml:"APIKey"`
	WaitlistListSlug          string `yaml:"WaitlistListSlug"`
	NewsletterListSlug        string `yaml:"NewsletterListSlug"`
	UsersListSlug             string `yaml:"UsersListSlug"`
	OnboardingSequence        string `yaml:"OnboardingSequence"`
	TrialConversionSequence   string `yaml:"TrialConversionSequence"`
	EmailVerificationTemplate string `yaml:"EmailVerificationTemplate"`
	PasswordResetTemplate     string `yaml:"PasswordResetTemplate"`
	WelcomeTemplate           string `yaml:"WelcomeTemplate"`
}

type LeveeConfig struct {
	APIKey      string `yaml:"APIKey"`
	BaseURL     string `yaml:"BaseURL"`
	GRPCAddress string `yaml:"GRPCAddress"`
	Enabled     string `yaml:"Enabled"`

	CheckoutSuccessURL string `yaml:"CheckoutSuccessURL"`
	CheckoutCancelURL  string `yaml:"CheckoutCancelURL"`

	WaitlistListSlug    string `yaml:"WaitlistListSlug"`
	NewsletterListSlug  string `yaml:"NewsletterListSlug"`
	UsersListSlug       string `yaml:"UsersListSlug"`
	ProductHuntListSlug string `yaml:"ProductHuntListSlug"`

	WaitlistNurtureSequence string `yaml:"WaitlistNurtureSequence"`
	OnboardingSequence      string `yaml:"OnboardingSequence"`
	TrialConversionSequence string `yaml:"TrialConversionSequence"`

	FreeProductSlug        string `yaml:"FreeProductSlug"`
	ProMonthlyProductSlug  string `yaml:"ProMonthlyProductSlug"`
	ProYearlyProductSlug   string `yaml:"ProYearlyProductSlug"`
	TeamMonthlyProductSlug string `yaml:"TeamMonthlyProductSlug"`
	TeamYearlyProductSlug  string `yaml:"TeamYearlyProductSlug"`

	EmailVerificationTemplate     string `yaml:"EmailVerificationTemplate"`
	PasswordResetTemplate         string `yaml:"PasswordResetTemplate"`
	PasswordChangedTemplate       string `yaml:"PasswordChangedTemplate"`
	AccountDeletedTemplate        string `yaml:"AccountDeletedTemplate"`
	WaitlistConfirmTemplate       string `yaml:"WaitlistConfirmTemplate"`
	SubscriptionWelcomeTemplate   string `yaml:"SubscriptionWelcomeTemplate"`
	SubscriptionRenewalTemplate   string `yaml:"SubscriptionRenewalTemplate"`
	SubscriptionCancelledTemplate string `yaml:"SubscriptionCancelledTemplate"`
	PaymentFailedTemplate         string `yaml:"PaymentFailedTemplate"`
	TrialEndingTemplate           string `yaml:"TrialEndingTemplate"`
	CostAlertTemplate             string `yaml:"CostAlertTemplate"`
}

type OAuthConfig struct {
	GoogleEnabled      string `yaml:"GoogleEnabled"`
	GoogleClientID     string `yaml:"GoogleClientID"`
	GoogleClientSecret string `yaml:"GoogleClientSecret"`
	GitHubEnabled      string `yaml:"GitHubEnabled"`
	GitHubClientID     string `yaml:"GitHubClientID"`
	GitHubClientSecret string `yaml:"GitHubClientSecret"`
	CallbackBaseURL    string `yaml:"CallbackBaseURL"`
}

type FeaturesConfig struct {
	OrganizationsEnabled string `yaml:"OrganizationsEnabled"`
	NotificationsEnabled string `yaml:"NotificationsEnabled"`
	OAuthEnabled         string `yaml:"OAuthEnabled"`
}

// ========== Helper Methods ==========

// IsProductionMode returns true if production mode is enabled.
func (c Config) IsProductionMode() bool {
	return parseBool(c.App.ProductionMode, false)
}

// IsCSRFEnabled returns true if CSRF protection is enabled.
func (c Config) IsCSRFEnabled() bool {
	return parseBool(c.Security.CSRFEnabled, true)
}

// IsCSRFSecureCookie returns true if CSRF cookies should be secure.
func (c Config) IsCSRFSecureCookie() bool {
	return parseBool(c.Security.CSRFSecureCookie, true)
}

// IsRateLimitEnabled returns true if rate limiting is enabled.
func (c Config) IsRateLimitEnabled() bool {
	return parseBool(c.Security.RateLimitEnabled, true)
}

// IsSecurityHeadersEnabled returns true if security headers are enabled.
func (c Config) IsSecurityHeadersEnabled() bool {
	return parseBool(c.Security.EnableSecurityHeaders, true)
}

// IsForceHTTPS returns true if HTTPS should be enforced.
func (c Config) IsForceHTTPS() bool {
	return parseBool(c.Security.ForceHTTPS, false)
}

// IsAnalyticsEnabled returns true if analytics is enabled.
func (c Config) IsAnalyticsEnabled() bool {
	return parseBool(c.Analytics.Enabled, true)
}

// IsAnalyticsDebug returns true if analytics debug mode is enabled.
func (c Config) IsAnalyticsDebug() bool {
	return parseBool(c.Analytics.Debug, false)
}

// IsSubscriptionEnabled returns true if subscription features are enabled.
func (c Config) IsSubscriptionEnabled() bool {
	return parseBool(c.Subscription.Enabled, true)
}

// IsEnforceQuotas returns true if usage quotas should be enforced.
func (c Config) IsEnforceQuotas() bool {
	return parseBool(c.Subscription.EnforceQuotas, true)
}

// IsAIEnabled returns true if AI features are enabled.
func (c Config) IsAIEnabled() bool {
	return parseBool(c.AI.Enabled, true)
}

// IsCostAlertsEnabled returns true if cost alerts are enabled.
func (c Config) IsCostAlertsEnabled() bool {
	return parseBool(c.CostAlerts.Enabled, true)
}

// IsLeveeEnabled returns true if Levee integration is enabled.
func (c Config) IsLeveeEnabled() bool {
	return parseBool(c.Levee.Enabled, true)
}

// IsGoogleOAuthEnabled returns true if Google OAuth is enabled.
func (c Config) IsGoogleOAuthEnabled() bool {
	return parseBool(c.OAuth.GoogleEnabled, false)
}

// IsGitHubOAuthEnabled returns true if GitHub OAuth is enabled.
func (c Config) IsGitHubOAuthEnabled() bool {
	return parseBool(c.OAuth.GitHubEnabled, false)
}

// IsOrganizationsEnabled returns true if organizations feature is enabled.
func (c Config) IsOrganizationsEnabled() bool {
	return parseBool(c.Features.OrganizationsEnabled, true)
}

// IsNotificationsEnabled returns true if notifications feature is enabled.
func (c Config) IsNotificationsEnabled() bool {
	return parseBool(c.Features.NotificationsEnabled, true)
}

// IsOAuthEnabled returns true if OAuth feature is enabled.
func (c Config) IsOAuthEnabled() bool {
	return parseBool(c.Features.OAuthEnabled, false)
}

// Product defines a subscription product with its prices.
type Product struct {
	Slug            string   `yaml:"slug" json:"slug"`
	Name            string   `yaml:"name" json:"name"`
	Description     string   `yaml:"description" json:"description,omitempty"`
	Features        []string `yaml:"features" json:"features,omitempty"`
	Prices          []Price  `yaml:"prices" json:"prices,omitempty"`
	Default         bool     `yaml:"default" json:"default,omitempty"`
	StripeProductID string   `yaml:"stripeProductId" json:"stripeProductId,omitempty"`
}

// Price defines a single pricing option for a product.
type Price struct {
	Slug          string `yaml:"slug" json:"slug"`
	Amount        int64  `yaml:"amount" json:"amount"`
	Currency      string `yaml:"currency" json:"currency"`
	Interval      string `yaml:"interval" json:"interval"`
	IntervalCount int64  `yaml:"intervalCount" json:"intervalCount"`
	TrialDays     int64  `yaml:"trialDays" json:"trialDays,omitempty"`
	StripePriceID string `yaml:"stripePriceId" json:"stripePriceId,omitempty"`
}
