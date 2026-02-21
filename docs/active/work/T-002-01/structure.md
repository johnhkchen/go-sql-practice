# Structure: T-002-01 define-collections

## File Operations

### Created Files

**migrations/collections.go**
- New package directory: `migrations/`
- New file implementing collection definitions
- Exports single public function: `Register()`

### Modified Files

**main.go**
- Add import: `"github.com/jchen/go-sql-practice/migrations"`
- Insert migration registration before `app.Start()`
- No other changes

**go.mod**
- Fix PocketBase dependency from indirect to direct
- Will be auto-updated by go mod tidy

## Package Architecture

### migrations Package

**Public Interface**:
```go
package migrations

// Register registers all migrations with the PocketBase app
func Register(app core.App)
```

**Internal Structure**:
```go
// createCollections creates tags and links collections
func createCollections(txApp core.App) error
```

Single public function, implementation details hidden.

## Module Boundaries

### Dependencies

```
main
  └── migrations (internal package)
       └── pocketbase/core (external)
```

**Import Structure**:
- `migrations` imports PocketBase core types
- `main` imports `migrations` package
- No circular dependencies

### Data Flow

```
main.go
  ↓ (calls)
migrations.Register()
  ↓ (binds to)
app.OnBootstrap()
  ↓ (executes)
createCollections()
  ↓ (creates)
Database Collections
```

## Component Structure

### migrations/collections.go Components

**1. Import Block**
```go
import (
    "fmt"
    "github.com/pocketbase/pocketbase/core"
    "github.com/pocketbase/pocketbase/forms"
)
```

**2. Register Function** (lines ~10-20)
- Public entry point
- Binds to OnBootstrap event
- Wraps createCollections

**3. createCollections Function** (lines ~25-150)
- Main implementation
- Creates both collections
- Handles relation setup

**4. createTagsCollection Helper** (lines ~155-190)
- Builds tags collection
- Configures name and slug fields
- Returns configured collection

**5. createLinksCollection Helper** (lines ~195-280)
- Builds links collection
- Configures all fields
- Sets up relations
- Returns configured collection

## Field Definitions Detail

### Tags Collection Structure

```go
collection := core.NewBaseCollection("tags")
collection.Type = core.CollectionTypeBase

// Fields array
collection.Fields.Add(&core.TextField{
    Name:     "name",
    Required: true,
    Min:      1,
    Max:      100,
})

collection.Fields.Add(&core.TextField{
    Name:     "slug",
    Required: true,
    Unique:   true,
    Min:      1,
    Max:      100,
    Pattern:  "^[a-z0-9]+(?:-[a-z0-9]+)*$",
})
```

### Links Collection Structure

```go
collection := core.NewBaseCollection("links")
collection.Type = core.CollectionTypeBase

// URL field
collection.Fields.Add(&core.URLField{
    Name:     "url",
    Required: true,
})

// Title field
collection.Fields.Add(&core.TextField{
    Name:     "title",
    Required: true,
    Min:      1,
    Max:      500,
})

// Description field
collection.Fields.Add(&core.TextField{
    Name:     "description",
    Required: false,
    Max:      2000,
})

// View count field
collection.Fields.Add(&core.NumberField{
    Name:     "view_count",
    Required: false,
    Min:      0,
    OnlyInt:  true,
})

// Tags relation (set after tags collection exists)
collection.Fields.Add(&core.RelationField{
    Name:         "tags",
    Required:     false,
    CollectionId: tagsId, // From saved tags collection
    MaxSelect:    100,
})

// Created by relation
collection.Fields.Add(&core.RelationField{
    Name:         "created_by",
    Required:     false,
    CollectionId: usersId, // From users collection lookup
    MaxSelect:    1,
})
```

## Error Handling Structure

### Error Wrapping Pattern

```go
if err := operation(); err != nil {
    return fmt.Errorf("context: %w", err)
}
```

Applied at each level:
1. Collection save operations
2. Collection lookups
3. Event handler return

### Error Messages

- "failed to create tags collection: %w"
- "failed to find tags collection: %w"
- "failed to find users collection: %w"
- "failed to create links collection: %w"

## Integration Points

### main.go Changes

**Before** (current):
```go
func main() {
    app := pocketbase.New()

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

**After** (with migrations):
```go
func main() {
    app := pocketbase.New()

    // Register migrations
    migrations.Register(app)

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

Minimal, single-line addition.

## Transaction Boundaries

### Bootstrap Event

The `OnBootstrap` event runs in a transaction:
- All operations succeed or all fail
- No partial state possible
- Automatic rollback on error

### Collection Save Order

1. Save tags collection (transaction)
2. Get tags collection ID
3. Lookup users collection ID
4. Save links collection (transaction)

Each save is atomic, bootstrap wraps all.

## File Organization

### Directory Layout

```
go-sql-practice/
├── main.go               # Modified: adds migration registration
├── migrations/
│   └── collections.go    # New: all collection definitions
├── go.mod               # Auto-updated by go mod tidy
└── go.sum               # Auto-updated by go mod tidy
```

### Code Organization in collections.go

```go
// Package documentation
package migrations

// Imports

// Public functions
func Register(app core.App)

// Private implementation
func createCollections(txApp core.App) error
func createTagsCollection() *core.Collection
func createLinksCollection(tagsId, usersId string) *core.Collection
```

Clear top-down organization.

## Validation Rules

### Built-in Validations

**Tags**:
- name: required, 1-100 chars
- slug: required, unique, pattern match

**Links**:
- url: required, valid URL format
- title: required, 1-500 chars
- description: optional, max 2000 chars
- view_count: optional, >= 0, integer
- tags: optional, valid relation
- created_by: optional, valid relation

### Database Constraints

Handled by PocketBase:
- Unique index on tags.slug
- Foreign key relations
- Required field constraints
- Type validations

## Configuration Details

### Collection Settings

Both collections:
- Type: `core.CollectionTypeBase` (not auth)
- ListRule: nil (admin only initially)
- ViewRule: nil (admin only initially)
- CreateRule: nil (admin only initially)
- UpdateRule: nil (admin only initially)
- DeleteRule: nil (admin only initially)

Security rules intentionally left for future tickets.

## Testing Boundaries

### Unit Test Points

If tests were added:
- `createTagsCollection()` - returns valid collection
- `createLinksCollection()` - returns valid collection with relations
- Field configurations correct

### Integration Test Points

- Full migration runs without error
- Collections exist after bootstrap
- Fields correctly configured
- Relations properly set

## Next Phase Preparation

The Plan phase will sequence:
1. Create migrations directory
2. Write collections.go file
3. Update main.go
4. Fix go.mod dependency
5. Test migration execution
6. Verify in admin UI