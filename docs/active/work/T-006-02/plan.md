# T-006-02: sync-api-routes - Plan

## Overview

This document sequences the implementation steps for replacing the existing sync session API endpoints with new ones that match the acceptance criteria. The plan follows an incremental approach to minimize risk while ensuring exact compliance with specifications.

## Implementation Steps

### Step 1: Update Route Definitions
**Objective**: Change HTTP route paths to match acceptance criteria
**Duration**: 5 minutes
**Verification**: Routes registered at correct paths

**Actions**:
1. Modify `registerSyncSessions()` function in `routes/sync_sessions.go`
2. Change `POST "/api/sync-sessions"` to `POST "/api/sync/create"`
3. Change `PATCH "/api/sync-sessions/:id"` to `POST "/api/sync/:id/progress"`

**Verification Command**:
```bash
go run main.go &
curl -X POST http://localhost:8090/api/sync/create
curl -X POST http://localhost:8090/api/sync/test-id/progress
pkill -f "go run main.go"
```

**Success Criteria**:
- No 404 errors for the new routes
- Application starts without route registration errors

### Step 2: Update Request/Response Models
**Objective**: Change response format to match acceptance criteria
**Duration**: 10 minutes
**Verification**: Response structure matches specification

**Actions**:
1. Modify `CreateSessionResponse` struct:
   - Change `ID` field to `SessionID` with json tag `"session_id"`
   - Remove `Progress` field from create response
2. Keep `UpdateProgressRequest` struct unchanged (already correct)

**Verification**:
- Create session response contains `session_id`, `admin_url`, `viewer_url`
- Create session response does not contain `progress` field
- Update progress request accepts `{"progress": 0.5}` format

**Test Commands**:
```bash
# Create session and check response format
curl -X POST http://localhost:8090/api/sync/create | jq '.session_id'
curl -X POST http://localhost:8090/api/sync/create | jq '.admin_url'
curl -X POST http://localhost:8090/api/sync/create | jq '.viewer_url'
```

### Step 3: Update Authentication Method
**Objective**: Change from header-based to query parameter authentication
**Duration**: 15 minutes
**Verification**: Token validation works with query parameters

**Actions**:
1. Modify `handleUpdateProgress()` function:
   - Remove `X-Admin-Token` header extraction
   - Add query parameter extraction: `e.Request.URL.Query().Get("token")`
   - Update error message from "Missing admin token" to account for query params

**Verification**:
- Progress update with `?token=valid_token` succeeds
- Progress update without token returns 403/401
- Progress update with invalid token returns 403

**Test Commands**:
```bash
# Create session to get token
SESSION=$(curl -s -X POST http://localhost:8090/api/sync/create)
ID=$(echo $SESSION | jq -r '.session_id')
TOKEN=$(echo $SESSION | jq -r '.admin_url' | grep -o 'token=[^&]*' | cut -d= -f2)

# Test token validation
curl -X POST "http://localhost:8090/api/sync/$ID/progress?token=$TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"progress": 0.5}'
```

### Step 4: Update URL Format in Responses
**Objective**: Change admin and viewer URL formats to match specification
**Duration**: 10 minutes
**Verification**: URLs match acceptance criteria format

**Actions**:
1. Modify `handleCreateSession()` response building:
   - Change admin URL from `/admin/%s?token=%s` to `/sync/%s/control?token=%s`
   - Change viewer URL from `/viewer/%s` to `/sync/%s`

**Verification**:
- Admin URL format: `/sync/{id}/control?token={token}`
- Viewer URL format: `/sync/{id}`

**Test Commands**:
```bash
curl -X POST http://localhost:8090/api/sync/create | jq '.admin_url' | grep "/sync/.*/control?token="
curl -X POST http://localhost:8090/api/sync/create | jq '.viewer_url' | grep "/sync/"
```

### Step 5: Update Error Handling and Status Codes
**Objective**: Ensure error responses match acceptance criteria
**Duration**: 10 minutes
**Verification**: Correct HTTP status codes for all error conditions

**Actions**:
1. Verify status codes in `handleUpdateProgress()`:
   - 403 for invalid/missing token
   - 404 for session not found
   - 400 for invalid progress values
   - 200 for successful update
2. Update error messages to be consistent with API design

**Verification**:
- Missing token returns 403 or 401 (acceptance criteria allows either)
- Invalid token returns 403
- Non-existent session returns 404
- Progress < 0 or > 1 returns 400
- Valid request returns 200

**Test Commands**:
```bash
# Test error conditions
curl -w "%{http_code}" -X POST "http://localhost:8090/api/sync/invalid-id/progress" -d '{"progress": 0.5}'  # Should be 404
curl -w "%{http_code}" -X POST "http://localhost:8090/api/sync/$ID/progress?token=invalid" -d '{"progress": 0.5}'  # Should be 403
curl -w "%{http_code}" -X POST "http://localhost:8090/api/sync/$ID/progress?token=$TOKEN" -d '{"progress": 1.5}'  # Should be 400
```

### Step 6: Integration Testing
**Objective**: Verify end-to-end functionality with PocketBase realtime
**Duration**: 15 minutes
**Verification**: Complete workflow functions correctly

**Actions**:
1. Test complete session lifecycle:
   - Create session
   - Update progress multiple times
   - Verify realtime broadcasts (if possible to test)
2. Test boundary conditions:
   - Progress values at 0.0 and 1.0
   - Token validation edge cases
   - Concurrent progress updates

