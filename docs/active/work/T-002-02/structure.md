# Structure: T-002-02 seed-data-migration

## File Changes

### New Files

#### migrations/seed.go
```go
package migrations

import (
    "fmt"
    "github.com/pocketbase/pocketbase/core"
)

// seedData creates initial seed data for development
func seedData(txApp core.App) error
```

**Purpose**: Contains all seed data logic
**Exports**: `seedData` function (package-internal)
**Dependencies**: pocketbase/core for record operations

### Modified Files

#### migrations/collections.go
```go
// In Register function, after createCollections:
if err := createCollections(e.App); err != nil {
    return err
}
// ADD: Seed data after collections exist
if err := seedData(e.App); err != nil {
    return err
}
```

**Change**: Chain seed data creation after collections
**Impact**: Ensures collections exist before seeding

## Module Structure

### Package: migrations

**Functions**:
- `Register(app core.App)` - Entry point (existing)
- `createCollections(txApp core.App) error` - Creates schema (existing)
- `seedData(txApp core.App) error` - Creates seed data (new)

**Internal Structure of seedData**:
1. Check for existing seed data
2. Create tags
3. Create links with relations

## Data Structures

### Tag Seed Data
```go
type tagSeed struct {
    name string
    slug string
}

var tagSeeds = []tagSeed{
    {"Go", "golang"},
    {"JavaScript", "javascript"},
    {"Database", "database"},
    {"DevOps", "devops"},
    {"Frontend", "frontend"},
    {"Backend", "backend"},
    {"Testing", "testing"},
    {"Architecture", "architecture"},
}
```

### Link Seed Data
```go
type linkSeed struct {
    url         string
    title       string
    description string
    viewCount   int
    tagSlugs    []string
}

var linkSeeds = []linkSeed{
    {
        url:         "https://go.dev/doc/",
        title:       "Go Documentation",
        description: "Official Go programming language documentation",
        viewCount:   42,
        tagSlugs:    []string{"golang", "backend"},
    },
    // ... 9 more entries
}
```

## Function Boundaries

### seedData Function Structure
```go
func seedData(txApp core.App) error {
    // 1. Idempotency check
    if seedDataExists(txApp) {
        return nil
    }

    // 2. Create tags
    tagMap, err := createSeedTags(txApp)
    if err != nil {
        return err
    }

    // 3. Create links
    return createSeedLinks(txApp, tagMap)
}
```

### Helper Functions
```go
// Check if seed data already exists
func seedDataExists(txApp core.App) bool

// Create all seed tags, return slug->ID map
func createSeedTags(txApp core.App) (map[string]string, error)

// Create all seed links with relations
func createSeedLinks(txApp core.App, tagMap map[string]string) error
```

## Error Boundaries

### Error Message Format
```
"failed to [action] [subject]: %w"
```

Examples:
- "failed to find tags collection: %w"
- "failed to create tag golang: %w"
- "failed to create link Go Documentation: %w"

### Error Propagation
- All errors bubble up to Register function
- Server startup fails on seed error
- Original error preserved with %w

## Interface Boundaries

### PocketBase API Usage
```go
// Collection retrieval
txApp.FindCollectionByNameOrId(string) (*core.Collection, error)

// Record creation
core.NewRecord(*core.Collection) *core.Record

// Field setting
record.Set(string, interface{})

// Record persistence
txApp.Save(*core.Record) error

// Record query (for idempotency)
txApp.FindFirstRecordByData(collection, field, value) (*core.Record, error)
```

### No External APIs
- All data hardcoded in source
- No file I/O required
- No network calls

## Data Flow

### Execution Flow
```
app.Start()
  └── OnServe event
      └── Register() binding
          └── createCollections()
          └── seedData()
              ├── Check existing
              ├── Create tags
              └── Create links
```

### Data Dependencies
```
Tags (must exist first)
  └── Links (reference tag IDs)
      └── Tag relations (array of IDs)
```

### Tag ID Mapping
```
Tag slug → Tag record → Tag ID → Link relation
```

## Complete Link Seeds Structure

```go
var linkSeeds = []linkSeed{
    {
        url:         "https://go.dev/doc/",
        title:       "Go Documentation",
        description: "Official Go programming language documentation",
        viewCount:   42,
        tagSlugs:    []string{"golang", "backend"},
    },
    {
        url:         "https://react.dev/",
        title:       "React",
        description: "The library for web and native user interfaces",
        viewCount:   38,
        tagSlugs:    []string{"javascript", "frontend"},
    },
    {
        url:         "https://www.postgresql.org/docs/",
        title:       "PostgreSQL Documentation",
        description: "The world's most advanced open source database",
        viewCount:   25,
        tagSlugs:    []string{"database", "backend"},
    },
    {
        url:         "https://hub.docker.com/",
        title:       "Docker Hub",
        description: "Container image library and community",
        viewCount:   31,
        tagSlugs:    []string{"devops"},
    },
    {
        url:         "https://go.dev/doc/tutorial/add-a-test",
        title:       "Testing in Go",
        description: "Learn how to write unit tests in Go",
        viewCount:   18,
        tagSlugs:    []string{"golang", "testing", "backend"},
    },
    {
        url:         "https://developer.mozilla.org/",
        title:       "MDN Web Docs",
        description: "Resources for developers, by developers",
        viewCount:   45,
        tagSlugs:    []string{"javascript", "frontend"},
    },
    {
        url:         "https://kubernetes.io/docs/",
        title:       "Kubernetes Documentation",
        description: "Production-grade container orchestration",
        viewCount:   22,
        tagSlugs:    []string{"devops", "architecture"},
    },
    {
        url:         "https://astro.build/",
        title:       "Astro",
        description: "The web framework for content-driven websites",
        viewCount:   15,
        tagSlugs:    []string{"frontend", "javascript"},
    },
    {
        url:         "https://pocketbase.io/docs/",
        title:       "PocketBase",
        description: "Open source backend in 1 file",
        viewCount:   8,
        tagSlugs:    []string{"database", "backend", "golang"},
    },
    {
        url:         "https://docs.github.com/actions",
        title:       "GitHub Actions Documentation",
        description: "Automate your workflow from idea to production",
        viewCount:   12,
        tagSlugs:    []string{"devops", "testing"},
    },
}
```

## Constraints

### Must Maintain
- Idempotency (safe to run multiple times)
- Order (tags before links)
- Error handling consistency

### Must Not
- Delete existing non-seed data
- Modify existing seed data
- Create duplicates

### Performance
- Runs on every server start
- Should complete in < 1 second
- ~18 database operations total