#!/bin/bash
#
# GoSaaS Installer & Setup Script
# ================================
#
# This single script handles both remote installation and local setup.
#
# REMOTE INSTALLATION (via curl):
#   curl -fsSL https://raw.githubusercontent.com/almatuck/gosaas/main/install.sh | bash -s -- myapp
#
# LOCAL SETUP (after git clone):
#   ./install.sh myapp
#
# The script auto-detects which mode to run based on context.
#

set -e

# ============================================================
# Configuration
# ============================================================
OLD_NAME="gosaas"
OLD_NAME_PROPER="GoSaaS"
VERSION="1.0.0"
REPO_URL="https://github.com/almatuck/gosaas.git"

# ============================================================
# Colors & Output Helpers
# ============================================================
if [[ -t 1 ]]; then
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  YELLOW='\033[1;33m'
  BLUE='\033[0;34m'
  BOLD='\033[1m'
  NC='\033[0m'
else
  RED='' GREEN='' YELLOW='' BLUE='' BOLD='' NC=''
fi

print_step()    { echo -e "${BLUE}▶${NC} $1"; }
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_warning() { echo -e "${YELLOW}⚠${NC} $1"; }
print_error()   { echo -e "${RED}✗${NC} $1"; }
print_info()    { echo -e "  $1"; }

die() { print_error "$1"; exit 1; }

# ============================================================
# Help
# ============================================================
show_help() {
  cat << EOF
${BOLD}GoSaaS Installer v${VERSION}${NC}

${BOLD}REMOTE INSTALLATION:${NC}
  curl -fsSL https://raw.githubusercontent.com/almatuck/gosaas/main/install.sh | bash -s -- <app-name>

${BOLD}LOCAL SETUP (after git clone):${NC}
  ./install.sh <app-name> [options]

${BOLD}ARGUMENTS:${NC}
  app-name      Your app name (lowercase, letters and numbers only)
                Example: myapp, coolsaas, app123

${BOLD}OPTIONS:${NC}
  --dir PATH    Install directory (remote mode only, default: current dir)
  --no-deps     Skip installing dependencies
  --no-env      Skip creating .env file
  --no-start    Don't offer to start the app after setup
  --force       Overwrite existing .env file
  --help        Show this help message

${BOLD}EXAMPLES:${NC}
  # Remote: Install in ./myapp
  curl ... | bash -s -- myapp

  # Remote: Install in specific directory
  curl ... | bash -s -- myapp --dir ~/projects

  # Local: Set up after manual clone
  git clone https://github.com/almatuck/gosaas.git myapp
  cd myapp
  ./install.sh myapp

  # Local: Just rename, skip deps
  ./install.sh myapp --no-deps

${BOLD}WHAT IT DOES:${NC}
  1. Checks prerequisites (Go, Node.js, pnpm)
  2. Clones repository (remote mode only)
  3. Renames project from 'gosaas' to your app name
  4. Generates secure JWT secret
  5. Creates .env with working defaults
  6. Installs all dependencies
  7. Initializes git repository

EOF
  exit 0
}

# ============================================================
# Parse Arguments
# ============================================================
NEW_NAME=""
INSTALL_DIR=""
INSTALL_DEPS=true
CREATE_ENV=true
FORCE_ENV=false
OFFER_START=true

while [[ $# -gt 0 ]]; do
  case $1 in
    --help|-h)
      show_help
      ;;
    --dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --no-deps)
      INSTALL_DEPS=false
      shift
      ;;
    --no-env)
      CREATE_ENV=false
      shift
      ;;
    --no-start)
      OFFER_START=false
      shift
      ;;
    --force)
      FORCE_ENV=true
      shift
      ;;
    -*)
      die "Unknown option: $1 (use --help for usage)"
      ;;
    *)
      if [[ -z "$NEW_NAME" ]]; then
        NEW_NAME="$1"
      else
        die "Too many arguments. Use --help for usage."
      fi
      shift
      ;;
  esac
done

# ============================================================
# Detect Mode: Remote (curl) vs Local (already cloned)
# ============================================================
# We're in "local mode" if we can find gosaas project files
if [[ -f "go.mod" ]] && [[ -f "${OLD_NAME}.go" || -f "${OLD_NAME}.api" || -d "internal" ]]; then
  MODE="local"
else
  MODE="remote"
fi

