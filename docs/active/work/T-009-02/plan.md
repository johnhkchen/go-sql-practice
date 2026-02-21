# T-009-02 Plan: SQL Injection and View Endpoint Fix

## Implementation Strategy

This plan sequences the security fixes and functionality enhancements in a way that minimizes risk and enables incremental verification. Each step can be committed independently and includes specific verification criteria.

## Pre-Implementation Setup

### Step 0: Environment Verification
**Duration**: 5 minutes
**Artifacts**: None

**Actions**:
1. Verify current working directory is `/home/jchen/repos/go-sql-practice`
2. Confirm Go environment and PocketBase dependency versions
3. Run existing tests to establish baseline: `go test ./routes/...`
4. Check current git status and ensure clean working state

**Verification Criteria**:
- All existing tests pass
- No compilation errors
- Go mod dependencies are satisfied
- Git working directory is clean or only contains expected changes

## Phase 1: Security Fixes (High Priority)

### Step 1: Fix SQL Injection in links_search_simple.go
**Duration**: 15 minutes
**Artifacts**: Modified `routes/links_search_simple.go`

**Actions**:
1. Add `"github.com/pocketbase/dbx"` import
2. Modify `handleSearchSimple()` function:
   - Replace direct `fmt.Sprintf` on line 82 with named placeholder system
   - Build `queryParams map[string]interface{}`
   - Use `db.NewQuery(query).Bind(dbx.Params(queryParams)).Rows()`
3. Update `escapeLikePattern2()` function:
   - Remove single quote escaping (handled by parameter binding)
   - Keep LIKE wildcard escaping for pattern matching
4. Test SQL injection prevention manually with malicious inputs

**Specific Code Changes**:
```go
// Before (line 82):
whereClauses = append(whereClauses, fmt.Sprintf("(title LIKE '%s' OR description LIKE '%s')", searchPattern, searchPattern))

// After:
whereClauses = append(whereClauses, "(title LIKE {searchPattern} OR description LIKE {searchPattern})")
queryParams["searchPattern"] = escapeLikePattern2(params.Q)
```

**Verification Criteria**:
- Code compiles without errors
- Existing search functionality works correctly
- SQL injection attempts fail safely
- Search results accuracy maintained
- Response format unchanged

**Test Commands**:
```bash
go test ./routes/ -run TestLinksSearchSimple
curl "http://localhost:8090/api/links/search-simple?q=test"
curl "http://localhost:8090/api/links/search-simple?q=%27;DROP%20TABLE%20links;--"
```

### Step 2: Fix SQL Injection in links_search.go
**Duration**: 25 minutes
**Artifacts**: Modified `routes/links_search.go`

**Actions**:
1. Add `"github.com/pocketbase/dbx"` import if not present
2. Fix `searchLinks()` function (lines ~200-230):
   - Replace manual parameter replacement loop with parameter map
   - Convert `args []interface{}` to named parameters
   - Update query placeholders from `?` to `{param1}`, `{param2}`, etc.
3. Fix `countSearchResults()` function (lines ~280-310):
   - Apply same parameter binding approach
   - Maintain identical WHERE clause logic
4. Fix `fetchTagsForLinks()` function (lines ~330-360):
   - Handle dynamic IN clause construction with named parameters
   - Create parameters for each link ID in the IN clause

**Specific Code Changes**:
```go
// Before (manual replacement):
for _, arg := range args {
    argStr := fmt.Sprintf("%v", arg)
    switch arg.(type) {
    case int:
        finalQuery = strings.Replace(finalQuery, "?", argStr, 1)
    default:
        argStr = strings.ReplaceAll(argStr, "'", "''")
        finalQuery = strings.Replace(finalQuery, "?", "'"+argStr+"'", 1)
    }
}

// After (parameter binding):
queryParams := make(map[string]interface{})
queryWithParams := query
for i, arg := range args {
    paramName := fmt.Sprintf("param%d", i+1)
    queryParams[paramName] = arg
    queryWithParams = strings.Replace(queryWithParams, "?", "{"+paramName+"}", 1)
}
rows, err := db.NewQuery(queryWithParams).Bind(dbx.Params(queryParams)).Rows()
```

**Verification Criteria**:
- Code compiles without errors
- All search functionality works (text search, tag filtering, pagination)
- SQL injection attempts blocked
- Search performance within acceptable range (< 10% regression)
- Response format and data accuracy maintained

**Test Commands**:
```bash
go test ./routes/ -run TestLinksSearch
curl "http://localhost:8090/api/links/search?q=test&tag=example"
curl "http://localhost:8090/api/links/search?q=%27UNION%20SELECT%20*%20FROM%20users--"
```

## Phase 2: Functionality Enhancement

### Step 3: Enhance links_view.go to Return Full Records
**Duration**: 20 minutes
**Artifacts**: Modified `routes/links_view.go`

**Actions**:
1. Add `"time"` import for timestamp formatting
2. Create `LinkViewResponse` struct for response formatting
3. Modify `handleLinksView()` function:
   - Keep existing UPDATE query (already secure with parameter binding)
   - Add record retrieval using `app.FindRecordById("links", linkId)`
   - Convert PocketBase record to response format
   - Handle tags separately if needed (check if tags are required in response)
4. Update error handling for record retrieval failures

**Specific Code Changes**:
```go
// Add response struct:
type LinkViewResponse struct {
    ID          string   `json:"id"`
    URL         string   `json:"url"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    ViewCount   int      `json:"view_count"`
    Tags        []string `json:"tags"`
    Created     string   `json:"created"`
    Updated     string   `json:"updated"`
}

