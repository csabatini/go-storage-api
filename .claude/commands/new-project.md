---
description: Create a new project with all scaffolding rules applied
argument-hint: <path> [profile-or-options...]
allowed-tools: Bash, Write, Read, AskUserQuestion
---

# New Project Scaffold

Create a new project with all best practices from the Claude Code Mastery Guides.

**Arguments:** $ARGUMENTS

## Argument Parsing

### Step 0 — Read the config file

Before parsing arguments, read `claude-mastery-project.conf` (in the starter kit root or `~/.claude/claude-mastery-project.conf` as fallback).

Extract the `[global]` section for `root_dir`. This is the default parent directory for new projects.

### Step 0.0 — Global Claude Config (one-time setup)

Check if the user already has the global Claude config installed:

```bash
# Check if global CLAUDE.md exists
ls ~/.claude/CLAUDE.md 2>/dev/null
```

**If `~/.claude/CLAUDE.md` does NOT exist:**
- ASK: "You don't have a global CLAUDE.md yet. Want me to install the Claude Code Mastery global config to `~/.claude/`? This sets up security rules, hooks, and standards that apply to ALL your projects. (This is a one-time setup.)"
- If yes: copy `global-claude-md/CLAUDE.md` → `~/.claude/CLAUDE.md` and `global-claude-md/settings.json` → `~/.claude/settings.json`
- Also copy hooks: `mkdir -p ~/.claude/hooks && cp .claude/hooks/block-secrets.py ~/.claude/hooks/ && cp .claude/hooks/verify-no-secrets.sh ~/.claude/hooks/`

**If `~/.claude/CLAUDE.md` DOES exist:**
- ASK: "You already have a global CLAUDE.md. Want me to check if the starter kit version has anything new to merge in?"
- If yes: diff the two files and show what's different. Let the user decide what to merge.
- If no: skip and continue.

**This step typically only happens once.** After the first install, the global config persists across all projects.

### Step 0.1 — Resolve the project path

The **first argument** is the project name or path. Resolve it using `root_dir`:

1. **Explicit path** (starts with `./`, `../`, `~/`, or `/`) → use as-is
   - `/new-project ~/code/my-app` → creates at `~/code/my-app`
   - `/new-project ./my-app` → creates at `./my-app`

2. **Just a name** (no path separators) → prepend `root_dir` from `[global]`
   - Config has `root_dir = ~/projects`
   - `/new-project my-app` → creates at `~/projects/my-app`
   - `/new-project tims-api` → creates at `~/projects/tims-api`

3. **No argument at all** → ASK the user for the project name, then prepend `root_dir`

Everything after the project name/path is shorthand options or a profile name.

### Shorthand Arguments (after the path/name)

Parse remaining $ARGUMENTS for these keywords:

**Profiles:** `clean`, `default`, `api`, `static-site`, `quick`, `enterprise` (from `claude-mastery-project.conf`)
**Special:** `clean` — Claude infrastructure only, zero coding opinions (see Clean Mode below)
**Project types:** `webapp`, `api`, `fullstack`, `cli`
**Frameworks:** `vite`, `react`, `next`, `nextjs`, `astro`, `fastify`, `express`, `hono`
**Options:** `seo`, `ssr`, `tailwind`, `prisma`, `docker`, `ci`, `multiregion`
**Hosting:** `dokploy`, `vercel`, `static`
**Database:** `mongo`, `postgres`, `sqlite`
**Analytics:** `rybbit`
**MCP servers:** `playwright`, `context7`, `rulecatch`
**NPM extras:** `ai-pooler` (installs @rulecatch/ai-pooler)
**Package managers:** `pnpm`, `npm`, `bun`

Examples:
- `/new-project my-app` — creates at ~/projects/my-app (from root_dir), asks questions
- `/new-project my-app clean` — Claude infrastructure only, no coding opinions
- `/new-project my-app default` — creates at ~/projects/my-app with default profile
- `/new-project my-app fullstack next seo tailwind pnpm` — ~/projects/my-app, skips all questions
- `/new-project ./custom-path/my-app api fastify` — explicit path, ignores root_dir
- `/new-project ~/code/my-app default` — explicit path, uses default profile
- `/new-project my-app fullstack next mongo playwright context7 rulecatch` — full stack

Any keyword not provided = ask the user.

---

## Clean Mode — `clean`

**If `clean` is detected in arguments, skip ALL of Steps 1-2 below and follow this section instead.**

