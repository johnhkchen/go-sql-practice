# T-009-01 Plan: Fix Static File Serving

## Implementation Steps

### Step 1: Enable Basic Static Serving
**Objective**: Remove the blocking mechanisms and restore basic functionality

**Tasks**:
1. Remove early return from `routes/static.go` line 36
2. Uncomment Echo v5 imports in `routes/static.go` (lines 5, 9, 10)
3. Uncomment `registerStatic(e)` call in `routes/routes.go` line 20
4. Test basic compilation and startup

**Verification**:
- `go build` compiles without errors
- Application starts without panics
- No obvious runtime errors in logs

**Commit**: "fix: remove early return blocking static file serving"

### Step 2: Add Error Handling and Validation
**Objective**: Make static serving robust and provide clear feedback

**Tasks**:
1. Add frontend asset validation at start of `registerStatic()`
2. Add error handling around middleware registration
3. Add logging statements for debugging
4. Implement graceful degradation when assets missing

**Code Changes**:
```go
// In routes/static.go
func registerStatic(e *core.ServeEvent) {
    frontendFS, err := frontend.GetFrontendFS()
    if err != nil {
        e.App.Logger().Error("Frontend assets not available", "error", err)
        return
    }

    e.App.Logger().Info("Initializing static file serving")

    spaFilesystem := &spaFS{fs: frontendFS}

    // Rest of implementation...
}
```

**Verification**:
- Error messages appear when frontend assets missing
- Success messages appear when static serving enabled
- No panics or crashes during startup

**Commit**: "feat: add error handling and validation to static serving"

### Step 3: Enhance Route Registration
**Objective**: Make route registration conditional and informative

**Tasks**:
1. Add conditional registration based on `frontend.FrontendExists()`
2. Add appropriate logging for both success and failure cases
3. Ensure graceful operation in both development and production

**Code Changes**:
```go
// In routes/routes.go
if frontend.FrontendExists() {
    registerStatic(e)
} else {
    e.App.Logger().Warn("Frontend assets not found, static serving disabled")
}
```

**Verification**:
- Appropriate log messages in both asset-available and asset-missing scenarios
- No errors when frontend not built
- Static serving enabled when assets available

**Commit**: "feat: add conditional static route registration"

### Step 4: Build Frontend Assets for Testing
**Objective**: Create test assets to validate static serving works

**Tasks**:
1. Run `make frontend` to build Astro assets
2. Verify embedded assets are created correctly
3. Check that `internal/frontend/frontend/dist/client/` contains built files

**Verification**:
- Frontend build completes successfully
- `internal/frontend/frontend/dist/client/` contains `index.html` and assets
- `frontend.FrontendExists()` returns true
- Embedded filesystem includes expected files

**Commit**: Not needed (build artifacts)

### Step 5: Test Route Protection
**Objective**: Verify existing API and admin routes still work

**Tasks**:
1. Start application with built frontend
2. Test `/api/health` endpoint works
3. Test `/_/` admin interface loads
4. Verify root `/` serves frontend content
5. Test non-existent paths fallback to `index.html`

**Test Cases**:
- `GET /` → Returns Astro frontend (200 OK)
- `GET /api/health` → Returns JSON health response (200 OK)
- `GET /_/` → Returns PocketBase admin interface (200 OK)
- `GET /nonexistent` → Returns `index.html` content (200 OK, SPA fallback)

**Verification**:
- All route types work as expected
- No route conflicts or shadowing
- SPA client-side routing supported

**Commit**: "test: verify route protection and SPA functionality"

### Step 6: Test Error Scenarios
**Objective**: Verify graceful handling of edge cases

**Tasks**:
1. Test startup without frontend assets built
2. Test static serving after assets are deleted
3. Verify error messages are clear and actionable
4. Test recovery after assets are restored

**Test Cases**:
- Start with missing frontend → Graceful degradation, clear logs
- Build frontend while running → Static serving remains disabled (expected)
- Invalid embedded assets → Error handling prevents crashes

**Verification**:
- No crashes or panics in any scenario
- Clear, helpful error messages
- Graceful degradation maintains API/admin functionality

**Commit**: "test: validate error handling and edge cases"

## Testing Strategy

### Unit Testing
- **Scope**: Limited - most functionality is integration-level
- **Focus**: Error handling logic in `registerStatic()`
- **Approach**: Mock PocketBase logger and frontend filesystem

