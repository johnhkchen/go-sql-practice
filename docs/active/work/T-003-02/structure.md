# Structure: View Count Endpoint (T-003-02)

## File Changes Overview

Based on the design decision to implement atomic SQL increment with PocketBase record retrieval, this structure defines the minimal set of file changes required to implement `POST /api/links/:id/view`.

## File Modifications

### 1. New File: `routes/links_view.go`

**Purpose**: Contains the complete implementation of the view count increment endpoint.

**Location**: `routes/links_view.go`

**Content Structure**:
```go
package routes

import (
    // Standard imports
    "net/http"

    // PocketBase imports
    "github.com/pocketbase/pocketbase/core"
    "github.com/pocketbase/pocketbase/daos"
    "github.com/pocketbase/pocketbase/tools/rest"
)

// Registration function following established pattern
func registerLinksView(e *core.ServeEvent) {
    e.Router.POST("/api/links/:id/view", handleLinksView(e.App))
}

// Handler function with error handling and response formatting
func handleLinksView(app core.App) func(*core.RequestEvent) error {
    // Implementation details defined in plan phase
}
```

**Public Interface**:
- `registerLinksView(e *core.ServeEvent)` - Registration function called from routes.go
- HTTP endpoint: `POST /api/links/:id/view`

**Internal Architecture**:
- Parameter extraction from URL path (`:id`)
- Atomic SQL execution using `app.DB()`
- Error handling for various failure modes
- Record retrieval using PocketBase DAO operations
- JSON response formatting

**Dependencies**:
- PocketBase core framework
- Standard HTTP and JSON packages
- Database access through `app.DB()`
- DAO operations through `daos` package

### 2. Modified File: `routes/routes.go`

**Current Content Analysis**:
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        registerHealth(e)
        registerSyncSessions(e)
        registerStats(e)
        registerLinksSearch(e)
        registerLinksSearchSimple(e)
        return e.Next()
    })
}
```

**Required Change**:
Add single line `registerLinksView(e)` in the registration sequence.

**Modified Content**:
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        registerHealth(e)
        registerSyncSessions(e)
        registerStats(e)
        registerLinksSearch(e)
        registerLinksSearchSimple(e)
        registerLinksView(e)    // <- NEW LINE
        return e.Next()
    })
}
```

**Impact**: Minimal change following established pattern.

## No File Changes Required

### Database Schema
- `links` collection already exists with `view_count` field
- Field definition supports NULL values (handled by COALESCE in SQL)
- No migration required

### Migration Files
- Collection creation migration already exists
- Seed data already populates sample view_count values
- No new migrations needed

### Configuration Files
- No changes to main.go required
- No changes to go.mod dependencies (using existing PocketBase)
- No new external dependencies

## Module Boundaries and Dependencies

### Route Module (`routes/`)
**Responsibility**: HTTP request handling, parameter extraction, response formatting

**Public Interface**:
- Registration functions following `register*()` pattern
- Integration with PocketBase's `OnServe()` hook system

**Dependencies**:
- PocketBase core framework
- Standard library (http, json)
- Database access through PocketBase App interface

**Internal Organization**:
- One file per major endpoint/feature
- Consistent naming: `links_view.go` for links view functionality
- Handler functions follow `handle*()` naming pattern

### Database Layer
**Responsibility**: Data persistence and atomic operations

**Interface Used**:
- `app.DB()` for direct SQL execution
- PocketBase DAO for record operations
- Standard SQL parameterized queries for safety

**Atomic Operation Boundary**:
- Single `UPDATE` statement for atomic increment
- Separate `SELECT` via DAO for record retrieval
- No application-level transactions required

## Implementation Architecture

### Request Flow
1. **Route Registration**: `registerLinksView()` registers endpoint with Echo router
2. **Request Handling**: Handler extracts ID parameter from URL path
3. **Atomic Increment**: Direct SQL execution atomically increments view_count
4. **Result Validation**: Check affected rows to detect missing links
5. **Record Retrieval**: Use PocketBase DAO to fetch updated record
6. **Response Generation**: Return JSON with full record or error

