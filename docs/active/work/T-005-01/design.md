# Design: Go Integration Tests for Custom API Endpoints

## Problem Statement

The go-sql-practice application has three custom API endpoints that need comprehensive integration testing:
1. Search endpoint (`GET /api/links/search`) - complex query logic with pagination
2. View count endpoint (`POST /api/links/:id/view`) - atomic updates with error handling
3. Stats endpoint (`GET /api/stats`) - aggregated data queries

These endpoints require full HTTP stack testing with realistic data scenarios, while maintaining test isolation and avoiding external dependencies.

## Testing Architecture Options

### Option 1: Standard Go HTTP Testing with Test Server

**Approach**: Use `httptest.NewServer()` with a real PocketBase instance configured for testing.

**Pros**:
- Standard Go testing patterns
- Full HTTP stack coverage
- Easy request/response validation
- Familiar to most Go developers

**Cons**:
- Requires complex PocketBase app lifecycle management
- Database cleanup between tests is error-prone
- No built-in support for PocketBase's migration system
- Difficult to test error conditions in database layer

**Implementation**:
```go
func setupTestServer(t *testing.T) *httptest.Server {
    app := pocketbase.NewWithConfig(&pocketbase.Config{...})
    // Complex setup required
}
```

### Option 2: PocketBase Native Testing with In-Memory Database

**Approach**: Use PocketBase's built-in testing capabilities with `pocketbase.NewWithConfig()` and in-memory SQLite.

**Pros**:
- Leverages PocketBase's testing infrastructure
- In-memory database provides perfect isolation
- Migration system works natively
- Full integration with PocketBase's DAO layer
- Realistic test environment matching production

**Cons**:
- Less familiar to developers without PocketBase experience
- Requires understanding PocketBase's internal testing patterns
- Documentation is limited compared to standard Go testing

**Implementation**:
```go
func setupTestApp(t *testing.T) *pocketbase.PocketBase {
    app := pocketbase.NewWithConfig(&pocketbase.Config{
        DefaultDataDir: "",  // In-memory
    })
    // Direct access to app.TestDataPath() and app.TestMailer()
}
```

### Option 3: Hybrid Approach with Mocked Dependencies

**Approach**: Mock the database layer while testing HTTP handlers directly.

**Pros**:
- Fast test execution
- Precise control over database responses
- Easy to test error conditions

**Cons**:
- Doesn't test actual SQL queries
- Complex mock setup for realistic scenarios
- Breaks integration testing principles
- Maintenance burden for mock synchronization

**Implementation**:
```go
type MockDB struct{}
func (m *MockDB) NewQuery(sql string) dbx.Query { ... }
```

## Decision: Option 2 - PocketBase Native Testing

**Chosen Approach**: PocketBase native testing with in-memory database configuration.

**Rationale**:

1. **Integration Fidelity**: Testing against an actual PocketBase instance with real SQLite database ensures our tests validate the complete integration stack, including PocketBase's internal query execution, JSON handling, and response formatting.

2. **Test Isolation**: In-memory SQLite databases provide perfect isolation between tests with automatic cleanup on app destruction.

3. **Migration Compatibility**: PocketBase's migration system works seamlessly with test apps, allowing us to use the same collection definitions and seed data patterns used in production.

4. **Error Testing**: Real database operations allow testing actual error conditions (constraint violations, malformed queries, etc.) rather than simulated errors.

5. **Maintenance**: Using the same code paths as production reduces the gap between test and production behavior, minimizing maintenance overhead.

**Rejected Alternatives**:
- Option 1 was rejected due to complex lifecycle management and cleanup challenges
- Option 3 was rejected as it doesn't provide true integration testing value

## Test Structure Design

### Package Organization

**Chosen Structure**: `routes/routes_test.go` with helper functions
- Co-locates tests with the code being tested
- Follows Go testing conventions
- Provides access to unexported functions for detailed testing

**Alternative Considered**: Separate `routes_test` package
- Rejected due to loss of access to internal functions and types
- Would require exposing internal APIs just for testing

### Test Helper Architecture

