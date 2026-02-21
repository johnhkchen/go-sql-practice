# Design: View Count Endpoint (T-003-02)

## Problem Statement

Implement `POST /api/links/:id/view` to atomically increment a link's view count. The key challenge is ensuring atomicity under concurrent requests while maintaining consistency with existing PocketBase patterns and the project architecture.

## Design Options Analysis

### Option 1: Pure SQL Atomic Update

**Approach**: Use raw SQL `UPDATE` statement to atomically increment view_count, then fetch the updated record using PocketBase record operations.

**Implementation**:
```go
// Atomic increment
db.NewQuery("UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?").
    Exec(linkId)

// Fetch updated record
record := dao.FindRecordById("links", linkId)
```

**Pros**:
- Guaranteed atomicity via single SQL operation
- Handles NULL view_count case with COALESCE
- Leverages SQLite's row-level locking
- Follows pattern from stats.go for direct SQL usage
- No race conditions even under high concurrency

**Cons**:
- Requires two database operations (update + fetch)
- Bypasses PocketBase's built-in update hooks/validation
- Need to handle case where UPDATE affects 0 rows (link not found)

**Risk Assessment**: Low. SQLite atomic updates are well-established and reliable.

### Option 2: PocketBase Record Update with Transaction

**Approach**: Use PocketBase record operations within an explicit database transaction.

**Implementation**:
```go
dao.RunInTransaction(func(txDao *daos.Dao) error {
    record := txDao.FindRecordById("links", linkId)
    viewCount := record.GetInt("view_count") + 1
    record.Set("view_count", viewCount)
    return txDao.SaveRecord(record)
})
```

**Pros**:
- Uses PocketBase's native record operations
- Triggers any defined hooks/validation
- Single transaction ensures atomicity
- More "idiomatic" PocketBase code

**Cons**:
- Read-modify-write pattern - potential race condition window
- More complex error handling
- Transaction overhead
- Research shows no existing patterns using explicit transactions in this codebase

**Risk Assessment**: Medium. Race conditions possible despite transaction wrapping.

### Option 3: PocketBase Update Hook with SQL

**Approach**: Use PocketBase's `OnRecordUpdateRequest` hook to intercept update attempts and perform atomic SQL.

**Implementation**:
```go
app.OnRecordUpdateRequest("links").BindFunc(func(e *core.RecordUpdateRequestEvent) error {
    if isViewIncrementRequest(e) {
        return performAtomicIncrement(e)
    }
    return e.Next()
})
```

**Pros**:
- Integrates with PocketBase event system
- Could potentially be reused for other atomic operations
- Maintains PocketBase's request/response patterns

**Cons**:
- Overcomplicates a simple increment operation
- Hook registration complexity
- Difficult to isolate to specific endpoint
- No existing hook patterns in codebase
- Deviates from the direct route registration pattern

**Risk Assessment**: High. Adds unnecessary complexity without clear benefits.

### Option 4: Custom Stored Procedure/Function

**Approach**: Create SQLite function or use RETURNING clause to increment and return in single operation.

**Implementation**:
```sql
UPDATE links SET view_count = COALESCE(view_count, 0) + 1
WHERE id = ?
RETURNING *
```

**Pros**:
- Single database round-trip
- Atomic increment and fetch
- Efficient performance

**Cons**:
- RETURNING clause support varies by SQLite version
- Would need to manually map returned row to PocketBase record format
- Deviates from established patterns in codebase
- More complex result processing

**Risk Assessment**: Medium. SQLite version compatibility concerns.

## Design Decision

**Selected Option: Option 1 - Pure SQL Atomic Update**

**Rationale**:

1. **Atomicity Guarantee**: The single `UPDATE` statement provides true atomicity without race conditions. Even under high concurrency, SQLite's row-level locking ensures each increment is isolated.

2. **Codebase Consistency**: The existing `stats.go` demonstrates direct SQL usage for performance-critical operations. This approach aligns with established patterns.

3. **Simplicity**: Clear, straightforward implementation that's easy to understand, test, and maintain.

4. **NULL Handling**: `COALESCE(view_count, 0) + 1` elegantly handles the case where view_count might be NULL, defaulting to 0 before increment.