### Error Handling Architecture
**Layered Error Handling**:
- **Parameter Level**: Invalid ID format validation
- **Database Level**: SQL execution errors, constraint violations
- **Application Level**: Link not found, serialization errors
- **HTTP Level**: Appropriate status codes and JSON error responses

**Error Response Format**:
```go
// Success (200)
{
  "id": "...",
  "url": "...",
  // ... full PocketBase record
}

// Not Found (404)
{
  "error": "Link not found"
}

// Server Error (500)
{
  "error": "Internal server error"
}
```

### Data Flow Architecture

**Incoming Request**:
```
HTTP POST /api/links/abc123/view
↓
Echo Router (parameter extraction)
↓
Handler Function (validation)
```

**Database Operations**:
```
Atomic SQL Update:
UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?
↓
Check Affected Rows (0 = not found)
↓
PocketBase DAO Record Fetch:
dao.FindRecordById("links", "abc123")
```

**Response Generation**:
```
Record Object
↓
PocketBase JSON Serialization
↓
HTTP Response (200 or 404)
```

## Component Interfaces

### HTTP Interface
```go
// Endpoint registration
func registerLinksView(e *core.ServeEvent)

// Request handler signature
func handleLinksView(app core.App) func(*core.RequestEvent) error
```

### Database Interface
```go
// Direct SQL execution
result := app.DB().NewQuery(
    "UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?"
).Exec(linkId)

// Record retrieval
record := dao.FindRecordById("links", linkId)
```

### Response Interface
```go
// Success response
return e.JSON(http.StatusOK, record)

// Error response
return e.JSON(http.StatusNotFound, map[string]string{
    "error": "Link not found",
})
```

## Change Ordering Dependencies

### Implementation Sequence
1. **Create `routes/links_view.go`** - Complete implementation first
2. **Modify `routes/routes.go`** - Add registration call
3. **Build and Test** - Verify compilation and functionality

**Rationale**: Creating the implementation file first ensures all dependencies and imports are resolved before registering the route. This prevents compilation errors during development.

### Testing Dependencies
- Implementation must be complete before integration testing
- Database must contain seed data for realistic testing
- HTTP server must be running for endpoint testing

## Security Boundaries

### Input Validation
- URL parameter extraction handled by Echo router
- SQL injection prevention via parameterized queries
- No user input in SQL strings

### Access Control
- No authentication required (per acceptance criteria)
- Public endpoint accessible to all users
- No authorization checks implemented

### Data Protection
- Only public link data exposed in responses
- No sensitive user data in view increment operation
- Standard PocketBase record serialization maintains data consistency

## Performance Characteristics

### Database Operations
- **Single UPDATE**: O(1) operation on indexed ID field
- **Single SELECT**: O(1) operation on indexed ID field
- **Total**: 2 database operations per request

### Memory Usage
- **Minimal**: Single record object in memory
- **Garbage Collection**: Short-lived objects, minimal GC pressure

### Concurrency
- **Thread-Safe**: Atomic SQL operations handle concurrent requests
- **No Shared State**: Handler functions are stateless
- **Database Locking**: SQLite row-level locking ensures consistency

## Rollback Strategy

### Deployment Safety
- **Additive Changes**: New file + single line addition
- **No Breaking Changes**: Existing functionality unaffected
- **Easy Rollback**: Remove registration line and delete new file

### Database Safety
- **No Schema Changes**: Uses existing view_count field
- **No Data Migration**: Works with existing records
- **Backward Compatible**: Existing records continue to function normally

This structure provides a clear blueprint for implementing the atomic view count increment while maintaining the existing codebase patterns and ensuring minimal risk of introducing bugs or breaking changes.