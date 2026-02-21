# Research: Fix Integration Tests

## Problem Statement

The integration tests in `routes/routes_test.go` and `routes/links_search_test.go` use a mocked `makeRequest()` function that returns hardcoded responses instead of executing real HTTP requests against the PocketBase application. This defeats the purpose of integration testing by ensuring tests always pass regardless of actual endpoint functionality.

## Current Test Infrastructure Analysis

### Existing Test Files

1. **`routes/routes_test.go`** - Core test infrastructure (160 lines)
   - Contains `setupTestApp()` for in-memory PocketBase instance creation
   - Contains stubbed `makeRequest()` function at lines 109-128
   - Includes helper functions like `parseJSONResponse()` and `assertErrorResponse()`
   - Has basic tests: `TestSetup`, `TestDatabase`

2. **`routes/links_search_test.go`** - Search endpoint tests (480 lines)
   - Extensive test scenarios: basic query, tag filtering, pagination, SQL injection
   - All tests call the mocked `makeRequest()` function
   - Tests include realistic data creation and validation logic
   - Contains proper test data builders like `createTestLinks()`, `buildSearchURL()`

3. **`routes/path_utils_test.go`** - Utility function tests (working correctly)

### Mock Implementation Issue

The core problem is in `routes/routes_test.go:109-128`:

```go
func makeRequest(app *pocketbase.PocketBase, method, url string, body io.Reader) (*http.Response, error) {
    // TODO: Implement proper HTTP testing for PocketBase v0.36.5
    // For now, return a mock response to test the infrastructure
    // Real implementation will need proper router integration

    // Suppress unused parameter warnings
    _ = app
    _ = method
    _ = url
    _ = body

    // Mock success response for testing infrastructure
    mockResponse := &http.Response{
        StatusCode: 200,
        Header:     make(http.Header),
        Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
    }

    return mockResponse, nil
}
```

### PocketBase Architecture Understanding

From examining the codebase:

1. **Route Registration**: Routes are registered via `routes.Register(app)` in `main.go`
2. **Route Pattern**: Routes use `app.OnServe().BindFunc()` to register handlers
3. **Handler Pattern**: Each route uses `e.Router.GET()` or `e.Router.POST()` in the serve event
4. **Request Pattern**: Handlers receive `*core.RequestEvent` and return errors

Example from `routes/links_search.go:48-50`:
```go
func registerLinksSearch(e *core.ServeEvent) {
    e.Router.GET("/api/links/search", func(ev *core.RequestEvent) error {
        return handleSearch(ev, e.App)
    })
}
```

### Test Infrastructure Status

**Working Components:**
- ✅ `setupTestApp()` creates functional in-memory PocketBase instances
- ✅ Migrations run correctly (`runTestMigrations()`)
- ✅ Database operations work (confirmed 10 seed links)
- ✅ Test data builders (`createTestLinks()`, `buildSearchURL()`)
- ✅ JSON response parsing (`parseJSONResponse()`)
- ✅ Error response validation (`assertErrorResponse()`)

**Broken Component:**
- ❌ `makeRequest()` returns mock data instead of executing real requests

## PocketBase v0.36.5 HTTP Testing Constraints

From examining the codebase and previous implementation work (T-005-01), several approaches were considered:

### Option 1: httptest.NewServer() (Previously Rejected)
- **Issue**: Complex PocketBase app lifecycle management
- **Complexity**: Requires starting actual HTTP server with cleanup challenges
- **Verdict**: Too complex for testing needs

### Option 2: httptest.NewRecorder() (Partially Explored)
- **Approach**: Direct handler invocation with response recording
- **Challenge**: PocketBase routes are registered within OnServe() events
- **Challenge**: Requires echo router setup and manual event triggering

### Option 3: App.OnServe() Simulation (Most Viable)
- **Approach**: Create mock ServeEvent and trigger route registration
- **Approach**: Use httptest.NewRecorder() with configured echo router
- **Benefits**: Stays within PocketBase's architecture patterns
- **Benefits**: Tests actual route handlers without network layer

## Endpoint Testing Requirements

### Primary Test Targets

1. **Search Endpoint** (`GET /api/links/search`)
   - Location: `routes/links_search.go`
   - Handler: `handleSearch()`
   - Features: Text search, tag filtering, pagination, SQL injection protection
   - Current tests: 7 comprehensive scenarios in `links_search_test.go`

2. **View Count Endpoint** (`POST /api/links/:id/view`)
   - Location: `routes/links_view.go`
   - Handler: `handleLinksView()`
   - Features: Atomic view count increment, error handling
   - Current tests: None implemented yet

3. **Stats Endpoint** (`GET /api/stats`)
   - Location: `routes/stats.go`
   - Handler: `handleStats()`
   - Features: Aggregated statistics queries
   - Current tests: None implemented yet

4. **Sync Sessions** (`POST /api/sync`, `PUT /api/sync/:id/progress`)
   - Location: `routes/sync_sessions.go`
   - Handlers: `handleCreateSyncSession()`, `handleUpdateProgress()`
   - Current tests: None implemented yet