**Core Helper Function**:
```go
func setupTestApp(t *testing.T) (*pocketbase.PocketBase, func()) {
    app := pocketbase.NewWithConfig(&pocketbase.Config{
        DefaultDataDir: "",  // Forces in-memory SQLite
    })

    // Initialize collections and seed data
    if err := runMigrations(app); err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }

    cleanup := func() {
        app.OnTerminate().Trigger(&core.TerminateEvent{App: app})
    }

    return app, cleanup
}
```

**Test Data Strategy**:
- Minimal, focused seed data per test
- Avoid shared fixtures that couple tests
- Use builder pattern for test data construction

### HTTP Request Testing Pattern

**Request Helper**:
```go
func makeRequest(app *pocketbase.PocketBase, method, path string, body io.Reader) (*httptest.ResponseRecorder, error) {
    req := httptest.NewRequest(method, path, body)
    rec := httptest.NewRecorder()

    // Create a mock ServeEvent to trigger route registration
    app.OnServe().Trigger(&core.ServeEvent{
        App:    app,
        Router: echo.New(),
    })

    // Execute request through PocketBase's router
    app.Echo().ServeHTTP(rec, req)
    return rec, nil
}
```

**Response Validation**:
- JSON unmarshaling with struct validation
- Status code assertions
- Header validation where relevant

## Test Coverage Strategy

### Search Endpoint Tests

**Core Scenarios**:
1. **Basic Search**: Text query with expected results
2. **Tag Filtering**: Single tag filter with result validation
3. **Combined Filters**: Text + tag filtering
4. **Pagination**: Multiple pages with correct counts
5. **Empty Results**: Queries that return no matches
6. **Parameter Validation**: Invalid page/perPage values
7. **SQL Injection Protection**: Malicious input handling

**Test Data Requirements**:
- Links with varying titles and descriptions for text search
- Links with different tag combinations
- Sufficient volume for pagination testing (25+ links)

### View Count Endpoint Tests

**Core Scenarios**:
1. **Successful Increment**: POST to existing link, verify count increase
2. **Nonexistent Link**: POST to invalid ID returns 404
3. **Atomic Updates**: Concurrent increments maintain consistency
4. **Zero to One**: Links with null view_count initialize correctly
5. **Response Format**: Updated record matches expected structure

**Test Data Requirements**:
- Links with various initial view_count values (0, null, positive)
- Known link IDs for positive tests
- Invalid link IDs for error testing

### Stats Endpoint Tests

**Core Scenarios**:
1. **Complete Response**: All stats fields populated correctly
2. **Data Accuracy**: Counts match seed data precisely
3. **Top Items**: Most viewed links and top tags in correct order
4. **Empty Database**: Graceful handling of zero data
5. **Response Structure**: JSON schema validation

**Test Data Requirements**:
- Predictable link and tag counts
- Links with known view_count distribution
- Tags with known link association counts

## Error Handling Strategy

### Database Error Simulation

PocketBase's in-memory database allows testing real database errors:
- Constraint violations through malformed data
- Query failures through invalid SQL (if exposed)
- Connection errors through app shutdown timing

### HTTP Error Validation

Test error response format consistency:
```go
type ErrorResponse struct {
    Error string `json:"error"`
}
```

All endpoints return errors in this format, tests should validate both status codes and error message structure.

### Concurrent Testing

For view count endpoint, test concurrent increments using goroutines:
```go
func TestLinksViewConcurrent(t *testing.T) {
    // Launch multiple goroutines incrementing same link
    // Verify final count equals number of increments
}
```

## Test Execution Environment

### Go Test Integration

Tests integrate with standard `go test` command and Makefile `test` target:
- No external dependencies required
- Runs in CI/CD pipelines without setup
- Compatible with coverage reporting tools

### Performance Considerations

In-memory SQLite provides excellent test performance:
- Database creation/destruction in milliseconds
- Full test suite should complete in seconds
- No I/O bottlenecks from disk operations

### Debugging Support

Test helpers include debugging utilities:
- Database state inspection functions
- Request/response logging for failed tests
- SQL query logging during test development

This design provides comprehensive integration test coverage while maintaining test isolation, performance, and maintainability. The PocketBase-native approach ensures tests accurately reflect production behavior while leveraging the framework's built-in testing capabilities.