# AI Development Instructions

This document provides instructions for AI coding assistants working on this codebase. Follow these rules exactly.

## Project Overview

**GoSaaS** is a full-stack SaaS boilerplate with:
- **Backend**: Go-Zero framework (Go 1.25+)
- **Frontend**: SvelteKit 2 + Svelte 5 + Tailwind v4
- **Database**: SQLite (standalone) or Levee (managed platform)
- **Payments**: Stripe

The frontend is statically built and embedded into a single Go binary.

## Architecture

```
├── gosaas.api                 # API definition (routes, request/response types)
├── gosaas.go                  # Entry point
├── etc/gosaas.yaml            # Backend config + Products/Pricing
├── internal/
│   ├── handler/               # AUTO-GENERATED - DO NOT EDIT
│   ├── types/                 # AUTO-GENERATED - DO NOT EDIT
│   ├── logic/                 # Business logic - EDIT HERE
│   ├── svc/                   # Service context
│   ├── db/                    # SQLite (standalone mode)
│   ├── local/                 # Local auth/billing (standalone mode)
│   ├── config/                # Config structs
│   └── middleware/            # CORS, JWT, rate limiting
├── app/                       # SvelteKit frontend
│   ├── src/
│   │   ├── routes/            # File-based routing
│   │   │   ├── (www)/         # Marketing pages (public)
│   │   │   ├── (auth)/        # Auth pages (login, register)
│   │   │   ├── (app)/         # App pages (authenticated)
│   │   │   └── (minimal)/     # Minimal layout pages
│   │   └── lib/
│   │       ├── config/site.ts # Branding/SEO - single source of truth
│   │       ├── api/           # AUTO-GENERATED TypeScript client
│   │       ├── stores/        # Svelte stores (auth, subscription)
│   │       ├── components/    # UI components
│   │       └── utils/         # Utilities (seo.ts, etc.)
│   └── static/                # Static assets
└── docs/                      # Documentation
```

## Critical Rules

### 1. Code Generation

**NEVER run `goctl` commands directly.** Always use:
```bash
make gen
```

This generates:
- Go handlers in `internal/handler/`
- Go types in `internal/types/`
- TypeScript API client in `app/src/lib/api/`

### 2. Hot Reloading

**Do NOT restart the server.** We use `air` for hot reloading:
```bash
make air
```

### 3. Package Manager

**Use pnpm only.** Never use npm or yarn:
```bash
cd app && pnpm install
cd app && pnpm dev
cd app && pnpm build
```

### 4. Styling

**All styles go in `app/src/app.css` or component-specific CSS files.**
- NEVER use inline styles
- NEVER use `<style>` blocks in Svelte files
- Use Tailwind v4 utility classes

### 5. Svelte 5 Syntax

Use Svelte 5 runes, NOT Svelte 4 syntax:

```svelte
<!-- CORRECT - Svelte 5 -->
<script lang="ts">
  let count = $state(0);
  let doubled = $derived(count * 2);
  let { name, age }: { name: string; age: number } = $props();

  $effect(() => {
    console.log('count changed:', count);
  });
</script>

<!-- WRONG - Svelte 4 -->
<script lang="ts">
  export let name;  // DON'T USE
  let count = 0;    // DON'T USE for reactive state
  $: doubled = count * 2;  // DON'T USE
</script>
```

For slots, use snippets:
```svelte
<!-- CORRECT - Svelte 5 -->
<script lang="ts">
  import type { Snippet } from 'svelte';
  let { children }: { children: Snippet } = $props();
</script>
{@render children()}

<!-- WRONG - Svelte 4 -->
<slot />
```

### 6. Go Code Style

**Idiomatic Go only.** One function with parameters, not multiple variations:

```go
// CORRECT
func Register(token string, opts ...Option) error

// WRONG - Don't create multiple functions
func Register() error
func RegisterWithToken(token string) error
```

### 7. Dual-Mode Support

All logic handlers must support both standalone and Levee modes:

```go
func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
    // Check mode first
    if l.svcCtx.UseLocal() {
        return l.loginLocal(req)
    }

    // Levee mode
    if l.svcCtx.Levee == nil {
        return nil, fmt.Errorf("auth service not configured")
    }
    // ... Levee implementation
}
```

Key methods:
- `l.svcCtx.UseLocal()` - Returns true for standalone mode (SQLite + Stripe)
- `l.svcCtx.UseLevee()` - Returns true for Levee mode
- `l.svcCtx.DB` - SQLite database (nil in Levee mode)
- `l.svcCtx.Levee` - Levee SDK client (nil in standalone mode)

### 8. Minimal Changes

- **NEVER remove code that appears unused** - It may be called from frontend or other services
- **NEVER add features beyond what's requested** - No over-engineering
- **NEVER add comments to code you didn't change**
- **Ask before deleting** - If something seems unused, ask first

