# Plan Phase - T-010-02: Extract Shared Go Utilities

## Implementation Steps

### Step 1: Create Token Utilities File
**Action**: Create `routes/tokens.go` with shared token functions
**Verification**:
- File exists at `routes/tokens.go`
- Contains `TokenLength` constant
- Contains `GenerateToken()` function
- Contains `ValidateToken()` function
- `go vet routes/tokens.go` passes

### Step 2: Update presentations.go
**Action**: Modify `routes/presentations.go` to use shared token utilities
**Changes**:
1. Remove `TokenLength` constant (line 14)
2. Remove `generateAdminToken()` function (lines 84-90)
3. Update `handleStartLive()` to call `GenerateToken()` instead of `generateAdminToken()`
**Verification**:
- No duplicate token generation code remains
- `go build ./routes` succeeds
- `go test ./routes -run TestPresentation` passes (if tests exist)

### Step 3: Update sync_sessions.go Part 1 - Token Functions
**Action**: Replace token functions with shared utilities
**Changes**:
1. Remove `syncTokenLength` constant (line 15)
2. Remove `generateSyncAdminToken()` function (lines 150-156)
3. Remove `validateSyncToken()` function (lines 167-172)
4. Update `handleCreateSession()` to call `GenerateToken()`
5. Update `handleUpdateProgress()` to call `ValidateToken()`
**Verification**:
- No duplicate token code remains
- Code compiles

### Step 4: Fix Security - Remove Admin Token from Response
**Action**: Modify `sync_sessions.go` response in `handleUpdateProgress()`
**Changes**:
1. Remove `"admin_token"` field from response map (line 143)
2. Keep only `id`, `progress`, `created`, `updated` fields
**Verification**:
- Response no longer contains `admin_token`
- `go test ./routes -run TestSyncSession` passes

### Step 5: Fix HTTP Status Codes in stats.go
**Action**: Replace raw integers with http.Status* constants
**Changes**:
1. Add `"net/http"` to imports
2. Add `statsTopN = 5` constant after line 9
3. Replace all `500` with `http.StatusInternalServerError`
4. Replace `200` with `http.StatusOK`
**Verification**:
- No raw HTTP status codes remain
- `go build ./routes` succeeds

### Step 6: Fix HTTP Status Codes in links_view.go
**Action**: Replace raw integers with http.Status* constants
**Changes**:
1. Add `"net/http"` to imports
2. Replace `400` with `http.StatusBadRequest`
3. Replace `500` with `http.StatusInternalServerError`
4. Replace `404` with `http.StatusNotFound`
5. Replace `200` with `http.StatusOK`
**Verification**:
- No raw HTTP status codes remain
- Code compiles

### Step 7: Fix Error Handling in links_view.go
**Action**: Handle the ignored RowsAffected error
**Changes**:
1. Change `rowsAffected, _` to `rowsAffected, err`
2. Add error check and logging
3. Set `rowsAffected = 0` on error
**Verification**:
- Error is properly handled
- Logging call is present

### Step 8: Add Error Logging in links_search.go
**Action**: Add proper logging for tag fetch errors
**Changes**:
1. Add `app.Logger().Error("Failed to fetch tags", "error", err)` at line 91
**Verification**:
- Error logging is present
- Code still handles error gracefully

### Step 9: Run Tests
**Action**: Execute all test suites
**Commands**:
```bash
go test ./routes/...
go vet ./...
go build .
```
**Verification**:
- All tests pass
- No vet warnings
- Build succeeds

### Step 10: Manual Verification
**Action**: Test key endpoints manually
**Test Cases**:
1. Create presentation session - should generate token
2. Create sync session - should generate token
3. Update sync progress - should validate token, no admin_token in response
4. Get stats - should return data with proper status codes
5. View link - should increment count with proper status codes
6. Search links - should work with tag fetching
**Verification**:
- All endpoints work as expected
- No token exposed in responses
- Proper HTTP status codes returned

## Testing Strategy

### Unit Tests
- Existing tests should pass without modification
- Token generation maintains same output format
- Token validation maintains same behavior

### Integration Tests
The following flows need verification:
1. **Presentation Flow**: Start live → Get status → Stop live
2. **Sync Session Flow**: Create session → Update progress → Validate token
3. **Stats Flow**: Get stats with all metrics
4. **Links Flow**: Search → View → Increment count

### Regression Testing
Check for:
- Token length consistency (64 hex characters)
- Token uniqueness (no collisions in reasonable samples)
- Timing-safe comparison still works
- HTTP responses maintain same structure (except admin_token removal)

## Rollback Plan

If issues arise:
1. Revert all changes via git
2. Individual file rollback possible since changes are isolated
3. Can temporarily restore duplicate functions if needed

## Risk Assessment

**Low Risk Changes**:
- HTTP status code replacements (cosmetic)
- Adding error logging (additive)
- Adding documentation constant (comment-like)

**Medium Risk Changes**:
- Consolidating token functions (behavior must match exactly)
- Removing admin_token from response (API change)

**Mitigation**:
- Careful testing of token generation/validation
- Verify no client code depends on admin_token in response

## Commit Strategy

Suggested atomic commits:
1. "refactor: create shared token utilities in routes/tokens.go"
2. "refactor: use shared token functions in presentations and sync_sessions"
3. "fix: remove admin_token from sync progress response"
4. "refactor: replace raw HTTP status codes with constants"
5. "fix: add error handling and logging for swallowed errors"
6. "docs: document statsTopN constant for clarity"

Each commit should be independently buildable and testable.