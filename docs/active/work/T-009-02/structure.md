# T-009-02 Structure: SQL Injection and View Endpoint Fix

## Overview

This structure defines the specific file-level changes needed to eliminate SQL injection vulnerabilities and restore full functionality to the links view endpoint. All changes maintain API compatibility while implementing secure parameterized queries.

## File Modifications

### 1. routes/links_search_simple.go (MODIFIED)

**Purpose**: Eliminate SQL injection vulnerability through proper parameterized queries

**Key Changes**:
- Replace `fmt.Sprintf` string interpolation with `dbx.Params` binding
- Modify `handleSearchSimple()` function to use secure query patterns
- Update `escapeLikePattern2()` to return unescaped pattern (escaping handled by binding)
- Preserve all existing API contracts and response formats

**Modified Functions**:
```go
// handleSearchSimple() - lines 52-138
- Remove: fmt.Sprintf string building on line 82
- Add: dbx.Params parameter collection
- Change: Query execution to use .Bind() method

// escapeLikePattern2() - lines 178-187
- Remove: Single quote escaping (handled by parameter binding)
- Keep: LIKE wildcard escaping (\%, \_, \\)
- Change: Return unescaped pattern for parameter binding
```

**New Imports**:
- Add: `"github.com/pocketbase/dbx"` (for dbx.Params)

**Internal Architecture**:
- Query building: Collect WHERE conditions with named placeholders
- Parameter collection: Build `map[string]interface{}` for all dynamic values
- Execution: `db.NewQuery(query).Bind(dbx.Params(queryParams)).Rows()`

### 2. routes/links_search.go (MODIFIED)

**Purpose**: Eliminate manual parameter replacement with proper parameterized queries

**Key Changes**:
- Remove all manual `strings.Replace(finalQuery, "?", argStr, 1)` operations
- Replace with `dbx.Params` binding throughout
- Maintain existing query logic and response formats
- Fix three vulnerable functions: `searchLinks()`, `countSearchResults()`, `fetchTagsForLinks()`

**Modified Functions**:
```go
// searchLinks() - lines ~200-230
- Remove: Manual parameter replacement loop
- Add: Named placeholder system {param1}, {param2}, etc.
- Change: Query execution to use parameter binding

// countSearchResults() - lines ~280-310
- Remove: Manual parameter replacement loop
- Add: Same parameter binding pattern as searchLinks()
- Maintain: Same WHERE clause building logic

// fetchTagsForLinks() - lines ~330-360
- Remove: Manual parameter replacement for IN clause
- Add: Named placeholders for dynamic IN clause construction
- Change: Build parameterized IN clause with {param1}, {param2}, ...
```

**Parameter Binding Strategy**:
- Convert `args []interface{}` to `queryParams map[string]interface{}`
- Use sequential parameter names: `param1`, `param2`, etc.
- Replace `?` placeholders with `{param1}`, `{param2}` in query string
- Execute with: `db.NewQuery(query).Bind(dbx.Params(queryParams))`

### 3. routes/links_view.go (MODIFIED)

**Purpose**: Return full updated link record instead of just rows_affected count

**Key Changes**:
- Add record retrieval after successful UPDATE operation
- Create proper response structure for full link data
- Maintain existing UPDATE query (already secure with parameterized binding)
- Add error handling for record retrieval operations

**Modified Functions**:
```go
// handleLinksView() - lines 19-63
- Keep: Existing UPDATE query (lines 41-42) - already secure
- Add: Record retrieval using app.FindRecordById()
- Add: Response structure conversion from PocketBase record
- Change: Return full link record instead of rows_affected

// Response format change:
- Remove: map[string]interface{}{"success": true, "rows_affected": count}
- Add: Complete link record with all fields
```

**New Response Structure**:
```go
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
```

**Record Retrieval Pattern**:
```go
// After successful UPDATE
record, err := app.FindRecordById("links", linkId)
if err != nil {
    // Handle record not found or database error
}

// Convert to response format
response := LinkViewResponse{
    ID:          record.Id,
    URL:         record.GetString("url"),
    Title:       record.GetString("title"),
    Description: record.GetString("description"),
    ViewCount:   record.GetInt("view_count"),
    Created:     record.Created.Time().Format(time.RFC3339),
    Updated:     record.Updated.Time().Format(time.RFC3339),
}

// Fetch and attach tags separately (if needed)
```

