# Design: Fix Integration Tests

## Problem Restatement

Integration tests currently use a mocked `makeRequest()` function that returns static JSON responses, making tests vacuous. We need real HTTP request execution through PocketBase's routing system to validate actual endpoint functionality.

## Solution Approaches

### Approach 1: Direct Echo Router Invocation

**Concept**: Bypass PocketBase's event system and directly invoke Echo router with httptest.

**Implementation**:
```go
func makeRequest(app *pocketbase.PocketBase, method, url string, body io.Reader) (*http.Response, error) {
    req := httptest.NewRequest(method, url, body)
    rec := httptest.NewRecorder()

    // Get Echo instance directly
    echo := app.Echo()
    echo.ServeHTTP(rec, req)

    return rec.Result(), nil
}
```

**Pros**:
- Simple implementation
- Standard httptest patterns
- Direct HTTP execution

**Cons**:
- Routes not registered (app.Echo() is nil before OnServe)
- Bypasses PocketBase middleware
- Missing authentication and context setup

**Verdict**: ❌ Won't work - routes aren't registered without OnServe event

### Approach 2: OnServe Event Simulation

**Concept**: Trigger OnServe event to register routes, then use Echo router.

**Implementation**:
```go
func makeRequest(app *pocketbase.PocketBase, method, url string, body io.Reader) (*http.Response, error) {
    // Trigger route registration once
    if !routesRegistered {
        e := &core.ServeEvent{
            App:    app,
            Router: echo.New(),
        }
        app.OnServe().Trigger(e)
        routesRegistered = true
    }

    req := httptest.NewRequest(method, url, body)
    rec := httptest.NewRecorder()
    app.Echo().ServeHTTP(rec, req)

    return rec.Result(), nil
}
```

**Pros**:
- Uses PocketBase's official route registration
- Maintains middleware chain
- Routes registered properly

**Cons**:
- Requires managing global state (routesRegistered flag)
- Echo instance might not persist correctly
- Complex lifecycle management

**Verdict**: ⚠️ Possible but fragile

### Approach 3: Per-Test Router Setup

**Concept**: Create fresh router for each test, register routes, execute requests.

**Implementation**:
```go
func makeRequest(app *pocketbase.PocketBase, method, url string, body io.Reader) (*http.Response, error) {
    router := echo.New()

    // Create and trigger serve event
    serveEvent := &core.ServeEvent{
        App:    app,
        Router: router,
    }

    // Register routes
    routes.Register(app)
    app.OnServe().Trigger(serveEvent)

    // Execute request
    req := httptest.NewRequest(method, url, body)
    rec := httptest.NewRecorder()
    router.ServeHTTP(rec, req)

    return rec.Result(), nil
}
```

**Pros**:
- Clean test isolation
- No global state
- Fresh router per request

**Cons**:
- Performance overhead (route registration per request)
- Might trigger side effects in OnServe handlers
- Router middleware might be incomplete

**Verdict**: ✅ Most robust approach

### Approach 4: Test-Specific Router Wrapper

**Concept**: Create a test helper that manages router lifecycle at test level.

**Implementation**:
```go
type TestContext struct {
    App    *pocketbase.PocketBase
    Router *echo.Echo
}

func setupTestContext(t *testing.T) (*TestContext, func()) {
    app := setupTestApp(t)
    router := echo.New()

    // Register routes once
    serveEvent := &core.ServeEvent{
        App:    app,
        Router: router,
    }
    routes.Register(app)
    app.OnServe().Trigger(serveEvent)

    ctx := &TestContext{
        App:    app,
        Router: router,
    }

    cleanup := func() {
        // Cleanup logic
    }

    return ctx, cleanup
}

func (tc *TestContext) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
    req := httptest.NewRequest(method, url, body)
    rec := httptest.NewRecorder()
    tc.Router.ServeHTTP(rec, req)
    return rec.Result(), nil
}
```

**Pros**:
- Routes registered once per test
- Clean separation of concerns
- Efficient execution

**Cons**:
- Requires refactoring all existing tests
- More complex test structure
- Breaking change to test API

**Verdict**: ⚠️ Better architecture but high refactoring cost

## Decision: Approach 3 - Per-Test Router Setup

### Rationale

1. **Simplicity**: Drop-in replacement for existing mock function
2. **Compatibility**: No changes needed to existing test cases
3. **Isolation**: Each request gets clean router state
4. **Reliability**: Follows PocketBase's intended patterns

### Implementation Details

The chosen approach will:

1. Create a fresh Echo router for each request
2. Trigger OnServe event to register all routes
3. Execute the HTTP request through the router
4. Return the actual HTTP response

### Performance Optimization

To address performance concerns:

1. **Lazy Registration**: Cache router setup per app instance
2. **Minimal Middleware**: Skip unnecessary middleware in tests
3. **Direct Routing**: Bypass authentication for test requests

### Error Handling Strategy

The implementation must handle:

1. **Route Not Found**: Return 404 responses correctly
2. **Handler Panics**: Capture and report as test failures
3. **Database Errors**: Propagate actual database errors

### Validation Approach

To prove tests aren't vacuous:

1. **Negative Test**: Add test that expects 404 for non-existent route
2. **Error Test**: Add test that triggers actual database error
3. **Mutation Test**: Temporarily break an endpoint and verify test fails

## Test Coverage Expansion

### Phase 1: Fix makeRequest() (Critical)
- Replace mock implementation
- Verify existing search tests still pass
- Add negative test to prove real execution

### Phase 2: View Count Tests (High Priority)
```go
func TestLinksView_Success(t *testing.T)
func TestLinksView_NotFound(t *testing.T)
func TestLinksView_Concurrent(t *testing.T)
```

### Phase 3: Stats Endpoint Tests (Medium Priority)
```go
func TestStats_Complete(t *testing.T)
func TestStats_EmptyDatabase(t *testing.T)
```

### Phase 4: Sync Session Tests (Medium Priority)
```go
func TestSyncCreate_Success(t *testing.T)
func TestSyncProgress_Update(t *testing.T)
func TestSyncProgress_InvalidToken(t *testing.T)
```

### Phase 5: Presentation Tests (Low Priority)
```go
func TestPresentations_Get(t *testing.T)
func TestPresentations_Create(t *testing.T)
```

## Risk Mitigation

### Risk 1: Router Registration Side Effects
**Mitigation**: Test in isolated goroutines if needed

### Risk 2: Performance Degradation
**Mitigation**: Benchmark before/after, optimize if >2x slower

### Risk 3: Unexpected PocketBase Behavior
**Mitigation**: Add detailed logging during initial implementation

## Success Metrics

1. **All existing tests pass**: No regression in test behavior
2. **Real endpoint execution**: Verified by intentional failure test
3. **Performance maintained**: Test suite completes in <5 seconds
4. **Coverage expanded**: At least 5 new endpoint tests added

## Implementation Sequence

1. Create new `makeRequest()` implementation
2. Run existing search tests to verify compatibility
3. Add proof-of-execution test (404 scenario)
4. Add view count endpoint tests
5. Add stats endpoint tests
6. Add sync session tests
7. Clean up and optimize

The design prioritizes compatibility and simplicity while ensuring real HTTP execution. The per-request router setup provides clean isolation at acceptable performance cost.