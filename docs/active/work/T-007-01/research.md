# Research: presentations-collection (T-007-01)

## Context

This ticket requires creating a PocketBase collection for presentations. Each presentation is a named sequence of steps that maps to segments of the 0-to-1 progress range. When a presentation goes live, it links to a sync_session from S-006 (implemented in T-006-01).

## Current State

### Project Overview

The project is a Go-based PocketBase application managing collections for links, tags, sync sessions, and now presentations. Key architecture:
- PocketBase v0.36.5 as core framework
- Programmatic migration system via `migrations/collections.go`
- Custom routes via `routes/` package for specialized logic
- SQLite database backend (via PocketBase)
- Astro frontend application (separate concern)

### Existing Collections

Three collections currently exist:

1. **tags** (from T-002-01):
   - `name`: text, required, 1-100 chars
   - `slug`: text, required, unique index, pattern-validated
   - Public read access, no write rules

2. **links** (from T-002-01):
   - `url`: URL field, required
   - `title`: text, required, 1-500 chars
   - `description`: text, optional, max 2000 chars
   - `view_count`: number, integer, min 0
   - `tags`: relation to tags collection, max 100
   - `created_by`: relation to users (if exists), max 1
   - Public read access, no write rules

3. **sync_sessions** (from T-006-01):
   - `progress`: number, float, min 0, max 1
   - `admin_token`: text, required, 64 chars (hex-encoded 32 bytes)
   - Public view/list access
   - No direct create/update/delete rules (nil)
   - Custom routes handle creation and admin updates

### Migration Architecture

From `migrations/collections.go`:

```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        createCollections(e.App)  // Creates tags, links
        createSyncSessionsCollection(e.App)  // Creates sync_sessions
        seedData(e.App)  // Seeds initial data
    })
}
```

Key patterns:
- Idempotent checks: `FindCollectionByNameOrId()` before creation
- Separate functions per collection group
- Error propagation with fmt.Errorf wrapping
- Field definitions using PocketBase core types
- API rules set via pointer strings (nil = no access, "" = public)

### Field Type Usage

From existing implementations:

```go
// TextField
collection.Fields.Add(&core.TextField{
    Name:     "name",
    Required: true,
    Min:      1,
    Max:      100,
    Pattern:  "^[a-z0-9]+(?:-[a-z0-9]+)*$",  // Optional regex
})

// NumberField
minVal := 0.0
maxVal := 1.0
collection.Fields.Add(&core.NumberField{
    Name:     "progress",
    Required: false,
    Min:      &minVal,
    Max:      &maxVal,
    OnlyInt:  false,  // true for integers only
})

// RelationField
collection.Fields.Add(&core.RelationField{
    Name:         "created_by",
    Required:     false,
    CollectionId: usersCollection.Id,
    MaxSelect:    1,  // 1 for single, >1 for multiple
})
```

### Routes Integration

From `routes/sync_sessions.go`:
- Custom endpoints for specialized operations
- Token generation using crypto/rand
- Constant-time token comparison for security
- Record creation/update via app.Save(record)
- Response formatting with structured types

## Requirements Breakdown

### Field Requirements

Per acceptance criteria:
1. **name** (text, required): Human-readable presentation name
2. **step_count** (number, required, min 1): Total number of steps
3. **step_labels** (json, optional): Array of strings labeling each step
4. **active_session** (relation to sync_sessions, optional): Currently live sync session
5. **created_by** (relation to users, optional): Owner

### Progress Calculation

From requirement: "A presentation with `step_count: 5` implies progress breakpoints at `[0.0, 0.25, 0.5, 0.75, 1.0]`"

Analysis:
- 5 steps create 4 intervals: 0→0.25, 0.25→0.5, 0.5→0.75, 0.75→1.0
- Formula: For step i of n steps, progress = i / (n - 1) where i ∈ [0, n)
- Special case: 1 step = single point at 0 or 1? (needs clarification)

### API Rules Translation

"Authenticated users to create/edit their own presentations, anyone to view"

PocketBase rule syntax:
- `ListRule: ""` - Anyone can list
- `ViewRule: ""` - Anyone can view
- `CreateRule: "@request.auth.id != ''"` - Authenticated users only
- `UpdateRule: "created_by = @request.auth.id"` - Only owner
- `DeleteRule: "created_by = @request.auth.id"` - Only owner

## Technical Analysis

### JSON Field Investigation

PocketBase supports JSONField type:
```go
collection.Fields.Add(&core.JSONField{
    Name:     "step_labels",
    Required: false,
})
```

Considerations:
- Stores any valid JSON structure
- No built-in schema validation
- Frontend/routes must validate array structure
- Could store invalid data if not careful

### Relations Setup

**active_session** relation:
- References sync_sessions collection
- Need to get sync_sessions collection ID first
- Optional field (null when not presenting)
- MaxSelect: 1 for single relation

**created_by** relation:
- References users collection
- Must check if users collection exists
- Pattern from links collection: skip if not exists
- MaxSelect: 1 for single relation

### Integration Approach

Following existing pattern from `createSyncSessionsCollection`:
1. Create separate function `createPresentationsCollection()`
2. Call from main `createCollections()` after sync_sessions
3. Check existence first
4. Get required collection IDs for relations
5. Handle missing users collection gracefully

## Risks and Edge Cases

1. **JSON Validation**: No database-level validation for step_labels array structure or length

2. **Users Collection**: May not exist, requiring conditional field addition

3. **Step Count Changes**: If step_count modified, step_labels could be out of sync

4. **Single Step Case**: Unclear how 1 step maps to progress (0? 1? both?)

5. **Active Session Lifecycle**: No requirements for managing active_session state changes

6. **Concurrent Presentations**: Multiple presentations could link to same session

## File Structure

Current:
```
migrations/
  collections.go         # 197 lines, contains createCollections, createSyncSessionsCollection
  seed.go               # Data seeding logic
routes/
  routes.go             # Main route registration
  sync_sessions.go      # Session management endpoints
  health.go             # Health check
main.go                 # Entry point
```

To modify:
- `migrations/collections.go`: Add createPresentationsCollection function

## Next Phase Guidance

Design phase should resolve:
1. JSON validation strategy for step_labels
2. Single step progress mapping behavior
3. Whether custom routes needed for presentation management
4. Active session state management approach
5. Error handling for relation mismatches
6. Whether to add metadata fields (timestamps, status, etc.)