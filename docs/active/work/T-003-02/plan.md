# Plan: View Count Endpoint (T-003-02)

## Implementation Strategy

Based on the research and design phases, this plan sequences the implementation of `POST /api/links/:id/view` into discrete, verifiable steps that can be executed and committed independently.

## Execution Overview

The implementation will be completed in 4 main steps:
1. **Environment Setup & Code Review** - Understand existing patterns
2. **Core Implementation** - Create the new route handler
3. **Integration** - Register the route with the existing system
4. **Testing & Verification** - Ensure correctness and functionality

Each step is designed to be atomic and committable, allowing for easy rollback if issues are discovered.

## Step-by-Step Implementation Plan

### Step 1: Environment Setup & Code Review
**Duration**: 5-10 minutes
**Goal**: Understand the existing codebase patterns and prepare for implementation

**Tasks**:
1. **Review Existing Route Files**
   - Examine `routes/health.go` or similar for import patterns
   - Examine `routes/stats.go` for SQL execution patterns
   - Confirm the parameter extraction pattern used in existing routes

2. **Verify Database Schema**
   - Confirm `links` collection has `view_count` field
   - Check field type and constraints
   - Examine existing seed data to understand current view_count values

3. **Test Current System**
   - Build and run the application
   - Verify existing `/api/links/records` endpoint works
   - Confirm route registration system is functional

**Verification Criteria**:
- [ ] Application builds without errors
- [ ] Existing routes respond correctly
- [ ] `links` collection accessible via API
- [ ] Route patterns understood

**Deliverables**: Understanding of implementation context

### Step 2: Core Implementation - Create Route Handler
**Duration**: 15-20 minutes
**Goal**: Create `routes/links_view.go` with complete endpoint implementation

**Tasks**:
1. **Create Base File Structure**
   ```go
   // routes/links_view.go
   package routes

   import (
       "net/http"
       "github.com/pocketbase/pocketbase/core"
       "github.com/pocketbase/pocketbase/tools/rest"
   )
   ```

2. **Implement Registration Function**
   ```go
   func registerLinksView(e *core.ServeEvent) {
       e.Router.POST("/api/links/:id/view", handleLinksView(e.App))
   }
   ```

3. **Create Handler Function Scaffold**
   - Function signature: `func handleLinksView(app core.App) func(*core.RequestEvent) error`
   - Parameter extraction: `id := e.Request.PathParam("id")`
   - Basic error response structure

4. **Implement Atomic SQL Update**
   - SQL: `UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?`
   - Execute using `app.DB().NewQuery(sql).Exec(id)`
   - Check `result.RowsAffected()` for link existence

5. **Implement Record Retrieval**
   - Use `dao.New(app.DB()).FindRecordById("links", id)`
   - Handle potential errors (should not occur if update succeeded)

6. **Implement Response Formatting**
   - Success: Return record as JSON with 200 status
   - Not Found: Return `{"error": "Link not found"}` with 404 status
   - Server Error: Return generic error with 500 status

**Verification Criteria**:
- [ ] File compiles without errors
- [ ] All imports resolved correctly
- [ ] Function signatures match established patterns
- [ ] SQL syntax is correct
- [ ] Error handling covers all cases

**Deliverables**: Complete `routes/links_view.go` implementation

**Commit Point**: "feat: implement links view count handler"

### Step 3: Integration - Register Route
**Duration**: 5 minutes
**Goal**: Integrate the new route with the existing route registration system

**Tasks**:
1. **Modify `routes/routes.go`**
   - Add `registerLinksView(e)` call in the `Register` function
   - Place after existing route registrations
   - Maintain consistent formatting and ordering

2. **Verify Build**
   - Compile the application to catch any import or integration errors
   - Resolve any compilation issues

**Current Code**:
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        registerHealth(e)
        registerSyncSessions(e)
        registerStats(e)
        registerLinksSearch(e)
        registerLinksSearchSimple(e)
        return e.Next()
    })
}
```

**Target Code**:
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        registerHealth(e)
        registerSyncSessions(e)
        registerStats(e)
        registerLinksSearch(e)
        registerLinksSearchSimple(e)
        registerLinksView(e)          // NEW LINE
        return e.Next()
    })
}
```

**Verification Criteria**:
- [ ] Application builds successfully
- [ ] No compilation errors
- [ ] All existing routes still functional
- [ ] New route registered without conflicts

**Deliverables**: Integrated route registration

**Commit Point**: "feat: register links view count endpoint"

### Step 4: Testing & Verification
**Duration**: 15-20 minutes
**Goal**: Comprehensive testing to ensure endpoint functions correctly under various conditions

**Tasks**:

1. **Basic Functionality Test**
   - Start the application
   - Create test link via PocketBase admin or seed data
   - Send `POST /api/links/{id}/view` request
   - Verify response contains updated view_count
   - Verify response format matches expected JSON structure

2. **Error Condition Testing**
   - Test with non-existent link ID (should return 404)
   - Test with malformed ID (should be handled gracefully)
   - Verify error response format matches specification

