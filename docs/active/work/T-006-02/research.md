# Research: sync-api-routes (T-006-02)

## Context

This ticket requires implementing custom Go routes for creating and controlling sync sessions. The create endpoint generates a session with a random admin token and returns both admin and viewer URLs. The progress endpoint lets the admin push new progress values, which PocketBase's realtime then broadcasts to subscribers.

## Current State

### Existing Infrastructure

The project has the following structure in place:
- PocketBase v0.36.5 application with migrations and custom routes
- `sync_sessions` collection created by T-006-01 with `progress` (float 0-1) and `admin_token` fields
- Routes package exists with registration pattern established by T-001-02
- `routes/sync_sessions.go` file already exists but doesn't match ticket requirements

### Current Implementation Analysis

The existing `routes/sync_sessions.go` file (169 lines) implements:
- `POST /api/sync-sessions` for session creation
- `PATCH /api/sync-sessions/:id` for progress updates
- Admin token in `X-Admin-Token` header (not query param)
- Uses crypto/rand for token generation (32 bytes, hex encoded)
- Constant-time token comparison for security
- Returns different response structure than specified

Key differences from requirements:
1. **Endpoint paths**: Uses `/api/sync-sessions` instead of `/api/sync/create` and `/api/sync/:id/progress`
2. **HTTP method**: Uses PATCH instead of POST for progress updates
3. **Token location**: Uses header instead of query parameter
4. **Response structure**: Different field names (`id` vs `session_id`)
5. **URL format**: Returns `/admin/` and `/viewer/` paths instead of `/sync/` paths

### Route Registration Pattern

From `routes/routes.go` and `routes/health.go`:
- Routes are registered via `OnServe` hook
- Individual route files provide registration functions
- Routes can be registered directly on `e.Router` (Echo v5)
- Middleware binding pattern available for intercepting existing routes

### Collection Schema

From `migrations/collections.go`:
- `sync_sessions` collection fields:
  - `progress`: NumberField, min 0, max 1, not required (allows default)
  - `admin_token`: TextField, required, min/max 64 chars (hex-encoded 32 bytes)
- Public read access via API rules (ListRule and ViewRule set to "")
- No create/update/delete rules (handled by custom routes)

### Security Implementation

Current code shows good security practices:
- Cryptographically secure random token generation using crypto/rand
- Constant-time comparison for token validation (crypto/subtle)
- 32-byte tokens (256 bits of entropy)
- Hex encoding for safe transmission

## Technical Requirements Analysis

### Endpoint Specifications

**POST /api/sync/create**:
- Creates new sync_session record
- Generates random admin token server-side
- Returns JSON with specific structure:
  - `session_id` (not `id`)
  - `admin_url` with format `/sync/{id}/control?token={token}`
  - `viewer_url` with format `/sync/{id}`

**POST /api/sync/:id/progress**:
- Accepts JSON body with `progress` field (0-1 float)
- Token validation via query parameter `?token={admin_token}`
- Updates session record's progress field
- Returns appropriate HTTP status codes:
  - 200: Success with updated session
  - 403: Invalid token
  - 404: Session not found
  - 400: Progress out of range

### URL Structure Requirements

The ticket specifies specific URL formats:
- Admin URL: `/sync/{session_id}/control?token={admin_token}`
- Viewer URL: `/sync/{session_id}`

These appear to be frontend routes, not backend API endpoints. The backend returns these as convenience strings for the frontend to use.

### Token Requirements

- Cryptographically random string
- Example given: "32 hex chars" (which would be 16 bytes)
- Current implementation uses 32 bytes (64 hex chars) - more secure
- No authentication required beyond token for admin operations
- Viewer endpoints are fully public

## Implementation Considerations

### Route Path Conflicts

Need to ensure new routes don't conflict with existing:
- Current: `/api/sync-sessions` and `/api/sync-sessions/:id`
- Required: `/api/sync/create` and `/api/sync/:id/progress`
- These are distinct paths, so both can coexist if needed

### Token in Query vs Header

Requirements specify query parameter for token:
- Pros: Simpler for basic HTTP clients, visible in URLs
- Cons: Can leak in logs, referrer headers, browser history
- Current header approach is more secure but not what's specified

### Response Format Mapping

Need to transform response structure:
- Current uses `ID` field from PocketBase record
- Required uses `session_id` field
- URLs need different path structure

### Method Choice

Requirements specify POST for progress updates:
- POST typically for creating resources
- PATCH more semantically correct for partial updates
- Requirements are explicit about POST

## Dependencies and Constraints

### Completed Dependencies
- T-001-02: Custom health route pattern established
- T-006-01: sync_sessions collection created

### File Structure
```
routes/
  routes.go           # Main registration function
  health.go          # Health endpoint example
  sync_sessions.go   # Current implementation (needs modification)
```

### PocketBase Integration
- Uses Echo v5 router under the hood
- Request/response handled via `core.RequestEvent`
- Database operations via `app.FindRecordById` and `app.Save`
- Collection lookup via `app.FindCollectionByNameOrId`

## Risks and Decisions

### Backwards Compatibility
- Existing implementation might be in use
- Could maintain both endpoint sets temporarily
- Or replace entirely per requirements

### Security Trade-offs
- Query parameter tokens are less secure than headers
- Requirements are explicit about query parameters
- Should follow requirements despite security implications

### Error Response Format
- Current implementation returns JSON error objects
- Requirements don't specify error format
- Should maintain consistency with existing patterns

### Frontend URL Generation
- Admin and viewer URLs are frontend routes
- Backend generates these as convenience
- Frontend must handle these URL patterns appropriately

## Next Steps

The Design phase should address:
1. Whether to replace or supplement existing endpoints
2. Maintaining security with query parameter tokens
3. Error response format standards
4. Testing strategy for both endpoints
5. Migration path if replacing existing code
6. Documentation for API consumers