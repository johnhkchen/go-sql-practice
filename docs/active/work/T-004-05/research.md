# Research: Tag Page Implementation (T-004-05)

## Objective
Research the codebase structure to build a tag page at `/tags/[slug]` that displays all links for a given tag, using Astro's dynamic routing and the existing search endpoint.

## Ticket Context
- **ID**: T-004-05
- **Type**: task
- **Dependencies**: T-004-01 (astro-layout-and-nav), T-003-01 (search-endpoint)
- **Goal**: Dynamic route showing filtered links by tag with proper fallbacks

## Current Frontend Architecture

### Astro Setup
- **Version**: 5.17.3
- **Output**: Static (configured in astro.config.mjs)
- **Adapter**: Node.js standalone mode
- **Build Format**: Directory-based
- **Node Requirement**: >=24

### Directory Structure
```
frontend/src/
├── components/
│   ├── Navigation.astro           # Responsive nav bar
│   └── StatsSummary.astro         # Complex stats display component
├── layouts/
│   └── BaseLayout.astro           # Main layout with CSS variables
├── pages/
│   ├── index.astro                # Simple welcome page
│   ├── stats.astro                # Stats page using StatsSummary
│   └── watch/
│       └── [id].astro             # Dynamic route example
└── styles/
    └── animations.css             # Pulse/shimmer animations
```

### BaseLayout Analysis
**File**: `/frontend/src/layouts/BaseLayout.astro`

**Props Interface**:
```typescript
interface Props {
  title: string;
  description?: string;
}
```

**CSS Architecture**:
- Global CSS variables for theming (colors, spacing, layout)
- Mobile-first responsive design
- Flexbox layout with sticky header
- CSS custom properties: --color-*, --space-*, --max-width
- Breakpoint: 767px for mobile/desktop distinction

**Key Variables**:
```css
--color-bg: #ffffff
--color-text: #333333
--color-primary: #111111
--color-border: #e5e5e5
--color-footer: #f5f5f5
--space-xs/sm/md/lg/xl: 0.25rem to 3rem
--max-width: 1200px
--header-height: 60px
```

### Navigation Component Analysis
**File**: `/frontend/src/components/Navigation.astro`

**Current Links**:
- Home (`/`)
- Stats (`/stats`)

**Mobile Behavior**:
- Hamburger menu with checkbox toggle
- Hidden nav links revealed on mobile
- CSS-only implementation (no JavaScript)
- Smooth transitions for menu states

**Styling Pattern**:
- Uses CSS variables from BaseLayout
- Hover effects with ::after pseudo-elements
- Responsive with mobile-specific styles

### Dynamic Routing Pattern
**Reference**: `/frontend/src/pages/watch/[id].astro`

**Key Patterns Identified**:
1. **SSR Configuration**: `export const prerender = false`
2. **Parameter Access**: `const { id } = Astro.params`
3. **Error Handling**: Redirect to `/404` for missing params
4. **API Integration**: Fetch data with timeout and error states
5. **State Management**: Multiple view states (live/waiting/error)
6. **Progressive Enhancement**: Client-side JavaScript for real-time updates

**Error Handling Strategy**:
- 404 errors → Custom error message
- Network errors → Generic connection message
- Timeout errors → "Request took too long"
- Server errors → "Something went wrong"

## API Endpoint Analysis

### Search Endpoint
**Endpoint**: `GET /api/links/search`
**File**: `/routes/links_search.go`

**Parameters**:
```go
type SearchParams struct {
    Q       string  // Full-text search query
    Tag     string  // Tag slug filter
    Page    int     // Page number (default: 1)
    PerPage int     // Items per page (default: 20, max: 100)
}
```

**Response Structure**:
```go
type SearchResponse struct {
    Items      []LinkItem `json:"items"`
    Page       int        `json:"page"`
    PerPage    int        `json:"perPage"`
    TotalItems int        `json:"totalItems"`
}

type LinkItem struct {
    ID          string   `json:"id"`
    URL         string   `json:"url"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    ViewCount   int      `json:"view_count"`
    Tags        []string `json:"tags"`  // Array of tag slugs
    Created     string   `json:"created"`
    Updated     string   `json:"updated"`
}
```

**SQL Implementation**:
- Raw SQL queries using app.DB() for learning purposes
- JOIN with tags table via `t.slug = ?` parameter
- Tag relationships through JSON array in `links.tags`
- Parameterized queries with manual string interpolation
- Pagination with LIMIT/OFFSET
- Separate count query for totalItems

**Tag Filtering Logic**:
```sql
JOIN json_each(l.tags) AS jt ON 1=1
JOIN tags t ON t.id = jt.value
WHERE t.slug = ?
```

## Data Model

### Database Collections
**From**: `/migrations/collections.go`

**Tags Collection**:
- `name`: TextField (required, 1-100 chars)
- `slug`: TextField (required, unique index, regex pattern)
- Pattern: `^[a-z0-9]+(?:-[a-z0-9]+)*$`
- Public read access (ListRule/ViewRule: "")

**Links Collection**:
- `url`: URLField (required)
- `title`: TextField (required)
- `description`: TextField (optional)
- `view_count`: NumberField (default: 0)
- `tags`: JSON array of tag IDs
- Created/updated timestamps (automatic)

**Relationship Pattern**:
- Links → Tags: Many-to-many via JSON array
- Tag filtering uses SQL JSON functions
- No foreign key constraints (NoSQL-style)

## Component Patterns

### StatsSummary Analysis
**File**: `/frontend/src/components/StatsSummary.astro`

**Relevant Patterns for Link Cards**:
1. **Ranked Lists**: `.ranked-item` CSS class with flex layout
2. **Card Structure**: Background, border, padding, hover effects
3. **Link Styling**: Primary color, hover underlines
4. **Loading States**: Skeleton animations, shimmer effects
5. **Error Handling**: Retry buttons, empty state messages
6. **Responsive Design**: Grid layouts that adapt to mobile

**Card CSS Pattern**:
```css
.ranked-item {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding: var(--space-sm);
    background-color: var(--color-bg);
    border-radius: 4px;
    transition: background-color 0.2s ease;
}

