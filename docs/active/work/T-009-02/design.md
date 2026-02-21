# T-009-02 Design: SQL Injection and View Endpoint Fix

## Overview

This design addresses critical security vulnerabilities and functionality gaps in two Go endpoint handlers. The primary focus is eliminating SQL injection vulnerabilities while restoring proper functionality to the links view endpoint.

## Problem Analysis

### Critical Security Issues Identified

1. **Direct String Interpolation in links_search_simple.go (Line 82)**
   - Uses `fmt.Sprintf` to embed user input directly into SQL queries
   - Even with `escapeLikePattern2()` escaping, format string attacks remain possible
   - Attack vector: Format specifiers like `%s`, `%d` can manipulate the query structure

2. **Manual Parameter Replacement in links_search.go (Lines 210-224, 291-305, 336-349)**
   - Defeats parameterized queries by manually replacing `?` placeholders with string values
   - String replacement approach is fundamentally vulnerable to SQL injection
   - Current "escaping" using `strings.ReplaceAll(argStr, "'", "''")` is insufficient

3. **Missing Record Return in links_view.go**
   - Currently returns only `rows_affected` count after incrementing view_count
   - Should return the full updated link record per acceptance criteria
   - Contains TODO comment acknowledging the missing functionality

## Design Decisions

### 1. Security Fix Strategy: Proper Parameterized Queries

**Decision**: Replace all manual string interpolation with PocketBase's native parameter binding using `dbx.Params`

**Rationale**:
- PocketBase v0.36.5 provides robust parameter binding through `app.DB().NewQuery(sql).Bind(dbx.Params{...})`
- This approach is already demonstrated correctly in links_view.go:42 for the UPDATE query
- Eliminates all SQL injection vectors by letting the database driver handle parameter escaping
- Maintains query performance through prepared statement optimization

**Rejected Alternatives**:
- Enhanced manual escaping: Still vulnerable to edge cases and new attack vectors
- Input sanitization: Can never be comprehensive enough for security-critical code
- Query builders: Would require significant architectural changes beyond scope

### 2. Code Consolidation Strategy

**Decision**: Fix both search endpoints rather than removing the simple one

**Rationale**:
- Both endpoints serve different use cases (simple vs comprehensive search)
- API consumers may depend on both endpoint contracts
- Fixing is safer than removal without usage analysis
- Demonstrates proper patterns for future development

### 3. Links View Endpoint Enhancement

**Decision**: Use PocketBase's `app.FindRecordById()` pattern to return full record after update

**Rationale**:
- Pattern is well-established in the codebase (presentations.go, sync_sessions.go)
- Atomic operation: UPDATE followed by SELECT using database transaction semantics
- Maintains API contract expectations for record-returning endpoints
- Provides complete link data including computed fields and relationships

### 4. Response Format Preservation

**Decision**: Maintain existing JSON response structures for all endpoints

**Rationale**:
- Ensures backward compatibility with existing API consumers
- No breaking changes to public API contracts
- Focus on security fixes rather than API redesign

## Technical Implementation Approach

### SQL Injection Fix Pattern

Replace manual string building with parameterized queries:

**Before (Vulnerable)**:
```go
searchPattern := escapeLikePattern2(params.Q)
whereClauses = append(whereClauses, fmt.Sprintf("(title LIKE '%s' OR description LIKE '%s')", searchPattern, searchPattern))
```

**After (Secure)**:
```go
whereClauses = append(whereClauses, "(title LIKE {searchPattern} OR description LIKE {searchPattern})")
queryParams["searchPattern"] = "%" + escapedPattern + "%"
```

Then execute with: `db.NewQuery(query).Bind(dbx.Params(queryParams)).Rows()`

### Links View Enhancement Pattern

**Current Implementation**:
```go
result, err := app.DB().NewQuery(sql).Bind(dbx.Params{"linkId": linkId}).Execute()
return e.JSON(200, map[string]interface{}{
    "success": true,
    "rows_affected": rowsAffected,
})
```

**Enhanced Implementation**:
```go
// 1. Execute UPDATE (existing secure code)
result, err := app.DB().NewQuery(sql).Bind(dbx.Params{"linkId": linkId}).Execute()

// 2. Fetch updated record
record, err := app.FindRecordById("links", linkId)

// 3. Return full record in proper format
return e.JSON(200, LinkViewResponse{
    ID: record.Id,
    URL: record.GetString("url"),
    Title: record.GetString("title"),
    // ... other fields
})
```

### Parameter Binding Architecture

All dynamic SQL queries will follow this secure pattern:

1. **Query Structure**: Use named placeholders `{paramName}` in SQL strings
2. **Parameter Collection**: Build `map[string]interface{}` for all dynamic values
3. **Execution**: Use `dbx.Params(queryParams)` with `Bind()` method
4. **Type Safety**: Let PocketBase handle type conversion and escaping

This eliminates:
- Manual string concatenation/replacement
- Format string vulnerabilities
- SQL injection attack vectors
- Type conversion errors

## Database Interaction Patterns

### Search Queries (Both Endpoints)
- **Text Search**: `title LIKE {pattern} OR description LIKE {pattern}`
- **Tag Filtering**: `l.tags LIKE '%' || {tagId} || '%'` (JSON array search)
- **Pagination**: `LIMIT {limit} OFFSET {offset}`
- **Ordering**: Static string (no parameters needed)

### Count Queries
- Same WHERE conditions as main query
- Same parameter binding approach
- Single integer result

### Record Retrieval (Links View)
- UPDATE with parameter binding (already implemented correctly)
- Followed by `app.FindRecordById("links", linkId)`
- Convert PocketBase record to API response format

## Error Handling Strategy

### SQL Injection Prevention
- All user input treated as potentially malicious
- No direct string interpolation in any query
- Parameter binding for all dynamic values
- Type validation before database operations

### Database Error Handling
- Query execution errors return HTTP 500 with generic message
- Record not found errors return HTTP 404
- Parameter validation errors return HTTP 400
- No database details exposed to client

### Graceful Degradation
- Empty search results return empty arrays, not null
- Missing optional fields default to appropriate zero values
- Malformed parameters rejected with clear error messages

## Testing Strategy

### Security Testing
- SQL injection attempt testing with malicious payloads
- Format string attack attempts
- Special character handling verification
- Parameter boundary testing

### Functionality Testing
- Search result accuracy and pagination
- View count increment verification
- Record retrieval completeness
- API response format validation

### Regression Testing
- Existing test suite compatibility
- API contract preservation verification
- Performance impact assessment

## Migration Considerations

### Backward Compatibility
- All existing API endpoints maintain same URLs
- Response formats remain unchanged
- Query parameter interfaces preserved
- No breaking changes to public APIs

### Deployment Safety
- Changes are additive security fixes
- No database schema modifications required
- No configuration changes needed
- Zero-downtime deployment compatible

## Success Criteria

1. **Security**: Zero SQL injection vulnerabilities in all endpoints
2. **Functionality**: Links view endpoint returns complete updated records
3. **Compatibility**: All existing tests pass without modification
4. **Performance**: No significant query performance regression
5. **Maintainability**: Code follows established PocketBase patterns

This design eliminates all identified security vulnerabilities while preserving API compatibility and improving endpoint functionality through proper use of PocketBase v0.36.5 APIs.