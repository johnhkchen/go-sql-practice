# Design: Delete Dead Code and Duplicates

## Approach Options

### Option 1: Incremental Removal (Function by Function)
**Approach**: Remove individual dead functions first, then files
**Pros**:
- Easier rollback if issues arise
- Can verify each change independently
- Minimal risk of breaking dependencies

**Cons**:
- More commits required
- Slower overall process
- Temporary inconsistent state

### Option 2: Batch Removal (All at Once)
**Approach**: Delete all identified dead code in single operation
**Pros**:
- Clean, atomic change
- No intermediate states
- Single verification pass

**Cons**:
- Harder to isolate issues if build breaks
- Large diff for review
- All-or-nothing approach

### Option 3: Category-Based Removal
**Approach**: Group removals by type (files, functions, comments)
**Pros**:
- Logical grouping of changes
- Balanced risk/reward
- Clear progression

**Cons**:
- Still requires multiple passes
- Some interdependencies between categories

## Selected Approach: Category-Based Removal

### Rationale
1. **Logical Grouping**: Changes are easier to review when grouped by type
2. **Risk Management**: File deletions are highest impact, do them first
3. **Dependency Order**: Remove consumers before providers
4. **Testing Strategy**: Each category can be tested independently

### Removal Order

#### Phase 1: Delete Duplicate Files
- `routes/links_search_simple.go` (entire file)
- Update `routes/routes.go` to remove registration
- **Verification**: Build succeeds, `/api/links/search` works

#### Phase 2: Remove Dead Functions
- `validateToken()` from `presentations.go`
- `stepToProgress()` from `presentations.go`
- `TestData` struct from `routes_test.go`
- `assertErrorResponse()` from `routes_test.go`
- **Verification**: Tests pass, no compilation errors

#### Phase 3: Inline Single-Use Export
- Remove `FrontendExists()` from `internal/frontend/embed.go`
- Update `routes/routes.go` to check error directly
- **Verification**: Static serving still works

#### Phase 4: Clean Up Comments and Unreachable Code
- Remove commented import from `presentations.go`
- Remove stale comment from `presentations.go`
- Fix unreachable block in `static.go`
- **Verification**: Linting passes

#### Phase 5: Finalize Dependencies
- Run `go mod tidy`
- **Verification**: No changes to go.mod (all deps still used)

## Implementation Details

### Replacing FrontendExists() Call
Current code in `routes/routes.go:21-25`:
```go
if frontend.FrontendExists() {
    registerStatic(e)
} else {
    e.App.Logger().Warn("Frontend assets not found, static serving disabled")
}
```

New approach:
```go
registerStatic(e)
```

The `registerStatic()` function already handles the error internally with early return and logging. The conditional check is redundant.

### Fixing static.go Unreachable Code
Current problematic code (lines 73-79):
```go
// If not a ReadSeeker, fall back to copying the content
ev.Response.Header().Set("Content-Type", "application/octet-stream")
_, err = ev.Response.Write([]byte("file serving not supported"))
```

Better approach:
```go
// If not a ReadSeeker, continue to next handler
return ev.Next()
```

This maintains middleware chain integrity instead of returning an error message.

## Risk Assessment

### Low Risk Items
- Deleting `links_search_simple.go`: Duplicate functionality, no unique features
- Removing unused test helpers: No test dependencies on them
- Removing comments: Documentation only
- `go mod tidy`: Automated tool, safe operation

### Medium Risk Items
- Removing `FrontendExists()`: Requires understanding `registerStatic()` error handling
- Modifying `static.go`: Affects error path, but unlikely to execute

### Mitigations
1. **Testing**: Run full test suite after each phase
2. **Build Verification**: Compile after each file change
3. **Endpoint Testing**: Manually verify `/api/links/search` works
4. **Git Safety**: Each phase in separate commit for easy revert

## Success Criteria

### Functional
- All endpoints continue to work
- Tests pass without modification
- Build succeeds without warnings

### Code Quality
- No unused exports
- No unreachable code
- No commented imports
- Clean `go vet` output

### Metrics
- Lines removed: ~400 (191 from links_search_simple.go alone)
- Functions removed: 4
- Exports removed: 1
- Comments cleaned: 2

## Rejected Alternatives

### Alternative 1: Keep FrontendExists() as Public API
**Rejected because**: Single use, adds no value over direct error check

### Alternative 2: Keep validateToken() for Future Use
**Rejected because**: YAGNI principle, validateSyncToken() already exists for the same purpose

### Alternative 3: Fix links_search_simple.go Instead of Delete
**Rejected because**: Complete duplicate, maintenance burden, not in API spec

## Rollback Plan
Each phase creates atomic commits. If issues arise:
1. Identify problematic phase from git log
2. `git revert <commit>` for that phase
3. Re-evaluate and adjust approach
4. Continue with remaining phases