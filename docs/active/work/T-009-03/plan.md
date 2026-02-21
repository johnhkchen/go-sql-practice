# T-009-03 Plan: Fix PathValue Routing

## Implementation Sequence

This document breaks down the implementation into ordered, atomic steps that can be executed and verified independently. Each step includes verification criteria and commit-worthy units of work.

## Step 1: Create Path Utilities Module
**Duration**: 15 minutes
**Files**: `routes/path_utils.go` (new)
**Commit**: "Add path parameter extraction utilities for PocketBase routing"

### Implementation
1. Create `routes/path_utils.go` with two utility functions:
   ```go
   func extractPathParam(path, segment string) string
   func extractLastPathParam(path, beforeSegment string) string
   ```

2. Implement robust path parsing logic:
   - Split path by "/"
   - Find target segment index
   - Return next segment as parameter value
   - Handle edge cases (empty path, missing segments, malformed URLs)

3. Add comprehensive documentation with examples

### Verification
- [ ] File compiles without errors
- [ ] Functions handle valid paths correctly: `/api/presentations/abc123/status` → `"abc123"`
- [ ] Functions handle invalid paths gracefully: `/api/presentations//status` → `""`
- [ ] Functions handle missing segments: `/api/presentations/status` → `""`
- [ ] Functions handle empty input: `""` → `""`

### Test Command
```bash
# Verify compilation
go build ./routes/...
```

## Step 2: Add Unit Tests for Utilities
**Duration**: 20 minutes
**Files**: `routes/path_utils_test.go` (new)
**Commit**: "Add comprehensive tests for path parameter extraction utilities"

### Implementation
1. Create test file with comprehensive test cases:
   - `TestExtractPathParam_ValidPaths`
   - `TestExtractPathParam_InvalidPaths`
   - `TestExtractPathParam_EdgeCases`
   - `TestExtractLastPathParam_ValidPaths`
   - `TestExtractLastPathParam_EdgeCases`

2. Test all scenarios identified in structure document:
   - Valid presentation paths: `/api/presentations/abc123/status`
   - Valid sync paths: `/api/sync/def456/progress`
   - Valid links paths: `/api/links/ghi789/view`
   - Empty segments: `/api/presentations//status`
   - Missing parameters: `/api/presentations/status`
   - Malformed paths: `/api//presentations/abc123`

### Verification
- [ ] All tests pass
- [ ] Test coverage includes edge cases
- [ ] Tests document expected behavior clearly

### Test Command
```bash
go test ./routes/... -v
```

## Step 3: Fix Presentations Routing
**Duration**: 10 minutes
**Files**: `routes/presentations.go` (modify)
**Commit**: "Fix presentations routing by replacing PathValue with manual parsing"

### Implementation
1. Update three handler functions:
   - `handleGetStatus()` line ~185
   - `handleStopLive()` line ~219
   - `handleStartLive()` line ~277

2. Replace in each function:
   ```go
   // Old:
   presentationID := e.Request.PathValue("id")

   // New:
   presentationID := extractPathParam(e.Request.URL.Path, "presentations")
   ```

3. Verify error handling remains unchanged

### Verification
- [ ] File compiles without errors
- [ ] No changes to function signatures
- [ ] No changes to error response formats
- [ ] All three handlers updated consistently

### Test Command
```bash
go build ./routes/...
# Manual verification: all presentation routes should work with proper IDs
```

## Step 4: Fix Sync Sessions Routing
**Duration**: 5 minutes
**Files**: `routes/sync_sessions.go` (modify)
**Commit**: "Fix sync sessions routing by replacing PathValue with manual parsing"

### Implementation
1. Update `handleUpdateProgress()` function at line ~88:
   ```go
   // Old:
   sessionID := e.Request.PathValue("id")

   // New:
   sessionID := extractPathParam(e.Request.URL.Path, "sync")
   ```

2. Verify error handling remains unchanged

### Verification
- [ ] File compiles without errors
- [ ] Function signature unchanged
- [ ] Error response format unchanged
- [ ] Parameter extraction works for sync paths

### Test Command
```bash
go build ./routes/...
# Manual verification: sync progress endpoint should work with proper session IDs
```

## Step 5: Integration Testing
**Duration**: 15 minutes
**Files**: None (testing only)
**Commit**: None (verification step)

### Testing Strategy
1. **Presentation Endpoints**:
   - GET `/api/presentations/{valid-id}/status` → 200 OK or 404 Not Found
   - POST `/api/presentations/{valid-id}/live` → depends on auth/state
   - POST `/api/presentations/{valid-id}/stop` → depends on auth/state
   - GET `/api/presentations/{invalid-id}/status` → 400 Bad Request

2. **Sync Endpoints**:
   - POST `/api/sync/{valid-id}/progress` → depends on auth/token
   - POST `/api/sync/{invalid-id}/progress` → 400 Bad Request

