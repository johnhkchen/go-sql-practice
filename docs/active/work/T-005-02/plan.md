# Plan: Build Orchestration for Go-SQL-Practice

## Implementation Steps

### Step 1: Foundation Setup
**Objective**: Add configuration variables and improve existing target structure
**Files**: `Makefile` (lines 1-10)
**Duration**: 15 minutes
**Dependencies**: None

**Actions**:
1. Add configuration variables at top of Makefile:
   - `BINARY_NAME := go-sql-practice`
   - `FRONTEND_DIR := frontend`
   - `DIST_DIR := $(FRONTEND_DIR)/dist`
   - `SERVER_PORT := 127.0.0.1:8090`

2. Update existing `.PHONY` declaration to include new targets:
   - Add `help`, `validate-build` to `.PHONY` list

**Verification**:
- `make build` continues to work unchanged
- Variables can be referenced in subsequent targets
- No functional changes to existing behavior

**Commit**: "feat: add Makefile configuration variables and expand .PHONY declaration"

### Step 2: Enhanced Clean Target
**Objective**: Improve cleanup to remove all build artifacts including pb_data
**Files**: `Makefile` (clean target)
**Duration**: 10 minutes
**Dependencies**: Step 1

**Actions**:
1. Expand clean target to include:
   ```makefile
   clean:
       @echo "Cleaning build artifacts..."
       rm -rf $(DIST_DIR)
       rm -f $(BINARY_NAME)
       rm -rf pb_data
       @echo "Clean complete"
   ```

**Verification**:
- Create test artifacts: `touch go-sql-practice`, `mkdir -p frontend/dist`, `mkdir -p pb_data`
- Run `make clean`
- Verify all artifacts are removed
- Verify clean output messages appear

**Commit**: "feat: enhance clean target to remove all build artifacts including pb_data"

### Step 3: Backend Target Validation
**Objective**: Add build validation to backend target to prevent silent failures
**Files**: `Makefile` (backend target)
**Duration**: 20 minutes
**Dependencies**: Step 1

**Actions**:
1. Add pre-build validation to backend target:
   ```makefile
   backend: frontend
       @echo "Building backend..."
       @if [ ! -d "$(DIST_DIR)" ]; then \
           echo "Error: Frontend not built. Run 'make frontend' first."; \
           exit 1; \
       fi
       go build -o $(BINARY_NAME)
       @echo "Backend build complete"
   ```

**Verification**:
- Test failure case: `make clean && make backend` (should fail with clear message)
- Test success case: `make frontend && make backend` (should succeed)
- Verify error message is clear and actionable
- Verify binary is created on success

**Commit**: "feat: add frontend validation to backend build target"

### Step 4: Frontend Target Enhancement
**Objective**: Add output validation to frontend target
**Files**: `Makefile` (frontend target)
**Duration**: 15 minutes
**Dependencies**: Step 1

**Actions**:
1. Enhance frontend target with validation:
   ```makefile
   frontend:
       @echo "Building frontend..."
       cd $(FRONTEND_DIR) && npm ci && npm run build
       @if [ ! -d "$(DIST_DIR)/client" ]; then \
           echo "Error: Frontend build failed - client directory not created"; \
           exit 1; \
       fi
       @echo "Frontend build complete"
   ```

**Verification**:
- Test with clean frontend: `rm -rf frontend/dist && make frontend`
- Verify `frontend/dist/client` directory exists after success
- Test failure handling (simulate by making frontend/dist read-only)
- Verify clear error messages

**Commit**: "feat: add output validation to frontend build target"

### Step 5: Help Target Implementation
**Objective**: Add comprehensive help documentation for all targets
**Files**: `Makefile` (new help target)
**Duration**: 15 minutes
**Dependencies**: None

**Actions**:
1. Add help target:
   ```makefile
   help:
       @echo "Available targets:"
       @echo "  build       - Build frontend and backend (full pipeline)"
       @echo "  frontend    - Install deps and build Astro frontend"
       @echo "  backend     - Build Go binary (requires frontend)"
       @echo "  clean       - Remove all build artifacts"
       @echo "  dev         - Start development server"
       @echo "  test        - Run Go tests"
       @echo "  validate    - Validate build artifacts"
       @echo "  help        - Show this help message"
   ```

2. Make help the default target:
   ```makefile
   .DEFAULT_GOAL := help
   ```

**Verification**:
- Run `make help` and verify all targets are documented
- Run `make` (no target) and verify help is displayed
- Verify descriptions are accurate and helpful

**Commit**: "feat: add comprehensive help target with usage documentation"

### Step 6: Build Validation Target
**Objective**: Add standalone build validation for CI/automation use
**Files**: `Makefile` (new validate-build target)
**Duration**: 20 minutes
**Dependencies**: Steps 1-4

**Actions**:
1. Add validate-build target:
   ```makefile
   validate-build: build
       @echo "Validating build artifacts..."
       @if [ ! -f "$(BINARY_NAME)" ]; then \
           echo "Error: Binary not created"; exit 1; \
       fi
       @if [ ! -d "$(DIST_DIR)" ]; then \
           echo "Error: Frontend assets missing"; exit 1; \
       fi
       @echo "Build validation successful"
   ```

**Verification**:
- Test full pipeline: `make clean && make validate-build`
- Test validation failure: `make build && rm go-sql-practice && make validate-build`
- Verify exit codes (0 for success, 1 for failure)
- Test as part of main build: `make build` should call validate-build

**Commit**: "feat: add build validation target for CI/automation"

