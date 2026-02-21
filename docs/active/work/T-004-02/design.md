# T-004-02 Design: Home Page Link List

## Overview

This design implements the main home page that displays all links with search functionality using Astro's island architecture. The page will provide a static initial load with progressive enhancement via client-side JavaScript for search interactivity.

## Design Options Analysis

### Option 1: Pure Server-Side Rendering (SSR)
**Approach:** Similar to tags page - fetch data server-side, render static HTML
**Pros:**
- Works without JavaScript
- Fast initial load
- SEO-friendly
- Consistent with existing tags page pattern
**Cons:**
- Search requires full page reloads
- Poor UX for filtering
- No real-time interactivity

### Option 2: Pure Client-Side App (SPA)
**Approach:** Empty initial page, fetch everything via JavaScript
**Pros:**
- Rich interactive experience
- No page reloads for search
- Can implement advanced filtering
**Cons:**
- Blank page without JavaScript
- Poor SEO
- Slow initial load
- Breaks accessibility principles

### Option 3: Hybrid SSR + Client Islands (CHOSEN)
**Approach:** Server-render initial content, enhance with client-side search island
**Pros:**
- Progressive enhancement philosophy
- Works with and without JavaScript
- Fast initial load with rich interactions
- SEO-friendly with accessibility
- Matches Astro's intended architecture
**Cons:**
- More complex implementation
- Need to manage state synchronization

## Chosen Architecture: Hybrid SSR + Islands

### Core Principles
1. **Progressive Enhancement:** Page works without JavaScript, enhanced with it
2. **Static First:** Initial page load renders full link list server-side
3. **Interactive Second:** Search functionality added as client-side island
4. **Consistent Patterns:** Follow existing component and styling conventions

### Technical Implementation Strategy

**Server-Side Foundation:**
- `index.astro` fetches initial links via PocketBase Collections API
- Renders complete HTML with LinkCard components
- Includes fallback states (loading, error, empty)
- Search form rendered but non-functional without JS

**Client-Side Enhancement:**
- Search island component hydrated with `client:visible` or `client:load`
- Intercepts search form submission
- Calls `/api/links/search` endpoint
- Updates DOM with filtered results
- Maintains URL state for bookmarkable searches

## Component Architecture

### Page Structure
```
index.astro (SSR)
├── BaseLayout.astro
├── SearchInterface.astro (Island - client:visible)
│   ├── SearchForm
│   ├── FilterControls (future)
│   └── ResultsCount
└── LinksList.astro
    └── LinkCard.astro (multiple instances)
```

### Data Flow
1. **Initial Load:** Server fetches all links → render LinkCard grid
2. **Search Interaction:** Client-side island intercepts form → API call → DOM update
3. **URL Management:** Search terms reflected in URL parameters
4. **State Management:** Simple component state, no external store needed

### Search Interface Design

**Search Form Elements:**
- Text input for query (title/description search)
- Submit button (functional without JS)
- Clear/reset option
- Results count display
- Loading indicator during search

**Search Behavior:**
- **No JS:** Form submits to `/?q=search-term` → SSR handles query
- **With JS:** Form intercepted → AJAX search → DOM update
- **Debounced Input:** 300ms delay on typing for better UX
- **URL Sync:** Search state reflected in browser URL

### API Integration Strategy

**Primary API:** `GET /api/links/search`
- Query parameters: `q` (search), `tag` (filter), `page`, `per_page`
- Same endpoint used by tags page - proven reliable
- Handles both text search and tag filtering

**Fallback API:** `GET /api/collections/links/records`
- PocketBase built-in endpoint for initial load
- More reliable for server-side rendering
- Used when no search parameters present

**Error Handling:**
- Network timeouts (5s limit)
- Server errors (5xx responses)
- Invalid responses
- Graceful degradation to static content

## UI/UX Design

### Layout Structure
```
Header: Site Navigation (existing)
Main Content:
  ├── Page Title + Subtitle
  ├── Search Interface (prominent, centered)
  ├── Results Summary ("Showing X links")
  ├── Links Grid (responsive, LinkCard components)
  └── Pagination Controls (future enhancement)
Footer: (existing)
```

### Search Interface Positioning
- **Prominent Placement:** Below page title, above content
- **Center Aligned:** Focus user attention
- **Responsive Design:** Full width on mobile, constrained on desktop
- **Visual Priority:** Distinct styling to indicate interactivity

### Results Display
- **Grid Layout:** Reuse existing LinkCard component and grid CSS
- **Responsive Columns:** 1 col mobile, 2-3 desktop (match tags page)
- **Smooth Transitions:** Fade in/out for search results
- **Loading States:** Skeleton cards during search
- **Empty States:** Friendly message for no results

