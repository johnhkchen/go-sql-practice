# T-004-02 Structure: Home Page Link List

## Architecture Overview

This implementation follows the hybrid SSR + client islands pattern established in the design. The structure creates minimal new files while maximizing reuse of existing components and patterns.

## File-Level Changes

### New Files

#### `/frontend/src/components/SearchInterface.astro`
**Purpose:** Client-side search island component
**Type:** Astro component with client-side JavaScript
**Dependencies:** None (vanilla JS)
**Interface:**
```typescript
interface Props {
  initialQuery?: string;
  initialResults?: LinkItem[];
  initialCount?: number;
  className?: string;
}
```

**Responsibilities:**
- Render search form HTML structure
- Handle client-side search interactions
- Manage search state (query, loading, results, errors)
- Update DOM with search results
- Sync URL parameters with search state
- Provide loading and error states

#### `/frontend/src/components/LinksList.astro`
**Purpose:** Container for links display logic
**Type:** Server-side Astro component
**Dependencies:** LinkCard.astro (existing)
**Interface:**
```typescript
interface Props {
  links: LinkItem[];
  totalCount: number;
  searchQuery?: string;
  isSearchResult?: boolean;
  className?: string;
}
```

**Responsibilities:**
- Render grid layout container
- Handle empty states (no links, no search results)
- Display results count and context messaging
- Iterate over links array and render LinkCard components
- Provide consistent grid responsive behavior

### Modified Files

#### `/frontend/src/pages/index.astro` (MAJOR MODIFICATION)
**Current State:** Simple welcome page with BaseLayout
**New Structure:**
```astro
---
// Server-side data fetching logic
import BaseLayout from '../layouts/BaseLayout.astro';
import SearchInterface from '../components/SearchInterface.astro';
import LinksList from '../components/LinksList.astro';

// API configuration and data fetching
// Error handling and response processing
// Search query parameter handling
---

<BaseLayout title="Link Bookmarks">
  <div class="home-container">
    <header class="page-header">
      <h1>Link Bookmarks</h1>
      <p>Discover and search through your collected links</p>
    </header>

    <SearchInterface
      client:visible
      initialQuery={searchQuery}
      initialResults={isSearchRequest ? links : []}
      initialCount={totalCount}
    />

    <LinksList
      links={links}
      totalCount={totalCount}
      searchQuery={searchQuery}
      isSearchResult={isSearchRequest}
    />
  </div>
</BaseLayout>

<style>
  /* Page-specific styles */
</style>
```

**Key Changes:**
- Remove simple welcome content
- Add server-side data fetching (similar to tags page pattern)
- Import and compose new components
- Handle both initial load and search scenarios
- Add error handling and empty states
- Include page-specific styling

#### `/frontend/src/components/LinkCard.astro` (MINOR MODIFICATION)
**Current State:** Fully implemented with proper interface
**Modifications:**
- Make tags clickable for navigation (if not already implemented)
- Ensure proper handling of search highlighting (future enhancement)
- Verify interface compatibility with search API response

**Changes Required:**
```astro
<!-- In tags section -->
{link.tags.map(tag => (
  <a href={`/tags/${tag}`} class="link-card-tag" key={tag}>{tag}</a>
))}
```

**Rationale:** LinkCard component is already well-structured with proper interface. Only need to ensure tag navigation works correctly.

## Component Boundaries & Responsibilities

### Component Hierarchy
```
BaseLayout (existing)
└── index.astro (page controller)
    ├── SearchInterface (client island)
    │   ├── Search Form (HTML form with submit handling)
    │   ├── Loading Indicator (CSS-only spinner)
    │   ├── Results Count Display (dynamic content)
    │   └── Error Display (conditional message)
    └── LinksList (server-side container)
        ├── Empty State Display (conditional)
        ├── Results Context Header (count, query info)
        └── LinkCard[] (existing component array)
```

### Data Flow Architecture

#### Server-Side Flow (SSR)
```
1. index.astro receives request
2. Parse URL parameters (q, page, etc.)
3. Determine API endpoint (search vs. collections)
4. Fetch initial data with timeout/error handling
5. Process response into LinkItem[] format
6. Pass data to LinksList and SearchInterface
7. Render complete HTML with hydration markers
```

#### Client-Side Flow (Island Hydration)
```
1. SearchInterface hydrates with client:visible
2. Initialize state from props (initialQuery, initialResults)
3. Bind event handlers to form elements
4. On search input: debounce → fetch → update DOM
5. Update URL parameters without page reload
6. Maintain focus management and accessibility
```

## Public Interface Definitions

### LinkItem Interface (Shared)
```typescript
// Used across LinkCard, SearchInterface, LinksList
interface LinkItem {
  id: string;
  title: string;
  url: string;
  description: string;
  tags: string[];
  created_at: string;
  view_count: number;
}
```

