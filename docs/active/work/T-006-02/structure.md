# T-006-02: sync-api-routes - Structure

## Overview

This document defines the file-level changes, architecture, and component boundaries for implementing the sync API routes. The structure replaces existing sync session endpoints with new ones that match the exact acceptance criteria specifications.

## File Changes Summary

### Modified Files

**1. `routes/sync_sessions.go`**
- **Change Type**: Major modification
- **Scope**: Replace existing route definitions and request/response handling
- **Impact**: Breaking change to existing API endpoints

**2. `routes/routes.go`**
- **Change Type**: No changes required
- **Rationale**: Registration function `registerSyncSessions(e)` remains the same

### No New Files Required

The existing file structure is sufficient. Creating a separate file would fragment related functionality unnecessarily.

### No Deleted Files

All changes are contained within the existing routing structure.

## Module Boundaries

### Core Components

**1. Route Registration (`registerSyncSessions`)**
```go
func registerSyncSessions(e *core.ServeEvent)
```
- **Responsibility**: Register HTTP route handlers with PocketBase router
- **Input**: PocketBase ServeEvent for route binding
- **Output**: Registered route handlers
- **Dependencies**: PocketBase core router

**2. Session Creation Handler (`handleCreateSession`)**
```go
func handleCreateSession(e *core.RequestEvent, app core.App) error
```
- **Responsibility**: Create new sync session with admin token
- **Input**: HTTP request (POST with empty body)
- **Output**: JSON response with session details
- **Dependencies**: PocketBase app, crypto/rand, sync_sessions collection

**3. Progress Update Handler (`handleUpdateProgress`)**
```go
func handleUpdateProgress(e *core.RequestEvent, app core.App) error
```
- **Responsibility**: Update session progress with token validation
- **Input**: HTTP request (POST with JSON body and query token)
- **Output**: JSON status response
- **Dependencies**: PocketBase app, sync_sessions collection

### Utility Functions

**1. Token Generation (`generateAdminToken`)**
```go
func generateAdminToken() (string, error)
```
- **Responsibility**: Generate cryptographically secure random token
- **Input**: None
- **Output**: 64-character hex string
- **Dependencies**: crypto/rand, encoding/hex
- **Reuse**: Keep existing implementation (proven secure)

**2. Token Validation (`validateToken`)**
```go
func validateToken(provided, stored string) bool
```
- **Responsibility**: Constant-time token comparison
- **Input**: Two token strings
- **Output**: Boolean match result
- **Dependencies**: crypto/subtle
- **Reuse**: Keep existing implementation (security-critical)

**3. Progress Validation (`validateProgress`)**
```go
func validateProgress(progress float64) error
```
- **Responsibility**: Validate progress value range
- **Input**: Float64 progress value
- **Output**: Error if invalid, nil if valid
- **Dependencies**: fmt
- **Reuse**: Keep existing implementation

## Public Interfaces

### HTTP API Surface

**Endpoint 1: Create Session**
- **Method**: POST
- **Path**: `/api/sync/create`
- **Request Body**: Empty
- **Response**: JSON with session_id, admin_url, viewer_url
- **Status Codes**: 201 (success), 500 (server error)

**Endpoint 2: Update Progress**
- **Method**: POST
- **Path**: `/api/sync/:id/progress`
- **Query Parameters**: `token=<admin_token>`
- **Request Body**: JSON with progress field
- **Response**: JSON status message
- **Status Codes**: 200 (success), 400 (bad request), 403 (forbidden), 404 (not found), 500 (server error)

### Data Models

**Create Session Response**
```go
type CreateSessionResponse struct {
    SessionID string `json:"session_id"`
    AdminURL  string `json:"admin_url"`
    ViewerURL string `json:"viewer_url"`
}
```

**Update Progress Request**
```go
type UpdateProgressRequest struct {
    Progress float64 `json:"progress"`
}
```

**Error Response**
```go
type ErrorResponse struct {
    Error string `json:"error"`
}
```

**Success Response**
```go
type SuccessResponse struct {
    Message string `json:"message"`
}
```

## Internal Organization

### Package Structure
```
routes/
├── sync_sessions.go    # Modified - all sync session functionality
├── routes.go          # Unchanged - route registration
├── health.go          # Unchanged
├── stats.go           # Unchanged
└── ...                # Other route files unchanged
```