3. **Error Cases**:
   - POST `/api/presentations//status` → 400 Bad Request (empty ID)
   - POST `/api/sync//progress` → 400 Bad Request (empty ID)

### Verification
- [ ] Valid IDs extract correctly
- [ ] Invalid IDs result in appropriate 400 errors
- [ ] Empty IDs result in appropriate 400 errors
- [ ] No 500 server errors from routing issues
- [ ] Response formats unchanged from previous behavior

### Test Commands
```bash
# If test server is available:
curl -X GET "http://localhost:8090/api/presentations/test123/status"
curl -X POST "http://localhost:8090/api/sync/test456/progress"
```

## Step 6: Refactor Links View (Optional)
**Duration**: 10 minutes
**Files**: `routes/links_view.go` (modify)
**Commit**: "Refactor links view to use shared path parameter utilities"

### Implementation
1. Replace manual path parsing in `handleLinksView()` (lines 22-32):
   ```go
   // Old: Manual parsing with for loop
   path := e.Request.URL.Path
   parts := strings.Split(path, "/")
   // ... for loop logic

   // New: Use utility
   linkId := extractPathParam(e.Request.URL.Path, "links")
   ```

2. Remove unused variables and simplify logic
3. Maintain identical error handling behavior

### Verification
- [ ] File compiles without errors
- [ ] Identical behavior to previous implementation
- [ ] Code is cleaner and more consistent
- [ ] Links view endpoint continues to work

### Test Command
```bash
go build ./routes/...
# Manual verification: POST /api/links/{id}/view should work identically
```

## Step 7: Cleanup and Final Verification
**Duration**: 5 minutes
**Files**: `test_pathvalue.go` (delete)
**Commit**: "Remove temporary test file from PathValue investigation"

### Implementation
1. Delete temporary test file created during investigation
2. Verify no references to the deleted file exist
3. Run final comprehensive test

### Verification
- [ ] Temporary file removed
- [ ] All routes compile successfully
- [ ] All affected endpoints respond correctly
- [ ] No regressions in existing functionality

### Test Command
```bash
rm test_pathvalue.go
go build ./...
go test ./routes/... -v
```

## Rollback Strategy

Each step is designed to be atomic and reversible:

### Step-by-Step Rollback
1. **Step 7**: `git checkout HEAD~1 -- test_pathvalue.go`
2. **Step 6**: `git checkout HEAD~1 -- routes/links_view.go`
3. **Step 5**: No files to rollback (testing only)
4. **Step 4**: `git checkout HEAD~1 -- routes/sync_sessions.go`
5. **Step 3**: `git checkout HEAD~1 -- routes/presentations.go`
6. **Step 2**: `rm routes/path_utils_test.go`
7. **Step 1**: `rm routes/path_utils.go`

### Complete Rollback
```bash
git reset --hard HEAD~4  # Assuming 4 commits were made
```

## Risk Mitigation

### Compilation Safety
- Each step verifies compilation before proceeding
- New utilities have no external dependencies
- Modified functions maintain identical signatures

### Runtime Safety
- Utility functions handle all edge cases gracefully
- Error handling behavior preserved in all handlers
- No new failure modes introduced

### Testing Safety
- Unit tests validate utility behavior before integration
- Integration testing verifies end-to-end functionality
- Manual verification confirms no regressions

## Success Metrics

### Functional Requirements
- [ ] `GET /api/presentations/{id}/status` extracts `{id}` correctly
- [ ] `POST /api/presentations/{id}/live` extracts `{id}` correctly
- [ ] `POST /api/presentations/{id}/stop` extracts `{id}` correctly
- [ ] `POST /api/sync/{id}/progress` extracts `{id}` correctly

### Quality Requirements
- [ ] All new code has unit test coverage
- [ ] No regressions in existing functionality
- [ ] Error handling maintains existing behavior
- [ ] Code follows existing patterns and style

### Performance Requirements
- [ ] Path parsing overhead is negligible (< 1ms per request)
- [ ] No memory leaks from string operations
- [ ] Response times unchanged from baseline

## Dependencies and Blockers

### Internal Dependencies
- None - all changes are self-contained within routes package

### External Dependencies
- Standard library `strings` package (already in use)
- PocketBase core types (already in use)

### Potential Blockers
- None anticipated - solution uses only standard Go features
- Manual testing may require running application locally

## Estimated Timeline

**Total Duration**: 80 minutes (1 hour 20 minutes)
- Step 1: 15 minutes
- Step 2: 20 minutes
- Step 3: 10 minutes
- Step 4: 5 minutes
- Step 5: 15 minutes
- Step 6: 10 minutes (optional)
- Step 7: 5 minutes

**Critical Path**: Steps 1-5 are essential; Steps 6-7 are cleanup/optimization.

**Parallelization**: None - steps must be executed sequentially due to dependencies.