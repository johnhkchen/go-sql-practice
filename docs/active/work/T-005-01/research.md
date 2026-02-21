# Research: Go Integration Tests for Custom API Endpoints

## Codebase Overview

The go-sql-practice project is a PocketBase-based Go application that provides a link management system with custom API endpoints. The project structure follows a modular pattern with clear separation of concerns.

### Project Structure

```
├── main.go                    # Application entry point, PocketBase setup
├── go.mod/go.sum             # Go module dependencies (Go 1.26, PocketBase v0.36.5)
├── Makefile                  # Build automation with `test: go test ./...`
├── routes/                   # Custom HTTP route handlers
│   ├── routes.go            # Route registration central hub
│   ├── stats.go             # Statistics endpoint implementation
│   ├── links_search.go      # Search endpoint with pagination
│   ├── links_view.go        # View count increment endpoint
│   └── [other routes]       # Additional routes (health, sync, presentations, static)
├── migrations/              # Database schema and seed data
│   ├── collections.go       # PocketBase collection definitions
│   └── seed.go             # Development seed data
└── internal/                # Internal packages
    └── frontend/embed.go    # Frontend static files embedding
```

### Target API Endpoints for Testing

Based on ticket dependencies (T-003-01, T-003-02, T-003-03), three custom API endpoints need testing:

1. **Search Endpoint**: `GET /api/links/search`
   - File: `routes/links_search.go`
   - Features: Query parameter parsing, text search, tag filtering, pagination
   - Parameters: `q` (text search), `tag` (slug filter), `page`, `perPage`
   - Response: JSON with items array, pagination metadata

2. **Stats Endpoint**: `GET /api/stats`
   - File: `routes/stats.go`
   - Features: Aggregated statistics from database
   - Response: Total counts (links, tags, views), top tags, most viewed links

3. **View Count Endpoint**: `POST /api/links/:id/view`
   - File: `routes/links_view.go`
   - Features: Atomic increment of view_count, 404 for nonexistent links
   - Response: Updated link record JSON

### Database Schema

The application uses PocketBase's collection system with SQLite as the underlying database:

#### Collections:
- **links**: Primary content collection
  - Fields: `url` (URL), `title` (text), `description` (text), `view_count` (number), `tags` (relation to tags), `created_by` (relation to users)
  - Indexes: Default PocketBase indexes

- **tags**: Tag taxonomy for links
  - Fields: `name` (text), `slug` (text, unique)
  - Unique constraint: `idx_tags_slug ON tags (slug)`

- **sync_sessions**: Session management (additional)
- **presentations**: Presentation data (additional)

#### Seed Data Pattern:
- 8 predefined tags (golang, javascript, database, devops, frontend, backend, testing, architecture)
- 10 example links with varying view counts (8-45) and tag associations
- Seed check mechanism: Tests for "golang" tag existence

### Current Testing Infrastructure

**Current State**: No existing Go test files found (`*_test.go`)
- Makefile includes `test: go test ./...` target
- No test packages or helper utilities present
- No existing PocketBase testing patterns in codebase

**PocketBase Testing Approach**:
PocketBase applications support in-memory testing via `pocketbase.NewWithConfig()` with in-memory database configuration. This allows full HTTP handler stack testing without external processes.

### Route Registration Pattern

All routes are registered through a centralized system:
1. `main.go` calls `routes.Register(app)`
2. `routes/routes.go` contains `Register()` function that binds to `app.OnServe()`
3. Individual route files export registration functions (e.g., `registerStats`, `registerLinksSearch`, `registerLinksView`)
4. Each endpoint uses PocketBase's `core.RequestEvent` for HTTP handling

### Dependencies and Constraints

**Go Dependencies**:
- `github.com/pocketbase/pocketbase v0.36.5` - Core framework
- `github.com/labstack/echo/v5 v5.0.4` - HTTP router (PocketBase dependency)
- Standard library for testing support

**Database Dependencies**:
- SQLite via PocketBase's dbx abstraction
- JSON field support for tag relations (`json_each()` queries)
- Migration system handles collection creation automatically

**Testing Requirements**:
- In-memory PocketBase app instance setup
- Collection creation before test execution
- Seed data population for realistic test scenarios
- HTTP request simulation to test full stack
- No external PocketBase server dependency

### API Endpoint Implementation Patterns

**Common Patterns Across Endpoints**:
1. Parameter parsing from `http.Request`
2. Validation with structured error responses
3. Database queries using `app.DB().NewQuery()`
4. JSON response formatting via `e.JSON(statusCode, data)`
5. Error handling with appropriate HTTP status codes

**Search Endpoint Specifics**:
- Complex SQL queries with optional JOINs for tag filtering
- Manual SQL parameter interpolation (security consideration)
- Pagination with LIMIT/OFFSET
- Separate count query for total items
- Post-processing to attach tags to results

**Stats Endpoint Specifics**:
- Multiple independent SQL queries (total counts, top items)
- Aggregation queries with GROUP BY and ORDER BY
- Structured response with nested data types

**View Count Endpoint Specifics**:
- Atomic SQL UPDATE with COALESCE for null handling
- RowsAffected check for existence validation
- PocketBase DAO integration for response formatting

### Current Gaps for Testing

1. **Test Package Structure**: No established pattern for test organization
2. **Test Helpers**: No utilities for PocketBase app setup or data seeding
3. **HTTP Testing**: No patterns for request/response validation
4. **Data Fixtures**: Need structured test data separate from development seeds
5. **Test Database**: Need in-memory database configuration pattern
6. **Assertion Libraries**: Standard `testing` package only

### Integration Points

**PocketBase Core Integration**:
- App lifecycle management (`pocketbase.New()`, `app.Start()`)
- Database access via `app.DB()`
- Collection management through `app.FindCollectionByNameOrId()`
- Record operations via `daos.New(app.DB())`

**HTTP Layer Integration**:
- Router access via `core.ServeEvent`
- Request handling through `core.RequestEvent`
- Middleware chain participation (`e.Next()`)

**Migration System Integration**:
- Collection creation in `migrations.Register()`
- Seed data population automatically on startup
- Existence checks prevent duplicate seeding

This research provides the foundation for designing integration tests that can exercise the complete HTTP handling stack, from request parsing through database operations to response formatting, while maintaining isolation through in-memory database instances.