3. **Atomicity Verification**
   - Send multiple concurrent requests to same link ID
   - Verify view_count increments correctly (no lost updates)
   - Can be tested with simple script or manual rapid requests

4. **Edge Case Testing**
   - Test with link that has NULL view_count (should increment to 1)
   - Test with link that has view_count = 0 (should increment to 1)
   - Test with link that has existing view_count (should increment by 1)

5. **Integration Testing**
   - Verify existing endpoints still work correctly
   - Verify new endpoint doesn't affect other functionality
   - Test application startup and shutdown

**Test Script Example**:
```bash
# Assume application running on localhost:8090

# Test successful increment
curl -X POST http://localhost:8090/api/links/{existing_id}/view

# Test non-existent link
curl -X POST http://localhost:8090/api/links/nonexistent/view

# Test concurrent requests
for i in {1..10}; do
  curl -X POST http://localhost:8090/api/links/{existing_id}/view &
done
wait
```

**Verification Criteria**:
- [ ] Successful increment returns 200 with updated record
- [ ] Non-existent link returns 404 with error message
- [ ] Concurrent requests don't lose increments
- [ ] NULL view_count handled correctly (COALESCE works)
- [ ] Response format matches PocketBase record structure
- [ ] Existing functionality unaffected

**Deliverables**: Verified working implementation

**Commit Point**: "test: verify links view count endpoint functionality"

## Testing Strategy Details

### Unit Testing Approach
While formal unit tests are not strictly required for this practice application, the testing approach will verify:

1. **Correctness**: Endpoint returns expected data format
2. **Atomicity**: Concurrent requests handled properly
3. **Error Handling**: Appropriate responses for error conditions
4. **Integration**: No interference with existing functionality

### Manual Testing Protocol

**Test Case 1: Successful Increment**
```bash
# Get initial view_count
curl "http://localhost:8090/api/collections/links/records/{id}"

# Increment view_count
curl -X POST "http://localhost:8090/api/links/{id}/view"

# Verify increment occurred
curl "http://localhost:8090/api/collections/links/records/{id}"
```

**Test Case 2: Non-existent Link**
```bash
curl -X POST "http://localhost:8090/api/links/invalid_id/view"
# Expected: 404 with {"error": "Link not found"}
```

**Test Case 3: Concurrent Requests**
```bash
# Run multiple increments simultaneously
for i in {1..5}; do
  curl -X POST "http://localhost:8090/api/links/{id}/view" &
done
wait

# Verify all increments were counted
curl "http://localhost:8090/api/collections/links/records/{id}"
```

### Performance Testing
Basic performance characteristics to verify:
- Response time under normal load (< 100ms expected)
- No memory leaks under repeated requests
- Database connection handling (no connection exhaustion)

## Risk Management

### Potential Issues and Mitigation

**Issue 1: Race Conditions**
- **Risk**: Concurrent requests lose increments
- **Mitigation**: Atomic SQL operation handles concurrency
- **Verification**: Concurrent request testing

**Issue 2: Database Errors**
- **Risk**: SQL execution fails, endpoint returns 500
- **Mitigation**: Proper error handling and graceful degradation
- **Verification**: Error condition testing

**Issue 3: Integration Conflicts**
- **Risk**: New route conflicts with existing routes or patterns
- **Mitigation**: Follow established patterns, minimal changes
- **Verification**: Integration testing of existing endpoints

**Issue 4: Data Consistency**
- **Risk**: view_count becomes inconsistent or corrupted
- **Mitigation**: COALESCE handles NULL values, atomic operations
- **Verification**: Edge case testing with NULL values

### Rollback Procedures

**If Step 2 Fails**: Delete `routes/links_view.go`, no other changes needed
**If Step 3 Fails**: Revert changes to `routes/routes.go`, keep implementation file for debugging
**If Step 4 Fails**: Keep implementation, investigate issues, fix incrementally

## Quality Assurance

### Code Quality Standards
- Follow existing code patterns and naming conventions
- Use proper error handling and logging
- Maintain consistency with PocketBase idioms
- Include appropriate comments for complex logic

### Documentation Requirements
- Update progress.md with implementation progress
- Document any deviations from the plan
- Record test results and verification outcomes

### Success Criteria
The implementation will be considered complete when:
- [ ] All acceptance criteria from the ticket are met
- [ ] All verification criteria from each step are satisfied
- [ ] No existing functionality is broken
- [ ] Performance is acceptable for the use case
- [ ] Error handling is robust and user-friendly

## Post-Implementation

### Monitoring Recommendations
For production deployment (beyond this practice application):
- Monitor endpoint response times
- Track error rates and types
- Consider rate limiting for abuse prevention
- Log increment patterns for analytics

### Future Enhancement Opportunities
- Rate limiting per IP address
- Authentication requirements
- Bulk increment operations
- View count analytics and reporting
- Caching for frequently accessed links

This plan provides a systematic approach to implementing the view count endpoint while maintaining code quality and minimizing implementation risks.