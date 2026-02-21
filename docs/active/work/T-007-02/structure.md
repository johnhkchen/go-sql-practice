# Structure: presentation-api-routes (T-007-02)

## Overview

This document defines the file-level changes, architecture boundaries, and implementation structure for the presentation API routes. The implementation follows established patterns while adding new functionality for presentation lifecycle management.

## File Changes

### New Files

**`routes/presentations.go`** (~180 lines)
- Location: `/home/jchen/repos/go-sql-practice/routes/presentations.go`
- Purpose: Main implementation file for presentation lifecycle routes
- Dependencies: PocketBase core, existing sync session utilities
- Exports: `registerPresentations(e *core.ServeEvent)` function

### Modified Files

**`routes/routes.go`** (minimal change)
- Action: Add `registerPresentations(e)` call to Register function
- Location: Line ~12, after `registerSyncSessions(e)`
- Integration point for new route registration

**`routes/sync_sessions.go`** (potential extraction)
- Action: Extract shared utilities if needed (token generation, validation)
- Alternative: Keep utilities private, duplicate minimal code
- Decision: Keep utilities private to avoid premature abstraction

### No Changes Required

**`main.go`** - Entry point unchanged, routes automatically registered
**`migrations/collections.go`** - Collections already exist from T-007-01
**Frontend files** - New routes complement existing APIs

## Architecture Boundaries

### Route Responsibility

**`routes/presentations.go`**
```
Responsibilities:
├── Presentation lifecycle control (start/stop live sessions)
├── Enhanced status information with computed fields
├── Authentication and authorization for presentation owners
└── Integration with sync session creation/management

Boundaries:
├── Does NOT handle sync session progress updates (routes/sync_sessions.go)
├── Does NOT handle presentation CRUD (PocketBase built-in API)
├── Does NOT handle frontend routing (separate Astro application)
└── Does NOT manage viewer connections (PocketBase realtime)
```

**Integration Points**
```
routes/presentations.go
├── Uses: sync session creation patterns from routes/sync_sessions.go
├── Calls: PocketBase collection APIs (app.FindRecordById, app.Save)
├── Returns: JSON responses compatible with frontend expectations
└── Registers: Route handlers via routes/routes.go integration
```

### Data Flow Architecture

```
Request Flow:
1. HTTP Request → PocketBase Router
2. Router → Authentication Middleware
3. Middleware → routes/presentations.go handlers
4. Handlers → PocketBase Collection APIs
5. Collection APIs → SQLite Database
6. Database → Response JSON
7. Response JSON → Client

Dependencies:
presentations.go → PocketBase Core → SQLite
presentations.go → Authentication Context
presentations.go → sync_sessions patterns (code patterns only)
```

## Module Structure

### `routes/presentations.go` Internal Organization

```go
package routes

// === IMPORTS ===
import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/pocketbase/pocketbase/core"
)

// === TYPES AND CONSTANTS ===
type StartLiveResponse struct { /* ... */ }
type StopLiveResponse struct { /* ... */ }
type StatusResponse struct { /* ... */ }

const TokenLength = 32 // Match sync_sessions.go

// === PUBLIC INTERFACE ===
func registerPresentations(e *core.ServeEvent)

// === ROUTE HANDLERS ===
func handleStartLive(e *core.RequestEvent, app core.App) error
func handleStopLive(e *core.RequestEvent, app core.App) error
func handleGetStatus(e *core.RequestEvent, app core.App) error

// === HELPER FUNCTIONS ===
func getAuthenticatedUser(e *core.RequestEvent) (*core.Record, error)
func checkPresentationOwnership(presentation, user *core.Record) error
func generateAdminToken() (string, error)
func progressToStep(progress float64, stepCount int) int
func stepToProgress(stepIndex int, stepCount int) float64

// === UTILITIES ===
func buildStartLiveResponse(session, presentation *core.Record, token string) StartLiveResponse
func buildStatusResponse(presentation, session *core.Record) StatusResponse
```

### Code Organization Principles

1. **Public Interface First**: Export only `registerPresentations`
2. **Handler Separation**: One handler per endpoint for clarity
3. **Helper Functions**: Extract reusable logic (auth, ownership, calculations)
4. **Type Safety**: Structured request/response types
5. **Error Consistency**: Standardized error response format

## Interface Definitions

### Public Interface

```go
// routes/presentations.go exports
func registerPresentations(e *core.ServeEvent)
```

**Registration Contract:**
- Called from `routes/routes.go` during server startup
- Registers three HTTP endpoints with PocketBase router
- No return value, panics on critical errors during registration

### HTTP API Interface

**POST /api/presentations/:id/live**
```
Request:
- Method: POST
- Path: /api/presentations/{presentation_id}/live
- Headers: Authorization (PocketBase auth token)
- Body: None

Response (201 Created):
{
  "session_id": "string",
  "admin_url": "string",
  "viewer_url": "string",
  "step_count": number,
  "step_labels": ["string"]
}

Error Responses:
- 401: {"error": "Authentication required"}
- 403: {"error": "Not authorized"}
- 404: {"error": "Presentation not found"}
- 409: {"error": "Presentation already live"}
- 500: {"error": "Internal server error"}
```

**POST /api/presentations/:id/stop**
```
Request:
- Method: POST
- Path: /api/presentations/{presentation_id}/stop
- Headers: Authorization (PocketBase auth token)
- Body: None

Response (200 OK):
{
  "message": "Presentation stopped"
}

Error Responses:
- 401: {"error": "Authentication required"}
- 403: {"error": "Not authorized"}
- 404: {"error": "Presentation not found"}
- 409: {"error": "Presentation not live"}
- 500: {"error": "Internal server error"}
```

