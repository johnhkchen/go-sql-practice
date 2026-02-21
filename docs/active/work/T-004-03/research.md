# Research: Link Detail Page Implementation (T-004-03)

## Codebase Overview

### Current State Analysis

The project is a Go-based link bookmarks application using PocketBase as the backend with an Astro frontend. Dependencies T-004-01 (base layout) and T-003-02 (view count endpoint) are complete and functional.

### Astro Dynamic Route Structure

**Existing Dynamic Route Pattern**: `/frontend/src/pages/watch/[id].astro`
- Uses `export const prerender = false` for SSR
- Extracts ID via `Astro.params` destructuring: `const { id } = Astro.params`
- Implements comprehensive error handling with specific error types:
  - `notfound`: 404 responses from API
  - `network`: Connection failures
  - `server`: 5xx HTTP errors
  - `timeout`: AbortController timeout after 5s
- Uses fetch with AbortController for timeout handling
- API calls target `${API_BASE}/api/collections/presentations/records/${id}`
- Redirect pattern: `return Astro.redirect('/404')` for missing params

**Error Handling Architecture**:
- Centralized error messages object with title/message pairs
- Conditional rendering based on error state
- Graceful fallbacks with back-to-home links
- Consistent error container styling

### Link Data Structure & API Endpoints

**Database Schema** (from `/migrations/collections.go`):
```
links collection:
- id (auto-generated)
- url (URLField, required)
- title (TextField, required, max 500 chars)
- description (TextField, optional, max 2000 chars)
- view_count (NumberField, optional, integer, min 0)
- tags (RelationField to tags collection, max 100)
- created_by (RelationField to users, optional)
- created/updated (auto timestamps)
```

**Tags Schema**:
```
tags collection:
- id (auto-generated)
- name (TextField, 1-100 chars)
- slug (TextField, regex: ^[a-z0-9]+(?:-[a-z0-9]+)*$, unique)
```

**Existing API Endpoints**:
- `GET /api/links/search` - Full search with pagination (complex query builder)
- `GET /api/links/search-simple` - Simplified search variant
- `POST /api/links/:id/view` - **READY TO USE** - Atomically increments view count, returns updated record
- `GET /api/stats` - Aggregate statistics
- PocketBase default: `GET /api/collections/links/records/:id` - Direct record access

**Link Data Format** (from search endpoints):
```typescript
interface LinkItem {
  id: string
  url: string
  title: string
  description: string
  view_count: number
  tags: string[] // Array of tag slugs
  created: string
  updated: string
}
```

### Tag Display & Routing Patterns

**Tag Structure**: Links store tag relations as JSON array of tag IDs, but API responses convert these to tag slugs for display.

**Tag Fetching Logic** (from `/routes/links_search.go`):
- Uses `fetchTagsForLinks()` function that JOINs through `json_each(l.tags)`
- Query: `SELECT l.id, t.slug FROM links l JOIN json_each(l.tags) AS jt ON 1=1 JOIN tags t ON t.id = jt.value`
- Returns `map[string][]string` (link_id -> tag_slugs)

**Tag Link Pattern**: Based on search endpoint, tags link to `/tags/:slug` would use tag filtering:
- Search API supports `?tag=slug-name` parameter
- Frontend would need to implement tag pages or filter interfaces

### Frontend Architecture & Patterns

**Layout System**:
- `BaseLayout.astro` provides HTML shell with global styles
- CSS custom properties system for theming
- Navigation component with sticky header
- Responsive design with mobile-first approach

**Component Patterns**:
- Complex state management in `StatsSummary.astro` using TypeScript classes
- Client-side controllers for API interaction
- Loading states with skeleton animations
- Progressive enhancement with fallbacks

**API Communication**:
- Hardcoded API base: `http://127.0.0.1:8094` (in StatsSummary) and `http://localhost:8090` (in watch page)
- Environment variable support: `import.meta.env.PUBLIC_API_URL`
- Fetch with error handling and timeout patterns
- JSON response parsing with validation

**State Management**: No global state management - each component handles its own state with TypeScript classes and DOM manipulation.

### Error Handling Approaches

**SSR Error Handling** (from watch page):
- Fetch errors caught in try/catch during render
- Error states stored in component variables
- Conditional rendering based on error type
- User-friendly error messages with specific guidance

**Client-Side Error Handling** (from StatsSummary):
- Error states in component state objects
- DOM updates via class-based controllers
- Retry mechanisms with manual refresh buttons
- Screen reader announcements via aria-live regions

### Frontend Dependencies & Build

**Astro Configuration**: Uses standard Astro setup with TypeScript support
- `.astro/types.d.ts` - Auto-generated type definitions
- `frontend/src/styles/animations.css` - Shared animation library
- Scoped styling within components
- CSS custom properties for theming consistency

**Animation Library**: Existing animations include:
- `pulse`, `breathing`, `shimmer` keyframes
- Utility classes: `.animate-pulse`, `.animate-breathing`, `.animate-shimmer`
- `prefers-reduced-motion` accessibility support

### Database View Count Implementation

**View Count Endpoint** (`/routes/links_view.go`):
- **ATOMIC**: Uses `UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?`
- Returns 404 if link doesn't exist (checks `RowsAffected`)
- Returns full updated record using `dao.FindRecordById("links", linkId)`
- No authentication required
- Handles concurrent requests safely

### Constraints & Assumptions

**API Access Patterns**:
- PocketBase collections have public read access (`ListRule` and `ViewRule` set to `""`)
- No authentication required for view operations
- Direct collection access available via `/api/collections/:collection/records/:id`

**Data Consistency**:
- Tag relationships stored as JSON arrays in links table
- Tag display requires JOIN queries to resolve slugs from IDs
- View counts initialize to NULL, handled by COALESCE in increment operation

**Frontend Limitations**:
- No client-side routing library - relies on browser navigation
- No centralized loading/error state management
- Component communication via DOM events or direct method calls

**Performance Considerations**:
- Search endpoints use complex JOINs for tag resolution
- No caching layer identified in current implementation
- Timeout set to 5 seconds for API calls

### Relevant Files for Implementation

**Backend**:
- `/routes/routes.go` - Route registration (may need single link endpoint)
- `/routes/links_view.go` - View count increment (ready)
- `/migrations/collections.go` - Database schema reference

**Frontend**:
- `/frontend/src/pages/watch/[id].astro` - Dynamic route pattern reference
- `/frontend/src/layouts/BaseLayout.astro` - Layout to extend
- `/frontend/src/components/StatsSummary.astro` - Complex component pattern
- `/frontend/src/styles/animations.css` - Animation utilities

**Infrastructure**:
- PocketBase handles `/api/collections/links/records/:id` automatically
- Tag resolution may require custom endpoint or client-side additional call