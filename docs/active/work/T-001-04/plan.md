# Implementation Plan: Embed Astro in Go Binary

## Overview
This plan implements embedding Astro frontend files into the Go binary and serving them via PocketBase. Each step is designed to be independently verifiable and committable.

## Step 1: Create embed.go Foundation
**Goal**: Establish file embedding infrastructure

**Actions**:
1. Create `/embed.go` with basic embed directive
2. Add `getFrontendFS()` function that returns sub-filesystem
3. Add build validation to ensure frontend/dist exists

**Verification**:
- `go build` compiles successfully
- `getFrontendFS()` returns valid filesystem interface
- Build fails gracefully if frontend/dist missing

**Commit**: "feat: add embed infrastructure for frontend files"

**Code**:
```go
// embed.go
package main

import (
    "embed"
    "io/fs"
    "fmt"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS

func getFrontendFS() (fs.FS, error) {
    subFS, err := fs.Sub(frontendFiles, "frontend/dist")
    if err != nil {
        return nil, fmt.Errorf("failed to create frontend sub-filesystem: %w", err)
    }
    return subFS, nil
}

func frontendExists() bool {
    _, err := getFrontendFS()
    return err == nil
}
```

**Dependencies**: Current frontend build output must exist

## Step 2: Create Static Route Foundation
**Goal**: Set up route registration structure without serving logic

**Actions**:
1. Create `/routes/static.go` with placeholder registration function
2. Add call to `registerStatic(e)` in `/routes/routes.go`
3. Implement basic "not implemented" response

**Verification**:
- `go build` compiles successfully
- Server starts without errors
- Requests to `/` return "not implemented" response
- Admin UI `/_/` still works
- API routes `/api/*` still work

**Commit**: "feat: add static file route registration placeholder"

**Code**:
```go
// routes/static.go
package routes

import (
    "net/http"
    "github.com/pocketbase/pocketbase/core"
)

func registerStatic(e *core.ServeEvent) {
    e.Router.GET("/*", func(c echo.Context) error {
        return c.String(http.StatusNotImplemented, "Static serving not implemented")
    })
}
```

**Dependencies**: Step 1 completed

## Step 3: Implement SPA Filesystem Wrapper
**Goal**: Create filesystem wrapper that handles SPA routing

**Actions**:
1. Create `spaFS` struct that wraps embed.FS
2. Implement `Open()` method with fallback to index.html
3. Add path exclusion logic for API and admin routes

**Verification**:
- Unit test: `spaFS.Open("index.html")` succeeds
- Unit test: `spaFS.Open("nonexistent.js")` falls back to index.html
- Unit test: API paths are not handled by SPA logic

**Commit**: "feat: implement SPA filesystem wrapper with fallback logic"

**Code**:
```go
// Addition to routes/static.go
type spaFS struct {
    fs fs.FS
}

func (s *spaFS) Open(name string) (fs.File, error) {
    // Try exact match first
    file, err := s.fs.Open(name)
    if err == nil {
        return file, nil
    }

    // Fallback to index.html for SPA routing
    if !strings.HasPrefix(name, "api/") && !strings.HasPrefix(name, "_/") {
        return s.fs.Open("index.html")
    }

    return nil, err
}
```

**Dependencies**: Step 2 completed

## Step 4: Connect Static Serving
**Goal**: Replace placeholder with actual file serving

**Actions**:
1. Update `registerStatic()` to use Echo's StaticFS middleware
2. Configure middleware to use spaFS wrapper
3. Set appropriate cache headers

**Verification**:
- GET `/` returns frontend index.html
- GET `/_/` still serves admin UI
- GET `/api/*` routes still work
- GET `/nonexistent-route` serves index.html (SPA fallback)
- Static files have appropriate cache headers

**Commit**: "feat: implement static file serving with Echo StaticFS"

**Code**:
```go
// Update to routes/static.go registerStatic function
func registerStatic(e *core.ServeEvent) {
    frontendFS, err := getFrontendFS()
    if err != nil {
        e.App.Logger().Error("Failed to get frontend filesystem", "error", err)
        return
    }

    spaFilesystem := &spaFS{fs: frontendFS}
    e.Router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
        Root:       "/",
        Filesystem: http.FS(spaFilesystem),
        Browse:     false,
    }))
}
```

