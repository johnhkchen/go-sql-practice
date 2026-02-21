# Research: T-009-03 - fix-pathvalue-routing

## Problem Statement

`routes/presentations.go` and `routes/sync_sessions.go` use `e.Request.PathValue("id")` to extract route parameters, but the routes are defined using Echo-style syntax (`:id`) instead of Go 1.22 syntax (`{id}`). This mismatch causes `PathValue()` to return empty strings, causing all presentation and sync endpoints to silently fail.

## Route Definition vs Parameter Extraction Mismatch

### Current Problematic Patterns

**presentations.go:185, 219, 277:**
```go
// Route defined with Echo syntax
e.Router.GET("/api/presentations/:id/status", ...)

// Handler trying to use Go 1.22 PathValue
presentationID := e.Request.PathValue("id")  // Returns empty string
```

**sync_sessions.go:88:**
```go
// Route defined with Echo syntax
e.Router.POST("/api/sync/:id/progress", ...)

// Handler trying to use Go 1.22 PathValue
sessionID := e.Request.PathValue("id")  // Returns empty string
```

### Working Example in Codebase

**links_view.go** already recognized this issue and implemented manual path parsing:
```go
// Route uses Echo syntax
e.Router.POST("/api/links/:id/view", ...)

// Handler uses manual parsing instead of PathValue
path := e.Request.URL.Path
parts := strings.Split(path, "/")
// Manual extraction logic...
```

Comment in lines 21-22: "Parse path manually since PathParam may not be available"

## PocketBase v0.36.5 Router Architecture

### Historical Context
- PocketBase v0.23+ migrated from Echo router to Go 1.22's standard HTTP router
- Route syntax changed from Echo's `:param` to Go 1.22's `{param}`
- Parameter extraction changed from `c.PathParam("name")` to `e.Request.PathValue("name")`

### Current Implementation
- PocketBase v0.36.5 uses Go 1.22's `net/http.ServeMux` internally
- Routes are registered through `e.Router.GET/POST()` methods
- Parameter extraction should use `e.Request.PathValue()` for Go 1.22 syntax
- Old Echo syntax `:param` doesn't work with `PathValue()`

### Dependencies
```go
// go.mod line 3
go 1.26

// go.mod line 8
github.com/pocketbase/pocketbase v0.36.5

// go.mod line 6
github.com/labstack/echo/v5 v5.0.4  // Still present but not used for routing
```

## Affected Endpoints

### Presentation Endpoints (routes/presentations.go)
1. **GET /api/presentations/:id/status** (line 47)
   - Handler: `handleGetStatus` (line 183)
   - PathValue call: line 185
   - Impact: Status checks fail silently

2. **POST /api/presentations/:id/stop** (line 52)
   - Handler: `handleStopLive` (line 217)
   - PathValue call: line 219
   - Impact: Cannot stop live presentations

3. **POST /api/presentations/:id/live** (line 57)
   - Handler: `handleStartLive` (line 275)
   - PathValue call: line 277
   - Impact: Cannot start live presentations

### Sync Session Endpoints (routes/sync_sessions.go)
1. **POST /api/sync/:id/progress** (line 40)
   - Handler: `handleUpdateProgress` (line 85)
   - PathValue call: line 88
   - Impact: Progress updates fail silently

### Working Endpoints
- All routes without path parameters work correctly
- `/api/sync/create` works (no path parameters)
- `/api/stats`, `/api/links/search*` work (no path parameters)

## Route Registration Flow

1. **Entry Point**: `routes/routes.go:Register()`
2. **Registration**: `app.OnServe().BindFunc()` callback
3. **Route Setup**: Individual `register*()` functions called:
   - `registerPresentations(e)` (line 14)
   - `registerSyncSessions(e)` (line 13)
4. **Handler Binding**: Each route bound to handler function

## Testing Infrastructure

### Current State
- Basic test setup exists in `routes/routes_test.go`
- `setupTestApp()` creates in-memory PocketBase instance
- `makeRequest()` function exists but not fully implemented (lines 108-128)
- No integration tests for path parameter extraction

### Verification Challenges
- PocketBase testing requires proper router integration
- Current mock implementation doesn't test real routing
- Need proper HTTP request simulation for path parameter testing

## Solution Approaches

### Option 1: Update Route Syntax (Recommended)
Change route definitions from Echo syntax to Go 1.22 syntax:
```go
// From:
e.Router.GET("/api/presentations/:id/status", ...)

// To:
e.Router.GET("/api/presentations/{id}/status", ...)
```

### Option 2: Manual Path Parsing (Existing Pattern)
Follow the pattern in `links_view.go`:
```go
path := e.Request.URL.Path
parts := strings.Split(path, "/")
// Extract ID manually
```

### Option 3: Query Parameter Fallback
Extract from query parameters:
```go
id := e.Request.URL.Query().Get("id")
```

## Code Patterns and Conventions

### Handler Function Structure
```go
func handleXxx(e *core.RequestEvent, app core.App) error {
    // 1. Extract parameters
    id := e.Request.PathValue("param")

    // 2. Validate input
    if id == "" { return e.JSON(400, error) }

    // 3. Authenticate if needed
    user, err := getAuthenticatedUser(e)

    // 4. Database operations
    record, err := app.FindRecordById("collection", id)

    // 5. Return JSON response
    return e.JSON(200, response)
}
```

### Error Response Format
Consistent across all handlers:
```go
return e.JSON(statusCode, map[string]string{
    "error": "Error message",
})
```

## Security Considerations

### Authentication Patterns
- `getAuthenticatedUser(e)` extracts auth from request context
- `checkPresentationOwnership()` validates user permissions
- Admin token validation uses constant-time comparison

### Token Validation
```go
func validateToken(provided, stored string) bool {
    return subtle.ConstantTimeCompare([]byte(provided), []byte(stored)) == 1
}
```

## Dependencies and Constraints

### Runtime Dependencies
- Go 1.26 (supports PathValue method)
- PocketBase v0.36.5 (uses Go 1.22 router)
- Standard library HTTP routing

### Collection Dependencies
Routes depend on these PocketBase collections:
- `presentations` - for presentation lifecycle
- `sync_sessions` - for progress tracking
- Must exist for handlers to work

### Migration Considerations
- Route syntax change is breaking for any direct API clients
- Frontend code may need URL updates if hard-coded
- Consider backward compatibility during transition

## Implementation Files

### Core Files
- `routes/presentations.go` - 359 lines, 4 path parameter calls
- `routes/sync_sessions.go` - 172 lines, 1 path parameter call
- `routes/routes.go` - 30 lines, registration logic

### Supporting Files
- `routes/routes_test.go` - testing infrastructure
- `routes/links_view.go` - working manual parsing example

### Documentation
- `go.mod` - dependency versions
- This research document

## Risk Assessment

### High Impact Issues
1. **Silent Failures**: `PathValue("")` returns empty string, no error
2. **User Experience**: All presentation management broken
3. **API Reliability**: Sync progress updates fail

### Low Risk Areas
- Routes without parameters unaffected
- Database operations work correctly
- Authentication mechanisms intact

## Next Phase Inputs

For Design phase, key considerations:
1. Route syntax migration strategy (breaking vs non-breaking)
2. Testing approach for path parameter validation
3. Backward compatibility requirements
4. Error handling improvements (fail fast vs silent)
5. Frontend impact assessment