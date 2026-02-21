# Structure: Search Endpoint Implementation (T-003-01)

## File Structure

### New Files

#### routes/links_search.go
Primary handler file containing the search endpoint implementation.

**Exports:**
- `registerLinksSearch(*core.ServeEvent)` - Registration function

**Internal Functions:**
- `handleSearch(*core.RequestEvent, core.App) error` - Main handler
- `parseSearchParams(*http.Request) SearchParams` - Parameter extraction
- `validateSearchParams(SearchParams) error` - Input validation
- `escapeLikePattern(string) string` - SQL LIKE escaping
- `executeSearchQuery(core.App, SearchParams) ([]LinkItem, error)` - Query execution
- `executeCountQuery(core.App, SearchParams) (int, error)` - Total count
- `fetchTagsForLinks(core.App, []string) (map[string][]string, error)` - Tag lookup

**Types:**
```go
type SearchParams struct {
    Q       string
    Tag     string
    Page    int
    PerPage int
}

type SearchResponse struct {
    Items      []LinkItem `json:"items"`
    Page       int        `json:"page"`
    PerPage    int        `json:"perPage"`
    TotalItems int        `json:"totalItems"`
}

type LinkItem struct {
    ID          string   `json:"id"`
    URL         string   `json:"url"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    ViewCount   int      `json:"view_count"`
    Tags        []string `json:"tags"`
}
```

### Modified Files

#### routes/routes.go
Add registration call for the new search endpoint.

**Changes:**
- Add `registerLinksSearch(e)` call in the `Register` function

## Module Organization

### Package Structure
```
routes/
├── routes.go          # Main registration (MODIFIED)
├── health.go          # Existing health endpoint
├── sync_sessions.go   # Existing sync sessions
└── links_search.go    # NEW: Search endpoint
```

### Import Dependencies
```go
// routes/links_search.go
import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "github.com/pocketbase/pocketbase/core"
)
```

## Function Boundaries

### Public API
- `GET /api/links/search` - Single public endpoint

### Internal Architecture

#### Layer 1: HTTP Handler
- Request parsing
- Parameter validation
- Response formatting
- Error handling

#### Layer 2: Query Builders
- SQL query construction
- Parameter binding
- LIKE pattern escaping

#### Layer 3: Database Access
- Direct SQL execution via `app.DB()`
- Result set scanning
- Transaction management (if needed)

## Data Flow

1. **Request Entry**
   ```
   HTTP Request → registerLinksSearch → handleSearch
   ```

2. **Parameter Processing**
   ```
   handleSearch → parseSearchParams → validateSearchParams
   ```

3. **Query Execution**
   ```
   handleSearch → executeSearchQuery → app.DB().Query()
                → executeCountQuery → app.DB().QueryRow()
                → fetchTagsForLinks → app.DB().Query()
   ```

4. **Response Assembly**
   ```
   Query Results → LinkItem structs → SearchResponse → JSON
   ```

## SQL Query Structure

### Main Search Query
```sql
SELECT DISTINCT l.id, l.url, l.title, l.description, l.view_count
FROM links l
LEFT JOIN links_tags lt ON l.id = lt.link_id
LEFT JOIN tags t ON lt.tag_id = t.id
WHERE [conditions]
ORDER BY l.created DESC
LIMIT ? OFFSET ?
```

### Count Query
```sql
SELECT COUNT(DISTINCT l.id)
FROM links l
LEFT JOIN links_tags lt ON l.id = lt.link_id
LEFT JOIN tags t ON lt.tag_id = t.id
WHERE [conditions]
```

### Tags Fetch Query
```sql
SELECT lt.link_id, t.slug
FROM links_tags lt
JOIN tags t ON lt.tag_id = t.id
WHERE lt.link_id IN (?)
```

## Interface Contracts

### HTTP Request
```
GET /api/links/search?q=search_term&tag=tag_slug&page=1&perPage=20
```

### HTTP Response
```json
{
    "items": [...],
    "page": 1,
    "perPage": 20,
    "totalItems": 42
}
```

## Error Boundaries

### Validation Errors (400)
- Invalid page number (< 1)
- Invalid perPage (< 1 or > 100)
- Malformed parameters

### Database Errors (500)
- Connection failures
- Query execution errors
- Unexpected schema issues

### Success Cases (200)
- Results found
- Empty results (items: [])

## Security Boundaries

### Input Sanitization
- All SQL queries use parameterized statements
- LIKE wildcards escaped in user input
- Tag slug validated against pattern

### Resource Limits
- Max perPage: 100
- Default perPage: 20
- Query timeout: Rely on database defaults

## Testing Structure

### Unit Test Files
- `routes/links_search_test.go` - Handler tests

### Test Coverage Areas
1. Parameter parsing and validation
2. SQL query construction
3. LIKE pattern escaping
4. Pagination logic
5. Error handling

## Configuration Constants

```go
const (
    DefaultPage    = 1
    DefaultPerPage = 20
    MaxPerPage     = 100
)
```

## Ordering Dependencies

1. First: Create `routes/links_search.go` with all functionality
2. Then: Modify `routes/routes.go` to register the endpoint
3. Finally: Test the complete integration

This structure maintains separation of concerns while keeping the implementation focused and testable.