# T-001-03: astro-project-init - Research

## Current State

### Repository Structure
The project is a Go-based application with the following existing structure:
- Root directory contains documentation and configuration files
- `.flox/` directory with environment configuration (manifest.toml)
- `docs/` directory with active tickets, stories, and work artifacts
- No `frontend/` directory exists yet
- No Go source files exist yet (go.mod not created)
- No existing build system beyond the flox environment setup

### Development Environment
Flox environment is configured with:
- Go 1.26.0 (via `go_1_26` package)
- Node.js 24.13.0 (via `nodejs_24` package)
- GNU Make (via `gnumake` package)
- Environment activated via `flox activate` command

### Project Context
This ticket is part of story S-001 (project-scaffold) which aims to:
- Create a Go binary embedding PocketBase as a library
- Expose custom API routes alongside PocketBase's built-in ones
- Serve an Astro frontend from the embedded filesystem
- Produce a single binary via `go build` containing everything

### Dependencies
- Depends on T-001-00 (pin-tool-versions) which is marked as done
- T-001-00 has pinned Node.js 24.13.0 and expects Astro 5.17.3
- T-001-01 will create the Go module structure
- T-001-02 will create the Makefile for build orchestration

### Related Tickets Analysis
- T-001-01 (create-go-module): Will create `go.mod` and `cmd/server/main.go`
- T-001-02 (create-makefile): Will orchestrate the build process
- T-001-04 (embed-frontend): Will use `embed` package to include `frontend/dist/`
- T-001-05 through T-001-10: Will implement the Go backend and PocketBase integration

## Astro Framework Requirements

### Version Constraints
- Astro 5.17.3 is the target version (stable as of 2026-02-21)
- Node.js 24.13.0 is available via flox environment
- Must configure for static output for embedding into Go binary

### Static Site Generation
- Astro supports static output mode via `output: 'static'` in config
- Build artifacts go to `dist/` directory by default
- Static output allows embedding via Go's `embed` package

### Project Structure Expectations
Standard Astro project structure:
- `astro.config.mjs` - Configuration file
- `package.json` - Dependencies and scripts
- `src/pages/` - Page components
- `src/layouts/` - Layout components (optional)
- `src/components/` - Reusable components (optional)
- `public/` - Static assets (optional)
- `dist/` - Build output (gitignored)

## Integration Points

### Build Pipeline
- Frontend build must complete before Go embed in T-001-04
- Makefile (T-001-02) will orchestrate the build sequence
- `npm run build` should produce `frontend/dist/` artifacts

### Serving Strategy
- Go binary will serve the embedded static files
- Root path (`/`) will serve the Astro frontend
- API routes will be at `/api/*` paths
- PocketBase admin UI typically at `/_/` or `/admin/`

### Development Workflow
- Frontend development can use Astro's dev server
- Production builds embed the static output
- Hot reload during development via Astro CLI
- Final integration testing with full Go binary

## File System Constraints

### Git Ignore Requirements
- `frontend/node_modules/` must be excluded (large, regeneratable)
- `frontend/dist/` must be excluded (build artifacts)
- `.gitignore` needs updates or creation

### Path Conventions
- Frontend at `frontend/` subdirectory (not root level)
- Keeps separation between Go backend and JS frontend
- Allows independent package management

## Technical Considerations

### Node.js Engine Requirement
- Package.json should specify `"engines": {"node": ">=24"}`
- Ensures compatibility with flox-provided Node.js version
- Prevents issues with older Node versions

### Astro Configuration
- Must set `output: 'static'` for embedding
- Default dev server port is 4321
- Build command is typically `astro build`
- Dev command is typically `astro dev`

### Testing Requirements
- Index page at `src/pages/index.astro` as smoke test
- Should render "Link Bookmarks" heading
- Validates basic Astro setup and build process

## Assumptions and Constraints

### Assumptions
- Astro 5.17.3 is compatible with Node.js 24.13.0
- Static output mode meets all frontend requirements
- No server-side rendering (SSR) needed
- No edge/serverless deployment requirements

### Constraints
- Must use exact Astro version 5.17.3 (not v6 beta)
- Must produce embeddable static files
- Must work within flox environment
- Frontend isolation in `frontend/` subdirectory

## Open Questions
- Asset optimization strategy for embedded files
- Cache busting approach for static assets
- Frontend routing strategy (client-side vs server-side)
- API client generation or manual fetch calls