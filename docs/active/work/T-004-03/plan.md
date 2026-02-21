# Implementation Plan: Link Detail Page (T-004-03)

## Overview

This plan sequences the implementation of `/frontend/src/pages/links/[id].astro` following the established patterns from `/watch/[id].astro`. The approach uses direct PocketBase collection access for SSR with client-side enhancement for view counting and tag resolution.

## Implementation Strategy

### Core Approach
- **SSR First**: Fetch link data during server-side rendering for immediate content
- **Progressive Enhancement**: Client-side view counting and tag resolution enhance the experience
- **Pattern Replication**: Follow `/watch/[id].astro` patterns exactly for consistency
- **Error Handling**: Reuse established error categorization and user experience patterns

### Key Dependencies (Already Available)
- ✅ Backend API: `POST /api/links/:id/view` endpoint (T-003-02 completed)
- ✅ Frontend Layout: `BaseLayout.astro` (T-004-01 completed)
- ✅ PocketBase Collections: Auto-generated endpoints for links and tags
- ✅ CSS System: Custom properties and animation utilities

## Step-by-Step Implementation Plan

### Step 1: Create Basic SSR Structure and Dynamic Routing
**Duration**: 30-45 minutes
**Complexity**: Low
**Dependencies**: None

#### Implementation Tasks
1. Create `/frontend/src/pages/links/[id].astro` file
2. Add SSR section with basic parameter extraction and validation
3. Implement redirect to `/404` for missing parameters
4. Add basic TypeScript interfaces for Link and Error states
5. Set up API base URL configuration with environment variable support

#### Code Structure
```astro
---
import BaseLayout from '../../layouts/BaseLayout.astro'

// Type definitions
interface LinkItem {
  id: string
  url: string
  title: string
  description: string
  view_count: number
  tags: string[]
  created: string
  updated: string
}

interface ErrorState {
  type: 'notfound' | 'network' | 'server' | 'timeout' | null
  title?: string
  message?: string
}

// Parameter validation
const { id } = Astro.params
if (!id) {
  return Astro.redirect('/404')
}

// API configuration
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090'

// Placeholder for data fetching (Step 2)
let link: LinkItem | null = null
let error: ErrorState | null = null

const pageTitle = 'Link Details'
---

<BaseLayout title={pageTitle}>
  <div class="page-container">
    <h1>Link Detail Page - Step 1 Complete</h1>
    <p>Link ID: {id}</p>
  </div>
</BaseLayout>
```

#### Acceptance Criteria
- [ ] File creates successfully at correct path
- [ ] Route responds at `/links/123` with placeholder content
- [ ] Invalid URLs (no ID) redirect to `/404`
- [ ] TypeScript interfaces defined correctly
- [ ] API_BASE configuration works with environment variables

#### Testing Approach
- **Manual**: Navigate to `/links/test-id` and verify placeholder displays
- **Manual**: Navigate to `/links/` (no ID) and verify redirect to `/404`
- **Manual**: Check browser developer tools for TypeScript compilation errors

---

### Step 2: Implement Link Data Fetching and Error Handling
**Duration**: 45-60 minutes
**Complexity**: Medium
**Dependencies**: Step 1 complete

#### Implementation Tasks
1. Add comprehensive fetch logic with AbortController timeout (5 seconds)
2. Implement error categorization matching `/watch/[id].astro` patterns exactly
3. Add error message definitions for link-specific contexts
4. Handle successful response parsing and validation
5. Update page title generation based on fetched data

#### Code Changes
Replace placeholder data fetching with:
```astro
const TIMEOUT_MS = 5000

// Error message definitions
const errorMessages = {
  notfound: {
    title: 'Link Not Found',
    message: 'This link does not exist or has been removed. Please check the URL and try again.'
  },
  network: {
    title: 'Connection Failed',
    message: 'Unable to connect to the server. Please check your internet connection and try again.'
  },
  server: {
    title: 'Server Error',
    message: 'The server encountered an error. Please try again in a few moments.'
  },
  timeout: {
    title: 'Request Timeout',
    message: 'The request took too long to complete. Please try again.'
  }
}

// Data fetching with comprehensive error handling
try {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), TIMEOUT_MS)

  const response = await fetch(`${API_BASE}/api/collections/links/records/${id}`, {
    signal: controller.signal
  })

  clearTimeout(timeoutId)

  if (response.status === 404) {
    error = { type: 'notfound', ...errorMessages.notfound }
  } else if (!response.ok) {
    if (response.status >= 500) {
      error = { type: 'server', ...errorMessages.server }
    } else {
      error = { type: 'network', ...errorMessages.network }
    }
  } else {
    const fetchedLink = await response.json()
    // Basic validation
    if (fetchedLink.id && fetchedLink.url && fetchedLink.title) {
      link = fetchedLink
    } else {
      error = { type: 'server', ...errorMessages.server }
    }
  }
} catch (fetchError: any) {
  if (fetchError.name === 'AbortError') {
    error = { type: 'timeout', ...errorMessages.timeout }
  } else {
    error = { type: 'network', ...errorMessages.network }
  }
}

// Update page title
const pageTitle = link ? `${link.title} - Link Details` : 'Link Not Found'
```

