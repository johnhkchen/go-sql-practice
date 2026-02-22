# Structure Phase - T-010-02: Extract Shared Go Utilities

## File Operations Overview

### CREATE: `routes/tokens.go`
New file for shared token utilities.

### MODIFY: `routes/presentations.go`
Update to use shared token functions.

### MODIFY: `routes/sync_sessions.go`
Update to use shared token functions and fix admin_token leak.

### MODIFY: `routes/stats.go`
Replace raw HTTP status codes and document magic number.

### MODIFY: `routes/links_view.go`
Replace raw HTTP status codes and fix error handling.

### MODIFY: `routes/links_search.go`
Add proper error logging.

## Detailed File Structure

### CREATE: `routes/tokens.go`

```go
package routes

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/hex"
)

// TokenLength defines the byte length of tokens before encoding
const TokenLength = 32

// GenerateToken generates a cryptographically secure random token
func GenerateToken() (string, error) {
    bytes := make([]byte, TokenLength)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

// ValidateToken performs constant-time comparison of tokens
func ValidateToken(provided, stored string) bool {
    if len(provided) != len(stored) {
        return false
    }
    return subtle.ConstantTimeCompare([]byte(provided), []byte(stored)) == 1
}
```

### MODIFY: `routes/presentations.go`

**Remove**:
- Line 14: `const TokenLength = 32`
- Lines 84-90: `generateAdminToken()` function

**Change**:
- Line 297: `adminToken, err := generateAdminToken()`
  → `adminToken, err := GenerateToken()`

**Imports remain unchanged** (already has crypto/subtle imported)

### MODIFY: `routes/sync_sessions.go`

**Remove**:
- Line 15: `const syncTokenLength = 32`
- Lines 150-156: `generateSyncAdminToken()` function
- Lines 167-172: `validateSyncToken()` function

**Change**:
- Line 48: `token, err := generateSyncAdminToken()`
  → `token, err := GenerateToken()`
- Line 123: `if !validateSyncToken(adminToken, storedToken)`
  → `if !ValidateToken(adminToken, storedToken)`

**Change response at lines 140-146**:
```go
return e.JSON(http.StatusOK, map[string]interface{}{
    "id":       record.Id,
    "progress": record.GetFloat("progress"),
    "created":  record.GetDateTime("created").Time(),
    "updated":  record.GetDateTime("updated").Time(),
    // admin_token removed from response
})
```

### MODIFY: `routes/stats.go`

**Add after line 9** (after sqlTotalLinks const):
```go
// statsTopN defines the number of top items returned by stats queries
const statsTopN = 5
```

**Add import**:
```go
import (
    "net/http"  // Add this
    "github.com/pocketbase/dbx"
    "github.com/pocketbase/pocketbase/core"
)
```

**Change HTTP status codes**:
- Line 76: `500` → `http.StatusInternalServerError`
- Line 82: `500` → `http.StatusInternalServerError`
- Line 88: `500` → `http.StatusInternalServerError`
- Line 94: `500` → `http.StatusInternalServerError`
- Line 100: `500` → `http.StatusInternalServerError`
- Line 111: `200` → `http.StatusOK`

### MODIFY: `routes/links_view.go`

**Add import**:
```go
import (
    "net/http"  // Add this
    "strings"
    "time"
    "github.com/pocketbase/dbx"
    "github.com/pocketbase/pocketbase/core"
)
```

**Change HTTP status codes**:
- Line 25: `400` → `http.StatusBadRequest`
- Line 34: `500` → `http.StatusInternalServerError`
- Line 42: `404` → `http.StatusNotFound`
- Line 50: `500` → `http.StatusInternalServerError`
- Line 67: `200` → `http.StatusOK`

**Fix error handling at line 40**:
```go
// Before:
rowsAffected, _ := result.RowsAffected()

// After:
rowsAffected, err := result.RowsAffected()
if err != nil {
    app.Logger().Error("Failed to get rows affected", "error", err)
    rowsAffected = 0
}
```

### MODIFY: `routes/links_search.go`

**Change lines 88-93**:
```go
tagMap, err := fetchTagsForLinks(app, linkIDs)
if err != nil {
    // Log error but don't fail the request
    app.Logger().Error("Failed to fetch tags", "error", err)
    tagMap = make(map[string][]string)
}
```

## Module Boundaries

### Token Module (`tokens.go`)
**Public API**:
- `TokenLength` constant
- `GenerateToken()` function
- `ValidateToken()` function

**Internal Dependencies**:
- Uses only standard library packages
- No PocketBase dependencies

### Consumer Modules
**presentations.go**:
- Consumes: `GenerateToken()`
- Does not need `ValidateToken()` currently

**sync_sessions.go**:
- Consumes: `GenerateToken()`, `ValidateToken()`

## Import Graph

```
routes/tokens.go
├── crypto/rand
├── crypto/subtle
└── encoding/hex

routes/presentations.go
├── routes/tokens (NEW)
├── crypto/subtle (already imported, keep for future use)
├── net/http (already imported)
└── [other existing imports]

routes/sync_sessions.go
├── routes/tokens (NEW)
├── net/http (already imported)
└── [other existing imports minus removed crypto imports]

routes/stats.go
├── net/http (NEW)
└── [other existing imports]

routes/links_view.go
├── net/http (NEW)
└── [other existing imports]

routes/links_search.go
└── [no new imports, uses app.Logger() already available]
```

## Testing Impact

- Existing tests should continue to pass
- Token generation behavior unchanged (same algorithm, length, encoding)
- Token validation behavior unchanged (same comparison logic)
- HTTP responses unchanged (except admin_token removal from one endpoint)
- No new test files required in this ticket
- Consider adding `tokens_test.go` as enhancement (not in ticket scope)