Clean mode gives the user every piece of Claude Code infrastructure without imposing ANY opinions about how they should code, what language to use, what framework to pick, or how to structure their source code.

### What `clean` creates

```
project/
├── CLAUDE.md              # Security rules ONLY (see below)
├── CLAUDE.local.md        # Personal overrides template
├── .claude/
│   ├── settings.json      # Hooks configuration
│   ├── commands/
│   │   ├── review.md
│   │   ├── commit.md
│   │   ├── progress.md
│   │   ├── test-plan.md
│   │   ├── architecture.md
│   │   ├── new-project.md
│   │   ├── security-check.md
│   │   ├── optimize-docker.md
│   │   ├── create-e2e.md
│   │   └── worktree.md
│   ├── skills/
│   │   ├── code-review/SKILL.md
│   │   └── create-service/SKILL.md
│   ├── agents/
│   │   ├── code-reviewer.md
│   │   └── test-writer.md
│   └── hooks/
│       ├── block-secrets.py
│       ├── lint-on-save.sh
│       └── verify-no-secrets.sh
├── project-docs/
│   ├── ARCHITECTURE.md
│   ├── INFRASTRUCTURE.md
│   └── DECISIONS.md
├── tests/
│   ├── CHECKLIST.md
│   └── ISSUES_FOUND.md
├── .env                   # Empty (NEVER commit)
├── .env.example           # Template with NODE_ENV and PORT
├── .gitignore             # Standard ignores
├── .dockerignore          # Standard ignores
└── README.md              # Minimal project readme
```

### What `clean` does NOT create

- No `src/` directory — user decides their own structure
- No `package.json` — user picks their own language, runtime, and package manager
- No `tsconfig.json` — user may not even use TypeScript
- No `vitest.config.ts` or `playwright.config.ts` — user picks their own test tools
- No database wrapper or `scripts/db-query.ts` — user picks their own database
- No content builder — user decides if they need one
- No SEO templates — user decides their own approach
- No port assignments — user decides their own ports
- No framework-specific configs — user picks their own framework

### Clean CLAUDE.md content

The CLAUDE.md for `clean` mode contains ONLY universal, non-opinionated rules:

```markdown
# CLAUDE.md — Project Instructions

---

## Critical Rules

### 0. NEVER Publish Sensitive Data

- NEVER commit passwords, API keys, tokens, or secrets to git/npm/docker
- NEVER commit `.env` files — ALWAYS verify `.env` is in `.gitignore`
- Before ANY commit: verify no secrets are included
- NEVER output secrets in suggestions, logs, or responses

### 1. NEVER Hardcode Credentials

- ALWAYS use environment variables for secrets
- NEVER put API keys, passwords, or tokens directly in code
- NEVER hardcode connection strings — use environment variables from .env

### 2. ALWAYS Ask Before Deploying

- NEVER auto-deploy, even if the fix seems simple
- NEVER assume approval — wait for explicit "yes, deploy"
- ALWAYS ask before deploying to production

---

## When Something Seems Wrong

Before jumping to conclusions:

- Missing UI element? → Check feature gates BEFORE assuming bug
- Empty data? → Check if services are running BEFORE assuming broken
- 404 error? → Check service separation BEFORE adding endpoint
- Auth failing? → Check which auth system BEFORE debugging
- Test failing? → Read the error message fully BEFORE changing code

---

## Project Documentation

| Document | Purpose | When to Read |
|----------|---------|--------------|
| `project-docs/ARCHITECTURE.md` | System overview & data flow | Before architectural changes |
| `project-docs/INFRASTRUCTURE.md` | Deployment details | Before environment changes |
| `project-docs/DECISIONS.md` | Architectural decisions | Before proposing alternatives |

**ALWAYS read relevant docs before making cross-service changes.**

---

## Workflow Preferences

- Quality over speed — if unsure, ask before executing
- One task, one chat — `/clear` between unrelated tasks
- When testing: queue observations, fix in batch (not one at a time)

---

## Naming — NEVER Rename Mid-Project

If you must rename packages, modules, or key variables:

1. Create a checklist of ALL files and references first
2. Use IDE semantic rename (not search-and-replace)
3. Full project search for old name after renaming
4. Check: .md files, .txt files, .env files, comments, strings, paths
5. Start a FRESH Claude session after renaming
```

### Clean mode steps

