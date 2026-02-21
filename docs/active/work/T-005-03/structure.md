# Structure: T-005-03 GitHub Actions CI Setup

## File-Level Changes Overview

This implementation requires creating one new file and modifying one existing file to enable GitHub Actions CI while maintaining local development compatibility.

## New Files to Create

### `.github/workflows/ci.yml`
**Purpose**: GitHub Actions workflow definition for continuous integration
**Location**: `.github/workflows/ci.yml` (new directory and file)
**Size**: ~80-100 lines
**Content**: Complete CI workflow with setup, caching, build, and test steps

## Modified Files

### `Makefile`
**Purpose**: Add conditional Flox detection for CI compatibility
**Location**: `Makefile` (existing file modification)
**Changes**: Modify `backend` and `test` targets to work without Flox
**Impact**: 4-6 lines modified, maintaining backward compatibility

## Directory Structure Changes

### New Directory Creation
```
.github/                    # New directory
└── workflows/             # New subdirectory
    └── ci.yml             # New workflow file
```

**Rationale**: Standard GitHub Actions convention for workflow storage

### No Directory Removals
No existing directories or files will be removed or relocated.

## Component Architecture Definition

### GitHub Actions Workflow Components

#### 1. Workflow Metadata
**Location**: Top of `.github/workflows/ci.yml`
**Purpose**: Define workflow name, triggers, and permissions
**Structure**:
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
```

#### 2. Environment Setup Section
**Purpose**: Configure Go and Node.js environments
**Components**:
- Checkout action
- Go 1.26 setup with caching
- Node.js 24 setup with caching
- Dependency cache restoration

#### 3. Build Execution Section
**Purpose**: Execute the build pipeline
**Components**:
- Frontend build via `make frontend`
- Backend build via `make backend`
- Build validation via `make validate-build`

#### 4. Test Execution Section
**Purpose**: Run the test suite
**Components**:
- Test execution via `make test`
- Test result reporting

#### 5. CI-Specific Validation Section
**Purpose**: Additional validations for CI environment
**Components**:
- Binary smoke test
- Build artifact verification

### Makefile Modification Architecture

#### 1. Flox Detection Function
**Purpose**: Determine if Flox is available in current environment
**Implementation**: Shell command detection pattern
**Usage**: Shared by multiple targets

#### 2. Backend Target Modification
**Current**: `flox activate -- go build -o $(BINARY_NAME)`
**New Structure**:
```make
backend: frontend
    @if command -v flox >/dev/null 2>&1; then \
        flox activate -- go build -o $(BINARY_NAME); \
    else \
        go build -o $(BINARY_NAME); \
    fi
```

#### 3. Test Target Modification
**Current**: `flox activate -- go test ./...`
**New Structure**:
```make
test:
    @if command -v flox >/dev/null 2>&1; then \
        flox activate -- go test ./...; \
    else \
        go test ./...; \
    fi
```

## Public Interface Definition

### Workflow Interface
**Trigger Interface**:
- Push to `main` branch → Full CI execution
- Pull request targeting `main` → Full CI execution
- Manual workflow dispatch → Full CI execution (optional)

**Output Interface**:
- Build status (success/failure)
- Test results (pass/fail with details)
- Build artifacts (temporarily available)
- Workflow logs and timing information

**Integration Points**:
- GitHub PR status checks
- Branch protection rule compatibility
- GitHub Actions API integration

### Makefile Interface Changes
**Backward Compatibility**: All existing make targets work unchanged in Flox environments
**New Capability**: Same targets work in non-Flox environments (CI)
**Interface Stability**: No changes to target names, arguments, or expected outputs

## Internal Organization

### Workflow Step Organization

#### Setup Phase (Steps 1-4)
1. **actions/checkout@v4**: Repository checkout
2. **actions/setup-go@v4**: Go 1.26 installation with module caching
3. **actions/setup-node@v4**: Node.js 24 installation with npm caching
4. **Cache restoration**: Restore cached dependencies if available

#### Build Phase (Steps 5-7)
5. **Frontend build**: `make frontend` - npm install and Astro build
6. **Backend build**: `make backend` - Go compilation with embedded assets
7. **Build validation**: `make validate-build` - artifact verification

#### Test Phase (Step 8)
8. **Test execution**: `make test` - Go test suite execution

#### Validation Phase (Steps 9-10)
9. **Smoke test**: Verify binary starts successfully
10. **Artifact check**: Verify no uncommitted changes from build

### Caching Strategy Organization

#### Go Module Cache
**Cache Key**: `go-${{ runner.os }}-${{ hashFiles('go.sum') }}`
**Restore Keys**: `go-${{ runner.os }}-`
**Paths**: Go module cache directory (auto-detected)

#### npm Dependency Cache
**Cache Key**: `npm-${{ runner.os }}-${{ hashFiles('frontend/package-lock.json') }}`
**Restore Keys**: `npm-${{ runner.os }}-`
**Paths**: `frontend/node_modules`

#### Cache Lifecycle
- **Save**: After successful dependency installation
- **Restore**: Before dependency installation
- **Invalidation**: Automatic on lock file changes

## Change Ordering Requirements

### Implementation Sequence
1. **First**: Create `.github/workflows/` directory structure
2. **Second**: Create `ci.yml` workflow file with complete configuration
3. **Third**: Modify Makefile to add Flox detection
4. **Fourth**: Test workflow execution

**Rationale**: Directory and workflow must exist before Makefile changes are tested in CI

### Deployment Sequence
1. **Commit 1**: Add `.github/workflows/ci.yml` with workflow definition
2. **Commit 2**: Modify Makefile for CI compatibility
3. **Optional**: Combine into single commit if testing locally first

### Testing Order
1. **Local testing**: Verify Makefile changes work with and without Flox
2. **CI testing**: Push to test branch and verify workflow execution
3. **Integration testing**: Create test PR to validate trigger behavior

## Module Boundaries

### GitHub Actions Module
**Scope**: CI/CD automation and workflow management
**Responsibilities**:
- Environment setup and tool installation
- Dependency caching and restoration
- Build pipeline orchestration
- Test execution and result reporting
- Status reporting to GitHub

**Dependencies**: GitHub Actions platform, Ubuntu runner environment

### Build System Module
**Scope**: Cross-platform build orchestration
**Responsibilities**:
- Environment detection (Flox vs direct tools)
- Build step coordination
- Artifact validation
- Development workflow support

**Dependencies**: Make, Go toolchain, Node.js toolchain, optional Flox

### Separation of Concerns
- **CI Configuration**: Defined entirely in `.github/workflows/ci.yml`
- **Build Logic**: Remains in Makefile with CI compatibility layer
- **Project Code**: No changes to application source code
- **Environment Management**: Abstracted behind Makefile interface

## Validation Strategy

### Build Validation
1. **Artifact existence**: Binary and frontend assets created
2. **Artifact functionality**: Binary starts without errors
3. **Build reproducibility**: Clean builds produce consistent results
4. **Environment independence**: Builds work with or without Flox

### Integration Validation
1. **Trigger validation**: Workflow runs on correct events
2. **Status reporting**: Results properly reported to GitHub
3. **Caching effectiveness**: Cache hit/miss rates acceptable
4. **Performance validation**: Build times within acceptable limits

### Compatibility Validation
1. **Local development**: Existing Flox-based workflow unchanged
2. **CI environment**: New direct tool workflow functional
3. **Cross-platform**: CI workflow portable to different runners
4. **Version compliance**: Exact tool versions maintained

This structure provides a clear blueprint for implementing GitHub Actions CI while maintaining full compatibility with existing development workflows and build systems.