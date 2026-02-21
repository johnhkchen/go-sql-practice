# Plan: Stats Endpoint (T-003-03)

## Implementation Steps

### Step 1: Create stats.go with type definitions

Create `routes/stats.go` with:
- Package declaration and imports
- StatsResponse, TagStats, LinkStats structs with JSON tags
- Empty registerStats function

**Verification:** File compiles with `go build`

### Step 2: Implement basic endpoint registration

Add:
- registerStats function that registers GET /api/stats
- handleGetStats stub that returns empty response

**Verification:**
- Endpoint responds with 200 and empty JSON
- `curl localhost:8090/api/stats` returns `{}`

### Step 3: Add route registration to routes.go

Modify `routes/routes.go`:
- Add `registerStats(e)` call after other registrations

**Verification:**
- Server starts without errors
- `/api/stats` endpoint is accessible

### Step 4: Implement getTotalLinks query

Add:
- SQL constant for COUNT(*) query
- getTotalLinks function
- Integration in handleGetStats

**Verification:**
- Returns correct count from seeded data (10 links)
- Returns 0 if database is cleared

### Step 5: Implement getTotalTags query

Add:
- SQL constant for tags COUNT
- getTotalTags function
- Integration in handleGetStats

**Verification:**
- Returns correct count from seeded data (8 tags)

### Step 6: Implement getTotalViews query

Add:
- SQL constant with COALESCE for NULL handling
- getTotalViews function using SUM
- Integration in handleGetStats

**Verification:**
- Returns correct sum from seeded data
- Returns 0 for empty database

### Step 7: Implement getMostViewed query

Add:
- SQL constant with ORDER BY and LIMIT
- getMostViewed function returning []LinkStats
- Handle NULL view_count with COALESCE

**Verification:**
- Returns top 5 links ordered by view_count DESC
- MDN Web Docs (45 views) should be first
- Returns empty array for empty database

### Step 8: Implement getTopTags query

Add:
- Complex SQL with JSON relation handling
- getTopTags function with tag-link count logic
- Use SQLite JSON functions for relation matching

**Verification:**
- Returns top 5 tags by link count
- Each tag shows correct link_count
- Returns empty array for empty database

### Step 9: Complete handler with all queries

Update handleGetStats to:
- Call all query functions
- Assemble complete StatsResponse
- Handle any database errors with 500 response

**Verification:**
- Complete response matches acceptance criteria
- All fields populated correctly

### Step 10: Test empty database scenario

- Temporarily clear database or test with fresh instance
- Verify endpoint returns 200 with zeroed/empty values

**Verification:**
- total_links: 0
- total_tags: 0
- total_views: 0
- top_tags: []
- most_viewed: []

### Step 11: Test with production-like data

- Verify performance with larger dataset if available
- Check SQL query execution times
- Confirm no memory issues

**Verification:**
- Response time < 100ms for seeded data
- No memory spikes

### Step 12: Error handling verification

- Test database connection failures
- Verify 500 responses for errors
- Check error logging

**Verification:**
- Graceful error responses
- No panics or crashes

## Testing Strategy

### Unit Testing
- Mock dbx.DB interface for query functions
- Test each query function independently
- Verify SQL syntax is valid

### Integration Testing
- Test against real PocketBase instance
- Verify with seeded data
- Test empty database case
- Verify JSON response structure

### Manual Testing
```bash
# Test the endpoint
curl -X GET http://localhost:8090/api/stats | jq

# Expected response structure
{
  "total_links": 10,
  "total_tags": 8,
  "total_views": 234,
  "top_tags": [...],
  "most_viewed": [...]
}
```

## Commit Strategy

1. Initial stats endpoint structure
2. Add simple count queries (links, tags, views)
3. Add most_viewed query with ordering
4. Add complex top_tags query with JSON handling
5. Complete error handling and empty DB support

## Risk Mitigation

### JSON Relation Complexity
- Start with simple substring matching
- Optimize with JSON functions if needed
- Have fallback query approach ready

### Performance
- Add query timeouts if needed
- Consider caching for production (not in scope)
- Monitor query execution plans

### Empty Database
- Test early and often with empty DB
- Ensure all COALESCE statements work
- Verify empty arrays are returned correctly

## Dependencies

- Completed: T-001-02 (route pattern)
- Completed: T-002-02 (seed data)
- No blocking dependencies

## Success Criteria

- [x] Endpoint returns 200 status
- [x] All five statistics are computed via SQL
- [x] Empty database returns valid response
- [x] Response matches exact JSON structure
- [x] No records loaded into memory
- [x] Performance acceptable (< 100ms)

## Summary

The plan breaks implementation into testable steps, starting with structure and simple queries, then adding complexity. Each step has clear verification criteria. The most complex part (top_tags with JSON relations) is isolated in Step 8, allowing earlier steps to provide value even if that proves challenging.