1. Resolve project path (same as Step 0 / 0.1 above)
2. Create the project directory
3. Copy ALL `.claude/` contents from the starter kit (commands, skills, agents, hooks, settings.json)
4. Create project-docs/ with ARCHITECTURE.md, INFRASTRUCTURE.md, DECISIONS.md templates
5. Create tests/ with CHECKLIST.md and ISSUES_FOUND.md
6. Create the clean CLAUDE.md (security rules only, as shown above)
7. Create CLAUDE.local.md template
8. Create .env (empty), .env.example (NODE_ENV + PORT only)
9. Create .gitignore and .dockerignore with standard entries
10. Create a minimal README.md
11. Initialize git and create initial commit
12. Report what was created

### Clean verification checklist

- [ ] `.claude/` directory with all commands, skills, agents, hooks
- [ ] `.claude/settings.json` with hooks wired up
- [ ] `CLAUDE.md` has ONLY security rules (no TypeScript, no ports, no quality gates)
- [ ] `project-docs/` has all three templates
- [ ] `tests/` has CHECKLIST.md and ISSUES_FOUND.md
- [ ] `.env` exists (empty)
- [ ] `.env.example` exists
- [ ] `.gitignore` includes .env, node_modules/, dist/, CLAUDE.local.md
- [ ] `.dockerignore` exists
- [ ] NO `package.json`, `tsconfig.json`, or framework configs created
- [ ] NO `src/` directory created
- [ ] Git initialized with initial commit

**After creating a `clean` project, the user can add their own language, framework, and structure — Claude will follow the security rules and use the slash commands without imposing any coding patterns.**

---

## Step 1 — Ask the User (skip questions answered by arguments)

For any choices NOT provided via arguments, ask the user (use AskUserQuestion):

### Question 1: Project Type
"What type of project are you building?"
- **Web App** — Frontend with UI (SPA or SSR)
- **API** — Backend REST/GraphQL service
- **Full-Stack** — Frontend + backend in one repo
- **CLI Tool** — Command-line application

### Question 2: Framework (based on project type)

**If Web App or Full-Stack:**
"Which framework do you want to use?"
- **Vite + React** — Fastest HMR, lightweight, great for SPAs (Recommended)
- **Next.js (App Router)** — SSR, server components, built-in routing
- **Astro** — Content-first, island architecture, great for marketing/docs sites

**If API:**
"Which framework do you want to use?"
- **Fastify** — Fastest Node.js HTTP framework, built-in validation (Recommended)
- **Express** — Most popular, largest ecosystem
- **Hono** — Ultra-lightweight, edge-ready

**If CLI Tool:**
- Use **Commander.js** + **TypeScript** (no framework question needed)

### Question 3: SSR Requirement (Web App / Full-Stack only, skip if Next.js or Astro already chosen)
"Do you need server-side rendering (SSR)?"
- **No (SPA)** — Client-side only, simpler deployment (Recommended for dashboards/apps)
- **Yes (SSR)** — SEO-critical pages, faster first paint (Recommended for public-facing sites)

If they chose Vite + React and want SSR, switch to **Next.js (App Router)** or add **vite-plugin-ssr**.

### Question 4: Package Manager
"Which package manager?"
- **pnpm** — Fast, disk-efficient (Recommended)
- **npm** — Default, universal
- **bun** — Fastest, newer ecosystem

### Question 5: Hosting / Deployment
"Where will this be deployed?" (skip if `dokploy`, `vercel`, or `static` in arguments)
- **Dokploy on Hostinger VPS** — Self-hosted Docker containers with Dokploy management (Recommended for full control)
- **Vercel** — Zero-config for Next.js / static sites
- **Static hosting** — GitHub Pages, Netlify, Cloudflare Pages
- **None / Decide later** — Skip deployment scaffolding

### Question 6: Extras (multi-select)
"What extras do you want to include?"
- **Tailwind CSS** — Utility-first CSS framework
- **Prisma** — Type-safe database ORM
- **Docker** — Containerized deployment (auto-included with Dokploy)
- **GitHub Actions CI** — Automated testing pipeline
- **Multi-region** — US + EU deployment (Dokploy only)

## Step 2 — Create the Project

Based on answers, scaffold the project:

1. Create project directory
2. Initialize with chosen framework and package manager
3. Install TypeScript + Vitest (ALWAYS, non-negotiable)
4. Create ALL required files (see below)
5. Apply framework-specific rules
6. Apply SEO requirements (if web project)
7. Initialize git repository
8. Create initial commit: "Initial project scaffold"
9. Display verification checklist

