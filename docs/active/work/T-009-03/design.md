# T-009-03 Design: Fix PathValue Routing

## Problem Analysis

The current codebase has a routing compatibility issue between PocketBase v0.36.5 and Go 1.22's `PathValue()` method. Specifically:

### Current State
- Routes are defined using Echo-style syntax: `/api/presentations/:id/status`
- Handler code attempts to extract parameters using Go 1.22's `PathValue("id")`
- PocketBase v0.36.5 uses Echo v5 router, not Go 1.22's native router
- `PathValue()` returns empty strings, causing silent failures

### Evidence Found
1. **presentations.go:185**: `presentationID := e.Request.PathValue("id")`
2. **sync_sessions.go:88**: `sessionID := e.Request.PathValue("id")`
3. **Route definitions**: All use `:id` syntax (Echo style)
4. **Existing workaround**: `links_view.go` manually parses URL path due to this same issue

### Root Cause
PocketBase v0.36.5 still uses Echo v5 router internally, despite having some Go 1.22-compatible APIs. The `PathValue()` method only works with Go's native `net/http` router, not with Echo's routing system.

## Design Options

### Option 1: Use Echo's Parameter Extraction
**Approach**: Replace `PathValue()` calls with Echo's parameter extraction methods.

**Implementation**:
```go
// Instead of:
presentationID := e.Request.PathValue("id")

// Use Echo's method via RequestEvent:
presentationID := e.Request.URL.Query().Get("id") // Wrong for path params

// Or extract from RequestEvent context (requires investigation)
```

**Pros**:
- Minimal code changes
- Works with current PocketBase version

**Cons**:
- Requires understanding Echo's internal parameter extraction
- May not be stable API
- Need to investigate RequestEvent structure

### Option 2: Manual Path Parsing (Current Pattern)
**Approach**: Parse URL path manually as done in `links_view.go`.

**Implementation**:
```go
func extractIDFromPath(path string, segment string) string {
    parts := strings.Split(path, "/")
    for i, part := range parts {
        if part == segment && i+1 < len(parts) {
            return parts[i+1]
        }
    }
    return ""
}

// Usage:
presentationID := extractIDFromPath(e.Request.URL.Path, "presentations")
```

**Pros**:
- Already proven to work in `links_view.go`
- No dependency on router internals
- Consistent across all endpoints
- Simple and transparent

**Cons**:
- More verbose
- Manual URL parsing logic
- Less elegant than native parameter extraction

### Option 3: Update Route Syntax to Go 1.22 Format
**Approach**: Change route definitions from `:id` to `{id}` syntax.

**Implementation**:
```go
// Change from:
e.Router.GET("/api/presentations/:id/status", handler)

// To:
e.Router.GET("/api/presentations/{id}/status", handler)
```

**Pros**:
- Would work with Go 1.22 native router
- More modern syntax

**Cons**:
- **CRITICAL**: PocketBase v0.36.5 likely doesn't support `{id}` syntax yet
- High risk of breaking all routes
- May require PocketBase version upgrade

### Option 4: Use PocketBase's Route Parameter API
**Approach**: Investigate if PocketBase exposes Echo's parameters through RequestEvent.

**Implementation**:
```go
// Investigate if available:
presentationID := e.PathParam("id")  // Unknown if this exists
// or
presentationID := e.Param("id")      // Unknown if this exists
```

**Pros**:
- Clean API if available
- Router-agnostic

**Cons**:
- Unknown if this API exists in current PocketBase version
- Requires API investigation

## Recommended Solution

**Choose Option 2: Manual Path Parsing**

### Rationale

1. **Proven Working**: `links_view.go` demonstrates this approach works
2. **Low Risk**: No dependency on router internals or version changes
3. **Consistent**: Can standardize across all affected endpoints
4. **Maintainable**: Simple, readable code that's easy to debug
5. **Future-Proof**: Will continue working regardless of PocketBase router changes

### Implementation Strategy

1. Create a shared utility function for path parameter extraction
2. Replace `PathValue()` calls in both `presentations.go` and `sync_sessions.go`
3. Add simple test to verify parameter extraction works correctly
4. Consider refactoring `links_view.go` to use the shared utility

### Code Structure

```go
// In routes/utils.go (new file)
func extractPathParam(path, segment string) string {
    // Robust path parsing implementation
}

// Updated handlers will use:
presentationID := extractPathParam(e.Request.URL.Path, "presentations")
sessionID := extractPathParam(e.Request.URL.Path, "sync")
```

### Testing Strategy

Add verification that route parameters are correctly extracted by:
1. Creating a simple test that calls the utility function with known paths
2. Testing edge cases (empty paths, missing segments, malformed URLs)
3. Verifying that existing endpoints continue to work

## Alternative Consideration

If time permits, investigate Option 4 to see if PocketBase exposes proper parameter extraction. However, Option 2 should be implemented first as the safe fallback.

## Risk Assessment

- **Low Risk**: Manual parsing approach is simple and battle-tested
- **Medium Impact**: Affects multiple critical endpoints, but fix is straightforward
- **High Priority**: Silent failures make this critical to fix immediately

## Success Criteria

1. All presentation endpoints work: `/api/presentations/:id/live`, `/stop`, `/status`
2. All sync endpoints work: `/api/sync/:id/progress`
3. Route parameters are correctly extracted in all cases
4. Simple test validates the parameter extraction
5. No regressions in existing functionality