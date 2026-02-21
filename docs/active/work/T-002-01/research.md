# Research: T-002-01 define-collections

## Current State

### Project Structure

The Go module is initialized with PocketBase v0.36.5 as per T-001-01 (completed):

```
go-sql-practice/
├── main.go              # PocketBase app entry point
├── go.mod               # Module: github.com/jchen/go-sql-practice
├── go.sum               # Dependency checksums
├── pb_data/             # PocketBase data directory (auto-created)
└── docs/
    └── active/
        └── tickets/     # T-002-01 defines collections requirement
```

### Existing Code

**main.go** (lines 1-15):
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

Simple PocketBase initialization without any custom configuration or migrations yet.

### PocketBase Version Analysis

From go.mod:
- Direct dependency: `github.com/pocketbase/pocketbase v0.36.5`
- Status: Listed as indirect (need to fix this)

The PocketBase dependency is incorrectly marked as indirect. This needs correction for proper dependency management.

## Requirements Analysis

### Collection: tags

**Fields Required**:
- `name` (text, required) - Display name for the tag
- `slug` (text, required, unique) - URL-friendly identifier

**Considerations**:
- Unique constraint on slug for URL routing
- System fields (id, created, updated) auto-added by PocketBase
- Base collection type (not auth)

### Collection: links

**Fields Required**:
- `url` (url, required) - The bookmark URL
- `title` (text, required) - Link title
- `description` (text, optional) - Longer description
- `tags` (relation, multiple) - Many-to-many with tags collection
- `view_count` (number, default 0) - Track popularity
- `created_by` (relation to users) - Owner tracking

**Considerations**:
- URL field type provides validation
- Tags relation requires tags collection to exist first
- created_by assumes users collection exists (PocketBase default)
- view_count needs default value of 0

### Migration Requirements

From acceptance criteria:
1. Migration files in `migrations/` package
2. Migrations registered in main.go
3. Run automatically on startup
4. Creates both collections in fresh database
5. Collections visible in admin UI

## PocketBase Migration Patterns

### Migration System Architecture

PocketBase provides two migration approaches:

1. **JavaScript Migrations** - Interpreted at runtime
2. **Go Migrations** - Compiled into binary

For this project, Go migrations are required (per ticket specification).

### Go Migration Structure

Standard PocketBase Go migration pattern:

```go
package migrations

import (
    "github.com/pocketbase/pocketbase/core"
    m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
    m.Register(func(txApp core.App) error {
        // Migration logic here
        return nil
    }, func(txApp core.App) error {
        // Optional rollback logic
        return nil
    }, "migration_name")
}
```

Key points:
- Uses init() for auto-registration
- First func is "up" migration
- Second func is "down" migration (optional)
- Third param is migration identifier

### Collection Creation API

From PocketBase v0.36.x documentation pattern:

```go
collection := core.NewBaseCollection("name")
collection.Fields.Add(&core.TextField{...})
collection.Fields.Add(&core.RelationField{...})
err := txApp.Save(collection)
```

Field types available:
- `core.TextField` - For text fields
- `core.URLField` - For URL validation
- `core.NumberField` - For numeric values
- `core.RelationField` - For relations

## Dependency Graph

### This Ticket Dependencies

**Depends on**:
- T-001-01 (init-go-module) - **COMPLETED**
  - Provides PocketBase app instance
  - main.go structure ready

**Blocks**:
- T-002-02 (seed-database) - Needs collections to exist
- T-003-01 (list-links-api) - Needs links collection
- T-003-02 (crud-link-api) - Needs links collection
- T-003-03 (tags-api) - Needs tags collection

### Collection Dependencies

Internal dependencies:
1. `tags` collection must be created first
2. `links` collection depends on `tags` (for relation)
3. `links` assumes `users` collection exists (PocketBase default)

## File System Observations

### Current Go Files

Only main.go exists. No package structure yet.

### Migration Package Location

Standard Go project structure suggests:
```
migrations/
├── register.go       # Package entry point
└── 001_collections.go # First migration
```

Alternative flat structure:
```
migrations.go  # Single file with all migrations
```

Given the simple scope, either works.

## PocketBase Core API

### Relevant Types (v0.36.5)

From examining the dependency:
- `core.App` - Application interface
- `core.Collection` - Collection model
- `core.NewBaseCollection()` - Collection constructor
- `core.TextField` - Text field type
- `core.URLField` - URL field type
- `core.NumberField` - Number field type
- `core.RelationField` - Relation field type

### Field Configuration

Common field properties:
- `Name` - Field identifier
- `Required` - Validation flag
- `Min/Max` - Length/value constraints
- `Unique` - Uniqueness constraint

Relation-specific:
- `CollectionId` - Target collection
- `MaxSelect` - For multiple relations

## Constraints and Assumptions

### Constraints

1. Must use Go migrations (not JS)
2. Must be in `migrations/` package
3. Must register in main.go
4. Collections must match exact specifications
5. Must work on fresh database

### Assumptions

1. Users collection exists (PocketBase default)
2. Standard PocketBase admin UI available at /_/
3. SQLite database (PocketBase default)
4. No custom validation rules needed initially
5. No custom indexes beyond unique constraints

## Migration Execution Flow

1. App starts with `pocketbase.New()`
2. Migrations registered before `app.Start()`
3. PocketBase checks `_migrations` table
4. Runs pending migrations in order
5. Updates `_migrations` table
6. Collections available immediately

## Testing Considerations

### Verification Steps

1. Delete pb_data directory for fresh start
2. Run application
3. Check console for migration output
4. Access admin UI at localhost:8090/_/
5. Verify both collections exist
6. Check field configurations

### Edge Cases

- Running with existing pb_data
- Collections already exist
- Partial migration failure
- Invalid field configurations

## Next Phase Preparation

The Design phase will need to decide:
- Package structure (separate files vs single)
- Migration naming convention
- Error handling approach
- Registration pattern in main.go
- Field validation rules
- Relation configuration details