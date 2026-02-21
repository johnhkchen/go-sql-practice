# Design: Link Detail Page Implementation (T-004-03)

## Overview

This design document outlines the implementation approach for the link detail page at `/links/[id]`, building on the established patterns identified in the research phase. The page will display comprehensive link information, handle view count incrementation, and provide robust error handling following existing architectural patterns.

## Viable Implementation Approaches

### Approach 1: Direct PocketBase Collection Access (Recommended)

**Architecture**: Follow the established pattern from `/watch/[id].astro` with direct PocketBase API calls
- **Data Fetching**: Use `GET /api/collections/links/records/:id` for link data during SSR
- **View Increment**: Client-side `POST /api/links/:id/view` call on page load
- **Tag Resolution**: Single additional API call to resolve tag IDs to slugs
- **State Management**: Component-level state for view count updates

**Advantages**:
- Leverages proven patterns from watch page implementation
- Uses existing PocketBase auto-generated endpoints
- Clear separation between SSR data fetching and client-side interactions
- Minimal custom backend changes required

**Implementation Strategy**:
1. SSR phase: Fetch link record and handle errors identically to watch page
2. Client-side enhancement: Call view increment endpoint and update DOM
3. Tag resolution: Additional fetch to resolve tag IDs to display slugs
4. Progressive enhancement: Page works without JavaScript, enhanced with it

### Approach 2: Custom Single Endpoint

**Architecture**: Create dedicated `/api/links/:id/detail` endpoint combining all data needs
- **Data Fetching**: Single custom endpoint returning link + resolved tags
- **View Increment**: Embedded in detail endpoint or separate call
- **State Management**: Simplified since all data comes pre-resolved

**Advantages**:
- Single network request for complete data
- Server-side tag resolution eliminates client-side complexity
- Potentially better performance due to reduced API calls

**Disadvantages**:
- Requires new backend endpoint development
- Deviates from established PocketBase patterns
- Additional testing surface area
- Breaks consistency with existing dynamic page patterns

### Approach 3: Search Endpoint Repurposing

**Architecture**: Use existing `/api/links/search?id=:id` functionality
- **Data Fetching**: Leverage search endpoint with ID filter
- **Tag Resolution**: Built-in via search endpoint's tag fetching logic
- **View Increment**: Separate client-side call

**Advantages**:
- Reuses existing complex tag resolution logic
- Tags come pre-resolved as slugs
- No new backend development

**Disadvantages**:
- Semantic mismatch (search for detail view)
- Potential performance overhead from search complexity
- May return unexpected data structure for single records
- Search endpoint optimized for multiple results, not single record

## Assessment Against Codebase Reality

### Pattern Consistency Analysis

**Approach 1 Alignment**:
- ✅ **Dynamic Route Structure**: Perfectly matches `/watch/[id].astro` pattern with `export const prerender = false`
- ✅ **Error Handling**: Can reuse identical error categorization and messaging system
- ✅ **API Communication**: Follows established timeout, AbortController, and fetch patterns
- ✅ **Layout Integration**: Uses existing `BaseLayout.astro` consistently
- ✅ **State Management**: Aligns with component-level state approach from `StatsSummary.astro`

**Tag Resolution Reality**:
- Research shows links store tag relations as JSON arrays of tag IDs
- Search endpoints use `fetchTagsForLinks()` with complex JOIN logic
- Direct collection access returns raw tag IDs, requiring resolution
- Client-side resolution adds complexity but maintains architectural consistency

**View Count Integration**:
- `/routes/links_view.go` endpoint is production-ready and atomic
- Returns full updated record including new view count
- Perfect for optimistic UI updates with actual data confirmation

### Performance Considerations

**Network Requests**:
- Approach 1: 2-3 requests (link data + view increment + optional tag resolution)
- Approach 2: 1-2 requests (unified endpoint + optional view increment)
- Research shows existing 5-second timeout provides adequate buffer

**Data Loading Strategy**:
- SSR provides immediate content for SEO and accessibility
- Client-side enhancement handles view tracking without blocking render
- Tag resolution can be progressive (show IDs, then enhance with names)

## Selected Approach: Direct PocketBase Collection Access

### Rationale