#### Acceptance Criteria
- [ ] Valid link IDs load link data correctly
- [ ] Invalid link IDs show appropriate 404 error
- [ ] Network failures show network error with retry guidance
- [ ] Server errors (5xx) show server error messages
- [ ] Requests timeout after 5 seconds with timeout error
- [ ] Page title updates dynamically based on link data
- [ ] All error messages match established patterns from watch page

#### Testing Approach
- **Manual**: Test with valid link ID from database
- **Manual**: Test with non-existent link ID (should show 404)
- **Manual**: Test with server disconnected (should show network error)
- **Network**: Use browser dev tools to simulate slow network (should timeout)
- **API**: Use curl to verify backend endpoint responses

#### Debugging Strategy
- Check browser console for fetch errors
- Verify API_BASE URL matches running backend
- Check PocketBase admin panel for existing link records
- Use browser network tab to inspect HTTP requests/responses

---

### Step 3: Add Template Structure and Basic Styling
**Duration**: 60-75 minutes
**Complexity**: Medium
**Dependencies**: Step 2 complete

#### Implementation Tasks
1. Replace placeholder template with conditional rendering for error/success states
2. Add complete link detail template structure
3. Implement error container styling identical to watch page
4. Add comprehensive CSS using established custom property system
5. Implement responsive design patterns and accessibility features

#### Template Structure
```astro
<BaseLayout title={pageTitle}>
  <div class="page-container">
    {error ? (
      <!-- Error state (identical to watch page) -->
      <div class="error-container">
        <div class="error-content">
          <h1 class="error-title">{error.title}</h1>
          <p class="error-message">{error.message}</p>
          <div class="error-actions">
            <a href="/" class="home-link">← Back to Home</a>
            <button onclick="window.location.reload()" class="retry-button">Try Again</button>
          </div>
        </div>
      </div>
    ) : link && (
      <!-- Success state -->
      <article class="link-detail">
        <header class="link-header">
          <h1 class="link-title">{link.title}</h1>
          <div class="link-meta">
            <span class="view-count" data-link-id={link.id} data-current-count={link.view_count}>
              <span class="view-number">{link.view_count || 0}</span>
              <span class="view-label">views</span>
            </span>
            <time class="link-date" datetime={link.created}>
              {new Date(link.created).toLocaleDateString()}
            </time>
          </div>
        </header>

        <div class="link-content">
          <div class="link-url-section">
            <label class="section-label">URL:</label>
            <a href={link.url}
               class="link-url"
               target="_blank"
               rel="noopener noreferrer">
              {link.url}
            </a>
          </div>

          {link.description && (
            <div class="link-description-section">
              <label class="section-label">Description:</label>
              <p class="link-description">{link.description}</p>
            </div>
          )}

          <div class="link-tags-section">
            <label class="section-label">Tags:</label>
            <div class="tags-container" data-link-id={link.id}>
              {link.tags && link.tags.length > 0 ? (
                link.tags.map(tagId => (
                  <span class="tag-pill loading" data-tag-id={tagId}>
                    <span class="tag-content">{tagId}</span>
                    <span class="tag-loading" aria-hidden="true">…</span>
                  </span>
                ))
              ) : (
                <span class="no-tags">No tags assigned</span>
              )}
            </div>
          </div>
        </div>
      </article>
    )}
  </div>
</BaseLayout>
```