## Required Files (EVERY Project)

- `.env` — Empty, for secrets (NEVER commit)
- `.env.example` — Template with placeholder values
- `.gitignore` — Must include: .env, .env.*, node_modules/, dist/, CLAUDE.local.md
- `.dockerignore` — Must include: .env, .git/, node_modules/
- `README.md` — Project overview (reference env vars, don't hardcode)
- `CLAUDE.md` — Must include: project overview, tech stack, build/test/dev commands, architecture, port assignments
- `tsconfig.json` — Strict mode enabled, `noUncheckedIndexedAccess: true`

## Required Directory Structure

```
project/
├── src/
├── tests/
├── project-docs/
│   ├── ARCHITECTURE.md
│   ├── INFRASTRUCTURE.md
│   └── DECISIONS.md
├── .claude/
│   ├── commands/
│   ├── skills/
│   └── agents/
└── scripts/
    ├── db-query.ts          # (MongoDB only) Test Query Master
    └── queries/             # (MongoDB only) Individual dev/test query files
```

## MongoDB Test Query System (projects with `mongo` database)

When the project uses MongoDB, ALWAYS scaffold the db-query system:

1. Create `scripts/db-query.ts` — the master index/CLI runner
2. Create `scripts/queries/` directory for individual query files
3. Add the db-query rules to the project's `CLAUDE.md`

**The rule that MUST be in every MongoDB project's CLAUDE.md:**

> ALL ad-hoc / test / dev database queries go through `scripts/db-query.ts`.
> When asked to look something up in the database:
> 1. Create a query file in `scripts/queries/<name>.ts`
> 2. Register it in `scripts/db-query.ts`
> 3. NEVER create standalone scripts or inline queries in `src/`

This prevents Claude from scattering random query scripts all over the project.

## TypeScript + Vitest + Playwright (ALWAYS)

Every project MUST have Vitest for unit tests and Playwright for E2E tests.

### vitest.config.ts
```typescript
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node', // or 'jsdom' for web
    include: ['tests/unit/**/*.test.ts', 'tests/integration/**/*.test.ts'],
    exclude: ['tests/e2e/**/*'],
  },
});
```

### playwright.config.ts
```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: [['html'], ['list']],
  use: {
    baseURL: 'http://localhost:4000', // TEST port, not dev port
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
  ],
  webServer: [
    {
      command: 'pnpm dev:test:website',
      port: 4000,
      reuseExistingServer: !process.env.CI,
      timeout: 30_000,
    },
  ],
});
```

### package.json test scripts (REQUIRED in every project)
```json
{
  "scripts": {
    "dev:test:website": "PORT=4000 tsx watch src/index.ts",
    "dev:test:api": "PORT=4010 tsx watch src/index.ts",
    "test": "pnpm test:unit && pnpm test:e2e",
    "test:unit": "vitest run",
    "test:unit:watch": "vitest",
    "test:coverage": "vitest run --coverage",
    "test:e2e": "pnpm test:kill-ports && playwright test",
    "test:e2e:ui": "pnpm test:kill-ports && playwright test --ui",
    "test:e2e:headed": "pnpm test:kill-ports && playwright test --headed",
    "test:e2e:report": "playwright show-report",
    "test:kill-ports": "lsof -ti:4000,4010,4020 | xargs kill -9 2>/dev/null || true"
  }
}
```

**CRITICAL: `test:kill-ports` runs BEFORE every E2E test command.** This prevents "port already in use" failures. Never skip this step.

### tsconfig.json (minimum)
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitOverride": true,
    "esModuleInterop": true,
    "skipLibCheck": true
  }
}
```

## Node.js Entry Point Requirements

Add to EVERY Node.js entry point. If the project uses MongoDB, use `gracefulShutdown` to close pools before exit:

```typescript
// WITH MongoDB (projects using src/core/db/)
import { gracefulShutdown } from '@/core/db/index.js';

process.on('SIGTERM', () => gracefulShutdown(0));
process.on('SIGINT', () => gracefulShutdown(0));
process.on('uncaughtException', (err) => {
  console.error('Uncaught Exception:', err);
  gracefulShutdown(1);
});
process.on('unhandledRejection', (reason) => {
  console.error('Unhandled Rejection:', reason);
  gracefulShutdown(1);
});
```

```typescript
// WITHOUT MongoDB (no database or non-Mongo projects)
process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Rejection:', reason);
  process.exit(1);
});

