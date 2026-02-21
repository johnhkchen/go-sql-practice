# Plan: Fix Integration Tests

## Implementation Steps

### Step 1: Fix Core makeRequest() Function

**Files**: `routes/routes_test.go`

**Actions**:
1. Uncomment `"net/http/httptest"` import
2. Add required imports: `"github.com/labstack/echo/v5"`
3. Replace mock `makeRequest()` function (lines 109-128) with real implementation
4. Test with simple health check endpoint

**Verification**:
- Run: `go test -v ./routes -run TestSetup`
- Expect: Test still passes with new implementation

**Commit**: "fix: replace mock makeRequest with real HTTP execution"

### Step 2: Add Proof of Execution Test

**Files**: `routes/routes_test.go`

**Actions**:
1. Add `TestMakeRequest_RealExecution` function
2. Test valid endpoint (health check) returns 200
3. Test invalid endpoint returns 404
4. Test that response bodies are different

**Verification**:
- Run: `go test -v ./routes -run TestMakeRequest_RealExecution`
- Expect: Both assertions pass

**Commit**: "test: add proof of real HTTP execution"

### Step 3: Validate Existing Search Tests

**Files**: None (validation only)

**Actions**:
1. Run all search tests with new `makeRequest()`
2. Debug any failures (likely JSON response format differences)
3. Adjust response parsing if needed

**Verification**:
- Run: `go test -v ./routes -run TestLinksSearch`
- Expect: All 7 search test scenarios pass

**Commit**: "test: verify search tests work with real requests"

### Step 4: Implement View Count Tests

**Files**: `routes/links_view_test.go` (new)

**Actions**:
1. Create new test file
2. Implement `TestLinksView_Success` - increment existing link
3. Implement `TestLinksView_NotFound` - attempt on non-existent ID
4. Implement `TestLinksView_InvalidID` - malformed ID format

**Verification**:
- Run: `go test -v ./routes -run TestLinksView`
- Expect: All view count tests pass

**Commit**: "test: add view count endpoint integration tests"

### Step 5: Implement Stats Endpoint Tests

**Files**: `routes/stats_test.go` (new)

**Actions**:
1. Create new test file
2. Define `StatsTestResponse` struct for unmarshaling
3. Implement `TestStats_Complete` - verify all fields populated
4. Implement `TestStats_Accuracy` - verify counts match seed data

**Verification**:
- Run: `go test -v ./routes -run TestStats`
- Expect: Stats tests pass with correct counts

**Commit**: "test: add stats endpoint integration tests"

### Step 6: Implement Sync Session Tests

**Files**: `routes/sync_sessions_test.go` (new)

**Actions**:
1. Create new test file
2. Implement `TestSyncCreate_Success` - create new session
3. Implement `TestSyncProgress_ValidUpdate` - update progress
4. Implement `TestSyncProgress_InvalidToken` - wrong token rejection

**Verification**:
- Run: `go test -v ./routes -run TestSync`
- Expect: Sync session tests pass

**Commit**: "test: add sync session endpoint integration tests"

### Step 7: Add Concurrent Test for View Count

**Files**: `routes/links_view_test.go`

**Actions**:
1. Implement `TestLinksView_Concurrent`
2. Use goroutines to increment same link 10 times concurrently
3. Verify final count equals initial + 10

**Verification**:
- Run: `go test -v ./routes -run TestLinksView_Concurrent -race`
- Expect: No race conditions, correct final count

**Commit**: "test: add concurrent view count increment test"

### Step 8: Performance Validation

**Files**: `routes/routes_test.go`

**Actions**:
1. Add `BenchmarkMakeRequest` function
2. Benchmark request execution time
3. Document performance metrics

**Verification**:
- Run: `go test -bench=. ./routes`
- Expect: Request time <10ms, full suite <5s

**Commit**: "test: add performance benchmarks"

### Step 9: Cleanup and Documentation

**Files**: Multiple

**Actions**:
1. Remove TODO comments from `routes/routes_test.go`
2. Update any misleading comments
3. Ensure all test files have package documentation

**Verification**:
- Run: `go test ./...`
- Expect: All tests pass

**Commit**: "docs: clean up test documentation and TODOs"

### Step 10: Final Validation

**Files**: `docs/active/work/T-009-05/progress.md`

**Actions**:
1. Run complete test suite
2. Intentionally break an endpoint to verify test catches it
3. Fix the endpoint and verify tests pass again
4. Document results in progress.md

**Verification**:
- Run: `go test -v ./...`
- Expect: Full test suite passes
- Break endpoint, verify test fails
- Fix endpoint, verify test passes

**Commit**: "test: complete integration test fixes for T-009-05"

## Testing Strategy

### Unit Test Approach
Each endpoint gets focused tests:
- **Positive cases**: Expected successful operations
- **Negative cases**: Error conditions and edge cases
- **Validation cases**: Parameter validation and constraints

### Integration Test Approach
Tests verify complete request/response cycle:
- HTTP request routing
- Parameter extraction
- Database operations
- Response formatting

### Regression Prevention
- Keep all existing test scenarios
- Add new tests for uncovered paths
- Document expected behavior

## Risk Management

### Risk: Router Registration Fails
**Mitigation**: Add detailed error logging in `makeRequest()`
**Fallback**: Revert to mock temporarily if blocking

### Risk: Tests Too Slow
**Mitigation**: Benchmark after each step
**Fallback**: Cache router if >10ms per request

### Risk: Unexpected Response Format
**Mitigation**: Log actual vs expected for debugging
**Fallback**: Adjust response parsing helpers

## Success Criteria Validation

### Criterion 1: Real HTTP Requests
- ✓ Step 1 implements httptest.NewRecorder()
- ✓ Step 2 proves real execution with 404 test

### Criterion 2: Search Tests Exercise Real Endpoint
- ✓ Step 3 validates all search tests work
- ✓ Response data from actual database queries

### Criterion 3: Basic Tests for Other Endpoints
- ✓ Steps 4-6 add tests for view, stats, sync
- ✓ Minimum 2 tests per endpoint

### Criterion 4: go test ./... Passes
- ✓ Step 10 validates complete suite
- ✓ No regressions introduced

### Criterion 5: Tests Detect Broken Endpoints
- ✓ Step 10 includes break/fix validation
- ✓ Proves tests aren't vacuous

## Time Estimates

- Step 1: 15 minutes (core fix)
- Step 2: 10 minutes (proof test)
- Step 3: 20 minutes (validation/debugging)
- Step 4: 20 minutes (view count tests)
- Step 5: 15 minutes (stats tests)
- Step 6: 20 minutes (sync tests)
- Step 7: 15 minutes (concurrent test)
- Step 8: 10 minutes (benchmarks)
- Step 9: 10 minutes (cleanup)
- Step 10: 15 minutes (validation)

**Total Estimate**: 2.5 hours

## Rollback Plan

If implementation encounters blockers:
1. Keep test infrastructure improvements
2. Document specific PocketBase limitations
3. Create issue for future investigation
4. Partial implementation better than none

The plan provides incremental value with each step, allowing partial completion if needed while prioritizing the most critical fixes first.