# Plan: T-005-03 GitHub Actions CI Setup

## Implementation Strategy

This plan implements GitHub Actions CI through a two-stage approach: first modifying the Makefile for CI compatibility, then creating the workflow configuration. Each step is independently testable and can be committed atomically.

## Implementation Steps

### Step 1: Prepare Build System for CI
**Objective**: Modify Makefile to support both Flox (local development) and direct tool execution (CI)
**Duration**: 15 minutes
**Dependencies**: None

**Tasks**:
1. Add conditional Flox detection to `backend` target in Makefile
2. Add conditional Flox detection to `test` target in Makefile
3. Test locally with and without Flox available

**Implementation Details**:
```make
backend: frontend
	@if command -v flox >/dev/null 2>&1; then \
		echo "Building with Flox environment..."; \
		flox activate -- go build -o $(BINARY_NAME); \
	else \
		echo "Building with system Go..."; \
		go build -o $(BINARY_NAME); \
	fi

test:
	@if command -v flox >/dev/null 2>&1; then \
		flox activate -- go test ./...; \
	else \
		go test ./...; \
	fi
```

**Testing Strategy**:
- Local test with Flox: `make clean && make build && make test`
- Local test without Flox: Temporarily rename flox binary, repeat tests
- Verify all existing functionality preserved

**Success Criteria**:
- Build succeeds in both Flox and non-Flox environments
- Test execution works in both environments
- No regression in local development workflow

**Commit Point**: "feat: add Makefile CI compatibility for direct tool execution"

### Step 2: Create GitHub Actions Workflow Structure
**Objective**: Set up workflow file and directory structure
**Duration**: 10 minutes
**Dependencies**: None

**Tasks**:
1. Create `.github/workflows/` directory
2. Create basic `ci.yml` file with workflow metadata
3. Verify GitHub recognizes the workflow (draft PR or push to branch)

**Implementation Details**:
```yaml
name: CI
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Placeholder
        run: echo "Workflow structure created"
```

**Testing Strategy**:
- Create test branch with workflow file
- Push to trigger workflow execution
- Verify workflow appears in GitHub Actions tab

**Success Criteria**:
- Workflow file recognized by GitHub
- Basic workflow executes successfully
- Repository Actions tab shows workflow

**Commit Point**: "ci: add GitHub Actions workflow structure"

### Step 3: Implement Environment Setup
**Objective**: Configure Go and Node.js environments with caching
**Duration**: 20 minutes
**Dependencies**: Step 2

**Tasks**:
1. Add repository checkout step
2. Configure Go 1.26 setup with module caching
3. Configure Node.js 24 setup with npm caching
4. Test environment setup without build execution

**Implementation Details**:
```yaml
steps:
  - uses: actions/checkout@v4

  - name: Set up Go 1.26
    uses: actions/setup-go@v4
    with:
      go-version: '1.26'
      cache: true

  - name: Set up Node.js 24
    uses: actions/setup-node@v4
    with:
      node-version: '24'
      cache: 'npm'
      cache-dependency-path: 'frontend/package-lock.json'
```

**Testing Strategy**:
- Add environment verification steps (go version, node version)
- Test cache creation and restoration
- Measure setup time with and without cache

**Success Criteria**:
- Go 1.26 correctly installed and cached
- Node.js 24 correctly installed and cached
- Environment verification passes

**Commit Point**: "ci: add Go 1.26 and Node.js 24 environment setup with caching"

### Step 4: Implement Build Pipeline
**Objective**: Execute the complete build process using existing Makefile targets
**Duration**: 15 minutes
**Dependencies**: Step 1, Step 3

**Tasks**:
1. Add frontend build step (`make frontend`)
2. Add backend build step (`make backend`)
3. Add build validation step (`make validate-build`)
4. Test complete build pipeline in CI

**Implementation Details**:
```yaml
- name: Build frontend
  run: make frontend

- name: Build backend
  run: make backend

- name: Validate build
  run: make validate-build
```

