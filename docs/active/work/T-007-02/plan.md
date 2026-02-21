# Plan: presentation-api-routes (T-007-02)

## Overview

This plan sequences the implementation of presentation API routes into discrete, testable steps. Each step builds incrementally while maintaining system stability and enabling atomic commits.

## Implementation Strategy

### Approach
- **Incremental Development**: Build routes one at a time, testing each thoroughly
- **Risk-First Ordering**: Implement lowest-risk components first (status endpoint)
- **Foundation-First**: Create shared utilities before specific handlers
- **Test-Driven**: Verify each component before moving to the next
- **Atomic Commits**: Each step should be independently committable

### Testing Strategy
- **Manual Testing**: Use curl/Postman to test each endpoint
- **Integration Testing**: Test with actual PocketBase database and auth
- **Error Path Testing**: Verify error conditions return correct responses
- **Frontend Integration**: Test compatibility with existing frontend polling

## Implementation Steps

### Step 1: Foundation Setup
**Goal**: Create basic file structure and utilities

**Tasks**:
1. Create `routes/presentations.go` with package structure and imports
2. Implement shared constants and type definitions
3. Implement utility functions (progress calculations, response builders)
4. Implement authentication helpers

**Files Created**:
- `routes/presentations.go` (partial, ~60 lines)

**Verification**:
- File compiles without errors
- Utility functions work with test inputs
- Authentication helpers return expected results with mock data

**Commit**: "feat: add presentation routes foundation and utilities"

### Step 2: Status Endpoint (GET /status)
**Goal**: Implement read-only status endpoint - lowest risk, no state changes

**Tasks**:
1. Implement `handleGetStatus` function
2. Implement `buildStatusResponse` helper
3. Add route registration for status endpoint only
4. Update `routes/routes.go` to call `registerPresentations`

**Files Modified**:
- `routes/presentations.go` (+40 lines)
- `routes/routes.go` (+1 line)

**Verification**:
- GET /api/presentations/:id/status returns correct data for existing presentations
- Returns 404 for non-existent presentations
- Correctly computes is_live, progress, current_step fields
- No authentication required (public read access)

**Test Cases**:
```bash
# Test existing presentation without active session
curl -X GET http://localhost:8090/api/presentations/{id}/status

# Test existing presentation with active session
# (manually create session, link to presentation first)
curl -X GET http://localhost:8090/api/presentations/{id}/status

# Test non-existent presentation
curl -X GET http://localhost:8090/api/presentations/invalid/status
```

**Commit**: "feat: implement GET /api/presentations/:id/status endpoint"

### Step 3: Stop Endpoint (POST /stop)
**Goal**: Implement session stopping - simple state change, lower risk than creation

**Tasks**:
1. Implement `handleStopLive` function
2. Add authentication and ownership validation
3. Add route registration for stop endpoint
4. Test with various auth scenarios

**Files Modified**:
- `routes/presentations.go` (+45 lines)

**Verification**:
- POST /api/presentations/:id/stop clears active_session field
- Requires valid authentication token
- Requires presentation ownership
- Returns 409 if presentation not live
- Returns appropriate error codes for auth failures

**Test Cases**:
```bash
# Test stopping live presentation (need to create session first)
curl -X POST http://localhost:8090/api/presentations/{id}/stop \
  -H "Authorization: Bearer {auth_token}"

# Test stopping non-live presentation
curl -X POST http://localhost:8090/api/presentations/{id}/stop \
  -H "Authorization: Bearer {auth_token}"

# Test unauthorized access
curl -X POST http://localhost:8090/api/presentations/{id}/stop \
  -H "Authorization: Bearer {wrong_token}"

# Test unauthenticated access
curl -X POST http://localhost:8090/api/presentations/{id}/stop
```

**Commit**: "feat: implement POST /api/presentations/:id/stop endpoint"

### Step 4: Live Endpoint (POST /live)
**Goal**: Implement session creation - highest complexity and risk

**Tasks**:
1. Implement `generateAdminToken` function (copied from sync_sessions.go pattern)
2. Implement `handleStartLive` function with session creation logic
3. Add route registration for live endpoint
4. Test complete workflow: start → status → stop

**Files Modified**:
- `routes/presentations.go` (+60 lines)

**Verification**:
- POST /api/presentations/:id/live creates new sync session
- Updates presentation.active_session to link new session
- Returns session_id, admin_url, viewer_url, step metadata
- Requires authentication and ownership
- Returns 409 if presentation already live
- Generated URLs follow expected patterns

**Test Cases**:
```bash
# Test starting presentation live
curl -X POST http://localhost:8090/api/presentations/{id}/live \
  -H "Authorization: Bearer {auth_token}"

# Test starting already-live presentation
curl -X POST http://localhost:8090/api/presentations/{id}/live \
  -H "Authorization: Bearer {auth_token}"

# Test unauthorized start
curl -X POST http://localhost:8090/api/presentations/{id}/live \
  -H "Authorization: Bearer {wrong_token}"

# Test complete workflow
# 1. Start live
# 2. Check status (should show is_live: true)
# 3. Stop live
# 4. Check status (should show is_live: false)
```

**Commit**: "feat: implement POST /api/presentations/:id/live endpoint"