**Dependencies**: Step 3 completed

## Step 5: Create Build Automation
**Goal**: Provide automated build process

**Actions**:
1. Create `Makefile` with frontend, backend, and combined build targets
2. Add clean and development targets
3. Include error handling for missing dependencies

**Verification**:
- `make frontend` builds Astro successfully
- `make backend` compiles Go binary
- `make build` runs both in correct order
- `make clean` removes all build artifacts
- Build fails gracefully if npm not available

**Commit**: "feat: add Makefile for automated build process"

**Code**:
```make
# Makefile
.PHONY: build frontend backend clean dev test

build: frontend backend

frontend:
	@echo "Building frontend..."
	cd frontend && npm ci && npm run build

backend:
	@echo "Building backend..."
	go build -o go-sql-practice

clean:
	@echo "Cleaning build artifacts..."
	rm -rf frontend/dist
	rm -f go-sql-practice

dev:
	@echo "Starting development server..."
	./go-sql-practice serve --http="127.0.0.1:8090"

test:
	go test ./...
```

**Dependencies**: Step 4 completed

## Step 6: End-to-End Verification
**Goal**: Verify all acceptance criteria are met

**Actions**:
1. Run full build process from clean state
2. Test all specified endpoints and functionality
3. Verify admin UI and API routes remain functional
4. Test SPA routing behavior

**Verification Checklist**:
- [ ] `make build` succeeds from clean checkout
- [ ] Binary contains embedded frontend files
- [ ] GET `http://localhost:8090/` shows Astro index page
- [ ] GET `http://localhost:8090/_/` shows PocketBase admin UI
- [ ] GET `http://localhost:8090/api/health` works (existing endpoint)
- [ ] Frontend SPA routing works (random paths serve index.html)
- [ ] Binary is self-contained (works without frontend/ directory)

**Commit**: "feat: complete Astro frontend embedding in Go binary"

**Dependencies**: Step 5 completed

## Testing Strategy

### Unit Tests
**File**: `embed_test.go`
```go
func TestGetFrontendFS(t *testing.T) {
    fs, err := getFrontendFS()
    assert.NoError(t, err)
    assert.NotNil(t, fs)
}

func TestSPAFallback(t *testing.T) {
    // Test SPA filesystem wrapper
}
```

### Integration Tests
**File**: `routes/static_test.go`
```go
func TestStaticRouteRegistration(t *testing.T) {
    // Test route registration doesn't break existing routes
}

func TestFrontendServing(t *testing.T) {
    // Test actual HTTP requests to frontend endpoints
}
```

### Build Tests
**Commands**:
```bash
# Test clean build
make clean && make build

# Test individual targets
make frontend
make backend

# Test error cases
rm -rf frontend/dist && make backend  # Should fail gracefully
```

## Risk Mitigation

### Risk: Binary Size Growth
**Mitigation**: Monitor binary size, consider build optimization flags
**Detection**: Compare binary sizes before/after implementation

### Risk: Route Conflicts
**Mitigation**: Register static routes last, comprehensive testing
**Detection**: Test all existing endpoints after implementation

### Risk: Build Process Complexity
**Mitigation**: Clear error messages, documented prerequisites
**Detection**: Test on clean environment without dependencies

### Risk: Development Workflow Impact
**Mitigation**: Maintain separate dev server capability
**Detection**: Ensure hot reload still works for frontend development

## Rollback Plan

If issues arise:
1. **Step 1-2**: Remove files, revert route registration
2. **Step 3-4**: Disable static serving, return placeholder
3. **Step 5-6**: Use manual build process, remove Makefile

Each commit is designed to be individually revertable without breaking the application.

## Performance Considerations

1. **Build Time**: Frontend build adds ~30s to build process
2. **Binary Size**: Expect 1-5MB increase depending on assets
3. **Memory Usage**: Embedded files loaded into memory at startup
4. **Request Performance**: Static serving should be fast (no disk I/O)

## Dependencies

**External**:
- Node.js and npm (for frontend build)
- Make (for build automation)

**Internal**:
- Current PocketBase setup must remain functional
- Existing frontend build output in `frontend/dist/`
- Go 1.21+ (for embed package features)

This plan ensures each step is independently testable and the implementation progresses safely without breaking existing functionality.