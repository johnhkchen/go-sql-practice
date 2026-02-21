# Research: Embed Astro in Go Binary

## Current State

### Go Application Structure

The Go application is built on PocketBase with a simple structure:
- **main.go:11-23**: Initializes PocketBase app, registers migrations and custom routes, starts server
- **go.mod:3**: Using Go 1.26 with PocketBase v0.36.5
- **routes/routes.go:8-16**: Route registration using OnServe hook with ServeEvent pattern
- **routes/health.go, sync_sessions.go**: Custom API endpoints already registered

### PocketBase Framework

PocketBase is serving on port 8094 (testing instance):
- Admin UI accessible at `/_/` path (confirmed working)
- Has built-in static file serving capabilities via Echo framework
- API routes mounted at `/api/*` by default
- Uses Echo v5 internally (dependency in go.mod:6)

### Frontend Structure

Astro frontend exists in `frontend/` directory:
- **package.json:7-9**: Build command outputs to `dist/` directory
- **astro.config.mjs:3-8**: Configured for static output, directory format
- **Current build output**: Single `index.html` file in `frontend/dist/`
- **index.html**: Complete single-page application with inline CSS, Link Bookmarks branding

## Key Observations

### Go Embed Directive

Go's embed package allows embedding files at compile time:
- `//go:embed` directive embeds files into the binary
- Can embed entire directory trees with patterns like `dist/*`
- Files become part of the binary, no external dependencies
- Access via `embed.FS` filesystem interface

### PocketBase Static Serving

PocketBase provides multiple ways to serve static files:
- `OnServe()` hook provides access to Echo router
- Echo has `Static()` and `StaticFS()` methods for file serving
- Can mount embedded filesystem using `echo.StaticFS()`
- Routes are evaluated in registration order

### Current Port Conflict

Port 8090 is occupied by unrelated service (qBittorrent):
- PocketBase instance running on 8094 for testing
- Need to ensure clean port 8090 for final deployment
- Multiple background processes attempting port 8090 (failing)

## Constraints and Dependencies

### Build Order Dependency

1. Frontend must be built before Go compilation
2. Astro build outputs to `frontend/dist/`
3. Go embed reads files at compile time
4. Binary contains snapshot of dist at build time

### Path Routing Priority

PocketBase route precedence:
1. Admin UI at `/_/` (protected, must not override)
2. API routes at `/api/*` (critical functionality)
3. Custom routes registered via OnServe
4. Static files should be catch-all for remaining paths

### File System Structure

Current frontend build output:
- Single HTML file (may expand with assets)
- No asset directory yet (will appear with more complex builds)
- All paths relative to dist root

## Technical Requirements

### Embed Implementation

Need to:
1. Import `embed` package in Go file
2. Declare embedded filesystem variable with directive
3. Point to `frontend/dist` directory (relative to module root)
4. Ensure dist exists before compilation

### Static Server Configuration

Must configure:
1. Mount embedded FS to Echo router
2. Set as fallback route (not override API/admin)
3. Handle index.html for SPA routing
4. Support future asset directories

### Build Process

Requires two-step build:
1. `npm run build` in frontend directory
2. `go build` in project root
3. Order is critical - frontend first

## Identified Files and Modules

### Files to Modify
- New Go file for embed (e.g., `embed.go` or add to `main.go`)
- Update route registration to include static serving
- Possible Makefile for build automation

### Files to Reference
- `routes/routes.go`: Pattern for route registration
- `main.go`: Application initialization point
- `frontend/dist/`: Source for embedded files

## Assumptions

1. PocketBase's Echo instance is accessible via ServeEvent
2. StaticFS can work with embed.FS interface
3. Frontend will remain SPA (single-page application)
4. Build will happen on development machine with npm available

## Risks and Considerations

1. Binary size will increase by size of frontend assets
2. Frontend changes require full Go rebuild
3. Development workflow needs hot reload solution
4. Must preserve PocketBase admin and API functionality
5. Path conflicts between frontend routes and API routes