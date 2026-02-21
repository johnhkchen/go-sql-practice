# Structure: Custom Health Route Implementation

## File Structure

### New Files
```
routes/
├── routes.go        # Package entry point with Register function
└── health.go        # Health endpoint handler
```

### Modified Files
```
main.go              # Add routes.Register(app) call
```

## File Specifications

### routes/routes.go

**Purpose**: Package entry point that registers all custom routes with PocketBase.

**Exports**:
- `Register(app core.App)` - Public function called from main

**Imports**:
```go
import (
    "github.com/pocketbase/pocketbase/core"
)
```

**Structure**:
```go
package routes

// Register registers all custom routes with the PocketBase app
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        // Register individual routes
        registerHealth(e)

        // Continue middleware chain
        return e.Next()
    })
}
```

### routes/health.go

**Purpose**: Implements the health check endpoint handler.

**Exports**: None (internal to package)

**Imports**:
```go
import (
    "github.com/labstack/echo/v5"
    "github.com/pocketbase/pocketbase/core"
)
```

**Functions**:
- `registerHealth(e *core.ServeEvent)` - Registers the health route
- `healthHandler(c echo.Context) error` - Handles health requests

**Structure**:
```go
package routes

// registerHealth registers the health check endpoint
func registerHealth(e *core.ServeEvent) {
    e.Router.GET("/api/health", healthHandler)
}

// healthHandler responds to health check requests
func healthHandler(c echo.Context) error {
    return c.JSON(200, map[string]string{
        "status": "ok",
    })
}
```

### main.go (Modified)

**Changes**: Add routes registration after migrations

**Import Addition**:
```go
import (
    "github.com/jchen/go-sql-practice/routes"
)
```

**Code Addition** (after line 14):
```go
// Register custom routes
routes.Register(app)
```

## Module Dependencies

### Direct Dependencies
- `github.com/pocketbase/pocketbase/core` - For App and ServeEvent types
- `github.com/labstack/echo/v5` - For Context type (transitive via PocketBase)

### Transitive Dependencies
Echo v5 is already a dependency of PocketBase, so no new dependencies are added to go.mod.

## Package Interfaces

### routes Package

**Public API**:
```go
// Register attaches all custom routes to the PocketBase app
func Register(app core.App)
```

**Internal Functions**:
```go
// registerHealth sets up the health endpoint
func registerHealth(e *core.ServeEvent)

// healthHandler processes health check requests
func healthHandler(c echo.Context) error
```

## Data Flow

1. `main()` calls `routes.Register(app)`
2. `Register()` binds to OnServe hook
3. When server starts, hook fires with ServeEvent
4. `registerHealth()` adds route to Echo router
5. HTTP request to `/api/health` triggers `healthHandler()`
6. Handler returns JSON response with status 200

## Error Handling

### Registration Phase
- OnServe hook returns error on failure (propagated from Next())
- Registration errors would prevent server startup

### Request Phase
- Health handler always returns success (no error cases)
- Echo handles JSON serialization errors internally

## Naming Conventions

### Package Name
- `routes` - Plural, indicating collection of route handlers

### Function Names
- `Register` - Public, imperative verb
- `registerHealth` - Private, camelCase with verb prefix
- `healthHandler` - Private, resource + "Handler" suffix

### File Names
- `routes.go` - Package main file
- `health.go` - Feature-specific file

## Type Usage

### PocketBase Types
- `core.App` - Application instance interface
- `core.ServeEvent` - Server lifecycle event

### Echo Types
- `echo.Context` - Request/response context

### Built-in Types
- `map[string]string` - JSON response body
- `error` - Standard error interface

## Architectural Boundaries

### Package Boundaries
- `routes` package is internal to the application
- Only `Register()` is exported
- Individual handlers are private

### Dependency Direction
```
main.go → routes → pocketbase/core
                 → echo/v5
```

### Separation of Concerns
- `main.go`: Application initialization
- `routes/routes.go`: Route registration orchestration
- `routes/health.go`: Specific endpoint implementation
- `migrations/`: Database schema management (unchanged)

## Extension Points

### Adding New Routes
1. Create new file in `routes/` (e.g., `metrics.go`)
2. Add register function (e.g., `registerMetrics()`)
3. Call from `Register()` in `routes.go`

### Adding Middleware
Can wrap handlers or use Echo's middleware:
```go
e.Router.GET("/api/health", healthHandler, middlewareFunc)
```

### Route Groups
Future routes can use groups:
```go
api := e.Router.Group("/api")
api.GET("/health", healthHandler)
```

## Summary

The structure establishes a clean `routes` package following the pattern set by `migrations`.
Files are organized by feature, with clear public/private boundaries. The implementation
requires creating two new files and a minor modification to main.go. No new external
dependencies are introduced.