process.on('uncaughtException', (error) => {
  console.error('Uncaught Exception:', error);
  process.exit(1);
});
```

## Mandatory SEO (ALL Web Projects)

Every web project MUST include these SEO fundamentals. This is non-negotiable for any page that serves HTML.

### 1. HTML Meta Tags (in layout/head)

```html
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Page Title — Site Name</title>
  <meta name="description" content="Concise page description (150-160 chars)">
  <meta name="robots" content="index, follow">
  <link rel="canonical" href="https://example.com/current-page">

  <!-- Open Graph (Facebook, LinkedIn, Discord) -->
  <meta property="og:type" content="website">
  <meta property="og:title" content="Page Title">
  <meta property="og:description" content="Page description">
  <meta property="og:image" content="https://example.com/og-image.png">
  <meta property="og:url" content="https://example.com/current-page">
  <meta property="og:site_name" content="Site Name">

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary_large_image">
  <meta name="twitter:title" content="Page Title">
  <meta name="twitter:description" content="Page description">
  <meta name="twitter:image" content="https://example.com/og-image.png">
</head>
```

### 2. JSON-LD Structured Data (schema.org)

EVERY web project must include at minimum an Organization or WebSite schema:

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "Your Site Name",
  "url": "https://example.com",
  "description": "Site description",
  "publisher": {
    "@type": "Organization",
    "name": "Your Organization",
    "logo": {
      "@type": "ImageObject",
      "url": "https://example.com/logo.png"
    }
  }
}
</script>
```

For specific page types, add the appropriate schema:
- **Article pages:** `@type: "Article"` with author, datePublished, dateModified
- **Product pages:** `@type: "Product"` with price, availability, reviews
- **FAQ pages:** `@type: "FAQPage"` with question/answer pairs
- **How-to pages:** `@type: "HowTo"` with steps
- **Breadcrumbs:** `@type: "BreadcrumbList"` on all pages with navigation depth

### 3. Technical SEO Files

Create these in the project root (or public directory):

**robots.txt:**
```
User-agent: *
Allow: /
Sitemap: https://example.com/sitemap.xml
```

**sitemap.xml** (or generate dynamically):
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2025-01-01</lastmod>
    <priority>1.0</priority>
  </url>
</urlset>
```

### 4. Performance SEO

- Images MUST use WebP format with `alt` attributes
- Include `<link rel="preconnect">` for external domains (fonts, analytics, CDNs)
- Set proper cache headers for static assets
- Ensure Largest Contentful Paint (LCP) < 2.5 seconds

### 5. Framework-Specific SEO

**Next.js:**
- Use `metadata` export in layout.tsx / page.tsx (App Router)
- Use `generateMetadata()` for dynamic pages
- JSON-LD via `<script>` in layout or use `next-seo` package
- next/image for automatic WebP conversion and lazy loading
- Automatic sitemap generation with `next-sitemap`

**Vite + React (SPA):**
- Use `react-helmet-async` for dynamic `<head>` management
- For SEO-critical SPAs, consider prerendering with `vite-plugin-ssr` or `prerender-spa-plugin`
- NOTE: SPAs have inherent SEO limitations — if SEO is critical, recommend SSR

**Astro:**
- Built-in `<head>` management in `.astro` layouts
- Automatic sitemap with `@astrojs/sitemap`
- Built-in image optimization

## Framework-Specific Rules

### Vite + React
- Use Vite's built-in HMR (no config needed)
- Add `@vitejs/plugin-react` or `@vitejs/plugin-react-swc` (SWC is faster)
- Use path aliases: `"@/*": ["./src/*"]` in tsconfig
- Vitest shares Vite config — zero extra setup

### Next.js (App Router)
- Use App Router (NOT Pages Router)
- Create `src/app/` directory structure
- Use Server Components by default, `"use client"` only when needed
- Strict mode in next.config
- Use `metadata` export for SEO (not `<Head>`)

### Fastify
- Use `@fastify/type-provider-typebox` for schema validation
- Register routes as plugins for encapsulation
- Use `fastify-swagger` for auto-generated API docs
- All routes under `/api/v1/` prefix

### Astro
- Use content collections for structured content
- Islands architecture: interactive components only where needed
- Built-in image optimization with `<Image>` component

### Python
- Create `pyproject.toml` (not setup.py)
- Use `src/` layout
- Include `requirements.txt` AND `requirements-dev.txt`

### Docker
- Multi-stage builds ALWAYS
- Never run as root (create service-specific user)
- Include health checks
- COPY package.json first for layer caching
- For monorepos: build shared packages first, copy dist into deployed node_modules

### Docker Multi-Stage Template
```dockerfile
# Stage 1: Builder
FROM node:20-alpine AS builder
WORKDIR /app

