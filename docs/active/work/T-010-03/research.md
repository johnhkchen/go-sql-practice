# Research: Frontend Shared Types and API Client
## Ticket T-010-03

### Overview
This ticket addresses systematic duplication across the Astro frontend, where every SSR page independently defines API configuration, fetch patterns, and TypeScript interfaces. The frontend currently uses copy-paste patterns across 9+ pages, making maintenance error-prone and inflating file sizes.

### Current Frontend Structure

#### Directory Layout
```
frontend/src/
├── components/      # Astro components
├── layouts/         # Layout components
├── pages/          # SSR pages
├── styles/         # CSS styles
└── utils/          # Utility functions
```

**Note**: No `lib/` or `types/` directories currently exist. These will need to be created.

### Duplication Analysis

#### 1. API Configuration Constants

**Pattern Found**: Every SSR page defines identical constants:
```typescript
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
const FETCH_TIMEOUT = 5000;
```

**Occurrences**:
- `pages/index.astro:15-16` - FETCH_TIMEOUT = 5000
- `pages/tags/[slug].astro:16-17` - FETCH_TIMEOUT = 5000
- `pages/links/[id].astro:15-16` - FETCH_TIMEOUT = 5000
- `pages/present/index.astro:9-10` - FETCH_TIMEOUT = 5000
- `pages/present/[id].astro:20-21` - FETCH_TIMEOUT = 10000 (longer timeout)
- `pages/sync/[id].astro:12-13` - FETCH_TIMEOUT = 5000
- `pages/sync/[id]/control.astro:23-24` - FETCH_TIMEOUT = 5000
- `pages/watch/[id].astro:16-17` - FETCH_TIMEOUT = 5000
- `pages/stats.astro` - Needs verification
- Client-side in `pages/links/[id].astro:311` duplicates API_BASE

#### 2. AbortController Timeout Pattern

**Pattern Found**: Identical fetch timeout implementation across all pages:
```typescript
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT);
// ... fetch ...
clearTimeout(timeoutId);
```

This pattern appears in every page that fetches data.

#### 3. Interface Definitions

**LinkItem Interface**:
Multiple incompatible definitions exist:

- `pages/index.astro:8` - Imports from LinkCard.astro
- `pages/tags/[slug].astro:20-28` - Inline definition with `created_at`
- `pages/links/[id].astro:19-28` - Inline definition with `created`, `updated`
- `components/SearchInterface.astro:433-441` - Script block definition with `created_at`
- `components/LinkCard.astro` - Original export (needs verification)

**Field Inconsistencies**:
- Creation date: `created_at` vs `created`
- Update tracking: Some have `updated`, others don't
- Different field orders

**SearchResponse Interface**:
- `pages/index.astro:27-33` - Inline definition
- `pages/tags/[slug].astro:30-36` - Duplicate inline definition
- `components/SearchInterface.astro:443-449` - Script block duplicate

**PocketBaseResponse Interface**:
- `pages/index.astro:19-25` - Generic version
- `pages/present/index.astro:34-40` - Presentation-specific

**Stats Interfaces**:
- `components/StatsSummary.astro:3-24` - Exported interfaces (StatsData, StatsState)
- Used only in stats.astro page

**Presentation Interfaces**:
- `pages/present/index.astro:13-32` - Presentation, PresentationStatus

**SearchState Interface**:
- `components/SearchInterface.astro:451-457` - Client-side state management

### Critical Issues

#### 1. Variable Declaration Order Bug (control.astro)
**Location**: `pages/sync/[id]/control.astro`
- Line 19: Uses `error = 'missing_token'`
- Line 28: Declares `let error = null`
- **Bug**: Variable used before declaration (hoisting issue)

#### 2. Inconsistent Error Handling
Each page has its own error message objects with slight variations in structure and messaging.

#### 3. Mixed Import Patterns
- Some pages import LinkItem from LinkCard.astro
- Others define it inline
- Creates coupling between components and pages

### API Fetch Patterns

#### Common Fetch Structure
All pages follow similar fetch patterns:
1. Create AbortController with timeout
2. Make fetch request with signal
3. Clear timeout on success
4. Parse JSON response
5. Handle errors (404, timeout, network, server)

#### Endpoints Used
- `/api/collections/links/records` - PocketBase collections API
- `/api/links/search` - Custom search API
- `/api/collections/presentations/records` - Presentations
- `/api/collections/sync_sessions/records` - Sync sessions
- `/api/collections/tags/records` - Tags
- `/api/presentations/{id}/status` - Presentation status
- `/api/links/{id}/view` - View count increment

### Component Dependencies

#### SearchInterface.astro
- Defines own LinkItem, SearchResponse, SearchState interfaces in script block
- Uses client-side JavaScript for progressive enhancement
- Duplicates interfaces already in SSR pages

#### StatsSummary.astro
- Exports StatsData and StatsState interfaces
- Has custom API resolution helper
- Only component that properly exports types

#### LinkCard.astro
- Exports LinkItem interface (imported by index.astro)
- Central definition but not used consistently

### Environment Configuration
- Uses `import.meta.env.PUBLIC_API_URL` for API base URL
- Fallback to `http://localhost:8090`
- Some client scripts duplicate this logic

### Build System
- Astro SSR with `prerender = false` on all dynamic pages
- TypeScript support enabled
- No explicit type checking configuration found

### Existing Utils Directory
`frontend/src/utils/` exists but usage unclear. Could house new `api.ts` utilities.

### Summary of Required Changes

1. **Create Type Definitions** (`types/api.ts`)
   - Consolidate 5+ LinkItem definitions
   - Unify SearchResponse interface
   - Standardize PocketBaseResponse
   - Move Stats and Presentation interfaces

2. **Create API Utilities** (`lib/api.ts`)
   - Shared API_BASE constant
   - Shared FETCH_TIMEOUT constant
   - Generic apiFetch helper with AbortController

3. **Refactor 9+ Pages**
   - Remove inline type definitions
   - Import from shared modules
   - Fix control.astro variable declaration bug

4. **Update Components**
   - Remove duplicate interfaces from SearchInterface script
   - Re-export types from types/api.ts in StatsSummary

Total affected files: ~13 (9 pages, 3 components, 2 new files)