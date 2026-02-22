# Research: Delete Dead Code and Duplicates

## Dead Code Locations Identified

### 1. `routes/links_search_simple.go` (Entire File - 191 lines)
**Status**: Complete duplicate of `links_search.go`
- Registered at `/api/links/search-simple` endpoint
- Called from `routes/routes.go:17` via `registerLinksSearchSimple(e)`
- Duplicates all structures and functions with "2" suffix:
  - `SearchParams2`, `SearchResponse2`, `LinkItem2`
  - `handleSearchSimple`, `parseSearchParams2`, `validateSearchParams2`, `escapeLikePattern2`
- Missing features compared to `links_search.go`:
  - No tag filtering implementation (tag parameter ignored)
  - No proper total count (uses result length instead of COUNT query)
  - No timestamps in response (created/updated fields)
  - Simpler query without JOINs

### 2. `routes/presentations.go` Dead Functions

#### Line 174: `validateToken()`
**Status**: Unused function
- Returns boolean for token comparison
- Uses `subtle.ConstantTimeCompare` for security
- No calls found in codebase (grep confirms)
- Similar function `validateSyncToken()` exists in `sync_sessions.go` and is actively used

#### Line 106: `stepToProgress()`
**Status**: Unused function
- Inverse of `progressToStep()` (line 95)
- Converts step index to progress value (0-1)
- No calls found in Go code
- Frontend has own implementation in `frontend/src/utils/stepConversion.ts`
- `progressToStep()` IS used in `buildStatusResponse()` at line 165

#### Line 181: Stale comment
**Status**: Obsolete comment
- States "Placeholder handler functions - to be implemented in subsequent steps"
- Functions below ARE implemented: `handleGetStatus`, `handleStopLive`, `handleStartLive`

#### Line 7: Commented import
**Status**: Unused import
- `// "encoding/json"  // TODO: Re-enable if JSON operations are added`
- No JSON marshal/unmarshal operations in the file
- All JSON handling done by PocketBase framework

### 3. `routes/routes_test.go` Dead Code

#### Line 17-22: `TestData` struct
**Status**: Unused struct
- Defined but never instantiated
- Referenced in `stats_test.go` imports but not used there either
- Likely leftover from test setup refactoring

#### Line 157-168: `assertErrorResponse()` function
**Status**: Unused function
- Helper for validating error responses
- Never called in any test
- Similar validation done inline in individual tests

### 4. `internal/frontend/embed.go` Dead Export

#### Line 20-23: `FrontendExists()` function
**Status**: Used but can be inlined
- Exported function, only called once in `routes/routes.go:21`
- Simple wrapper around `GetFrontendFS()` error check
- Can be replaced by direct error check in caller

### 5. `routes/static.go` Unreachable Code

#### Lines 74-79: Fallback block
**Status**: Unreachable code
- After `if readSeeker, ok := file.(io.ReadSeeker); ok` check fails
- Returns raw string "file serving not supported"
- All `fs.File` implementations from `embed.FS` implement `io.ReadSeeker`
- Better to call `ev.Next()` if cast fails (unlikely scenario)

## Cross-References and Dependencies

### Import Dependencies
- `links_search_simple.go` imports same packages as `links_search.go`
- No unique imports that would be lost

### Route Registration Flow
1. `main.go` calls `routes.Register(app)`
2. `routes.Register()` sets up `OnServe` hook
3. Hook calls individual `register*` functions including:
   - `registerLinksSearchSimple()` - to be removed
   - `registerStatic()` - conditional on `FrontendExists()`

### Test Infrastructure
- `makeRequest()` in `routes_test.go` IS actively used by `TestMakeRequest_RealExecution`
- Uses `echo.New()` from `labstack/echo/v5` - dependency will remain
- `httptest` package used for test execution

## Module Dependencies

### Direct Dependencies (go.mod)
- `github.com/labstack/echo/v5` - Used by `routes_test.go` for `echo.New()`
- `github.com/pocketbase/dbx` - Used throughout for database queries
- `github.com/pocketbase/pocketbase` - Core framework

All three direct dependencies are actively used and will remain after cleanup.

## Migration Risks

### Low Risk Changes
1. Deleting `links_search_simple.go` - duplicate functionality exists
2. Removing unused test helpers - no test dependencies
3. Removing commented imports - already non-functional
4. Removing stale comments - documentation only

### Medium Risk Changes
1. Removing `FrontendExists()` - requires updating caller logic
2. Modifying static.go fallback - affects error path (rarely executed)

### No Risk from `go mod tidy`
- All direct dependencies actively used
- Transitive dependencies managed automatically

## Verification Requirements

### Build Verification
- `go build .` must succeed
- No unresolved imports

### Test Verification
- `go test ./...` must pass
- `TestMakeRequest_RealExecution` specifically validates routing

### Lint Verification
- `go vet ./...` must pass
- No unused variables or imports

### Runtime Verification
- `/api/links/search` endpoint remains functional
- Static file serving continues to work
- Frontend embedding detection works without `FrontendExists()`