# ============================================================
# Validation
# ============================================================
if [[ -z "$NEW_NAME" ]]; then
  echo ""
  echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║${NC}              🚀 GoSaaS - Ship Your SaaS Fast              ${GREEN}║${NC}"
  echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
  echo ""
  echo -e "${BOLD}Usage:${NC}"
  echo ""
  if [[ "$MODE" == "remote" ]]; then
    echo "  curl -fsSL https://raw.githubusercontent.com/almatuck/gosaas/main/install.sh | bash -s -- <app-name>"
  else
    echo "  ./install.sh <app-name>"
  fi
  echo ""
  echo -e "${BOLD}Example:${NC}"
  echo "  ${MODE == 'remote' && 'curl ... | bash -s -- myapp' || './install.sh myapp'}"
  echo ""
  echo "Run with --help for more options."
  exit 1
fi

# Validate app name
if [[ ! "$NEW_NAME" =~ ^[a-z][a-z0-9]*$ ]]; then
  die "Invalid app name: '$NEW_NAME'

App name must:
  • Start with a lowercase letter
  • Contain only lowercase letters and numbers
  • Be a valid Go package name

Examples: myapp, coolsaas, app123"
fi

# ============================================================
# OS Detection for sed compatibility
# ============================================================
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed_inplace() { sed -i '' "$@"; }
else
  sed_inplace() { sed -i "$@"; }
fi

# ============================================================
# Start
# ============================================================
echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║${NC}              🚀 GoSaaS - Ship Your SaaS Fast              ${GREEN}║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  ${BOLD}App name:${NC} $NEW_NAME"
echo -e "  ${BOLD}Mode:${NC}     $MODE"
echo ""

# ============================================================
# Step 1: Check Prerequisites
# ============================================================
print_step "Checking prerequisites..."

MISSING=""
WARNINGS=""

# Check Go
if command -v go &> /dev/null; then
  GO_VERSION=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
  print_success "Go $GO_VERSION"
else
  MISSING="$MISSING Go"
fi

# Check Node.js
if command -v node &> /dev/null; then
  NODE_VERSION=$(node --version 2>/dev/null)
  print_success "Node.js $NODE_VERSION"
else
  MISSING="$MISSING Node.js"
fi

# Check pnpm (try to install if missing)
if command -v pnpm &> /dev/null; then
  PNPM_VERSION=$(pnpm --version 2>/dev/null)
  print_success "pnpm $PNPM_VERSION"
elif command -v npm &> /dev/null; then
  print_warning "pnpm not found - will install via npm"
  WARNINGS="$WARNINGS pnpm"
else
  MISSING="$MISSING pnpm"
fi

# Check git
if command -v git &> /dev/null; then
  print_success "git $(git --version 2>/dev/null | awk '{print $3}')"
else
  MISSING="$MISSING git"
fi

if [[ -n "$MISSING" ]]; then
  echo ""
  print_error "Missing required tools:$MISSING"
  echo ""
  echo "Please install:"
  [[ "$MISSING" == *"Go"* ]] && echo "  • Go: https://go.dev/dl/"
  [[ "$MISSING" == *"Node"* ]] && echo "  • Node.js: https://nodejs.org/"
  [[ "$MISSING" == *"pnpm"* ]] && echo "  • pnpm: npm install -g pnpm"
  [[ "$MISSING" == *"git"* ]] && echo "  • git: https://git-scm.com/"
  exit 1
fi

# Install pnpm if needed
if [[ "$WARNINGS" == *"pnpm"* ]]; then
  print_info "Installing pnpm..."
  if npm install -g pnpm 2>/dev/null; then
    print_success "pnpm installed"
  else
    print_warning "Could not install pnpm globally (may need sudo)"
  fi
fi

echo ""

# ============================================================
# Step 2: Clone Repository (Remote Mode Only)
# ============================================================
if [[ "$MODE" == "remote" ]]; then
  print_step "Cloning GoSaaS repository..."

  # Determine target directory
  if [[ -n "$INSTALL_DIR" ]]; then
    TARGET="$INSTALL_DIR/$NEW_NAME"
  else
    TARGET="./$NEW_NAME"
  fi

  # Check if target exists
  if [[ -d "$TARGET" ]]; then
    die "Directory already exists: $TARGET

