# T-006-02: sync-api-routes - Design

## Overview

This ticket implements custom Go API routes for sync session management, building on the existing sync_sessions collection and routing patterns. The implementation must support session creation with admin tokens and progress updates with token-based authorization.

## Current State Analysis

### Existing Architecture

The codebase follows a clean PocketBase-based architecture:

1. **Main Entry Point**: `main.go` registers migrations and routes, then starts PocketBase
2. **Route Management**: Centralized in `routes/routes.go` with individual route files
3. **Collection Management**: Automated migrations in `migrations/collections.go` create required schema
4. **Route Pattern**: Each feature has its own file (e.g., `health.go`, `sync_sessions.go`)

### Sync Sessions Infrastructure

The sync_sessions collection already exists with:
- `progress`: NumberField (0.0-1.0, float)
- `admin_token`: TextField (64 chars, required for hex-encoded 32 bytes)
- Public read access (ListRule and ViewRule set to "")
- No create/update/delete rules (handled by custom routes)

### Existing Implementation Gap

Current `sync_sessions.go` implements:
- `POST /api/sync-sessions` (create)
- `PATCH /api/sync-sessions/:id` (update progress)

But the ticket requires:
- `POST /api/sync/create`
- `POST /api/sync/:id/progress`

The URL structure and some response formats don't match the acceptance criteria.

## Design Options Evaluated

### Option 1: Extend Existing sync_sessions.go

**Approach**: Add the required endpoints to the existing file alongside current endpoints.

**Pros**:
- Reuses existing logic for token generation and validation
- Maintains consistency with current patterns
- Minimal code duplication

**Cons**:
- Creates two different URL schemes for the same functionality
- Potential confusion with multiple endpoints serving similar purposes
- Response format mismatches (current vs. required)

### Option 2: Replace Existing Implementation

**Approach**: Modify existing routes to match ticket requirements exactly.

**Pros**:
- Single source of truth for sync session APIs
- Exact match to acceptance criteria
- Clean, consistent interface

**Cons**:
- May break existing frontend dependencies
- Requires understanding current usage patterns
- Higher risk if other components depend on existing URLs

### Option 3: New Implementation with Coexistence

**Approach**: Create new routes matching the ticket spec while leaving existing routes unchanged.

**Pros**:
- Zero risk to existing functionality
- Exact compliance with acceptance criteria
- Clear separation of concerns

**Cons**:
- Code duplication between similar endpoints
- Multiple ways to achieve the same result
- Potential maintenance burden

## Chosen Approach: Option 2 - Replace Existing Implementation

**Rationale**: Based on the acceptance criteria's specificity and the fact that this appears to be refining the API design rather than adding parallel functionality, replacing the existing implementation is the most appropriate choice.

**Key Factors**:
1. The ticket has specific URL structure requirements that don't match current implementation
2. Response format specifications suggest this is the intended final API design
3. No indication that existing endpoints need to be preserved
4. Dependencies listed (T-001-02, T-006-01) suggest this builds on established foundations

## Technical Design

### Route Structure

Replace existing routes with:
```go
// POST /api/sync/create - Create new session
e.Router.POST("/api/sync/create", func(ev *core.RequestEvent) error {
    return handleCreateSession(ev, e.App)
})

// POST /api/sync/:id/progress - Update progress with query param token
e.Router.POST("/api/sync/:id/progress", func(ev *core.RequestEvent) error {
    return handleUpdateProgress(ev, e.App)
})
```

### Response Format Changes

**Create Session Response**:
```json
{
    "session_id": "abc123",
    "admin_url": "/sync/abc123/control?token=<admin_token>",
    "viewer_url": "/sync/abc123"
}
```

Key changes from existing:
- Use "session_id" instead of "id"
- Admin URL format: `/sync/{id}/control?token={token}` (not `/admin/{id}?token={token}`)
- Viewer URL format: `/sync/{id}` (not `/viewer/{id}`)
- Remove progress field from create response

### Authentication Changes

**Token Validation**:
- Change from header-based (`X-Admin-Token`) to query parameter (`?token=<admin_token>`)
- Maintain constant-time comparison security
- Keep same cryptographic token generation (32 hex chars)

### Error Handling

Maintain existing error response patterns:
- 403 for invalid token
- 404 for session not found
- 400 for invalid progress values
- 500 for server errors

### Code Reuse Strategy

**Preserve Security Functions**:
- `generateAdminToken()` - Keep as-is, proven crypto implementation
- `validateToken()` - Keep constant-time comparison logic
- `validateProgress()` - Keep range validation (0.0-1.0)

**Modify Request Handling**:
- Update token extraction from query params vs headers
- Adjust response formatting to match specifications
- Change URL path handling for new structure

## Integration Considerations

### Frontend Dependencies

The admin and viewer URLs must match frontend routing expectations:
- `/sync/{id}/control?token={token}` for admin interface
- `/sync/{id}` for viewer interface

### PocketBase Realtime

Progress updates will continue to trigger PocketBase's realtime subscriptions since we're updating the same collection records. No changes needed to realtime broadcast behavior.

### Database Schema

No schema changes required - existing sync_sessions collection fields match requirements:
- `progress` field handles float values 0.0-1.0
- `admin_token` field stores 64-char hex strings
- Collection rules allow public read access for viewer functionality

## Testing Strategy

### Unit Tests
- Token generation randomness and format
- Progress validation boundary conditions
- Constant-time token comparison

### Integration Tests
- Session creation flow with proper response format
- Progress update flow with query param authentication
- Error scenarios (invalid tokens, missing sessions, out-of-range progress)

### Security Tests
- Token brute force resistance
- Side-channel attack resistance in token comparison
- SQL injection resistance in ID handling

## Risk Assessment

**Low Risk**:
- Building on established PocketBase patterns
- Reusing proven security components
- Well-defined acceptance criteria

**Medium Risk**:
- Frontend compatibility depends on URL format changes
- Query parameter vs header authentication change may affect clients

**Mitigation**:
- Verify no production dependencies exist on current endpoints
- Test authentication flow thoroughly
- Document API changes clearly

## Performance Considerations

**Database Impact**: Minimal - same collection operations as existing implementation

**Authentication Overhead**: Query parameter parsing is marginally faster than header parsing

**Memory Usage**: Same token generation and validation logic maintains current memory profile

## Deployment Notes

Since this modifies existing API endpoints, deployment should:
1. Ensure frontend code expects new URL formats
2. Test admin token flow with query parameters
3. Verify realtime subscriptions continue working
4. Confirm response format compatibility

The implementation maintains backward compatibility at the data level (same collection, same fields) while changing the HTTP API surface.