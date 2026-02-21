# Research: View Count Endpoint (T-003-02)

## Overview

The task requires implementing `POST /api/links/:id/view` to atomically increment a link's view count using SQL atomic operations. This research maps the existing codebase to understand patterns, constraints, and implementation approach.

## Codebase Architecture

### Project Structure
- **Main application**: `main.go` - PocketBase app with custom migrations and routes
- **Migrations**: `migrations/` - Collection definitions and seed data
- **Custom routes**: `routes/` - Custom API endpoints beyond PocketBase defaults
- **Dependencies**: PocketBase framework with Echo router and SQLite database

### Database Schema

#### Links Collection (`collections.go:87-155`)
The `links` collection is already defined with:
- `url` (URLField, required)
- `title` (TextField, required, 1-500 chars)
- `description` (TextField, optional, max 2000 chars)
- **`view_count` (NumberField, optional, min=0, integer only)** - Key field for this task
- `tags` (RelationField to tags collection)
- `created_by` (RelationField to users, optional)

#### Seed Data (`seed.go:36-108`)
Sample links have `view_count` values ranging from 8-45, confirming the field exists and is populated.

### Route Registration Pattern

#### Registration Flow (`routes.go:8-20`)
```go
func Register(app core.App) {
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        registerHealth(e)
        registerSyncSessions(e)
        registerStats(e)
        registerLinksSearch(e)
        registerLinksSearchSimple(e)
        // New registerLinksView(e) would go here
        return e.Next()
    })
}
```

#### Route Handler Pattern
Based on existing routes, the pattern is:
1. Register route with HTTP method and path
2. Handler function receives `*core.RequestEvent` and `core.App`
3. Extract path parameters and validate
4. Execute database operations using `app.DB()`
5. Return JSON response with appropriate HTTP status

### Existing Database Patterns

#### Direct SQL Operations (`stats.go`)
The stats endpoint demonstrates raw SQL usage:
- Uses `app.DB()` to get database connection
- Executes queries with `db.NewQuery(sql).Row(&result)`
- Handles both single-row results and row iteration
- Error handling returns 500 status with JSON error

#### PocketBase Record Operations
PocketBase provides high-level record operations, but for atomic increment we need raw SQL to avoid read-modify-write race conditions.

### URL Parameter Extraction

From existing code patterns (not directly shown but inferred from PocketBase/Echo router), path parameters like `:id` are accessed via:
- `e.Request.PathParam("id")` or similar Echo router method

## Dependencies Analysis

### T-001-02: Custom Health Route (DONE)
✅ Establishes the route registration pattern in `routes/` package that we'll follow.

### T-002-02: Seed Data Migration (DONE)
✅ Confirms `links` collection with `view_count` field exists and is populated.

## Database Constraints

### SQLite Atomic Operations
- PocketBase uses SQLite as the database
- SQLite supports atomic `UPDATE` operations
- `UPDATE links SET view_count = view_count + 1 WHERE id = ?` is atomic
- No explicit transaction needed for single atomic operation

### Data Types
- `view_count` is defined as integer-only NumberField
- SQLite stores as INTEGER type
- Default value handling: field is optional, need to handle NULL case

## API Response Requirements

### Success Case (200)
- Return full link record with updated `view_count`
- Include all fields: id, url, title, description, view_count, tags, etc.
- Use PocketBase record serialization for consistency

### Error Case (404)
- Link ID doesn't exist
- Return JSON error: `{"error": "Link not found"}`

### Potential Error Cases
- Invalid ID format (return 400)
- Database errors (return 500)

## Security Considerations

- No authentication required (specified in acceptance criteria)
- Public endpoint, potential for abuse
- No rate limiting specified
- ID parameter validation needed to prevent SQL injection (PocketBase likely handles this)

## Implementation Approach

Based on research, the implementation will:

1. **Route Registration**: Follow existing pattern in `routes.go`
2. **Handler Function**: Create `registerLinksView()` and handler
3. **Parameter Extraction**: Extract `:id` from URL path
4. **Atomic SQL**: Use `UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?`
5. **Record Retrieval**: Fetch updated record using PocketBase record operations
6. **Response**: Return full record JSON or 404 error

## File Changes Required

1. **New file**: `routes/links_view.go` - Handler implementation
2. **Modified file**: `routes/routes.go` - Add registration call
3. **No changes**: Database schema (view_count field exists)
4. **No changes**: Migrations (links collection exists)

## Testing Considerations

- Test concurrent requests to verify atomicity
- Test with non-existent link ID (404)
- Test with invalid ID format
- Verify view_count increments correctly
- Verify returned record includes all expected fields

## Atomicity Verification

The SQL operation `UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?` is atomic because:
- Single SQL statement
- SQLite row-level locking
- `COALESCE` handles NULL case (defaults to 0)
- No read-modify-write cycle in application code
- Even under concurrent load, each increment is isolated

## API Contract

```
POST /api/links/:id/view
Response 200: {
  "id": "abc123",
  "url": "https://example.com",
  "title": "Example Link",
  "description": "...",
  "view_count": 43,  // incremented value
  "tags": ["tag1", "tag2"],
  "created": "2024-01-01T00:00:00Z",
  "updated": "2024-01-02T00:00:00Z"
}

Response 404: {
  "error": "Link not found"
}
```