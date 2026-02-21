# Design: T-005-03 GitHub Actions CI Setup

## Design Problem Statement

Create a GitHub Actions workflow that validates the project builds cleanly in a fresh environment while maintaining compatibility with the existing Makefile-based build system. The main challenge is replacing Flox environment management with direct tool installation for CI execution.

## Design Alternatives Evaluation

### Alternative 1: Direct Tool Installation (Recommended)
**Approach**: Install Go 1.26 and Node.js 24 directly in GitHub Actions using standard setup actions, modify Makefile targets to work without Flox wrapper.

**Pros**:
- Standard GitHub Actions pattern
- Fast execution (no Flox overhead)
- Excellent caching support via GitHub Actions
- Portable to other CI systems
- Simple maintenance

**Cons**:
- Requires modifying existing Makefile to support both Flox and CI environments
- Slight divergence between local dev (Flox) and CI (direct) environments

**Implementation Complexity**: Low

### Alternative 2: Flox in CI
**Approach**: Install Flox in GitHub Actions and use existing build commands unchanged.

**Pros**:
- Zero changes to existing build system
- Perfect environment parity between dev and CI

**Cons**:
- Flox not officially supported in GitHub Actions
- Complex setup and installation
- Slower execution due to environment activation overhead
- Dependency on external tool availability
- Limited caching opportunities

**Implementation Complexity**: High

### Alternative 3: Docker-based Build
**Approach**: Create Dockerfile with Flox and build environment, run builds in container.

**Pros**:
- Complete environment isolation
- Reproducible builds
- Could reuse for local development

**Cons**:
- Significant infrastructure changes required
- Slower execution due to container overhead
- Additional complexity for maintenance
- Over-engineering for current requirements

**Implementation Complexity**: Very High

## Recommended Design: Direct Tool Installation

Based on evaluation criteria (simplicity, maintainability, performance), **Alternative 1** is selected.

### Design Rationale

1. **Alignment with GitHub Actions Best Practices**: Using `actions/setup-go` and `actions/setup-node` is the standard approach
2. **Performance**: Direct tool installation is faster than environment wrappers
3. **Maintainability**: Fewer moving parts, standard patterns
4. **Caching**: Excellent support for Go modules and npm dependencies
5. **Future-proofing**: Works with any CI system, not tied to Flox availability

### Key Design Decisions

#### 1. Makefile Compatibility Strategy
**Decision**: Add conditional Flox detection to existing Makefile targets

**Implementation**: Modify `backend` and `test` targets to use direct `go` command when Flox is not available:
```make
backend: frontend
    @if command -v flox >/dev/null 2>&1; then \
        flox activate -- go build -o $(BINARY_NAME); \
    else \
        go build -o $(BINARY_NAME); \
    fi

test:
    @if command -v flox >/dev/null 2>&1; then \
        flox activate -- go test ./...; \
    else \
        go test ./...; \
    fi
```

**Rationale**: Maintains backward compatibility while enabling CI execution without changes to workflow commands.

#### 2. Workflow Trigger Configuration
**Decision**: Trigger on push to `main` and all pull requests

**Configuration**:
```yaml
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
```

**Rationale**: Provides validation for all changes while avoiding excessive CI runs on feature branches.

#### 3. Caching Strategy
**Decision**: Implement caching for Go modules and npm dependencies

**Go Module Caching**:
- Cache key: Hash of `go.sum`
- Cache path: Go module cache directory
- Fallback: Partial match on `go.sum` prefix

**npm Dependency Caching**:
- Cache key: Hash of `frontend/package-lock.json`
- Cache path: `frontend/node_modules`
- Fallback: Partial match on package-lock prefix

**Rationale**: Reduces build time from ~2 minutes to ~30-45 seconds on cache hits.

#### 4. Build Validation Strategy
**Decision**: Use existing Makefile targets with additional CI-specific validations

**Validation Steps**:
1. Execute `make build` (includes frontend, backend, validate-build)
2. Execute `make test`
3. CI-specific: Verify no uncommitted changes after build
4. CI-specific: Attempt to start binary (smoke test)

**Rationale**: Leverages existing validation logic while adding CI-appropriate checks.

## Workflow Architecture

### Job Structure
**Single Job Design**: All build steps in one job for simplicity

**Job Sequence**:
1. **Setup**: Checkout code, setup Go 1.26, setup Node.js 24
2. **Cache**: Restore cached dependencies
3. **Build**: Execute full build pipeline via Makefile
4. **Test**: Execute test suite via Makefile
5. **Validate**: Perform CI-specific validations

### Environment Configuration

#### Go Configuration
```yaml
- uses: actions/setup-go@v4
  with:
    go-version: '1.26'
    cache: true
```

#### Node.js Configuration
```yaml
- uses: actions/setup-node@v4
  with:
    node-version: '24'
    cache: 'npm'
    cache-dependency-path: 'frontend/package-lock.json'
```

### Error Handling Strategy
**Build Failures**: Workflow fails immediately on build errors
**Test Failures**: Workflow fails after collecting test results
**Cache Failures**: Continue without cache (graceful degradation)

## Implementation Requirements

### Makefile Modifications Required
1. **Conditional Flox Detection**: Add environment detection to `backend` and `test` targets
2. **Error Handling**: Ensure proper exit codes for CI
3. **Path Handling**: Ensure relative paths work correctly in CI environment

### Workflow File Structure
**File Location**: `.github/workflows/ci.yml`
**Key Components**:
- Workflow name and triggers
- Job definition with Ubuntu runner
- Step sequence with proper error handling
- Caching configuration
- Build and test execution

### Dependencies and Constraints

#### Version Pinning
- Go version: Exactly `1.26` (per ticket requirements)
- Node version: Exactly `24` (per ticket requirements and package.json)
- Actions versions: Pin to specific major versions with auto-updates enabled

#### Build Environment
- Runner: `ubuntu-latest`
- No additional system dependencies required
- Standard GitHub Actions environment sufficient

## Success Criteria

### Functional Requirements
1. **Build Success**: Complete build pipeline executes successfully
2. **Test Execution**: All tests pass in CI environment
3. **Fresh Environment**: Build works from clean checkout
4. **Error Detection**: Failures properly reported and block PR merge

### Performance Requirements
1. **Build Time**: Complete workflow under 3 minutes (fresh), under 1 minute (cached)
2. **Resource Usage**: Standard GitHub Actions limits sufficient
3. **Caching Effectiveness**: >50% build time reduction on cache hits

### Quality Requirements
1. **Reliability**: Consistent results across multiple runs
2. **Maintainability**: Clear workflow structure, commented configuration
3. **Debugging**: Adequate logging for troubleshooting failures

## Migration Strategy

### Implementation Phases
1. **Phase 1**: Modify Makefile for Flox compatibility
2. **Phase 2**: Create GitHub Actions workflow
3. **Phase 3**: Test workflow with sample PR
4. **Phase 4**: Enable required status checks

### Rollback Plan
If CI implementation fails:
1. Remove `.github/workflows/ci.yml`
2. Revert Makefile changes if necessary
3. Build system returns to development-only validation

### Testing Approach
1. **Local Testing**: Verify Makefile changes work with and without Flox
2. **CI Testing**: Create test branch and validate workflow execution
3. **Integration Testing**: Test with actual PR to validate trigger behavior

This design provides a pragmatic solution that maintains existing build system functionality while enabling robust CI validation through GitHub Actions.