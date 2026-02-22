# Plan: T-010-07 - Fix Go Build Failures

## Implementation Steps

### Step 1: Verify Current Build Failures
**Action**: Run go build to confirm the exact errors we're fixing
**Command**: `go build .`
**Expected Output**:
```
routes/links_view.go:5:2: "strings" imported and not used
routes/presentations.go:4:2: "crypto/subtle" imported and not used
```
**Verification**: Build fails with these two import errors

### Step 2: Remove Unused Import from presentations.go
**File**: `routes/presentations.go`
**Action**: Remove line 4 containing `"crypto/subtle"`
**Edit Operation**:
- Find: Line 4 with `"crypto/subtle"`
- Delete: The entire line
**Verification**: File should have one fewer line, imports should look clean

### Step 3: Remove Unused Import from links_view.go
**File**: `routes/links_view.go`
**Action**: Remove line 5 containing `"strings"`
**Edit Operation**:
- Find: Line 5 with `"strings"`
- Delete: The entire line
**Verification**: File should have one fewer line, imports should look clean

### Step 4: Verify Build Now Succeeds
**Action**: Run go build again
**Command**: `go build .`
**Expected**: Build should complete successfully
**Fallback**: If other errors appear, they should not be import-related

### Step 5: Check Test Compilation Error
**Action**: Attempt to compile tests to see the missing import error
**Command**: `go test ./routes -c`
**Expected Error**: `routes/links_view_test.go:201:32: undefined: pocketbase`
**Note**: Other test errors may appear but we focus on this one

### Step 6: Add Missing Import to links_view_test.go
**File**: `routes/links_view_test.go`
**Action**: Add `github.com/pocketbase/pocketbase` import
**Edit Operation**:
- Find: Import block with `"github.com/pocketbase/pocketbase/core"`
- Add: `"github.com/pocketbase/pocketbase"` on the line before it
**Verification**: Import block should have both pocketbase imports

### Step 7: Final Build Verification
**Action**: Run full build
**Command**: `go build .`
**Expected**: Success with no output
**Required**: Must pass before proceeding

### Step 8: Run go vet
**Action**: Run go vet to check for issues
**Command**: `go vet ./...`
**Expected**: No import-related errors
**Note**: Other warnings may exist but no import issues

### Step 9: Verify Test Compilation
**Action**: Compile the tests (not run them)
**Command**: `go test ./routes -c`
**Expected**: Compilation succeeds (may have other errors, but must compile)
**Success Criteria**: The undefined pocketbase error is gone

### Step 10: Clean Up Test Binary
**Action**: Remove the compiled test binary
**Command**: `rm routes.test`
**Reason**: Don't leave build artifacts in the repository

## Testing Strategy

### Unit Testing
Not applicable - we're only changing imports, not logic.

### Compilation Testing
Primary validation through successful compilation:
- `go build .` must succeed
- `go test ./routes -c` must compile

### Integration Testing
No integration tests needed - import changes don't affect runtime behavior.

## Rollback Plan

If any step fails unexpectedly:
1. Revert all three files to their original state
2. Re-examine the actual error messages
3. Adjust the plan based on actual vs expected errors

Simple rollback since we're only modifying three import sections.

## Commit Strategy

Single atomic commit after all changes:
- Message: "fix: remove unused imports and add missing pocketbase import"
- Files: 3 files changed
- Verification: Run build and vet before committing

## Risk Assessment

### Low Risk
- Changes are purely import-related
- No logic modifications
- Easy to verify correctness

### Potential Issues
1. **Import ordering**: Go has conventions but they're not strict
   - Mitigation: Follow existing file patterns
2. **Hidden usage**: Import might be used in build tags we don't see
   - Mitigation: Research phase confirmed these are truly unused
3. **Test dependencies**: Other test failures might cascade
   - Mitigation: Focus only on import-related compilation errors

## Success Criteria Checklist

- [ ] `go build .` succeeds
- [ ] `go vet ./...` passes with no import errors
- [ ] `go test ./routes -c` compiles successfully
- [ ] No functional code was modified
- [ ] All three import issues from ticket are resolved

## Time Estimate

Total time: ~5 minutes
- Step 1-3: 2 minutes (remove unused imports)
- Step 4-6: 2 minutes (add missing import)
- Step 7-10: 1 minute (verification)