**Primary Factors**:
1. **Pattern Consistency**: Exactly matches established `/watch/[id].astro` implementation
2. **Risk Minimization**: Uses proven, stable patterns with existing error handling
3. **Development Speed**: Minimal new backend code required
4. **Maintainability**: Future developers can easily understand and modify

**Architecture Decision**:
- Follow watch page SSR pattern for link data fetching
- Use existing view increment endpoint for atomic view counting
- Implement client-side tag resolution as progressive enhancement
- Maintain component-level state management consistent with `StatsSummary.astro`

### Tag Resolution Strategy

**Three-Phase Approach**:
1. **Initial Render**: Display link with tag IDs as fallback pills
2. **Tag Fetch**: Client-side resolution of tag IDs to `{id, name, slug}` objects
3. **Enhancement**: Replace tag ID pills with properly linked tag name pills

**API Call Strategy**:
```javascript
// Fetch individual tag details for resolution
const tagPromises = link.tags.map(tagId =>
  fetch(`/api/collections/tags/records/${tagId}`)
);
const tagResponses = await Promise.allSettled(tagPromises);
```

**Progressive Enhancement Benefits**:
- Page functional without JavaScript (shows tag IDs)
- Graceful degradation if tag resolution fails
- Follows established accessibility patterns

### Error Handling Architecture

**Reuse Existing Patterns**:
- Identical error categorization: `notfound`, `network`, `server`, `timeout`
- Same error message structure and styling
- Consistent back-to-home navigation pattern

**Link-Specific Adaptations**:
```javascript
const errorMessages = {
  notfound: {
    title: 'Link Not Found',
    message: 'This link does not exist or has been removed.'
  },
  // ... other error types remain identical
};
```

### View Count Update Strategy

**Atomic Operation Flow**:
1. Page loads with SSR-fetched link data (including current view count)
2. Client-side script calls `POST /api/links/:id/view` on load
3. Endpoint atomically increments and returns updated record
4. UI updates with new view count from response

**Implementation Pattern**:
```javascript
// Based on StatsSummary.astro class-based approach
class ViewCountUpdater {
  async incrementAndUpdate() {
    try {
      const response = await fetch(`/api/links/${linkId}/view`, {
        method: 'POST',
        signal: AbortController.signal // 5s timeout
      });
      const updatedLink = await response.json();
      this.updateViewCountDisplay(updatedLink.view_count);
    } catch (error) {
      // Graceful failure - page still functional
    }
  }
}
```

## Rejected Approaches and Rationale

### Custom Unified Endpoint (Approach 2)

**Rejection Reasons**:
- **Development Overhead**: Requires new backend endpoint, testing, and documentation
- **Pattern Divergence**: Creates inconsistency with established dynamic route patterns
- **Over-Engineering**: Solves performance problem that doesn't exist (2-3 API calls acceptable)
- **Maintenance Burden**: Additional custom code to maintain vs. leveraging PocketBase conventions

### Search Endpoint Repurposing (Approach 3)

**Rejection Reasons**:
- **Semantic Mismatch**: Search endpoint designed for queries, not single record retrieval
- **Complexity Overhead**: Search includes pagination, filtering logic unnecessary for detail view
- **API Misuse**: Using search functionality for non-search purposes reduces code clarity
- **Future Brittleness**: Search endpoint may evolve independently, breaking detail page

## Implementation Structure

### Component Architecture

**File Structure**:
```
frontend/src/pages/links/[id].astro  # Main dynamic route file
```

**Component Breakdown**:
1. **SSR Section**: Link data fetching, error handling, initial render prep
2. **Template Section**: Conditional rendering based on data/error states
3. **Style Section**: Link-specific styling following existing CSS custom property patterns
4. **Script Section**: Client-side view counting and tag resolution enhancement

### Data Flow Design

**SSR Phase**:
1. Extract `id` from `Astro.params`
2. Fetch link data via `GET /api/collections/links/records/${id}`
3. Handle fetch errors using established error categorization
4. Prepare data for template rendering

**Client-Side Enhancement Phase**:
1. Initialize view count updater on page load
2. Call view increment endpoint asynchronously
3. Resolve tag IDs to tag objects for proper pill display
4. Update DOM with enhanced data

### Error State Management

**SSR Error Handling**:
- Reuse exact error detection and categorization from watch page
- Adapt error messages for link context
- Maintain consistent error container styling and navigation

