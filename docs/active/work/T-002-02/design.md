# Design: T-002-02 seed-data-migration

## Idempotency Strategy

### Option 1: Check Each Record Individually
```go
// Check if "golang" tag exists
if _, err := app.FindFirstRecordByData("tags", "slug", "golang"); err != nil {
    // Create it
}
```

**Pros**:
- Granular control
- Can update existing records
- Handles partial seed state

**Cons**:
- Many database queries
- Complex for relations
- Verbose code

**Verdict**: Too complex for seed data

### Option 2: All-or-Nothing Check
```go
// Check if any seed tag exists
if _, err := app.FindFirstRecordByData("tags", "slug", "golang"); err == nil {
    return nil // Assume all seed data exists
}
```

**Pros**:
- Simple logic
- Fast (single check)
- Clear state

**Cons**:
- Can't handle partial state
- Manual deletion breaks it

**Verdict**: **SELECTED** - Appropriate for dev seed data

### Option 3: Delete and Recreate
```go
// Delete all seed records first
// Then create fresh
```

**Pros**:
- Always fresh data
- Simple logic

**Cons**:
- Destroys user modifications
- Breaks foreign keys
- Not truly idempotent

**Verdict**: Destructive, inappropriate

## File Organization

### Option A: Add to collections.go
```go
// In existing migrations/collections.go
func seedData(txApp core.App) error {
    // Seed logic
}
```

**Pros**:
- Single file for all initial setup
- Runs in same transaction

**Cons**:
- File becomes large
- Mixes concerns

**Verdict**: Acceptable but not ideal

### Option B: Separate seed.go File
```go
// New migrations/seed.go
package migrations

func seedData(txApp core.App) error {
    // Seed logic
}
```

**Pros**:
- Clean separation
- Easy to find/modify
- Follows single responsibility

**Cons**:
- Another file to maintain

**Verdict**: **SELECTED** - Clean architecture

### Option C: Separate Package
```go
// seeds/initial.go
package seeds
```

**Pros**:
- Maximum separation

**Cons**:
- Over-engineered
- Package for one file

**Verdict**: Unnecessary complexity

## Registration Pattern

### Option 1: Chain in Register Function
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        if err := createCollections(e.App); err != nil {
            return err
        }
        return seedData(e.App)
    })
}
```

**Pros**:
- Sequential guarantee
- Single registration

**Cons**:
- Couples migrations

**Verdict**: **SELECTED** - Simple and correct

### Option 2: Separate Event Handler
```go
app.OnServe().BindFunc(createCollections)
app.OnServe().BindFunc(seedData)
```

**Pros**:
- Independent handlers

**Cons**:
- Order not guaranteed
- Could run in parallel

**Verdict**: Race condition risk

## Seed Data Content

### Tags to Create
```go
tags := []struct {
    name string
    slug string
}{
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

**Rationale**:
- Common development categories
- Mix of languages and concepts
- Enough for realistic queries
- Natural groupings for links

### Links to Create
```go
links := []struct {
    url         string
    title       string
    description string
    viewCount   int
    tagSlugs    []string
}{
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
    // ... 8 more links
}
```

**Rationale**:
- Real documentation sites
- Varied view counts (0-50)
- Multiple tags per link
- Realistic titles/descriptions

## Implementation Details

### Tag Creation
```go
tagsCollection, _ := txApp.FindCollectionByNameOrId("tags")
tagMap := make(map[string]string) // slug -> ID

for _, t := range tags {
    record := core.NewRecord(tagsCollection)
    record.Set("name", t.name)
    record.Set("slug", t.slug)

    if err := txApp.Save(record); err != nil {
        return fmt.Errorf("failed to create tag %s: %w", t.slug, err)
    }
    tagMap[t.slug] = record.Id
}
```

### Link Creation with Relations
```go
linksCollection, _ := txApp.FindCollectionByNameOrId("links")

for _, l := range links {
    record := core.NewRecord(linksCollection)
    record.Set("url", l.url)
    record.Set("title", l.title)
    record.Set("description", l.description)
    record.Set("view_count", l.viewCount)

    // Convert tag slugs to IDs
    tagIds := []string{}
    for _, slug := range l.tagSlugs {
        if id, ok := tagMap[slug]; ok {
            tagIds = append(tagIds, id)
        }
    }
    record.Set("tags", tagIds)

    if err := txApp.Save(record); err != nil {
        return fmt.Errorf("failed to create link %s: %w", l.title, err)
    }
}
```

## Error Handling

### Consistent Pattern
```go
return fmt.Errorf("failed to [action] [item]: %w", err)
```

### Check Existence First
```go
collection, err := txApp.FindCollectionByNameOrId("tags")
if err != nil {
    return fmt.Errorf("tags collection not found: %w", err)
}
```

## Complete Seed Data Set

### All Links (10 total)
1. Go Documentation - golang, backend - 42 views
2. React - javascript, frontend - 38 views
3. PostgreSQL Docs - database, backend - 25 views
4. Docker Hub - devops - 31 views
5. Testing Go Code - golang, testing, backend - 18 views
6. MDN Web Docs - javascript, frontend - 45 views
7. Kubernetes Documentation - devops, architecture - 22 views
8. Astro - frontend, javascript - 15 views
9. PocketBase - database, backend, golang - 8 views
10. GitHub Actions - devops, testing - 12 views

**Distribution**:
- View counts: 8 to 45 (varied)
- Tags per link: 1 to 3
- All tags used at least once
- Realistic documentation URLs

## Selected Design Summary

**File Structure**:
- New `migrations/seed.go` file
- Called from existing Register function

**Idempotency**:
- Check if "golang" tag exists
- Skip all seeding if found
- Assumes complete seed set

**Implementation**:
1. Check for existing seed data
2. Create all tags first
3. Store tag IDs in map
4. Create links with relations
5. Use consistent error messages

**Data**:
- 8 descriptive tags
- 10 realistic documentation links
- Proper many-to-many relations
- Varied view counts

**Rationale**:
- Clean separation of concerns
- Fast idempotency check
- Realistic test data
- Follows existing patterns