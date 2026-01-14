# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**GoSaaS** - Ship your SaaS in 10 days. Full-stack boilerplate with Go-Zero, SvelteKit, and optional Levee integration.

Two modes of operation:

1. **Standalone Mode** (default): SQLite + direct Stripe integration - zero external dependencies
2. **Levee Mode**: Full platform features via Levee SDK (auth, billing, email, LLM)

Install via:

```bash
curl -fsSL https://raw.githubusercontent.com/almatuck/gosaas/main/install.sh | bash -s -- myapp
```

## Build & Development Commands

```bash
# Backend (Go-Zero API)
make build           # Build the API binary
make run             # Build and run the API (localhost:8888)
make test            # Run all Go tests
make air             # Run with hot reload

# Frontend (SvelteKit, in /app directory)
cd app && pnpm dev   # Start dev server (localhost:5173)
cd app && pnpm build # Production build
cd app && pnpm check # Type checking

# Code Generation (CRITICAL: Never run goctl directly!)
make gen             # Generate Go handlers/types + TypeScript API client from .api file

# Project Setup
make setup NAME=myapp  # Rename project from 'gosaas' to 'myapp'
make deps              # Download dependencies
make dev-setup         # Full dev environment setup
```

## Architecture Overview

**Backend (Go-Zero):**

- `<appname>.api` - API definition file (routes, request/response types)
- `<appname>.go` - Service entry point
- `etc/<appname>.yaml` - Configuration (uses env vars: `${VAR_NAME}`)
- `internal/handler/` - Auto-generated HTTP handlers (DO NOT EDIT DIRECTLY)
- `internal/types/` - Auto-generated request/response types (DO NOT EDIT DIRECTLY)
- `internal/logic/` - Business logic implementation (EDIT HERE for endpoint logic)
- `internal/svc/` - Service context (initializes either Levee or local services)
- `internal/db/` - SQLite database setup and migrations (standalone mode)
- `internal/local/` - Local auth and billing services (standalone mode)
- `internal/config/` - Configuration structs
- `internal/middleware/` - CORS, JWT, security headers, rate limiting

**Frontend (SvelteKit + Svelte 5 + Tailwind v4):**

- `app/src/routes/` - SvelteKit file-based routing
- `app/src/lib/api/` - Auto-generated TypeScript API client
- `app/src/lib/` - Shared components and utilities

## Dual-Mode Architecture

### Standalone Mode (LEVEE_ENABLED=false)

Default mode - no external services required.

- **Database**: SQLite with WAL mode (`internal/db/sqlite.go`)
- **Auth**: Local JWT auth with bcrypt passwords (`internal/local/auth.go`)
- **Billing**: Direct Stripe integration (`internal/local/billing.go`)
- **Email**: Not included (add SMTP if needed)

Tables: `users`, `user_preferences`, `refresh_tokens`, `subscriptions`, `leads`

### Levee Mode (LEVEE_ENABLED=true)

Full platform features via Levee SDK.

- **Auth**: Levee handles registration, login, password reset
- **Billing**: Levee manages Stripe subscriptions
- **Email**: Levee email lists and sequences
- **LLM**: Cost-tracked AI gateway

## Logic Handler Pattern

All handlers support both modes via `UseLocal()` check:

```go
func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
    // Use local auth when Levee is disabled
    if l.svcCtx.UseLocal() {
        return l.loginLocal(req)
    }

    // Use Levee when enabled
    if l.svcCtx.Levee == nil {
        return nil, fmt.Errorf("auth service not configured")
    }

    // Levee implementation...
}
```

Key service context methods:

- `l.svcCtx.UseLocal()` - Returns true if using SQLite + direct Stripe
- `l.svcCtx.UseLevee()` - Returns true if using Levee SDK
- `l.svcCtx.Auth` - Local auth service (nil if using Levee)
- `l.svcCtx.Billing` - Local billing service (nil if using Levee)
- `l.svcCtx.DB` - SQLite database (nil if using Levee)
- `l.svcCtx.Levee` - Levee SDK client (nil if using local)

## Environment Variables

**Required (both modes):**

- `ACCESS_SECRET` - JWT secret (generate: `openssl rand -hex 32`)
- `APP_BASE_URL` - Frontend URL for redirects

**Standalone Mode:**

- `LEVEE_ENABLED=false` (default)
- `SQLITE_PATH` - SQLite database path
- `STRIPE_SECRET_KEY` - Stripe API key
- `STRIPE_PUBLISHABLE_KEY` - Stripe publishable key
- `STRIPE_WEBHOOK_SECRET` - For Stripe webhooks
- `STRIPE_PRO_MONTHLY_PRICE_ID` - Stripe price ID for Pro plan

**Levee Mode:**

- `LEVEE_ENABLED=true`
- `LEVEE_API_KEY` - API key from Levee dashboard
- `LEVEE_BASE_URL` - Your Levee instance URL

## Frontend Structure

**SvelteKit Route Groups:**

- `(www)` - Marketing pages (landing, pricing, privacy, terms)
- `(auth)` - Authentication pages (login, register, reset-password)
- `(app)` - Authenticated app pages (dashboard, account)

**Stores (Svelte 5):**

- `auth` - Authentication state, login/register/logout functions
- `subscription` - Subscription state, checkout/billing functions
- `currentUser` - Derived from auth store
- `isAuthenticated` - Derived boolean from auth store

## Critical Rules

- **NEVER run goctl commands directly** - Always use `make gen`
- **Use air for hot reloading** - No need to restart the server during development
- **pnpm only** - Do not use npm or yarn for frontend
- **Styles in app.css only** - Never inline styles or use `<style>` blocks in Svelte files
- **Always build before git push** - Run `make build` and `cd app && pnpm build`
- **Idiomatic Go** - One function with parameters, not multiple function variations
- **Minimal changes** - Preserve existing functionality, never remove code that appears unused without asking
- **Svelte 5 runes** - Use `$state`, `$derived`, `$props`, `$effect` (not Svelte 4 syntax)
- **Support both modes** - When modifying logic handlers, ensure both Levee and local paths work

## Testing

```bash
# Go tests
go test -v -run TestFunctionName ./internal/logic/auth/

# Frontend tests (Vitest)
cd app && pnpm test:unit           # Run all tests
cd app && pnpm test:unit -- MyTest # Run specific test
```

## Deployment

**Single Binary:** The app compiles to a single binary with embedded frontend. Deploy the binary + `.env` file anywhere.

**Docker:** Use `docker compose up` for development, `deploy/compose.yaml` for production.

**Platforms:** Works on Fly.io, Railway, DigitalOcean, any VPS.
