# Plan: Search Endpoint Implementation (T-003-01)

## Implementation Steps

### Step 1: Create Core Types and Constants
**File:** `routes/links_search.go`

Create the basic types and constants needed for the endpoint:
- SearchParams struct
- SearchResponse struct
- LinkItem struct
- Constants (DefaultPage, DefaultPerPage, MaxPerPage)

**Verification:** File compiles with `go build`

### Step 2: Implement Parameter Parsing
**File:** `routes/links_search.go`

Implement `parseSearchParams` function:
- Extract query parameters from request
- Set defaults for page and perPage
- Return SearchParams struct

**Verification:** Unit test for various parameter combinations

### Step 3: Implement Parameter Validation
**File:** `routes/links_search.go`

Implement `validateSearchParams` function:
- Check page >= 1
- Check perPage between 1 and MaxPerPage
- Return error for invalid params

**Verification:** Unit test for validation edge cases

### Step 4: Implement LIKE Pattern Escaping
**File:** `routes/links_search.go`

Implement `escapeLikePattern` function:
- Escape % characters
- Escape _ characters
- Escape \ characters
- Add surrounding % for partial match

**Verification:** Unit test with special characters

### Step 5: Implement Main Search Query
**File:** `routes/links_search.go`

Implement `executeSearchQuery` function:
- Build SQL query dynamically based on params
- Execute parameterized query
- Scan results into LinkItem structs
- Handle database errors

**Verification:** Manual test with curl against running server

### Step 6: Implement Count Query
**File:** `routes/links_search.go`

Implement `executeCountQuery` function:
- Build count query with same WHERE conditions
- Execute and return total count
- Handle database errors

**Verification:** Ensure count matches filtered results

### Step 7: Implement Tag Fetching
**File:** `routes/links_search.go`

Implement `fetchTagsForLinks` function:
- Query tags for given link IDs
- Return map of link_id to tag slugs
- Handle empty results

**Verification:** Test with links having multiple tags

### Step 8: Implement Main Handler
**File:** `routes/links_search.go`

Implement `handleSearch` function:
- Parse and validate parameters
- Execute search query with pagination
- Execute count query
- Fetch tags for results
- Assemble and return response

**Verification:** End-to-end test of complete flow

### Step 9: Implement Registration Function
**File:** `routes/links_search.go`

Implement `registerLinksSearch` function:
- Register GET route at /api/links/search
- Connect to handleSearch function

**Verification:** Route appears in server logs

### Step 10: Register Route in Main Router
**File:** `routes/routes.go`

Modify Register function:
- Add call to registerLinksSearch(e)

**Verification:** Endpoint accessible via HTTP

### Step 11: Integration Testing
**Manual Testing via curl**

Test all acceptance criteria:
1. Text search: `curl "localhost:8091/api/links/search?q=golang"`
2. Tag filter: `curl "localhost:8091/api/links/search?tag=database"`
3. Combined: `curl "localhost:8091/api/links/search?q=doc&tag=golang"`
4. Pagination: `curl "localhost:8091/api/links/search?page=2&perPage=5"`
5. Empty results: `curl "localhost:8091/api/links/search?q=nonexistent"`

**Verification:** All acceptance criteria pass

### Step 12: SQL Injection Testing
**Manual Testing**

Test SQL injection prevention:
1. Test with quotes: `?q=test' OR '1'='1`
2. Test with semicolons: `?q=test; DROP TABLE links`
3. Test with comments: `?q=test--`

**Verification:** No SQL errors, queries properly escaped

## Testing Strategy

### Unit Tests
Location: `routes/links_search_test.go`

1. **Parameter Parsing Tests**
   - Default values
   - Custom values
   - Missing parameters

2. **Validation Tests**
   - Invalid page numbers
   - Invalid perPage values
   - Boundary conditions

3. **Escaping Tests**
   - Special SQL characters
   - Unicode characters
   - Empty strings

### Integration Tests
Manual testing with curl commands:

1. **Search Functionality**
   - Title matching
   - Description matching
   - Case-insensitive search

2. **Tag Filtering**
   - Existing tags
   - Non-existent tags
   - Multiple tags on same link

3. **Pagination**
   - First page
   - Last page
   - Beyond last page

4. **Combined Queries**
   - Search + tag
   - Search + pagination
   - All parameters together

### Error Cases

1. **Database Errors**
   - Simulate connection failure (stop server)
   - Test recovery

2. **Invalid Input**
   - Negative page numbers
   - Huge perPage values
   - Malformed parameters

## Commit Strategy

### Commit 1: Core implementation
- Add routes/links_search.go with all functionality
- Message: "feat: implement GET /api/links/search endpoint with SQL queries"

### Commit 2: Route registration
- Modify routes/routes.go
- Message: "feat: register search endpoint in router"

### Commit 3: Testing (if test files added)
- Add routes/links_search_test.go
- Message: "test: add unit tests for search endpoint"

## Rollback Plan

If issues arise:
1. Remove registration from routes/routes.go
2. Delete routes/links_search.go
3. Server continues to run with existing endpoints

## Performance Considerations

### Current Implementation
- LIKE queries adequate for seed data (10 links)
- Sub-second response expected

### Future Optimizations
- Add index on title/description (future)
- Consider FTS5 for large datasets (future)
- Add caching layer if needed (future)

## Security Checklist

- [x] All queries use parameter binding
- [x] LIKE patterns properly escaped
- [x] Input validation in place
- [x] Error messages don't leak schema
- [x] Rate limiting (rely on PocketBase defaults)

## Dependencies

- No new Go dependencies required
- Uses only standard library and PocketBase
- Compatible with existing SQLite database

## Success Criteria

Endpoint is complete when:
1. All acceptance criteria pass
2. No SQL injection vulnerabilities
3. Proper error handling
4. Clean code structure
5. Follows existing patterns