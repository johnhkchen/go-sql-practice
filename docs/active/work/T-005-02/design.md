# Design: Build Orchestration for Go-SQL-Practice

## Problem Definition

The current build orchestration in the Makefile is functional but incomplete. Key gaps include:

1. **Development Workflow**: `make dev` doesn't provide hot reload, just starts the server
2. **Incomplete Cleanup**: `make clean` misses `pb_data/` directory
3. **No Build Validation**: Silent failures when frontend assets are missing
4. **Missing Documentation**: No help target or target descriptions
5. **Sequential Dependencies**: No parallel execution optimization
6. **No File Watching**: Manual rebuild required for development changes

## Viable Approaches

### Option 1: Enhanced Makefile

Extend the existing 23-line Makefile with additional targets and improved logic.

**Pros:**
- Builds on existing foundation
- Universal tool availability (make is standard on Unix systems)
- Simple syntax, easy to understand
- No new dependencies
- Already has `.PHONY` declarations and basic structure

**Cons:**
- Limited programming constructs (no arrays, complex conditions)
- Platform-specific (Unix-only, struggles on Windows)
- No built-in file watching capability
- Verbose syntax for complex logic
- Poor error handling primitives

### Option 2: Taskfile (Task runner)

Replace Makefile with a Taskfile.yml using the Task runner tool.

**Pros:**
- YAML syntax is more readable than Makefile
- Built-in file watching with `task --watch`
- Better variable handling and templating
- Cross-platform support (Windows, macOS, Linux)
- Rich dependency management and parallel execution
- Built-in checksum-based change detection

**Cons:**
- Introduces new dependency (requires `task` binary installation)
- Less universal than make (not always pre-installed)
- Learning curve for team members familiar with make
- Adds complexity to development environment setup

### Option 3: Shell Scripts + Makefile

Hybrid approach with shell scripts for complex logic, thin Makefile as interface.

**Pros:**
- Combines make's universality with shell's power
- Complex logic isolated in dedicated scripts
- Maintains familiar `make` interface
- Better error handling in shell scripts
- Can implement custom file watching

**Cons:**
- Multiple files to maintain (Makefile + scripts)
- Shell scripting portability issues (bash vs sh)
- More complex architecture
- Duplicated logic between make and scripts

### Option 4: Pure Go Build Tool

Custom Go program to handle build orchestration.

**Pros:**
- Leverages existing Go expertise in project
- Cross-platform by default
- Rich programming constructs
- Can integrate with Go toolchain directly
- Type safety and error handling

**Cons:**
- Significant development overhead
- Another binary to maintain and distribute
- Reinventing established build tool patterns
- Team must learn custom tool instead of standard tools
- Overkill for current scope

## Design Decision: Enhanced Makefile

**Selected Approach: Option 1 - Enhanced Makefile**

### Rationale

The enhanced Makefile approach is optimal for this project because:

1. **Incremental Improvement**: Builds on the existing working foundation without disruption
2. **Zero New Dependencies**: No additional tools required for development environment
3. **Universal Compatibility**: Works on all Unix-like systems where Go development occurs
4. **Team Familiarity**: Developers already understand and use the existing Makefile
5. **Scope Alignment**: The requirements fit well within Makefile capabilities

### Rejected Options Analysis

**Taskfile (Option 2)** was rejected because:
- The file watching requirement can be met with simple process management
- Adding a new dependency (task binary) introduces setup complexity
- The YAML syntax benefit doesn't justify the dependency cost
- Cross-platform support isn't currently needed (Unix development environment)

**Shell Scripts + Makefile (Option 3)** was rejected because:
- The current logic complexity doesn't justify hybrid architecture
- Maintenance burden increases with multiple files
- Shell portability issues add unnecessary complexity

**Pure Go Build Tool (Option 4)** was rejected because:
- Massive overkill for current requirements
- Development time investment doesn't match value delivered
- Standard tools exist for this exact purpose

## Enhanced Makefile Design

