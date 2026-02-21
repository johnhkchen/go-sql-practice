# Structure: Stats Page Implementation (T-004-04)

## Files to Create

### 1. Main Stats Page
**File**: `frontend/src/pages/stats.astro`
**Purpose**: Main stats page route handler
**Type**: New file
**Dependencies**: BaseLayout.astro, StatsSummary.astro

### 2. Stats Summary Component
**File**: `frontend/src/components/StatsSummary.astro`
**Purpose**: Interactive island for stats data fetching and display
**Type**: New file
**Dependencies**: None (vanilla JavaScript)

## Files to Modify

None. All functionality is contained in new files that integrate with existing components via established patterns.

## Component Architecture

### StatsSummary.astro - Interactive Island
```astro
---
// Component frontmatter (TypeScript)
export interface StatsData {
  total_links: number;
  total_tags: number;
  total_views: number;
  top_tags: Array<{
    name: string;
    slug: string;
    link_count: number;
  }>;
  most_viewed: Array<{
    id: string;
    title: string;
    url: string;
    view_count: number;
  }>;
}

export interface StatsState {
  loading: boolean;
  error: string | null;
  data: StatsData | null;
}
---

<!-- Static HTML template with loading state -->
<div id="stats-container">
  <!-- Summary cards grid -->
  <!-- Top tags section -->
  <!-- Most viewed section -->
</div>

<script>
  // Client-side JavaScript for data fetching and DOM updates
</script>

<style>
  /* Component-scoped styles */
</style>
```

### stats.astro - Main Page
```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import StatsSummary from '../components/StatsSummary.astro';

const title = "Statistics Overview";
const description = "View aggregate statistics for links, tags, and views.";
---

<BaseLayout title={title} description={description}>
  <h1>Statistics Overview</h1>
  <StatsSummary />
</BaseLayout>
```

## Public Interfaces

### StatsSummary Component Interface
**External Interface**:
- No props required (self-contained data fetching)
- Exposes refresh capability via custom events (optional)
- Provides loading, error, and success visual states

**Internal Interface**:
```typescript
interface StatsController {
  // State management
  state: StatsState;
  setState(newState: Partial<StatsState>): void;

  // Data operations
  fetchStats(): Promise<void>;
  refreshStats(): Promise<void>;

  // DOM operations
  renderSummaryCards(data: StatsData): void;
  renderTopTags(tags: StatsData['top_tags']): void;
  renderMostViewed(links: StatsData['most_viewed']): void;
  renderLoadingState(): void;
  renderErrorState(error: string): void;
}
```

### API Integration Interface
**Endpoint**: `GET /api/stats`
**Expected Response**:
```typescript
interface StatsApiResponse {
  total_links: number;
  total_tags: number;
  total_views: number;
  top_tags: Array<{
    name: string;
    slug: string;
    link_count: number;
  }>;
  most_viewed: Array<{
    id: string;
    title: string;
    url: string;
    view_count: number;
  }>;
}
```

**Error Handling**:
- Network errors: Retry mechanism
- HTTP errors: Display user-friendly message
- Parse errors: Graceful degradation

## Module Boundaries

### Component Boundary: StatsSummary
**Responsibilities**:
- API data fetching and caching
- State management for loading/error/data states
- DOM manipulation for dynamic content updates
- User interaction handling (refresh)

**Dependencies**:
- Native fetch API
- Native DOM APIs
- CSS custom properties from BaseLayout

**Isolation**:
- Self-contained: no external state management
- No dependencies on other custom components
- Communicates only with backend API

### Page Boundary: stats.astro
**Responsibilities**:
- Route handling for `/stats`
- Page metadata (title, description)
- Component composition and layout

**Dependencies**:
- BaseLayout (existing)
- StatsSummary (new)

## CSS Architecture

### Custom Property Additions
```css
/* Add to BaseLayout.astro global styles */
:root {
  /* Stats-specific spacing */
  --stats-card-padding: var(--space-lg);
  --stats-card-gap: var(--space-md);

  /* Stats-specific sizing */
  --stats-number-size: 2.5rem;
  --stats-label-size: 0.875rem;
  --stats-card-radius: 8px;

  /* Stats-specific colors */
  --stats-card-bg: var(--color-bg);
  --stats-card-border: var(--color-border);
  --stats-card-shadow: rgba(0, 0, 0, 0.1);

  /* Loading animation */
  --skeleton-base: #f0f0f0;
  --skeleton-highlight: #e0e0e0;
}
```

### Component-Scoped Styles
**StatsSummary.astro styles**:
- Grid layouts for responsive summary cards
- Flex layouts for ranked lists
- Loading skeleton animations
- Error state styling
- Hover and focus states