5. **Error Handling**: Can easily detect if no rows were affected by the UPDATE to return 404.

**Rejected Options**:
- **Option 2**: Read-modify-write creates race condition window despite transaction
- **Option 3**: Unnecessarily complex for a simple atomic operation
- **Option 4**: SQLite version compatibility risks and complex result mapping

## Implementation Architecture

### Database Operation Flow
1. Extract `id` parameter from URL path
2. Execute atomic SQL: `UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?`
3. Check affected rows count (0 = link not found)
4. If updated, fetch complete record using PocketBase DAO
5. Return JSON response with full record or 404 error

### Error Handling Strategy
- **Invalid ID format**: Return 400 Bad Request (if needed)
- **Link not found**: UPDATE affects 0 rows → Return 404 with JSON error
- **Database errors**: Return 500 with generic error message
- **JSON serialization errors**: Return 500 (unlikely with PocketBase record)

### Response Format Design
Follow PocketBase's standard record JSON format to maintain API consistency:
- Include all fields: id, url, title, description, view_count, tags, created, updated
- Use PocketBase's built-in JSON serialization
- Success: 200 with full record
- Not found: 404 with `{"error": "Link not found"}`

## File Organization

### New File: `routes/links_view.go`
Contains the complete implementation:
- `registerLinksView(e *core.ServeEvent)` function
- Route handler with parameter extraction
- Atomic SQL execution
- Record fetching and response formatting
- Error handling for all cases

### Modified File: `routes/routes.go`
Add single line: `registerLinksView(e)` in the route registration sequence.

## Concurrency Safety Analysis

The chosen design is concurrency-safe because:

1. **Atomic SQL Operation**: `UPDATE ... SET view_count = COALESCE(view_count, 0) + 1` is a single SQL statement that SQLite executes atomically.

2. **Row-Level Locking**: SQLite's default behavior provides row-level locking during the UPDATE operation.

3. **No Race Conditions**: There's no read-modify-write cycle in application code that could create race conditions.

4. **COALESCE Safety**: Handles NULL values atomically within the SQL operation itself.

Even with 100 concurrent requests to the same link ID, each increment will be processed sequentially by SQLite, ensuring no increments are lost.

## API Contract Definition

```http
POST /api/links/:id/view

Success Response (200):
{
  "id": "record_id_here",
  "url": "https://example.com",
  "title": "Example Link Title",
  "description": "Optional description text",
  "view_count": 42,
  "tags": ["tag1", "tag2"],
  "created": "2024-01-01T00:00:00.000Z",
  "updated": "2024-01-02T00:00:00.000Z",
  "created_by": "user_id_or_null"
}

Error Response (404):
{
  "error": "Link not found"
}
```

## Testing Strategy

### Unit Testing Approach
1. Test successful increment with existing link
2. Test 404 response with non-existent link ID
3. Test NULL view_count handling (first increment sets to 1)
4. Test response format matches PocketBase record structure

### Integration Testing
1. Verify endpoint registration and routing
2. Test with real database and HTTP requests
3. Verify JSON response format

### Concurrency Testing
1. Multiple simultaneous requests to same link ID
2. Verify all increments are counted (no lost updates)
3. Performance testing under load

## Migration Considerations

No database migration required:
- `links` collection already exists with `view_count` field
- Seed data already populates view_count values
- Field definition allows NULL values (handled by COALESCE)

## Security Analysis

### Input Validation
- Link ID parameter validation handled by PocketBase/Echo router
- SQL injection prevention via parameterized queries
- No authentication required (per acceptance criteria)

### Potential Abuse Vectors
- Unlimited increment requests (no rate limiting specified)
- Could artificially inflate view counts
- Public endpoint accessible to all users

### Mitigation Notes
Rate limiting and abuse prevention are out of scope for this ticket but should be considered for production deployment.

## Performance Characteristics

- **Database Operations**: 1 UPDATE + 1 SELECT per request
- **Memory Usage**: Minimal - single record in memory
- **Latency**: Low - simple SQL operations on indexed ID field
- **Scalability**: Good - atomic operations scale well with SQLite

The design prioritizes correctness (atomicity) while maintaining reasonable performance for the expected usage patterns of a practice application.