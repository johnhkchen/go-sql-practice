# Structure: Build Orchestration for Go-SQL-Practice

## File Changes Overview

### Modified Files

#### `/Makefile` (MODIFIED)
**Current**: 23 lines with basic targets (build, frontend, backend, clean, dev, test)
**Changes**: Expand to ~50-60 lines with enhanced targets and validation

**Modifications**:
- Add build validation logic to `backend` target
- Enhance `clean` target to include `pb_data/` directory
- Replace simple `dev` target with proper development server startup
- Add new targets: `help`, `validate-build`
- Add dependency validation between targets
- Improve error messages and user feedback
- Add comments explaining each target's purpose

### New Files

#### None
No new files are created. All enhancements are contained within the existing Makefile.

### Deleted Files

#### None
No files are removed as part of this enhancement.

## Module Boundaries and Organization

### Makefile Target Architecture

#### Core Build Targets
```makefile
# Production build pipeline
build: frontend backend validate-build

# Development workflow
frontend: [dependency validation] → npm install → npm build → [output validation]
backend: [frontend validation] → go build → [binary validation]
```

#### Development Targets
```makefile
# Development workflow
dev: build → start server
help: documentation display
```

#### Maintenance Targets
```makefile
# Cleanup and validation
clean: remove all build artifacts including pb_data/
test: run Go test suite
validate-build: verify all artifacts exist and are valid
```

### Target Dependencies

```
build
├── frontend (parallel eligible)
│   ├── validate node/npm available
│   ├── npm ci (frontend/package-lock.json → frontend/node_modules/)
│   ├── npm run build (frontend/src/ → frontend/dist/)
│   └── validate frontend/dist/client/ exists
└── backend (depends on frontend)
    ├── validate frontend/dist/ exists
    ├── go build (*.go + embedded frontend/dist → go-sql-practice)
    └── validate binary exists

validate-build (depends on build)
├── validate binary exists
├── validate frontend assets exist
└── optional: smoke test binary startup

dev (depends on build)
└── start ./go-sql-practice serve

clean (independent)
├── rm -rf frontend/dist/
├── rm -f go-sql-practice
└── rm -rf pb_data/

test (independent)
└── go test ./...

help (independent)
└── display target documentation
```

### Internal Organization

#### Error Handling Strategy
- **Pre-flight Checks**: Validate prerequisites before expensive operations
- **Post-build Validation**: Verify outputs exist after creation
- **Graceful Failures**: Clear error messages with actionable guidance
- **Exit Codes**: Proper exit codes for CI/automation integration

#### Variable Management
```makefile
# Configuration variables
BINARY_NAME := go-sql-practice
FRONTEND_DIR := frontend
DIST_DIR := $(FRONTEND_DIR)/dist
SERVER_PORT := 127.0.0.1:8090
```

#### Target Implementation Patterns

**Validation Pattern**: Each target validates its inputs and outputs
```makefile
target-name: dependencies
	@echo "Starting target-name..."
	# Pre-flight validation
	@if [ condition ]; then echo "Error: reason"; exit 1; fi
	# Main operation
	command
	# Post-operation validation
	@if [ ! -condition ]; then echo "Error: reason"; exit 1; fi
	@echo "target-name complete"
```

**Dependency Declaration**: Explicit dependencies with validation
```makefile
backend: frontend
	@echo "Building backend..."
	@if [ ! -d "frontend/dist" ]; then \
		echo "Error: Frontend not built. Run 'make frontend' first."; \
		exit 1; \
	fi
	go build -o $(BINARY_NAME)
```

## Component Boundaries

### Build System Responsibilities

#### Makefile Core (Enhanced)
**Responsibilities**:
- Orchestrate build pipeline sequencing
- Validate build prerequisites and outputs
- Provide clear user interface (targets + help)
- Handle error conditions gracefully
- Manage build artifact lifecycle

**Public Interface**:
- `make build` - Complete build pipeline
- `make dev` - Development server startup
- `make test` - Test execution
- `make clean` - Artifact cleanup
- `make help` - Documentation
- `make validate-build` - Build verification

**Internal Boundaries**:
- Frontend build logic (npm operations)
- Backend build logic (Go operations)
- Validation logic (pre/post-build checks)
- Development server management

#### Frontend Build Component
**Responsibilities**:
- Install Node.js dependencies via npm
- Execute Astro build process
- Validate build outputs exist

**Interface**:
- Input: `frontend/src/`, `frontend/package.json`, `frontend/package-lock.json`
- Output: `frontend/dist/client/`, `frontend/dist/server/`
- Commands: `npm ci`, `npm run build`

#### Backend Build Component
**Responsibilities**:
- Compile Go binary with embedded frontend assets
- Validate frontend assets available before build
- Verify binary creation post-build

**Interface**:
- Input: `*.go` files, `frontend/dist/` (via embed)
- Output: `go-sql-practice` binary
- Commands: `go build -o go-sql-practice`

### External Boundaries

#### Development Environment
**Requirements**:
- Unix-like system (Linux, macOS)
- Make utility available
- Node.js v24+ and npm
- Go 1.26+ toolchain

#### Runtime Environment
**Artifacts**:
- `go-sql-practice` binary (self-contained)
- `pb_data/` directory (managed by PocketBase)
- No external dependencies required

#### CI/CD Integration Points
**Entry Points**:
- `make build` - Production builds
- `make test` - Test execution
- `make validate-build` - Build verification
- `make clean` - Artifact cleanup

**Exit Codes**:
- 0: Success
- 1: Build failure
- 2: Validation failure (defined convention)

## Implementation Ordering

### Phase 1: Core Enhancements (Critical Path)
1. **Backend target validation**: Add frontend dependency check
2. **Clean target expansion**: Include `pb_data/` removal
3. **Help target addition**: User documentation and guidance
4. **Variable definition**: Centralize configuration values

### Phase 2: Validation Infrastructure
1. **Build validation target**: Post-build verification
2. **Error message improvement**: Clear, actionable feedback
3. **Dependency validation**: Pre-flight checks for each target

### Phase 3: Development Experience
1. **Dev target enhancement**: Proper server startup with error handling
2. **Independent target verification**: Ensure each target runs standalone
3. **Documentation completion**: Inline comments and help text

### Implementation Constraints

#### Backward Compatibility
- All existing `make` commands continue to work
- Existing behavior preserved, only enhanced
- No breaking changes to established workflow

#### Minimal Disruption
- Single file modification (Makefile only)
- No new dependencies introduced
- No changes to project structure or build outputs

#### Incremental Enhancement
- Each phase can be implemented and tested independently
- Rollback capability maintained throughout implementation
- Gradual feature introduction without "big bang" changes

This structure provides the blueprint for implementing the enhanced build orchestration system while maintaining simplicity and building incrementally on the existing foundation.