### 9. Build Before Push

Always verify before committing:
```bash
make build
cd app && pnpm check && pnpm build
```

## Configuration System

Three config files, each with a specific purpose:

| File | Purpose | Contains |
|------|---------|----------|
| `app/src/lib/config/site.ts` | Frontend branding/SEO | Site name, tagline, social links, legal info |
| `etc/gosaas.yaml` | Backend config + pricing | Products, prices, features, backend settings |
| `.env` | Secrets only | API keys, JWT secret, database credentials |

### Branding Changes

Edit `app/src/lib/config/site.ts`:
```typescript
export const site = {
  name: 'YourApp',
  tagline: 'Your tagline',
  url: 'https://yourapp.com',
  // ...
}
```

### Pricing Changes

Edit `etc/gosaas.yaml`:
```yaml
Products:
  - slug: "pro"
    name: "Pro"
    prices:
      - slug: "monthly"
        amount: 2900  # cents
        interval: "month"
```

Pricing is read from YAML at build time and embedded in static HTML.

## Adding New Endpoints

1. Define in `gosaas.api`:
```
@server(
    prefix: /api/v1
    middleware: JwtAuth
)
service gosaas {
    @handler GetWidget
    get /widgets/:id (GetWidgetRequest) returns (GetWidgetResponse)
}

type GetWidgetRequest {
    Id string `path:"id"`
}

type GetWidgetResponse {
    Id   string `json:"id"`
    Name string `json:"name"`
}
```

2. Run `make gen`

3. Implement in `internal/logic/getwidgetlogic.go`

4. Frontend types auto-available in `$lib/api`

## File Locations

| Task | Location |
|------|----------|
| Add API endpoint | `gosaas.api` → `make gen` → `internal/logic/` |
| Add page | `app/src/routes/` |
| Add component | `app/src/lib/components/` |
| Add store | `app/src/lib/stores/` |
| Change branding | `app/src/lib/config/site.ts` |
| Change pricing | `etc/gosaas.yaml` |
| Add middleware | `internal/middleware/` |
| Database migrations | `internal/db/` |

## Route Groups

| Group | Purpose | Layout |
|-------|---------|--------|
| `(www)` | Marketing pages | Full header/footer |
| `(auth)` | Login, register, reset password | Minimal, centered |
| `(app)` | Authenticated app pages | App shell with sidebar |
| `(minimal)` | Email confirmations, errors | Minimal layout |

## Stores

| Store | Purpose |
|-------|---------|
| `auth` | Authentication state, login/logout functions |
| `subscription` | Subscription state, checkout functions |
| `currentUser` | Derived from auth store |
| `isAuthenticated` | Derived boolean |

## Environment Variables

### Required (Both Modes)
```bash
ACCESS_SECRET=xxx        # JWT secret (openssl rand -hex 32)
APP_BASE_URL=xxx         # Frontend URL
```

### Standalone Mode
```bash
LEVEE_ENABLED=false      # or unset
SQLITE_PATH=./data/gosaas.db
STRIPE_SECRET_KEY=sk_xxx
STRIPE_PUBLISHABLE_KEY=pk_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
```

### Levee Mode
```bash
LEVEE_ENABLED=true
LEVEE_API_KEY=lvk_xxx
LEVEE_BASE_URL=https://api.levee.sh
```

## Testing

```bash
# Go tests
go test -v ./internal/logic/...
go test -v -run TestFunctionName ./internal/logic/auth/

# Frontend tests
cd app && pnpm test:unit
cd app && pnpm test:unit -- MyTest
```

## Common Patterns

### API Client Usage (Frontend)
```typescript
import { login, getMe, createCheckout } from '$lib/api';

// Auth
const { token } = await login({ email, password });

// User
const user = await getMe();

// Subscription
const { checkoutUrl } = await createCheckout({ planName: 'pro' });
```

### Protected Routes
```svelte
<script lang="ts">
  import { isAuthenticated } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  $effect(() => {
    if (!$isAuthenticated) {
      goto('/auth/login');
    }
  });
</script>
```

### Form Handling
```svelte
<script lang="ts">
  let loading = $state(false);
  let error = $state('');

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    loading = true;
    error = '';
    try {
      await someApiCall();
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
</script>
```

## Do NOT

- Run `goctl` directly (use `make gen`)
- Use npm or yarn (use `pnpm`)
- Use Svelte 4 syntax (`export let`, `$:`, `<slot>`)
- Add inline styles or `<style>` blocks
- Create multiple function variations in Go
- Remove code without asking
- Over-engineer solutions
- Skip the build check before commits
- Create `.cursorrules`, `.windsurfrules`, or similar files
