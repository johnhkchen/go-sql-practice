# T-009-01 Structure: Fix Static File Serving

## File Modifications

### Modified Files

#### `routes/static.go`
**Changes:**
- Remove early return on line 36
- Uncomment Echo v5 imports (lines 5, 9, 10)
- Add comprehensive error handling around middleware setup
- Add logging for debugging static serving issues
- Add validation for frontend filesystem availability

**Structure:**
```go
func registerStatic(e *core.ServeEvent) {
    // Validate frontend assets exist
    frontendFS, err := frontend.GetFrontendFS()
    if err != nil {
        e.App.Logger().Error("Frontend assets not available", "error", err)
        return // Graceful degradation
    }

    // Log static serving initialization
    e.App.Logger().Info("Initializing static file serving")

    // Create SPA filesystem with existing logic
    spaFilesystem := &spaFS{fs: frontendFS}

    // Register static middleware with error handling
    err = e.Router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
        Root:       "/",
        Filesystem: http.FS(spaFilesystem),
        Browse:     false,
    }))

    if err != nil {
        e.App.Logger().Error("Failed to register static middleware", "error", err)
        return
    }

    e.App.Logger().Info("Static file serving enabled")
}
```

**Rationale:** Preserve existing SPA logic while adding safety measures and proper error handling.

#### `routes/routes.go`
**Changes:**
- Uncomment `registerStatic(e)` call on line 20
- Add conditional registration based on frontend availability
- Add logging for route registration status

**Structure:**
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        // Register individual routes (existing)
        registerHealth(e)
        registerSyncSessions(e)
        registerPresentations(e)
        registerStats(e)
        registerLinksSearch(e)
        registerLinksSearchSimple(e)
        registerLinksView(e)

        // Register static file serving with availability check
        if frontend.FrontendExists() {
            registerStatic(e)
        } else {
            e.App.Logger().Warn("Frontend assets not found, static serving disabled")
        }

        // Continue middleware chain
        return e.Next()
    })
}
```

**Rationale:** Enable static serving while maintaining graceful degradation when frontend assets are missing.

### No New Files Required

The existing file structure is sufficient. The implementation leverages:
- Existing `internal/frontend/embed.go` for asset management
- Existing `spaFS` struct for SPA routing
- Existing error handling patterns in the codebase

### No Files Deleted

All existing files remain. The solution builds upon existing infrastructure rather than replacing it.

## Component Architecture

### Static Serving Flow
```
HTTP Request "/"
    ↓
PocketBase Router
    ↓
Static Middleware (if registered)
    ↓
spaFS.Open(path)
    ↓
[Exact file exists?] → Yes → Serve file
    ↓ No
[Protected path (/api/, /_/)?] → Yes → Return error
    ↓ No
Serve index.html (SPA fallback)
```

### Error Handling Chain
```
registerStatic() called
    ↓
[Frontend assets exist?] → No → Log warning, return gracefully
    ↓ Yes
[Middleware registration succeeds?] → No → Log error, return
    ↓ Yes
Static serving enabled
```

### Route Protection Logic
Existing `spaFS.Open()` method handles route protection:
- `/api/*` paths: Return error (let API handlers process)
- `/_/*` paths: Return error (preserve PocketBase admin UI)
- Other paths: Serve file or fallback to index.html

## Interface Boundaries

### Public Interfaces (Unchanged)
- `routes.Register(app core.App)`: Main registration entry point
- `frontend.GetFrontendFS()`: Filesystem access
- `frontend.FrontendExists()`: Asset availability check

### Internal Interfaces (Enhanced)
- `registerStatic(e *core.ServeEvent)`: Enhanced error handling
- `spaFS.Open(name string)`: Existing SPA routing logic (unchanged)

### Error Interfaces
- Logging via PocketBase's logger (`e.App.Logger()`)
- Graceful degradation (return without error on asset unavailability)
- Error propagation for middleware registration failures

## Module Dependencies

### Existing Dependencies (Preserved)
- `github.com/pocketbase/pocketbase/core`: ServeEvent, App interfaces
- `github.com/jchen/go-sql-practice/internal/frontend`: Asset management
- `github.com/labstack/echo/v5`: HTTP router and middleware (re-enabled)
- `github.com/labstack/echo/v5/middleware`: Static serving middleware (re-enabled)
- `net/http`: File serving interfaces (re-enabled)

### No New Dependencies
The solution works entirely within the existing dependency set.

## Implementation Order

### Phase 1: Enable Core Static Serving
1. Modify `routes/static.go`:
   - Remove early return
   - Uncomment imports
   - Add basic error handling
2. Modify `routes/routes.go`:
   - Uncomment `registerStatic(e)` call
3. Test basic functionality

### Phase 2: Add Robust Error Handling
1. Enhance `routes/static.go`:
   - Add frontend validation
   - Add comprehensive logging
   - Add middleware registration error handling
2. Enhance `routes/routes.go`:
   - Add conditional registration
   - Add asset availability logging
3. Test error scenarios

### Phase 3: Validation and Testing
1. Test route protection (API, admin UI)
2. Test SPA fallback behavior
3. Test missing asset scenarios
4. Verify build integration

## Testing Strategy

### Manual Testing Approach
1. **Without frontend assets**: Verify graceful degradation
2. **With frontend assets**: Verify static serving works
3. **API route protection**: Ensure `/api/health` still works
4. **Admin UI protection**: Ensure `/_/` admin interface works
5. **SPA routing**: Verify client-side routes fallback to index.html

### Integration Points
- Build system: `make frontend` → `make backend` → test static serving
- Route conflicts: Verify no interference with existing API routes
- Asset embedding: Confirm embedded assets are served correctly

## Rollback Strategy

### If Implementation Fails
1. Re-add early return to `registerStatic()`
2. Re-comment imports in `static.go`
3. Re-comment `registerStatic(e)` call in `routes.go`

### Failure Detection
- Middleware registration errors logged
- Route conflicts detected via testing
- Performance issues identified via monitoring

### Recovery Approach
Each phase can be rolled back independently, allowing precise identification of failure points.

## Security Considerations

### Route Protection Maintained
- Existing `spaFS` logic preserves API and admin route protection
- No changes to authentication or authorization patterns
- Static assets served with appropriate MIME types

### Asset Validation
- Frontend asset existence validated before serving
- No directory traversal risks (embedded filesystem)
- SPA fallback contained to non-protected paths

## Performance Impact

### Minimal Overhead
- Static middleware adds one layer to request processing
- File serving handled by Go's efficient http.FileServer
- Embedded assets served from memory (no disk I/O)

### Caching Considerations
- Browser caching handled by standard HTTP headers
- No server-side caching required for embedded assets
- SPA assets typically include cache-busting hashes

## Conclusion

This structure provides a safe, incremental approach to enabling static file serving while preserving all existing functionality. The solution builds on the existing codebase architecture and adds appropriate safety measures without introducing complexity or new dependencies.

The modular design allows for independent testing and rollback of each component, ensuring reliable delivery of the critical static serving functionality.