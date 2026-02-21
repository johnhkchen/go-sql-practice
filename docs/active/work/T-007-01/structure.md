# Structure: Presentations Collection (T-007-01)

## File Changes

### Modified Files

#### `/home/jchen/repos/go-sql-practice/migrations/collections.go`

**Location:** After line 196, before closing brace of file

**Addition:** New function `createPresentationsCollection`

```go
// createPresentationsCollection creates the presentations collection
func createPresentationsCollection(txApp core.App) error {
    // Implementation here
}
```

**Location:** Within `createCollections` function, after line 146 (after sync_sessions creation)

**Modification:** Add call to new function

```go
// Create sync_sessions collection
if err := createSyncSessionsCollection(txApp); err != nil {
    return err
}

// Create presentations collection
return createPresentationsCollection(txApp)
```

### No New Files Required

All changes contained within existing migration file structure.

## Component Structure

### createPresentationsCollection Function

**Signature:**
```go
func createPresentationsCollection(txApp core.App) error
```

**Internal Structure:**
1. **Existence Check Block**
   - Check if presentations collection exists
   - Return nil if exists (idempotent)

2. **Dependency Verification Block**
   - Find sync_sessions collection (required)
   - Find users collection (optional)
   - Store collection IDs for relations

3. **Collection Creation Block**
   - Create base collection named "presentations"
   - Set collection type to CollectionTypeBase

4. **Field Definition Block**
   - name field (TextField, required)
   - step_count field (NumberField, required, min 1, integer)
   - step_labels field (JSONField, optional)
   - active_session field (RelationField to sync_sessions)
   - created_by field (RelationField to users, if exists)

5. **API Rules Block**
   - ListRule: public ("")
   - ViewRule: public ("")
   - CreateRule: authenticated ("@request.auth.id != ''")
   - UpdateRule: owner only ("@request.auth.id = created_by")
   - DeleteRule: owner only ("@request.auth.id = created_by")

6. **Save Block**
   - Save collection with error handling
   - Return wrapped error on failure

## Field Specifications

### name Field
```go
&core.TextField{
    Name:     "name",
    Required: true,
    Min:      1,
    Max:      255,
}
```

### step_count Field
```go
stepCountMin := 1.0
&core.NumberField{
    Name:     "step_count",
    Required: true,
    Min:      &stepCountMin,
    OnlyInt:  true,
}
```

### step_labels Field
```go
&core.JSONField{
    Name:     "step_labels",
    Required: false,
}
```

### active_session Field
```go
&core.RelationField{
    Name:         "active_session",
    Required:     false,
    CollectionId: syncSessionsCollection.Id,
    MaxSelect:    1,
}
```

### created_by Field (conditional)
```go
if usersCollection != nil {
    &core.RelationField{
        Name:         "created_by",
        Required:     false,
        CollectionId: usersCollection.Id,
        MaxSelect:    1,
    }
}
```

## Import Requirements

No new imports needed. Existing imports sufficient:
- `github.com/pocketbase/pocketbase/core`
- `github.com/pocketbase/pocketbase/tools/types`
- `fmt` for error wrapping

## Error Handling Structure

### Check Points
1. Collection existence check → return nil
2. sync_sessions lookup → return error if missing
3. users lookup → continue if missing
4. Collection save → return wrapped error

### Error Messages
- "failed to find sync_sessions collection: %w"
- "failed to create presentations collection: %w"

## Integration Points

### With sync_sessions
- Requires sync_sessions collection to exist
- Creates relation field pointing to sync_sessions
- Enables linking active presentation to live session

### With users
- Optional dependency on users collection
- If exists, adds created_by field for ownership
- If missing, continues without ownership tracking

### With Migration Flow
- Called from createCollections after sync_sessions
- Maintains transaction safety via txApp parameter
- Returns error to halt migration on failure

## Ordering Dependencies

1. tags collection (exists)
2. links collection (exists)
3. sync_sessions collection (must exist before presentations)
4. presentations collection (this implementation)

## API Access Pattern

### Public Operations
- List all presentations
- View any presentation

### Authenticated Operations
- Create new presentation (any authenticated user)
- Update own presentation (creator only)
- Delete own presentation (creator only)

### System Operations
- Link/unlink active_session (via custom routes, not direct API)

## Validation Boundaries

### Database Level
- name: required, min 1 char
- step_count: required, min 1, integers only
- step_labels: valid JSON structure
- Relations: foreign key constraints

### Application Level (not in migration)
- step_labels array length matches step_count
- Progress calculation from step_count
- Active session management logic