**GET /api/presentations/:id/status**
```
Request:
- Method: GET
- Path: /api/presentations/{presentation_id}/status
- Headers: None required
- Body: None

Response (200 OK):
{
  "id": "string",
  "name": "string",
  "step_count": number,
  "step_labels": ["string"],
  "is_live": boolean,
  "progress": number | null,
  "current_step": number | null
}

Error Responses:
- 404: {"error": "Presentation not found"}
- 500: {"error": "Internal server error"}
```

### Internal Interfaces

**Authentication Interface**
```go
func getAuthenticatedUser(e *core.RequestEvent) (*core.Record, error)
func checkPresentationOwnership(presentation, user *core.Record) error
```

**Utility Interface**
```go
func generateAdminToken() (string, error)
func progressToStep(progress float64, stepCount int) int
func stepToProgress(stepIndex int, stepCount int) float64
```

**Response Building Interface**
```go
func buildStartLiveResponse(session, presentation *core.Record, token string) StartLiveResponse
func buildStatusResponse(presentation, session *core.Record) StatusResponse
```

## Component Boundaries

### Internal Components

**Authentication Component**
```
Purpose: User authentication and authorization
Scope: Request-level auth token validation, ownership checks
Dependencies: PocketBase auth context
Interface: getAuthenticatedUser, checkPresentationOwnership
```

**Session Management Component**
```
Purpose: Sync session creation and linking
Scope: Create sessions, generate tokens, update presentations
Dependencies: PocketBase collections API, crypto/rand
Interface: generateAdminToken, session creation logic
```

**Status Computation Component**
```
Purpose: Live status and progress calculations
Scope: Progress-to-step mapping, response enhancement
Dependencies: Presentation and session data
Interface: progressToStep, buildStatusResponse
```

**Route Handling Component**
```
Purpose: HTTP request/response processing
Scope: Route registration, request parsing, response formatting
Dependencies: All internal components
Interface: handleStartLive, handleStopLive, handleGetStatus
```

### External Dependencies

**PocketBase Core**
```
Used: Record finding, creation, updates
Interface: app.FindRecordById, app.Save, core.NewRecord
Boundary: Database operations only, no business logic
```

**HTTP Layer**
```
Used: Request parsing, response generation
Interface: e.Request, e.JSON, e.Router
Boundary: Protocol handling only, no business logic
```

**Authentication Layer**
```
Used: User context extraction
Interface: e.Request.Context()
Boundary: Auth token validation only
```

## File Dependencies

### Direct Dependencies

```
routes/presentations.go
├── github.com/pocketbase/pocketbase/core (PocketBase APIs)
├── encoding/json (JSON marshaling)
├── net/http (HTTP status codes)
├── crypto/rand (admin token generation)
├── crypto/subtle (constant-time comparison)
├── encoding/hex (token encoding)
└── fmt (error formatting)
```

### Code Pattern Dependencies

```
routes/presentations.go uses patterns from:
├── routes/sync_sessions.go (token generation, validation, response structure)
├── routes/health.go (route registration pattern)
└── migrations/collections.go (PocketBase record operations)
```

### Import Strategy

**Standard Library**
- crypto/rand, crypto/subtle: Security-critical token operations
- encoding/json, encoding/hex: Data serialization
- fmt: Error message formatting
- net/http: HTTP status constants

**Third-Party**
- github.com/pocketbase/pocketbase/core: Core PocketBase functionality
- No additional third-party dependencies required

## Testing Structure

### Unit Test Boundaries

**`routes/presentations_test.go`** (future)
```
Test Categories:
├── Authentication helpers (getAuthenticatedUser, checkPresentationOwnership)
├── Utility functions (progressToStep, stepToProgress)
├── Response builders (buildStartLiveResponse, buildStatusResponse)
└── Error handling (various error conditions)

Mock Dependencies:
├── PocketBase record objects (*core.Record)
├── HTTP request/response (*core.RequestEvent)
└── Authentication context
```

### Integration Test Points

**Route Registration**
- Verify routes registered correctly in router
- Check endpoint availability and method matching

**Database Integration**
- Real PocketBase app with test collections
- Actual record creation/update operations
- Transaction behavior verification

**Authentication Integration**
- PocketBase auth middleware integration
- Context extraction and validation
- Ownership checks with real user records

## Deployment Structure

### Build Integration

**No Build Changes Required**
- Go module system handles new file automatically
- No additional build steps or dependencies
- Standard `go build` includes new routes package

### Configuration Dependencies

**Environment Variables**
- No new environment variables required
- URL generation uses request context and constants
- Admin token generation self-contained

**Database Schema**
- Collections already exist from T-007-01
- No migrations required
- Uses existing presentation/sync_session relationship

## Change Ordering

### Implementation Order

1. **Create `routes/presentations.go`** with basic structure
2. **Implement authentication helpers** (self-contained)
3. **Implement utility functions** (progress calculations)
4. **Implement GET /status endpoint** (read-only, lowest risk)
5. **Implement POST /stop endpoint** (simple state change)
6. **Implement POST /live endpoint** (most complex, creates resources)
7. **Update `routes/routes.go`** to register new routes
8. **Test integration** with existing system

### Verification Steps

**After Each Handler Implementation**
- Verify endpoint responds with correct status codes
- Test error conditions return proper JSON
- Validate integration with PocketBase collections

**After Full Implementation**
- Test complete workflow: start live → check status → stop live
- Verify frontend compatibility with new APIs
- Check no regression in existing functionality

This structure provides clear boundaries, minimal coupling, and follows established patterns while enabling the new presentation lifecycle functionality.