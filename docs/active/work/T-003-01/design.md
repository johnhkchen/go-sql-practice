# Design: Search Endpoint Implementation (T-003-01)

## Objective
Design the `GET /api/links/search` endpoint that provides full-text search on links with optional tag filtering, implementing proper SQL queries via PocketBase's DAO layer.

## Design Options

### Option 1: Direct SQL with app.DB()
Use PocketBase's raw database connection to execute parameterized SQL queries directly.

**Pros:**
- Full control over SQL query structure
- Can optimize queries precisely
- Direct access to SQLite features

**Cons:**
- Manual result mapping to JSON
- Need to handle all escaping manually
- Bypasses PocketBase's record system

### Option 2: PocketBase DAO with Query Builder
Use PocketBase's DAO query builder methods for safer query construction.

**Pros:**
- Built-in SQL injection protection
- Automatic result mapping
- Consistent with PocketBase patterns

**Cons:**
- Limited to DAO's query builder capabilities
- May not support complex SQL operations
- Documentation is limited

### Option 3: Hybrid - Raw SQL with Record Mapping
Execute raw SQL queries but use PocketBase's record system for result handling.

**Pros:**
- SQL flexibility with PocketBase conveniences
- Can leverage existing record serialization
- Maintains consistency with collection structure

**Cons:**
- Requires understanding both raw SQL and PocketBase internals
- Potential for impedance mismatch

## Selected Approach: Option 1 - Direct SQL

### Rationale
1. **Learning Objective**: The ticket explicitly states this is "the primary practice target for writing SQL queries"
2. **Full Control**: Need complex JOINs for tag filtering by slug
3. **Performance**: Can optimize queries without abstraction overhead
4. **Transparency**: Clear SQL makes debugging and learning easier

## Implementation Design

### Query Structure

#### Base Query for Links with Tags
```sql
SELECT DISTINCT
    l.id,
    l.url,
    l.title,
    l.description,
    l.view_count,
    l.created,
    l.updated
FROM links l
LEFT JOIN links_tags lt ON l.id = lt.link_id
LEFT JOIN tags t ON lt.tag_id = t.id
```

#### Where Clause Construction
- Text search: `WHERE (l.title LIKE ? OR l.description LIKE ?)`
- Tag filter: `WHERE t.slug = ?`
- Combined: Both conditions with AND

#### Pagination
- Use `LIMIT ? OFFSET ?` for pagination
- Separate COUNT query for totalItems

### Request Parameter Handling

```go
type SearchParams struct {
    Q       string // search query
    Tag     string // tag slug filter
    Page    int    // page number (default: 1)
    PerPage int    // items per page (default: 20, max: 100)
}
```

### Response Structure

```go
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
    Tags        []string `json:"tags"` // tag slugs
}
```

### Security Considerations

1. **SQL Injection Prevention**
   - Use parameterized queries exclusively
   - Escape LIKE wildcards in user input
   - Validate all numeric parameters

2. **Input Validation**
   - Sanitize search query (remove/escape %, _, \)
   - Validate tag slug format
   - Enforce pagination limits

### Query Optimization

1. **Index Usage**
   - Leverage existing unique index on tags.slug
   - Consider full-text index for future optimization

2. **Query Efficiency**
   - Use DISTINCT to avoid duplicate rows from JOIN
   - Separate query for tags to avoid N+1 problem
   - Cache tag mappings if needed

## Rejected Alternatives

### Why Not FTS5?
- Requires migration to create virtual tables
- Overkill for current requirements
- LIKE queries sufficient for MVP

### Why Not PocketBase List API?
- Can't do complex SQL joins via REST API
- Would require client-side filtering
- Defeats purpose of practicing SQL

### Why Not ORM/Query Builder?
- Learning objective is raw SQL practice
- PocketBase's query builder not well documented
- Need precise control over JOIN logic

## Error Handling

1. **Database Errors**: Return 500 with generic message
2. **Invalid Parameters**: Return 400 with specific validation errors
3. **No Results**: Return 200 with empty items array
4. **Tag Not Found**: Return results as if no tag filter (empty result set)

## Testing Strategy

1. **Unit Tests**: Parameter validation, SQL escaping
2. **Integration Tests**: Query execution with test data
3. **Edge Cases**:
   - Empty search query
   - Non-existent tag
   - Special characters in search
   - Pagination boundaries

## Performance Considerations

1. **Current**: LIKE queries adequate for seed data size
2. **Future**: Consider FTS5 migration if dataset grows
3. **Caching**: Not needed for MVP, but tag lookup could be cached

## Migration Path

This design allows future migration to FTS5 or other search solutions:
1. Query logic isolated in single handler
2. Response format remains consistent
3. Can add search provider abstraction later