### Additional Endpoints

5. **Presentations** (`GET /api/presentations/:id`, `POST /api/presentations`)
   - Location: `routes/presentations.go`
   - Current tests: None implemented yet

6. **Health Check** (`GET /api/health`)
   - Location: `routes/health.go`
   - Simple endpoint, likely working correctly

## Test Data Dependencies

### Required Collections
- `links` - Core entity with title, URL, description, view_count, tags relationships
- `tags` - Tag entities with name and slug fields
- `sync_sessions` - Session tracking with tokens and progress
- `presentations` - Presentation data with metadata

### Seed Data Analysis
From `TestDatabase` output, the test database contains 10 links from migrations. This provides a foundation for testing but specific test scenarios may need additional controlled data.

### Test Data Isolation
Current `setupTestApp()` creates isolated in-memory databases for each test, ensuring proper test isolation without cleanup concerns.

## Dependencies and Imports

### Current Dependencies (Working)
- `github.com/pocketbase/pocketbase v0.36.5` - Core framework
- `github.com/pocketbase/dbx v1.12.0` - Database operations
- `github.com/labstack/echo/v5 v5.0.4` - HTTP router (PocketBase dependency)

### Testing Dependencies (Available)
- Standard `testing` package
- `net/http` for request/response types
- `net/http/httptest` for response recording (commented out but available)
- `encoding/json` for response parsing

## Technical Constraints

### PocketBase v0.36.5 Specific Issues
From previous implementation work, several API changes required fixes:
- Path parameter extraction methods changed
- Database query parameter binding syntax updated
- Response writer patterns modified
- Function naming conflicts resolved

### Router Integration Challenges
The main technical challenge is integrating httptest with PocketBase's event-driven route registration:
1. Routes are registered during OnServe() event handling
2. Each route receives a `*core.RequestEvent` instead of raw `*http.Request`
3. Echo router needs proper initialization and event simulation

### Test Execution Environment
- Tests must run without external PocketBase server
- Must integrate with `go test ./...` command
- Must work in CI/CD environments
- Performance target: Complete test suite under 5 seconds

## Success Criteria Analysis

### Acceptance Criteria Mapping

1. **"makeRequest() uses httptest to make real HTTP requests"**
   - Current: Mocked responses
   - Required: Real HTTP execution through PocketBase router
   - Approach: httptest.NewRecorder() with ServeEvent simulation

2. **"Tests in links_search_test.go actually exercise real search endpoint"**
   - Current: 7 comprehensive test scenarios with mock responses
   - Required: Same scenarios against real search endpoint
   - Challenge: Response validation logic already correct, just need real responses

3. **"Add basic tests for other endpoints"**
   - Current: Only search endpoint tests exist
   - Required: View count, stats, sync create, sync progress tests
   - Scope: 4-5 additional endpoints with basic scenarios

4. **"go test ./... passes with real endpoint tests"**
   - Current: Tests pass with mocks
   - Required: Tests pass with real endpoint calls
   - Constraint: Must not break existing test infrastructure

5. **"At least one test fails if endpoint is broken"**
   - Current: All tests always pass (mocked)
   - Required: Tests detect actual endpoint failures
   - Validation: Need mechanism to prove tests are not vacuous

## Implementation Complexity Assessment

### Low Complexity Areas
- ✅ Test infrastructure exists and works
- ✅ Test data builders are complete
- ✅ Response parsing and validation helpers ready
- ✅ In-memory PocketBase setup proven functional

### Medium Complexity Areas
- 🔄 httptest.NewRecorder() integration with PocketBase
- 🔄 ServeEvent simulation for route registration
- 🔄 Echo router setup and request routing

### High Complexity Areas
- ⚠️ Understanding PocketBase v0.36.5 internal routing mechanics
- ⚠️ RequestEvent creation and parameter passing
- ⚠️ Ensuring test isolation with stateful route registration

## Risk Assessment

### Technical Risks
1. **PocketBase API Compatibility**: v0.36.5 routing internals may have undocumented constraints
2. **Test Isolation**: Route registration might leak state between tests
3. **Performance**: Real HTTP execution might slow down test suite significantly

### Mitigation Strategies
1. **Incremental Approach**: Fix one endpoint at a time, starting with simplest (health check)
2. **Fallback Plan**: If httptest integration fails, document specific limitations
3. **Performance Monitoring**: Measure test execution time and optimize if needed

### Success Probability
- **High (90%)**: Basic httptest integration with simple GET endpoints
- **Medium (70%)**: Complex POST endpoints with request body parsing
- **Medium (60%)**: All endpoints working with full test coverage under performance targets

The research indicates this is primarily an implementation challenge rather than an architectural problem. The test infrastructure foundation is solid, and the main work involves replacing the mocked `makeRequest()` function with real HTTP execution through PocketBase's routing system.