# Research: presentation-api-routes (T-007-02)

## Context

This ticket implements custom Go routes for presentation lifecycle management, specifically focusing on live session workflows that PocketBase's built-in CRUD API cannot handle. The routes create a bridge between presentations and sync sessions established in earlier tickets.

## Current State

### Project Architecture

The project is a Go-based PocketBase v0.36.5 application with:
- PocketBase as core framework providing built-in REST APIs and realtime subscriptions
- Programmatic migration system in `migrations/collections.go`
- Custom route registration in `routes/` package for specialized business logic
- SQLite database backend via PocketBase
- Astro frontend with SSR enabled for presentation viewer interfaces

### Existing Collections

From dependencies T-006-02 and T-007-01, three collections exist:

1. **sync_sessions** (T-006-01):
   - `progress`: float, 0-1 range, represents presentation progress
   - `admin_token`: 64-char hex string for admin access control
   - Created/updated via custom routes only (no direct PocketBase API access)
   - Public read access for real-time progress broadcasting

2. **presentations** (T-007-01):
   - `name`: text, human-readable presentation name
   - `step_count`: integer ≥ 1, total number of presentation steps
   - `step_labels`: JSON array, optional labels for each step
   - `active_session`: relation to sync_sessions, null when not live
   - `created_by`: relation to users for ownership
   - Uses PocketBase built-in API for CRUD operations

3. **links**, **tags**: Unrelated to presentation functionality

### Existing Custom Route Patterns

From `routes/sync_sessions.go` (T-006-02):

**Architecture Pattern:**
```go
// Registration in routes/routes.go
func registerSyncSessions(e *core.ServeEvent) {
    e.Router.POST("/api/sync/create", handleCreateSession)
    e.Router.POST("/api/sync/:id/progress", handleUpdateProgress)
}

// Handler pattern
func handleCreateSession(e *core.RequestEvent, app core.App) error {
    // 1. Generate secure tokens
    // 2. Find collection
    // 3. Create record with core.NewRecord()
    // 4. Set fields with record.Set()
    // 5. Save with app.Save(record)
    // 6. Return JSON response
}
```

**Security Pattern:**
- Admin tokens: 32 random bytes, hex-encoded to 64 chars
- Constant-time token comparison using `subtle.ConstantTimeCompare`
- Query parameter authentication: `?token=<admin_token>`
- Validation with structured error responses

**Response Pattern:**
```go
type CreateSessionResponse struct {
    SessionID string `json:"session_id"`
    AdminURL  string `json:"admin_url"`
    ViewerURL string `json:"viewer_url"`
}
```

### Frontend Integration Points

From `frontend/src/pages/watch/[id].astro`:

**Expected APIs:**
- Presentations fetched via: `GET /api/collections/presentations/records/{id}` (PocketBase built-in)
- Presentation state determined by `active_session` field (null = waiting, non-null = live)
- Real-time polling checks presentation status for live transitions

**URL Patterns:**
- Viewer URL: `/watch/{presentation_id}` → frontend route
- Admin control URLs: Referenced in sync session creation but implementation unclear

## Requirements Analysis

### New Route Specifications

Per acceptance criteria, three routes needed:

1. **`POST /api/presentations/:id/live`** - Start live session:
   - Input: presentation ID from URL path
   - Process: Create sync_session, link to presentation
   - Output: session info with URLs and step metadata

2. **`POST /api/presentations/:id/stop`** - End live session:
   - Input: presentation ID from URL path
   - Process: Clear presentation.active_session
   - Output: confirmation response

3. **`GET /api/presentations/:id/status`** - Get presentation with live state:
   - Input: presentation ID from URL path
   - Process: Fetch presentation + compute live status
   - Output: presentation data + computed live fields

### Step-to-Progress Mapping Logic

Per requirements: `step_progress = step_index / (step_count - 1)` for step_count > 1, or `0.0` for single-step.

**Examples:**
- 1 step: progress 0.0 (single point)
- 2 steps: [0.0, 1.0] (start and end)
- 5 steps: [0.0, 0.25, 0.5, 0.75, 1.0]

**Formula Implementation:**
```go
func stepToProgress(stepIndex, stepCount int) float64 {
    if stepCount <= 1 {
        return 0.0
    }
    return float64(stepIndex) / float64(stepCount - 1)
}
```

### Authentication Strategy

**Key Question:** No explicit authentication requirements in T-007-02, but:
- Presentations have ownership via `created_by` field
- PocketBase API rules require authentication for presentation updates
- Sync session creation currently has no auth requirements

**Analysis Options:**
1. No authentication (anyone can start any presentation live)
2. Ownership-based (only presentation owner can control)
3. Admin token system (like sync sessions)

Given PocketBase API rules on presentations, ownership-based makes most sense.

## Integration Requirements

### Database Operations Needed

**For /live endpoint:**
1. Fetch presentation by ID (validate exists)
2. Create new sync_session record (progress=0, generate admin_token)
3. Update presentation.active_session = new session ID
4. Return structured response with URLs and metadata

