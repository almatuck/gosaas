# Quick Start Guide

Launch your SaaS in under 10 minutes.

## The Fastest Way (One Command)

```bash
curl -fsSL https://raw.githubusercontent.com/almatuck/gosaas/main/install.sh | bash -s -- myapp
```

This single command:
1. Checks you have Go, Node.js, and pnpm installed
2. Clones the repository
3. Renames everything from `gosaas` to `myapp`
4. Generates a secure JWT secret
5. Creates `.env` with working defaults
6. Installs all dependencies
7. Offers to start the app

**Then visit http://localhost:5173** - you're done!

## Prerequisites

The installer checks for these automatically:

- **Go 1.25+** - [Download](https://go.dev/dl/)
- **Node.js 20+** - [Download](https://nodejs.org/)
- **pnpm** - Auto-installed if npm is available

## Alternative: Manual Setup

If you prefer to clone manually:

```bash
git clone https://github.com/almatuck/gosaas.git myapp
cd myapp
./install.sh myapp
make dev
```

The install script auto-detects whether it's running remotely (via curl) or locally (after clone) and does the right thing.

## Two Modes of Operation

GoSaaS works out of the box in **Standalone Mode**. No configuration needed to start developing.

### Standalone Mode (Default)

- **Database**: SQLite (auto-created)
- **Auth**: Built-in JWT authentication
- **Billing**: Direct Stripe integration
- **Zero external dependencies**

This is set up automatically by the installer.

### Levee Mode (Optional)

For advanced features (email sequences, CMS, AI gateway), edit `.env`:

```bash
LEVEE_ENABLED=true
LEVEE_API_KEY=lvk_your_api_key_here  # Get from https://dashboard.levee.sh
LEVEE_BASE_URL=https://api.levee.sh
```

See [Levee Integration Guide](./LEVEE_INTEGRATION.md) for details.

---

## Customizing Your App

### Step 1: Configure Your Brand (5 minutes)

Edit `app/src/lib/config/site.ts`:

```typescript
export const site = {
  name: 'MyApp',
  tagline: 'Your awesome tagline',
  description: 'A brief description for SEO',
  url: 'https://myapp.com',  // Your production URL
  supportEmail: 'support@myapp.com',
  ogImage: '/images/og-default.png',
  twitter: '@myapp',
  social: {
    twitter: 'https://twitter.com/myapp',
    github: 'https://github.com/myapp',
    linkedin: '',
    discord: ''
  },
  legal: {
    companyName: 'MyApp, Inc.',
    companyAddress: '123 Main St, City, State 12345',
    privacyUrl: '/privacy',
    termsUrl: '/terms'
  }
}
```

### Step 2: Configure Products & Pricing

Edit `etc/gosaas.yaml` (or `etc/myapp.yaml` after renaming):

```yaml
Products:
  - slug: "free"
    name: "Free"
    description: "Get started for free"
    default: true
    features:
      - "5 projects"
      - "Basic analytics"
      - "Community support"
    prices:
      - slug: "default"
        amount: 0
        currency: "usd"
        interval: "month"

  - slug: "pro"
    name: "Pro"
    description: "For professionals"
    features:
      - "Unlimited projects"
      - "Advanced analytics"
      - "Priority support"
    prices:
      - slug: "monthly"
        amount: 2900  # $29.00
        currency: "usd"
        interval: "month"
        trialDays: 14
      - slug: "yearly"
        amount: 29000  # $290.00
        currency: "usd"
        interval: "year"
        trialDays: 14
```

### Step 3: Set Up Stripe Webhooks (For Payments)

### Development

```bash
# Install Stripe CLI
brew install stripe/stripe-cli/stripe

# Login
stripe login

# Forward webhooks (run in separate terminal)
stripe listen --forward-to localhost:8888/api/webhook/stripe
```

Copy the webhook secret (`whsec_xxx`) to your `.env`:

```bash
STRIPE_WEBHOOK_SECRET=whsec_xxx
```

### Production

1. Go to [Stripe Dashboard → Webhooks](https://dashboard.stripe.com/webhooks)
2. Add endpoint: `https://yourdomain.com/api/webhook/stripe`
3. Select events: `checkout.session.completed`, `customer.subscription.*`, `invoice.*`
4. Copy signing secret to production `.env`

---

## Running the App

### Easiest Way

```bash
make dev
```

This starts both backend and frontend using Docker (or natively if Docker isn't available).

### Manual (Two Terminals)

**Terminal 1 - Backend:**
```bash
make air
```

**Terminal 2 - Frontend:**
```bash
cd app && pnpm dev
```

---

## Verify Everything Works

1. Open http://localhost:5173
2. Click "Get Started Free" to register
3. Check the dashboard at http://localhost:5173/app
4. Try the pricing page at http://localhost:5173/pricing
5. Test checkout flow (use Stripe test card: `4242 4242 4242 4242`)

---

## Customize the UI

### Colors & Theme

Edit `app/src/app.css`:

```css
:root {
  --color-primary: #your-brand-color;
  --color-accent: #your-accent-color;
}
```

### Logo & Favicon

Replace files in `app/static/`:
- `favicon.svg` - Browser tab icon
- `logo.svg` - Site logo

### Landing Page

Edit `app/src/routes/(www)/+page.svelte`:
- Update hero section
- Modify feature cards
- Customize FAQ

## What's Next?

- [Configuration Guide](./CONFIGURATION.md) - All settings explained
- [Deployment Guide](./DEPLOYMENT.md) - Production deployment with auto-SSL
- [Stripe Integration](./STRIPE.md) - Payment setup
- [Levee Integration](./LEVEE_INTEGRATION.md) - Platform features (optional)
- [Customization Guide](./CUSTOMIZATION.md) - Theming and components

## Troubleshooting

### Build Errors

```bash
# Clear caches and rebuild
cd app
rm -rf node_modules .svelte-kit
pnpm install
pnpm build
```

### Port Already in Use

```bash
# Check ports
lsof -i :8888
lsof -i :5173

# Kill process or change ports in .env
```

### Stripe Webhooks Not Working

1. Ensure Stripe CLI is running: `stripe listen --forward-to localhost:8888/api/webhook/stripe`
2. Check `STRIPE_WEBHOOK_SECRET` matches the CLI output
3. Verify endpoint URL is correct

### Database Issues

```bash
# Reset SQLite database
rm -f ./data/myapp.db
make air  # Will recreate on startup
```
