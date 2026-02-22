# Structure: Frontend Shared Types and API Client
## Ticket T-010-03

### New File Structure

#### 1. Create: `frontend/src/types/api.ts`

**Purpose**: Single source of truth for all API-related TypeScript interfaces

**Exports**:
```typescript
// Link-related types
export interface LinkItem
export interface SearchResponse
export interface PocketBaseResponse<T = LinkItem>

// Stats types
export interface StatsData
export interface StatsState

// Presentation types
export interface Presentation
export interface PresentationStatus

// Search state
export interface SearchState
```

**Module Boundaries**:
- Pure type definitions only (no runtime code)
- All interfaces exported for maximum reusability
- Generic PocketBaseResponse for flexibility
- No imports from other modules (self-contained)

#### 2. Create: `frontend/src/lib/api.ts`

**Purpose**: Shared API utilities and configuration

**Exports**:
```typescript
// Constants
export const API_BASE: string
export const FETCH_TIMEOUT: number
export const FETCH_TIMEOUT_LONG: number

// Error class
export class ApiError extends Error

// Fetch helper
export async function apiFetch<T>(): Promise<T>
```

**Module Boundaries**:
- Runtime utilities only (executable code)
- No component-specific logic
- Environment variable access centralized here
- Error handling standardized

### Modified Files Structure

#### 3. Modify: `frontend/src/pages/index.astro`

**Changes**:
- Remove lines 15-16 (API_BASE, FETCH_TIMEOUT constants)
- Remove lines 19-33 (interface definitions)
- Add import: `import { API_BASE, FETCH_TIMEOUT, apiFetch } from '../lib/api'`
- Add import: `import type { LinkItem, SearchResponse, PocketBaseResponse } from '../types/api'`
- Keep LinkCard import but remove type import from it
- Update fetch pattern to use apiFetch (optional first pass)

#### 4. Modify: `frontend/src/pages/tags/[slug].astro`

**Changes**:
- Remove lines 16-17 (API_BASE, FETCH_TIMEOUT)
- Remove lines 20-36 (LinkItem, SearchResponse interfaces)
- Add imports for types and api utilities
- Update fetch pattern

#### 5. Modify: `frontend/src/pages/links/[id].astro`

**Changes**:
- Remove lines 15-16 (constants)
- Remove lines 19-28 (LinkItem interface)
- Add imports
- Update field references: `created` → `created_at`
- Remove line 311 (duplicate API_BASE in script block)

#### 6. Modify: `frontend/src/pages/present/index.astro`

**Changes**:
- Remove lines 9-10 (constants)
- Remove lines 13-40 (Presentation, PresentationStatus, PocketBaseResponse)
- Add imports
- Generic PocketBaseResponse<Presentation> for type safety

#### 7. Modify: `frontend/src/pages/present/[id].astro`

**Changes**:
- Remove lines 20-21 (constants)
- Add imports including FETCH_TIMEOUT_LONG for longer operations
- Keep local timeout override pattern if needed

#### 8. Modify: `frontend/src/pages/sync/[id].astro`

**Changes**:
- Remove lines 12-13 (constants)
- Add imports

#### 9. Modify: `frontend/src/pages/sync/[id]/control.astro`

**Critical Fix**:
- Move line 28 (`let error = null`) to line 18 (before first use)
- Remove lines 23-24 (constants after fix)
- Add imports

#### 10. Modify: `frontend/src/pages/watch/[id].astro`

**Changes**:
- Remove lines 16-17 (constants)
- Add imports

#### 11. Modify: `frontend/src/pages/stats.astro`

**Changes**:
- Add imports if API_BASE is used
- Import StatsData, StatsState types

#### 12. Modify: `frontend/src/components/SearchInterface.astro`

**Changes in script block**:
- Remove lines 433-457 (interface definitions)
- Add: `import type { LinkItem, SearchResponse, SearchState } from '../types/api'`
- Interfaces now imported as types for client-side code

#### 13. Modify: `frontend/src/components/StatsSummary.astro`

**Changes**:
- Remove lines 3-24 (interface exports)
- Add: `export type { StatsData, StatsState } from '../types/api'`
- Re-export pattern maintains component API

#### 14. Modify: `frontend/src/components/LinkCard.astro`

**Changes**:
- Remove LinkItem export if it exists
- Import from types/api if needed internally
- Update any components importing LinkItem from here

### Import Graph

```
types/api.ts (no imports)
    ↑
    ├── lib/api.ts (imports types)
    ├── pages/*.astro (import both)
    ├── components/SearchInterface.astro (import types)
    ├── components/StatsSummary.astro (re-export types)
    └── components/LinkCard.astro (import types if needed)
```

### Module Boundaries

#### Types Module (`types/api.ts`)
- **Responsibility**: Type definitions only
- **Dependencies**: None
- **Consumers**: All pages and components
- **Contract**: Stable interface definitions

#### API Module (`lib/api.ts`)
- **Responsibility**: API utilities and configuration
- **Dependencies**: types/api (for error types)
- **Consumers**: All pages making API calls
- **Contract**: Consistent fetch behavior

#### Pages
- **Responsibility**: SSR rendering and data fetching
- **Dependencies**: types/api, lib/api
- **Contract**: Use shared utilities, no duplicate definitions

#### Components
- **Responsibility**: Reusable UI logic
- **Dependencies**: types/api (no direct api.ts dependency)
- **Contract**: Type-safe props, optional re-exports

### File Creation Order

1. Create `frontend/src/types/api.ts` first (no dependencies)
2. Create `frontend/src/lib/api.ts` second (depends on types)
3. Update pages in any order (all independent)
4. Update components last (may affect pages)

### Validation Points

After creating new files:
- TypeScript compilation should succeed
- No runtime errors (types are compile-time only)

After updating each page:
- Page should still render
- API calls should succeed
- No TypeScript errors

After updating components:
- Components should render
- Parent pages should not break
- Type checking should pass

### Architecture Decisions

**Why separate types/ and lib/?**
- types/ is pure TypeScript (no runtime)
- lib/ contains executable code
- Clear separation of concerns
- Better tree-shaking potential

**Why not use frontend/src/utils/?**
- utils/ is generic utilities
- lib/ is specifically for library-like modules
- api.ts is more library than utility
- Follows common convention

**Why re-export from StatsSummary?**
- Maintains backwards compatibility
- Component already exports these types
- Smooth migration path
- Can be removed in future refactor

### Breaking Changes

None expected if implementation follows this structure. All changes maintain existing public APIs while reorganizing internals.