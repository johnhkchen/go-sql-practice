# Plan: T-001-01 init-go-module

## Implementation Steps

### Step 1: Initialize Go Module

**Command**:
```bash
go mod init github.com/jchen/go-sql-practice
```

**Verification**:
- `go.mod` exists in root
- Contains `module github.com/jchen/go-sql-practice`
- File is valid (no syntax errors)

**Commit**: No (module incomplete)

### Step 2: Set Go Version

**Command**:
```bash
go mod edit -go=1.26
```

**Verification**:
- `go.mod` contains `go 1.26`
- Version matches Flox pin

**Commit**: No (module incomplete)

### Step 3: Add PocketBase Dependency

**Command**:
```bash
go get github.com/pocketbase/pocketbase@v0.36.5
```

**Verification**:
- `go.mod` contains `require github.com/pocketbase/pocketbase v0.36.5`
- `go.sum` exists and is populated
- No download errors

**Commit**: Yes - "Initialize Go module with PocketBase v0.36.5"

### Step 4: Create main.go

**Action**: Write file with content:
```go
package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
```

**Verification**:
- File exists at `/home/jchen/repos/go-sql-practice/main.go`
- No syntax errors (check with `go vet`)
- Imports resolve correctly

**Commit**: Yes - "Add main.go entry point with PocketBase initialization"

### Step 5: Build Verification

**Command**:
```bash
go build -o go-sql-practice
```

**Verification**:
- Build completes without errors
- Binary `go-sql-practice` exists
- Binary is executable
- Size approximately 30-40MB

**Commit**: No (binary not tracked)

### Step 6: Runtime Verification

**Command**:
```bash
./go-sql-practice serve
```

**Verification**:
- Server starts without errors
- Logs show "Server started at: http://127.0.0.1:8090"
- Process responds to Ctrl+C gracefully
- Creates `pb_data` directory

**Commit**: No (runtime test only)

### Step 7: Admin UI Verification

**Action**: Access `http://localhost:8090/_/` in browser or via curl

**Command** (alternative):
```bash
curl -I http://localhost:8090/_/
```

**Verification**:
- Returns HTTP 200 OK
- Admin UI loads (if browser test)
- No console errors

**Commit**: No (verification only)

### Step 8: Update .gitignore

**Action**: Add PocketBase runtime artifacts to `.gitignore`

**Content to add**:
```
# PocketBase
pb_data/
*.db
*.db-shm
*.db-wal

# Go build artifacts
go-sql-practice
*.exe
```

**Verification**:
- `.gitignore` exists or is created
- Contains PocketBase patterns
- `git status` doesn't show `pb_data/`

**Commit**: Yes - "Add PocketBase and Go artifacts to .gitignore"

## Testing Strategy

### Unit Tests

Not applicable - main package with single function.

### Integration Tests

Manual verification steps:

1. **Clean build test**:
   ```bash
   rm -rf pb_data go-sql-practice
   go build
   ```
   Expected: Builds successfully

2. **Server start test**:
   ```bash
   ./go-sql-practice serve &
   sleep 2
   curl -s http://localhost:8090/_/ | grep -q "PocketBase"
   pkill go-sql-practice
   ```
   Expected: Admin UI accessible

3. **Dependency test**:
   ```bash
   go list -m all | grep pocketbase
   ```
   Expected: Shows `github.com/pocketbase/pocketbase v0.36.5`

### Verification Criteria

**Build Success**:
- [ ] `go build` completes without errors
- [ ] Binary size between 25-45MB
- [ ] No missing dependency errors

**Runtime Success**:
- [ ] Server starts on port 8090
- [ ] Admin UI accessible at `/_/`
- [ ] No panic or fatal errors on startup
- [ ] Clean shutdown with Ctrl+C

**Module Correctness**:
- [ ] Module path is `github.com/jchen/go-sql-practice`
- [ ] Go version is `1.26`
- [ ] PocketBase version is exactly `v0.36.5`
- [ ] go.sum contains all dependency hashes

## Error Recovery

### Common Issues and Fixes

1. **Port already in use**:
   - Error: "bind: address already in use"
   - Fix: Kill existing process or use different port

2. **Module cache issues**:
   - Error: "checksum mismatch"
   - Fix: `go clean -modcache` then retry

3. **Permission denied**:
   - Error: Cannot create pb_data
   - Fix: Check directory permissions

4. **Build failures**:
   - Error: Package not found
   - Fix: `go mod download` then retry

## Commit Strategy

Three atomic commits:

1. **"Initialize Go module with PocketBase v0.36.5"**
   - Files: `go.mod`, `go.sum`
   - After: Step 3

2. **"Add main.go entry point with PocketBase initialization"**
   - Files: `main.go`
   - After: Step 4

3. **"Add PocketBase and Go artifacts to .gitignore"**
   - Files: `.gitignore`
   - After: Step 8

## Progress Tracking

Will update `progress.md` after each step with:
- Step completed
- Any deviations from plan
- Issues encountered
- Time spent

## Dependencies and Prerequisites

**Required**:
- Go 1.26 (available via Flox)
- Internet connection (for go get)
- Port 8090 available
- Write permissions in project directory

**Assumed Complete**:
- T-001-00 (Flox setup) - VERIFIED
- Git repository initialized - VERIFIED
- RDSPI phases 1-4 complete - YES

## Success Metrics

**Primary**:
- Binary builds successfully
- PocketBase starts and serves admin UI
- All acceptance criteria met

**Secondary**:
- Clean code with standard formatting
- No warnings from `go vet`
- Reasonable binary size (<45MB)
- Fast startup time (<1 second)