# Install package manager
RUN corepack enable && corepack prepare pnpm@latest --activate

# Install dependencies (cached layer)
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Build args for Next.js (baked at build time)
ARG NEXT_PUBLIC_RYBBIT_SITE_ID
ARG NEXT_PUBLIC_RYBBIT_URL
ENV NEXT_PUBLIC_RYBBIT_SITE_ID=$NEXT_PUBLIC_RYBBIT_SITE_ID
ENV NEXT_PUBLIC_RYBBIT_URL=$NEXT_PUBLIC_RYBBIT_URL

# Copy source and build
COPY . .
RUN pnpm build

# Stage 2: Runner
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Non-root user
RUN addgroup --system --gid 1001 appgroup && \
    adduser --system --uid 1001 appuser
USER appuser

# Copy built artifacts
COPY --from=builder --chown=appuser:appgroup /app/dist ./dist
COPY --from=builder --chown=appuser:appgroup /app/node_modules ./node_modules
COPY --from=builder --chown=appuser:appgroup /app/package.json ./

EXPOSE 3000
CMD ["node", "dist/server.js"]
```

## Dokploy on Hostinger VPS (if selected)

When Dokploy is selected as the hosting target, scaffold a complete deployment pipeline:

### Deployment Architecture
```
Code → Docker Build → Local Test → Docker Hub → Dokploy (webhook) → Live
```

### Required Environment Variables (.env.example additions)

```bash
# Dokploy Deployment
DOKPLOY_URL=http://your-vps-ip:3000/api
DOKPLOY_API_KEY=your_dokploy_api_key
DOKPLOY_APP_ID=your_application_id
DOKPLOY_REFRESH_TOKEN=your_webhook_refresh_token

# Docker Hub
DOCKER_HUB_USER=your_docker_username
DOCKER_IMAGE_NAME=your_docker_username/your_app_name

# Region (if multi-region)
DEPLOY_REGION=us
```

### Deployment Script: scripts/deploy.sh

Create this deployment script:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Load environment
source .env

IMAGE="$DOCKER_IMAGE_NAME:latest"
TAG="${1:-latest}"

echo "=== Building Docker image ==="
docker build -t "$IMAGE" .

echo "=== Testing locally ==="
docker run -d -p 3000:3000 --name deploy-test "$IMAGE"
sleep 5

if ! curl -sf http://localhost:3000 > /dev/null; then
  echo "ERROR: Local test FAILED. Aborting deployment."
  docker logs deploy-test
  docker stop deploy-test && docker rm deploy-test
  exit 1
fi

echo "Local test PASSED."
docker stop deploy-test && docker rm deploy-test

echo "=== Pushing to Docker Hub ==="
docker push "$IMAGE"

echo "=== Deploying via Dokploy ==="
RESPONSE=$(curl -s -X POST \
  -H "x-api-key: $DOKPLOY_API_KEY" \
  -H "Content-Type: application/json" \
  "$DOKPLOY_URL/application.deploy" \
  -d "{\"applicationId\":\"$DOKPLOY_APP_ID\"}")

echo "Dokploy response: $RESPONSE"
echo "=== Deployment complete ==="
```

### Dokploy API Reference (for CLAUDE.md)

Add these to the project's CLAUDE.md when Dokploy is selected:

```markdown
## Deployment Commands

### Deploy (build, test, push, deploy)
bash scripts/deploy.sh

### Dokploy API (direct)
# List all projects
curl -s -H "x-api-key: $DOKPLOY_API_KEY" "$DOKPLOY_URL/project.all"

# Deploy application
curl -s -X POST -H "x-api-key: $DOKPLOY_API_KEY" -H "Content-Type: application/json" \
  "$DOKPLOY_URL/application.deploy" -d '{"applicationId":"APP_ID"}'

# Redeploy (rebuild)
curl -s -X POST -H "x-api-key: $DOKPLOY_API_KEY" -H "Content-Type: application/json" \
  "$DOKPLOY_URL/application.redeploy" -d '{"applicationId":"APP_ID"}'

# Start / Stop
curl -s -X POST -H "x-api-key: $DOKPLOY_API_KEY" -H "Content-Type: application/json" \
  "$DOKPLOY_URL/application.start" -d '{"applicationId":"APP_ID"}'

# Webhook deploy (no auth needed — use refresh token)
curl -X POST http://your-vps-ip:3000/api/deploy/REFRESH_TOKEN
```