#### Styling Implementation
Add comprehensive CSS section using custom properties:
```astro
<style>
  /* Page layout */
  .page-container {
    max-width: var(--content-max-width, 800px);
    margin: 0 auto;
    padding: var(--spacing-md, 1rem);
  }

  /* Error states (identical to watch page) */
  .error-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 400px;
    text-align: center;
  }

  .error-content {
    max-width: 500px;
    padding: var(--spacing-xl, 2rem);
    background: var(--surface-color, white);
    border-radius: var(--radius-lg, 12px);
    box-shadow: var(--shadow-md, 0 4px 6px rgba(0,0,0,0.1));
  }

  .error-title {
    font-size: var(--text-xl, 1.25rem);
    font-weight: 600;
    margin: 0 0 var(--spacing-md, 1rem) 0;
    color: var(--text-primary, #111827);
  }

  .error-message {
    color: var(--text-secondary, #6b7280);
    margin: 0 0 var(--spacing-lg, 1.5rem) 0;
    line-height: 1.5;
  }

  .error-actions {
    display: flex;
    gap: var(--spacing-md, 1rem);
    justify-content: center;
    flex-wrap: wrap;
  }

  .home-link,
  .retry-button {
    padding: var(--spacing-sm, 0.5rem) var(--spacing-md, 1rem);
    border-radius: var(--radius-md, 6px);
    text-decoration: none;
    font-weight: 500;
    transition: all 0.2s;
  }

  .home-link {
    background: var(--primary-color, #3b82f6);
    color: white;
  }

  .retry-button {
    background: var(--surface-color, white);
    color: var(--text-primary, #111827);
    border: 1px solid var(--border-color, #d1d5db);
    cursor: pointer;
  }

  /* Link detail layout */
  .link-detail {
    background: var(--surface-color, white);
    border-radius: var(--radius-lg, 12px);
    padding: var(--spacing-xl, 2rem);
    box-shadow: var(--shadow-sm, 0 2px 4px rgba(0,0,0,0.1));
    border: 1px solid var(--border-color, #e5e7eb);
  }

  .link-header {
    border-bottom: 2px solid var(--border-color, #e5e7eb);
    padding-bottom: var(--spacing-lg, 1.5rem);
    margin-bottom: var(--spacing-lg, 1.5rem);
  }

  .link-title {
    font-size: var(--text-2xl, 1.5rem);
    font-weight: 700;
    margin: 0 0 var(--spacing-md, 1rem) 0;
    color: var(--text-primary, #111827);
    line-height: 1.3;
  }

  .link-meta {
    display: flex;
    gap: var(--spacing-lg, 1.5rem);
    align-items: center;
    font-size: var(--text-sm, 0.875rem);
    color: var(--text-secondary, #6b7280);
  }

  .view-count {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs, 0.25rem);
  }

  .view-number {
    font-weight: 600;
    color: var(--text-primary, #111827);
  }

  /* Content sections */
  .link-url-section,
  .link-description-section,
  .link-tags-section {
    margin-bottom: var(--spacing-lg, 1.5rem);
  }

  .section-label {
    display: block;
    font-weight: 600;
    margin-bottom: var(--spacing-sm, 0.5rem);
    color: var(--text-primary, #111827);
    font-size: var(--text-sm, 0.875rem);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .link-url {
    display: inline-block;
    color: var(--link-color, #3b82f6);
    text-decoration: underline;
    word-break: break-all;
    padding: var(--spacing-sm, 0.5rem);
    border-radius: var(--radius-md, 6px);
    transition: background-color 0.2s;
    font-weight: 500;
  }

  .link-url:hover {
    background-color: var(--link-hover-bg, #eff6ff);
    text-decoration-thickness: 2px;
  }

  .link-description {
    color: var(--text-primary, #111827);
    line-height: 1.6;
    margin: 0;
  }

  /* Tags styling */
  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-sm, 0.5rem);
  }

  .tag-pill {
    display: inline-flex;
    align-items: center;
    padding: var(--spacing-xs, 0.25rem) var(--spacing-sm, 0.5rem);
    background: var(--tag-bg, #f3f4f6);
    border: 1px solid var(--tag-border, #d1d5db);
    border-radius: var(--radius-full, 9999px);
    font-size: var(--text-xs, 0.75rem);
    font-weight: 500;
    transition: all 0.2s;
  }

  .tag-pill.loading {
    opacity: 0.7;
  }

  .tag-pill.resolved {
    background: var(--tag-resolved-bg, #dbeafe);
    border-color: var(--tag-resolved-border, #93c5fd);
    cursor: pointer;
  }

  .tag-pill.resolved:hover {
    background: var(--tag-hover-bg, #bfdbfe);
    transform: translateY(-1px);
  }

  .tag-loading {
    margin-left: var(--spacing-xs, 0.25rem);
    font-size: var(--text-xs, 0.75rem);
    opacity: 0.6;
    animation: pulse 1.5s ease-in-out infinite;
  }

  .no-tags {
    color: var(--text-secondary, #6b7280);
    font-style: italic;
  }

  /* Responsive design */
  @media (max-width: 768px) {
    .link-meta {
      flex-direction: column;
      align-items: flex-start;
      gap: var(--spacing-sm, 0.5rem);
    }

    .link-url {
      word-break: break-all;
    }

    .tags-container {
      gap: var(--spacing-xs, 0.25rem);
    }

    .error-actions {
      flex-direction: column;
    }
  }

  @media (max-width: 480px) {
    .page-container {
      padding: var(--spacing-sm, 0.5rem);
    }

    .link-detail {
      padding: var(--spacing-lg, 1.5rem);
    }
  }

  /* Accessibility */
  @media (prefers-reduced-motion: reduce) {
    .tag-loading {
      animation: none;
    }

    .tag-pill,
    .link-url,
    .home-link,
    .retry-button {
      transition: none;
    }

    .tag-pill.resolved:hover {
      transform: none;
    }
  }

  /* Focus styles */
  .link-url:focus,
  .home-link:focus,
  .retry-button:focus {
    outline: 2px solid var(--focus-color, #3b82f6);
    outline-offset: 2px;
  }
</style>
```

