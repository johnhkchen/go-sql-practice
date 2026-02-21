# Structure: Go Integration Tests for Custom API Endpoints

## File Structure Overview

The integration tests will be implemented as a single comprehensive test file within the `routes` package to maintain co-location with the tested code and provide access to internal functions.

```
routes/
├── routes.go              # Existing - route registration hub
├── links_search.go        # Existing - search endpoint implementation
├── links_view.go          # Existing - view count endpoint implementation
├── stats.go               # Existing - stats endpoint implementation
└── routes_test.go         # NEW - comprehensive integration tests
```

## New File: `routes/routes_test.go`

### Package Declaration and Imports

```go
package routes

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "sync"
    "testing"

    "github.com/labstack/echo/v5"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/core"
    "github.com/pocketbase/pocketbase/migrations"
)
```

**Dependencies Rationale**:
- Standard `testing` and `net/http/httptest` for Go testing patterns
- `sync` for concurrent testing scenarios
- PocketBase imports for app setup and migration handling
- Echo for HTTP routing (PocketBase dependency)

### Test Helper Infrastructure

#### Core App Setup Helper

```go
// setupTestApp creates an in-memory PocketBase app for testing
func setupTestApp(t *testing.T) (*pocketbase.PocketBase, func()) {
    app := pocketbase.NewWithConfig(&pocketbase.Config{
        DefaultDataDir: "", // Forces in-memory SQLite database
    })

    // Initialize the app without starting HTTP server
    if err := app.Bootstrap(); err != nil {
        t.Fatalf("Failed to bootstrap test app: %v", err)
    }

    // Run migrations to create collections
    if err := runTestMigrations(app); err != nil {
        t.Fatalf("Failed to run test migrations: %v", err)
    }

    // Register routes
    router := echo.New()
    serveEvent := &core.ServeEvent{App: app, Router: router}
    Register(serveEvent)

    cleanup := func() {
        if app.DB() != nil {
            app.DB().Close()
        }
    }

    return app, cleanup
}
```

#### Migration Helper

```go
// runTestMigrations runs the production migrations in test environment
func runTestMigrations(app *pocketbase.PocketBase) error {
    // Import production migrations
    return migrations.Register(app, "")
}
```

#### HTTP Request Helper

```go
// makeRequest executes an HTTP request against the test app
func makeRequest(app *pocketbase.PocketBase, method, url string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    rec := httptest.NewRecorder()
    app.OnServe().Trigger(&core.ServeEvent{
        App:    app,
        Router: app.Echo(),
    })

    app.Echo().ServeHTTP(rec, req)
    return rec.Result(), nil
}
```

#### JSON Response Helper

```go
// parseJSONResponse unmarshals HTTP response body into target struct
func parseJSONResponse(t *testing.T, resp *http.Response, target interface{}) {
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }

    if err := json.Unmarshal(body, target); err != nil {
        t.Fatalf("Failed to unmarshal JSON response: %v\nBody: %s", err, string(body))
    }
}
```

#### Test Data Builders

```go
// TestData holds references to created test data
type TestData struct {
    LinkIDs []string
    TagIDs  []string
    TagSlugs []string
}

// seedTestData creates predictable test data for search and stats tests
func seedTestData(t *testing.T, app *pocketbase.PocketBase) *TestData {
    // Implementation creates links and tags with known relationships
}

// createTestLink creates a single link for targeted testing
func createTestLink(t *testing.T, app *pocketbase.PocketBase, title, url string, viewCount int, tagSlugs []string) string {
    // Implementation creates link and returns ID
}

// createTestTag creates a single tag for testing
func createTestTag(t *testing.T, app *pocketbase.PocketBase, name, slug string) string {
    // Implementation creates tag and returns ID
}
```

### Test Function Organization

#### Search Endpoint Tests

**Test Function**: `TestLinksSearch`
**Sub-tests** (using `t.Run`):
- `"BasicTextSearch"` - Text query with expected matches
- `"TagFilter"` - Single tag filtering
- `"CombinedFilters"` - Text + tag filtering
- `"Pagination"` - Multiple pages with correct counts
- `"EmptyResults"` - Queries returning no matches
- `"InvalidPagination"` - Bad page/perPage parameters
- `"SQLInjectionProtection"` - Malicious input handling