### Step 5: Integration Testing
**Goal**: Verify complete system integration and frontend compatibility

**Tasks**:
1. Test full presentation lifecycle workflow
2. Verify frontend `/watch/{id}` page works with new APIs
3. Test error handling across all endpoints
4. Verify no regression in existing functionality
5. Test concurrent operations (multiple users, rapid start/stop)

**Files Modified**:
- None (testing only)

**Verification**:
- Complete workflow: create presentation → start live → check status → stop live
- Frontend page correctly shows waiting → live → waiting states
- Error conditions return consistent, actionable messages
- Existing sync session routes continue to work
- Multiple presentations can be live simultaneously

**Test Scenarios**:
1. **Happy Path**: Complete lifecycle with valid auth
2. **Auth Failures**: Various authentication/authorization failures
3. **State Conflicts**: Attempt invalid state transitions
4. **Frontend Integration**: Browser-based testing of /watch pages
5. **Concurrent Access**: Multiple users managing same presentation

**Commit**: "test: verify presentation API routes integration and compatibility"

### Step 6: Documentation and Cleanup
**Goal**: Final polish and documentation

**Tasks**:
1. Add code comments and documentation strings
2. Clean up any TODO comments or temporary code
3. Verify error messages are user-friendly and consistent
4. Add any missing error handling

**Files Modified**:
- `routes/presentations.go` (documentation improvements)

**Verification**:
- Code is well-documented and self-explaining
- Error messages are clear and actionable
- No technical debt or temporary workarounds remain
- Code follows Go conventions and project patterns

**Commit**: "docs: add documentation and finalize presentation routes"

## Testing Checklist

### Unit Testing (Manual)
- [ ] `progressToStep` function works with various inputs
- [ ] `stepToProgress` function works with various inputs
- [ ] `buildStatusResponse` creates correct response structure
- [ ] `buildStartLiveResponse` creates correct response structure
- [ ] Authentication helpers work with valid/invalid auth

### Integration Testing
- [ ] GET /status works with existing presentations
- [ ] GET /status returns 404 for non-existent presentations
- [ ] POST /stop requires authentication
- [ ] POST /stop requires ownership
- [ ] POST /stop works only on live presentations
- [ ] POST /live creates sessions and links presentations
- [ ] POST /live requires authentication and ownership
- [ ] POST /live prevents duplicate live sessions

### End-to-End Testing
- [ ] Complete workflow: start → status → stop
- [ ] Frontend `/watch/{id}` shows correct states
- [ ] Multiple presentations can be live simultaneously
- [ ] Error conditions return appropriate HTTP status codes
- [ ] Existing sync session functionality unchanged

### Edge Case Testing
- [ ] Presentation with 1 step (progress calculation)
- [ ] Presentation with large step count
- [ ] Rapid start/stop operations
- [ ] Invalid presentation IDs
- [ ] Malformed auth tokens
- [ ] Network timeouts/interruptions

## Rollback Strategy

### Step-by-Step Rollback
Each step is independently committable, allowing selective rollback:

1. **Step 6 Issues**: Remove documentation, keep functionality
2. **Step 5 Issues**: Fix integration problems without removing endpoints
3. **Step 4 Issues**: Remove /live endpoint, keep status and stop working
4. **Step 3 Issues**: Remove /stop endpoint, keep status working
5. **Step 2 Issues**: Remove route registration, keep utility functions
6. **Step 1 Issues**: Delete entire file, no system impact

### Emergency Rollback
If critical issues discovered:
1. Comment out `registerPresentations(e)` in routes.go
2. System reverts to pre-implementation state
3. All existing functionality preserved

## Risk Mitigation

### High-Risk Areas
1. **Session Creation Logic**: Complex multi-step operation with potential for orphaned sessions
2. **Authentication Integration**: PocketBase auth context extraction could fail
3. **Concurrent Access**: Multiple users starting same presentation simultaneously

### Mitigation Strategies
1. **Session Creation**: Test creation/linking as separate operations, implement validation
2. **Authentication**: Extensive testing of auth context extraction and error handling
3. **Concurrent Access**: Test rapid operations, document expected behavior

### Monitoring Points
- Database consistency (presentations.active_session links to valid sessions)
- Orphaned session detection (sessions not linked to any presentation)
- Authentication failure rates (potential API misuse)

## Success Criteria

### Functional Requirements Met
- [ ] POST /api/presentations/:id/live creates sessions and returns URLs
- [ ] POST /api/presentations/:id/stop clears active sessions
- [ ] GET /api/presentations/:id/status returns enhanced presentation data
- [ ] All endpoints return 404 for non-existent presentations
- [ ] Authentication and authorization work correctly

### Technical Requirements Met
- [ ] Code follows existing patterns and conventions
- [ ] Error handling is consistent with existing routes
- [ ] No regression in existing functionality
- [ ] Performance is acceptable for expected load
- [ ] Code is maintainable and well-documented

### Integration Requirements Met
- [ ] Frontend compatibility maintained
- [ ] URL patterns work as expected
- [ ] Response formats match acceptance criteria
- [ ] PocketBase collection APIs used correctly
- [ ] Real-time updates continue to work

This plan provides a clear path from start to finish with built-in verification points and risk mitigation strategies.