# T-004-02 Plan: Home Page Link List

## Implementation Sequence

This plan breaks down the implementation into ordered, testable steps that can be executed and verified independently. Each step builds incrementally on the previous work while maintaining a working state.

## Step 1: Create LinksList Component
**Duration:** ~20 minutes
**Files:** `frontend/src/components/LinksList.astro`

### Tasks:
1. Create component file with TypeScript interfaces
2. Implement server-side logic for empty states and result counts
3. Add HTML structure with results header and links grid
4. Include scoped CSS for grid layout and empty states
5. Test component with sample data in isolation

### Implementation Details:
```astro
---
import LinkCard from './LinkCard.astro';

interface Props {
  links: LinkItem[];
  totalCount: number;
  searchQuery?: string;
  isSearchResult?: boolean;
  className?: string;
}

const { links, totalCount, searchQuery, isSearchResult = false, className = '' } = Astro.props;
---

<!-- Results context header -->
<!-- Empty state handling -->
<!-- Links grid with LinkCard iteration -->
```

### Testing Criteria:
- Component renders with empty array (shows empty state)
- Component renders with sample links (shows grid)
- Grid layout is responsive (manual browser test)
- CSS scoping works correctly (no style leakage)

### Commit Message:
```
feat: add LinksList component for link grid display

- Server-side component with empty state handling
- Responsive grid layout matching tags page pattern
- Results count and context messaging
- Ready for integration with index.astro

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Step 2: Create SearchInterface Foundation
**Duration:** ~30 minutes
**Files:** `frontend/src/components/SearchInterface.astro`

### Tasks:
1. Create component with TypeScript interfaces for props
2. Implement HTML form structure with search input
3. Add loading, error, and results count display areas
4. Include comprehensive CSS for form styling and states
5. Add non-functional form (server-side only initially)

### Implementation Details:
```astro
---
interface Props {
  initialQuery?: string;
  initialResults?: LinkItem[];
  initialCount?: number;
  className?: string;
}
---

<!-- Search form with proper accessibility -->
<!-- Loading indicator (hidden initially) -->
<!-- Error display area (hidden initially) -->
<!-- Results count display -->
```

### Testing Criteria:
- Form renders correctly with all elements
- CSS styling matches design system variables
- Form submission works (even if non-functional)
- Accessibility: proper labels, focus management
- Responsive design on mobile/desktop

### Commit Message:
```
feat: add SearchInterface component foundation

- HTML form structure with search input
- Loading and error state placeholders
- Responsive CSS styling with design system
- Accessible form elements and labeling
- Ready for client-side enhancement

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Step 3: Update Index Page with Server-Side Data Fetching
**Duration:** ~40 minutes
**Files:** `frontend/src/pages/index.astro`

### Tasks:
1. Add server-side imports for new components
2. Implement data fetching logic (similar to tags page pattern)
3. Add error handling and timeout logic
4. Handle both initial load and search query scenarios
5. Replace existing content with new component composition
6. Add page-specific CSS for layout containers

### Implementation Details:
```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import SearchInterface from '../components/SearchInterface.astro';
import LinksList from '../components/LinksList.astro';

// URL parameter processing
const searchQuery = Astro.url.searchParams.get('q') || '';
const isSearchRequest = !!searchQuery;

// API configuration
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
const FETCH_TIMEOUT = 5000;

// Data fetching with error handling
let links: LinkItem[] = [];
let totalCount = 0;
let error: string | null = null;

try {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT);

  const endpoint = isSearchRequest
    ? `${API_BASE}/api/links/search?q=${encodeURIComponent(searchQuery)}`
    : `${API_BASE}/api/collections/links/records`;

  const response = await fetch(endpoint, {
    signal: controller.signal,
    headers: { 'Accept': 'application/json' }
  });

  clearTimeout(timeoutId);
  // Process response...
} catch (err) {
  // Error handling...
}
---
```