**Testing Strategy**:
- Test build with cache miss (fresh dependencies)
- Test build with cache hit (restored dependencies)
- Verify build artifacts created correctly
- Test build failure scenarios

**Success Criteria**:
- Frontend builds successfully with npm cache utilization
- Backend builds successfully with Go module cache
- Build validation passes
- Build artifacts available for download

**Commit Point**: "ci: implement complete build pipeline with validation"

### Step 5: Implement Test Execution
**Objective**: Execute Go test suite and report results
**Duration**: 10 minutes
**Dependencies**: Step 1, Step 4

**Tasks**:
1. Add test execution step (`make test`)
2. Configure test result reporting
3. Test with passing and failing test scenarios

**Implementation Details**:
```yaml
- name: Run tests
  run: make test
```

**Testing Strategy**:
- Execute tests with all passing
- Temporarily introduce failing test to verify failure detection
- Verify test output properly captured and displayed

**Success Criteria**:
- Tests execute successfully in CI environment
- Test failures properly reported and fail the workflow
- Test output visible in workflow logs

**Commit Point**: "ci: add Go test execution with proper error reporting"

### Step 6: Add CI-Specific Validations
**Objective**: Add additional validations specific to CI environment
**Duration**: 10 minutes
**Dependencies**: Step 5

**Tasks**:
1. Add binary smoke test (verify binary starts)
2. Add build cleanliness check (no uncommitted changes)
3. Test validation steps

**Implementation Details**:
```yaml
- name: Smoke test binary
  run: |
    ./go-sql-practice --help || (echo "Binary smoke test failed" && exit 1)

- name: Check for uncommitted changes
  run: |
    if [ -n "$(git status --porcelain)" ]; then
      echo "Build process generated uncommitted changes:"
      git status --porcelain
      exit 1
    fi
```

**Testing Strategy**:
- Test binary smoke test with valid binary
- Test build cleanliness with clean and dirty builds
- Verify proper error reporting

**Success Criteria**:
- Binary smoke test passes with successful binary
- Build cleanliness check detects any uncommitted changes
- Validation failures properly reported

**Commit Point**: "ci: add binary smoke test and build cleanliness validation"

### Step 7: Optimize Workflow Performance
**Objective**: Enhance caching strategy and optimize execution time
**Duration**: 15 minutes
**Dependencies**: Step 6

**Tasks**:
1. Fine-tune cache keys for optimal hit rates
2. Add cache statistics reporting
3. Optimize step execution order
4. Add workflow timing analysis

**Implementation Details**:
```yaml
- name: Cache Go modules
  uses: actions/cache@v3
  with:
    path: |
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-

- name: Cache npm dependencies
  uses: actions/cache@v3
  with:
    path: frontend/node_modules
    key: ${{ runner.os }}-npm-${{ hashFiles('frontend/package-lock.json') }}
    restore-keys: |
      ${{ runner.os }}-npm-
```

**Testing Strategy**:
- Measure build time with cold cache
- Measure build time with warm cache
- Test cache invalidation on dependency changes
- Monitor cache hit/miss rates

**Success Criteria**:
- Cache hit provides >50% build time reduction
- Cache invalidation works correctly on dependency changes
- Workflow completes under 3 minutes (cold), under 1 minute (warm)

**Commit Point**: "ci: optimize workflow performance with enhanced caching strategy"

### Step 8: Integration Testing and Documentation
**Objective**: Comprehensive testing and workflow documentation
**Duration**: 20 minutes
**Dependencies**: Step 7

**Tasks**:
1. Test complete workflow with actual PR
2. Test workflow triggers (push to main, PR creation)
3. Add workflow documentation comments
4. Update project documentation if needed

**Implementation Details**:
```yaml
name: CI
# This workflow validates the project builds cleanly in a fresh environment
# Triggers on push to main and pull requests targeting main
# Uses existing Makefile targets for consistency with local development
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
```