### Function Organization within sync_sessions.go

**1. Type Definitions (Top)**
- Constants (TokenLength, ProgressMin, ProgressMax)
- Struct definitions for request/response models

**2. Public Registration Function**
- `registerSyncSessions()` - Entry point for route registration

**3. HTTP Handlers (Request Flow Order)**
- `handleCreateSession()` - POST /api/sync/create
- `handleUpdateProgress()` - POST /api/sync/:id/progress

**4. Utility Functions (Alphabetical)**
- `generateAdminToken()` - Token generation
- `validateProgress()` - Progress validation
- `validateToken()` - Token comparison

## Implementation Dependencies

### External Dependencies
- `github.com/pocketbase/pocketbase/core` - PocketBase core functionality
- `crypto/rand` - Secure random number generation
- `crypto/subtle` - Constant-time comparison
- `encoding/hex` - Token encoding
- `encoding/json` - JSON request/response handling
- `fmt` - Error formatting
- `net/http` - HTTP status codes

### Internal Dependencies
- `sync_sessions` collection - Must exist (created by migrations)
- PocketBase realtime system - For progress broadcast (automatic)

## Change Ordering Requirements

### Critical Dependencies
1. **Database schema must exist** - sync_sessions collection (already satisfied by T-006-01)
2. **No frontend coordination required** - API changes are breaking but self-contained

### Implementation Sequence
1. **Update route definitions** - Change URL paths first
2. **Modify request handling** - Update token extraction method
3. **Update response formatting** - Match acceptance criteria
4. **Update error handling** - Ensure correct status codes
5. **Test integration** - Verify with existing collection

### Rollback Considerations
- Single file modification enables quick rollback
- No schema changes reduce rollback complexity
- Breaking API changes require frontend coordination for rollback

## Integration Points

### PocketBase Integration
- **Collection Access**: Uses `app.FindCollectionByNameOrId("sync_sessions")`
- **Record Operations**: Uses `app.Save(record)` and `app.FindRecordById()`
- **Route Registration**: Uses `e.Router.POST()` for endpoint binding

### Frontend Integration
- **URL Format Changes**: Admin URL format changes from `/admin/{id}?token={token}` to `/sync/{id}/control?token={token}`
- **Authentication Method**: Token moves from header to query parameter
- **Response Format**: session_id field name changes from "id"

### Realtime Integration
- **Progress Broadcasting**: Automatic through PocketBase realtime when record.Set("progress") called
- **Subscription Targets**: Frontend subscriptions to sync_sessions collection remain valid
- **Event Format**: Progress update events maintain same data structure

## Security Boundaries

### Authentication Perimeter
- **Public Endpoints**: Session creation (no auth required)
- **Protected Endpoints**: Progress updates (admin token required)
- **Token Scope**: Admin tokens are session-scoped, not global

### Data Access Patterns
- **Read Access**: Public read access to sync_sessions collection
- **Write Access**: Restricted to custom API endpoints only
- **Token Storage**: Admin tokens stored in database, validated server-side

### Security Validations
- **Input Validation**: Progress values constrained to 0.0-1.0 range
- **Token Validation**: Constant-time comparison prevents timing attacks
- **SQL Injection**: PocketBase ORM handles parameter sanitization
- **CORS**: Handled by PocketBase default CORS policy

## Performance Characteristics

### Memory Usage
- **Token Generation**: 32 bytes random + 64 bytes hex encoding per session
- **Request Processing**: Single database query per request (lookup or save)
- **Response Size**: ~150 bytes for create response, ~50 bytes for update response

### Database Impact
- **Create Session**: 1 INSERT operation
- **Update Progress**: 1 SELECT + 1 UPDATE operation
- **Indexes**: Primary key lookups only (session ID)

### Concurrent Access
- **Session Creation**: Safe concurrent creation (unique IDs)
- **Progress Updates**: Safe concurrent updates to same session (last write wins)
- **Realtime Broadcasting**: PocketBase handles concurrent subscriptions

This structure maintains the existing architectural patterns while implementing the exact requirements specified in the acceptance criteria. The single-file modification approach minimizes risk while ensuring complete compliance with the ticket requirements.