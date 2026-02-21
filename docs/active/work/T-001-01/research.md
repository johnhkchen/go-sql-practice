# Research: T-001-01 init-go-module

## Current State

### Repository Structure

The repository is a fresh project with no Go code yet. Core structure:

```
go-sql-practice/
├── .flox/              # Flox package manager configuration (exists)
├── .lisa/              # Lisa workflow automation (exists)
├── docs/
│   ├── active/
│   │   ├── tickets/    # 18 tickets defined
│   │   ├── stories/    # S-001 through S-005 defined
│   │   └── work/       # T-001-00 artifacts exist
│   ├── archive/
│   └── knowledge/
│       └── rdspi-workflow.md
└── CLAUDE.md           # Project instructions
```

No Go files exist. No `go.mod` or `go.sum` present. This is the absolute first code to be written.

### Environment Configuration

From T-001-00 (completed), the following tools are pinned via Flox:

- **Go 1.26.0**: Available at `.flox/run/x86_64-linux.go-sql-practice.dev/bin/go`
- **Node.js 24.13.0**: For future Astro frontend work
- **GNU Make**: For build orchestration

The `.flox/env/manifest.toml` defines these in the `[install]` section:
- `go_1_26` at version `1.26.0`
- `nodejs_24` at version `24.13.0`
- `gnumake`

Go environment variables:
- `GOPATH`: `/home/jchen/go`
- `GOMODCACHE`: `/home/jchen/go/pkg/mod`

### Ticket Dependencies

**T-001-01** (this ticket) depends on:
- T-001-00 (pin-tool-versions) - **DONE**

**T-001-01** blocks:
- T-001-02 (custom-health-route) - Needs the PocketBase app instance
- T-001-03 (astro-frontend-setup) - Needs Go module structure
- T-001-04 (embed-static-files) - Needs main.go and module
- T-002-01 (define-collections) - Needs PocketBase instance
- All subsequent tickets

### Requirements Analysis

From the ticket acceptance criteria:

1. **Module initialization**:
   - Module path: `github.com/jchen/go-sql-practice`
   - Go version: `1.26`

2. **Dependencies**:
   - PocketBase: `github.com/pocketbase/pocketbase v0.36.5`
   - This must be a direct dependency, not indirect

3. **Entry point**:
   - `main.go` file in root
   - Creates PocketBase app instance
   - Calls `Start()` method

4. **Verification**:
   - `go build` produces a binary
   - Binary starts PocketBase server
   - Admin UI accessible at `/_/`

### PocketBase Integration Patterns

PocketBase v0.36.5 is used as a Go library, not a CLI tool. Standard pattern:

1. **Import**: `github.com/pocketbase/pocketbase`
2. **Instantiation**: `pocketbase.New()` or `pocketbase.NewWithConfig()`
3. **Configuration**: Optional data directory, debug mode, etc.
4. **Lifecycle hooks**: OnServe, OnBootstrap for customization
5. **Start**: `app.Start()` begins the server

The app will embed PocketBase's:
- Admin UI (accessible at `/_/`)
- REST API endpoints
- Auth system
- Database (SQLite by default)

### Related Story Context

From S-001 (project-scaffold), the goal is a single binary that:
- Starts PocketBase with admin UI and REST API
- Registers custom Go routes (T-001-02)
- Serves Astro frontend from embedded filesystem (T-001-03, T-001-04)

This ticket establishes the foundation: the Go module and minimal PocketBase integration.

### Constraints and Assumptions

**Constraints**:
- Must use Go 1.26 as pinned by Flox
- Must use PocketBase v0.36.5 exactly
- Module path must be `github.com/jchen/go-sql-practice`
- Cannot use alternative module paths or versions

**Assumptions**:
- Default PocketBase configuration is acceptable initially
- SQLite database (PocketBase default) is sufficient
- Data directory can be `pb_data` (conventional)
- No custom configuration needed yet

### File System Observations

- Working directory: `/home/jchen/repos/go-sql-practice`
- Git repository initialized (`.git` exists)
- Main branch active
- No existing Go files to preserve or migrate
- Clean slate for implementation

### Module Ecosystem

The `go.mod` will track:
- Direct dependency: `github.com/pocketbase/pocketbase v0.36.5`
- Indirect dependencies: All PocketBase requirements (Echo framework, SQLite driver, etc.)

The `go.sum` will contain cryptographic hashes for all dependencies.

Standard Go module commands:
- `go mod init` - Create module
- `go get` - Add dependencies
- `go mod tidy` - Clean up
- `go build` - Compile binary

### Next Phase Preparation

After research completes, the Design phase will need to decide:
- Exact `main.go` structure
- PocketBase configuration approach
- Directory layout for future packages (routes/, migrations/, etc.)
- Error handling pattern
- Logging setup

The Structure phase will specify:
- File creation order
- Package organization
- Public vs internal boundaries

The Plan phase will sequence:
- Module initialization steps
- Dependency addition
- Code writing
- Testing approach
- Verification criteria