.ranked-item:hover {
    background-color: var(--color-footer);
}
```

### No Existing Link Card Component
**Finding**: No dedicated link card component exists yet. The home page (`/pages/index.astro`) is minimal with just a welcome message.

**Implication**: Need to create new LinkCard component based on patterns from StatsSummary.

## Dependency Analysis

### T-004-01: Astro Layout and Navigation (Completed)
**Deliverables Created**:
- `/frontend/src/layouts/BaseLayout.astro` - Main layout shell
- `/frontend/src/components/Navigation.astro` - Responsive navigation
- CSS variable system for consistent theming
- Mobile-responsive patterns

**Design Decisions**:
- Minimal scoped components (no external CSS frameworks)
- CSS variables for design tokens
- Progressive enhancement approach
- Mobile-first responsive design

### T-003-01: Search Endpoint (Completed)
**Deliverables Created**:
- `GET /api/links/search` endpoint with tag filtering
- Raw SQL implementation for learning purposes
- Parameterized queries with security measures
- Pagination support (page/perPage)
- JSON response format with metadata

**Key Features for Tag Page**:
- Tag slug filtering: `?tag=database`
- Combines with text search: `?q=intro&tag=golang`
- Pagination for large result sets
- Returns link metadata and tag arrays

## Technical Constraints

### Astro Configuration
- **Static Output**: Need SSR for dynamic routes (`prerender: false`)
- **No Hydration**: Minimal JavaScript usage preferred
- **Directory Format**: Generated routes follow `/tags/[slug]/index.html`

### API Configuration
- **Base URL**: Environment-dependent (PUBLIC_API_URL or localhost:8090)
- **Timeout Handling**: 5-second timeout pattern established
- **Error States**: Consistent error message structure

### Performance Considerations
- **Search Endpoint**: Uses LIKE queries (adequate for current scale)
- **Tag Lookup**: Separate query to avoid N+1 problems
- **Pagination**: Default 20 items, max 100 per page

## Patterns for Reuse

### Dynamic Route Implementation
**Template** (from watch/[id].astro):
```astro
---
export const prerender = false;
const { slug } = Astro.params;
if (!slug) return Astro.redirect('/404');

// API fetch with timeout + error handling
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
let data = null, error = null;

try {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 5000);

  const response = await fetch(`${API_BASE}/api/links/search?tag=${slug}`, {
    signal: controller.signal
  });

  if (response.ok) {
    data = await response.json();
  } else if (response.status === 404) {
    error = 'notfound';
  }
} catch (err) {
  error = err.name === 'AbortError' ? 'timeout' : 'network';
}
---
```

### Component Architecture
**Suggested LinkCard Component**:
- Follow StatsSummary's `.ranked-item` pattern
- Include hover states and accessibility
- Display: title, description, URL, view count, tags
- Link to external URL with proper attributes

### Styling Integration
- Use existing CSS variables from BaseLayout
- Follow mobile-first responsive patterns
- Implement loading states with skeleton animations
- Use consistent spacing (--space-* variables)

## Remaining Unknowns

### Tag Validation
- **Question**: Should non-existent tags show 404 or empty results?
- **Current API**: Returns empty results (200 status)
- **Recommendation**: Follow API behavior for consistency

### Navigation Integration
- **Question**: Add "Tags" link to Navigation component?
- **Current State**: Only Home and Stats links exist
- **Consideration**: May be out of scope for this ticket

### Search Integration
- **Question**: Add search box to tag pages?
- **Current Scope**: Focus on tag-filtered display only
- **Future Enhancement**: Could combine with text search

## Implementation Path

### Phase 1: Route Structure
1. Create `/frontend/src/pages/tags/[slug].astro`
2. Add SSR configuration and parameter extraction
3. Implement API integration with search endpoint
4. Add timeout and error handling patterns

### Phase 2: Component Development
1. Create LinkCard component based on StatsSummary patterns
2. Implement responsive grid layout for link cards
3. Add loading states with skeleton animations
4. Style using existing CSS variables

### Phase 3: Error Handling
1. 404 fallback for invalid tag slugs
2. Empty state when no links match tag
3. Network error handling with retry options
4. Loading states during API calls

### Phase 4: Integration
1. Test with existing seeded data
2. Verify responsive behavior across devices
3. Validate accessibility patterns
4. Performance testing with pagination

## Success Criteria Validation

✅ **Dynamic Route**: `/tags/[slug]` pattern established
✅ **API Integration**: `/api/links/search?tag=:slug` endpoint available
✅ **Component Reuse**: StatsSummary patterns identified for LinkCard
✅ **Error Handling**: 404 and empty state patterns defined
✅ **Responsive Design**: BaseLayout and CSS variables ready
✅ **Accessibility**: Screen reader patterns available from StatsSummary

The codebase provides all necessary foundations for implementing the tag page feature with consistent design patterns and robust error handling.