### Tag Integration
- **Clickable Tags:** Navigate to `/tags/:slug` (existing behavior)
- **Search Refinement:** Clicking tag adds to search query (future enhancement)
- **Visual Consistency:** Reuse LinkCard tag styling

## Technical Implementation Details

### File Structure
```
frontend/src/pages/index.astro           # Main page (SSR + islands)
frontend/src/components/SearchInterface.astro   # Search island component
frontend/src/components/LinksList.astro         # Links display wrapper
frontend/src/components/LinkCard.astro          # (existing, reused)
```

### Search Island Component
```typescript
// SearchInterface.astro
interface SearchState {
  query: string;
  isLoading: boolean;
  results: LinkItem[];
  totalCount: number;
  error: string | null;
}

// Key features:
- Form handling with preventDefault
- Debounced input (300ms)
- URL parameter sync
- Loading/error states
- Results count display
- Clear functionality
```

### Server-Side Data Fetching
```typescript
// index.astro server-side logic
const searchQuery = Astro.url.searchParams.get('q') || '';
const isSearchRequest = !!searchQuery;

let links: LinkItem[] = [];
let totalCount = 0;
let error: string | null = null;

if (isSearchRequest) {
  // Use search API for server-side search
  response = await fetch(`/api/links/search?q=${encodeURIComponent(searchQuery)}`);
} else {
  // Use collections API for initial load
  response = await fetch(`/api/collections/links/records`);
}
```

### CSS Strategy
- **Reuse Existing:** Leverage BaseLayout variables and LinkCard styles
- **Search Styling:** New CSS for search interface, consistent with design system
- **Responsive Grid:** Follow tags page pattern (350px min column width)
- **Loading States:** Simple pulse animation for search loading
- **Focus Management:** Clear focus indicators for search form

### JavaScript Approach
- **Vanilla JS:** No external dependencies, keep bundle small
- **Progressive Enhancement:** Enhance existing HTML form
- **Event Delegation:** Efficient event handling
- **History API:** Update URL without page reload
- **Accessibility:** Maintain focus, announce results to screen readers

## Error Handling & Edge Cases

### Network Issues
- **Timeout Handling:** 5 second limit with abort controller
- **Retry Logic:** Manual retry button on failures
- **Fallback Content:** Show initial server-rendered content
- **User Feedback:** Clear error messages with suggested actions

### Empty States
- **No Results:** Friendly message with search suggestions
- **No Links:** Guide user to add content (future)
- **Loading State:** Show skeleton cards during fetch

### URL State Management
- **Search Queries:** Reflected in URL parameters
- **Bookmarkable:** Users can share search URLs
- **History Navigation:** Back/forward buttons work correctly
- **Clean URLs:** Remove empty parameters

### Accessibility Considerations
- **Keyboard Navigation:** All functionality accessible via keyboard
- **Screen Reader:** Proper ARIA labels and live regions
- **Focus Management:** Clear focus indicators and logical tab order
- **High Contrast:** Tested in high contrast mode
- **Reduced Motion:** Respect user motion preferences

## Performance Optimization

### Initial Load Performance
- **Static Generation:** Fast server-rendered HTML
- **Critical CSS:** Inline styles for above-fold content
- **Preload APIs:** Hint browser about likely API calls
- **Image Optimization:** Defer non-critical LinkCard assets

### Search Performance
- **Debounced Input:** Prevent excessive API calls
- **Request Deduplication:** Cancel previous requests
- **Caching Strategy:** Client-side cache for repeated searches
- **Pagination:** Limit initial results, load more on demand

### Bundle Size
- **No External Dependencies:** Pure vanilla JavaScript
- **Tree Shaking:** Only include used Astro features
- **CSS Scoping:** Leverage Astro's scoped styles

## Testing Strategy

### Manual Testing
- **JavaScript Disabled:** Verify core functionality works
- **Slow Network:** Test timeout handling and loading states
- **Various Screen Sizes:** Responsive design verification
- **Keyboard Only:** Complete accessibility testing

### Integration Points
- **API Compatibility:** Verify search API responses match LinkCard interface
- **URL Handling:** Test various search parameter combinations
- **Component Integration:** Ensure LinkCard handles dynamic data correctly

This design provides a solid foundation that balances performance, accessibility, and user experience while following established patterns in the codebase. The hybrid approach ensures the page works for all users while providing enhanced functionality for those with JavaScript enabled.