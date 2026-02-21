# Plan for T-010-05: Infrastructure Cleanup

## Implementation Steps

### Step 1: Clean Git-Tracked Test Artifacts
**Actions**:
1. Update `.gitignore` to add `*.test` pattern
2. Update `.gitignore` to add `routes/pb_data/` pattern
3. Remove `routes.test` from git tracking and delete file
4. Verify no other `.test` files exist

**Verification**:
- Run `git status` to confirm routes.test is deleted
- Run `find . -name "*.test"` to verify no test binaries remain
- Create a test binary and verify it's ignored by git

**Commit**: "chore: add .gitignore entries and remove tracked test artifacts"

### Step 2: Remove Stray Directories
**Actions**:
1. Delete `/frontend/frontend/` directory recursively
2. Delete `/routes/pb_data/` directory recursively
3. Verify main frontend source remains intact

**Verification**:
- Confirm `/frontend/src/` still exists and is complete
- Verify `routes/` directory has no pb_data subdirectory
- Run `ls -la` to confirm cleanup

**Commit**: "chore: remove duplicate directories and test artifacts"

### Step 3: Fix Makefile Documentation
**Actions**:
1. Update line 74 help text from "validate" to "validate-build"
2. Test `make help` displays correct information

**Verification**:
- Run `make help` and verify output matches targets
- Confirm all listed targets exist

**Commit**: "fix: correct Makefile help text for validate-build target"

### Step 4: Add Quality Check Targets
**Actions**:
1. Update `.PHONY` line to include `lint vet`
2. Add `lint` target with gofmt check
3. Add `vet` target with go vet command
4. Test both new targets work correctly

**Verification**:
- Run `make lint` - should pass if code is formatted
- Run `make vet` - should analyze code successfully
- Intentionally misformat a file and verify lint catches it

**Commit**: "feat: add lint and vet targets to Makefile"

### Step 5: Create Frontend Configuration Example
**Actions**:
1. Create `/frontend/.env.example` with PUBLIC_API_URL
2. Verify it matches the working `.env` file

**Verification**:
- Compare `.env` and `.env.example` content
- Confirm both use port 8090

**Commit**: "docs: add .env.example for frontend configuration"

### Step 6: Update Project Documentation
**Actions**:
1. Replace placeholder in CLAUDE.md with actual description
2. Keep description concise and accurate

**Verification**:
- Read CLAUDE.md to confirm placeholder is gone
- Verify description accurately reflects the project

**Commit**: "docs: update CLAUDE.md with project description"

### Step 7: Enhance Error Logging
**Actions**:
1. Add logging before each error return in routes/stats.go
2. Use app.Logger().Error() with descriptive context
3. Maintain existing error response format

**Verification**:
- Run the application
- Trigger a stats endpoint error (if possible)
- Verify logs show detailed error but response stays generic

**Commit**: "feat: add server-side error logging to stats endpoint"

### Step 8: Update CI Workflow
**Actions**:
1. Update `actions/setup-go@v4` to `v5`
2. Add go vet step after tests
3. Add gofmt check step
4. Add npm audit step after frontend build

**Verification**:
- Review CI YAML for syntax correctness
- Push to branch to trigger CI run
- Monitor CI execution for new steps

**Commit**: "ci: update workflow with quality checks and dependency audit"

## Testing Strategy

### Local Testing
- Run full build: `make clean && make build`
- Execute tests: `make test`
- Try new targets: `make lint && make vet`
- Start dev server: `make dev`
- Verify frontend loads correctly

### CI Testing
- Create feature branch for changes
- Push after each commit group
- Monitor GitHub Actions for success
- Fix any CI failures before proceeding

### Manual Verification
- Check git status after .gitignore updates
- Inspect directories after deletion
- Read help output after Makefile changes
- Review logs after adding error logging

## Rollback Scenarios

### If Build Breaks
- Git revert the specific commit
- Focus on Makefile changes as likely cause
- Test each target individually

### If CI Fails
- Check new go vet and gofmt steps first
- Verify npm audit threshold isn't too strict
- Can temporarily comment out new steps

### If Frontend Issues
- Verify frontend/dist still builds correctly
- Check embed.go still finds assets
- Ensure .env.example doesn't override .env

## Success Criteria

1. No .test files in git
2. No duplicate directories in project
3. make help shows accurate information
4. make lint and make vet work
5. CI runs with all quality checks
6. .env.example documents configuration
7. CLAUDE.md has real description
8. Errors logged server-side
9. All existing functionality preserved
10. Build and tests still pass