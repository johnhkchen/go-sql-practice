# Design Phase - T-010-02: Extract Shared Go Utilities

## Design Options

### Option 1: Create `routes/tokens.go` with Shared Utilities

**Approach**: Create a new file dedicated to token management that consolidates all token-related functionality.

**Structure**:
```go
package routes

const TokenLength = 32

func GenerateToken() (string, error) { ... }
func ValidateToken(provided, stored string) bool { ... }
```

**Pros**:
- Clear separation of concerns
- Single source of truth for token logic
- Follows existing pattern (`path_utils.go` for path utilities)
- Easy to test in isolation
- Minimal changes to existing code

**Cons**:
- Adds a new file to the routes package
- Requires updating imports in two files

### Option 2: Create `internal/tokens` Package

**Approach**: Create a separate internal package for token management.

**Structure**:
```
internal/
└── tokens/
    ├── tokens.go
    └── tokens_test.go
```

**Pros**:
- Better package-level separation
- Could be reused outside routes package if needed
- More formal API boundary

**Cons**:
- Over-engineering for just 2-3 functions
- Introduces new package dependency
- More complex import paths
- Project doesn't have `internal/` pattern established

### Option 3: Keep Duplication, Add Comments

**Approach**: Leave the duplicate functions but add comments explaining the duplication.

**Pros**:
- No code changes required
- Each file remains self-contained

**Cons**:
- Violates DRY principle
- Maintenance burden (must update in multiple places)
- Risk of divergence over time
- Doesn't address the core issue

## Selected Design: Option 1

### Rationale

Option 1 is the optimal choice because:

1. **Follows Established Patterns**: The codebase already has `path_utils.go` for shared path extraction utilities. Creating `tokens.go` follows this same pattern.

2. **Minimal Disruption**: Only requires changes to two files that consume the token functions, with straightforward function renames.

3. **Right-Sized Solution**: Not over-engineered (like Option 2) but properly addresses the duplication (unlike Option 3).

4. **Test-Friendly**: Can create `tokens_test.go` following the existing test pattern.

## Detailed Design Decisions

### Token Functions API

```go
// TokenLength defines the byte length of tokens before encoding
const TokenLength = 32

// GenerateToken generates a cryptographically secure random token
func GenerateToken() (string, error)

// ValidateToken performs constant-time comparison of tokens
func ValidateToken(provided, stored string) bool
```

**Design Choices**:
- Export functions with capitalized names (Go convention for public APIs)
- Keep same token length (32 bytes) for compatibility
- Maintain hex encoding for tokens
- Use generic names (`GenerateToken` not `GenerateAdminToken`) since usage context varies

### HTTP Status Code Constants

**Approach**: Replace raw integers with `http.Status*` constants directly in place.

**Example**:
```go
// Before
e.JSON(500, map[string]string{"error": "..."})

// After
e.JSON(http.StatusInternalServerError, map[string]string{"error": "..."})
```

**Design Choice**: Don't create custom constants or wrappers. Use standard library constants directly for clarity and convention.

### Error Handling Improvements

#### `links_search.go` Logging
**Approach**: Add proper logging using PocketBase's logger.

```go
if err != nil {
    app.Logger().Error("Failed to fetch tags", "error", err)
    tagMap = make(map[string][]string)
}
```

#### `links_view.go` Error Check
**Approach**: Check error and log if present, but continue operation.

```go
rowsAffected, err := result.RowsAffected()
if err != nil {
    app.Logger().Error("Failed to get rows affected", "error", err)
    // Continue with rowsAffected = 0 assumption
}
```

### Security Fix: Remove Admin Token from Response

**Approach**: Build custom response map excluding sensitive fields.

```go
// Remove "admin_token" from the response map
return e.JSON(http.StatusOK, map[string]interface{}{
    "id":       record.Id,
    "progress": record.GetFloat("progress"),
    "created":  record.GetDateTime("created").Time(),
    "updated":  record.GetDateTime("updated").Time(),
    // admin_token explicitly excluded
})
```

### Magic Number Extraction

**Approach**: Define constant near the SQL queries.

```go
const (
    statsTopN = 5  // Number of top items to return

    sqlMostViewed = `...LIMIT 5`
    sqlTopTags = `...LIMIT 5`
)
```

**Design Choice**: Keep the LIMIT in the SQL as-is since these are const strings. The `statsTopN` constant documents the intent. Alternative would require fmt.Sprintf which adds complexity for minimal benefit.

## Migration Strategy

1. Create `tokens.go` with shared utilities
2. Update `presentations.go` to use shared functions
3. Update `sync_sessions.go` to use shared functions
4. Fix HTTP status codes in `stats.go`
5. Fix HTTP status codes in `links_view.go`
6. Add error logging in `links_search.go`
7. Fix error handling in `links_view.go`
8. Remove admin_token from response in `sync_sessions.go`
9. Add statsTopN constant documentation in `stats.go`

## Testing Strategy

- Create `tokens_test.go` for token utilities
- Verify existing tests still pass after changes
- No new tests needed for status code changes (cosmetic)
- Existing tests should catch any regression from admin_token removal

## Rejected Alternatives

- **Custom HTTP status wrapper functions**: Unnecessary abstraction
- **Logging framework change**: Out of scope, use existing PocketBase logger
- **Token interface abstraction**: Over-engineering for current needs
- **Database migration for token storage**: Not required by ticket scope