**Test Structure Pattern**:
```go
func TestLinksSearch(t *testing.T) {
    app, cleanup := setupTestApp(t)
    defer cleanup()

    // Seed comprehensive test data
    testData := seedTestData(t, app)

    t.Run("BasicTextSearch", func(t *testing.T) {
        // Test implementation
    })

    // Additional sub-tests...
}
```

#### View Count Endpoint Tests

**Test Function**: `TestLinksView`
**Sub-tests**:
- `"SuccessfulIncrement"` - POST to existing link
- `"NonexistentLink"` - POST to invalid ID returns 404
- `"AtomicUpdates"` - Concurrent increments
- `"ZeroToOne"` - Null view_count initialization
- `"ResponseFormat"` - Validate updated record structure

**Concurrent Test Pattern**:
```go
t.Run("AtomicUpdates", func(t *testing.T) {
    const numIncrements = 10
    var wg sync.WaitGroup

    for i := 0; i < numIncrements; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Make increment request
        }()
    }

    wg.Wait()
    // Verify final count
})
```

#### Stats Endpoint Tests

**Test Function**: `TestStats`
**Sub-tests**:
- `"CompleteResponse"` - All fields populated correctly
- `"DataAccuracy"` - Counts match seed data
- `"TopItems"` - Correct ordering of most viewed/top tags
- `"EmptyDatabase"` - Graceful zero-data handling
- `"ResponseSchema"` - JSON structure validation

### Error Response Testing

#### Shared Error Response Structure

```go
// ErrorResponse matches the common error format across all endpoints
type ErrorResponse struct {
    Error string `json:"error"`
}

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

### Test Data Architecture

#### Predictable Test Data Design

**Search Test Data**:
- 8 tags matching production seed (golang, javascript, etc.)
- 15+ links with varying title/description combinations
- Controlled tag associations for filtering tests
- Sufficient volume for pagination (3+ pages)

**Stats Test Data**:
- Known quantities: 5 links, 3 tags, predictable view counts
- Links with view counts: [45, 23, 12, 8, 0]
- Tags with link counts: golang(3), javascript(2), database(1)

**View Count Test Data**:
- Links with various initial states: null, 0, positive values
- Known link IDs for positive tests
- Non-existent link IDs for 404 tests

#### Data Isolation Strategy

Each test function creates its own isolated data set:
- No shared global fixtures
- Independent test data prevents coupling
- Cleanup handled by in-memory database destruction

### File Dependencies and Integration

#### Import Dependencies

The test file depends on:
- **Internal**: All route handler functions and types from same package
- **PocketBase**: Core app functionality and migration system
- **Standard Library**: Testing, HTTP, JSON utilities
- **Echo**: HTTP router for request handling

#### Migration Integration

Tests reuse production migrations:
- `migrations/collections.go` - Collection definitions
- `migrations/seed.go` - Reference for data structure (not used directly in tests)

#### Route Registration Integration

Tests use the existing `Register()` function from `routes.go`:
- Ensures same route registration path as production
- Maintains consistency with main application setup
- Tests actual route patterns and handlers

### Test Execution Integration

#### Makefile Integration

The existing `Makefile` `test` target automatically includes the new test file:
```makefile
test:
    go test ./...
```

No modifications to build system required.

#### CI/CD Compatibility

Test design supports standard Go testing workflows:
- No external dependencies (database, services)
- Fast execution through in-memory database
- Compatible with coverage reporting tools
- Supports parallel test execution

### Performance Considerations

#### In-Memory Database Benefits

- Database creation/destruction: <10ms per test
- Query execution: Microsecond range
- Full test suite execution: <5 seconds expected

#### Resource Management

- Each test gets isolated database instance
- Cleanup functions prevent resource leaks
- Concurrent tests supported through isolation

### Debugging and Maintenance

#### Debug Helpers

```go
// dumpTestData logs current database state for debugging
func dumpTestData(t *testing.T, app *pocketbase.PocketBase) {
    // Implementation logs table contents
}

// logRequest logs HTTP requests during test development
func logRequest(t *testing.T, req *http.Request, resp *http.Response) {
    // Implementation logs request/response details
}
```

#### Maintenance Hooks

- Test data builders centralize test data creation
- Shared helpers reduce code duplication
- Error assertion helpers ensure consistent validation

This structure provides comprehensive integration test coverage while maintaining clean organization, performance, and maintainability. The single-file approach keeps related tests together while leveraging Go's sub-test organization for logical grouping.