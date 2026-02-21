# Structure: Stats Endpoint (T-003-03)

## File Changes

### New Files

**routes/stats.go**
- Package: `routes`
- Purpose: Stats endpoint handler and SQL queries
- Exports: None (internal registration only)
- Structure:
  ```
  - Type definitions:
    - StatsResponse struct
    - TagStats struct
    - LinkStats struct
  - Functions:
    - registerStats(*core.ServeEvent)
    - handleGetStats(*core.RequestEvent, core.App) error
    - getTotalLinks(db) (int64, error)
    - getTotalTags(db) (int64, error)
    - getTotalViews(db) (int64, error)
    - getTopTags(db) ([]TagStats, error)
    - getMostViewed(db) ([]LinkStats, error)
  ```

### Modified Files

**routes/routes.go**
- Add `registerStats(e)` call in Register function
- Location: After `registerSyncSessions(e)` line

## Component Architecture

### Data Types

```go
type StatsResponse struct {
    TotalLinks  int64       `json:"total_links"`
    TotalTags   int64       `json:"total_tags"`
    TotalViews  int64       `json:"total_views"`
    TopTags     []TagStats  `json:"top_tags"`
    MostViewed  []LinkStats `json:"most_viewed"`
}

type TagStats struct {
    Name      string `json:"name"`
    Slug      string `json:"slug"`
    LinkCount int64  `json:"link_count"`
}

type LinkStats struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    URL       string `json:"url"`
    ViewCount int64  `json:"view_count"`
}
```

### Function Signatures

```go
// Route registration
func registerStats(e *core.ServeEvent)

// Main handler
func handleGetStats(e *core.RequestEvent, app core.App) error

// Query functions (all take dbx.DB interface)
func getTotalLinks(db dbx.DB) (int64, error)
func getTotalTags(db dbx.DB) (int64, error)
func getTotalViews(db dbx.DB) (int64, error)
func getTopTags(db dbx.DB) ([]TagStats, error)
func getMostViewed(db dbx.DB) ([]LinkStats, error)
```

### Database Access

- All queries use `app.Dao().DB()` for read operations
- Each query function receives the db interface for testability
- No transactions needed (read-only operations)
- Use `db.NewQuery(sql)` pattern from PocketBase's dbx layer

### SQL Query Storage

SQL queries defined as constants at package level:
```go
const (
    sqlTotalLinks = "SELECT COUNT(*) as total FROM links"
    sqlTotalTags  = "SELECT COUNT(*) as total FROM tags"
    sqlTotalViews = "SELECT COALESCE(SUM(view_count), 0) as total FROM links"
    sqlTopTags    = `...` // Complex query with JSON
    sqlMostViewed = `...` // ORDER BY with LIMIT
)
```

### Error Boundaries

- Database errors: Caught at query function level, returned to handler
- Handler aggregates all errors and returns 500 if any fail
- Empty results: Not an error, return zero/empty values
- JSON marshaling: Handled by ev.JSON(), unlikely to fail with our types

### Module Dependencies

```go
import (
    "github.com/pocketbase/pocketbase/core"
    "github.com/pocketbase/dbx"  // For DB interface
)
```

No new external dependencies required.

## Interface Contracts

### HTTP Interface

```
GET /api/stats
Response: 200 OK
{
  "total_links": 10,
  "total_tags": 8,
  "total_views": 234,
  "top_tags": [
    {"name": "Go", "slug": "golang", "link_count": 3},
    ...
  ],
  "most_viewed": [
    {"id": "abc123", "title": "MDN Web Docs", "url": "...", "view_count": 45},
    ...
  ]
}
```

### Internal Interfaces

The handler coordinates between:
1. HTTP layer (core.RequestEvent)
2. Database layer (app.Dao().DB())
3. Response serialization (ev.JSON())

Each query function is independent and can execute in parallel if needed.

## Ordering Constraints

1. Must import required packages first
2. Define types before using them
3. Define SQL constants before query functions
4. Register route after defining handler
5. Add registration call to routes.go last

## Design Decisions

### Query Isolation
Each statistic gets its own query function for:
- Clear testing boundaries
- Potential parallel execution
- Independent error handling

### SQL as Constants
SQL queries defined as constants rather than inline for:
- Clarity and maintenance
- Potential reuse
- Easier testing/mocking

### JSON Handling for Relations
Using SQLite's JSON functions with LIKE for tag matching:
- Works with PocketBase's array storage
- Avoids complex JSON parsing
- SQLite-native solution

### Error Strategy
Return 500 for any database error rather than partial results:
- Consistent error behavior
- Clear failure signal
- Prevents confusing partial data

## Summary

The structure follows established patterns in the codebase: a dedicated route file with type definitions, a registration function, and a main handler that coordinates database queries. Query logic is decomposed into testable functions, with SQL stored as constants for clarity.