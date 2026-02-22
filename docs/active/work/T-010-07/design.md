# Design: T-010-07 - Fix Go Build Failures

## Problem Restatement

Three import-related issues prevent the Go backend from compiling:
1. Unused import of `crypto/subtle` in `routes/presentations.go`
2. Unused import of `strings` in `routes/links_view.go`
3. Missing import of `github.com/pocketbase/pocketbase` in `routes/links_view_test.go`

## Solution Approaches

### Approach 1: Manual Import Cleanup
Remove the two unused imports and add the missing import. This is the most direct approach.

**Pros:**
- Minimal change footprint
- No risk of affecting functionality
- Easy to verify correctness

**Cons:**
- None for this specific case

### Approach 2: Use goimports Tool
Run `goimports` to automatically fix all import issues.

**Pros:**
- Automated solution
- Would fix all three issues at once

**Cons:**
- May reformat imports in ways that differ from project style
- Could potentially add or remove imports beyond the three identified issues
- Less control over exact changes

### Approach 3: Use go fmt with -s flag
Use `go fmt -s` to simplify the code and remove unused imports.

**Pros:**
- Standard Go tooling

**Cons:**
- Only removes unused imports, won't add missing ones
- May make other simplifications beyond import cleanup

## Decision: Approach 1 - Manual Import Cleanup

### Rationale

Manual import cleanup is the best approach because:

1. **Precision**: We know exactly what needs to be changed - two imports to remove, one to add
2. **Control**: No risk of unintended changes to code formatting or other imports
3. **Traceability**: Clear audit trail of exactly what was changed for this ticket
4. **Minimal Risk**: No possibility of breaking working code or changing behavior
5. **Project Conformance**: Maintains existing import organization and formatting

The automated approaches would work but introduce unnecessary uncertainty for such a simple fix.

## Implementation Details

### File: routes/presentations.go
- Remove line 4: `"crypto/subtle"`
- Maintain the remaining imports in their current order
- The blank line after imports should remain

### File: routes/links_view.go
- Remove line 5: `"strings"`
- Maintain the remaining imports in their current order
- The blank line after imports should remain

### File: routes/links_view_test.go
- Add import `github.com/pocketbase/pocketbase` to the import block
- Position it appropriately with other third-party imports
- Based on Go conventions, it should go after standard library imports

## Validation Strategy

After making the changes:
1. Run `go build .` - should succeed with no errors
2. Run `go vet ./...` - should pass with no import-related errors
3. Run `go test ./routes -c` - should compile without errors (tests may fail, but must compile)

## Risks and Mitigations

### Risk 1: Import Organization
**Risk**: The added import might not follow project import grouping conventions
**Mitigation**: Check other test files to see the import pattern and follow it

### Risk 2: Hidden Dependencies
**Risk**: Removing an import might break something not immediately obvious
**Mitigation**: The research phase confirmed these imports are truly unused through code analysis

## Rejected Alternatives

### Why not use goimports?
While goimports would fix the issues, it operates on all files it touches and may reformat import blocks differently than the project style. For three specific import changes, manual editing provides more control.

### Why not fix other test errors?
The ticket scope is specifically about import-related build failures from T-010-01 and T-010-02. Other test compilation errors (stats_test.go, routes_test.go) are out of scope and likely require separate tickets.

## Success Criteria

The implementation is successful when:
1. `go build .` completes without import errors
2. `go vet ./...` passes
3. `go test ./routes -c` compiles (even if tests fail at runtime)