**Verification**:
- Session creation → progress update → success workflow
- Progress broadcasts trigger PocketBase realtime (collection updates visible)
- All acceptance criteria satisfied

**Test Workflow**:
```bash
# Complete integration test
SESSION=$(curl -s -X POST http://localhost:8090/api/sync/create)
echo "Created session: $SESSION"

ID=$(echo $SESSION | jq -r '.session_id')
TOKEN=$(echo $SESSION | jq -r '.admin_url' | sed 's/.*token=\([^&]*\).*/\1/')

echo "Testing progress updates..."
curl -X POST "http://localhost:8090/api/sync/$ID/progress?token=$TOKEN" \
     -H "Content-Type: application/json" -d '{"progress": 0.0}'

curl -X POST "http://localhost:8090/api/sync/$ID/progress?token=$TOKEN" \
     -H "Content-Type: application/json" -d '{"progress": 0.5}'

curl -X POST "http://localhost:8090/api/sync/$ID/progress?token=$TOKEN" \
     -H "Content-Type: application/json" -d '{"progress": 1.0}'

echo "Integration test complete"
```

## Testing Strategy

### Unit Tests
**Scope**: Individual function validation
**Location**: Inline tests or separate test file (if project uses testing)

**Test Cases**:
1. `generateAdminToken()`:
   - Returns 64-character hex string
   - Different calls return different tokens
   - No predictable patterns

2. `validateToken()`:
   - Identical strings return true
   - Different strings return false
   - Different length strings return false
   - Constant-time execution (not easily testable)

3. `validateProgress()`:
   - Values 0.0 and 1.0 are valid
   - Values < 0.0 return error
   - Values > 1.0 return error
   - Values between 0.0-1.0 are valid

### Integration Tests
**Scope**: HTTP endpoint behavior
**Method**: Manual testing with curl (automated if testing framework available)

**Test Scenarios**:
1. **Happy Path**:
   - Create session → get valid response format
   - Update progress with valid token → success
   - Multiple progress updates → all succeed

2. **Error Conditions**:
   - Update non-existent session → 404
   - Update with invalid token → 403
   - Update with no token → 403/401
   - Update with invalid progress → 400

3. **Edge Cases**:
   - Progress exactly 0.0 → success
   - Progress exactly 1.0 → success
   - Very long session IDs → handled gracefully
   - Malformed JSON in progress update → 400

### Security Tests
**Scope**: Authentication and authorization validation
**Method**: Manual testing with various token scenarios

**Test Cases**:
1. **Token Security**:
   - Token brute force resistance (large token space)
   - Timing attack resistance (constant-time comparison)
   - Token uniqueness across sessions

2. **Input Validation**:
   - SQL injection attempts in session ID
   - XSS attempts in JSON payloads
   - Buffer overflow attempts with large payloads

### Performance Tests
**Scope**: Basic performance characteristics
**Method**: Simple load testing if needed

**Metrics**:
- Response time for session creation
- Response time for progress updates
- Memory usage with multiple concurrent sessions
- Database query performance

## Rollback Strategy

### Rollback Triggers
- Breaking changes affect frontend functionality
- Security vulnerabilities discovered
- Performance degradation observed
- Integration test failures

### Rollback Process
1. **Immediate**: Revert `routes/sync_sessions.go` to previous version
2. **Verification**: Ensure old endpoints function correctly
3. **Communication**: Notify of rollback and timeline for fixes

### Rollback Testing
```bash
# Verify old endpoints work after rollback
curl -X POST http://localhost:8090/api/sync-sessions  # Should work
curl -X PATCH http://localhost:8090/api/sync-sessions/test-id \
     -H "X-Admin-Token: test-token" \
     -d '{"progress": 0.5}'  # Should work
```

## Risk Mitigation

### High-Risk Changes
1. **URL Path Changes**: Breaking change for any existing clients
   - **Mitigation**: Coordinate with frontend team on deployment timing
   - **Detection**: Monitor 404 errors on old endpoints post-deployment

2. **Authentication Method Change**: Header vs query parameter
   - **Mitigation**: Thoroughly test token extraction logic
   - **Detection**: Monitor 403 errors for authentication failures

### Medium-Risk Changes
1. **Response Format Changes**: Field name changes may break clients
   - **Mitigation**: Verify frontend expects new field names
   - **Detection**: Monitor client-side JavaScript errors

2. **Database Operations**: Same collection, different access patterns
   - **Mitigation**: Test with existing data, verify no data corruption
   - **Detection**: Monitor database errors and data integrity

### Low-Risk Changes
1. **Utility Function Reuse**: Proven security functions unchanged
2. **PocketBase Integration**: Using established patterns
3. **Error Handling**: Similar to existing error patterns

## Success Metrics

### Functional Success
- All acceptance criteria satisfied
- No regression in existing functionality
- Complete session lifecycle works end-to-end

### Technical Success
- No increase in response times
- No increase in error rates
- Proper error status codes returned
- Security properties maintained

### Integration Success
- Frontend compatibility maintained (with coordination)
- PocketBase realtime continues functioning
- Database integrity preserved

## Deployment Checklist

- [ ] All implementation steps completed
- [ ] Integration tests pass
- [ ] Security validation complete
- [ ] Performance baseline established
- [ ] Frontend team notified of changes
- [ ] Rollback procedure tested
- [ ] Documentation updated (if applicable)

This plan provides a systematic approach to implementing the sync API route changes while maintaining system stability and security. Each step can be completed and verified independently, allowing for incremental progress and early detection of issues.