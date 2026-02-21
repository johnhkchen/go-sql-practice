# Research: Stats Endpoint (T-003-03)

## Context

Implementing `GET /api/stats` to return aggregate statistics using SQL aggregate functions (COUNT, SUM, GROUP BY, ORDER BY, LIMIT) against PocketBase's SQLite database.

## Current Codebase State

### Application Architecture

The application is built on PocketBase v0.36.5, using it as a framework rather than a standalone binary. Entry point is `main.go` which:
- Creates a PocketBase app instance
- Registers migrations (collections and seed data)
- Registers custom routes
- Starts the server

### Route Registration Pattern

Routes are defined in the `routes/` package:
- `routes/routes.go`: Main registration function using `app.OnServe()` hook
- `routes/health.go`: Custom health endpoint using middleware pattern with priority
- `routes/sync_sessions.go`: CRUD operations for sync sessions

The pattern establishes custom routes via:
```go
e.Router.POST("/api/path", func(ev *core.RequestEvent) error { ... })
e.Router.GET("/api/path", func(ev *core.RequestEvent) error { ... })
```

Response handling uses `ev.JSON(statusCode, data)` for JSON responses.

### Database Schema

Two main collections exist (created via migrations):

**tags**:
- `name`: string (required, 1-100 chars)
- `slug`: string (required, unique, kebab-case)
- Fields: id, created, updated (system fields)

**links**:
- `url`: URL field (required)
- `title`: string (required, 1-500 chars)
- `description`: string (optional, max 2000)
- `view_count`: integer (optional, min 0)
- `tags`: relation to tags (many-to-many, max 100)
- `created_by`: relation to users (optional)
- Fields: id, created, updated (system fields)

### Seed Data

The migrations/seed.go creates:
- 8 tags: golang, javascript, database, devops, frontend, backend, testing, architecture
- 10 links with varying view_counts (0-45) and multiple tag associations

### Database Access Patterns

PocketBase provides multiple database access methods:

1. **Record-based API** (sync_sessions.go uses this):
   - `app.FindCollectionByNameOrId(name)`
   - `app.FindRecordById(collection, id)`
   - `core.NewRecord(collection)`
   - `app.Save(record)`

2. **Raw SQL via DAO**:
   - `app.DB()`: Returns concurrent connection pool
   - `app.NonconcurrentDB()`: Single connection for writes
   - `app.Dao().DB()`: Access to dbx.DB for raw queries

From web search, raw SQL execution pattern:
```go
results, err := app.Dao().DB().
    NewQuery("SELECT * FROM table").
    Execute()
```

### SQLite Table Structure

PocketBase collections map to SQLite tables with naming convention:
- Collection "tags" → table "tags"
- Collection "links" → table "links"
- System fields: id (15-char string), created (datetime), updated (datetime)
- Relations stored as JSON arrays of IDs

### Required Statistics

The endpoint must compute:
1. `total_links`: COUNT(*) from links
2. `total_tags`: COUNT(*) from tags
3. `total_views`: SUM(view_count) from links
4. `top_tags`: Tags with link counts, sorted by count DESC, limit 5
5. `most_viewed`: Links sorted by view_count DESC, limit 5

### SQL Challenges

- Relations in PocketBase are stored as JSON arrays, complicating JOIN operations
- Tag link counts require parsing JSON relation fields or using JSON functions
- Need to handle NULL values in view_count (coalesce to 0)
- Empty database should return zeroed counts and empty arrays

### Dependencies

This ticket depends on:
- T-001-02: Custom health route pattern (completed)
- T-002-02: Seed data migration (completed)

Both are done, establishing the route pattern and providing data to query.

### File Organization

Existing pattern shows:
- Route files in `routes/` package
- Each feature gets its own file (health.go, sync_sessions.go)
- Registration happens in routes.go via register* functions
- Constants and types defined at top of route files

## Summary

The codebase provides a clear pattern for adding custom routes and accessing the database. The main challenge will be writing efficient SQL queries against PocketBase's schema, particularly for computing tag link counts given the JSON-stored relations. Raw SQL via `app.Dao().DB()` appears to be the appropriate approach for aggregate queries.