Remove it first or choose a different name."
  fi

  # Create parent directory if needed
  if [[ -n "$INSTALL_DIR" ]]; then
    mkdir -p "$INSTALL_DIR" 2>/dev/null || die "Cannot create directory: $INSTALL_DIR"
  fi

  # Clone
  if ! git clone --depth 1 "$REPO_URL" "$TARGET" 2>/dev/null; then
    die "Failed to clone repository.

Check your internet connection and try again."
  fi

  # Enter directory
  cd "$TARGET" || die "Cannot enter directory: $TARGET"

  # Remove git history (fresh start)
  rm -rf .git

  print_success "Repository cloned to $(pwd)"
  echo ""
fi

# ============================================================
# Step 3: Verify We're in a GoSaaS Project
# ============================================================
if [[ ! -f "go.mod" ]]; then
  die "Not in a GoSaaS project directory.

Expected to find: go.mod

Make sure you're running this from the project root."
fi

# Check if already set up with a different name
CURRENT_MODULE=$(grep "^module " go.mod 2>/dev/null | awk '{print $2}')
if [[ "$CURRENT_MODULE" != "$OLD_NAME" ]] && [[ "$CURRENT_MODULE" != "$NEW_NAME" ]]; then
  die "Project already renamed to '$CURRENT_MODULE'.

To use a different name, start fresh:
  rm -rf $(pwd)
  curl -fsSL $REPO_URL/install.sh | bash -s -- $NEW_NAME"
fi

# ============================================================
# Step 4: Rename Project Files
# ============================================================

# Generate proper case name (capitalize first letter)
NEW_NAME_PROPER="$(echo "${NEW_NAME:0:1}" | tr '[:lower:]' '[:upper:]')${NEW_NAME:1}"

if [[ -f "${OLD_NAME}.go" ]] || [[ -f "${OLD_NAME}.api" ]] || [[ -f "etc/${OLD_NAME}.yaml" ]]; then
  print_step "Renaming project to '$NEW_NAME'..."

  # Main files
  [[ -f "${OLD_NAME}.go" ]] && mv "${OLD_NAME}.go" "${NEW_NAME}.go" && print_info "${OLD_NAME}.go → ${NEW_NAME}.go"
  [[ -f "${OLD_NAME}.api" ]] && mv "${OLD_NAME}.api" "${NEW_NAME}.api" && print_info "${OLD_NAME}.api → ${NEW_NAME}.api"
  [[ -f "etc/${OLD_NAME}.yaml" ]] && mv "etc/${OLD_NAME}.yaml" "etc/${NEW_NAME}.yaml" && print_info "etc/${OLD_NAME}.yaml → etc/${NEW_NAME}.yaml"

  print_success "Files renamed"
else
  print_success "Files already renamed"
fi

# ============================================================
# Step 5: Replace ALL occurrences in ALL files
# ============================================================
print_step "Replacing all occurrences of 'gosaas' and 'GoSaaS'..."

# Function to replace in a single file (both lowercase and proper case)
replace_in_file() {
  local file="$1"
  if [[ -f "$file" ]]; then
    # Replace GoSaaS (proper case) with NewName (proper case)
    if grep -q "${OLD_NAME_PROPER}" "$file" 2>/dev/null; then
      sed_inplace "s/${OLD_NAME_PROPER}/${NEW_NAME_PROPER}/g" "$file"
    fi
    # Replace gosaas (lowercase) with newname (lowercase)
    if grep -q "${OLD_NAME}" "$file" 2>/dev/null; then
      sed_inplace "s/${OLD_NAME}/${NEW_NAME}/g" "$file"
    fi
  fi
}

# Go files
print_info "Updating Go files..."
while IFS= read -r -d '' file; do
  replace_in_file "$file"
done < <(find . -type f -name "*.go" ! -path "./vendor/*" -print0 2>/dev/null)

# API definition file
if [[ -f "${NEW_NAME}.api" ]]; then
  replace_in_file "${NEW_NAME}.api"
fi

# YAML files (config, compose, etc)
print_info "Updating YAML files..."
while IFS= read -r -d '' file; do
  replace_in_file "$file"
done < <(find . -type f \( -name "*.yaml" -o -name "*.yml" \) -print0 2>/dev/null)

# Markdown files (docs, README, CLAUDE.md, etc)
print_info "Updating documentation..."
while IFS= read -r -d '' file; do
  replace_in_file "$file"
done < <(find . -type f -name "*.md" -print0 2>/dev/null)

