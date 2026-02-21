# T-009-02 Research: SQL Injection and View Endpoint Fix

## Overview

This ticket addresses two critical security and functionality issues in the Go endpoint handlers:
1. SQL injection vulnerability in `links_search_simple.go` (line 82)
2. Missing functionality to return updated records after view count increment in `links_view.go`

## Current Codebase Analysis

### Project Structure
- **Framework**: PocketBase v0.36.5 with Go 1.26
- **Database**: SQLite with PocketBase ORM
- **Architecture**: Routes registered through `routes.Register()` in main.go
- **Collections**: Links, tags, sync_sessions, presentations (defined in migrations/collections.go)

### SQL Injection Issues Identified

#### 1. routes/links_search_simple.go (Line 82)
**Vulnerability**: Direct string interpolation using `fmt.Sprintf`
```go
searchPattern := escapeLikePattern2(params.Q)
whereClauses = append(whereClauses, fmt.Sprintf("(title LIKE '%s' OR description LIKE '%s')", searchPattern, searchPattern))
```

**Problem**: Even though `escapeLikePattern2()` escapes LIKE wildcards and single quotes, the use of `fmt.Sprintf` allows SQL injection through format string attacks. An attacker can input format specifiers like `%s`, `%d`, etc. to manipulate the query.

#### 2. routes/links_search.go (Lines 210-224, 291-305, 336-349)
**Vulnerability**: Manual string replacement of parameterized queries
```go
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
```

**Problem**: This approach defeats the purpose of parameterized queries. The manual string replacement is vulnerable to injection attacks and should use proper parameter binding.

### Links View Endpoint Issue

#### routes/links_view.go Analysis
**Current Implementation**:
- Uses `app.DB().NewQuery(sql).Bind(dbx.Params{"linkId": linkId}).Execute()` - CORRECT approach
- Returns only `rows_affected` count, not the updated record
- Contains TODO comment about "daos package doesn't exist in v0.36.5"

**Missing Functionality**: Should return the full updated link record after incrementing view_count

### Database Schema (from migrations/collections.go)
Links collection has fields:
- `id` (auto-generated)
- `url` (URLField, required)
- `title` (TextField, required, max 500 chars)
- `description` (TextField, optional, max 2000 chars)
- `view_count` (NumberField, optional, integer, min 0)
- `tags` (RelationField to tags collection)
- `created_by` (RelationField to users, optional)
- Standard PocketBase fields: `created`, `updated`

### PocketBase v0.36.5 API Patterns Found

#### Correct Record Operations:
1. **Find by ID**: `app.FindRecordById("collection", "id")`
   - Used in presentations.go:193, 235, 293
   - Used in sync_sessions.go:114
2. **Save Record**: `app.Save(record)`
   - Used in presentations.go:261, 337, 347
   - Used in sync_sessions.go:69, 133
3. **Proper Parameterized Queries**:
   - `app.DB().NewQuery(sql).Bind(dbx.Params{...}).Execute()`
   - Used correctly in links_view.go:42

#### Query Execution Patterns:
- **Single Value**: `db.NewQuery(query).Row(&variable)`
- **Multiple Rows**: `db.NewQuery(query).Rows()` then iterate with `rows.Next()`
- **Parameterized**: Use `Bind(dbx.Params{...})` NOT manual string interpolation

### Route Registration
Both routes are registered in routes/routes.go:
- `registerLinksSearch(e)` - line 16 (regular search)
- `registerLinksSearchSimple(e)` - line 17 (simple search)
- `registerLinksView(e)` - line 18

### Functional Overlap Analysis
**links_search.go vs links_search_simple.go**:
- Both provide search functionality with query params (q, tag, page, perPage)
- links_search.go is more comprehensive (tag filtering, better pagination)
- links_search_simple.go is simpler but has the same core purpose
- **Decision Factor**: links_search.go also has SQL injection issues, so both need fixing

### Security Context
Current applications using these endpoints are vulnerable to:
1. **SQL Injection** via search parameters in both search endpoints
2. **Data Leakage** potential through injection attacks
3. **DoS attacks** through malicious SQL queries

### Dependencies and Constraints
- Must use PocketBase v0.36.5 API (no daos package)
- Must maintain existing API contract (JSON response shapes)
- Must preserve public read access (per collection rules)
- Links collection accessible without authentication
- Database is SQLite (affects SQL syntax choices)

### Testing Infrastructure
- Test files exist: routes/links_search_test.go, routes/routes_test.go
- Uses testPocketBase setup with in-memory database
- Has helper functions for test data creation

## Key Findings

1. **Both search endpoints are vulnerable** - not just the simple one
2. **PocketBase provides proper parameterized query support** - should use `Bind(dbx.Params{})`
3. **Record retrieval API exists** - `app.FindRecordById()` pattern well established
4. **Simple endpoint may be redundant** - links_search.go has more features
5. **Existing patterns show correct approach** - links_view.go uses proper binding for UPDATE
6. **Database schema supports all required fields** for proper response formatting