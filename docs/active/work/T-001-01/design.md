# Design: T-001-01 init-go-module

## Implementation Options

### Option 1: Minimal PocketBase Setup

**Approach**: Create the simplest possible `main.go` that starts PocketBase with defaults.

```go
package main

import "github.com/pocketbase/pocketbase"

func main() {
    app := pocketbase.New()
    if err := app.Start(); err != nil {
        panic(err)
    }
}
```

**Pros**:
- Absolute minimum code
- Fast to implement
- Meets all acceptance criteria
- No configuration complexity

**Cons**:
- No error logging setup
- No hooks for future customization
- Would need refactoring for T-001-02

**Verdict**: Too minimal. Would require immediate refactoring.

### Option 2: Structured with Hooks

**Approach**: Include OnServe hook structure for future route registration.

```go
package main

import (
    "log"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        // Future: Register custom routes here
        return e.Next()
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

**Pros**:
- Ready for T-001-02 (custom routes)
- Proper error handling
- Shows the extension pattern

**Cons**:
- Slightly more complex than minimal
- Empty hook might seem unnecessary

**Verdict**: Good balance, but hook is premature.

### Option 3: Clean Entry Point

**Approach**: Simple main with proper error handling, no premature structure.

```go
package main

import (
    "log"
    "github.com/pocketbase/pocketbase"
)

func main() {
    app := pocketbase.New()

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

**Pros**:
- Clean and focused
- Proper error handling with log.Fatal
- Easy to extend later
- Standard Go patterns

**Cons**:
- None significant

**Verdict**: **SELECTED** - Best option for this ticket's scope.

## Module Initialization Strategy

### Option A: Manual Commands

Run commands sequentially:
```bash
go mod init github.com/jchen/go-sql-practice
go get github.com/pocketbase/pocketbase@v0.36.5
```

**Pros**: Fine control
**Cons**: Two-step process
**Verdict**: Standard approach

### Option B: Init Then Edit

```bash
go mod init github.com/jchen/go-sql-practice
# Edit go.mod to add require
go mod download
```

**Pros**: Explicit go.mod control
**Cons**: More complex, manual editing
**Verdict**: Unnecessary complexity

### Option C: Init with Version

```bash
go mod init github.com/jchen/go-sql-practice
go mod edit -go=1.26
go get github.com/pocketbase/pocketbase@v0.36.5
```

**Pros**: Explicit Go version setting
**Cons**: Extra step
**Verdict**: **SELECTED** - Most explicit and correct

## Data Directory Handling

### Option 1: Default Behavior

Let PocketBase use its default `./pb_data` directory.

**Pros**:
- Zero configuration
- Standard PocketBase convention
- Works immediately

**Cons**:
- Creates directory in project root
- Not in .gitignore yet

**Verdict**: **SELECTED** - Follow conventions, add to .gitignore later

### Option 2: Custom Directory

```go
app := pocketbase.NewWithConfig(&pocketbase.Config{
    DefaultDataDir: "./data",
})
```

**Pros**: Custom location
**Cons**: Non-standard, more code
**Verdict**: Unnecessary deviation

## File Organization

### Current Phase Scope

Only two files needed:
- `go.mod` - Module definition
- `main.go` - Entry point

### Future Considerations (NOT this ticket)

Future tickets will add:
- `routes/` - Custom API routes (T-001-02)
- `migrations/` - Database setup (T-002-01)
- `frontend/` - Astro app (T-001-03)
- `embed/` - Static file embedding (T-001-04)

## Error Handling Pattern

### Option 1: panic()

```go
if err := app.Start(); err != nil {
    panic(err)
}
```

**Pros**: Simple
**Cons**: No logging context
**Verdict**: Too abrupt

### Option 2: log.Fatal()

```go
if err := app.Start(); err != nil {
    log.Fatal(err)
}
```

**Pros**:
- Standard library
- Includes timestamp
- Proper exit code

**Cons**: None
**Verdict**: **SELECTED** - Standard Go practice

## Dependency Version Strategy

### PocketBase Version

**Specified**: `v0.36.5` (exact)
**Rationale**: Ticket requirement, ensures reproducibility

### Go Version

**Specified**: `1.26`
**Rationale**: Matches Flox tooling pin from T-001-00

## Testing Strategy

### Build Verification

```bash
go build -o go-sql-practice
./go-sql-practice serve
```

### Admin UI Check

Start server and verify:
- Server starts on port 8090 (default)
- `http://localhost:8090/_/` loads admin UI
- No errors in console

## Chosen Design

**Module initialization**:
1. Create module with explicit Go version
2. Add PocketBase dependency at v0.36.5
3. Let go.sum populate automatically

**main.go structure**:
- Package main
- Import log and pocketbase
- Create app with `pocketbase.New()`
- Start with error handling via `log.Fatal()`
- No premature abstractions

**Rationale**:
- Simplest code that fully satisfies requirements
- Standard Go patterns throughout
- Easy to extend in subsequent tickets
- No unnecessary complexity
- Clear error reporting

This design provides the minimal working foundation that other tickets need while avoiding premature structure that might need refactoring.