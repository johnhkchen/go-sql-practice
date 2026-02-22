# Plan: Delete Dead Code and Duplicates

## Implementation Steps

### Step 1: Delete Duplicate Link Search Implementation
**Files**: `routes/links_search_simple.go`, `routes/routes.go`
**Actions**:
1. Delete entire file `routes/links_search_simple.go`
2. Remove `registerLinksSearchSimple(e)` call from `routes/routes.go` line 17
**Verification**:
- `go build .` succeeds
- `go vet ./routes` passes
**Commit message**: "chore: remove duplicate links search implementation"

### Step 2: Remove Unused Functions from presentations.go
**File**: `routes/presentations.go`
**Actions**:
1. Delete `stepToProgress()` function (lines 106-111)
2. Delete `validateToken()` function (lines 173-179)
**Verification**:
- `go build .` succeeds
- `grep -r "stepToProgress\|validateToken" routes/` returns only references in docs
**Commit message**: "chore: remove unused functions from presentations"

### Step 3: Clean Up Comments and Imports in presentations.go
**File**: `routes/presentations.go`
**Actions**:
1. Remove commented import on line 7: `// "encoding/json"  // TODO: Re-enable if JSON operations are added`
2. Remove stale comment on line 181: `// Placeholder handler functions - to be implemented in subsequent steps`
**Verification**:
- `go vet ./routes` passes
- No commented imports remain
**Commit message**: "chore: clean up stale comments in presentations"

### Step 4: Remove Unused Test Helpers
**File**: `routes/routes_test.go`
**Actions**:
1. Delete `TestData` struct (lines 17-22)
2. Delete `assertErrorResponse()` function (lines 157-168)
**Verification**:
- `go test ./routes` passes
- All existing tests continue to work
**Commit message**: "chore: remove unused test helpers"

### Step 5: Simplify Frontend Existence Check
**Files**: `internal/frontend/embed.go`, `routes/routes.go`
**Actions**:
1. Delete `FrontendExists()` function from `internal/frontend/embed.go` (lines 20-23)
2. Replace conditional block in `routes/routes.go` (lines 21-25) with just `registerStatic(e)`
**Verification**:
- `go build .` succeeds
- Static file serving still works (manual test if frontend exists)
**Commit message**: "chore: remove unnecessary FrontendExists wrapper"

### Step 6: Fix Unreachable Code in static.go
**File**: `routes/static.go`
**Actions**:
1. Replace lines 73-79 (unreachable fallback) with `return ev.Next()`
**Verification**:
- `go build .` succeeds
- `go vet ./routes` passes
**Commit message**: "chore: fix unreachable code in static handler"

### Step 7: Clean Up Dependencies
**File**: `go.mod`, `go.sum`
**Actions**:
1. Run `go mod tidy`
**Verification**:
- Check `git diff go.mod` - expect no changes (all deps still used)
- `go build .` succeeds
**Commit message**: "chore: run go mod tidy"

## Testing Strategy

### Unit Test Verification
**Run after each step**:
```bash
go test ./routes -v
```

**Expected outcomes**:
- All tests pass
- No new failures introduced
- `TestMakeRequest_RealExecution` specifically validates routing still works

### Integration Test
**After Step 1 (links_search_simple removal)**:
```bash
# Test that the main search endpoint still works
curl -X GET "http://localhost:8080/api/links/search?q=test&page=1&perPage=10"
```

**After Step 5 (FrontendExists removal)**:
```bash
# Verify static files are still served (if frontend built)
curl -I http://localhost:8080/
```

### Linting Verification
**Run after each file modification**:
```bash
go vet ./...
golangci-lint run --disable-all --enable=deadcode,unused
```

## Rollback Procedures

### Per-Step Rollback
Each step creates an atomic commit. If issues arise:
```bash
# Identify problematic commit
git log --oneline -5

# Revert specific commit
git revert <commit-hash>

# Re-run verification
go build .
go test ./...
```

### Full Rollback
If multiple issues cascade:
```bash
# Reset to beginning of cleanup
git reset --hard <initial-commit>

# Re-evaluate approach
```

## Success Metrics

### Quantitative
- **Lines deleted**: ~415 total
  - links_search_simple.go: 191
  - Other files: ~224
- **Functions removed**: 5
- **Files deleted**: 1
- **Build time**: No regression
- **Test coverage**: Maintained or improved

### Qualitative
- Codebase clarity improved
- No duplicate implementations
- No unused exports
- Clean linter output

## Risk Mitigation

### High-Risk Steps
**Step 1** (Delete links_search_simple.go):
- Mitigation: Verify `/api/links/search` endpoint first
- Fallback: Can restore file if needed

**Step 5** (Remove FrontendExists):
- Mitigation: Understand registerStatic error handling first
- Fallback: Restore check if static serving fails

### Low-Risk Steps
Steps 2-4, 6-7 are low risk:
- Remove only unused code
- No functional changes
- Easy to verify

## Timeline

### Estimated Duration: 30 minutes
- Step 1: 5 minutes (delete file, update routes.go)
- Step 2: 5 minutes (remove two functions)
- Step 3: 2 minutes (remove comments)
- Step 4: 3 minutes (remove test helpers)
- Step 5: 5 minutes (remove export, update caller)
- Step 6: 3 minutes (fix unreachable code)
- Step 7: 2 minutes (go mod tidy)
- Verification: 5 minutes (run all tests)

## Post-Implementation

### Documentation
- No user-facing documentation changes needed
- Internal code is self-documenting after cleanup

### Monitoring
- Watch for any 404s on `/api/links/search-simple` (expected, as it wasn't in spec)
- Verify no regression in `/api/links/search` performance

### Follow-up
- Consider adding linter rule to prevent unused exports
- Document in team practices: regular dead code cleanup