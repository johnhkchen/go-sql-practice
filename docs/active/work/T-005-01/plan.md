# Plan: Go Integration Tests for Custom API Endpoints

## Implementation Sequence

The implementation will proceed through distinct phases, each building on the previous phase and providing immediate verification capabilities. Each step can be committed atomically and tested independently.

### Phase 1: Test Infrastructure Foundation

#### Step 1.1: Create Test File and Basic Setup
**Goal**: Establish test file structure and basic PocketBase test app setup

**Implementation**:
1. Create `routes/routes_test.go` with package declaration and imports
2. Implement `setupTestApp()` helper function for in-memory PocketBase creation
3. Implement `runTestMigrations()` helper to execute production migrations
4. Add basic cleanup functionality

**Validation**:
```bash
go test ./routes -run TestSetup -v
```

**Success Criteria**:
- Test file compiles without errors
- PocketBase app creates successfully with in-memory database
- Migrations execute without errors
- Cleanup prevents resource leaks

**Commit Message**: `test: add basic PocketBase test infrastructure`

#### Step 1.2: HTTP Request Testing Infrastructure
**Goal**: Implement HTTP request helpers and response parsing utilities

**Implementation**:
1. Implement `makeRequest()` helper for HTTP request execution
2. Implement `parseJSONResponse()` helper for response unmarshaling
3. Implement `assertErrorResponse()` helper for error validation
4. Add request/response logging for debugging

**Validation**:
```bash
go test ./routes -run TestHTTPHelpers -v
```

**Success Criteria**:
- HTTP requests execute against test app successfully
- JSON responses parse correctly
- Error responses validate consistently
- Debug logging works for troubleshooting

**Commit Message**: `test: add HTTP request and response helpers`

#### Step 1.3: Test Data Infrastructure
**Goal**: Create test data builders and seeding utilities

**Implementation**:
1. Implement `TestData` struct to hold test data references
2. Implement `createTestTag()` for individual tag creation
3. Implement `createTestLink()` for individual link creation
4. Implement basic data validation helpers

**Validation**:
```bash
go test ./routes -run TestDataCreation -v
```

**Success Criteria**:
- Tags create successfully with proper name/slug values
- Links create successfully with proper associations
- Test data helpers return valid IDs
- Database constraints work as expected

**Commit Message**: `test: add test data creation helpers`

### Phase 2: Search Endpoint Testing

#### Step 2.1: Basic Search Test Implementation
**Goal**: Implement core search functionality tests

**Implementation**:
1. Create `TestLinksSearch` function with app setup
2. Implement `"BasicTextSearch"` sub-test with text query validation
3. Implement `"EmptyResults"` sub-test with no-match scenarios
4. Add basic search response validation

**Test Data Needed**:
- 5 links with varying titles: "Go Programming Guide", "JavaScript Basics", "Database Design", "Testing Strategies", "API Development"
- No tags required for this step

**Validation**:
```bash
go test ./routes -run TestLinksSearch/BasicTextSearch -v
go test ./routes -run TestLinksSearch/EmptyResults -v
```

**Success Criteria**:
- Text search returns expected results
- Empty query handling works correctly
- Response format matches expected JSON structure
- Search result ordering is consistent

**Commit Message**: `test: add basic search endpoint tests`

#### Step 2.2: Tag Filtering and Combined Search Tests
**Goal**: Implement tag-based filtering and combined search scenarios

**Implementation**:
1. Create comprehensive test data with tags: golang, javascript, database
2. Implement `"TagFilter"` sub-test for single tag filtering
3. Implement `"CombinedFilters"` sub-test for text + tag combination
4. Add tag association validation

**Test Data Needed**:
- 8 links with controlled tag associations
- 3 tags with predictable link counts
- Mixed scenarios: links with multiple tags, links with no tags

**Validation**:
```bash
go test ./routes -run TestLinksSearch/TagFilter -v
go test ./routes -run TestLinksSearch/CombinedFilters -v
```

**Success Criteria**:
- Tag filtering returns only tagged links
- Combined filters work as logical AND operation
- Tag data populates correctly in response
- Edge cases (no tags, multiple tags) handle properly

**Commit Message**: `test: add search tag filtering and combined search tests`

