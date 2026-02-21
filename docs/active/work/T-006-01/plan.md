# Plan: sync-sessions-collection

## Implementation Steps

### Step 1: Add sync_sessions collection to migrations
**File**: `migrations/collections.go`
**Action**: Add `createSyncSessionsCollection` function and call it from `createCollections`

**Implementation**:
1. Define the new function after existing collection functions
2. Check if sync_sessions collection exists
3. Create collection with progress and admin_token fields
4. Set appropriate field validations (min/max for progress)
5. Configure API rules for public read access
6. Add call to new function in main `createCollections` function

**Verification**:
- Start the app with fresh database
- Check PocketBase admin UI shows sync_sessions collection
- Verify fields are correctly configured

**Commit**: "feat: add sync_sessions collection migration"

### Step 2: Create sync sessions route handler
**File**: `routes/sync_sessions.go` (new file)
**Action**: Implement custom API routes for session management

**Implementation**:
1. Create new file with package declaration and imports
2. Define response/request structs
3. Implement `generateAdminToken()` function using crypto/rand
4. Implement `handleCreateSession()`:
   - Generate admin token
   - Create session record via PocketBase DAO
   - Return URLs and session ID
5. Implement `handleUpdateProgress()`:
   - Extract and validate X-Admin-Token header
   - Find session by ID
   - Compare tokens
   - Update progress if valid
6. Implement `registerSyncSessions()` to register routes

**Verification**:
- File compiles without errors
- Functions have proper error handling

**Commit**: "feat: implement sync sessions API routes"

### Step 3: Register sync sessions routes
**File**: `routes/routes.go`
**Action**: Add registration call for sync sessions routes

**Implementation**:
1. Add call to `registerSyncSessions(e)` after `registerHealth(e)`
2. Ensure proper import if needed

**Verification**:
- App starts without errors
- Routes are accessible

**Commit**: "feat: register sync sessions routes"

### Step 4: Test session creation
**Action**: Manual test of POST /api/sync-sessions endpoint

**Test Steps**:
1. Start the PocketBase server
2. Use curl or Postman to POST to /api/sync-sessions
3. Verify response includes:
   - Session ID
   - Admin and viewer URLs
   - Initial progress of 0
4. Check database has record with admin_token

**Expected Result**:
- 201 Created status
- Valid JSON response with all required fields
- Token stored in database

**Verification Command**:
```bash
curl -X POST http://localhost:8090/api/sync-sessions
```

### Step 5: Test progress update with valid token
**Action**: Test PATCH endpoint with correct admin token

**Test Steps**:
1. Create a session (from Step 4)
2. Extract admin token from database or creation response
3. PATCH with valid token and progress value
4. Verify progress updated in database

**Expected Result**:
- 200 OK status
- Progress value updated
- PocketBase triggers realtime update

**Verification Command**:
```bash
curl -X PATCH http://localhost:8090/api/sync-sessions/{id} \
  -H "X-Admin-Token: {token}" \
  -H "Content-Type: application/json" \
  -d '{"progress": 0.5}'
```

### Step 6: Test access control
**Action**: Verify security boundaries

**Test Cases**:
1. Try update without token - expect 401
2. Try update with wrong token - expect 401
3. Try direct collection API update - expect failure
4. Verify GET doesn't expose admin_token

**Verification Commands**:
```bash
# Without token
curl -X PATCH http://localhost:8090/api/sync-sessions/{id} \
  -d '{"progress": 0.5}'

# View session (shouldn't see token)
curl http://localhost:8090/api/collections/sync_sessions/records/{id}
```

### Step 7: Test realtime subscriptions
**Action**: Verify PocketBase SSE works with collection

**Test Steps**:
1. Open PocketBase admin UI
2. Navigate to sync_sessions collection
3. Create session via API
4. Update progress via API
5. Verify admin UI shows real-time updates

**Expected Result**:
- Changes appear immediately in admin UI
- SSE connection established automatically

### Step 8: Integration test
**Action**: Full end-to-end test of the feature

**Test Scenario**:
1. Create new session
2. Save admin URL and viewer URL
3. Update progress multiple times
4. Verify each update succeeds
5. Attempt invalid operations (bad token, out of range progress)

**Success Criteria**:
- All valid operations succeed
- All invalid operations fail appropriately
- No data corruption or leaks

**Commit**: "test: verify sync sessions implementation"

## Testing Strategy

### Unit Tests (Future Enhancement)
- Token generation produces correct length
- Progress validation rejects out-of-range values
- Token comparison is timing-safe

### Integration Tests
- Session creation returns expected structure
- Updates work with valid tokens
- Updates fail with invalid tokens
- Viewer endpoint doesn't expose tokens

### Manual Testing Checklist
- [ ] Collection created on startup
- [ ] POST /api/sync-sessions creates session
- [ ] PATCH with valid token updates progress
- [ ] PATCH without token returns 401
- [ ] PATCH with wrong token returns 401
- [ ] GET doesn't expose admin_token
- [ ] Progress must be between 0 and 1
- [ ] Realtime subscriptions work

## Error Handling

Each step includes proper error handling:
- Database errors return 500
- Validation errors return 400
- Auth errors return 401
- Not found errors return 404

## Rollback Plan

If issues arise:
1. Remove route registration from routes.go
2. Comment out collection creation in migrations
3. Delete sync_sessions collection from database
4. Revert commits in reverse order

## Dependencies

- Step 1 can run independently
- Steps 2-3 must be done together for compilation
- Steps 4-8 are testing/verification only

## Time Estimates

- Step 1: 10 minutes (migration)
- Step 2: 30 minutes (route implementation)
- Step 3: 5 minutes (registration)
- Steps 4-8: 20 minutes (testing)
- Total: ~65 minutes

## Completion Criteria

The ticket is complete when:
1. Collection exists with correct schema
2. API routes are functional
3. Admin token validation works
4. Viewers can access without authentication
5. Realtime subscriptions function
6. All acceptance criteria are met