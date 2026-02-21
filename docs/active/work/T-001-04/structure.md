# Structure: Embed Astro in Go Binary

## Files to Create

### `/embed.go` - New File
This file will contain the embed directive and filesystem access functions.

**Purpose**: Isolate frontend file embedding from application logic
**Location**: Project root (same level as main.go)
**Size**: ~30 lines

**Components**:
- `//go:embed frontend/dist/*` directive
- `frontendFiles embed.FS` variable
- `getFrontendFS() (fs.FS, error)` function to provide sub-filesystem
- Optional: `frontendExists() bool` for build validation

### `/routes/static.go` - New File
This file will handle static file serving configuration for PocketBase.

**Purpose**: Route registration for embedded frontend files
**Location**: routes/ directory (consistent with existing pattern)
**Size**: ~60 lines

**Components**:
- `registerStatic(e *core.ServeEvent)` function
- Custom filesystem wrapper for SPA fallback
- Echo StaticFS middleware configuration
- Path exclusion logic for API and admin routes

### `/Makefile` - New File
Build automation for the two-phase build process.

**Purpose**: Coordinate frontend build then backend compilation
**Location**: Project root
**Size**: ~40 lines

**Components**:
- `build` target (frontend + backend)
- `frontend` target (npm run build)
- `backend` target (go build)
- `clean` target (remove artifacts)
- `dev` target (serve development mode)

## Files to Modify

### `/routes/routes.go` - Lines 12-13
Add call to `registerStatic(e)` in the route registration function.

**Change Type**: Addition
**Location**: After existing route registrations, before `return e.Next()`
**Impact**: Minimal, follows existing pattern

**Before**:
```go
registerHealth(e)
registerSyncSessions(e)
registerStats(e)
registerLinksSearch(e)

// Continue middleware chain
return e.Next()
```

**After**:
```go
registerHealth(e)
registerSyncSessions(e)
registerStats(e)
registerLinksSearch(e)
registerStatic(e)

// Continue middleware chain
return e.Next()
```

### `/main.go` - Import Section (Line 4-6)
No changes required. The embed.go file will be in the same package, so its exported functions are automatically available.

## Module Boundaries

### Frontend Module (`frontend/`)
**Responsibility**: Build static assets for embedding
**Interface**: Outputs files to `frontend/dist/`
**Dependencies**: Node.js, npm, Astro framework
**Boundary**: Files under `frontend/dist/` become input to embed

### Embed Module (`embed.go`)
**Responsibility**: Provide access to embedded frontend files
**Interface**: `getFrontendFS() fs.FS` function
**Dependencies**: Go embed package, frontend build output
**Boundary**: Exposes filesystem interface, hides embed implementation

### Static Serving Module (`routes/static.go`)
**Responsibility**: Configure PocketBase to serve embedded files
**Interface**: `registerStatic(e *core.ServeEvent)` function
**Dependencies**: embed module, PocketBase/Echo framework
**Boundary**: Handles HTTP routing, delegates file access to embed module

### Build Module (`Makefile`)
**Responsibility**: Orchestrate multi-stage build process
**Interface**: Make targets (build, frontend, backend, clean)
**Dependencies**: npm, go toolchain
**Boundary**: External process coordination

## Internal Organization

### Package Structure
All new Go code stays in `main` package:
- `embed.go`: File embedding
- `routes/static.go`: Route configuration
- No new packages introduced

### Function Visibility
- `getFrontendFS()`: Exported (used by routes package)
- Route handler functions: Unexported (internal to routes)
- Build targets: Public make interface

### Data Flow
1. **Build Time**: Makefile → npm build → go embed → binary
2. **Runtime**: HTTP request → routes → embed.FS → response

### Error Handling
- Build errors fail fast and loudly
- Runtime errors serve 404 for missing files
- Embed errors are compile-time failures

## Ordering of Changes

### Phase 1: Foundation
1. Create `embed.go` with basic structure
2. Create `routes/static.go` with placeholder
3. Test compilation (should succeed with empty implementation)

### Phase 2: Serving Logic
1. Implement filesystem wrapper in `routes/static.go`
2. Add route registration call to `routes/routes.go`
3. Test static serving with existing dist files

### Phase 3: Build Integration
1. Create `Makefile` with all targets
2. Test full build pipeline
3. Verify embedded files are served correctly

### Phase 4: Documentation
1. Update README with build instructions
2. Add comments to new code files

## File Dependencies

### Build Dependencies
```
Makefile → npm (frontend build) → go build (embed)
```

### Runtime Dependencies
```
HTTP request → routes/routes.go → routes/static.go → embed.go → frontend files
```

### Import Dependencies
```go
// embed.go
import (
    "embed"
    "io/fs"
)

// routes/static.go
import (
    "io/fs"
    "net/http"
    "strings"

    "github.com/pocketbase/pocketbase/core"
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)
```

## Architecture Constraints

### Must Not Break
- PocketBase admin UI at `/_/`
- API routes at `/api/*`
- Existing health and sync endpoints
- Database functionality

### Must Preserve
- Single binary deployment
- Development workflow compatibility
- Standard Go project structure

### Must Enable
- Frontend served at `/` root path
- SPA routing support (fallback to index.html)
- Static assets served with proper headers
- Build automation

## Testing Strategy

### Unit Tests
- Test `getFrontendFS()` function
- Test filesystem wrapper SPA fallback logic
- Mock embed.FS for isolated testing

### Integration Tests
- Test full request cycle through routes
- Verify path exclusions work correctly
- Test with actual built frontend files

### Build Tests
- Test Makefile targets work independently
- Test build fails appropriately when frontend missing
- Test clean removes all artifacts

This structure provides clear separation of concerns while maintaining integration with the existing PocketBase application. Each module has a single responsibility and well-defined interfaces.