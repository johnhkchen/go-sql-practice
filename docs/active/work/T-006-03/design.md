# T-006-03: Sync Admin Page Design

## Overview

Design for an Astro admin page at `/sync/[id]/control` that provides real-time control over sync sessions via a progress slider and displays shareable viewer URLs.

## Context from Research

### Dependencies Analysis
- **T-004-01 (BaseLayout)**: Provides the foundational layout (`BaseLayout.astro`) with consistent styling, navigation, and responsive design
- **T-006-02 (Sync API Routes)**: Implements the backend API endpoints (`POST /api/sync/:id/progress`) with token-based authentication

### Existing Patterns
- **Astro SSR**: Dynamic routes use `export const prerender = false` with server-side data fetching
- **Error Handling**: Comprehensive error states (network, timeout, server, notfound) with user-friendly messages
- **Client-Side Enhancement**: Script tags provide progressive enhancement, not client islands (`client:load`)
- **CSS Architecture**: Scoped styles with CSS custom properties, consistent spacing/typography
- **Accessibility**: ARIA attributes, screen reader support, keyboard navigation

### API Integration
- Sync session API at `/api/sync/:id/progress` requires admin token validation
- Token passed via query parameter: `?token=<admin_token>`
- Progress validation: float between 0.0 and 1.0
- Returns 403 for invalid tokens, 404 for missing sessions, 400 for invalid progress

## Design Options

### Option 1: Pure Client-Side Component with `client:load`

**Approach**: Create an Astro client island that handles all interactivity.

```astro
<SyncController client:load sessionId={id} token={token} />
```

**Pros**:
- Clean separation between static content and dynamic behavior
- Follows React/Vue-like patterns familiar to many developers
- Built-in client-side state management
- Astro handles hydration automatically

**Cons**:
- **Breaks existing codebase patterns**: No other components use `client:load`
- Requires framework dependency (React/Vue/Svelte)
- Larger bundle size due to framework overhead
- More complex build configuration
- Hydration mismatch risks

**Technical Implementation**:
- Create `SyncController.tsx` (React) or `SyncController.svelte`
- Manage slider state, debouncing, API calls internally
- Handle token validation and error states
- Requires framework-specific dependencies

### Option 2: Astro with Progressive Enhancement Script

**Approach**: Server-rendered Astro page with embedded `<script>` tag for interactivity.

```astro
<BaseLayout title="Sync Control">
  <!-- Server-rendered initial state -->
  <div id="sync-control">...</div>
</BaseLayout>
<script>
  // Progressive enhancement JavaScript
</script>
```

**Pros**:
- **Consistent with existing patterns**: Matches `StatsSummary.astro` and `[id].astro` approaches
- No framework dependencies
- Smaller bundle size
- Works without JavaScript (graceful degradation)
- Server-side error handling for initial load

**Cons**:
- More verbose JavaScript for state management
- Manual DOM manipulation required
- Need to handle component lifecycle manually

**Technical Implementation**:
- Server-side initial data fetch with error handling
- TypeScript class for client-side controller
- Manual debouncing implementation
- Direct fetch API usage for progress updates

### Option 3: Hybrid SSR + Client Island for Slider Only

**Approach**: SSR page with minimal client island just for the slider component.

```astro
<BaseLayout>
  <!-- Static content server-rendered -->
  <ProgressSlider client:load sessionId={id} token={token} initialProgress={progress} />
  <!-- Rest of page static -->
</BaseLayout>
```

**Pros**:
- Minimizes client-side JavaScript
- Keeps static content fast-loading
- Targeted interactivity

**Cons**:
- Still introduces framework dependency
- Splits logic between server and client
- Communication complexity between static and dynamic parts

## Evaluation Criteria

### 1. Consistency with Codebase
- **Option 1**: ❌ Introduces new patterns not used elsewhere
- **Option 2**: ✅ Follows existing `StatsSummary.astro` and `[id].astro` patterns
- **Option 3**: ❌ Partial inconsistency with framework usage

