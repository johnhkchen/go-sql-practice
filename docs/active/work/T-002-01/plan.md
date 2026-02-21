# Plan: T-002-01 define-collections

## Implementation Steps

### Step 1: Create migrations package directory

**Action**: Create the `migrations` directory in project root

**Command**:
```bash
mkdir migrations
```

**Verification**:
- Directory exists at `go-sql-practice/migrations/`
- Ready for Go files

**Commit**: No (structure only)

### Step 2: Fix PocketBase dependency

**Action**: Update go.mod to mark PocketBase as direct dependency

**Commands**:
```bash
go get github.com/pocketbase/pocketbase@v0.36.5
go mod tidy
```

**Verification**:
- go.mod shows PocketBase without "// indirect"
- go.sum updated

**Commit**: "fix: mark pocketbase as direct dependency"

### Step 3: Write migrations/collections.go

**Action**: Create the collections migration file with all field definitions

**File Content Structure**:
```go
package migrations

import (
    "fmt"
    "github.com/pocketbase/pocketbase/core"
)

func Register(app core.App) { ... }
func createCollections(txApp core.App) error { ... }
```

**Key Implementation Points**:
1. Register function binds to OnBootstrap
2. createCollections creates tags first, then links
3. Proper error wrapping at each step
4. Field configurations match specifications exactly

**Verification**:
- File compiles without errors
- All imports resolved
- No syntax issues

**Commit**: "feat: add migrations package with collection definitions"

### Step 4: Update main.go

**Action**: Import migrations package and call Register

**Changes**:
```go
import (
    "log"
    "github.com/jchen/go-sql-practice/migrations"
    "github.com/pocketbase/pocketbase"
)

func main() {
    app := pocketbase.New()

    // Register migrations
    migrations.Register(app)

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

**Verification**:
- Imports resolve correctly
- No compilation errors
- Migration registration before app.Start()

**Commit**: "feat: register migrations in main.go"

### Step 5: Test migration execution

**Action**: Run the application with fresh database

**Commands**:
```bash
# Remove existing data
rm -rf pb_data

# Build and run
go build -o go-sql-practice
./go-sql-practice serve
```

**Expected Output**:
- Server starts on :8090
- No error messages
- "Bootstrap completed" or similar log

**Verification Checklist**:
- [ ] Server starts without errors
- [ ] pb_data directory created
- [ ] No migration errors in console

### Step 6: Verify collections in admin UI

**Action**: Access PocketBase admin and check collections

**Steps**:
1. Open browser to `http://localhost:8090/_/`
2. Create admin account if first run
3. Navigate to Collections section
4. Verify both collections exist

**Verification Checklist**:

**Tags Collection**:
- [ ] Collection name: "tags"
- [ ] Field: name (text, required)
- [ ] Field: slug (text, required, unique)
- [ ] System fields present (id, created, updated)

**Links Collection**:
- [ ] Collection name: "links"
- [ ] Field: url (url, required)
- [ ] Field: title (text, required)
- [ ] Field: description (text, optional)
- [ ] Field: tags (relation to tags, multiple)
- [ ] Field: view_count (number)
- [ ] Field: created_by (relation to users)
- [ ] System fields present

**Commit**: No changes needed if verification passes

### Step 7: Create progress tracking file

**Action**: Document implementation progress

**File**: `docs/active/work/T-002-01/progress.md`

**Content**:
- Steps completed
- Any deviations from plan
- Verification results

**Commit**: "docs: add progress tracking for T-002-01"

### Step 8: Update ticket status

**Action**: Update ticket frontmatter to reflect completion

**File**: `docs/active/tickets/T-002-01.md`

**Changes**:
```yaml
phase: done
status: done
```

**Commit**: "feat: complete T-002-01 collection definitions"

## Testing Strategy

### Manual Testing Flow

1. **Clean State Test**:
   - Delete pb_data
   - Run application
   - Verify collections created
   - Check all fields

2. **Idempotency Test**:
   - Restart application
   - Verify no errors
   - Collections unchanged

3. **Field Validation Test**:
   - Try creating records via admin UI
   - Test required fields
   - Test unique constraint on slug

### Automated Testing (Future)

Not in current scope, but structure allows for:
```go
func TestCollectionMigration(t *testing.T) {
    app := testutils.NewTestApp()
    migrations.Register(app)
    // Assert collections exist
}
```

## Error Recovery

### Potential Issues and Solutions

1. **Migration fails on startup**:
   - Check error message
   - Fix field configuration
   - Delete pb_data and retry

2. **Collections missing fields**:
   - Verify field definitions in code
   - Check for typos in field names
   - Ensure all Add() calls present

3. **Relations not working**:
   - Confirm tags saved before links
   - Check collection ID retrieval
   - Verify users collection exists

4. **Unique constraint not enforced**:
   - Check Unique flag on slug field
   - Verify in admin UI settings
   - Test with duplicate values

## Rollback Plan

If issues arise:

1. **Revert code changes**:
   ```bash
   git checkout main
   git reset --hard HEAD~n  # n = number of commits
   ```

2. **Clean database**:
   ```bash
   rm -rf pb_data
   ```

3. **Fix issues and retry**

## Success Criteria

All acceptance criteria must be met:

- [x] Migration file(s) exist in `migrations/` package
- [x] tags collection with correct fields
- [x] links collection with correct fields
- [x] Migrations registered in main.go
- [x] Run automatically on startup
- [x] Collections visible in admin UI

## Commit Strategy

Small, atomic commits:

1. Fix dependency (go.mod)
2. Add migration package
3. Register in main.go
4. Documentation updates
5. Final ticket status update

Each commit should compile and run.

## Time Estimates

- Step 1-2: 2 minutes (setup)
- Step 3: 10 minutes (main implementation)
- Step 4: 2 minutes (integration)
- Step 5-6: 5 minutes (testing)
- Step 7-8: 3 minutes (documentation)

Total: ~22 minutes

## Dependencies Verification

Before starting:
- [x] T-001-01 completed (PocketBase initialized)
- [x] main.go exists and works
- [x] go.mod present
- [x] Admin UI accessible

All dependencies satisfied, ready to implement.