# SQL files
print_info "Updating SQL files..."
while IFS= read -r -d '' file; do
  replace_in_file "$file"
done < <(find . -type f -name "*.sql" -print0 2>/dev/null)

# Config files in root
for file in Makefile .air.toml Dockerfile .env.example go.mod go.sum compose.yaml; do
  replace_in_file "$file"
done

# Frontend files
if [[ -d "app" ]]; then
  print_info "Updating frontend files..."

  # Rename API files first
  [[ -f "app/src/lib/api/${OLD_NAME}.ts" ]] && mv "app/src/lib/api/${OLD_NAME}.ts" "app/src/lib/api/${NEW_NAME}.ts"
  [[ -f "app/src/lib/api/${OLD_NAME}Components.ts" ]] && mv "app/src/lib/api/${OLD_NAME}Components.ts" "app/src/lib/api/${NEW_NAME}Components.ts"

  # TypeScript files
  while IFS= read -r -d '' file; do
    replace_in_file "$file"
  done < <(find app -type f -name "*.ts" -print0 2>/dev/null)

  # Svelte files (all of them - routes, components, layouts)
  while IFS= read -r -d '' file; do
    replace_in_file "$file"
  done < <(find app -type f -name "*.svelte" -print0 2>/dev/null)

  # CSS files
  while IFS= read -r -d '' file; do
    replace_in_file "$file"
  done < <(find app -type f -name "*.css" -print0 2>/dev/null)

  # JSON files (package.json, manifest, etc)
  while IFS= read -r -d '' file; do
    replace_in_file "$file"
  done < <(find app -type f -name "*.json" -print0 2>/dev/null)

  # Static files (webmanifest, etc)
  while IFS= read -r -d '' file; do
    replace_in_file "$file"
  done < <(find app/static -type f -print0 2>/dev/null)

  # .env.example in app
  replace_in_file "app/.env.example"
fi

# Deploy folder
if [[ -d "deploy" ]]; then
  print_info "Updating deploy files..."
  while IFS= read -r -d '' file; do
    replace_in_file "$file"
  done < <(find deploy -type f -print0 2>/dev/null)
fi

# Scripts folder
if [[ -d "scripts" ]]; then
  print_info "Updating scripts..."
  while IFS= read -r -d '' file; do
    replace_in_file "$file"
  done < <(find scripts -type f -print0 2>/dev/null)
fi

print_success "All replacements complete"

