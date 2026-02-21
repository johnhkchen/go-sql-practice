# Design for T-010-05: Infrastructure Cleanup

## Design Approaches

### Approach 1: Minimal Touch
Only fix the critical issues: add `.gitignore` entries, delete stray directories, fix Makefile help.
- **Pros**: Fast, low risk, minimal testing needed
- **Cons**: Leaves quality gaps, CI remains outdated, no lint/vet enforcement

### Approach 2: Full Enhancement
Implement all improvements: CI modernization, quality checks, error handling improvements.
- **Pros**: Comprehensive improvement, better long-term maintenance
- **Cons**: Higher risk of CI breakage, more testing required

### Approach 3: Incremental Safety (Selected)
Fix all issues but prioritize safety: careful CI updates, preserve working build, add quality checks without breaking existing flow.
- **Pros**: Balanced risk, maintains stability, progressive enhancement
- **Cons**: Slightly more complex than minimal approach

## Selected Approach: Incremental Safety

This approach addresses all issues while maintaining system stability. Each change is independently safe.

## Implementation Decisions

### .gitignore Updates
**Decision**: Add both `*.test` and `routes/pb_data/`
- `*.test` prevents all Go test binaries
- `routes/pb_data/` handles specific test artifact location
- Both patterns ensure comprehensive coverage

### Directory Cleanup
**Decision**: Delete both problematic directories
- `frontend/frontend/src/` is confirmed duplicate
- `routes/pb_data/` is test artifact
- Safe to delete without data loss

### Embed Path Resolution
**Decision**: Maintain current workaround, document it
- Current state works (files in both locations)
- Changing embed path risks breaking builds
- Add comment explaining the dual location requirement
- Future refactor can address this properly

**Rejected**: Modify embed.go path
- Risk of breaking production builds
- Would require extensive testing

### Makefile Improvements
**Decision**: Add quality targets, fix help text
- Add `lint` using `gofmt -l`
- Add `vet` using `go vet`
- Fix help text for validate-build
- Keep existing targets unchanged

**Rejected**: Add embed path copy step
- Current workaround is working
- Adding complexity without immediate benefit

### CI Workflow Updates
**Decision**: Careful version updates and new checks
- Update to `setup-go@v5` (well-tested action)
- Add vet and gofmt checks after successful build
- Add npm audit with high severity threshold
- Keep Go 1.26 version (don't change working config)

**Rejected**: Add govulncheck
- May not be available in Go 1.26
- Can be added later if needed

### Frontend Configuration
**Decision**: Create .env.example with same content as .env
- Simple duplication of working configuration
- Documents the required variable

### Error Handling
**Decision**: Add server-side logging before generic errors
- Use app.Logger() for actual error details
- Maintain generic user-facing messages
- Improves debugging without exposing internals

**Rejected**: Include error details in responses
- Security risk (information disclosure)
- Against best practices

### CLAUDE.md Update
**Decision**: Add factual one-line description
- "Link bookmarks app with Go/PocketBase backend and Astro frontend"
- Accurate and concise

## Risk Mitigation

1. **CI Breakage**: Test each CI change in sequence, revert if needed
2. **Build Issues**: Preserve working Makefile targets, only add new ones
3. **Data Loss**: Verify directories are truly duplicates before deletion
4. **Test Failures**: Run full test suite after each change group

## Change Groups

1. **File Cleanup**: .gitignore, directory deletion, routes.test removal
2. **Makefile**: New targets, help fix
3. **CI Pipeline**: Version update, quality checks
4. **Configuration**: .env.example, CLAUDE.md
5. **Error Handling**: Logging additions

## Validation Strategy

- Each group independently testable
- Existing functionality preserved
- New functionality is additive
- No breaking changes to external interfaces

## Rollback Plan

If issues arise:
1. Git revert the specific commit
2. Each change group is atomic
3. CI changes can be reverted via PR
4. Local changes testable before commit