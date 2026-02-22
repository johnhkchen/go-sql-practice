# Design: Frontend Shared Types and API Client
## Ticket T-010-03

### Design Goals
1. Eliminate all duplicate type definitions and API configuration
2. Create a single source of truth for types and API utilities
3. Maintain TypeScript type safety across SSR and client code
4. Fix the control.astro variable declaration bug
5. Minimize breaking changes to existing functionality

### Design Options

#### Option 1: Minimal Extraction
Create shared modules with exact copies of existing patterns.

**Pros**:
- Zero risk of breaking changes
- Simple find-and-replace refactoring
- Maintains all existing behavior exactly

**Cons**:
- Doesn't fix inconsistencies (created vs created_at)
- Perpetuates suboptimal patterns
- Misses opportunity to improve error handling

#### Option 2: Normalized Types with Adapters
Create canonical types and adapter functions to handle variations.

**Pros**:
- Single consistent data model
- Type-safe transformations
- Future-proof design

**Cons**:
- Requires adapter logic for each variation
- More complex implementation
- Higher risk of runtime errors

#### Option 3: Progressive Enhancement (Selected) ✓
Create shared modules with the most complete/correct versions, update all consumers to match.

**Pros**:
- Fixes inconsistencies during refactor
- Single source of truth
- Simpler than adapters
- Improves overall code quality

**Cons**:
- Requires careful testing
- Must verify field mappings

### Design Decisions

#### 1. Type Definitions Structure (`types/api.ts`)

**LinkItem Interface**:
Use the most complete definition from research:
```typescript
export interface LinkItem {
  id: string;
  title: string;
  url: string;
  description: string;
  tags: string[];
  created_at: string;  // Standardize on created_at
  view_count: number;
  updated?: string;    // Optional for backwards compatibility
}
```

**Rationale**:
- `created_at` matches the search API response format
- `updated` is optional since not all endpoints return it
- This matches the majority usage pattern

**Other Interfaces**:
Export all interfaces to support both SSR and client usage:
```typescript
export interface SearchResponse { ... }
export interface PocketBaseResponse<T = LinkItem> { ... }
export interface StatsData { ... }
export interface StatsState { ... }
export interface Presentation { ... }
export interface PresentationStatus { ... }
export interface SearchState { ... }
```

Generic PocketBaseResponse allows reuse for different item types.

#### 2. API Client Design (`lib/api.ts`)

**Constants**:
```typescript
export const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
export const FETCH_TIMEOUT = 5000;
export const FETCH_TIMEOUT_LONG = 10000; // For complex operations
```

**Generic Fetch Helper**:
```typescript
export async function apiFetch<T>(
  url: string,
  options?: RequestInit & { timeout?: number }
): Promise<T>
```

**Design Features**:
- Generic type parameter for response type safety
- Optional custom timeout (defaults to FETCH_TIMEOUT)
- Automatic AbortController setup
- JSON parsing with type assertion
- Structured error throwing

**Error Handling**:
```typescript
export class ApiError extends Error {
  constructor(
    message: string,
    public code: 'timeout' | 'network' | 'server' | 'notfound',
    public status?: number
  ) {
    super(message);
  }
}
```

Structured errors allow consistent handling across pages.

#### 3. Migration Strategy

**Phase 1**: Create new files without breaking existing code
- Add `frontend/src/types/api.ts`
- Add `frontend/src/lib/api.ts`
- These can coexist with inline definitions initially

**Phase 2**: Update pages to use shared modules
- Import types instead of inline definitions
- Replace fetch patterns with apiFetch
- Fix control.astro bug during refactor

**Phase 3**: Clean up components
- Remove duplicate definitions from SearchInterface script
- Update StatsSummary to re-export from types/api

#### 4. Import Resolution

**Astro Compatibility**:
- Use relative imports in SSR sections: `import { LinkItem } from '../types/api'`
- Script blocks can use: `import type { LinkItem } from '../types/api'`
- Ensure TypeScript `moduleResolution` supports our structure

**Build Compatibility**:
- Types are compile-time only, no runtime impact
- API utilities are standard ES modules
- Compatible with Astro's SSR and client hydration

#### 5. Backwards Compatibility

**Field Mapping**:
- Pages using `created` will need updates to `created_at`
- This aligns with the API response format
- SearchInterface already uses `created_at`

**API Response Handling**:
- If PocketBase returns `created`, we map it in apiFetch
- Ensures consistent internal representation
- Single place to handle API variations

#### 6. control.astro Bug Fix

**Current Bug**:
```typescript
if (!token) {
  error = 'missing_token'; // Line 19 - uses before declaration
}
// ...
let error = null; // Line 28 - declaration
```

**Fix**:
Move declaration before first use:
```typescript
let error = null;
if (!token) {
  error = 'missing_token';
}
```

### Implementation Risks & Mitigations

**Risk 1**: Type mismatches break runtime
- **Mitigation**: Comprehensive testing of each endpoint
- **Mitigation**: Use optional fields for variations

**Risk 2**: Import path issues
- **Mitigation**: Test build after each file update
- **Mitigation**: Use consistent relative paths

**Risk 3**: Client-side script compatibility
- **Mitigation**: Ensure types are importable in script blocks
- **Mitigation**: Test progressive enhancement features

### Rejected Alternatives

**Monolithic API Class**:
- Too heavy for SSR pages
- Unnecessary object instantiation
- Goes against Astro's lightweight philosophy

**Per-Page Type Files**:
- Doesn't solve duplication
- More files to maintain
- No single source of truth

**GraphQL/tRPC Integration**:
- Overkill for current needs
- Requires backend changes
- Outside ticket scope

### Success Metrics
1. No duplicate interface definitions (grep returns single results)
2. All pages use shared API_BASE and FETCH_TIMEOUT
3. TypeScript build succeeds with no errors
4. control.astro variable bug is fixed
5. Frontend functionality unchanged

### Summary
The progressive enhancement approach provides the best balance of improvement and safety. By consolidating on the most complete type definitions and creating a reusable fetch helper, we eliminate duplication while improving code quality. The phased migration allows validation at each step.