### Core Improvements

#### 1. Development Server with Hot Reload
```makefile
dev: build
	@echo "Starting development server with auto-restart..."
	@while true; do \
		./go-sql-practice serve --http="127.0.0.1:8090" & \
		PID=$$!; \
		inotifywait -e modify -r . --exclude='(pb_data|dist|\.git)' 2>/dev/null || sleep 5; \
		kill $$PID 2>/dev/null || true; \
		wait $$PID 2>/dev/null || true; \
		echo "Rebuilding..."; \
		$(MAKE) build; \
	done
```

**Alternative**: Use simple process management without `inotifywait` dependency:
```makefile
dev: build
	@echo "Starting development server (use Ctrl+C to stop)..."
	./go-sql-practice serve --http="127.0.0.1:8090"

dev-watch:
	@echo "Use 'make dev' in another terminal after changes"
	@echo "Development server with auto-restart requires manual rebuild"
```

#### 2. Comprehensive Cleanup
```makefile
clean:
	@echo "Cleaning build artifacts..."
	rm -rf frontend/dist
	rm -f go-sql-practice
	rm -rf pb_data
	@echo "Clean complete"
```

#### 3. Build Validation
```makefile
backend: frontend
	@echo "Building backend..."
	@if [ ! -d "frontend/dist" ]; then \
		echo "Error: Frontend not built. Run 'make frontend' first."; \
		exit 1; \
	fi
	go build -o go-sql-practice

validate-build: build
	@echo "Validating build..."
	@if [ ! -f "go-sql-practice" ]; then \
		echo "Error: Binary not created"; exit 1; \
	fi
	@if [ ! -d "frontend/dist" ]; then \
		echo "Error: Frontend assets missing"; exit 1; \
	fi
	@echo "Build validation successful"
```

#### 4. Help Target
```makefile
help:
	@echo "Available targets:"
	@echo "  build     - Build frontend and backend"
	@echo "  frontend  - Install deps and build Astro frontend"
	@echo "  backend   - Build Go binary (requires frontend)"
	@echo "  clean     - Remove all build artifacts"
	@echo "  dev       - Start development server"
	@echo "  test      - Run Go tests"
	@echo "  validate  - Validate build artifacts"
	@echo "  help      - Show this help message"
```

#### 5. Independent Target Verification
Each target will include validation:
- `frontend`: Check that `frontend/dist/client/` exists
- `backend`: Verify binary exists and frontend assets embedded
- `test`: Ensure tests run even with no test files (exit gracefully)

### Implementation Strategy

#### Phase 1: Core Enhancements
1. Add build validation to `backend` target
2. Improve `clean` target to include `pb_data`
3. Add `help` target with target documentation
4. Add dependency validation between targets

#### Phase 2: Development Experience
1. Implement simple `dev` target (no file watching initially)
2. Add `validate-build` target for CI/build verification
3. Improve error messages and user feedback

#### Phase 3: Advanced Features (Optional)
1. Consider file watching for `dev-watch` target if team requests
2. Add parallel frontend/backend builds if build time becomes an issue
3. Add smoke tests for build validation

### Risk Mitigation

#### Build Validation Risk
**Risk**: Go embed directive fails silently if frontend assets missing
**Mitigation**: Explicit directory checks before `go build`

#### Development Workflow Risk
**Risk**: Manual rebuild overhead slows development
**Mitigation**: Simple `dev` target first, file watching as Phase 3 enhancement

#### Cross-Platform Risk
**Risk**: Makefile Unix-specific
**Mitigation**: Document Unix requirement, evaluate Taskfile if Windows support needed

### Success Metrics

1. **`make build` reliability**: Zero silent failures
2. **`make clean` completeness**: All artifacts removed
3. **`make dev` usability**: Single command starts development
4. **Independent targets**: Each target can run standalone
5. **Documentation**: `make help` covers all targets

This design addresses all acceptance criteria while maintaining simplicity and building incrementally on the existing foundation.