### API Response Interfaces
```typescript
// PocketBase Collections API response
interface PocketBaseResponse {
  items: LinkItem[];
  page: number;
  perPage: number;
  totalItems: number;
  totalPages: number;
}

// Search API response
interface SearchResponse {
  links: LinkItem[];
  page: number;
  per_page: number;
  total_count: number;
  total_pages: number;
}
```

### Component Props Interfaces
```typescript
// SearchInterface.astro
interface SearchInterfaceProps {
  initialQuery?: string;
  initialResults?: LinkItem[];
  initialCount?: number;
  className?: string;
}

// LinksList.astro
interface LinksListProps {
  links: LinkItem[];
  totalCount: number;
  searchQuery?: string;
  isSearchResult?: boolean;
  className?: string;
}

// LinkCard.astro (existing, verified)
interface LinkCardProps {
  link: LinkItem;
  showIndex?: number;
  className?: string;
}
```

## Internal Module Organization

### SearchInterface Component Structure
```
SearchInterface.astro
├── TypeScript Interface Definitions
├── Props Destructuring and Defaults
├── Server-side Setup (minimal)
├── HTML Structure
│   ├── Form Elements
│   ├── Loading States
│   ├── Error Display
│   └── Results Count
├── Scoped CSS Styles
└── Client-side Script Block
    ├── State Management
    ├── API Integration
    ├── DOM Manipulation
    ├── URL Management
    └── Event Handling
```

### LinksList Component Structure
```
LinksList.astro
├── TypeScript Interface Definitions
├── Props Processing and Validation
├── Conditional Logic (empty states, search context)
├── HTML Structure
│   ├── Results Header
│   ├── Empty State Templates
│   └── Links Grid Container
└── Scoped CSS Styles
    ├── Grid Layout Rules
    ├── Empty State Styling
    └── Results Context Styling
```

### index.astro Page Structure
```
index.astro
├── Import Statements
├── Server-side Logic Block
│   ├── URL Parameter Processing
│   ├── API Configuration
│   ├── Data Fetching with Error Handling
│   └── Response Processing
├── HTML Template
│   ├── Page Header
│   ├── Component Composition
│   └── Container Structure
└── Page-specific CSS
    ├── Layout Containers
    ├── Responsive Behavior
    └── Component Integration Styles
```

## File Dependencies

### Dependency Graph
```
index.astro
├── BaseLayout.astro (existing)
├── SearchInterface.astro (new)
└── LinksList.astro (new)
    └── LinkCard.astro (existing)

SearchInterface.astro
└── (no component dependencies - vanilla JS)

LinksList.astro
└── LinkCard.astro (existing)
```

### API Dependencies
```
index.astro → GET /api/collections/links/records (PocketBase)
index.astro → GET /api/links/search (custom endpoint)
SearchInterface.astro → GET /api/links/search (client-side)
```

### Style Dependencies
```
All components → BaseLayout.astro (CSS custom properties)
SearchInterface.astro → No external CSS dependencies
LinksList.astro → Inherits grid patterns from tags page
LinkCard.astro → (existing styles, no changes)
```

## Change Sequencing

### Phase 1: Foundation Components
1. Create `LinksList.astro` (server-side only, no JavaScript)
2. Create empty `SearchInterface.astro` (HTML form only, no interactivity)
3. Test components in isolation with sample data

### Phase 2: Page Integration
1. Modify `index.astro` to add server-side data fetching
2. Integrate LinksList component with real data
3. Add SearchInterface component (non-functional initially)
4. Test complete page render and navigation

### Phase 3: Client-Side Enhancement
1. Add JavaScript functionality to SearchInterface
2. Implement search API integration
3. Add DOM manipulation and state management
4. Test progressive enhancement (with/without JS)

### Phase 4: Polish and Testing
1. Update LinkCard tags to be clickable (if needed)
2. Add loading states and error handling
3. Implement URL parameter syncing
4. Comprehensive testing and accessibility review

## Integration Points

### Existing System Integration
- **Navigation:** Links to home page already exist in Navigation.astro
- **API Endpoints:** Reuse existing search API and PocketBase endpoints
- **Styling System:** Leverage existing CSS custom properties and patterns
- **Component Patterns:** Follow LinkCard and tags page conventions

### External Dependencies
- **None Added:** Implementation uses only existing Astro and browser APIs
- **API Compatibility:** Matches existing search API response format
- **CSS Framework:** Continues with vanilla CSS and custom properties

### Build Process Integration
- **No Changes Required:** Standard Astro component compilation
- **Static Generation:** Components support both SSR and static generation
- **JavaScript Bundling:** Astro handles client-side script bundling automatically

This structure provides a clear implementation path that builds incrementally on the existing codebase while introducing minimal new complexity. The component boundaries are clean, dependencies are minimal, and the progressive enhancement approach ensures robustness across different client capabilities.