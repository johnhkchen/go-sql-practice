# T-009-03 Structure: Fix PathValue Routing

## Overview

This document defines the file-level changes required to fix the PathValue routing compatibility issue. The solution implements manual path parsing to replace the non-functional `PathValue()` calls in PocketBase v0.36.5.

## Files to Modify

### 1. `routes/path_utils.go` (NEW FILE)
**Purpose**: Centralized utility functions for path parameter extraction

**Interface**:
```go
package routes

// extractPathParam extracts a parameter value from a URL path
// Example: extractPathParam("/api/presentations/abc123/status", "presentations") returns "abc123"
func extractPathParam(path, segment string) string

// extractLastPathParam extracts the last segment from a path (for simple cases)
// Example: extractLastPathParam("/api/sync/abc123/progress", "progress") returns "abc123"
func extractLastPathParam(path, beforeSegment string) string
```

**Dependencies**:
- `strings` (standard library)

**Module Boundaries**:
- Internal utility functions for the routes package only
- No external dependencies beyond standard library

### 2. `routes/presentations.go` (MODIFY)
**Changes Required**:

**Function: `handleGetStatus()` (line 184-215)**
- Replace: `presentationID := e.Request.PathValue("id")`
- With: `presentationID := extractPathParam(e.Request.URL.Path, "presentations")`
- Add import for the new utility

**Function: `handleStopLive()` (line 217-273)**
- Replace: `presentationID := e.Request.PathValue("id")`
- With: `presentationID := extractPathParam(e.Request.URL.Path, "presentations")`

**Function: `handleStartLive()` (line 275-359)**
- Replace: `presentationID := e.Request.PathValue("id")`
- With: `presentationID := extractPathParam(e.Request.URL.Path, "presentations")`

**Import Changes**:
- No new external imports required (utility function is in same package)

### 3. `routes/sync_sessions.go` (MODIFY)
**Changes Required**:

**Function: `handleUpdateProgress()` (line 86-147)**
- Replace: `sessionID := e.Request.PathValue("id")`
- With: `sessionID := extractPathParam(e.Request.URL.Path, "sync")`
- Add import for the new utility

**Import Changes**:
- No new external imports required (utility function is in same package)

### 4. `routes/links_view.go` (REFACTOR - OPTIONAL)
**Changes for Consistency**:
- Replace manual path parsing (lines 22-32) with call to shared utility
- Maintain identical behavior but use consistent approach

**Before**:
```go
path := e.Request.URL.Path
parts := strings.Split(path, "/")
var linkId string
for i, part := range parts {
    if part == "links" && i+1 < len(parts) {
        linkId = parts[i+1]
        break
    }
}
```

**After**:
```go
linkId := extractPathParam(e.Request.URL.Path, "links")
```

### 5. `test_pathvalue.go` (DELETE)
**Purpose**: Remove temporary test file created during investigation
**Justification**: This was created for investigation only and is not needed for production

## File Dependencies

```
routes/path_utils.go
    ↓ (imports)
routes/presentations.go
routes/sync_sessions.go
routes/links_view.go (optional refactor)
```

## Module Architecture

### Package: `routes`
**Public Interface** (no changes):
- All existing route registration functions remain unchanged
- Handler function signatures remain unchanged
- Response formats remain unchanged

**Internal Implementation**:
- New utility module `path_utils.go` provides parameter extraction
- Handler implementations updated to use utilities instead of `PathValue()`
- No changes to route definitions (still use `:id` syntax)

### Function Boundaries

**Public Functions** (unchanged):
- `registerPresentations()`
- `registerSyncSessions()`
- `registerLinksView()`

**Private Functions** (modified implementation):
- `handleGetStatus()` - updated parameter extraction
- `handleStopLive()` - updated parameter extraction
- `handleStartLive()` - updated parameter extraction
- `handleUpdateProgress()` - updated parameter extraction
- `handleLinksView()` - optionally refactored for consistency

**New Private Functions**:
- `extractPathParam()` - core utility function
- `extractLastPathParam()` - convenience utility for simple cases

## Implementation Ordering

The changes must be made in this specific order to avoid compilation errors:

### Phase 1: Create Utilities
1. Create `routes/path_utils.go` with utility functions
2. Verify compilation succeeds with new file

### Phase 2: Update Presentations
1. Modify `routes/presentations.go` to use utilities
2. Test that presentations endpoints work correctly

### Phase 3: Update Sync Sessions
1. Modify `routes/sync_sessions.go` to use utilities
2. Test that sync endpoints work correctly

### Phase 4: Optional Refactoring
1. Update `routes/links_view.go` for consistency
2. Delete temporary `test_pathvalue.go` file

### Phase 5: Verification
1. Run comprehensive tests on all affected endpoints
2. Verify no regressions in existing functionality

## Error Handling Strategy

### Utility Functions
- Return empty string for invalid/missing parameters
- Match existing error handling patterns in handlers
- No panics or exceptions - graceful degradation

### Handler Updates
- Maintain existing error response formats
- Preserve HTTP status codes for invalid parameters
- No changes to error message content

## Testing Strategy

### Unit Tests (routes/path_utils_test.go - NEW)
Test the utility functions with:
- Valid paths: `/api/presentations/abc123/status`
- Invalid paths: `/api/presentations//status`
- Missing segments: `/api/presentations/status`
- Empty paths: `""`
- Malformed URLs: `/api//presentations/abc123`

### Integration Testing
- Verify all affected endpoints continue to work
- Test with realistic presentation and session IDs
- Confirm error responses for invalid IDs remain unchanged

## Backward Compatibility

**Route Definitions**: No changes - still use `:id` syntax
**Response Formats**: No changes to JSON structures
**Error Messages**: Preserve existing error text
**HTTP Status Codes**: No changes to response codes
**API Contracts**: Full backward compatibility maintained

## Risk Mitigation

### Compilation Safety
- New utility file has no external dependencies
- All modified functions maintain same signatures
- No changes to public interfaces

### Runtime Safety
- Utility functions handle edge cases gracefully
- Empty parameter detection preserves existing error paths
- No new failure modes introduced

### Testing Safety
- Existing tests should continue to pass
- New utility tests validate edge cases
- Integration tests verify end-to-end functionality

## Success Criteria Mapping

1. **All presentation endpoints work**: Modified parameter extraction in all three handlers
2. **All sync endpoints work**: Modified parameter extraction in `handleUpdateProgress()`
3. **Route parameters correctly extracted**: New utility functions handle all cases
4. **Simple test validates extraction**: New unit tests for utility functions
5. **No regressions**: Maintain all existing behavior and response formats