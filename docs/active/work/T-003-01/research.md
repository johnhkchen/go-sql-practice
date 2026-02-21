# Research: Search Endpoint Implementation (T-003-01)

## Overview
Implementing `GET /api/links/search` endpoint with full-text search and optional tag filtering. This requires understanding PocketBase's DAO layer and how to write SQL queries against its underlying SQLite database.

## Existing Codebase Structure

### Application Architecture
- **main.go**: Entry point, creates PocketBase app, registers migrations and routes
- **routes/**: Custom route registration package
  - **routes.go**: Main registration function using `OnServe` hook
  - **health.go**: Health endpoint using middleware pattern with `hook.Handler`
  - **sync_sessions.go**: CRUD operations for sync sessions
- **migrations/**: Database setup and seed data
  - **collections.go**: Creates collections (tags, links, sync_sessions)
  - **seed.go**: Populates initial data for development

### Database Schema

#### Links Collection
- **url**: URL field, required
- **title**: Text field, required, 1-500 chars
- **description**: Text field, optional, max 2000 chars
- **view_count**: Number field, optional, integer >= 0
- **tags**: Relation field to tags collection, max 100 selections
- **created_by**: Optional relation to users collection

#### Tags Collection
- **name**: Text field, required, 1-100 chars
- **slug**: Text field, required, unique, kebab-case pattern
- Has unique index on slug field

### Current Route Patterns

#### Route Registration (routes/routes.go)
```go
app.OnServe().BindFunc(func(e *core.ServeEvent) error {
    registerHealth(e)
    registerSyncSessions(e)
    return e.Next()
})
```

#### HTTP Handler Pattern (routes/sync_sessions.go)
- Uses `e.Router.POST()`, `e.Router.PATCH()` for route definition
- Handler functions receive `*core.RequestEvent` and `core.App`
- Response via `e.JSON(statusCode, data)`
- Request parsing via `json.NewDecoder(e.Request.Body)`
- Path parameters via `e.Request.PathValue("id")`

### Data Access Patterns

#### Record Creation
```go
record := core.NewRecord(collection)
record.Set("field", value)
app.Save(record)
```

#### Record Retrieval
- By ID: `app.FindRecordById("collection", id)`
- By field: `app.FindFirstRecordByData("collection", "field", value)`
- Collection lookup: `app.FindCollectionByNameOrId("name")`

### Missing Components for Search

1. **No existing SQL query examples**: The codebase uses high-level PocketBase record APIs, not direct SQL
2. **No DAO usage**: Need to access PocketBase's DAO for custom SQL queries
3. **No pagination implementation**: Will need to implement page/perPage handling
4. **No query parameter parsing**: Need to extract and validate `q`, `tag`, `page`, `perPage` params

## PocketBase Context

### DAO Layer Access
PocketBase provides `app.DB()` or `app.DAO()` for direct database access:
- `app.DB()`: Returns the underlying *sql.DB instance
- `app.DAO()`: Returns the data access object with query builder

### Full-Text Search in SQLite
SQLite supports FTS (Full-Text Search) but PocketBase collections use regular tables. Options:
1. Use LIKE queries for simple text matching
2. Create FTS virtual tables (complex, requires migrations)
3. Use combination of LIKE and INSTR for flexible searching

### Relation Handling
Tags are stored as relation IDs in the links table. To filter by slug:
1. Join with tags table
2. Match tag.slug with the provided filter
3. Handle multiple tags per link

## Key Constraints

1. **SQL Injection Prevention**: Must use parameterized queries
2. **Tag Filtering by Slug**: Not by ID, requires join with tags table
3. **Response Format**: Strict JSON structure with items, page, perPage, totalItems
4. **Empty Results**: Return 200 with empty items array, not 404
5. **Pagination Defaults**: page=1, perPage=20 if not specified

## Dependencies Analysis

- **T-001-02** (custom-health-route): ✅ Complete - Established route registration pattern
- **T-002-02** (seed-data-migration): ✅ Complete - Database has test data

## Technical Considerations

1. **Query Building**: Need to dynamically build WHERE clauses based on presence of `q` and `tag` params
2. **Performance**: Full-text search on title/description without FTS index may be slow for large datasets
3. **Escaping**: LIKE patterns need proper escaping for user input (%, _, \)
4. **Count Query**: Need separate query for totalItems to support pagination
5. **Case Sensitivity**: SQLite LIKE is case-insensitive by default

## Files to Create/Modify

1. **routes/links_search.go**: New file for search endpoint handler
2. **routes/routes.go**: Add registration for new endpoint
3. No model files needed - using PocketBase's record system
4. No migration needed - using existing collections