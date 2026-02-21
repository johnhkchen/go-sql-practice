# T-009-01 Design: Fix Static File Serving

## Problem Analysis

The research phase revealed that static file serving is completely broken due to two intentional disabling mechanisms:
1. Early return in `registerStatic()` function (routes/static.go:36)
2. Commented out function call in main registration (routes/routes.go:20)

Both issues stem from unresolved PocketBase v0.36.5 compatibility concerns. The core implementation exists but has never been tested or enabled.

## Design Options

### Option 1: Enable Existing Implementation As-Is

Simply remove the early return and uncomment the registration call.

**Pros:**
- Minimal code changes
- Uses existing SPA filesystem wrapper
- Preserves original architectural intent

**Cons:**
- May not work with PocketBase v0.36.5 (untested compatibility)
- Could have hidden issues in Echo v5 middleware integration
- Risk of breaking existing API/admin routes

**Risk Assessment:** High - Blindly enabling untested code with known compatibility concerns

### Option 2: Research PocketBase v0.36.5 Static Serving Patterns

Investigate PocketBase v0.36.5 documentation and source to understand the correct static serving approach.

**Pros:**
- Ensures compatibility with current PocketBase version
- Follows official patterns and best practices
- Lower risk of breaking changes

**Cons:**
- Requires additional research phase
- May require significant code restructuring
- Could deviate from existing codebase patterns

**Risk Assessment:** Medium - Thorough but time-intensive

### Option 3: Incremental Fix with Validation

Remove the early return, uncomment registration, but add comprehensive error handling and validation.

**Pros:**
- Progressive approach with safety nets
- Can identify specific compatibility issues
- Maintains existing architecture while improving robustness
- Enables testing to validate functionality

**Cons:**
- More complex implementation
- Requires additional error handling code

**Risk Assessment:** Low - Safe, testable approach

### Option 4: Replace with Native PocketBase Static Serving

Use PocketBase's built-in static file serving capabilities instead of custom Echo middleware.

**Pros:**
- Guaranteed compatibility
- Leverages framework features
- Reduces custom code maintenance

**Cons:**
- May not support SPA routing requirements
- Could conflict with existing route structure
- Requires understanding PocketBase's static serving API

**Risk Assessment:** Medium - Framework-native but may lack SPA support

## Decision Matrix

| Criteria | Option 1 | Option 2 | Option 3 | Option 4 |
|----------|----------|----------|----------|----------|
| Implementation Speed | High | Low | Medium | Medium |
| Compatibility Risk | High | Low | Low | Medium |
| Code Complexity | Low | Medium | Medium | Low |
| Testing Ease | Medium | High | High | Medium |
| Maintainability | Medium | High | High | High |

## Chosen Solution: Option 3 - Incremental Fix with Validation

**Rationale:**

Option 3 provides the best balance of safety, speed, and maintainability. Here's why:

1. **Preserves Existing Architecture**: The current implementation shows good understanding of SPA routing requirements (fallback to index.html, API/admin path protection). Rather than discarding this work, we build upon it.

2. **Addresses Core Issues Safely**: By removing the early return and adding proper error handling, we can identify exactly what compatibility issues exist with PocketBase v0.36.5.

3. **Enables Incremental Testing**: We can validate each component (filesystem access, middleware registration, route protection) independently.

4. **Minimizes Risk**: Unlike Option 1, we add validation and error handling. Unlike Option 2, we don't delay the fix. Unlike Option 4, we maintain SPA functionality.

## Detailed Approach

### Phase 1: Enable Basic Static Serving
1. Remove early return from `registerStatic()`
2. Uncomment registration call in `routes.go`
3. Add error handling and logging around middleware setup
4. Test with minimal frontend assets

### Phase 2: Validate Route Protection
1. Verify API routes (`/api/*`) are not shadowed
2. Verify admin UI (`/_/`) remains accessible
3. Test SPA fallback behavior for client-side routes

### Phase 3: Handle Edge Cases
1. Add proper MIME type handling
2. Implement caching headers if needed
3. Add graceful fallback if frontend assets missing

### Implementation Strategy

**Error Handling Approach:**
- Wrap middleware registration in try-catch equivalent
- Log specific error messages for debugging
- Graceful degradation if static serving fails
- Clear error messages for missing frontend assets

**Testing Strategy:**
- Manual testing with `go run main.go`
- Verify each route type (static, API, admin) works
- Test both development (missing assets) and production (built assets) scenarios

**Rollback Plan:**
- If implementation fails, revert to early return pattern
- Document specific compatibility issues found
- Provides data for future attempts (Option 2 or 4)

## Expected Outcomes

**Success Criteria:**
- Root path `/` serves Astro frontend when assets exist
- API routes continue working without interference
- PocketBase admin UI remains accessible at `/_/`
- SPA client-side routing works correctly
- Graceful handling when frontend assets are missing

**Risk Mitigation:**
- Comprehensive error logging identifies specific failure points
- Incremental approach allows rollback at any step
- Maintains all existing functionality as fallback

## Technical Implementation Notes

### Middleware Registration Pattern
Current pattern uses Echo v5's `StaticWithConfig` middleware:
```go
e.Router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
    Root:       "/",
    Filesystem: http.FS(spaFilesystem),
    Browse:     false,
}))
```

This approach integrates with Echo's middleware chain and should be compatible with PocketBase v0.36.5's router usage.

### SPA Filesystem Wrapper
The existing `spaFS` struct provides intelligent fallback:
- Serves exact file matches when available
- Falls back to `index.html` for SPA routes
- Protects `/api/` and `/_/` paths from fallback

This design correctly handles modern SPA requirements and should be preserved.

### Asset Validation
Before serving, validate that embedded assets exist and contain expected files (at minimum `index.html`). This prevents runtime errors and provides clear feedback during development.

## Conclusion

Option 3 provides a pragmatic solution that respects the existing codebase while addressing the immediate problem safely. The incremental approach with proper error handling ensures we can identify and fix any PocketBase v0.36.5 compatibility issues without breaking existing functionality.

The solution maintains the architectural integrity of the codebase while enabling the critical functionality needed for the frontend to be served correctly.