#### Acceptance Criteria
- [ ] Error states display correctly with consistent styling
- [ ] Link data displays in structured, readable format
- [ ] All required fields show: title, URL, description, tags, view count, creation date
- [ ] Tags display as pills with loading indicators (placeholder IDs)
- [ ] URL links open in new tab with proper security attributes
- [ ] Responsive design works on mobile and desktop
- [ ] Error page matches watch page styling exactly
- [ ] Accessibility features work (focus states, screen reader support)

#### Testing Approach
- **Manual**: Test with various link data (with/without description, with/without tags)
- **Responsive**: Test on different screen sizes using browser dev tools
- **Accessibility**: Tab through page using keyboard navigation
- **Error States**: Test all error conditions from Step 2
- **Visual**: Compare error styling with watch page for consistency

---

### Step 4: Implement Client-Side View Counter Increment
**Duration**: 30-45 minutes
**Complexity**: Low-Medium
**Dependencies**: Step 3 complete

#### Implementation Tasks
1. Add client-side script section with ViewCountUpdater class
2. Implement automatic view increment on page load
3. Add optimistic UI updates with actual response data
4. Handle errors gracefully without breaking page functionality
5. Add proper TypeScript types for client-side code

#### Script Implementation
```astro
<script>
  interface LinkRecord {
    id: string
    view_count: number
  }

  class ViewCountUpdater {
    private linkId: string
    private viewCountElement: HTMLElement | null
    private viewNumberElement: HTMLElement | null

    constructor(linkId: string) {
      this.linkId = linkId
      this.viewCountElement = document.querySelector(`[data-link-id="${linkId}"][data-current-count]`)
      this.viewNumberElement = this.viewCountElement?.querySelector('.view-number') || null
    }

    async incrementView(): Promise<void> {
      try {
        const controller = new AbortController()
        setTimeout(() => controller.abort(), 5000)

        const API_BASE = window.location.origin.includes('localhost')
          ? 'http://localhost:8090'
          : window.location.origin

        const response = await fetch(`${API_BASE}/api/links/${this.linkId}/view`, {
          method: 'POST',
          signal: controller.signal
        })

        if (response.ok) {
          const updatedLink: LinkRecord = await response.json()
          this.updateDisplay(updatedLink.view_count)
        } else {
          console.debug('View count update failed:', response.status)
        }
      } catch (error) {
        // Silent failure - page remains functional
        console.debug('View count update error:', error)
      }
    }

    private updateDisplay(newCount: number): void {
      if (this.viewNumberElement) {
        // Add subtle animation for the update
        this.viewNumberElement.style.transition = 'opacity 0.3s ease'
        this.viewNumberElement.style.opacity = '0.5'

        setTimeout(() => {
          this.viewNumberElement!.textContent = newCount.toString()
          this.viewNumberElement!.style.opacity = '1'
        }, 150)
      }
    }

    // Method to get current displayed count for testing
    getCurrentDisplayedCount(): number {
      const text = this.viewNumberElement?.textContent || '0'
      return parseInt(text, 10) || 0
    }
  }

  // Initialize on DOM ready
  document.addEventListener('DOMContentLoaded', () => {
    const linkElement = document.querySelector('[data-link-id]') as HTMLElement
    if (!linkElement) {
      console.debug('No link element found for view counting')
      return
    }

    const linkId = linkElement.getAttribute('data-link-id')
    if (!linkId) {
      console.debug('No link ID found in data attribute')
      return
    }

    // Initialize view count updater
    const viewUpdater = new ViewCountUpdater(linkId)

    // Small delay to ensure page is fully loaded before incrementing
    setTimeout(() => {
      viewUpdater.incrementView()
    }, 500)

    // Expose for testing purposes
    (window as any).__viewUpdater = viewUpdater
  })
</script>
```

