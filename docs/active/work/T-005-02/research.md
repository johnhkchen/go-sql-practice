# Research: Build Orchestration (T-005-02)

## Project Architecture Overview

### Core Components
- **Backend**: Go application embedding PocketBase as a library
- **Frontend**: Astro static site generator configured for static output
- **Binary**: Single executable with embedded frontend assets
- **Database**: SQLite managed by PocketBase (pb_data/ directory)

### Current Build Infrastructure

#### Existing Makefile Structure
Located at root level with these targets:
- `build`: Orchestrates frontend → backend build sequence
- `frontend`: Builds Astro project using npm ci && npm run build
- `backend`: Compiles Go binary with `go build -o go-sql-practice`
- `clean`: Removes build artifacts (frontend/dist, binary)
- `dev`: Starts development server via `./go-sql-practice serve`
- `test`: Runs Go tests with `go test ./...`

### Frontend Build System

#### Astro Configuration (frontend/astro.config.mjs)
- Output mode: 'static' (suitable for embedding)
- Adapter: @astrojs/node standalone mode
- Build directory: frontend/dist/ with client/ and server/ subdirectories
- Assets directory: 'assets' with 'directory' format

#### Package Management
- Package manager: npm
- Node.js version requirement: >=24 (package.json engines)
- Build command: `astro build`
- Development command: `astro dev`
- Preview command: `astro preview`

#### Build Output Structure
```
frontend/dist/
├── client/           # Static assets for embedding
│   ├── index.html
│   ├── stats/index.html
│   └── assets/       # Bundled JS/CSS
└── server/           # Server-side files (not embedded)
```

### Go Backend Architecture

#### Module Structure
- Module path: `github.com/jchen/go-sql-practice`
- Go version: 1.26
- Main dependencies: PocketBase v0.36.5, Echo v5.0.4

#### Embedding Infrastructure
**Location**: `internal/frontend/embed.go`
- Uses `//go:embed ../../frontend/dist/*` directive
- Embedded filesystem accessible via `embed.FS`
- Sub-filesystem created for `frontend/dist/client` (line 13)
- Exported functions: `GetFrontendFS()`, `FrontendExists()`

#### Static File Serving
**Location**: `routes/static.go`
- SPA filesystem wrapper (`spaFS`) with index.html fallback
- Fallback logic excludes API routes (api/) and admin routes (_/)
- Integration with Echo middleware using `middleware.StaticConfig`
- Mounted at root path ("/") with browse disabled

#### Route Registration
**Location**: `routes/routes.go`
- Registers custom routes and static file serving
- Called from `main.go` during app initialization

### Development Environment

#### Flox Environment Management
- Environment name: "go-sql-practice"
- Go version: 1.26.0 (defined in .flox/env/manifest.toml)
- Node.js version: 24.x (defined in manifest)
- Currently not activated (go binary not in PATH)

#### Runtime Dependencies
- Node.js: Available via NVM (v24.3.0)
- npm: Available via global installation
- Go: Defined in Flox but not currently active

### Build Artifacts and Cleanup

#### Generated Files
- Go binary: `go-sql-practice` (35MB executable)
- Frontend build: `frontend/dist/` directory tree
- Database: `pb_data/` directory (excluded from git)
- Logs: `server.log` and various .log files

#### Gitignore Patterns
- Go binary and executables
- PocketBase runtime data (pb_data/)
- Frontend node_modules and dist
- Flox runtime directories (.flox/cache/, .flox/run/, .flox/log/)
- Lisa runtime (.lisa/, .lisa-layout.kdl)
- Log files and IDE directories

### Build Pipeline Dependencies

#### Dependency Chain
1. Frontend build must complete before Go compilation
2. Go embed reads from frontend/dist/client/ at compile time
3. Binary serves embedded assets at runtime
4. PocketBase manages database initialization and migrations

#### Current Pipeline Issues
- Go toolchain not available in current environment
- Makefile assumes go command is in PATH
- Dev target fails because binary doesn't exist after failed build
- No hot reload or file watching for development

### Test Infrastructure
- Go test framework used (`go test ./...`)
- No frontend-specific tests identified
- Test files present in routes/ and migrations/ directories
- Testing depends on Go toolchain availability

### Development Workflow Patterns
- Multiple background server instances running (ports 8091-8096)
- Development server expects binary to exist before starting
- No automatic rebuild on file changes
- Frontend hot reload would require separate Astro dev server

### Asset Embedding Details
- Embed path: `../../frontend/dist/*` relative to internal/frontend/
- Embedded subdirectory: `client/` (contains static files)
- Server subdirectory: Not embedded (server-side rendering artifacts)
- File access via standard fs.FS interface with SPA fallback logic

### Configuration Files
- Go module: go.mod, go.sum (dependency management)
- Frontend: package.json, package-lock.json (npm dependencies)
- Astro: astro.config.mjs (build configuration)
- Environment: .flox/env/manifest.toml (development tooling)
- Build: Makefile (orchestration targets)

### Constraints and Assumptions

#### Environment Constraints
- Requires Node.js 24+ for frontend builds
- Requires Go 1.26+ for backend compilation
- Flox environment provides tooling but may not be activated
- Development assumes Unix-like environment (Makefile syntax)

#### Architecture Constraints
- Static Astro output required for embedding
- Single binary deployment model
- PocketBase database tied to pb_data/ directory
- SPA routing requires index.html fallback logic

#### Build Order Requirements
- Frontend must build before backend compilation
- Go embed directive reads files at compile time
- Clean target must remove both frontend and backend artifacts