### Style Hierarchy
```
BaseLayout (global scope)
├─ CSS custom properties
├─ Typography base styles
└─ Container spacing rules

StatsSummary (scoped styles)
├─ Component-specific layouts
├─ Interactive state styles
└─ Animation definitions
```

## DOM Structure

### Static HTML Template (SSR)
```html
<div id="stats-container">
  <div class="stats-summary-cards">
    <div class="stats-card loading"><!-- Total Links --></div>
    <div class="stats-card loading"><!-- Total Tags --></div>
    <div class="stats-card loading"><!-- Total Views --></div>
  </div>

  <div class="stats-lists">
    <section class="stats-section">
      <h2>Top Tags</h2>
      <div class="stats-list loading"><!-- Top tags list --></div>
    </section>

    <section class="stats-section">
      <h2>Most Viewed</h2>
      <div class="stats-list loading"><!-- Most viewed list --></div>
    </section>
  </div>

  <div class="stats-actions">
    <button id="refresh-stats" class="stats-refresh-btn">Refresh</button>
  </div>
</div>
```

### Dynamic Content Updates (CSR)
```html
<!-- After successful data fetch -->
<div class="stats-card">
  <span class="stats-number">127</span>
  <span class="stats-label">Total Links</span>
</div>

<ol class="ranked-list">
  <li class="ranked-item">
    <span class="rank">1</span>
    <span class="name">JavaScript</span>
    <span class="count">23 links</span>
  </li>
  <!-- ... more items -->
</ol>
```

## JavaScript Module Organization

### StatsSummary Client Script Structure
```javascript
// State management
class StatsState {
  constructor() {
    this.loading = true;
    this.error = null;
    this.data = null;
  }
}

// API integration
class StatsAPI {
  static async fetchStats() {
    // Fetch implementation with error handling
  }
}

// DOM manipulation
class StatsRenderer {
  static renderSummaryCards(data) { /* ... */ }
  static renderTopTags(tags) { /* ... */ }
  static renderMostViewed(links) { /* ... */ }
  static renderLoadingState() { /* ... */ }
  static renderErrorState(error) { /* ... */ }
}

// Main controller
class StatsController {
  constructor() {
    this.state = new StatsState();
    this.init();
  }

  async init() {
    await this.fetchStats();
    this.setupEventListeners();
  }
}

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
  new StatsController();
});
```

## Integration Points

### BaseLayout Integration
**Method**: Import and use BaseLayout as wrapper
**Props Passed**: title, description
**CSS Integration**: Inherit global custom properties and typography

### Navigation Integration
**Method**: No changes needed
**Reason**: Navigation.astro already includes `/stats` link from T-004-01

### API Integration
**Method**: Client-side fetch to existing endpoint
**URL**: `${window.location.origin}/api/stats`
**Headers**: Accept: application/json
**Error Handling**: Network, HTTP, and parse error scenarios

## Build Integration

### Static Build Impact
**New Assets**:
- `dist/stats/index.html` (pre-rendered page)
- Inline JavaScript in HTML (no separate bundle)
- Inline CSS in HTML (scoped styles)

**Bundle Size Impact**:
- Estimated JavaScript: ~2KB (vanilla implementation)
- Estimated CSS: ~1KB (component styles)
- No external dependencies added

### Development Integration
**Dev Server**: Standard Astro dev server handles new route
**Hot Reload**: Component changes trigger partial page updates
**TypeScript**: Full type checking for component interfaces

## Testing Strategy

### Component Testing
**Unit Tests**: Not required for this implementation phase
**Manual Testing**: Browser-based verification of:
- Static HTML rendering
- Client-side data fetching
- Loading and error states
- Responsive layout behavior

### Integration Testing
**API Integration**: Verify fetch calls to `/api/stats`
**Navigation Integration**: Verify link navigation works
**Layout Integration**: Verify BaseLayout composition

## Deployment Considerations

### Static Build Output
**Page Route**: `dist/stats/index.html`
**Assets**: All styles and scripts inlined
**API Dependency**: Requires backend server for client-side API calls

### Runtime Dependencies
**Client-side**: Native web APIs only (fetch, DOM)
**Server-side**: Existing backend API at `/api/stats`

## Implementation Order

### Phase 1: Static Structure
1. Create `frontend/src/pages/stats.astro` with BaseLayout
2. Create `frontend/src/components/StatsSummary.astro` with static HTML
3. Add CSS for layout and loading states
4. Verify static rendering and responsive layout

### Phase 2: Client-side Functionality
5. Add JavaScript for API fetching
6. Implement DOM updates for data display
7. Add error handling and retry logic
8. Add refresh button functionality

### Phase 3: Polish and Testing
9. Add loading animations and micro-interactions
10. Test error scenarios and edge cases
11. Verify browser compatibility
12. Performance validation and optimization

This structure provides clear boundaries, minimal complexity, and maintainable code organization while meeting all acceptance criteria.