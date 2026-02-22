# Research Phase - T-010-02: Extract Shared Go Utilities

## Current State Analysis

### Token Generation Duplication

#### `presentations.go` Implementation
- **Constant**: `TokenLength = 32` (line 14)
- **Function**: `generateAdminToken()` (lines 84-90)
  - Creates 32-byte random token
  - Uses `crypto/rand.Read`
  - Returns hex-encoded string
  - Used in `handleStartLive()` for creating presentation sessions

#### `sync_sessions.go` Implementation
- **Constant**: `syncTokenLength = 32` (line 15)
- **Function**: `generateSyncAdminToken()` (lines 150-156)
  - Identical logic to `generateAdminToken()`
  - Creates 32-byte random token
  - Uses `crypto/rand.Read`
  - Returns hex-encoded string
  - Used in `handleCreateSession()` for sync sessions

#### Token Validation
- `sync_sessions.go` has `validateSyncToken()` (lines 167-172)
  - Uses `crypto/subtle.ConstantTimeCompare` for timing-safe comparison
  - Checks length equality first
- `presentations.go` imports `crypto/subtle` but doesn't use validation function
  - Direct validation appears embedded in handler logic

### HTTP Status Code Usage

#### `stats.go` Issues
- Line 76: `e.JSON(500, ...)`
- Line 82: `e.JSON(500, ...)`
- Line 88: `e.JSON(500, ...)`
- Line 94: `e.JSON(500, ...)`
- Line 100: `e.JSON(500, ...)`
- Line 111: `e.JSON(200, ...)`

#### `links_view.go` Issues
- Line 25: `e.JSON(400, ...)`
- Line 34: `e.JSON(500, ...)`
- Line 42: `e.JSON(404, ...)`
- Line 50: `e.JSON(500, ...)`
- Line 67: `e.JSON(200, ...)`

#### Correct Usage Examples
- `presentations.go` uses `http.StatusBadRequest`, `http.StatusNotFound`, etc.
- `sync_sessions.go` uses `http.StatusInternalServerError`, `http.StatusCreated`, etc.
- `links_search.go` uses `http.StatusBadRequest`, `http.StatusInternalServerError`, `http.StatusOK`

### Error Handling Gaps

#### `links_search.go`
- Line 88-93: Comment says "Log error but don't fail the request"
- Currently just sets empty `tagMap` without logging
- No actual logging implementation present
- `fetchTagsForLinks()` has proper `rows.Err()` check at line 330

#### `links_view.go`
- Line 40: `rowsAffected, _ := result.RowsAffected()`
- Error from `RowsAffected()` is swallowed with underscore
- No error logging for this ignored error

### Security Issue: Admin Token Exposure

#### `sync_sessions.go` handleUpdateProgress()
- Lines 140-146: Response includes full session data
- Line 143: `"admin_token": record.GetString("admin_token")`
- Admin token included in progress update response
- Exposes sensitive admin credential to client

### Magic Numbers

#### `stats.go`
- Line 21: `LIMIT 5` in `sqlMostViewed` query
- Line 34: `LIMIT 5` in `sqlTopTags` query
- Both queries return top 5 results
- Hard-coded in SQL string constants

## Package Structure

### Current Files
```
routes/
├── health.go                 # Health check endpoint
├── links_search.go           # Search functionality
├── links_search_test.go      # Search tests
├── links_view.go             # View count tracking
├── links_view_test.go        # View tests
├── path_utils.go             # Path parameter extraction utilities
├── path_utils_test.go        # Path utils tests
├── presentations.go          # Presentation lifecycle management
├── routes.go                 # Route registration orchestrator
├── routes_test.go            # Route tests
├── static.go                 # Static file serving
├── stats.go                  # Statistics endpoints
├── stats_test.go             # Stats tests
├── sync_sessions.go          # Sync session management
└── sync_sessions_test.go     # Sync session tests
```

### Import Patterns
- All files use `github.com/pocketbase/pocketbase/core`
- Database operations use `github.com/pocketbase/dbx`
- Token generation uses `crypto/rand` and `encoding/hex`
- Timing-safe comparisons use `crypto/subtle`
- HTTP status codes should use `net/http` constants

## Testing Infrastructure

### Existing Test Files
- `links_search_test.go`: Comprehensive search tests
- `links_view_test.go`: View count increment tests
- `path_utils_test.go`: Path extraction tests
- `routes_test.go`: Route registration tests
- `stats_test.go`: Statistics endpoint tests
- `sync_sessions_test.go`: Sync session tests

### Test Patterns
- Tests use standard Go testing package
- Mock PocketBase app instances for testing
- Test data includes sample links, tags, presentations
- HTTP request/response testing via PocketBase's test utilities

## Dependencies and Constraints

### PocketBase Integration
- All routes registered via `core.ServeEvent`
- Database access through `app.DB()`
- Record operations via `app.FindRecordById()`, `app.Save()`
- Collection access via `app.FindCollectionByNameOrId()`

### Shared Utilities Already Present
- `extractPathParam()` in `path_utils.go` - used across multiple files
- `extractLastPathParam()` in `path_utils.go` - alternative extraction pattern

### Database Schema Dependencies
- `links` table with `view_count`, `tags` fields
- `tags` table with `name`, `slug` fields
- `presentations` table with `active_session` field
- `sync_sessions` table with `progress`, `admin_token` fields

## Summary

The codebase has clear patterns of duplication in token management, inconsistent HTTP status code usage, missing error handling, and a security issue with token exposure. The existing structure shows good separation of concerns with dedicated files for different route groups. Test coverage appears comprehensive. The shared utilities extraction will follow established patterns from `path_utils.go`.