#### Acceptance Criteria
- [ ] View count increments automatically when page loads
- [ ] UI updates with new count from server response
- [ ] Failed view increments don't break page functionality
- [ ] View count animation provides smooth visual feedback
- [ ] API calls use appropriate base URL (localhost vs production)
- [ ] Timeout handling prevents hanging requests
- [ ] Console logging provides debugging information
- [ ] Class methods are testable via browser console

#### Testing Approach
- **Functional**: Load page and verify view count increases by 1
- **API**: Check browser network tab for POST request to `/api/links/:id/view`
- **Error Handling**: Disconnect backend and verify page remains functional
- **Animation**: Observe view count number animation during update
- **Console Testing**: Use `window.__viewUpdater.getCurrentDisplayedCount()` to verify state
- **Database**: Verify view count increments in PocketBase admin panel

#### Performance Considerations
- 500ms delay before API call allows page to fully render
- 5-second timeout prevents hanging requests
- Silent failure ensures user experience isn't disrupted
- Minimal DOM manipulation for smooth updates

---

### Step 5: Add Tag Resolution and Linking
**Duration**: 45-60 minutes
**Complexity**: Medium-High
**Dependencies**: Step 4 complete

#### Implementation Tasks
1. Add TagResolver class for progressive tag enhancement
2. Implement parallel tag resolution with Promise.allSettled
3. Convert resolved tags to clickable links pointing to `/tags/:slug`
4. Add visual feedback for loading, resolved, and failed states
5. Handle partial resolution gracefully (mix of names and IDs)

