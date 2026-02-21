# Design: presentation-api-routes (T-007-02)

## Overview

This design document evaluates implementation options for three presentation lifecycle API routes: starting live sessions, stopping them, and getting enhanced status information. The routes bridge presentations with sync sessions while maintaining security and consistency.

## Decision Points from Research

### 1. Authentication Strategy

**Options Evaluated:**

**A) No Authentication (Open Access)**
- Anyone can start/stop any presentation
- Pros: Simple implementation, matches current sync session pattern
- Cons: No access control, potential abuse, inconsistent with PocketBase presentation rules
- Verdict: **Rejected** - Inconsistent with existing presentation ownership model

**B) Ownership-Based Authentication**
- Only presentation owner (created_by) can control live sessions
- Pros: Consistent with PocketBase API rules, clear ownership model
- Cons: Requires auth context handling, excludes shared presentations
- Verdict: **Selected** - Best fit for existing architecture

**C) Admin Token System**
- Generate per-presentation admin tokens like sync sessions
- Pros: Flexible access control, could enable sharing
- Cons: Additional complexity, storage requirements, unclear lifecycle
- Verdict: **Rejected** - Overengineering for current requirements

**Decision: Ownership-Based Authentication**
- Extract auth context from PocketBase request
- Verify `presentation.created_by` matches authenticated user ID
- Return 401 for unauthenticated, 403 for unauthorized
- Consistent with existing PocketBase collection rules

### 2. URL Generation Strategy

**Options Evaluated:**

**A) Sync-Based URLs (Current Pattern)**
- Admin URL: `/sync/{session_id}/control?token={admin_token}`
- Viewer URL: `/sync/{session_id}`
- Pros: Matches existing sync session API responses
- Cons: Frontend uses `/watch/{presentation_id}`, creates URL inconsistency

**B) Presentation-Based URLs**
- Admin URL: `/presentations/{presentation_id}/control`
- Viewer URL: `/watch/{presentation_id}` (matches frontend)
- Pros: Consistent with frontend routing, presentation-centric
- Cons: Doesn't directly reference session, admin control unclear

**C) Hybrid Approach**
- Admin URL: `/sync/{session_id}/control?token={admin_token}` (session control)
- Viewer URL: `/watch/{presentation_id}` (presentation viewing)
- Pros: Each URL serves its specific purpose optimally
- Cons: Mixed URL patterns

**Decision: Hybrid Approach**
- Viewer URL: `/watch/{presentation_id}` - matches frontend expectation
- Admin URL: `/sync/{session_id}/control?token={admin_token}` - enables session control
- Rationale: Viewers care about presentations, admins control sessions

### 3. Viewer Count Implementation

**Options Evaluated:**

**A) Always Omit**
- Don't include viewer_count in responses
- Pros: Simple, no infrastructure needed
- Cons: Misses acceptance criteria "if available"

**B) Track via Counter Fields**
- Add viewer_count field to sync_sessions, update on access
- Pros: Persistent counts, simple queries
- Cons: Race conditions, stale data, not real-time

**C) PocketBase Realtime Metrics**
- Use PocketBase internal subscription tracking
- Pros: Real-time accuracy
- Cons: Not exposed in public API, implementation dependent

**Decision: Always Omit (Phase 1)**
- Acceptance criteria says "if available from PocketBase, else omit"
- PocketBase doesn't expose connection metrics
- Keep implementation simple, add later if needed
- Document this decision for future enhancement

### 4. Session Lifecycle Management

**Options Evaluated:**

**A) Persistent Sessions**
- Created sessions persist until manually deleted
- /stop clears presentation.active_session only
- Pros: History preservation, audit trail
- Cons: Storage growth, orphaned sessions

**B) Auto-Cleanup on Stop**
- /stop deletes sync session record entirely
- Pros: Clean storage, no orphans
- Cons: Lose session history/admin tokens

**C) Soft Delete Pattern**
- Add deleted/ended timestamp, filter in queries
- Pros: History + cleanup, flexible recovery
- Cons: More complex queries, additional fields

**Decision: Persistent Sessions**
- /stop only clears presentation.active_session relationship
- Preserves sync session records and admin tokens
- Enables potential reconnection or history features
- Storage concerns addressable later via cleanup jobs

