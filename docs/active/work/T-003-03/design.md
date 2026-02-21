# Design: Stats Endpoint (T-003-03)

## Requirements Recap

Implement `GET /api/stats` that returns:
- `total_links`: Total count of all links
- `total_tags`: Total count of all tags
- `total_views`: Sum of all view_counts
- `top_tags`: Top 5 tags by link count (name, slug, link_count)
- `most_viewed`: Top 5 links by view_count (id, title, url, view_count)

All via SQL queries, not loading records into memory.

## Design Options

### Option 1: PocketBase Record API with Aggregation

Use PocketBase's FindAllRecords with pagination to iterate and aggregate in Go.

**Pros:**
- Uses established PocketBase patterns
- Type-safe with Record models
- Handles JSON relations automatically

**Cons:**
- Violates requirement of SQL-only aggregation
- Inefficient for large datasets
- Loads all records into memory

**Verdict:** Rejected - explicitly violates "not by loading all records into memory" requirement.

### Option 2: Raw SQL with app.Dao().DB()

Execute raw SQL queries directly against SQLite using PocketBase's DAO layer.

**Pros:**
- Direct SQL execution as required
- Maximum performance
- Full control over query optimization
- Can use SQLite's JSON functions for relations

**Cons:**
- Need to handle PocketBase's JSON relation format
- Must map SQL results to Go structs manually
- Bypasses PocketBase's data layer abstractions

**Verdict:** Strong candidate - meets all requirements.

### Option 3: Hybrid - Simple Counts via SQL, Relations via Records

Use SQL for simple counts/sums, but load records for tag link counts.

**Pros:**
- Simpler handling of JSON relations
- Mix optimal approaches per query

**Cons:**
- Inconsistent approach
- Still loads some records into memory
- More complex code with two patterns

**Verdict:** Rejected - inconsistent and partially violates requirements.

## Chosen Approach: Raw SQL with app.Dao().DB()

### Rationale

Option 2 is the only approach that fully satisfies the requirement to compute everything via SQL. While it requires handling PocketBase's JSON relation format, SQLite provides JSON functions to work with this.

### SQL Query Designs

**total_links:**
```sql
SELECT COUNT(*) as total FROM links
```

**total_tags:**
```sql
SELECT COUNT(*) as total FROM tags
```

**total_views:**
```sql
SELECT COALESCE(SUM(view_count), 0) as total FROM links
```
Note: COALESCE handles NULL when no records exist.

**top_tags:**
The challenge: links.tags is stored as JSON array of tag IDs.
```sql
SELECT
    t.name,
    t.slug,
    (
        SELECT COUNT(*)
        FROM links l
        WHERE json_extract(l.tags, '$') LIKE '%' || t.id || '%'
    ) as link_count
FROM tags t
ORDER BY link_count DESC
LIMIT 5
```
Alternative using JSON functions:
```sql
SELECT
    t.name,
    t.slug,
    COUNT(DISTINCT l.id) as link_count
FROM tags t
LEFT JOIN links l ON json_extract(l.tags, '$') LIKE '%' || t.id || '%'
GROUP BY t.id, t.name, t.slug
ORDER BY link_count DESC
LIMIT 5
```

**most_viewed:**
```sql
SELECT
    id,
    title,
    url,
    COALESCE(view_count, 0) as view_count
FROM links
ORDER BY view_count DESC
LIMIT 5
```

### Response Structure

```go
type StatsResponse struct {
    TotalLinks  int64         `json:"total_links"`
    TotalTags   int64         `json:"total_tags"`
    TotalViews  int64         `json:"total_views"`
    TopTags     []TagStats    `json:"top_tags"`
    MostViewed  []LinkStats   `json:"most_viewed"`
}

type TagStats struct {
    Name       string `json:"name"`
    Slug       string `json:"slug"`
    LinkCount  int64  `json:"link_count"`
}

type LinkStats struct {
    ID         string `json:"id"`
    Title      string `json:"title"`
    URL        string `json:"url"`
    ViewCount  int64  `json:"view_count"`
}
```

### Error Handling

- Database connection failures: Return 500
- SQL execution errors: Return 500
- Empty database: Return 200 with zero/empty values
- Malformed data: Skip row, log warning

### Database Access Pattern

Based on research, use:
```go
db := app.Dao().DB()
rows, err := db.NewQuery(sql).Execute()
```

For concurrent reads, this should use the appropriate connection pool automatically. No writes are performed, so concurrency concerns are minimal.

## Trade-offs

**Chosen approach trades:**
- PocketBase abstraction for SQL control
- Type safety for performance
- Simplicity for requirement compliance

**Mitigations:**
- Extensive testing with empty/full databases
- Clear SQL comments explaining JSON handling
- Structured error messages for debugging

## Summary

Raw SQL via app.Dao().DB() is the clear choice to meet the "all data computed via SQL" requirement. The main complexity lies in handling PocketBase's JSON-stored relations for tag link counts, which SQLite's JSON functions can handle adequately.