#### Step 2.3: Pagination and Parameter Validation Tests
**Goal**: Complete search endpoint test coverage with pagination and error cases

**Implementation**:
1. Expand test data to 25+ links for pagination testing
2. Implement `"Pagination"` sub-test with multiple pages
3. Implement `"InvalidPagination"` sub-test with bad parameters
4. Implement `"SQLInjectionProtection"` sub-test with malicious input

**Test Data Needed**:
- 25 links with predictable ordering
- Variety of titles for SQL injection testing

**Validation**:
```bash
go test ./routes -run TestLinksSearch/Pagination -v
go test ./routes -run TestLinksSearch/InvalidPagination -v
go test ./routes -run TestLinksSearch/SQLInjectionProtection -v
```

**Success Criteria**:
- Pagination returns correct page counts and offsets
- Invalid parameters return appropriate 400 errors
- SQL injection attempts are safely handled
- All search tests pass consistently

**Commit Message**: `test: add search pagination and parameter validation tests`

### Phase 3: View Count Endpoint Testing

#### Step 3.1: Basic View Count Tests
**Goal**: Implement core view count increment functionality tests

**Implementation**:
1. Create `TestLinksView` function with focused test data
2. Implement `"SuccessfulIncrement"` sub-test with count verification
3. Implement `"NonexistentLink"` sub-test with 404 validation
4. Add response format validation

**Test Data Needed**:
- 3 links with known initial view counts: 0, 5, null
- 1 known non-existent link ID

**Validation**:
```bash
go test ./routes -run TestLinksView/SuccessfulIncrement -v
go test ./routes -run TestLinksView/NonexistentLink -v
```

**Success Criteria**:
- View counts increment by exactly 1
- Non-existent links return 404 status
- Response includes updated link record
- Database changes persist correctly

**Commit Message**: `test: add basic view count endpoint tests`

#### Step 3.2: Edge Cases and Atomic Updates
**Goal**: Test edge cases and concurrency scenarios

**Implementation**:
1. Implement `"ZeroToOne"` sub-test for null view_count initialization
2. Implement `"ResponseFormat"` sub-test for PocketBase record structure
3. Implement `"AtomicUpdates"` sub-test with goroutine concurrency
4. Add concurrent test synchronization

**Test Data Needed**:
- Links with null view_count values
- Single link for concurrency testing

**Validation**:
```bash
go test ./routes -run TestLinksView/ZeroToOne -v
go test ./routes -run TestLinksView/AtomicUpdates -v
```

**Success Criteria**:
- Null view_count initializes to 1 correctly
- Concurrent increments maintain data integrity
- Final counts match number of increments
- No race conditions or data corruption

**Commit Message**: `test: add view count edge cases and concurrency tests`

### Phase 4: Stats Endpoint Testing

#### Step 4.1: Basic Stats Tests
**Goal**: Implement core stats endpoint functionality tests

**Implementation**:
1. Create `TestStats` function with predictable test data
2. Implement `"CompleteResponse"` sub-test with all fields validation
3. Implement `"DataAccuracy"` sub-test with precise count verification
4. Add stats response structure validation

**Test Data Needed**:
- Exactly 5 links with view counts: [45, 23, 12, 8, 0]
- Exactly 3 tags with link associations: golang(3), javascript(2), database(1)
- Predictable totals for validation

**Validation**:
```bash
go test ./routes -run TestStats/CompleteResponse -v
go test ./routes -run TestStats/DataAccuracy -v
```

**Success Criteria**:
- All stats fields populate correctly
- Counts match test data exactly
- Response format matches expected JSON structure
- No missing or null fields

**Commit Message**: `test: add basic stats endpoint tests`

#### Step 4.2: Ordering and Edge Cases
**Goal**: Complete stats endpoint test coverage with ordering and empty data scenarios

**Implementation**:
1. Implement `"TopItems"` sub-test with ordering verification
2. Implement `"EmptyDatabase"` sub-test with zero data handling
3. Implement `"ResponseSchema"` sub-test with complete JSON validation
4. Add comprehensive stats validation helpers

**Test Data Needed**:
- Empty database scenario (separate test app)
- Controlled ordering data for top items verification