**Testing Strategy**:
- Create actual PR to test PR trigger
- Push to main branch to test push trigger
- Test workflow with various PR scenarios (draft, ready for review)
- Verify status checks properly reported

**Success Criteria**:
- All workflow triggers function correctly
- PR status checks integrate properly with GitHub UI
- Workflow provides clear success/failure feedback
- Documentation is clear and helpful

**Commit Point**: "ci: finalize GitHub Actions CI with comprehensive testing and documentation"

## Testing Strategy

### Unit Testing (Per Step)
Each step includes specific testing requirements to verify implementation correctness before proceeding to the next step.

### Integration Testing
**Local-CI Parity Testing**:
- Verify `make build` produces identical results locally (with Flox) and in CI
- Verify `make test` produces consistent results across environments
- Test edge cases: missing dependencies, corrupted cache, network issues

**End-to-End Testing**:
- Create test PR with intentional changes to verify full pipeline
- Test workflow with failing tests to verify error reporting
- Test workflow with build failures to verify proper failure handling

**Performance Testing**:
- Measure baseline workflow execution time
- Test cache effectiveness across multiple runs
- Monitor resource usage and optimization opportunities

### Rollback Testing
- Verify local development unaffected if CI changes removed
- Test Makefile rollback scenarios
- Verify no permanent changes to development environment

## Verification Criteria

### Functional Verification
1. **Build Success**: Complete build pipeline executes successfully in CI
2. **Test Execution**: All tests pass consistently in CI environment
3. **Environment Isolation**: Clean builds work from fresh checkout
4. **Error Detection**: Build and test failures properly detected and reported

### Performance Verification
1. **Build Time**: Fresh build completes under 3 minutes
2. **Cached Build**: Cached build completes under 1 minute
3. **Cache Efficiency**: Cache provides >50% time reduction on hits
4. **Resource Usage**: Workflow stays within GitHub Actions standard limits

### Integration Verification
1. **Trigger Accuracy**: Workflow triggers on correct events only
2. **Status Reporting**: Build status correctly reported to GitHub PR interface
3. **Branch Protection**: Compatible with GitHub branch protection rules
4. **Developer Experience**: Clear feedback on build success/failure

### Quality Verification
1. **Reliability**: Workflow produces consistent results across multiple runs
2. **Maintainability**: Workflow configuration is clear and well-documented
3. **Debuggability**: Sufficient logging for troubleshooting failures
4. **Compatibility**: Maintains full compatibility with local development workflow

## Risk Mitigation

### Build Failure Risks
**Risk**: Workflow fails due to environment differences
**Mitigation**: Extensive local testing of Makefile changes before CI implementation

**Risk**: Dependency caching issues cause build failures
**Mitigation**: Graceful fallback to full dependency installation on cache failures

### Performance Risks
**Risk**: Workflow execution time exceeds acceptable limits
**Mitigation**: Incremental optimization with performance monitoring at each step

**Risk**: Cache storage exceeds GitHub Actions limits
**Mitigation**: Monitor cache usage and implement cache cleanup if needed

### Integration Risks
**Risk**: CI requirements conflict with local development workflow
**Mitigation**: Maintain backward compatibility and test local workflow after each change

**Risk**: GitHub Actions service limitations affect workflow reliability
**Mitigation**: Design workflow within documented GitHub Actions best practices and limits

## Success Metrics

### Implementation Metrics
- All 8 implementation steps completed successfully
- Zero regressions in local development workflow
- All test scenarios pass
- Complete workflow executes successfully

### Performance Metrics
- Build time: <3 minutes (fresh), <1 minute (cached)
- Cache hit rate: >80% for typical development workflows
- Test execution time: <30 seconds
- Overall workflow reliability: >95% success rate

### Quality Metrics
- Workflow configuration passes GitHub Actions linting
- All acceptance criteria from ticket met
- Documentation complete and accurate
- Developer feedback positive

This plan provides a structured approach to implementing GitHub Actions CI while maintaining development workflow compatibility and ensuring robust validation of all project changes.