#### Tag Resolution Implementation
Add to the script section:
```astro
<script>
  // ... existing ViewCountUpdater code ...

  interface TagRecord {
    id: string
    name: string
    slug: string
  }

  class TagResolver {
    private linkId: string
    private tagsContainer: HTMLElement | null
    private tagPills: NodeListOf<HTMLElement>

    constructor(linkId: string) {
      this.linkId = linkId
      this.tagsContainer = document.querySelector(`[data-link-id="${linkId}"] .tags-container`)
      this.tagPills = this.tagsContainer?.querySelectorAll('.tag-pill[data-tag-id]') || new NodeList() as any
    }

    async resolveAllTags(): Promise<void> {
      if (!this.tagsContainer || this.tagPills.length === 0) {
        console.debug('No tags to resolve')
        return
      }

      console.debug(`Resolving ${this.tagPills.length} tags`)

      // Start all tag resolutions in parallel
      const resolutionPromises = Array.from(this.tagPills).map(pill =>
        this.resolveTagPill(pill)
      )

      // Wait for all to complete (including failures)
      const results = await Promise.allSettled(resolutionPromises)

      const successful = results.filter(r => r.status === 'fulfilled').length
      const failed = results.filter(r => r.status === 'rejected').length

      console.debug(`Tag resolution complete: ${successful} successful, ${failed} failed`)
    }

    private async resolveTagPill(pill: HTMLElement): Promise<void> {
      const tagId = pill.dataset.tagId
      if (!tagId) {
        throw new Error('No tag ID found')
      }

      try {
        const controller = new AbortController()
        setTimeout(() => controller.abort(), 3000)

        const API_BASE = window.location.origin.includes('localhost')
          ? 'http://localhost:8090'
          : window.location.origin

        const response = await fetch(`${API_BASE}/api/collections/tags/records/${tagId}`, {
          signal: controller.signal
        })

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}`)
        }

        const tag: TagRecord = await response.json()

        if (!tag.name || !tag.slug) {
          throw new Error('Invalid tag data')
        }

        this.updateTagPillWithLink(pill, tag)
      } catch (error) {
        console.debug(`Tag resolution failed for ${tagId}:`, error)
        this.markTagFailed(pill, tagId)
        throw error
      }
    }

    private updateTagPillWithLink(pill: HTMLElement, tag: TagRecord): void {
      const contentSpan = pill.querySelector('.tag-content')
      const loadingSpan = pill.querySelector('.tag-loading')

      if (contentSpan && loadingSpan) {
        // Create clickable tag link
        const tagLink = document.createElement('a')
        tagLink.href = `/tags/${tag.slug}`
        tagLink.textContent = tag.name
        tagLink.className = 'tag-link'
        tagLink.title = `View all links tagged with "${tag.name}"`

        // Add click tracking for analytics if needed
        tagLink.addEventListener('click', () => {
          console.debug(`Tag clicked: ${tag.name} (${tag.slug})`)
        })

        // Replace content and remove loading indicator
        contentSpan.replaceWith(tagLink)
        loadingSpan.remove()
      }

      // Update pill classes for styling
      pill.classList.remove('loading')
      pill.classList.add('resolved')
      pill.setAttribute('data-tag-slug', tag.slug)
    }

    private markTagFailed(pill: HTMLElement, tagId: string): void {
      const loadingSpan = pill.querySelector('.tag-loading')
      const contentSpan = pill.querySelector('.tag-content')

      if (loadingSpan) {
        loadingSpan.remove()
      }

      if (contentSpan) {
        contentSpan.textContent = `#${tagId}` // Prefix with # to indicate ID fallback
      }

      pill.classList.remove('loading')
      pill.classList.add('failed')
      pill.title = `Tag information unavailable (ID: ${tagId})`
    }

    // Testing method
    getResolvedTags(): Array<{slug: string, name: string}> {
      const resolved = Array.from(this.tagPills).filter(pill => pill.classList.contains('resolved'))
      return resolved.map(pill => ({
        slug: pill.getAttribute('data-tag-slug') || '',
        name: pill.querySelector('a')?.textContent || ''
      }))
    }
  }

  // Update DOMContentLoaded event listener
  document.addEventListener('DOMContentLoaded', () => {
    const linkElement = document.querySelector('[data-link-id]') as HTMLElement
    if (!linkElement) {
      console.debug('No link element found')
      return
    }

    const linkId = linkElement.getAttribute('data-link-id')
    if (!linkId) {
      console.debug('No link ID found in data attribute')
      return
    }

    // Initialize view count updater
    const viewUpdater = new ViewCountUpdater(linkId)

    // Initialize tag resolver
    const tagResolver = new TagResolver(linkId)

    // Start both enhancements
    setTimeout(() => {
      viewUpdater.incrementView()
      tagResolver.resolveAllTags()
    }, 500)

    // Expose for testing purposes
    (window as any).__viewUpdater = viewUpdater
    (window as any).__tagResolver = tagResolver
  })
</script>
```

#### Additional CSS for Tag States
Add to the style section:
```css
/* Tag link styling */
.tag-link {
  color: inherit;
  text-decoration: none;
}

.tag-pill.resolved .tag-link:hover {
  text-decoration: underline;
}

