# Design: Custom Health Route Implementation

## Design Options

### Option 1: Single Routes Package with Register Function
Create a `routes/` package that mirrors the `migrations/` pattern.
- `routes/routes.go` with `Register(app core.App)` function
- `routes/health.go` with health handler
- Register all routes in OnServe hook

**Pros:**
- Consistent with existing migrations pattern
- Clear separation of concerns
- Easy to extend with more routes

**Cons:**
- Might be overkill for single route
- Extra package for simple functionality

### Option 2: Inline Registration in Main
Add route registration directly in `main.go` after migrations.
- Handler function defined in main.go
- OnServe hook registered in main()

**Pros:**
- Simplest implementation
- No extra packages needed
- Everything visible in one place

**Cons:**
- main.go becomes cluttered with growth
- Breaks single responsibility principle
- Inconsistent with migrations pattern

### Option 3: Routes as Part of Migrations
Extend migrations package to handle routes too.
- Add route registration to existing `migrations.Register()`
- Keep health handler in migrations package

**Pros:**
- No new packages needed
- Leverages existing hook

**Cons:**
- Conflates migrations with routing
- Poor semantic meaning
- Violates package cohesion

### Option 4: API Package with Handlers
Create `api/` package with handler modules.
- `api/health/health.go` for health endpoint
- `api/register.go` for registration logic

**Pros:**
- Scalable structure for API growth
- Clear module boundaries
- RESTful organization

**Cons:**
- Over-engineered for current needs
- Deep nesting for single handler

## Decision: Option 1 - Routes Package

### Rationale

1. **Consistency**: Matches the established pattern from migrations package
2. **Extensibility**: Provides clear location for future custom routes
3. **Separation**: Keeps routing logic separate from app initialization
4. **Simplicity**: Single registration point called from main

### Why Not Others

- Option 2: While simpler initially, it doesn't scale and makes main.go a dumping ground
- Option 3: Semantically incorrect - migrations should handle database, not HTTP routes
- Option 4: Premature optimization - too much structure for current requirements

## Implementation Design

### Package Structure
```go
// routes/routes.go
package routes

func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        registerHealth(e)
        return e.Next()
    })
}
```

### Health Handler
```go
// routes/health.go
func registerHealth(e *core.ServeEvent) {
    e.Router.GET("/api/health", healthHandler)
}

func healthHandler(c echo.Context) error {
    return c.JSON(200, map[string]string{
        "status": "ok",
    })
}
```

### Integration Point
```go
// main.go
migrations.Register(app)
routes.Register(app)  // Add after migrations
```

## Technical Decisions

### Route Path
Use `/api/health` as specified. While `/api/` is PocketBase's namespace, custom routes
can coexist as long as they don't conflict with collection names.

### Response Format
Return minimal JSON as specified: `{"status": "ok"}`
No additional fields like timestamp, version, or metrics for now.

### Status Code
Return 200 OK for healthy status. No need for detailed health states initially.

### Authentication
Route registered directly on router without auth middleware, making it public.
This is appropriate for health endpoints used by load balancers and monitoring.

### Error Handling
Health endpoint should not fail. Even if internal issues exist, returning 200
indicates the service is responsive. Detailed health checks can be added later.

## Alternative Considered

### Using Echo Groups
Could create an API group: `api := e.Router.Group("/api")`
Decided against because:
- Single route doesn't warrant grouping
- Might conflict with PocketBase's internal /api group
- Direct registration is clearer for simple case

### Returning More Data
Could include version, uptime, database status, etc.
Decided against because:
- Not in requirements
- Minimal endpoint is standard for basic health checks
- Can be extended later if needed

## Migration Strategy

1. Create routes package structure
2. Implement health handler
3. Register in main.go after migrations
4. Test endpoint accessibility
5. Verify no authentication required

## Future Considerations

- Add `/api/ready` for readiness vs liveness
- Include database connectivity check
- Add metrics endpoint for Prometheus
- Support middleware for custom routes
- Route documentation generation

## Summary

The routes package pattern provides a clean, extensible solution that mirrors the existing
migrations approach. It keeps main.go clean while providing a clear location for future
custom routes. The implementation is straightforward and satisfies all requirements without
over-engineering.