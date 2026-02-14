# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Reference

```bash
# Development (hot reload - no restart needed)
make air              # Backend with hot reload
cd app && pnpm dev    # Frontend dev server

# Code generation (TypeScript only)
make gen              # Regenerate TypeScript API client from Go types/routes

# Testing
go test -v ./internal/handler/...                       # All Go tests
go test -v -run TestName ./internal/handler/auth/      # Single test
cd app && pnpm check                                   # TypeScript check
cd app && pnpm test:unit                               # Frontend tests

# Database (standalone mode)
make migrate-up       # Run pending migrations
make migrate-down     # Rollback last migration
make migrate-status   # Check migration status

# Before committing
make build && cd app && pnpm build
```

## Architecture

**Dual-mode system** - switches between standalone (SQLite + Stripe) and Levee (managed platform):

```
internal/types/types.go      → Request/response types (source of truth)
internal/handler/routes.go   → Route registrations (source of truth)
internal/handler/{group}/    → Handler + business logic (one file per endpoint)
internal/svc/                → ServiceContext with UseLocal()/UseLevee() mode check
internal/db/                 → SQLite (standalone mode)
internal/local/              → Local auth/billing services (standalone mode)
cmd/genapi/                  → TypeScript generator (reads types + routes, writes TS)

app/src/
├── routes/(www)/            → Marketing pages (public)
├── routes/(auth)/           → Auth pages (login, register)
├── routes/(app)/            → App pages (authenticated)
├── lib/api/                 → AUTO-GENERATED TypeScript client (gosaas.ts, gosaasComponents.ts)
├── lib/config/site.ts       → Branding/SEO (single source of truth)
└── lib/stores/              → Svelte stores (auth, subscription)
```

## Adding API Endpoints

1. Add types to `internal/types/types.go`:
```go
type GetWidgetRequest struct {
    Id string `path:"id"`
}
type GetWidgetResponse struct {
    Name string `json:"name"`
}
```

2. Add route to `internal/handler/routes.go` (inside the `r.Route("/api/v1", ...)` block):
```go
r.Get("/widgets/{id}", widget.GetWidgetHandler(svcCtx))
```

3. Create handler file `internal/handler/widget/getwidget.go`:
```go
package widget

import (
    "context"
    "log/slog"
    "net/http"
    "gosaas/internal/httpx"
    "gosaas/internal/svc"
    "gosaas/internal/types"
)

type GetWidgetLogic struct {
    logger *slog.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func GetWidgetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.GetWidgetRequest
        if err := httpx.Parse(r, &req); err != nil {
            httpx.ErrorResponse(w, err)
            return
        }
        l := &GetWidgetLogic{logger: slog.Default(), ctx: r.Context(), svcCtx: svcCtx}
        resp, err := l.GetWidget(&req)
        if err != nil {
            httpx.ErrorResponse(w, err)
        } else {
            httpx.OkJson(w, resp)
        }
    }
}

func (l *GetWidgetLogic) GetWidget(req *types.GetWidgetRequest) (*types.GetWidgetResponse, error) {
    if l.svcCtx.UseLocal() {
        // SQLite implementation
    }
    // Levee implementation
}
```

4. Run `make gen` to regenerate TypeScript

5. Frontend types auto-available: `import { getWidget } from '$lib/api'`

## Mode-Aware Logic Pattern

All handlers must support both modes:

```go
func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
    if l.svcCtx.UseLocal() {
        return l.loginLocal(req)  // SQLite + local JWT
    }
    if l.svcCtx.Levee == nil {
        return nil, fmt.Errorf("auth service not configured")
    }
    // Levee SDK implementation
}
```

Key methods: `l.svcCtx.UseLocal()`, `l.svcCtx.UseLevee()`, `l.svcCtx.DB`, `l.svcCtx.Auth`, `l.svcCtx.Billing`, `l.svcCtx.Levee`

## Critical Rules

- **`make gen` for TypeScript only** - Regenerates TS from Go types/routes. No goctl.
- **pnpm only** - Never npm or yarn
- **Styles in app.css only** - No inline styles or `<style>` blocks
- **Svelte 5 runes** - `$state`, `$derived`, `$props`, `$effect` (not Svelte 4 `export let`, `$:`, `<slot>`)
- **DaisyUI components** - Use DaisyUI classes for UI components (btn, card, modal, etc.)
- **Idiomatic Go** - One function with parameters, not multiple variations
- **Minimal changes** - Never remove code that appears unused without asking
- **Support both modes** - Logic handlers must work with UseLocal() and UseLevee()

## Configuration

| File | Purpose |
|------|---------|
| `app/src/lib/config/site.ts` | Branding, SEO, social links |
| `etc/gosaas.yaml` | Products, pricing, backend settings |
| `.env` | Secrets only (API keys, JWT secret) |

## /init Flow

When user runs `/init`, follow the interactive setup in `AI.md`:
1. Environment setup (install.sh)
2. Business discovery - ask what they want to build
3. Research (optional) - run 14-step validation from `RESEARCH-PLAN.md`, outputs to `./plan/`
4. Auto-customize - update site.ts, landing page, theme (app.css), pricing (gosaas.yaml)
5. Verify and launch - `make build && cd app && pnpm build`

## Admin Backoffice

Access admin dashboard at `/admin` (requires `ADMIN_USERNAME` and `ADMIN_PASSWORD` from .env).
Admin API routes use JWT + basic auth middleware (`internal/middleware/adminauth.go`).