### Multi-Region Support (if selected)

When `multiregion` is selected, scaffold for US + EU:

```bash
# .env.example additions for multi-region
DOKPLOY_URL_US=http://us-vps-ip:3000/api
DOKPLOY_API_KEY_US=your_us_api_key
DOKPLOY_APP_ID_US=your_us_app_id

DOKPLOY_URL_EU=http://eu-vps-ip:3000/api
DOKPLOY_API_KEY_EU=your_eu_api_key
DOKPLOY_APP_ID_EU=your_eu_app_id
```

**CRITICAL multi-region rules (add to CLAUDE.md):**
- US containers NEVER connect to EU databases, and vice versa
- Each container gets region-specific `MONGODB_URI` or `DATABASE_URL`
- `DEPLOY_REGION` env var must match the VPS region
- When pushing images: push `:latest` for US, push `:eu` tag for EU
- ALWAYS deploy to both regions — never leave them out of sync

### scripts/deploy-all.sh (multi-region)

```bash
#!/usr/bin/env bash
set -euo pipefail
source .env

IMAGE="$DOCKER_IMAGE_NAME"

# Build and test locally first
docker build -t "$IMAGE:latest" .
docker run -d -p 3000:3000 --name deploy-test "$IMAGE:latest"
sleep 5
curl -sf http://localhost:3000 > /dev/null || { echo "FAILED"; docker logs deploy-test; docker stop deploy-test; docker rm deploy-test; exit 1; }
docker stop deploy-test && docker rm deploy-test

# Push both tags
docker push "$IMAGE:latest"
docker tag "$IMAGE:latest" "$IMAGE:eu"
docker push "$IMAGE:eu"

# Deploy to both regions
echo "Deploying to US..."
curl -s -X POST -H "x-api-key: $DOKPLOY_API_KEY_US" -H "Content-Type: application/json" \
  "$DOKPLOY_URL_US/application.deploy" -d "{\"applicationId\":\"$DOKPLOY_APP_ID_US\"}"

echo "Deploying to EU..."
curl -s -X POST -H "x-api-key: $DOKPLOY_API_KEY_EU" -H "Content-Type: application/json" \
  "$DOKPLOY_URL_EU/application.deploy" -d "{\"applicationId\":\"$DOKPLOY_APP_ID_EU\"}"

echo "=== Both regions deployed ==="
```

## Analytics: Rybbit (if selected)

When `rybbit` is selected as the analytics provider, scaffold tracking into the project:

### Required Environment Variables

```bash
# .env.example additions
NEXT_PUBLIC_RYBBIT_SITE_ID=your_rybbit_site_id
NEXT_PUBLIC_RYBBIT_URL=https://app.rybbit.io
```

### Next.js Integration (layout.tsx)

```tsx
<head>
  {process.env.NEXT_PUBLIC_RYBBIT_SITE_ID && (
    <script
      src={`${process.env.NEXT_PUBLIC_RYBBIT_URL || 'https://app.rybbit.io'}/api/script.js`}
      data-site-id={process.env.NEXT_PUBLIC_RYBBIT_SITE_ID}
      defer
    />
  )}
</head>
```

### Vite / Astro / Static HTML Integration

```html
<script
  src="https://app.rybbit.io/api/script.js"
  data-site-id="YOUR_SITE_ID"
  defer
></script>
```

### Docker Build Args (for Next.js on Dokploy)

When using both Rybbit + Dokploy + Next.js, the Rybbit env vars must be passed as build args:

```dockerfile
ARG NEXT_PUBLIC_RYBBIT_SITE_ID
ARG NEXT_PUBLIC_RYBBIT_URL
ENV NEXT_PUBLIC_RYBBIT_SITE_ID=$NEXT_PUBLIC_RYBBIT_SITE_ID
ENV NEXT_PUBLIC_RYBBIT_URL=$NEXT_PUBLIC_RYBBIT_URL
```

### Important
- Each website MUST have its own unique Rybbit site ID
- Create a new site in the Rybbit dashboard at https://app.rybbit.io
- NEVER reuse site IDs across different projects
- After deployment, verify the script is present in the page source

