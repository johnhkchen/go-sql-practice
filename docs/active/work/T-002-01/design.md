# Design: T-002-01 define-collections

## Migration Package Structure

### Option 1: Separate Package with Multiple Files

```
migrations/
├── register.go         # Package initialization
├── 001_create_tags.go  # Tags collection
└── 002_create_links.go # Links collection
```

**Pros**:
- Clear separation of concerns
- Easy to add new migrations
- Standard Go package structure
- Each migration is isolated

**Cons**:
- More files to manage
- Overkill for just two collections

**Verdict**: Over-engineered for current scope

### Option 2: Single Migration File in Package

```
migrations/
└── collections.go  # Both collections in one migration
```

**Pros**:
- Simple package structure
- Related collections in same file
- Easy to see full schema

**Cons**:
- Less granular than separate files
- Harder to rollback individually

**Verdict**: Good balance for current needs

### Option 3: Inline in main.go

```go
// In main.go, before app.Start()
m.Register(func(txApp core.App) error {
    // Create collections here
}, nil, "create_collections")
```

**Pros**:
- No separate package needed
- Everything in one place
- Fastest to implement

**Cons**:
- main.go becomes cluttered
- Not scalable
- Violates separation of concerns

**Verdict**: Too simplistic, poor practice

### Selected: Option 2

**Rationale**: Clean package separation without over-engineering. Single file is appropriate for two related collections. Easy to extend later if needed.

## Collection Creation Order

### Option A: Tags First, Then Links

```go
// 1. Create tags collection
tagsCollection := core.NewBaseCollection("tags")
// configure...
txApp.Save(tagsCollection)

// 2. Create links collection with relation to tags
linksCollection := core.NewBaseCollection("links")
// configure with tags relation...
txApp.Save(linksCollection)
```

**Pros**:
- Logical dependency order
- Tags ID available for relation
- No forward reference issues

**Cons**: None

**Verdict**: **SELECTED** - Natural dependency order

### Option B: Both in Single Transaction

```go
// Create both, save at end
tagsCollection := core.NewBaseCollection("tags")
linksCollection := core.NewBaseCollection("links")
// Configure both...
txApp.Save(tagsCollection)
txApp.Save(linksCollection)
```

**Pros**: Atomic operation
**Cons**: Must handle relation setup carefully
**Verdict**: Unnecessary complexity

## Field Implementation Details

### Tags Collection Fields

**name field**:
```go
&core.TextField{
    Name:     "name",
    Required: true,
    Min:      1,
    Max:      100,
}
```
- Min 1 ensures not empty
- Max 100 reasonable limit

**slug field**:
```go
&core.TextField{
    Name:     "slug",
    Required: true,
    Min:      1,
    Max:      100,
    Pattern:  "^[a-z0-9]+(?:-[a-z0-9]+)*$",
}
```
- Pattern enforces URL-friendly format
- Lowercase with hyphens only

**Unique Index**:
```go
tagsCollection.Indexes = []string{
    "CREATE UNIQUE INDEX idx_tags_slug ON tags (slug)",
}
```

### Links Collection Fields

**url field**:
```go
&core.URLField{
    Name:     "url",
    Required: true,
    OnlyDomains: []string{}, // Allow all domains
}
```
- URLField provides validation
- No domain restrictions initially

**title field**:
```go
&core.TextField{
    Name:     "title",
    Required: true,
    Min:      1,
    Max:      500,
}
```

**description field**:
```go
&core.TextField{
    Name:     "description",
    Required: false,
    Max:      2000,
}
```
- Optional field
- Generous max length

**tags relation**:
```go
&core.RelationField{
    Name:         "tags",
    Required:     false,
    CollectionId: "", // Set after tags saved
    MaxSelect:    100,
}
```
- Many-to-many (MaxSelect > 1)
- Optional (can have untagged links)

**view_count field**:
```go
&core.NumberField{
    Name:     "view_count",
    Required: false,
    Min:      0,
    OnlyInt:  true,
}
```
- Integer only
- Non-negative
- Default handled separately

**created_by relation**:
```go
&core.RelationField{
    Name:         "created_by",
    Required:     false,
    CollectionId: "", // Users collection ID
    MaxSelect:    1,
}
```
- One-to-many (single user)
- Optional for anonymous links

## Registration Pattern

### Option 1: Init Function

```go
func init() {
    m.Register(createCollections, nil, "create_collections")
}
```

**Pros**: Automatic registration
**Cons**: Hidden side effects
**Verdict**: Common but implicit

### Option 2: Explicit Register Function

```go
func Register(app core.App) {
    app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
        return createCollections(e.App)
    })
}
```

**Pros**: Explicit, called from main
**Cons**: More verbose
**Verdict**: **SELECTED** - Clear and testable

### Option 3: Return Migration List

```go
func Migrations() []*m.Migration {
    return []*m.Migration{
        {Up: createCollections, ID: "create_collections"},
    }
}
```

**Pros**: Functional approach
**Cons**: Not standard PocketBase pattern
**Verdict**: Non-idiomatic

## Error Handling Strategy

### Collection Creation Errors

```go
if err := txApp.Save(collection); err != nil {
    return fmt.Errorf("failed to create %s collection: %w", name, err)
}
```
- Wrap errors with context
- Include collection name
- Preserve original error

### Relation Setup

```go
// Get tags collection to find its ID
tagsCollection, err := txApp.FindCollectionByNameOrId("tags")
if err != nil {
    return fmt.Errorf("failed to find tags collection: %w", err)
}
```
- Must lookup after save to get ID
- Clear error messages

## Default Values Strategy

### view_count Default

PocketBase doesn't support field-level defaults directly. Options:

**Option 1**: Database trigger (complex)
**Option 2**: Application-level hooks (later ticket)
**Option 3**: Document as 0, handle in API
**Selected**: Option 3 - Document default, implement in API layer

## Main.go Integration

### Current main.go:
```go
func main() {
    app := pocketbase.New()
    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Updated main.go:
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

Clean, minimal change. Migrations run automatically on bootstrap.

## Testing Strategy

### Manual Testing

1. Delete `pb_data/` directory
2. Run `go run main.go serve`
3. Navigate to `http://localhost:8090/_/`
4. Login to admin (create admin on first run)
5. Verify collections exist with correct fields

### Automated Testing

Could add migration tests, but beyond current scope. Document for future:
```go
func TestCreateCollections(t *testing.T) {
    app := tests.NewTestApp()
    migrations.Register(app)
    // Assert collections exist
}
```

## Selected Design Summary

**Package Structure**:
- New `migrations` package
- Single `collections.go` file
- Explicit `Register()` function

**Implementation Flow**:
1. Create tags collection with fields
2. Save tags to get ID
3. Create links collection with relations
4. Save links

**Field Configurations**:
- Appropriate min/max constraints
- URL validation on url field
- Slug pattern validation
- Relations properly configured

**Integration**:
- Clean registration in main.go
- Runs on app bootstrap
- Wrapped error messages

**Rationale**:
- Simple but extensible
- Clear separation of concerns
- Follows PocketBase patterns
- Easy to test and debug
- Handles dependencies correctly