# Plan: Custom Health Route Implementation

## Implementation Steps

### Step 1: Create routes package directory
**Action**: Create the `routes/` directory at project root
**Command**: `mkdir routes`
**Verification**: Directory exists at `./routes/`

### Step 2: Implement routes.go
**Action**: Create `routes/routes.go` with Register function
**File**: `routes/routes.go`
**Content**:
- Package declaration
- Import `github.com/pocketbase/pocketbase/core`
- Register function that binds to OnServe
- Call to registerHealth (forward declaration)
- Return e.Next() to continue chain
**Verification**: File compiles (will have undefined registerHealth initially)

### Step 3: Implement health.go
**Action**: Create `routes/health.go` with handler
**File**: `routes/health.go`
**Content**:
- Package declaration (routes)
- Import core and echo packages
- registerHealth function
- healthHandler function returning JSON
**Verification**: Package compiles without errors

### Step 4: Update main.go imports
**Action**: Add routes package import
**File**: `main.go`
**Change**: Add `"github.com/jchen/go-sql-practice/routes"` to imports
**Verification**: Import statement present

### Step 5: Register routes in main.go
**Action**: Call routes.Register after migrations
**File**: `main.go`
**Change**: Add `routes.Register(app)` after line 14
**Verification**: Function call present in main()

### Step 6: Build the application
**Action**: Compile the Go application
**Command**: `go build -o go-sql-practice`
**Verification**: Binary created successfully, no compilation errors

### Step 7: Test server startup
**Action**: Start the server and check for errors
**Command**: `./go-sql-practice serve`
**Verification**:
- Server starts without errors
- Logs show "Server started" or similar
- Can Ctrl+C to stop

### Step 8: Test health endpoint
**Action**: Make GET request to health endpoint
**Command**: `curl -i http://localhost:8090/api/health`
**Expected Response**:
```
HTTP/1.1 200 OK
Content-Type: application/json
{"status":"ok"}
```
**Verification**:
- Status code is 200
- JSON response matches expected
- No authentication required

### Step 9: Test endpoint without auth
**Action**: Verify endpoint works without authentication headers
**Command**: `curl http://localhost:8090/api/health`
**Verification**: Returns `{"status":"ok"}` without auth headers

### Step 10: Verify PocketBase routes still work
**Action**: Check that admin UI is still accessible
**Command**: Open browser to `http://localhost:8090/_/`
**Verification**: PocketBase admin UI loads correctly

## Testing Strategy

### Unit Tests
Not required for this simple handler initially, but could test:
- JSON response structure
- Status code
- No side effects

### Integration Tests
Manual testing via curl commands to verify:
1. Endpoint is accessible at correct path
2. Returns correct JSON structure
3. Returns 200 status code
4. Works without authentication
5. Doesn't break existing PocketBase functionality

### Smoke Tests
Quick verification after each step:
- After Step 6: Binary builds
- After Step 7: Server starts
- After Step 8: Endpoint responds
- After Step 10: Nothing broke

## Rollback Plan

If issues occur at any step:

### Compilation Errors
- Check import paths
- Verify package names match
- Ensure all functions are defined

### Runtime Errors
- Check OnServe hook implementation
- Verify e.Next() is called
- Check for route conflicts

### Route Not Found
- Verify route registration in OnServe
- Check path string for typos
- Ensure server is running on expected port

### Complete Rollback
1. Delete `routes/` directory
2. Remove routes import from main.go
3. Remove routes.Register() call
4. Rebuild and verify original functionality

## Commit Points

### Commit 1: Create routes package
After Steps 1-3:
```bash
git add routes/
git commit -m "feat: add routes package with health endpoint handler"
```

### Commit 2: Integrate routes with main
After Steps 4-5:
```bash
git add main.go
git commit -m "feat: register custom routes in main application"
```

### Commit 3: Complete implementation
After successful testing (Steps 6-10):
```bash
git add -A
git commit -m "feat: implement GET /api/health endpoint returning status ok"
```

## Time Estimates

- Step 1-3: 2 minutes (file creation)
- Step 4-5: 1 minute (main.go updates)
- Step 6: 30 seconds (build)
- Step 7-10: 3 minutes (testing)
- Total: ~7 minutes

## Success Criteria

All acceptance criteria from ticket are met:
- [x] routes/ package exists for custom route registration
- [x] GET /api/health returns 200 with {"status": "ok"}
- [x] Route registered via OnServe hook
- [x] Health endpoint accessible without authentication
- [x] Code compiles and binary starts without errors

## Notes

- Echo v5 is already a PocketBase dependency, no go.mod changes needed
- The route pattern established here can be reused for future endpoints
- Consider adding logging in future iterations
- No database checks in health endpoint per design decision