### Step 7: Enhanced Dev Target
**Objective**: Improve development server startup with proper build dependency
**Files**: `Makefile` (dev target)
**Duration**: 15 minutes
**Dependencies**: Steps 1-6

**Actions**:
1. Enhance dev target:
   ```makefile
   dev: build
       @echo "Starting development server..."
       @echo "Server will be available at http://$(SERVER_PORT)"
       @echo "Press Ctrl+C to stop"
       ./$(BINARY_NAME) serve --http="$(SERVER_PORT)"
   ```

**Verification**:
- Test with clean build: `make clean && make dev`
- Verify build runs automatically before server starts
- Test server accessibility at configured port
- Verify clean shutdown with Ctrl+C
- Test port configuration via SERVER_PORT variable

**Commit**: "feat: enhance dev target with automatic build and better UX"

### Step 8: Error Handling Improvements
**Objective**: Add consistent error handling and user feedback across all targets
**Files**: `Makefile` (all targets)
**Duration**: 25 minutes
**Dependencies**: Steps 1-7

**Actions**:
1. Add consistent error messages to all targets
2. Ensure proper exit codes for automation
3. Add progress indicators for long-running operations
4. Add inline comments for complex logic

**Verification**:
- Test each target individually for proper error handling
- Verify error messages are clear and actionable
- Test all targets in sequence: `make clean && make help && make build && make dev`
- Verify exit codes with `echo $?` after each operation

**Commit**: "feat: improve error handling and user feedback across all targets"

### Step 9: Integration Testing
**Objective**: Comprehensive testing of complete build pipeline
**Files**: None (testing only)
**Duration**: 30 minutes
**Dependencies**: Steps 1-8

**Actions**:
1. Test full clean build cycle:
   ```bash
   make clean
   make build
   make validate-build
   make test
   make dev  # test in background
   ```

2. Test individual target independence:
   - Each target should work when run standalone where applicable
   - Verify proper dependency handling

3. Test error scenarios:
   - Missing Node.js dependencies
   - Frontend build failures
   - Go build failures
   - Missing dependencies

**Verification**:
- All targets work as documented
- Error messages are helpful
- Build artifacts are created correctly
- Development server starts and serves embedded assets
- All acceptance criteria met

**Commit**: "test: verify complete build orchestration pipeline"

## Testing Strategy

### Unit Testing Approach

**Target-Level Testing**: Each Makefile target tested independently
- `make help` - Output format and completeness
- `make clean` - Artifact removal verification
- `make frontend` - Build output validation
- `make backend` - Dependency checking and binary creation
- `make build` - Pipeline orchestration
- `make validate-build` - Validation logic
- `make dev` - Server startup
- `make test` - Go test execution

**Error Condition Testing**: Verify failure modes and error messages
- Missing dependencies (Node.js, Go)
- Build failures (syntax errors, missing files)
- Validation failures (missing artifacts)
- Permission issues (read-only directories)

### Integration Testing Approach

**Full Pipeline Testing**: End-to-end build validation
1. Clean environment setup
2. Complete build execution
3. Artifact validation
4. Development server functionality
5. Cleanup verification

**Dependency Testing**: Verify target relationships
- `backend` depends on `frontend`
- `build` orchestrates both
- `validate-build` verifies `build` outputs
- `dev` depends on `build`

**Cross-Platform Testing**: Unix-like systems (Linux, macOS)
- Path handling (forward slashes)
- Shell command compatibility
- File permission handling

### Verification Criteria

**Functional Requirements**: All acceptance criteria satisfied
- ✅ `make build` runs full pipeline
- ✅ `make dev` starts development server
- ✅ `make test` runs Go tests
- ✅ `make clean` removes all artifacts
- ✅ Each target runs independently
- ✅ Phony targets and comments included

**Quality Requirements**: Professional build system standards
- Clear error messages with actionable guidance
- Consistent user feedback and progress indicators
- Proper exit codes for automation integration
- Comprehensive help documentation
- Backward compatibility with existing workflow

**Performance Requirements**: Reasonable build times
- Frontend build: <2 minutes (npm ci + build)
- Backend build: <30 seconds (Go compilation)
- Full pipeline: <3 minutes total
- Cleanup: <5 seconds

## Risk Mitigation

### High-Risk Areas

**Frontend Build Dependencies**: npm ci can fail with network issues
- **Mitigation**: Clear error messages, retry guidance
- **Verification**: Test with network simulation/poor connectivity

**Go Embed Validation**: Silent failures when frontend assets missing
- **Mitigation**: Explicit directory checks before go build
- **Verification**: Test with missing/corrupted frontend assets

**Development Server Conflicts**: Port conflicts with existing servers
- **Mitigation**: Configurable SERVER_PORT variable
- **Verification**: Test with ports already in use

### Medium-Risk Areas

**Cross-Platform Compatibility**: Makefile Unix-specific
- **Mitigation**: Document Unix requirement clearly
- **Verification**: Test on Linux and macOS

**Build Artifact Cleanup**: Incomplete cleanup leaving state
- **Mitigation**: Comprehensive clean target testing
- **Verification**: Multiple clean/build cycles

### Recovery Procedures

**Build Failure Recovery**:
1. `make clean` to reset state
2. Check error message for specific guidance
3. Verify dependencies (Node.js, Go versions)
4. Retry individual targets: `make frontend`, then `make backend`

**Development Environment Recovery**:
1. Kill any running servers: `pkill go-sql-practice`
2. Clean build artifacts: `make clean`
3. Fresh build: `make build`
4. Restart development: `make dev`

This plan provides a systematic approach to implementing the enhanced build orchestration while maintaining reliability and providing comprehensive testing coverage.