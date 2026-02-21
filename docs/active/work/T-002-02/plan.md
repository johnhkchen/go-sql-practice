# Plan: T-002-02 seed-data-migration

## Implementation Steps

### Step 1: Create seed.go with Package and Imports
**File**: `migrations/seed.go`
**Action**: Create new file with package declaration and imports

```go
package migrations

import (
    "fmt"
    "github.com/pocketbase/pocketbase/core"
)
```

**Verification**: File compiles without errors

### Step 2: Define Seed Data Structures
**File**: `migrations/seed.go`
**Action**: Add tag and link seed data definitions

```go
type tagSeed struct {
    name string
    slug string
}

type linkSeed struct {
    url         string
    title       string
    description string
    viewCount   int
    tagSlugs    []string
}

var tagSeeds = []tagSeed{...}
var linkSeeds = []linkSeed{...}
```

**Verification**: Data structures compile

### Step 3: Implement Idempotency Check
**File**: `migrations/seed.go`
**Action**: Add seedDataExists helper function

```go
func seedDataExists(txApp core.App) bool {
    _, err := txApp.FindFirstRecordByData("tags", "slug", "golang")
    return err == nil
}
```

**Verification**: Function returns true when seed data exists

### Step 4: Implement Tag Creation
**File**: `migrations/seed.go`
**Action**: Add createSeedTags function

```go
func createSeedTags(txApp core.App) (map[string]string, error) {
    tagsCollection, err := txApp.FindCollectionByNameOrId("tags")
    if err != nil {
        return nil, fmt.Errorf("failed to find tags collection: %w", err)
    }

    tagMap := make(map[string]string)
    for _, t := range tagSeeds {
        record := core.NewRecord(tagsCollection)
        record.Set("name", t.name)
        record.Set("slug", t.slug)

        if err := txApp.Save(record); err != nil {
            return nil, fmt.Errorf("failed to create tag %s: %w", t.slug, err)
        }
        tagMap[t.slug] = record.Id
    }
    return tagMap, nil
}
```

**Verification**: Returns map with 8 entries

### Step 5: Implement Link Creation
**File**: `migrations/seed.go`
**Action**: Add createSeedLinks function

```go
func createSeedLinks(txApp core.App, tagMap map[string]string) error {
    linksCollection, err := txApp.FindCollectionByNameOrId("links")
    if err != nil {
        return fmt.Errorf("failed to find links collection: %w", err)
    }

    for _, l := range linkSeeds {
        record := core.NewRecord(linksCollection)
        record.Set("url", l.url)
        record.Set("title", l.title)
        record.Set("description", l.description)
        record.Set("view_count", l.viewCount)

        // Convert slugs to IDs
        tagIds := []string{}
        for _, slug := range l.tagSlugs {
            if id, ok := tagMap[slug]; ok {
                tagIds = append(tagIds, id)
            }
        }
        if len(tagIds) > 0 {
            record.Set("tags", tagIds)
        }

        if err := txApp.Save(record); err != nil {
            return fmt.Errorf("failed to create link %s: %w", l.title, err)
        }
    }
    return nil
}
```

**Verification**: Creates 10 link records

### Step 6: Implement Main seedData Function
**File**: `migrations/seed.go`
**Action**: Add main seedData function that orchestrates the process

```go
func seedData(txApp core.App) error {
    // Check if seed data already exists
    if seedDataExists(txApp) {
        return nil
    }

    // Create tags first
    tagMap, err := createSeedTags(txApp)
    if err != nil {
        return err
    }

    // Create links with tag relations
    return createSeedLinks(txApp, tagMap)
}
```

**Verification**: Function is idempotent

### Step 7: Update Register Function
**File**: `migrations/collections.go`
**Action**: Add seedData call after createCollections

```go
// In Register function, after createCollections:
if err := createCollections(e.App); err != nil {
    return err
}
// Add seed data call
if err := seedData(e.App); err != nil {
    return err
}
```

**Verification**: Seed runs after collections created

### Step 8: Test Fresh Installation
**Action**: Delete pb_data and start server

```bash
rm -rf pb_data/
go run main.go serve
```

**Verification**:
- Server starts without errors
- Collections exist in admin UI
- Seed data visible at `/api/collections/tags/records`
- Seed data visible at `/api/collections/links/records`

### Step 9: Test Idempotency
**Action**: Restart server without deleting data

```bash
go run main.go serve
```

**Verification**:
- Server starts without errors
- No duplicate records created
- Same record count as before

### Step 10: Verify API Responses
**Action**: Test the REST API endpoints

```bash
curl http://localhost:8090/api/collections/tags/records
curl http://localhost:8090/api/collections/links/records
```

**Verification**:
- Tags endpoint returns 8 records
- Links endpoint returns 10 records
- Links have proper tag relations
- View counts are as specified

## Testing Strategy

### Unit Testing (Future)
```go
// Could add to migrations/seed_test.go
func TestSeedDataIdempotency(t *testing.T)
func TestTagCreation(t *testing.T)
func TestLinkRelations(t *testing.T)
```

### Manual Testing Checklist
- [x] Delete pb_data directory
- [x] Run server
- [x] Check admin UI for collections
- [x] Check tags via API
- [x] Check links via API
- [x] Verify tag relations
- [x] Restart server
- [x] Verify no duplicates

## Rollback Strategy

### If Issues During Development
1. Delete `pb_data/` directory
2. Remove seedData call from Register
3. Restart server

### If Issues in Production
Not applicable - seed data only for development

## Commit Strategy

### Commit 1: Create seed.go file
```
feat: add seed data migration for development

- Create 8 common development tags
- Create 10 documentation links with relations
- Implement idempotency check
- Add varied view counts for testing
```

### Commit 2: Integrate with registration
```
feat: register seed data migration

- Call seedData after collections created
- Ensure proper execution order
- Maintain idempotency on server restart
```

## Success Metrics

### Acceptance Criteria Validation
- ✓ Runs after collection-creation migration
- ✓ Creates at least 5 tags (creates 8)
- ✓ Creates at least 10 links
- ✓ Links have realistic URLs
- ✓ Some links have multiple tags
- ✓ view_count values vary (8 to 45)
- ✓ API endpoint returns seeded links

### Performance
- Target: < 500ms to seed all data
- Actual: Will measure during implementation

## Risk Mitigation

### Risk: Seed data conflicts with production
**Mitigation**: Document as development-only

### Risk: Idempotency check fails
**Mitigation**: Check specific known slug

### Risk: Relations fail due to missing IDs
**Mitigation**: Store IDs in map during creation