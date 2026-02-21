# Design: Embed Astro in Go Binary

## Design Options

### Option 1: Embed in main.go
Add embed directive directly to main.go and configure static serving there.

**Pros:**
- Minimal file changes
- All initialization in one place
- Simple to understand

**Cons:**
- Mixes concerns (embed logic with app initialization)
- main.go becomes larger
- Harder to test embed functionality separately

### Option 2: Separate embed.go file
Create dedicated `embed.go` in project root with embed directive and serving logic.

**Pros:**
- Clean separation of concerns
- Easier to test and modify independently
- Follows Go convention for embed files
- Can be conditionally compiled if needed

**Cons:**
- One more file to maintain
- Need to ensure initialization order

### Option 3: Add to routes package
Include embed and static serving in routes/static.go.

**Pros:**
- Groups all routing logic together
- Consistent with existing pattern
- Natural place for serving logic

**Cons:**
- Embed path becomes more complex (../frontend/dist)
- Couples routes package to frontend structure

## Static Serving Approach

### Approach A: Echo Static Middleware
Use Echo's built-in StaticFS with embed.FS directly.

**Pros:**
- Native Echo integration
- Well-tested middleware
- Handles caching headers automatically

**Cons:**
- May conflict with API routes if not configured carefully
- Less control over routing priority

### Approach B: Custom Handler
Write custom handler that serves from embed.FS with fallback logic.

**Pros:**
- Full control over routing logic
- Can implement SPA fallback precisely
- Clear precedence rules

**Cons:**
- More code to write and maintain
- Need to handle caching headers manually
- Risk of bugs in custom implementation

### Approach C: Middleware Chain
Use middleware to check paths and serve conditionally.

**Pros:**
- Integrates cleanly with existing middleware
- Can short-circuit for API/admin paths
- Flexible routing rules

**Cons:**
- More complex logic
- Performance overhead of checking each request

## Build Process Design

### Manual Two-Step
Document commands in README:
```bash
cd frontend && npm run build && cd ..
go build
```

**Pros:**
- Simple, no tooling needed
- Developers understand what's happening

**Cons:**
- Easy to forget frontend build
- No automation

### Makefile Automation
Create Makefile with build targets.

**Pros:**
- Single command builds everything
- Standard Go project pattern
- Can add other targets (test, clean)

**Cons:**
- Requires make installed
- Windows compatibility issues

### Go Generate
Use `//go:generate` to trigger frontend build.

**Pros:**
- Part of Go toolchain
- Runs with `go generate`

**Cons:**
- Requires npm in path
- Not all developers know go:generate

## Selected Design

### Architecture Decision: **Option 2 - Separate embed.go file**

Creating a dedicated `embed.go` file provides the cleanest separation of concerns while maintaining simplicity. This file will live in the project root alongside main.go, making the embed path straightforward (`frontend/dist/*`).

### Serving Decision: **Approach A - Echo Static Middleware with Custom Setup**

We'll use Echo's StaticFS but configure it carefully to avoid conflicts:
1. Register static serving last in the route chain
2. Use custom filesystem wrapper to handle SPA routing
3. Skip paths starting with `/api/` and `/_/`

### Build Decision: **Makefile Automation**

A Makefile provides the best balance of automation and transparency. We'll include:
- `make build` - builds frontend then backend
- `make frontend` - builds just frontend
- `make backend` - builds just backend
- `make clean` - removes build artifacts

## Implementation Details

### embed.go Structure
```go
package main

import (
    "embed"
    "io/fs"
    "net/http"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS

func getFrontendFS() (fs.FS, error) {
    return fs.Sub(frontendFiles, "frontend/dist")
}
```

### Route Registration
In routes/routes.go or new routes/static.go:
```go
func registerStatic(e *core.ServeEvent) {
    // Get sub-filesystem
    // Configure middleware
    // Register with Echo
}
```

### SPA Fallback Logic
For single-page app routing:
1. Try to serve exact file match
2. If not found and path doesn't start with /api or /_
3. Serve index.html instead
4. Let frontend router handle the path

### Build Order
1. Check if frontend/node_modules exists
2. Run npm install if needed
3. Run npm run build
4. Run go build
5. Output binary with embedded files

## Rationale

This design balances several concerns:

1. **Separation**: Keeping embed logic separate maintains code clarity
2. **Integration**: Using Echo's native features reduces custom code
3. **Automation**: Makefile provides reliable build process
4. **Compatibility**: Works with existing PocketBase structure
5. **Maintainability**: Clear boundaries between frontend and backend

The approach ensures that:
- Admin UI at `/_/` remains accessible
- API routes at `/api/*` continue working
- Frontend gets all other routes
- Binary is self-contained with all assets
- Build process is repeatable and documented

## Rejected Approaches

### Rejected: Custom HTTP Handler
Too much reimplementation of existing Echo features. Risk of security issues and missing edge cases.

### Rejected: Embed in routes package
The relative path complexity (`../frontend/dist`) makes this fragile. Better to keep embed at root level.

### Rejected: Go Generate
While elegant, not all developers are familiar with go:generate. Makefile is more universal and visible.

### Rejected: Dynamic File Serving
Serving files from disk at runtime defeats the purpose of creating a self-contained binary.