### 5. Error Handling Standards

**Consistency Requirements:**
- Match existing sync session route error format
- Use structured JSON responses
- Appropriate HTTP status codes
- Clear, actionable error messages

**Standard Error Format:**
```json
{
  "error": "Description of what went wrong"
}
```

**Status Code Mapping:**
- 400: Invalid input (malformed progress, etc.)
- 401: Authentication required
- 403: Not authorized (wrong owner)
- 404: Resource not found (presentation, session)
- 409: Conflict (already live, not live)
- 500: Server errors (database, system)

### 6. Transaction Handling

**Problem:** Live endpoint needs atomic operations:
1. Create sync_session
2. Update presentation.active_session

**Options Evaluated:**

**A) Sequential Operations**
- Create session, then update presentation
- Handle partial failure in application logic
- Pros: Simple implementation
- Cons: Potential orphaned sessions

**B) Database Transactions**
- Wrap operations in PocketBase transaction
- Pros: True atomicity
- Cons: More complex, requires transaction API understanding

**C) Optimistic Approach**
- Create session first, update presentation
- Accept orphaned sessions as acceptable trade-off
- Clean up via background job if needed
- Pros: Simple, performant
- Cons: Eventually consistent

**Decision: Sequential with Validation**
- Check presentation not already live before creating session
- Create sync_session record first
- Update presentation.active_session second
- Return error if either step fails
- Accept rare orphaned session risk (addressable later)

## Implementation Design

### Route Structure

**File:** `routes/presentations.go`
```go
package routes

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/pocketbase/pocketbase/core"
    // Reuse token utilities from sync_sessions
)

func registerPresentations(e *core.ServeEvent) {
    e.Router.POST("/api/presentations/:id/live", handleStartLive)
    e.Router.POST("/api/presentations/:id/stop", handleStopLive)
    e.Router.GET("/api/presentations/:id/status", handleGetStatus)
}
```

### Data Structures

```go
// Start live response
type StartLiveResponse struct {
    SessionID   string   `json:"session_id"`
    AdminURL    string   `json:"admin_url"`
    ViewerURL   string   `json:"viewer_url"`
    StepCount   int      `json:"step_count"`
    StepLabels  []string `json:"step_labels"`
}

// Stop live response
type StopLiveResponse struct {
    Message string `json:"message"`
}

// Status response (extends presentation data)
type StatusResponse struct {
    // Include all original presentation fields
    ID         string      `json:"id"`
    Name       string      `json:"name"`
    StepCount  int         `json:"step_count"`
    StepLabels []string    `json:"step_labels"`
    // Add computed fields
    IsLive      bool    `json:"is_live"`
    Progress    *float64 `json:"progress,omitempty"`
    CurrentStep *int     `json:"current_step,omitempty"`
    // ViewerCount omitted per decision
}
```

### Step-to-Progress Mapping

```go
func progressToStep(progress float64, stepCount int) int {
    if stepCount <= 1 {
        return 0
    }
    // Find closest step boundary
    stepProgress := progress * float64(stepCount - 1)
    return int(stepProgress + 0.5) // Round to nearest
}

func stepToProgress(stepIndex int, stepCount int) float64 {
    if stepCount <= 1 {
        return 0.0
    }
    return float64(stepIndex) / float64(stepCount - 1)
}
```

### Authentication Helper

```go
func getAuthenticatedUser(e *core.RequestEvent) (*core.Record, error) {
    // Extract auth from PocketBase context
    auth, ok := e.Request.Context().Value("auth").(*core.Record)
    if !ok || auth == nil {
        return nil, fmt.Errorf("authentication required")
    }
    return auth, nil
}

func checkPresentationOwnership(presentation, user *core.Record) error {
    if presentation.GetString("created_by") != user.Id {
        return fmt.Errorf("not authorized")
    }
    return nil
}
```

### Core Handler Logic

**Start Live (/live):**
1. Extract presentation ID from URL
2. Authenticate user and verify ownership
3. Fetch presentation record, validate exists
4. Check presentation not already live (active_session == null)
5. Generate admin token (reuse sync session logic)
6. Create sync_session record (progress=0, admin_token)
7. Update presentation.active_session = session ID
8. Build response with URLs and step metadata
9. Return 201 Created