.tag-pill.failed {
  background: var(--error-light-bg, #fef2f2);
  border-color: var(--error-light-border, #fecaca);
  color: var(--error-dark, #991b1b);
}

.tag-pill.failed:hover {
  background: var(--error-light-bg, #fef2f2);
  transform: none;
  cursor: help;
}

/* Loading animation enhancement */
.tag-loading {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 1; }
}
```

#### Acceptance Criteria
- [ ] Tags resolve from IDs to names automatically on page load
- [ ] Resolved tags become clickable links to `/tags/:slug`
- [ ] Failed tag resolutions show fallback ID with # prefix
- [ ] Visual states clearly distinguish loading/resolved/failed tags
- [ ] Multiple tags resolve in parallel for better performance
- [ ] Partial failures don't prevent other tags from resolving
- [ ] Tag links include accessibility attributes (title, etc.)
- [ ] Resolution process provides debugging information in console

#### Testing Approach
- **Functional**: Verify tags change from IDs to names after page load
- **Linking**: Click resolved tag links to verify correct `/tags/:slug` URLs
- **Error Handling**: Test with invalid tag IDs (should show fallback)
- **Performance**: Check network tab for parallel tag requests
- **Console Testing**: Use `window.__tagResolver.getResolvedTags()` to inspect state
- **Partial Failure**: Create links with mix of valid/invalid tag IDs
- **Database**: Verify tag resolution queries in PocketBase admin panel

---

### Step 6: Testing and Verification
**Duration**: 45-60 minutes
**Complexity**: Medium
**Dependencies**: Steps 1-5 complete

#### Comprehensive Testing Strategy

##### 6.1: Functional Testing
**SSR and Basic Functionality**:
1. **Valid Link Display**: Test with existing link from database
   - Verify all fields display correctly (title, URL, description, view count, tags)
   - Check page title updates dynamically
   - Confirm creation date formatting is readable

2. **Error Handling**: Test all error scenarios
   - Non-existent link ID → 404 error page
   - Invalid/malformed ID → appropriate error handling
   - Backend server stopped → network error
   - Slow network conditions → timeout error

3. **Edge Cases**: Test data variations
   - Links without descriptions → description section hidden
   - Links without tags → "No tags assigned" message
   - Links with view_count = null/0 → displays "0 views"
   - Very long URLs → proper text wrapping
   - Very long titles → responsive text handling

##### 6.2: Client-Side Enhancement Testing
**View Count Increment**:
1. **Basic Increment**: Load page and verify view count increases by 1
2. **Database Consistency**: Check PocketBase admin panel for updated count
3. **Multiple Views**: Refresh page several times, verify each increment
4. **API Failure Handling**: Stop backend, verify page remains functional
5. **Animation**: Observe smooth count update animation

**Tag Resolution**:
1. **Complete Resolution**: Links with all valid tag IDs resolve to names/links
2. **Partial Resolution**: Links with mix of valid/invalid tags handle gracefully
3. **No Tags**: Links without tags show appropriate message
4. **Failed Resolution**: Invalid tag IDs fallback to #tagId display
5. **Parallel Loading**: Multiple tags resolve simultaneously

##### 6.3: Integration Testing
**Navigation and Routing**:
1. **Inbound Links**: Navigate from other pages to `/links/:id`
2. **Tag Links**: Click resolved tag pills (will show 404 until T-004-04 complete)
3. **External Links**: Click main URL link opens in new tab
4. **Back Navigation**: Browser back button works correctly
5. **Direct URL Access**: Typing `/links/:id` in address bar loads correctly

**Performance and Loading**:
1. **Initial Load**: Page renders quickly with SSR content
2. **Enhancement Speed**: View count and tags enhance within reasonable time
3. **Network Efficiency**: Minimal API calls (1 for view, N for tag resolution)
4. **Timeout Handling**: Long API calls don't block user interface

##### 6.4: Cross-Browser and Device Testing
**Responsive Design**:
1. **Desktop**: Test on standard desktop browser (Chrome/Firefox)
2. **Mobile**: Test on mobile viewport using dev tools
3. **Tablet**: Test on tablet-sized viewport
4. **Text Scaling**: Test with browser zoom at 125%, 150%

**Accessibility**:
1. **Keyboard Navigation**: Tab through all interactive elements
2. **Screen Reader**: Test with screen reader (alt text, labels)
3. **Color Contrast**: Verify text meets accessibility standards
4. **Focus Indicators**: All focusable elements show clear focus states

#### Testing Checklist

**Pre-Implementation Verification**:
- [ ] Backend API `/api/links/:id/view` endpoint is working
- [ ] PocketBase has accessible links and tags collections
- [ ] At least 2-3 test links exist with various tag configurations
- [ ] Development server is running on correct port

**Step-by-Step Testing**:
- [ ] **Step 1**: Basic routing works, 404 redirect functions
- [ ] **Step 2**: Data fetching works, all error states display correctly
- [ ] **Step 3**: Template renders properly, styling matches design
- [ ] **Step 4**: View count increments and updates UI
- [ ] **Step 5**: Tags resolve to names and become clickable links
- [ ] **Step 6**: All integration points work correctly

**Final Acceptance Testing**:
- [ ] Page loads correctly from fresh browser session
- [ ] All visual elements match existing design patterns
- [ ] Error handling provides clear, helpful feedback
- [ ] Performance meets expectations (loads within 2-3 seconds)
- [ ] Mobile and desktop experiences are optimized
- [ ] Accessibility requirements are met
- [ ] Code follows established patterns from `/watch/[id].astro`

#### Debugging and Troubleshooting Guide

**Common Issues and Solutions**:

1. **API Connection Errors**:
   - Check API_BASE URL configuration
   - Verify PocketBase server is running
   - Check browser console for CORS errors
   - Confirm firewall/network isn't blocking requests

2. **Tag Resolution Failures**:
   - Verify tag IDs exist in tags collection
   - Check tag records have required name/slug fields
   - Test individual tag API endpoints manually
   - Inspect network requests for HTTP error codes

3. **View Count Not Updating**:
   - Check POST request succeeds in network tab
   - Verify endpoint returns updated record with new count
   - Test with different link IDs
   - Check for JavaScript errors in console

4. **Styling Issues**:
   - Verify CSS custom properties are defined in global theme
   - Check for CSS syntax errors in browser dev tools
   - Test with base styles disabled to isolate issues
   - Compare with watch page for consistency

5. **Performance Problems**:
   - Check for excessive API calls in network tab
   - Profile page load times in dev tools
   - Test with slow network conditions
   - Verify timeout settings are appropriate

#### Success Metrics
- **Functionality**: 100% of test cases pass
- **Performance**: Page loads within 3 seconds on typical connection
- **Usability**: Users can complete all primary actions without confusion
- **Reliability**: Error states provide clear recovery paths
- **Consistency**: Experience matches established patterns throughout application

## Risk Assessment and Mitigation

### High-Risk Areas

**1. Tag Resolution Complexity**
- **Risk**: Multiple API calls for tag resolution could fail partially
- **Mitigation**: Use Promise.allSettled for independent resolution, graceful fallbacks
- **Monitoring**: Console logging for resolution success/failure rates

**2. View Count Race Conditions**
- **Risk**: Concurrent requests could cause view count inconsistencies
- **Mitigation**: Backend endpoint uses atomic SQL UPDATE, client-side is optimistic
- **Verification**: Database integrity checks during testing

**3. API Endpoint Availability**
- **Risk**: Required PocketBase endpoints might have access restrictions
- **Mitigation**: Research confirms public read access, existing view endpoint works
- **Fallback**: Error handling provides user-friendly messages

### Medium-Risk Areas

**4. Performance Impact of Multiple Requests**
- **Risk**: Tag resolution adds N additional API calls per page load
- **Mitigation**: 3-second timeouts, parallel resolution, non-blocking enhancement
- **Monitoring**: Network tab analysis during testing

**5. Error State Coverage**
- **Risk**: New error scenarios not covered by existing patterns
- **Mitigation**: Comprehensive error categorization reused from watch page
- **Testing**: Systematic testing of all failure modes

### Low-Risk Areas

**6. CSS Integration**
- **Risk**: New styles conflict with existing theme system
- **Mitigation**: Exclusive use of established CSS custom properties
- **Verification**: Visual comparison with existing pages

**7. TypeScript Compatibility**
- **Risk**: Type mismatches cause compilation errors
- **Mitigation**: Interface definitions match existing patterns
- **Testing**: Browser console checks for TypeScript errors

## Time Estimates and Dependencies

### Total Estimated Time: 4.5-6 hours

**Critical Path Steps**:
1. Step 1 + Step 2: 75-105 minutes (SSR foundation)
2. Step 3: 60-75 minutes (UI implementation)
3. Step 4: 30-45 minutes (View counting)
4. Step 5: 45-60 minutes (Tag resolution)
5. Step 6: 45-60 minutes (Testing and verification)

**Parallel Work Opportunities**:
- Steps 4 and 5 can be developed in parallel after Step 3
- CSS refinements can happen alongside JavaScript development
- Testing can begin as early as Step 2 completion

**External Dependencies**:
- ✅ PocketBase collections API (available)
- ✅ View increment endpoint (T-003-02 complete)
- ✅ BaseLayout.astro (T-004-01 complete)
- ✅ CSS custom properties system (established)

## Implementation Notes

### Code Quality Standards
- **Pattern Consistency**: Exactly match `/watch/[id].astro` structure and error handling
- **TypeScript**: All interfaces defined, proper typing throughout
- **Accessibility**: WCAG 2.1 AA compliance, keyboard navigation, screen reader support
- **Performance**: Optimize for Core Web Vitals, minimal JavaScript execution
- **Documentation**: Inline comments for complex logic, console debugging support

### Future Considerations
- **Tag Page Integration**: Links to `/tags/:slug` prepared for T-004-04 implementation
- **Caching Strategy**: Consider adding client-side caching for resolved tags
- **Analytics Integration**: View count increment provides hook for user behavior tracking
- **Progressive Enhancement**: Page fully functional without JavaScript

This implementation plan provides a systematic approach to building the link detail page while leveraging existing patterns and ensuring robust functionality across all user scenarios.