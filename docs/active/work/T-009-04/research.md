# T-009-04 Research: Frontend API URL Configuration

## Overview

The frontend contains hardcoded API URLs that prevent it from working correctly when served from different environments. The issue primarily affects `StatsSummary.astro` but extends to multiple components across the codebase.

## Current Architecture

### Server Setup

**Go Backend (PocketBase)**
- Main server runs on port 8090 (configured in `Makefile:5` as `SERVER_PORT := 127.0.0.1:8090`)
- Uses PocketBase framework with custom routes registered in `routes/routes.go`
- Routes include: health, sync sessions, presentations, stats, links search, and views
- Frontend assets embedded via `internal/frontend/embed.go` from `frontend/dist/client`

**Frontend Build Process**
- Astro-based static frontend in `/frontend` directory
- Built using `astro build` to generate static files in `frontend/dist/client`
- Embedded into Go binary via `embed.FS` directive
- Served from same origin as API (same-origin deployment)

### API URL Patterns Found

**Inconsistent URL Configuration (Problem Areas)**

1. **StatsSummary.astro:141** - Hardcoded `http://127.0.0.1:8094/api/stats`
   - Only component using port 8094 instead of standard ports
   - Uses absolute URL instead of relative or configurable

2. **Most other components** - Two different approaches:
   - **Server-side (build time)**: `import.meta.env.PUBLIC_API_URL || 'http://localhost:8090'`
   - **Client-side (runtime)**: `(window as any).PUBLIC_API_URL || 'http://localhost:8090'`

**Files Using Proper Configuration Pattern:**
- `PresenterController.astro:34,1308` - Uses both server and client-side env vars
- `SearchInterface.astro:474` - Client-side env var
- `pages/tags/[slug].astro:16` - Server-side env var
- `pages/index.astro:15` - Server-side env var
- `pages/links/[id].astro:15,311` - Both server and client-side env vars
- `pages/watch/[id].astro:16` - Server-side env var
- `pages/sync/[id].astro:12` - Server-side env var
- `pages/present/[id].astro:20` - Server-side env var with window injection
- `pages/present/index.astro:9` - Server-side env var
- `pages/sync/[id]/control.astro:23` - Server-side env var

### Environment Configuration

**Development Environment**
- `.env` file in `frontend/` sets `PUBLIC_API_URL=http://127.0.0.1:8093`
- Note: This conflicts with actual server port (8090) creating development confusion
- Astro dev server runs separately from Go server during development

**Production Environment**
- Frontend built and embedded into Go binary
- Served from same origin, making relative URLs possible
- No runtime environment variable injection currently implemented

### Port Usage Analysis

**Observed Ports:**
- 8090: Standard server port (Makefile, most components)
- 8093: Development environment (.env file)
- 8094: Hardcoded in StatsSummary.astro (anomaly)

**Port Inconsistencies:**
- StatsSummary.astro hardcodes 8094, which doesn't match any standard config
- .env file uses 8093, different from Makefile's 8090
- Most components default to localhost:8090 when env var missing

## File Structure Analysis

### Component Distribution
- 17 total `.astro` files identified
- 13 files properly use environment variable pattern
- 1 file (StatsSummary.astro) uses hardcoded URL with wrong port
- Several files use both server-side and client-side API base configurations

### Key Components by Category

**Page-level Components:**
- `pages/index.astro` - Main listing page
- `pages/links/[id].astro` - Individual link view
- `pages/tags/[slug].astro` - Tag-filtered views
- `pages/stats.astro` - Statistics dashboard (includes StatsSummary)
- `pages/sync/[id].astro` - Sync session views
- `pages/present/` - Presentation system pages

**Shared Components:**
- `components/StatsSummary.astro` - **Problem component** with hardcoded URL
- `components/PresenterController.astro` - Complex component with proper env var usage
- `components/SearchInterface.astro` - Search functionality
- `components/Navigation.astro` - Site navigation
- `components/LinkCard.astro`, `LinksList.astro`, `GoLive.astro` - Supporting components

## Current Problems Identified

### StatsSummary.astro Specific Issues
- Line 141: `fetch('http://127.0.0.1:8094/api/stats')` hardcoded
- Wrong port (8094) doesn't match server configuration
- Absolute URL prevents same-origin deployment benefits
- No fallback or configuration mechanism

### Broader Configuration Issues
- Environment variable inconsistency (8093 vs 8090)
- Mix of server-side and client-side configuration approaches
- No unified API base URL management
- Development vs production environment mismatch

### Development Experience Problems
- Frontend .env points to 8093, server runs on 8090
- Developers need to manually align ports
- StatsSummary will fail with connection errors in all environments

## Codebase Constraints

### Astro Framework Constraints
- `import.meta.env` available at build time (SSR/SSG)
- Client-side scripts need runtime environment access via window object
- Static build output requires build-time configuration

### PocketBase Integration
- Custom routes registered through middleware system
- Frontend embedded as static assets in Go binary
- Same-origin serving means relative URLs are viable

### Current Patterns
- Established pattern: `import.meta.env.PUBLIC_API_URL || fallback`
- Window injection pattern: `window.PUBLIC_API_URL = API_BASE`
- Fallback to localhost:8090 in most components

## Related Files and Dependencies

**Core Configuration:**
- `frontend/.env` - Development environment variables
- `Makefile` - Server port configuration
- `astro.config.mjs` - Build configuration (static output)

**Backend Integration:**
- `internal/frontend/embed.go` - Asset embedding
- `routes/routes.go` - API route registration
- `routes/stats.go` - Stats endpoint implementation

**Component Dependencies:**
- StatsSummary.astro depends on `/api/stats` endpoint
- Other components use similar patterns for different endpoints
- Shared styling and component architecture