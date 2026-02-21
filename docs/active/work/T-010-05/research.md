# Research for T-010-05: Infrastructure Cleanup

## Current State Analysis

### .gitignore Configuration
Located at `/.gitignore` (31 lines). Current exclusions include:
- Go binaries: `go-sql-practice`, `*.exe`
- PocketBase runtime: `pb_data/`
- Frontend artifacts: `frontend/node_modules/`, `frontend/dist/`
- Flox, Lisa runtime directories
- IDE and OS files

**Missing**: No entry for `*.test` files. The file `routes.test` (36MB test binary) exists in root and is tracked by git.

### Directory Structure Anomalies

**Duplicate frontend directory**
- `/frontend/frontend/src/` exists as a stray nested directory
- Main frontend source is at `/frontend/src/`
- The duplicate appears to be from a build script error

**Test database artifacts**
- `/routes/pb_data/` contains SQLite database files (data.db, auxiliary.db with WAL/SHM files)
- Total size ~2.5MB of test artifacts
- Should be generated in temp directories, not in source tree

**Embed path mismatch**
- `internal/frontend/embed.go` expects files at `frontend/dist/*` (line 9)
- Actual embed path used: `frontend/dist/client` (line 13)
- Build copies to `internal/frontend/frontend/dist/` but embed.go looks for `frontend/dist/*`
- Current workaround: files exist in both locations

### Makefile Analysis
Located at `/Makefile` (75 lines). Key findings:

**Build process**
- `frontend` target builds Astro, validates `dist/client` exists
- `backend` target depends on frontend, builds Go binary
- No step to copy `frontend/dist` to `internal/frontend/frontend/dist/`

**Target naming inconsistency**
- Line 74: help text shows `validate` but actual target is `validate-build`
- Line 8: `.PHONY` correctly lists `validate-build`

**Missing targets**
- No `lint` target for Go formatting checks
- No `vet` target for Go static analysis
- Both are standard Go development practices

**Port configuration**
- `SERVER_PORT` set to `127.0.0.1:8090` (line 5)
- Used consistently in `dev` target

### CI Workflow
Located at `.github/workflows/ci.yml` (57 lines). Current setup:

**Versions**
- Uses `actions/setup-go@v4` (outdated, v5 available)
- Go version 1.26 (line 24)
- Node.js 24 (line 28)

**Build steps**
1. Checkout, setup Go/Node
2. Build frontend, backend
3. Validate build
4. Run tests
5. Smoke test binary
6. Check for uncommitted changes

**Missing checks**
- No `go vet` step
- No `gofmt` verification
- No frontend security audit (`npm audit`)
- No Go vulnerability check (`govulncheck`)

### Frontend Configuration

**Environment variables**
- `/frontend/.env` exists with `PUBLIC_API_URL=http://127.0.0.1:8090`
- No `.env.example` file to document required variables
- Port 8090 matches Makefile `SERVER_PORT`

### Error Handling in Routes

**Generic vs detailed errors**
Examined `/routes/stats.go`:
- Lines 76, 82, 88, 94, 100: Returns generic user-facing errors
- Pattern: `{"error": "Failed to get X"}`
- Does not expose internal error details (good security practice)
- No logging of actual errors for debugging (potential issue)

**Deleted file**
- `links_search_simple.go` mentioned in ticket is already deleted (per T-010-01)

### Project Documentation

**CLAUDE.md**
- Contains placeholder: "TODO: add a one-line project description here"
- Has basic directory conventions documented
- References RDSPI workflow

## Dependencies and Constraints

1. **Go embed directive**: Requires files at compile-time path specified in `//go:embed`
2. **Frontend build output**: Astro generates `dist/client/` structure
3. **PocketBase**: Uses SQLite databases, expects `pb_data/` runtime directory
4. **CI environment**: Ubuntu latest, no Flox available
5. **Testing**: Uses standard `go test`, creates test binaries with `.test` suffix

## File Access Patterns

- Build artifacts: Created in project root (`go-sql-practice`, `routes.test`)
- Frontend assets: Built to `frontend/dist/`, needed by Go embed
- Test databases: Created in working directory during test execution
- Runtime data: PocketBase creates `pb_data/` wherever binary runs

## Current Problems Summary

1. Test binaries (*.test) tracked in git
2. Stray duplicate directories from build errors
3. Test artifacts in source tree
4. Build doesn't ensure embed path alignment
5. Makefile help text mismatch
6. Missing standard Go quality targets (lint, vet)
7. Outdated CI action version
8. Missing CI quality checks
9. No .env.example documentation
10. Placeholder project description
11. Error logs not captured server-side