**For /stop endpoint:**
1. Fetch presentation by ID (validate exists)
2. Clear presentation.active_session = null
3. Return confirmation

**For /status endpoint:**
1. Fetch presentation by ID (validate exists)
2. Fetch related active_session if exists
3. Compute derived fields (is_live, current_step, viewer_count)
4. Return enhanced presentation data

### Response Format Analysis

From acceptance criteria for `/live`:
```json
{
  "session_id": "abc123",
  "admin_url": "/sync/abc123/control?token=xyz",
  "viewer_url": "/sync/abc123",
  "step_count": 5,
  "step_labels": ["Intro", "Main", "Demo", "Q&A", "Close"]
}
```

**URL Format Questions:**
- Admin URL path `/sync/{id}/control` - does this route exist?
- Viewer URL path `/sync/{id}` - does this route exist?
- Or should these be `/watch/{presentation_id}` based on frontend?

From acceptance criteria for `/status`:
```json
{
  // Standard presentation fields
  "id": "pres123",
  "name": "My Presentation",
  "step_count": 5,
  "step_labels": [...],
  // Computed live fields
  "is_live": true,
  "progress": 0.5,
  "current_step": 2,
  "viewer_count": 42  // if available
}
```

## File Structure Impact

Current route organization:
```
routes/
  routes.go           # Central registration
  health.go          # Health check override
  sync_sessions.go   # Session management (172 lines)
  stats.go           # Statistics endpoints
  links_search.go    # Link search functionality
  links_search_simple.go  # Simplified link search
```

**New file needed:** `routes/presentations.go`
- Follow existing patterns from sync_sessions.go
- Register via registerPresentations() in routes.go
- Estimated ~150-200 lines based on similar complexity

## Technical Challenges

### 1. URL Generation Strategy

Sync session responses reference admin/viewer URLs, but unclear which routes serve these:
- Current sync endpoints: `/api/sync/{id}/*`
- Frontend viewer route: `/watch/{presentation_id}`
- Admin control interface: Unknown location

**Resolution needed:** Clarify URL patterns for session control vs presentation viewing

### 2. Viewer Count Implementation

Acceptance criteria mentions "viewer_count (if available from PocketBase, else omit)"

PocketBase doesn't track active connections per collection. Options:
- Omit always (simplest)
- Track via custom counter (complex)
- Use PocketBase realtime subscription metrics (if accessible)

### 3. Authentication Integration

PocketBase auth context available via `e.Request.Context()`, but pattern not used in existing routes.

For ownership checks:
```go
// Get auth record from request context
auth, ok := e.Request.Context().Value("auth").(*core.Record)
if !ok || auth == nil {
    return e.JSON(401, map[string]string{"error": "Authentication required"})
}

// Check ownership
if presentation.GetString("created_by") != auth.Id {
    return e.JSON(403, map[string]string{"error": "Not authorized"})
}
```

### 4. Error Handling Consistency

Existing patterns use structured error responses:
```go
return e.JSON(http.StatusNotFound, map[string]string{
    "error": "Presentation not found",
})
```

Need consistent error messages across all three endpoints.

## Dependencies and Integration Points

### Required Imports
```go
import (
    "encoding/json"
    "net/http"
    "github.com/pocketbase/pocketbase/core"
    // Possibly crypto/rand for admin tokens (already in sync_sessions.go)
)
```

### Collection Access Patterns
Following sync_sessions.go model:
```go
// Find collection
collection, err := app.FindCollectionByNameOrId("presentations")
if err != nil {
    return e.JSON(500, map[string]string{"error": "Collection not found"})
}

// Find record
record, err := app.FindRecordById("presentations", presentationID)
if err != nil {
    return e.JSON(404, map[string]string{"error": "Presentation not found"})
}
```

### Cross-Collection Operations
Live endpoint needs both presentations and sync_sessions collections:
1. Create sync_session record
2. Update presentation record with new session ID
3. Both operations should be transactional (or handle partial failure)

## Risk Analysis

### 1. Race Conditions
Multiple concurrent requests to `/live` endpoint could create multiple sessions for same presentation.

Mitigation: Check existing active_session before creating new one.

### 2. Orphaned Sessions
If `/live` creates session but fails to update presentation, sync_session becomes orphaned.

Mitigation: Create session first, then update presentation. Cleanup job could remove orphaned sessions.

### 3. Session Lifecycle Management
No automatic cleanup of ended sessions - they persist forever.

Consideration: Should `/stop` delete session or just clear relationship?

### 4. Frontend URL Assumptions
Generated URLs assume specific frontend route structure that may change.

Mitigation: Make URL templates configurable or use environment variables.

## Next Phase Guidance

The Design phase should resolve:

1. **Authentication strategy** - ownership-based vs open vs token-based
2. **URL generation patterns** - admin/viewer URL formats and routes
3. **Viewer count implementation** - omit vs track vs estimate
4. **Error handling standards** - consistent message formats
5. **Transaction handling** - ensure data consistency across collections
6. **Session cleanup strategy** - automatic vs manual cleanup
7. **Frontend integration** - URL patterns and expected response formats