# ============================================================
# Step 8: Generate .env
# ============================================================
if [[ "$CREATE_ENV" == true ]]; then
  print_step "Creating configuration..."

  if [[ -f ".env" ]] && [[ "$FORCE_ENV" != true ]]; then
    print_success ".env exists (use --force to overwrite)"
  else
    # Generate secure secrets
    generate_secret() {
      local length=${1:-32}
      local secret=""
      if command -v openssl &> /dev/null; then
        secret=$(openssl rand -hex "$length" 2>/dev/null)
      fi
      if [[ -z "$secret" ]] && [[ -r /dev/urandom ]]; then
        secret=$(head -c "$length" /dev/urandom | xxd -p 2>/dev/null | tr -d '\n')
      fi
      if [[ -z "$secret" ]] && [[ -r /dev/urandom ]]; then
        secret=$(head -c "$length" /dev/urandom | od -A n -t x1 2>/dev/null | tr -d ' \n')
      fi
      echo "$secret"
    }

    ACCESS_SECRET=$(generate_secret 32)
    ADMIN_PASSWORD=$(generate_secret 16)

    if [[ -z "$ACCESS_SECRET" ]] || [[ ${#ACCESS_SECRET} -lt 32 ]]; then
      ACCESS_SECRET="CHANGE_THIS_$(date +%s)_INSECURE_PLEASE_REGENERATE"
      print_warning "Could not generate secure JWT secret - please regenerate with: openssl rand -hex 32"
    fi

    if [[ -z "$ADMIN_PASSWORD" ]] || [[ ${#ADMIN_PASSWORD} -lt 16 ]]; then
      ADMIN_PASSWORD="admin_$(date +%s)"
      print_warning "Could not generate secure admin password - please change in .env"
    fi

    cat > .env << ENVEOF
# ============================================================
# ${NEW_NAME_PROPER} Configuration
# Generated: $(date)
# ============================================================

# Core Settings
APP_BASE_URL=http://localhost:5173
APP_DOMAIN=localhost
PRODUCTION_MODE=false

# JWT Secret (auto-generated - keep this safe!)
ACCESS_SECRET=${ACCESS_SECRET}

# Backoffice Admin (for SaaS metrics at /backoffice)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=${ADMIN_PASSWORD}

# CORS
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:8888

# ============================================================
# Standalone Mode (Default) - SQLite + Stripe
# ============================================================
LEVEE_ENABLED=false
SQLITE_PATH=./data/${NEW_NAME}.db

# Stripe - Get keys from https://dashboard.stripe.com/apikeys
# Leave as-is for local development without payments
STRIPE_SECRET_KEY=sk_test_your_key_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# ============================================================
# Email (Optional)
# ============================================================
# SMTP_HOST=smtp.gmail.com
# SMTP_PORT=587
# SMTP_USER=your-email@gmail.com
# SMTP_PASS=your-app-password
# EMAIL_FROM_ADDRESS=noreply@${NEW_NAME}.com
# EMAIL_FROM_NAME=${NEW_NAME_PROPER}

# ============================================================
# Levee Mode (Optional) - Full platform features
# ============================================================
# LEVEE_ENABLED=true
# LEVEE_API_KEY=lvk_your_api_key
# LEVEE_BASE_URL=https://api.levee.sh
ENVEOF

    print_success ".env created with secure secret"
  fi
fi

# ============================================================
# Step 9: Install Dependencies
# ============================================================
if [[ "$INSTALL_DEPS" == true ]]; then
  print_step "Installing dependencies..."

  # Go
  if command -v go &> /dev/null; then
    if go mod download 2>/dev/null && go mod tidy 2>/dev/null; then
      print_success "Go dependencies"
    else
      print_warning "Go deps may need: go mod tidy"
    fi
  fi

  # Frontend
  if [[ -d "app" ]] && command -v pnpm &> /dev/null; then
    if (cd app && pnpm install --silent 2>/dev/null); then
      print_success "Frontend dependencies"
    else
      print_warning "Frontend deps may need: cd app && pnpm install"
    fi
  fi
fi

# ============================================================
# Step 10: Create Data Directory
# ============================================================
mkdir -p data 2>/dev/null || true

# ============================================================
# Step 11: Initialize Git (Remote Mode)
# ============================================================
if [[ "$MODE" == "remote" ]]; then
  print_step "Initializing git repository..."
  git init -q 2>/dev/null
  git add . 2>/dev/null
  git commit -q -m "Initial commit - ${NEW_NAME} from GoSaaS" 2>/dev/null
  print_success "Git repository initialized"
fi

# ============================================================
# Done!
# ============================================================
echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║${NC}                  ✅ Setup Complete!                        ${GREEN}║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  ${BOLD}Your app:${NC}  $NEW_NAME"
echo -e "  ${BOLD}Location:${NC}  $(pwd)"
echo ""
echo -e "  ${BOLD}Start development:${NC}"
echo ""
echo "    make dev              # Starts backend + frontend"
echo ""
echo -e "  ${BOLD}Or manually:${NC}"
echo ""
echo "    make air              # Backend with hot reload"
echo "    cd app && pnpm dev    # Frontend (another terminal)"
echo ""
echo -e "  ${BOLD}Then visit:${NC} http://localhost:5173"
echo ""
echo -e "  ${BOLD}Admin Backoffice:${NC}"
echo "    URL:      http://localhost:8888/backoffice"
echo "    Username: admin"
echo "    Password: (see ADMIN_PASSWORD in .env)"
echo ""
echo -e "  ${BOLD}Configure payments:${NC}"
echo "    1. Get keys from https://dashboard.stripe.com/apikeys"
echo "    2. Edit .env with your STRIPE_* values"
echo ""

# ============================================================
# Clean Up - Remove installer to prevent re-running
# ============================================================
# By this point we're always in the project directory
# Remove install.sh to prevent accidental re-installation
rm -f ./install.sh 2>/dev/null || true

# ============================================================
# Offer to Start
# ============================================================
if [[ "$OFFER_START" == true ]] && [[ -t 0 ]]; then
  echo -ne "  Start the app now? [Y/n] "
  read -r START_NOW
  if [[ ! "$START_NOW" =~ ^[Nn]$ ]]; then
    echo ""
    print_step "Starting ${NEW_NAME}..."
    echo ""
    exec make dev
  fi
fi
