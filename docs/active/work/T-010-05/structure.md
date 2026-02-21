# Structure for T-010-05: Infrastructure Cleanup

## File Modifications

### Modified Files

#### `.gitignore`
**Location**: `/.gitignore`
**Changes**: Append two new patterns
- Add `*.test` entry after line 3 (Go section)
- Add `routes/pb_data/` after line 6 (below main pb_data/)

#### `Makefile`
**Location**: `/Makefile`
**Changes**: Add targets and fix documentation
- Line 8: Update `.PHONY` to include `lint vet`
- After line 55: Insert new `lint` target
- After lint: Insert new `vet` target
- Line 74: Fix help text from "validate" to "validate-build"

#### `.github/workflows/ci.yml`
**Location**: `/.github/workflows/ci.yml`
**Changes**: Modernize and add quality checks
- Line 22: Update `actions/setup-go@v4` to `@v5`
- After line 43: Insert `go vet` step
- After vet: Insert `gofmt` check step
- After line 36: Insert `npm audit` step

#### `CLAUDE.md`
**Location**: `/CLAUDE.md`
**Changes**: Replace placeholder description
- Line 5: Replace TODO with actual description

#### `routes/stats.go`
**Location**: `/routes/stats.go`
**Changes**: Add logging before error returns
- Before line 76: Add error log
- Before line 82: Add error log
- Before line 88: Add error log
- Before line 94: Add error log
- Before line 100: Add error log

### Created Files

#### `frontend/.env.example`
**Location**: `/frontend/.env.example`
**Content**: Single environment variable example
```
PUBLIC_API_URL=http://127.0.0.1:8090
```

### Deleted Files

#### Test binary
**Location**: `/routes.test`
**Action**: Remove from git tracking and delete

#### Duplicate frontend directory
**Location**: `/frontend/frontend/`
**Action**: Recursively delete entire directory tree

#### Test database artifacts
**Location**: `/routes/pb_data/`
**Action**: Recursively delete entire directory

## Module Boundaries

### Build System
- **Makefile**: Central build orchestration
- **CI workflow**: Automated validation
- No changes to module interfaces

### Error Handling Layer
- **routes/stats.go**: Example of logging pattern
- Maintains existing public API contract
- Internal logging enhancement only

### Configuration
- **Frontend config**: Documented via .env.example
- **Project config**: Updated CLAUDE.md
- No runtime configuration changes

## Component Organization

### Quality Targets (Makefile)
```make
lint:
	@gofmt -l . | grep -v "^$$" && exit 1 || echo "Code is formatted"

vet:
	@go vet ./...
```

### CI Quality Steps
```yaml
- name: Run go vet
  run: go vet ./...

- name: Check formatting
  run: |
    if [ -n "$(gofmt -l .)" ]; then
      echo "Code needs formatting. Run: gofmt -w ."
      gofmt -l .
      exit 1
    fi

- name: Audit npm packages
  run: cd frontend && npm audit --audit-level=high
```

### Logging Pattern (routes/stats.go)
```go
if err != nil {
    app.Logger().Error("getTotalLinks failed", "error", err)
    return e.JSON(500, map[string]string{"error": "Failed to get total links"})
}
```

## Change Ordering

1. **Clean tracked files first**
   - Update .gitignore
   - Remove routes.test from git
   - Delete directories

2. **Update build system**
   - Modify Makefile
   - Create .env.example

3. **Enhance quality checks**
   - Update CI workflow
   - Add logging to routes

4. **Documentation last**
   - Update CLAUDE.md

## Interface Preservation

### Public APIs: Unchanged
- HTTP endpoints maintain same request/response format
- Error response structure preserved
- Build commands remain compatible

### Internal APIs: Enhanced
- Logging added but not required
- New make targets are optional
- CI checks are additive

## Testing Boundaries

### Unit Test Scope
- No new unit tests needed (no logic changes)
- Existing tests validate preservation

### Integration Test Scope
- `make lint` and `make vet` are self-testing
- CI changes tested by GitHub Actions

### Manual Validation
- Directory deletion verification
- .gitignore effectiveness
- Build still produces working binary