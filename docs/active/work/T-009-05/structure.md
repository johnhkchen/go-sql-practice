# Structure: Fix Integration Tests

## File Modifications

### Modified: `routes/routes_test.go`

**Location**: Lines 109-128 (replacing mock `makeRequest()` function)

**New Structure**:
```go
// makeRequest executes an HTTP request against the test app with real routing
func makeRequest(app *pocketbase.PocketBase, method, url string, body io.Reader) (*http.Response, error) {
    // Create fresh router for this request
    router := echo.New()

    // Create serve event with router
    serveEvent := &core.ServeEvent{
        App:    app,
        Router: router,
    }

    // Register all routes
    routes.Register(app)

    // Trigger OnServe to register routes with router
    app.OnServe().Trigger(serveEvent)

    // Create HTTP request and recorder
    req := httptest.NewRequest(method, url, body)
    rec := httptest.NewRecorder()

    // Execute request through router
    router.ServeHTTP(rec, req)

    // Return the recorded response
    return rec.Result(), nil
}
```

**Import Changes**:
- Uncomment: `"net/http/httptest"`
- Add: `"github.com/labstack/echo/v5"`
- Add: `"github.com/jchen/go-sql-practice/routes"`

### New Tests: `routes/routes_test.go` (append)

**Test: Proof of Real Execution**
```go
// TestMakeRequest_RealExecution proves requests are actually executed
func TestMakeRequest_RealExecution(t *testing.T) {
    app, cleanup := setupTestApp(t)
    defer cleanup()

    // Test 1: Valid endpoint returns success
    resp, err := makeRequest(app, "GET", "/api/health", nil)
    // Assert 200 status

    // Test 2: Invalid endpoint returns 404
    resp, err = makeRequest(app, "GET", "/api/nonexistent", nil)
    // Assert 404 status
}
```

### New Tests: `routes/links_view_test.go` (new file)

**File Structure**:
```go
package routes

import (
    "fmt"
    "net/http"
    "testing"
)

// Test view count increment
func TestLinksView_Success(t *testing.T)
func TestLinksView_NotFound(t *testing.T)
func TestLinksView_InvalidID(t *testing.T)
func TestLinksView_Concurrent(t *testing.T)
```

### New Tests: `routes/stats_test.go` (new file)

**File Structure**:
```go
package routes

import (
    "encoding/json"
    "net/http"
    "testing"
)

// StatsResponse struct for unmarshaling
type StatsTestResponse struct {
    TotalLinks int `json:"total_links"`
    TotalTags  int `json:"total_tags"`
    TotalViews int `json:"total_views"`
    // ... other fields
}

// Test stats endpoint
func TestStats_Complete(t *testing.T)
func TestStats_EmptyDatabase(t *testing.T)
```

### New Tests: `routes/sync_sessions_test.go` (new file)

**File Structure**:
```go
package routes

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
)

// Test sync session creation
func TestSyncCreate_Success(t *testing.T)
func TestSyncCreate_InvalidRequest(t *testing.T)

// Test progress updates
func TestSyncProgress_ValidUpdate(t *testing.T)
func TestSyncProgress_InvalidToken(t *testing.T)
func TestSyncProgress_OutOfRange(t *testing.T)
```

## Module Organization

### Test Helper Functions (in `routes/routes_test.go`)

1. **Core Infrastructure** (existing, no changes):
   - `setupTestApp()` - Creates in-memory PocketBase instance
   - `runTestMigrations()` - Runs production migrations
   - `parseJSONResponse()` - Unmarshals JSON responses
   - `assertErrorResponse()` - Validates error format

2. **Request Execution** (modified):
   - `makeRequest()` - Executes real HTTP requests through router

3. **Test Data Builders** (existing, no changes):
   - Located in `links_search_test.go`
   - `createTestLinks()` - Creates test link records
   - `buildSearchURL()` - Constructs search URLs

### Test Organization Pattern

Each endpoint gets its own test file following this pattern:
1. Package declaration and imports
2. Response struct definitions (if needed)
3. Helper functions specific to that endpoint
4. Test functions following naming convention: `Test{Endpoint}_{Scenario}`
5. Benchmark functions if performance-critical

### Dependency Management

**Internal Dependencies**:
- `routes.Register()` - Called in `makeRequest()` to register routes
- PocketBase app instance - Passed through from `setupTestApp()`
- Echo router - Created fresh per request

**External Dependencies**:
- No new external dependencies required
- All needed packages already in go.mod

## Data Flow Architecture

### Request Flow in Tests

1. **Test Setup**:
   ```
   TestFunction → setupTestApp() → in-memory PocketBase
   ```

2. **Request Execution**:
   ```
   Test → makeRequest() → Echo Router → Route Handler → PocketBase DAO → SQLite
   ```

3. **Response Validation**:
   ```
   HTTP Response → parseJSONResponse() → Struct → Assertions
   ```

### State Management

**Per-Test State**:
- Fresh PocketBase app instance
- Isolated in-memory database
- Clean route registration

**Per-Request State**:
- New Echo router instance
- Fresh ServeEvent
- Clean HTTP context

**Shared State**:
- None (ensures test isolation)

## Error Handling Structure

### Error Response Format
All endpoints use consistent error structure:
```go
type ErrorResponse struct {
    Error string `json:"error"`
}
```

### Error Test Categories

1. **Resource Not Found** (404):
   - Non-existent link ID
   - Invalid sync session ID
   - Unknown route

2. **Bad Request** (400):
   - Invalid pagination parameters
   - Malformed JSON body
   - Out-of-range values

3. **Internal Errors** (500):
   - Database failures
   - Panic recovery

## Performance Considerations

### Optimization Points

1. **Router Creation**:
   - Created per request (acceptable overhead ~1ms)
   - Could cache if performance becomes issue

2. **Route Registration**:
   - Happens per request (acceptable for <100 tests)
   - Could optimize with singleton pattern if needed

3. **Database Operations**:
   - In-memory SQLite (microsecond operations)
   - No optimization needed

### Benchmarking Structure

Add benchmarks for critical paths:
```go
func BenchmarkMakeRequest(b *testing.B) {
    app, cleanup := setupTestApp(b)
    defer cleanup()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        makeRequest(app, "GET", "/api/health", nil)
    }
}
```

## Migration Path

### Phase 1: Core Fix
1. Update imports in `routes/routes_test.go`
2. Replace `makeRequest()` function
3. Run existing tests to verify

### Phase 2: Validation
1. Add `TestMakeRequest_RealExecution`
2. Verify 404 responses work
3. Confirm existing search tests pass

### Phase 3: Expansion
1. Create `links_view_test.go`
2. Create `stats_test.go`
3. Create `sync_sessions_test.go`

### Phase 4: Cleanup
1. Remove TODO comments
2. Update documentation
3. Verify all tests pass

## File Layout Summary

```
routes/
├── routes_test.go           # Modified: Real makeRequest(), proof test
├── links_search_test.go     # Existing: No changes needed
├── links_view_test.go       # New: View count endpoint tests
├── stats_test.go            # New: Stats endpoint tests
├── sync_sessions_test.go    # New: Sync session tests
└── path_utils_test.go       # Existing: No changes needed
```

This structure maintains backward compatibility while enabling real HTTP testing. The modular organization allows incremental implementation and easy maintenance.