## AI-Pooler Setup (if @rulecatch/ai-pooler in npm list)

When the default profile or user selects ai-pooler:

```bash
# Install and initialize
npx @rulecatch/ai-pooler init --api-key=dc_your_key --region=us
```

Add to `.env.example`:
```bash
RULECATCH_API_KEY=dc_your_api_key_here
RULECATCH_REGION=us
```

## MCP Server Setup (if selected)

When MCP servers are selected, add them to the project setup:

```bash
# Context7 — Live documentation (eliminates outdated API answers)
claude mcp add context7 -- npx -y @upstash/context7-mcp@latest

# Playwright — E2E testing
claude mcp add playwright -- npx -y @anthropic-ai/playwright-mcp

# RuleCatch — AI development analytics & rule monitoring
npx @rulecatch/mcp-server init
```

Add selected MCP servers to the project's CLAUDE.md under a "## MCP Servers" section.

## Profile System: claude-mastery-project.conf

If the user passes `default` (or any profile name), read `claude-mastery-project.conf` from the project root. This file defines reusable presets so users don't re-type preferences.

### claude-mastery-project.conf Format

```ini
# Claude Mastery Project Configuration
# Define profiles with preset options for /new-project

[default]
type = fullstack
framework = next
hosting = dokploy
package_manager = pnpm
database = mongo
options = seo, tailwind, docker, ci
mcp = playwright, context7, rulecatch

[api]
type = api
framework = fastify
hosting = dokploy
package_manager = pnpm
database = mongo
options = docker, ci
mcp = context7, rulecatch

[static-site]
type = webapp
framework = astro
hosting = static
package_manager = pnpm
options = seo, tailwind
mcp = context7

[quick]
type = webapp
framework = vite
hosting = vercel
package_manager = pnpm
options = tailwind
mcp = context7
```

### How Profiles Work

1. Read `claude-mastery-project.conf` from project root (or `~/.claude/claude-mastery-project.conf` for global defaults)
2. Parse the named profile section
3. Apply all settings from the profile
4. Any additional arguments OVERRIDE profile settings
5. Missing settings from profile = ask the user

Examples:
- `/new-project my-app default` — uses [default] profile for everything
- `/new-project my-app api` — uses [api] profile
- `/new-project my-app default vercel` — uses [default] but overrides hosting to Vercel
- `/new-project my-app` — no profile, asks all questions

### Create Default Config

When scaffolding the starter kit itself, create `claude-mastery-project.conf` with the profiles above as starting templates. Users customize to their preferences.

## Verification Checklist

After creation, verify and report:

**Core files:**
- [ ] .env exists (empty)
- [ ] .env.example exists (with placeholders)
- [ ] .gitignore includes all required entries
- [ ] .dockerignore exists
- [ ] CLAUDE.md has all required sections (overview, stack, commands, ports)
- [ ] package.json has ALL required scripts (dev, build, test, test:e2e, test:kill-ports)
- [ ] Error handlers in entry point (gracefulShutdown for MongoDB projects)
- [ ] TypeScript strict mode enabled

**Testing:**
- [ ] vitest.config.ts created and configured
- [ ] playwright.config.ts created with test ports (4000/4010/4020) and webServer
- [ ] test:kill-ports script kills test ports BEFORE E2E runs
- [ ] tests/e2e/ directory exists
- [ ] tests/unit/ directory exists
- [ ] Example E2E test has minimum 3 assertions (URL, element, data)
- [ ] `pnpm test` runs unit + E2E in sequence

**Web projects:**
- [ ] SEO meta tags in layout/head
- [ ] JSON-LD structured data included
- [ ] robots.txt created

**Infrastructure:**
- [ ] Dockerfile with multi-stage build (Docker projects)
- [ ] scripts/deploy.sh created (Dokploy projects)
- [ ] Multi-region deploy script (if multiregion selected)

**Database (MongoDB projects):**
- [ ] src/core/db/index.ts — MongoDB wrapper
- [ ] scripts/db-query.ts — Test Query Master
- [ ] scripts/queries/ directory
- [ ] db-query rules in CLAUDE.md

**Content (if web project with articles/posts):**
- [ ] scripts/build-content.ts
- [ ] scripts/content.config.json
- [ ] content/ directory

**Extras:**
- [ ] MCP servers installed (if selected)
- [ ] claude-mastery-project.conf created (if using profiles)
- [ ] No file > 300 lines
- [ ] All independent awaits use Promise.all

Report any missing items.