**Validation**:
```bash
go test ./routes -run TestStats/TopItems -v
go test ./routes -run TestStats/EmptyDatabase -v
```

**Success Criteria**:
- Most viewed links ordered correctly (descending view count)
- Top tags ordered correctly (descending link count)
- Empty database returns zeros gracefully
- All stats tests pass consistently

**Commit Message**: `test: add stats ordering and edge case tests`

### Phase 5: Integration and Final Validation

#### Step 5.1: Comprehensive Test Data Integration
**Goal**: Implement `seedTestData()` for comprehensive cross-endpoint testing

**Implementation**:
1. Implement `seedTestData()` with realistic data relationships
2. Create comprehensive test scenario covering all endpoints
3. Add cross-endpoint data consistency validation
4. Optimize test data for performance

**Test Data Specification**:
- 15 links with varied content and view counts
- 5 tags with realistic associations
- Sufficient complexity for all test scenarios
- Predictable relationships for validation

**Validation**:
```bash
go test ./routes -v
```

**Success Criteria**:
- All endpoints work with shared comprehensive test data
- Cross-endpoint data consistency maintained
- Performance remains acceptable (<5 seconds full suite)
- No test flakiness or race conditions

**Commit Message**: `test: add comprehensive test data integration`

#### Step 5.2: Debug Helpers and Performance Optimization
**Goal**: Add debugging utilities and optimize test performance

**Implementation**:
1. Implement `dumpTestData()` for debugging failed tests
2. Implement `logRequest()` for HTTP debugging
3. Add performance monitoring and optimization
4. Create test documentation and examples

**Validation**:
```bash
go test ./routes -v -cover
go test ./routes -bench=.
```

**Success Criteria**:
- Debug helpers aid in troubleshooting
- Test coverage meets acceptance criteria (>90%)
- Performance benchmarks within acceptable ranges
- Documentation supports future maintenance

**Commit Message**: `test: add debug helpers and optimize performance`

### Phase 6: Final Integration Testing

#### Step 6.1: End-to-End Validation
**Goal**: Validate complete test suite meets all acceptance criteria

**Implementation**:
1. Run complete test suite with coverage reporting
2. Validate all acceptance criteria are met
3. Test integration with existing build system
4. Verify CI/CD compatibility

**Validation Commands**:
```bash
go test ./...                    # Full project test suite
go test ./routes -cover -v       # Routes with coverage
make test                        # Makefile integration
```

**Acceptance Criteria Validation**:
- ✅ Test file(s) in routes package
- ✅ Test helper creates in-memory PocketBase instance
- ✅ Search endpoint: basic query, tag filter, combined filters, empty results, pagination
- ✅ View count endpoint: successful increment, 404 for nonexistent ID, count increases
- ✅ Stats endpoint: expected shape, counts match seed data
- ✅ `go test ./...` passes with all tests green
- ✅ Tests do not require running PocketBase server

**Success Criteria**:
- All acceptance criteria verified
- Test suite runs consistently in CI/CD
- No external dependencies required
- Documentation complete and accurate

**Commit Message**: `test: complete integration tests for custom API endpoints`

## Testing Strategy

### Unit Test Coverage
- **Search endpoint**: 7 test scenarios covering all query parameters and error cases
- **View count endpoint**: 5 test scenarios covering increments, errors, and concurrency
- **Stats endpoint**: 5 test scenarios covering response completeness and edge cases
- **Total**: 17 comprehensive test scenarios

### Integration Points
- PocketBase app lifecycle management
- SQLite in-memory database operations
- HTTP request/response handling
- JSON serialization/deserialization
- Database migration system
- Route registration system

### Error Scenarios
- Invalid input parameters (search pagination, view count IDs)
- Database constraint violations
- Concurrent access patterns
- Empty/null data handling
- Malicious input protection

### Performance Requirements
- Individual tests: <100ms execution
- Full test suite: <5 seconds execution
- Memory usage: <50MB peak per test
- No resource leaks or cleanup issues

### Maintenance Considerations
- Test data builders centralize data creation logic
- Shared helpers reduce code duplication
- Debug utilities support troubleshooting
- Clear test naming supports future modifications

This plan provides a systematic approach to implementing comprehensive integration tests while maintaining incremental progress, clear validation steps, and atomic commit points.