**Stop Live (/stop):**
1. Extract presentation ID from URL
2. Authenticate user and verify ownership
3. Fetch presentation record, validate exists
4. Check presentation is live (active_session != null)
5. Clear presentation.active_session = null
6. Save presentation record
7. Return 200 OK with confirmation

**Get Status (/status):**
1. Extract presentation ID from URL
2. Fetch presentation record, validate exists
3. If active_session exists, fetch session record for progress
4. Compute is_live, current_step from progress and step_count
5. Build enhanced response with computed fields
6. Return 200 OK (no auth required - read-only)

## Integration Points

### Registration in routes.go
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        registerHealth(e)
        registerSyncSessions(e)
        registerPresentations(e)  // Add this line
        registerStats(e)
        registerLinksSearch(e)
        registerLinksSearchSimple(e)
        return e.Next()
    })
}
```

### Frontend Compatibility
- GET /status endpoint provides enhanced data for presentation pages
- Viewer URLs match existing `/watch/{id}` frontend route
- Admin URLs enable session control (may need new frontend routes)
- Response formats compatible with existing frontend polling logic

### Error Response Consistency
Match existing sync session patterns:
```go
return e.JSON(http.StatusNotFound, map[string]string{
    "error": "Presentation not found",
})
```

## Validation Rules

### Input Validation
- Presentation ID: Valid UUID format from URL parameter
- Authentication: Valid PocketBase auth token in request
- Ownership: authenticated user matches presentation.created_by

### State Validation
- Start Live: presentation.active_session must be null
- Stop Live: presentation.active_session must not be null
- Status: no state requirements (read-only)

### Response Validation
- Step labels array length should match step_count (best effort)
- Progress values stay within 0.0-1.0 range
- URLs follow documented patterns

## Security Considerations

### Authentication Flow
1. PocketBase middleware validates auth token
2. Auth record available in request context
3. Handler extracts auth record
4. Ownership check against presentation.created_by
5. Structured error responses (no token leakage)

### Admin Token Security
- Reuse existing 32-byte random generation from sync sessions
- Hex encoding for URL safety
- Constant-time comparison in sync session routes
- Tokens remain valid for session lifetime

### Access Control
- Presentation CRUD: PocketBase collection rules (owner-only updates)
- Live session control: Custom routes with ownership checks
- Status endpoint: Public read access (matches PocketBase rules)
- Sync session control: Existing admin token validation

## Performance Implications

### Database Operations
- Start Live: 3 operations (read presentation, create session, update presentation)
- Stop Live: 2 operations (read presentation, update presentation)
- Status: 2 operations (read presentation, read session if active)

### Caching Opportunities
- Presentation metadata rarely changes (could cache step_count, labels)
- Live status changes frequently (should not cache)
- Auth checks could benefit from request-scoped caching

### Concurrency Considerations
- Multiple start requests: Check active_session atomically
- Rapid start/stop cycles: Last operation wins
- Status polling: Read-only, no conflicts

## Future Enhancement Hooks

### Viewer Count Tracking
- Add viewer_count field to sync_sessions
- Increment on viewer page access
- Decrement on page leave/timeout
- Include in status response when available

### Session History/Analytics
- Preserved sessions enable usage analytics
- Track session duration, step progression
- Presentation performance metrics

### Batch Operations
- Start multiple presentations live
- Bulk status checks
- Administrative session management

## Risk Mitigation

### Orphaned Sessions
- Risk: Created session but failed presentation update
- Mitigation: Background cleanup job for orphaned sessions
- Detection: sessions not referenced by any presentation.active_session

### Authentication Edge Cases
- Risk: Auth context unavailable or malformed
- Mitigation: Explicit auth validation with clear error messages
- Fallback: Graceful degradation with 401 responses

### Frontend URL Changes
- Risk: Generated URLs don't match frontend routes
- Mitigation: Make URL templates configurable
- Documentation: Clear URL pattern contracts

This design provides a secure, consistent, and maintainable implementation that integrates well with existing patterns while addressing all acceptance criteria requirements.