### Testing Criteria:
- Page loads with initial link list (no search query)
- Search URLs work (`/?q=test` shows search results)
- Error handling works (test with invalid API URL)
- Empty states display correctly
- Server-side rendering works without JavaScript

### Commit Message:
```
feat: implement home page with server-side data fetching

- Replace welcome page with dynamic link list
- Server-side data fetching for initial load and search
- Error handling with timeout and abort controller
- Integration with LinksList and SearchInterface components
- Progressive enhancement ready for client-side features

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Step 4: Add Client-Side Search Functionality
**Duration:** ~50 minutes
**Files:** `frontend/src/components/SearchInterface.astro`

### Tasks:
1. Add client-side script block with hydration directive
2. Implement search form event handling and preventDefault
3. Add debounced input handling (300ms delay)
4. Implement fetch logic for search API calls
5. Add DOM manipulation for results updates
6. Include loading states and error handling
7. Implement URL parameter synchronization

### Implementation Details:
```astro
<SearchInterface
  client:visible
  initialQuery={searchQuery}
  initialResults={isSearchRequest ? links : []}
  initialCount={totalCount}
/>

<script>
// State management
interface SearchState {
  query: string;
  isLoading: boolean;
  results: LinkItem[];
  totalCount: number;
  error: string | null;
}

// Form handling with preventDefault
// Debounced input with setTimeout
// Fetch API with abort controller
// DOM manipulation for results update
// History API for URL management
// Accessibility: focus management and live regions
</script>
```

### Testing Criteria:
- Search form submission prevents default browser action
- Typing in search box triggers debounced API calls
- Loading states show during search requests
- Search results update DOM without page reload
- URL updates to reflect search query
- Error states display correctly (test with network issues)
- Accessibility: screen reader announcements work

### Commit Message:
```
feat: add client-side search functionality to SearchInterface

- Form interception with preventDefault
- Debounced search input (300ms delay)
- API integration with loading and error states
- DOM manipulation for seamless result updates
- URL synchronization with History API
- Accessibility features for screen readers

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Step 5: Update LinkCard Tags for Navigation
**Duration:** ~15 minutes
**Files:** `frontend/src/components/LinkCard.astro`

### Tasks:
1. Review current LinkCard tag implementation
2. Ensure tags are clickable links to `/tags/:slug`
3. Update CSS if needed for hover states
4. Test tag navigation works correctly

### Implementation Details:
```astro
<!-- Tags section -->
{link.tags.length > 0 && (
  <div class="link-card-tags" aria-label="Tags">
    {link.tags.map(tag => (
      <a href={`/tags/${tag}`} class="link-card-tag" key={tag}>
        {tag}
      </a>
    ))}
  </div>
)}
```

### Testing Criteria:
- Tags display as clickable links
- Navigation to tag pages works correctly
- Hover states provide visual feedback
- CSS styling remains consistent

### Commit Message:
```
feat: make LinkCard tags clickable for navigation

- Convert tag spans to anchor links
- Add navigation to /tags/:slug pages
- Maintain existing CSS styling and hover states
- Improve user experience with clickable tags

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Step 6: Progressive Enhancement Testing & Polish
**Duration:** ~30 minutes
**Files:** Various CSS and accessibility improvements

### Tasks:
1. Test complete functionality with JavaScript disabled
2. Test search functionality across different browsers
3. Verify responsive design on mobile/tablet/desktop
4. Test accessibility with keyboard navigation
5. Test error scenarios (network timeouts, server errors)
6. Add any missing loading states or polish
7. Update documentation if needed

### Testing Scenarios:

#### JavaScript Disabled Testing:
- Initial page load shows complete link list ✓
- Search form submission goes to `/?q=query` ✓
- Server-side search results display correctly ✓
- Navigation and basic functionality work ✓

#### JavaScript Enabled Testing:
- Search input triggers debounced API calls ✓
- Loading states show during searches ✓
- Results update without page reload ✓
- URL updates to reflect search state ✓
- Browser back/forward buttons work ✓

#### Error Handling Testing:
- Network timeout (disconnect internet) ✓
- Server error (stop backend server) ✓
- Invalid search queries ✓
- Empty search results ✓

#### Accessibility Testing:
- Keyboard navigation through all elements ✓
- Screen reader announcements for search results ✓
- Focus management during search interactions ✓
- High contrast mode compatibility ✓

#### Responsive Design Testing:
- Mobile layout (< 768px) ✓
- Tablet layout (768px - 1023px) ✓
- Desktop layout (> 1024px) ✓
- Grid layout adapts correctly ✓

### Final Commit Message:
```
polish: comprehensive testing and accessibility improvements

