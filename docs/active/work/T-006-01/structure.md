# Structure: sync-sessions-collection

## File Changes Overview

### Modified Files
1. `migrations/collections.go` - Add sync_sessions collection creation
2. `routes/routes.go` - Register sync sessions routes

### New Files
3. `routes/sync_sessions.go` - Custom API routes for session management

## Detailed File Structure

### 1. migrations/collections.go (MODIFY)

**Location**: Add new function after existing collection creation
**Function**: `createSyncSessionsCollection(txApp core.App) error`

**Structure**:
```go
func createSyncSessionsCollection(txApp core.App) error {
    // Check if collection exists
    // Return nil if exists

    // Create collection instance
    // Type: core.CollectionTypeBase
    // Name: "sync_sessions"

    // Add fields:
    // - progress: NumberField (min: 0, max: 1, default: 0)
    // - admin_token: TextField (required, min: 64, max: 64)

    // Set API rules:
    // - ListRule: "" (empty = public)
    // - ViewRule: "" (empty = public)
    // - CreateRule: nil (no direct creation)
    // - UpdateRule: nil (no direct updates)
    // - DeleteRule: nil (no deletion)

    // Save collection
    // Return error if save fails
}
```

**Integration point**: Call from `createCollections()` after existing collections

### 2. routes/routes.go (MODIFY)

**Location**: In `Register()` function, after `registerHealth(e)`
**Addition**: Single line to register sync sessions routes

```go
registerSyncSessions(e)
```

### 3. routes/sync_sessions.go (CREATE)

**Package**: `routes`

**Imports**:
- Standard: crypto/rand, encoding/base64, encoding/json, net/http
- PocketBase: pocketbase/core, pocketbase/forms, pocketbase/models
- Echo: labstack/echo/v5

**Constants**:
```go
const (
    TokenLength = 32 // bytes before encoding
    ProgressMin = 0.0
    ProgressMax = 1.0
)
```

**Types**:
```go
type CreateSessionResponse struct {
    ID         string  `json:"id"`
    AdminURL   string  `json:"admin_url"`
    ViewerURL  string  `json:"viewer_url"`
    Progress   float64 `json:"progress"`
}

type UpdateProgressRequest struct {
    Progress float64 `json:"progress"`
}
```

**Functions**:

```go
// registerSyncSessions(e *core.ServeEvent)
// - Register POST /api/sync-sessions
// - Register PATCH /api/sync-sessions/:id

// handleCreateSession(c echo.Context) error
// - Generate secure admin token
// - Create sync_sessions record
// - Return URLs without exposing token in response

// handleUpdateProgress(c echo.Context) error
// - Extract token from X-Admin-Token header
// - Validate token against stored value
// - Update progress field
// - Return success/error

// generateAdminToken() (string, error)
// - Use crypto/rand for secure generation
// - Base64 encode for URL safety
// - Return 44-char string (32 bytes base64)

// validateProgress(progress float64) error
// - Check range [0, 1]
// - Return error if out of bounds
```

## Module Organization

### Dependencies
The implementation relies on:
- PocketBase core for database operations
- Echo v5 for HTTP routing (already in go.mod)
- Standard library for crypto and encoding

### Data Flow

1. **Session Creation**:
   ```
   Client -> POST /api/sync-sessions
          -> generateAdminToken()
          -> Create record via PocketBase DAO
          -> Return URLs
   ```

2. **Progress Update**:
   ```
   Admin -> PATCH /api/sync-sessions/:id
         -> Validate X-Admin-Token header
         -> Find record by ID
         -> Compare tokens
         -> Update progress via DAO
         -> PocketBase broadcasts via SSE
   ```

3. **Viewer Access**:
   ```
   Viewer -> GET /api/collections/sync_sessions/records/:id
          -> PocketBase built-in handler
          -> Returns record (admin_token hidden by rules)
   ```

4. **Realtime Subscription**:
   ```
   Viewer -> SSE /api/realtime
          -> PocketBase built-in SSE
          -> Subscribes to sync_sessions collection
          -> Receives updates on progress change
   ```

## Interface Boundaries

### Public API
- `POST /api/sync-sessions` - Create new session
- `PATCH /api/sync-sessions/:id` - Update progress (requires token)
- Built-in PocketBase endpoints for viewing

### Internal Interfaces
- Migration system: Function added to existing pattern
- Route registration: Follows existing pattern in routes package
- Database access: Via PocketBase's DAO and forms

### Security Boundaries
- Admin token never exposed after creation
- Token validation happens in custom route handler
- Collection rules prevent direct record manipulation

## Error Handling

All functions return appropriate HTTP status codes:
- 201 Created - Session created successfully
- 200 OK - Progress updated
- 400 Bad Request - Invalid progress value
- 401 Unauthorized - Invalid or missing admin token
- 404 Not Found - Session ID doesn't exist
- 500 Internal Server Error - Database or system errors

## Testing Considerations

The structure supports testing at multiple levels:
- Unit tests for token generation and validation
- Integration tests for API endpoints
- Manual testing via PocketBase admin UI
- Frontend testing with real SSE subscriptions