// Modify handleLinksView() after UPDATE:
// Fetch updated record
record, err := app.FindRecordById("links", linkId)
if err != nil {
    return e.JSON(500, map[string]string{"error": "Failed to retrieve updated record"})
}

// Convert and return full record
response := LinkViewResponse{
    ID:          record.Id,
    URL:         record.GetString("url"),
    Title:       record.GetString("title"),
    Description: record.GetString("description"),
    ViewCount:   record.GetInt("view_count"),
    Created:     record.Created.Time().Format(time.RFC3339),
    Updated:     record.Updated.Time().Format(time.RFC3339),
}
```

**Verification Criteria**:
- View count increments correctly
- Full link record returned in response
- All link fields populated correctly
- Timestamps formatted properly
- Error handling for missing records works
- Original UPDATE security maintained

**Test Commands**:
```bash
go test ./routes/ -run TestLinksView
curl -X POST "http://localhost:8090/api/links/{existing-id}/view"
# Verify response contains full record, not just success message
```

## Phase 3: Testing and Validation

### Step 4: Security Testing
**Duration**: 15 minutes
**Artifacts**: Test documentation/notes

**Actions**:
1. Test common SQL injection payloads against all three endpoints:
   - `'; DROP TABLE links; --`
   - `' UNION SELECT * FROM users --`
   - `%' OR 1=1 --`
   - Format string attacks: `%s%s%s%d`
2. Verify parameterized queries prevent all injection attempts
3. Check that error messages don't reveal database structure
4. Confirm malicious inputs are handled gracefully

**Test Payloads**:
```bash
# Test against search-simple
curl "http://localhost:8090/api/links/search-simple?q=%27%3B%20DROP%20TABLE%20links%3B%20--"
curl "http://localhost:8090/api/links/search-simple?q=%25s%25s%25d"

# Test against search
curl "http://localhost:8090/api/links/search?q=%27%20UNION%20SELECT%20*%20FROM%20users%20--"

# Test against view (if it accepts parameters)
curl -X POST "http://localhost:8090/api/links/malicious%27input/view"
```

**Verification Criteria**:
- All SQL injection attempts fail safely
- Applications remain functional after injection attempts
- No database errors exposed to clients
- Database integrity maintained

### Step 5: Integration Testing
**Duration**: 10 minutes
**Artifacts**: Test results documentation

**Actions**:
1. Run complete test suite: `go test ./...`
2. Test end-to-end workflows:
   - Search for links with various queries
   - View links and verify count increments
   - Verify API responses match expected formats
3. Performance validation:
   - Compare search response times before/after changes
   - Verify no significant regression (< 10% slower)
4. API contract verification:
   - Confirm all existing API consumers would continue to work
   - Verify JSON response structures are preserved

**Test Commands**:
```bash
go test ./routes/... -v
go test ./... -race
go test ./... -bench=.
```

**Verification Criteria**:
- All existing tests pass
- No race conditions detected
- Performance within acceptable bounds
- API responses maintain backward compatibility

## Phase 4: Documentation and Cleanup

### Step 6: Code Cleanup and Documentation
**Duration**: 10 minutes
**Artifacts**: Clean code, updated comments

**Actions**:
1. Remove TODO comments about missing functionality in `links_view.go`
2. Add brief security-focused comments for parameter binding patterns
3. Ensure consistent error handling across all modified functions
4. Verify import statements are clean and minimal
5. Check for any remaining debug code or temporary changes

**Verification Criteria**:
- No TODO comments remain
- Code follows established patterns in the codebase
- Comments are helpful but not excessive
- All imports are used and necessary

### Step 7: Final Validation and Commit Preparation
**Duration**: 5 minutes
**Artifacts**: Clean git state, commit message

**Actions**:
1. Final test run: `go test ./routes/...`
2. Build verification: `go build`
3. Git status check - ensure only intended files modified
4. Prepare commit message following project conventions
5. Review all changes one final time

**Commit Message**:
```
fix: eliminate SQL injection vulnerabilities and enhance links view endpoint

- Replace fmt.Sprintf string interpolation with parameterized queries in links_search_simple.go
- Fix manual parameter replacement with proper dbx.Params binding in links_search.go
- Enhance links_view.go to return full updated record instead of just success message
- All changes maintain API compatibility while eliminating security vulnerabilities

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Verification Criteria**:
- All tests pass
- Application builds successfully
- Git diff shows only intended changes
- Ready for code review/deployment

## Risk Mitigation

### Rollback Strategy
- Each step creates a commit point for easy rollback
- Changes maintain API compatibility - no breaking changes
- Database schema unchanged - no migration rollback needed
- All existing functionality preserved

### Testing Strategy
- Unit tests for individual functions
- Integration tests for full endpoint workflows
- Security tests for injection prevention
- Performance tests for regression detection
- Manual testing with real payloads

### Monitoring Points
- Database query performance metrics
- API response time monitoring
- Error rate tracking for all three endpoints
- Security incident monitoring

## Success Metrics

### Security (Critical)
- Zero SQL injection vulnerabilities detected
- All user input properly parameterized
- Security scanner passes on all endpoints

### Functionality (High)
- Links view endpoint returns complete records
- Search functionality accuracy maintained
- API response formats preserved

### Performance (Medium)
- Query response times within 10% of baseline
- No memory leaks introduced
- Database connection handling efficient

### Code Quality (Medium)
- Code follows established PocketBase patterns
- Error handling consistent and appropriate
- Test coverage maintained or improved

This plan ensures systematic elimination of SQL injection vulnerabilities while enhancing functionality, with each step independently verifiable and safely reversible.