- Verify progressive enhancement works correctly
- Test responsive design across all breakpoints
- Validate accessibility with keyboard and screen readers
- Confirm error handling for all edge cases
- Home page feature complete and production ready

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Step 7: Update Ticket Phase to Implement
**Duration:** ~5 minutes
**Files:** `docs/active/tickets/T-004-02.md`

### Tasks:
1. Update ticket frontmatter phase from 'plan' to 'implement'
2. Create initial progress.md file for implementation tracking

### Implementation Details:
```markdown
phase: implement
```

### Commit Message:
```
docs: advance T-004-02 to implementation phase

Plan completed with detailed implementation steps.
Ready to execute home page link list feature.

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Testing Strategy

### Unit Testing Approach
**Manual Component Testing:**
- Test each component with sample data in isolation
- Verify component interfaces match expected props
- Check CSS scoping and responsive behavior
- Validate accessibility attributes and structure

### Integration Testing Approach
**Page-Level Testing:**
- Test complete user journeys (initial load → search → navigation)
- Verify API integration works with real backend
- Check component composition and data flow
- Validate error boundaries and fallback states

### Cross-Browser Compatibility
**Testing Matrix:**
- Chrome/Chromium (primary development browser)
- Firefox (alternative engine testing)
- Safari (WebKit differences)
- Mobile Safari (iOS specific issues)
- Chrome Mobile (Android specific issues)

### Performance Testing
**Load Time Verification:**
- Initial page load under 3 seconds on 3G
- Search response time under 1 second on good connection
- Client-side bundle size reasonable (< 50KB total)
- No unnecessary re-renders or API calls

### Accessibility Testing
**Manual Verification:**
- Keyboard-only navigation complete workflow
- Screen reader testing with NVDA/VoiceOver
- High contrast mode compatibility
- Reduced motion preference respect

### Progressive Enhancement Validation
**Degradation Testing:**
- Complete feature functionality without JavaScript
- Graceful enhancement when JavaScript loads
- No broken states during hydration
- Server-side search fallback works correctly

## Risk Mitigation

### API Dependency Risk
**Mitigation:** Server-side fallbacks and comprehensive error handling
- If search API fails, fall back to PocketBase collections API
- Clear error messages guide user to alternative actions
- Timeout handling prevents hanging requests

### Client-Side Complexity Risk
**Mitigation:** Incremental enhancement and vanilla JavaScript
- Core functionality works without JavaScript
- Minimal external dependencies (zero npm packages)
- State management kept simple and predictable

### Performance Risk
**Mitigation:** Lazy loading and debounced interactions
- Search interface only hydrates when visible (`client:visible`)
- Debounced input prevents excessive API calls
- Minimal JavaScript bundle with tree shaking

### Accessibility Risk
**Mitigation:** Progressive enhancement philosophy and comprehensive testing
- Server-rendered content accessible by default
- Client-side enhancements maintain accessibility
- Comprehensive keyboard and screen reader testing

This plan provides a clear path from current state to fully functional home page, with each step building incrementally while maintaining system stability. The testing strategy ensures robustness across all user scenarios and technical constraints.