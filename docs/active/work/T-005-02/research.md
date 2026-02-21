# Research: Build Orchestration for Go-SQL-Practice

## Project Structure

The go-sql-practice project is a hybrid Go/PocketBase application with an embedded Astro frontend.

### Root Level
- `main.go` - Entry point using PocketBase framework
- `go.mod/go.sum` - Go dependencies (PocketBase v0.36.5, Echo v5.0.4)
- `Makefile` - Existing basic build orchestration (23 lines)
- `.gitignore` - Excludes build artifacts (go-sql-practice binary, frontend/dist/, pb_data/)

### Backend Components
- `routes/` - Custom PocketBase route handlers (9 files)
  - `static.go` - SPA filesystem handler using embedded frontend assets
  - `routes.go` - Route registration entry point
  - API handlers for health, links, presentations, stats, sync
- `migrations/` - Database migration files
- `internal/frontend/embed.go` - Go embed directive for frontend assets
  - Embeds `../../frontend/dist/*` using `//go:embed`
  - Serves from `frontend/dist/client` subdirectory

### Frontend Components
- `frontend/` - Astro static site generator project
  - `package.json` - Node.js v24+ required, Astro v5.17.3
  - `astro.config.mjs` - Static output with directory format
  - `src/` - Astro source files (pages, components, layouts)
  - `dist/` - Build output (client/ and server/ subdirectories)
  - Build generates both client and server outputs

### Current Build Pipeline

The existing Makefile provides basic orchestration:

```makefile
build: frontend backend
frontend: cd frontend && npm ci && npm run build
backend: go build -o go-sql-practice
clean: rm -rf frontend/dist && rm -f go-sql-practice
dev: ./go-sql-practice serve --http="127.0.0.1:8090"
test: go test ./...
```

### Asset Embedding Mechanism

The Go binary embeds the frontend via:
1. Astro builds static assets to `frontend/dist/`
2. `internal/frontend/embed.go` embeds `frontend/dist/*` at compile time
3. `routes/static.go` serves embedded assets via Echo middleware
4. SPA routing handled by fallback to `index.html` for non-API routes

### Dependencies and Constraints

#### Frontend Dependencies
- Node.js v24+ (specified in package.json engines)
- npm for package management
- Astro v5.17.3 with @astrojs/node adapter
- Build outputs to `dist/client` (served) and `dist/server` (unused)

#### Backend Dependencies
- Go 1.26+ (specified in go.mod)
- PocketBase framework for database and API
- Echo v5 for HTTP routing
- No external build tools beyond standard Go toolchain

#### Runtime Dependencies
- SQLite databases in `pb_data/` (excluded from git)
- Compiled Go binary (`go-sql-practice`)
- No separate web server required (PocketBase handles HTTP)

### Current Gaps

#### Missing Makefile Features
- No `make dev` hot reload or file watching
- `make clean` doesn't remove `pb_data/`
- No help target or target documentation
- No independent target verification
- No error handling or build validation

#### Development Workflow Issues
- Manual restart required for Go code changes
- Frontend changes require full rebuild
- No integrated testing approach
- Build pipeline can fail silently

#### Asset Embedding Edge Cases
- Embed path hardcoded to `frontend/dist/client`
- No validation that frontend was built before Go compilation
- `go build` succeeds even with missing frontend assets (embed fails at runtime)
- Embed directive processes all files, including potentially unnecessary server output

### Testing Infrastructure

- `make test` runs `go test ./...` but no Go test files found in current codebase
- No frontend testing configured
- No integration testing between Go backend and embedded frontend
- No build verification or smoke tests

### Current Dependencies Dependencies Satisfied

The ticket depends on T-001-04, which appears to be satisfied based on the existing embedded frontend infrastructure in `internal/frontend/embed.go` and `routes/static.go`.

### Environment and Tooling

- Development environment assumes Unix-like system (Makefile)
- Uses standard tools: make, npm, go
- No containerization or cross-platform build considerations
- No CI/CD integration visible
- Multiple background servers running (ports 8090-8095) suggest active development

This research provides the foundation for designing an improved build orchestration system that addresses the identified gaps while building on the existing solid foundation.