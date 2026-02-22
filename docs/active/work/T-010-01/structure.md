# Structure: Delete Dead Code and Duplicates

## Files to Delete

### `routes/links_search_simple.go`
**Action**: DELETE entire file
**Lines removed**: 191
**Reason**: Complete duplicate of links_search.go functionality

## Files to Modify

### `routes/routes.go`
**Action**: MODIFY
**Changes**:
1. Remove line 17: `registerLinksSearchSimple(e)`
2. Remove lines 21-25 (FrontendExists check), replace with line 22: `registerStatic(e)`

**Before** (lines 16-25):
```go
registerLinksSearch(e)
registerLinksSearchSimple(e)
registerLinksView(e)

// Register static file serving with availability check
if frontend.FrontendExists() {
    registerStatic(e)
} else {
    e.App.Logger().Warn("Frontend assets not found, static serving disabled")
}
```

**After** (lines 16-19):
```go
registerLinksSearch(e)
registerLinksView(e)

// Register static file serving
registerStatic(e)
```

### `routes/presentations.go`
**Action**: MODIFY
**Changes**:
1. Remove line 7: commented import `// "encoding/json"  // TODO: Re-enable if JSON operations are added`
2. Remove lines 106-111: `stepToProgress()` function
3. Remove lines 173-179: `validateToken()` function
4. Remove line 181: comment `// Placeholder handler functions - to be implemented in subsequent steps`

**Functions to remove**:
```go
// stepToProgress converts a step index to progress value (0-1) using the formula:
// progress = step_index / (step_count - 1) for step_count > 1, else 0.0
func stepToProgress(stepIndex int, stepCount int) float64 {
    if stepCount <= 1 {
        return 0.0
    }
    return float64(stepIndex) / float64(stepCount-1)
}
```

```go
// validateToken performs constant-time token comparison
func validateToken(provided, stored string) bool {
    if len(provided) != len(stored) {
        return false
    }
    return subtle.ConstantTimeCompare([]byte(provided), []byte(stored)) == 1
}
```

### `routes/routes_test.go`
**Action**: MODIFY
**Changes**:
1. Remove lines 17-22: `TestData` struct definition
2. Remove lines 157-168: `assertErrorResponse()` function

**Structures/functions to remove**:
```go
// TestData holds references to created test data
type TestData struct {
    LinkIDs  []string
    TagIDs   []string
    TagSlugs []string
}
```

```go
// assertErrorResponse validates error responses consistently
func assertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedMessage string) {
    if resp.StatusCode != expectedStatus {
        t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
    }

    var errResp ErrorResponse
    parseJSONResponse(t, resp, &errResp)

    if errResp.Error != expectedMessage {
        t.Errorf("Expected error message %q, got %q", expectedMessage, errResp.Error)
    }
}
```

### `internal/frontend/embed.go`
**Action**: MODIFY
**Changes**:
1. Remove lines 20-23: `FrontendExists()` function

**Function to remove**:
```go
func FrontendExists() bool {
    _, err := GetFrontendFS()
    return err == nil
}
```

### `routes/static.go`
**Action**: MODIFY
**Changes**:
1. Replace lines 73-79 with single line: `return ev.Next()`

**Before** (lines 71-79):
```go
}

// If not a ReadSeeker, fall back to copying the content
ev.Response.Header().Set("Content-Type", "application/octet-stream")
_, err = ev.Response.Write([]byte("file serving not supported"))
if err != nil {
    return err
}
return nil
```

**After** (lines 71-72):
```go
}
// If not a ReadSeeker, continue to next handler
return ev.Next()
```

## Import Changes

### `routes/presentations.go`
- No import changes needed (commented import is being removed)

### `routes/routes.go`
- No import changes needed (still uses `frontend` package for `GetFrontendFS`)

### All other files
- No import changes needed

## Public API Changes

### Removed Endpoints
- `/api/links/search-simple` - Not in API specification

### Removed Exports
- `frontend.FrontendExists()` - Was public, only used once internally

### Internal Functions Removed
- `validateToken()` - Unused
- `stepToProgress()` - Unused
- `TestData` struct - Unused
- `assertErrorResponse()` - Unused

## Dependency Impact

### Direct Dependencies (remain unchanged)
- `github.com/labstack/echo/v5` - Still used by routes_test.go
- `github.com/pocketbase/dbx` - Still used throughout
- `github.com/pocketbase/pocketbase` - Core framework

### Transitive Dependencies
- Will be cleaned by `go mod tidy` if any become unused

## Test Impact

### Tests to Verify
- `TestMakeRequest_RealExecution` - Must continue to pass
- `TestSetup` - Must continue to pass
- `TestDatabase` - Must continue to pass
- All tests in `routes/` package

### No Test Changes Required
- Removed functions/structs were not used by any tests
- No test assertions depend on removed code

## Build Artifacts

### Expected Outcomes
- Binary size: Slight reduction (~5KB)
- Compilation time: Negligible improvement
- Test execution time: No change

## Line Count Impact

### Total Lines Removed: ~415
- `links_search_simple.go`: 191 lines (entire file)
- `presentations.go`: ~20 lines
- `routes_test.go`: ~24 lines
- `frontend/embed.go`: 4 lines
- `routes.go`: ~8 lines (net reduction)
- `static.go`: ~6 lines (net reduction)

### Net Reduction
- Files: 1 deleted
- Functions: 5 removed
- Exports: 1 removed
- Comments: 2 removed

## Validation Strategy

### After Each File Change
1. Run `go build .`
2. Check for compilation errors
3. Run `go vet ./...`

### After All Changes
1. Run `go test ./...`
2. Run `go mod tidy`
3. Verify no changes to go.mod
4. Manual test of `/api/links/search` endpoint
5. Manual test of static file serving

## File Organization

No changes to directory structure. Only content modifications and one file deletion.