### 2. Performance
- **Option 1**: ❌ Framework overhead, larger bundle, hydration cost
- **Option 2**: ✅ Minimal JavaScript, fast initial load, progressive enhancement
- **Option 3**: ⚠️ Smaller overhead than Option 1, but still has framework cost

### 3. Maintainability
- **Option 1**: ✅ Familiar component patterns for React/Vue developers
- **Option 2**: ⚠️ More verbose but follows existing project patterns
- **Option 3**: ❌ Split complexity between two paradigms

### 4. Technical Requirements Alignment
- **Option 1**: ✅ Can handle all requirements but with framework overhead
- **Option 2**: ✅ Can handle all requirements with existing tech stack
- **Option 3**: ✅ Can handle requirements but adds complexity

### 5. Bundle Size & Dependencies
- **Option 1**: ❌ Requires React/Vue/Svelte framework addition
- **Option 2**: ✅ No new dependencies, vanilla JavaScript
- **Option 3**: ❌ Requires framework for slider only

## Recommended Solution: Option 2 - Progressive Enhancement

### Rationale

**Primary Factors**:
1. **Codebase Consistency**: Matches existing patterns in `StatsSummary.astro` and watch pages
2. **No Framework Lock-in**: Stays within current vanilla JavaScript approach
3. **Performance**: Minimal JavaScript overhead, fast loading
4. **Graceful Degradation**: Works without JavaScript for basic functionality

**Technical Architecture**:

```astro
---
// Server-side: fetch initial session data, handle auth, errors
const { id } = Astro.params;
const token = Astro.url.searchParams.get('token');
// Initial data fetch and validation logic
---

<BaseLayout title="Sync Control - {session.name}">
  <!-- Static content: session info, viewer URL, copy button -->
  <div class="sync-control-container">
    <div class="session-info">...</div>
    <div class="progress-control" id="progress-control">
      <!-- Server-rendered slider with initial value -->
      <input type="range" id="progress-slider" value={progress} />
      <span id="progress-display">{progress}</span>
    </div>
    <div class="viewer-url-section">...</div>
  </div>
</BaseLayout>

<script>
  // TypeScript class following StatsSummary pattern
  class SyncController { ... }
</script>
```

### Key Implementation Details

1. **Server-Side Rendering**: Initial session validation, token verification, error states
2. **Progressive Enhancement**: JavaScript enhances the experience but page works without it
3. **Debounced Updates**: Throttle slider changes to max ~30 updates/second
4. **Error Handling**: Network failures, token expiration, session not found
5. **Accessibility**: ARIA labels, keyboard navigation, screen reader announcements
6. **Responsive Design**: Works on mobile and desktop using existing CSS patterns

### Rejected Alternatives

- **Option 1** rejected due to framework dependency and pattern inconsistency
- **Option 3** rejected due to complexity split and partial framework introduction

## Technical Considerations

### Slider Debouncing Strategy
```javascript
// Throttle approach: max frequency limiting
const throttle = (func, delay) => {
  let timeoutId;
  let lastExecTime = 0;
  return (...args) => {
    const currentTime = Date.now();
    if (currentTime - lastExecTime > delay) {
      func.apply(this, args);
      lastExecTime = currentTime;
    } else {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => {
        func.apply(this, args);
        lastExecTime = Date.now();
      }, delay - (currentTime - lastExecTime));
    }
  };
};
```

### Token Management
- Store token in page data for JavaScript access
- Implement token refresh mechanism if needed
- Handle 403 responses gracefully with user feedback

### URL Structure Validation
- Route: `/sync/[id]/control?token=<admin_token>`
- Validate session ID format on server-side
- Validate token presence and format
- Generate proper viewer URLs: `/sync/[id]`

### Copy-to-Clipboard Implementation
- Feature detection for Clipboard API
- Fallback to text selection + document.execCommand
- Visual feedback for successful/failed copy operations
- Accessibility announcements for screen readers

This design provides a robust, performant, and maintainable solution that aligns with the existing codebase architecture while meeting all the specified requirements.