### Integration Testing
- **Scope**: Primary testing approach
- **Focus**: End-to-end request handling for all route types
- **Approach**: Manual testing with real PocketBase instance

### Manual Test Scenarios

#### Scenario 1: Development Environment (No Frontend Built)
1. Start application: `go run main.go`
2. Verify graceful startup with warning logs
3. Test API routes work: `curl http://localhost:8090/api/health`
4. Test admin UI works: `http://localhost:8090/_/`
5. Test root returns nothing (expected): `curl http://localhost:8090/`

#### Scenario 2: Production Environment (Frontend Built)
1. Build frontend: `make frontend`
2. Start application: `go run main.go`
3. Verify successful startup with static serving enabled logs
4. Test all route types work correctly
5. Test SPA fallback behavior

#### Scenario 3: Asset Recovery
1. Start with missing assets (Scenario 1)
2. Build frontend while running
3. Verify static serving remains disabled (requires restart)
4. Restart application
5. Verify static serving now enabled (Scenario 2)

### Performance Testing
- **Load Testing**: Not required for this ticket
- **Memory Usage**: Monitor for asset embedding overhead
- **Response Times**: Verify static serving doesn't slow API routes

## Error Handling Strategy

### Error Categories

#### Expected Errors (Graceful Handling)
- Missing frontend assets during development
- Invalid embedded filesystem structure
- Middleware registration conflicts

#### Unexpected Errors (Logging + Recovery)
- PocketBase logger unavailable
- Echo router unavailable
- Memory issues during asset serving

### Logging Strategy
- **Info Level**: Static serving enabled/disabled status
- **Warn Level**: Missing assets (expected during development)
- **Error Level**: Unexpected failures that prevent static serving

### Recovery Strategy
- **Primary**: Graceful degradation - continue without static serving
- **Secondary**: Clear error messages for debugging
- **Fallback**: Maintain all existing API and admin functionality

## Rollback Plan

### Automatic Rollback Triggers
- Application fails to compile
- Application panics during startup
- API routes stop working
- Admin interface becomes inaccessible

### Rollback Steps
1. Revert to previous commit
2. Re-add early return to `routes/static.go:36`
3. Re-comment imports and registration call
4. Verify application returns to previous working state

### Rollback Verification
- Application compiles and starts successfully
- All existing API routes work
- Admin interface accessible
- Root path returns nothing (original behavior)

## Success Criteria

### Primary Success Criteria (Must Have)
- ✅ `registerStatic(e)` is uncommented and working
- ✅ `static.go` serves the embedded frontend files on `/`
- ✅ PocketBase admin UI still works at `/_/`
- ✅ API routes (`/api/*`) are not shadowed by static file serving
- ✅ `go build` compiles without errors
- ✅ Navigating to `http://localhost:8090/` shows the Astro index page

### Secondary Success Criteria (Should Have)
- Clear error messages when frontend assets missing
- Graceful degradation in development environment
- Comprehensive logging for debugging
- SPA client-side routing works correctly

### Quality Gates
- No runtime panics or crashes
- No degradation of existing functionality
- Clear documentation of behavior changes
- Maintainable code with proper error handling

## Dependencies and Constraints

### Build Dependencies
- Frontend must be built before testing static serving
- `make frontend` must complete successfully
- Embedded assets must be valid and accessible

### Runtime Dependencies
- PocketBase v0.36.5 compatibility maintained
- Echo v5 middleware system functioning correctly
- Embedded filesystem accessible and valid

### External Constraints
- No changes to PocketBase configuration
- No changes to build system or deployment process
- No impact on existing API contracts or admin interface

## Risk Mitigation

### High-Risk Areas
1. **Echo v5 Compatibility**: Monitor for deprecated middleware patterns
2. **Route Conflicts**: Verify static serving doesn't shadow critical routes
3. **Performance Impact**: Watch for memory usage increases

### Mitigation Strategies
1. **Incremental Implementation**: Each step is independently verifiable and rollbackable
2. **Comprehensive Testing**: Manual validation of all route types and scenarios
3. **Graceful Degradation**: Ensure core functionality works even if static serving fails
4. **Clear Documentation**: Artifact trail enables future maintenance and debugging

### Contingency Planning
- If Echo v5 patterns don't work: Research PocketBase-native static serving
- If route conflicts occur: Implement path-specific middleware registration
- If performance issues arise: Add conditional asset loading or caching

This plan provides a safe, incremental approach to restoring static file serving while maintaining all existing functionality and providing clear validation at each step.