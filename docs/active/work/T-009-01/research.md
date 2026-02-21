# T-009-01 Research: Fix Static File Serving

## Problem Statement

Static file serving is completely broken in the application. The `registerStatic()` function in `routes/static.go` has an early return on line 36 that bypasses all serving logic, and the function call is commented out in `routes/routes.go` on line 20. As a result, hitting the root URL `/` returns nothing instead of serving the Astro frontend.

## Codebase Architecture

### Entry Point
- `main.go`: Creates PocketBase app instance and registers migrations/routes
- Uses PocketBase v0.36.5 as a Go framework (not standalone binary)
- Application structure: `pocketbase.New()` → `migrations.Register()` → `routes.Register()` → `app.Start()`

### Route Registration System
- `routes/routes.go`: Central route registration via `OnServe()` hook
- Pattern: Each feature has its own `register*()` function called from main registration
- Active routes: health, sync sessions, presentations, stats, links (search/view)
- **Broken**: `registerStatic(e)` commented out on line 20 with PocketBase v0.36.5 compatibility note

### Static File Serving Components

#### Embedded Frontend Assets (`internal/frontend/embed.go`)
- Embeds frontend files using `//go:embed frontend/dist/*`
- `GetFrontendFS()`: Creates sub-filesystem from `frontend/dist/client`
- `FrontendExists()`: Checks if frontend assets are available
- Current state: Only contains `placeholder.txt` (frontend not built)

#### Static File Handler (`routes/static.go`)
- Contains SPA filesystem wrapper (`spaFS`) for client-side routing
- Fallback logic: serves `index.html` for non-API routes (excluding `/api/` and `/_/`)
- **Critical issue**: Early return on line 36 completely disables function
- Commented out Echo v5 middleware implementation for static serving
- Multiple TODO comments reference PocketBase v0.36.5 compatibility issues

#### Dependencies and Compatibility
- `go.mod`: Uses Echo v5.0.4 and PocketBase v0.36.5
- Imports commented out in `static.go`: Echo v5 middleware and net/http packages
- Build system: Makefile orchestrates frontend build → backend build → embedded assets

### Frontend Integration
- Astro-based frontend in `/frontend/` directory
- Build output expected at `frontend/dist/client/` (referenced in embed path)
- Current state: Build artifacts missing, only placeholder present
- Frontend build creates assets that get embedded via Go's embed directive

### PocketBase Integration
- Admin UI served at `/_/` (protected by spaFS logic)
- API routes at `/api/*` (custom routes registered separately)
- Built-in PocketBase routes for collections, auth, admin interface
- Static serving should not shadow these existing routes

## Current Static Serving Flow (Broken)

1. `main.go` → `routes.Register()` → `OnServe()` hook
2. Route registration skips `registerStatic(e)` (commented out)
3. If called, `registerStatic()` immediately returns (line 36)
4. No static file middleware registered with Echo router
5. Root path `/` has no handler → returns nothing

## Working Components
- PocketBase core functionality intact
- Custom API routes functioning (health, presentations, etc.)
- Admin UI accessible at `/_/`
- Frontend assets embeddable (when built)
- SPA routing logic implemented but unused

## Technical Context

### Echo v5 Integration
- PocketBase v0.36.5 uses Echo framework internally
- Static serving should use Echo's static middleware
- Current implementation attempts `middleware.StaticWithConfig()` pattern
- http.FS adapter pattern for serving embedded filesystem

### Filesystem Structure
- Embedded path: `frontend/dist/client` → Go embed.FS
- Serving path: root `/` with SPA fallback to `index.html`
- Protected paths: `/api/` and `/_/` (skip SPA fallback)

## Root Causes

1. **Immediate return**: Line 36 in `registerStatic()` prevents any static serving setup
2. **Commented registration**: Line 20 in `routes.go` disables static route registration
3. **Missing frontend assets**: Build output directory only contains placeholder
4. **Version compatibility**: All TODOs reference PocketBase v0.36.5 compatibility issues

The implementation exists but is completely disabled. The core logic (SPA filesystem wrapper, path protection, embed integration) appears sound but untested due to the early return.

## Dependencies and Constraints

- Must preserve PocketBase admin UI at `/_/`
- Must not interfere with `/api/*` custom routes
- Must work with Echo v5 middleware system
- Must serve embedded Astro build output on root path
- Frontend must be built before backend build (embed requirement)

## Testing Strategy Implications

- Integration tests need working static serving to verify full stack
- Frontend-backend integration currently impossible to test
- HTTP handler tests require static middleware registration
- Build validation needs to verify embedded assets are served

The research confirms the ticket description: static file serving is completely broken due to deliberate disabling (early return + commented registration), with all issues stemming from PocketBase v0.36.5 compatibility concerns that appear to have been left unresolved.