## Architecture Patterns

### Parameterized Query Pattern (Applied to All Files)

**Standard Structure**:
```go
// 1. Build query with named placeholders
query := "SELECT * FROM table WHERE field = {param1} AND other = {param2}"

// 2. Collect parameters
queryParams := map[string]interface{}{
    "param1": userInput1,
    "param2": userInput2,
}

// 3. Execute with binding
rows, err := db.NewQuery(query).Bind(dbx.Params(queryParams)).Rows()
```

**Benefits**:
- Complete SQL injection prevention
- Type-safe parameter handling
- Database driver optimization
- Clear separation of query structure and data

### Error Handling Pattern (Consistent Across Files)

**Database Operation Errors**:
```go
if err != nil {
    return e.JSON(http.StatusInternalServerError, map[string]string{
        "error": "Database operation failed",
    })
}
```

**Record Not Found Errors**:
```go
if record == nil {
    return e.JSON(http.StatusNotFound, map[string]string{
        "error": "Record not found",
    })
}
```

### Response Format Preservation

**Search Endpoints** (No Changes):
- Maintain existing JSON structure
- Keep pagination metadata
- Preserve field naming conventions
- Ensure empty results return `[]` not `null`

**View Endpoint** (Enhanced):
- Return complete link record instead of success message
- Include all standard link fields
- Format timestamps consistently
- Handle missing optional fields gracefully

## Module Boundaries

### Security Boundary
- **Input Validation**: Parameter validation at handler entry points
- **Query Construction**: Named placeholder system with parameter maps
- **Execution**: Parameterized binding only - no string interpolation
- **Output Sanitization**: Standard JSON marshalling (no raw SQL data)

### Data Access Boundary
- **PocketBase Integration**: Use established patterns from existing endpoints
- **Record Operations**: Follow `app.FindRecordById()` and `app.Save()` patterns
- **Transaction Handling**: Rely on PocketBase's built-in transaction management
- **Connection Management**: Use provided `app.DB()` instance

### API Contract Boundary
- **URL Routes**: No changes to existing endpoint paths
- **HTTP Methods**: Preserve existing method requirements
- **Request Parameters**: Maintain existing query parameter interfaces
- **Response Formats**: Keep JSON structure compatibility (with enhancement for view endpoint)

## Change Ordering Requirements

### 1. Independent Changes (Can be done in parallel)
- `links_search_simple.go` fixes
- `links_search.go` fixes
- `links_view.go` enhancements

### 2. Testing Verification Order
1. Unit test security fixes (SQL injection prevention)
2. Integration test API compatibility
3. End-to-end test complete user workflows
4. Performance test query execution speed

### 3. Deployment Considerations
- All changes maintain backward compatibility
- No database migrations required
- No configuration changes needed
- Changes can be deployed incrementally if needed

## File Dependencies

### Import Dependencies
- **Existing**: `github.com/pocketbase/pocketbase/core`
- **Added**: `github.com/pocketbase/dbx` (for all modified files)
- **Removed**: None

### Code Dependencies
- **links_search_simple.go**: Independent (no dependencies on other modified files)
- **links_search.go**: Independent (no dependencies on other modified files)
- **links_view.go**: Independent (no dependencies on other modified files)

### Testing Dependencies
- Existing test files: `routes/links_search_test.go`, `routes/routes_test.go`
- Test modifications required: Update to test parameterized query security
- New test cases: SQL injection attempt validation

## Success Metrics

### Security Verification
- Zero SQL injection vulnerabilities detected by security scanner
- All user input properly parameterized in database queries
- No direct string interpolation in any SQL query

### Functionality Verification
- All existing API tests pass without modification
- Links view endpoint returns complete record data
- Search endpoints maintain response format compatibility
- Database query performance within acceptable ranges (< 10% regression)

### Code Quality Verification
- All modified functions follow established PocketBase patterns
- Error handling consistent across all endpoints
- Parameter binding used consistently for all dynamic queries
- No TODO comments or temporary workarounds remaining

This structure provides a comprehensive blueprint for eliminating SQL injection vulnerabilities while enhancing the links view endpoint functionality, all while maintaining strict API compatibility.