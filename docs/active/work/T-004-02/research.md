# T-004-02 Research: Home Page Link List

## Overview

This ticket implements the main home page that displays all links with their tags and provides search functionality. The page leverages Astro's island architecture to provide static initial content with interactive search capabilities.

## Current Architecture

### Backend - PocketBase + Custom Routes

**Core Components:**
- PocketBase server with SQLite backend
- Custom Go routes registered via `routes/routes.go`
- Collections: `links`, `tags`, `sync_sessions`, `presentations`
- Static file serving through embedded frontend files

**Key Files:**
- `main.go:8-23` - PocketBase app initialization with custom routes
- `routes/routes.go:8-23` - Route registration system
- `migrations/collections.go:87-146` - Links and tags collections schema

**Data Model:**
```
links:
  - id (auto)
  - url (required, URL field)
  - title (required, text 1-500 chars)
  - description (optional, text max 2000 chars)
  - view_count (optional, number, min 0, integer)
  - tags (relation to tags collection, max 100)
  - created_by (relation to users if exists)
  - created, updated (auto timestamps)

tags:
  - id (auto)
  - name (required, text 1-100 chars)
  - slug (required, text 1-100 chars, unique, kebab-case pattern)
  - created, updated (auto timestamps)
```

### Frontend - Astro Static Site

**Architecture:**
- Astro v5.17.3 with Node.js adapter
- Static output mode (build-time rendering)
- Embedded in Go binary via `embed.go:9-10`
- Served through custom SPA filesystem handler

**Current Structure:**
```
frontend/src/
├── layouts/BaseLayout.astro     # Main HTML shell with global CSS
├── components/
│   ├── Navigation.astro         # Header navigation
│   └── StatsSummary.astro       # Stats display component
├── pages/
│   ├── index.astro              # Current placeholder home page
│   ├── stats.astro              # Statistics page
│   └── watch/[id].astro         # Dynamic watch page
└── styles/animations.css        # Animation utilities
```

**Existing Layout System:**
- `BaseLayout.astro` provides consistent HTML structure
- CSS custom properties for theming (colors, spacing, typography)
- Responsive design with mobile-first approach
- Navigation component with links to Home and Stats

## API Integration Points

### Search API - `GET /api/links/search`

**Implementation:** `routes/links_search.go:46-367`
- Full-text search across title and description fields
- Tag filtering by slug
- Pagination support (page, perPage params)
- SQL injection protection via parameter escaping
- Complex JOIN logic for tag relationships

**Response Format:**
```json
{
  "items": [
    {
      "id": "string",
      "url": "string",
      "title": "string",
      "description": "string",
      "view_count": 0,
      "tags": ["tag-slug", "another-tag"],
      "created": "timestamp",
      "updated": "timestamp"
    }
  ],
  "page": 1,
  "perPage": 20,
  "totalItems": 42
}
```

### PocketBase Collections API

**Direct Access:** PocketBase provides built-in REST API
- `GET /api/collections/links/records` - List all links
- Collection API follows PocketBase standard format
- Automatically handles pagination, filtering, sorting
- Public read access configured in migrations

## Dependencies Status

**T-004-01 (astro-layout-and-nav):** ✅ DONE
- BaseLayout.astro implemented with responsive design
- Navigation component with Home/Stats links
- Global CSS system with custom properties

**T-003-01 (search-endpoint):** ✅ DONE
- Search API fully implemented with full-text search
- Tag filtering capability
- Pagination and proper response format
- SQL injection protection

## Technical Constraints

### Astro Client-Side Interactivity
- Islands architecture requires explicit `client:*` directives
- `client:load` - Hydrates immediately on page load
- `client:visible` - Hydrates when component enters viewport
- Search functionality needs client-side JavaScript for interactivity

### Static File Serving
- Frontend files embedded at build time (`embed.go`)
- SPA routing handled by custom filesystem (`routes/static.go:32-45`)
- API routes excluded from SPA fallback logic

### Database Query Patterns
- Complex tag joins via JSON relationship fields
- Manual parameter escaping due to dbx limitations
- DISTINCT queries needed for proper tag filtering

## Current Home Page Status

**Existing Implementation:** `frontend/src/pages/index.astro:1-8`
- Minimal placeholder with BaseLayout
- Static "Welcome" message
- No link display or search functionality

## Frontend Component Patterns

**Existing Components Analysis:**
- `Navigation.astro` - Pure server-side component, no client logic
- `StatsSummary.astro` - Uses Astro fetch for server-side data
- Components use scoped CSS with Astro's `<style>` blocks
- No TypeScript interfaces defined for props (simple string/number props)

**CSS System:**
- CSS custom properties in BaseLayout for consistency
- Responsive design with single mobile breakpoint (767px)
- Clean, minimal aesthetic with subtle borders and spacing
- No external CSS frameworks or utility libraries

## Static vs Dynamic Content Strategy

**Server-Side Rendering (SSR):**
- Initial page load with full link list
- SEO-friendly static HTML generation
- Works without JavaScript

**Client-Side Enhancement:**
- Search interface as Astro island
- JavaScript-powered filtering/pagination
- Progressive enhancement approach

## File System Integration

**Build Process:**
- Astro builds to `frontend/dist/`
- Go embeds `frontend/dist/*` via `//go:embed`
- Static handler serves embedded files with SPA fallback
- Development uses separate `npm run dev` process

This research reveals a well-structured foundation with all dependencies complete. The home page implementation can leverage existing APIs, layout system, and component patterns while adding client-side interactivity for search functionality.