# Design: Presentations Collection (T-007-01)

## Options Considered

### Option 1: Standard PocketBase JSON Field
Use PocketBase's built-in JSON field type for step_labels.

**Pros:**
- Native PocketBase support (v0.36.5 includes JSONField)
- Automatic JSON validation
- Works with PocketBase admin UI
- Consistent with PocketBase patterns

**Cons:**
- Limited validation of internal structure
- Array length validation must be application-level

**Implementation:**
```go
collection.Fields.Add(&core.JSONField{
    Name:     "step_labels",
    Required: false,
})
```

### Option 2: Text Field with JSON String
Store step_labels as validated text containing JSON.

**Pros:**
- More control over validation
- Can add pattern matching

**Cons:**
- Manual JSON parsing required
- Not idiomatic PocketBase
- Breaks admin UI compatibility
- Additional serialization overhead

**Rejected:** Unnecessarily complex for simple array storage.

### Option 3: Separate Step Labels Collection
Create a related collection for step labels.

**Pros:**
- Strong typing per label
- Individual label management
- Relational integrity

**Cons:**
- Over-engineering for simple string array
- Multiple queries for single presentation
- Complex migration

**Rejected:** Excessive complexity for requirements.

## API Rules Design

### Option A: PocketBase Rule Expressions
Use PocketBase's built-in rule system for ownership.

**Implementation:**
```go
collection.CreateRule = types.Pointer("@request.auth.id != ''")
collection.UpdateRule = types.Pointer("@request.auth.id = created_by")
collection.DeleteRule = types.Pointer("@request.auth.id = created_by")
```

**Pros:**
- Native PocketBase authorization
- Declarative rules
- Automatic enforcement

**Cons:**
- Requires authenticated users
- Complex for non-owner updates

### Option B: Custom Route Authorization
Handle all write operations via custom routes.

**Pros:**
- Full control over authorization logic
- Can handle special cases
- Consistent with sync_sessions pattern

**Cons:**
- More code to maintain
- Bypasses PocketBase authorization

**Rejected:** Unnecessary when built-in rules suffice.

## Migration Integration

### Option I: Extend createCollections Function
Add presentations creation to existing function.

**Pros:**
- Single migration function
- Runs with other collections

**Cons:**
- Growing function complexity
- Harder to test individually

### Option II: Separate createPresentationsCollection Function
Create dedicated function, call from main migration.

**Pros:**
- Modular, testable
- Follows sync_sessions pattern
- Clear separation of concerns
- Easy to disable/rollback

**Cons:**
- Additional function call

**Selected:** Option II matches existing patterns.

## Field Validation Strategy

### step_count Validation
- Minimum: 1 (at least one step required)
- Maximum: None (let applications decide)
- Type: Integer only

### step_labels Validation
- Type: JSON array of strings
- Optional field (can be null/empty)
- Length validation: Application-level check against step_count
- No database-level array length enforcement

### active_session Relation
- Optional single relation to sync_sessions
- Null when not presenting
- Foreign key constraint automatic via PocketBase

### created_by Relation
- Optional single relation to users
- Handle missing users collection gracefully
- Used for ownership rules

## Selected Design

**Collection Structure:** Option 1 (Native JSON field)
- Use core.JSONField for step_labels
- Leverages PocketBase native capabilities
- Simple, maintainable

**API Rules:** Option A (PocketBase expressions)
- Public read access (list/view)
- Authenticated create
- Owner-only update/delete
- Falls back gracefully if no auth

**Migration:** Option II (Separate function)
- createPresentationsCollection() function
- Called from createCollections()
- Checks existence before creation
- Handles missing users collection

**Implementation Pattern:**
1. Check if presentations exists (skip if yes)
2. Check if sync_sessions exists (dependency)
3. Get sync_sessions collection ID
4. Check if users collection exists
5. Create presentations collection
6. Add all fields with proper types
7. Set API rules for public read, authenticated write
8. Save collection

## Error Handling

- Collection already exists: Skip silently (idempotent)
- sync_sessions missing: Return error (hard dependency)
- users collection missing: Continue without created_by field
- Field creation failures: Return wrapped errors
- Save failures: Return detailed error

## Testing Strategy

- Start app, verify collection created
- Check field types and constraints
- Test with/without users collection
- Verify API rules work correctly
- Confirm step_labels accepts JSON arrays
- Test active_session relation to sync_sessions

## Rationale

This design:
1. **Minimizes complexity** by using native PocketBase features
2. **Follows existing patterns** from sync_sessions implementation
3. **Maintains consistency** with current migration structure
4. **Handles edge cases** like missing users collection
5. **Enables future extensions** without breaking changes

The native JSON field eliminates custom serialization while PocketBase rules provide robust authorization without custom code. The modular migration function keeps code organized and testable.