**Client-Side Error Handling**:
- View count increment failures should be silent (page remains functional)
- Tag resolution failures show tag IDs as fallback
- Network errors don't break core page functionality

## Key Design Decisions

### 1. SSR vs Client-Side Data Strategy

**Decision**: SSR for link data, client-side for enhancements
**Reasoning**: Matches established pattern, provides immediate content, enables progressive enhancement

### 2. Tag Resolution Approach

**Decision**: Client-side progressive enhancement
**Reasoning**: Maintains architectural consistency, provides graceful degradation, acceptable performance

### 3. State Management Pattern

**Decision**: Component-level state following `StatsSummary.astro` class approach
**Reasoning**: Consistent with existing patterns, no global state needed for single-page interactions

### 4. View Count Update Timing

**Decision**: Post-load client-side increment
**Reasoning**: Non-blocking, uses existing atomic endpoint, provides visual feedback

### 5. Error Handling Consistency

**Decision**: Exact replication of watch page error patterns
**Reasoning**: User experience consistency, proven error handling, reduced development complexity

## Implementation Accelerators

### Reusable Patterns

**From `/watch/[id].astro`**:
- Complete SSR data fetching pattern with timeout/abort handling
- Error state categorization and messaging system
- Page title dynamic generation
- Responsive container styling patterns

**From `StatsSummary.astro`**:
- Class-based client-side state management
- Progressive DOM enhancement techniques
- Error state management for API calls
- Accessibility considerations (aria-live regions)

**From Existing Styles**:
- CSS custom property system for theming
- Animation utilities from `animations.css`
- Mobile-responsive patterns
- Accessibility features (prefers-reduced-motion)

### Component Templates

**Link Detail Display Structure**:
```astro
<div class="link-detail-container">
  <header class="link-header">
    <h1 class="link-title">{link.title}</h1>
    <div class="link-meta">
      <span class="view-count" data-count="{link.view_count}">
        {link.view_count} views
      </span>
    </div>
  </header>

  <div class="link-content">
    <div class="link-url">
      <a href="{link.url}" target="_blank" rel="noopener noreferrer">
        {link.url}
      </a>
    </div>

    <div class="link-description">
      {link.description}
    </div>

    <div class="link-tags">
      {link.tags.map(tagId => (
        <span class="tag-pill" data-tag-id="{tagId}">{tagId}</span>
      ))}
    </div>
  </div>
</div>
```

## Risk Mitigation

### Identified Risks from Research

**API Endpoint Availability**:
- **Risk**: PocketBase collection endpoint might not be publicly accessible
- **Mitigation**: Research confirms public read access (`ListRule: ""`, `ViewRule: ""`)

**Tag Resolution Complexity**:
- **Risk**: Client-side tag resolution adds multiple API calls and failure points
- **Mitigation**: Progressive enhancement ensures page functional without tag names, graceful fallback to tag IDs

**View Count Race Conditions**:
- **Risk**: Concurrent view count updates could cause inconsistencies
- **Mitigation**: Existing `/routes/links_view.go` uses atomic SQL UPDATE, proven safe

**Error Handling Coverage**:
- **Risk**: New error scenarios not covered by existing patterns
- **Mitigation**: Reuse exact error categorization from watch page, covers network, timeout, 404, and server errors

## Success Metrics

**Functional Requirements**:
- ✅ Dynamic route renders single link details correctly
- ✅ View count increments atomically on page load
- ✅ All required fields display: title, URL, description, tags, view count
- ✅ Tags render as pills with links to `/tags/:slug`
- ✅ 404 handling for non-existent links
- ✅ View count updates visually after increment

**Implementation Quality**:
- ✅ Follows established architectural patterns exactly
- ✅ Reuses existing error handling and styling systems
- ✅ Maintains performance characteristics of similar pages
- ✅ Provides progressive enhancement without JavaScript dependency
- ✅ Maintains accessibility standards from existing components

**Integration Success**:
- ✅ Uses existing view increment endpoint without modification
- ✅ Leverages PocketBase auto-generated endpoints
- ✅ Consistent with `BaseLayout.astro` and navigation patterns
- ✅ Compatible with established CSS custom property system

This design provides a clear, research-backed implementation path that leverages existing patterns while meeting all functional requirements with minimal risk and development overhead.