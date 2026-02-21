# Structure: Link Detail Page Implementation (T-004-03)

## File Changes Overview

### New Files Created

**`/frontend/src/pages/links/[id].astro`** - Main dynamic route file (~200-250 lines)
- Primary implementation file for link detail page
- Follows `/frontend/src/pages/watch/[id].astro` pattern exactly
- Contains SSR section, template section, styles section, and client-side script section

### Modified Files

**No existing files require modification** - Implementation leverages existing APIs and patterns without changes.

### Deleted Files

**None** - Pure additive implementation.

## Module Boundaries and Public Interfaces

### Data Access Layer

**PocketBase API Integration**:
```typescript
// SSR Data Fetching
GET /api/collections/links/records/:id → LinkRecord
GET /api/collections/tags/records/:tagId → TagRecord

// Client-Side Enhancement
POST /api/links/:id/view → UpdatedLinkRecord
```

**Data Interfaces** (defined inline in component):
```typescript
interface LinkItem {
  id: string
  url: string
  title: string
  description: string
  view_count: number
  tags: string[] // Array of tag IDs
  created: string
  updated: string
}

interface TagItem {
  id: string
  name: string
  slug: string
}

interface ErrorState {
  type: 'notfound' | 'network' | 'server' | 'timeout' | null
  title?: string
  message?: string
}
```

### Component Architecture

**`/frontend/src/pages/links/[id].astro`** exports:
- No public API (page component)
- Internal class: `ViewCountUpdater` for client-side state
- Internal class: `TagResolver` for progressive enhancement

**External Dependencies**:
- `BaseLayout.astro` - Layout wrapper (import)
- `/frontend/src/styles/animations.css` - Animation utilities (implicit)
- CSS custom properties from global theme system

## Internal Organization

### `/frontend/src/pages/links/[id].astro` Structure

**Section 1: SSR Logic** (~60 lines)
```astro
---
import BaseLayout from '../../layouts/BaseLayout.astro'

// Type definitions (inline)
interface LinkItem { ... }
interface ErrorState { ... }

// Parameter extraction and validation
const { id } = Astro.params
if (!id) {
  return Astro.redirect('/404')
}

// API configuration
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090'
const TIMEOUT_MS = 5000

// Data fetching with error handling
let link: LinkItem | null = null
let error: ErrorState | null = null

try {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), TIMEOUT_MS)

  const response = await fetch(`${API_BASE}/api/collections/links/records/${id}`, {
    signal: controller.signal
  })

  clearTimeout(timeoutId)

  if (response.status === 404) {
    error = { type: 'notfound', title: 'Link Not Found', message: '...' }
  } else if (!response.ok) {
    error = { type: 'server', title: 'Server Error', message: '...' }
  } else {
    link = await response.json()
  }
} catch (fetchError) {
  // Error categorization logic (identical to watch page)
  if (fetchError.name === 'AbortError') {
    error = { type: 'timeout', ... }
  } else {
    error = { type: 'network', ... }
  }
}

// Page metadata
const pageTitle = link ? `${link.title} - Link Details` : 'Link Not Found'
---
```

**Section 2: Template Structure** (~80 lines)
```astro
<BaseLayout title={pageTitle}>
  <div class="page-container">
    {error ? (
      <div class="error-container">
        <div class="error-content">
          <h1 class="error-title">{error.title}</h1>
          <p class="error-message">{error.message}</p>
          <a href="/" class="home-link">← Back to Home</a>
        </div>
      </div>
    ) : link && (
      <article class="link-detail">
        <header class="link-header">
          <h1 class="link-title">{link.title}</h1>
          <div class="link-meta">
            <span class="view-count" data-link-id={link.id} data-current-count={link.view_count}>
              <span class="view-number">{link.view_count}</span>
              <span class="view-label">views</span>
            </span>
            <time class="link-date" datetime={link.created}>
              {new Date(link.created).toLocaleDateString()}
            </time>
          </div>
        </header>

        <div class="link-content">
          <div class="link-url-section">
            <label class="url-label">URL:</label>
            <a href={link.url}
               class="link-url"
               target="_blank"
               rel="noopener noreferrer"
               data-link-id={link.id}>
              {link.url}
            </a>
          </div>

          {link.description && (
            <div class="link-description-section">
              <label class="description-label">Description:</label>
              <p class="link-description">{link.description}</p>
            </div>
          )}

          <div class="link-tags-section">
            <label class="tags-label">Tags:</label>
            <div class="tags-container" data-link-id={link.id}>
              {link.tags.length > 0 ? (
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

**Section 3: Styling** (~60 lines)
```astro
<style>
  /* CSS Custom Properties Integration */
  .page-container {
    max-width: var(--content-max-width, 800px);
    margin: 0 auto;
    padding: var(--spacing-md, 1rem);
  }

  /* Error state styling (identical to watch page) */
  .error-container { ... }
  .error-content { ... }
  .home-link { ... }

  /* Link detail layout */
  .link-detail {
    background: var(--surface-color, white);
    border-radius: var(--radius-md, 8px);
    padding: var(--spacing-lg, 1.5rem);
    box-shadow: var(--shadow-sm, 0 2px 4px rgba(0,0,0,0.1));
  }

  .link-header {
    border-bottom: 1px solid var(--border-color, #e5e7eb);
    padding-bottom: var(--spacing-md, 1rem);
    margin-bottom: var(--spacing-md, 1rem);
  }

  .link-title {
    font-size: var(--text-xl, 1.25rem);
    font-weight: 600;
    margin: 0 0 var(--spacing-sm, 0.5rem) 0;
    color: var(--text-primary, #111827);
  }

  .link-meta {
    display: flex;
    gap: var(--spacing-md, 1rem);
    align-items: center;
    font-size: var(--text-sm, 0.875rem);
    color: var(--text-secondary, #6b7280);
  }

  /* URL section */
  .link-url-section,
  .link-description-section,
  .link-tags-section {
    margin-bottom: var(--spacing-md, 1rem);
  }

  .url-label,
  .description-label,
  .tags-label {
    display: block;
    font-weight: 500;
    margin-bottom: var(--spacing-xs, 0.25rem);
    color: var(--text-primary, #111827);
  }

  .link-url {
    display: inline-block;
    color: var(--link-color, #3b82f6);
    text-decoration: underline;
    word-break: break-all;
    padding: var(--spacing-xs, 0.25rem) var(--spacing-sm, 0.5rem);
    border-radius: var(--radius-sm, 4px);
    transition: background-color 0.2s;
  }

  .link-url:hover {
    background-color: var(--link-hover-bg, #eff6ff);
  }

  /* Tags styling */
  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-xs, 0.25rem);
  }

  .tag-pill {
    display: inline-flex;
    align-items: center;
    padding: var(--spacing-xs, 0.25rem) var(--spacing-sm, 0.5rem);
    background: var(--tag-bg, #f3f4f6);
    border: 1px solid var(--tag-border, #d1d5db);
    border-radius: var(--radius-full, 9999px);
    font-size: var(--text-xs, 0.75rem);
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
  }

  .tag-loading {
    margin-left: var(--spacing-xs, 0.25rem);
    font-size: var(--text-xs, 0.75rem);
    opacity: 0.6;
    animation: pulse 1.5s ease-in-out infinite;
  }

  /* Responsive design */
  @media (max-width: 640px) {
    .link-meta {
      flex-direction: column;
      align-items: flex-start;
      gap: var(--spacing-xs, 0.25rem);
    }

    .link-url {
      word-break: break-all;
    }
  }

  /* Accessibility */
  @media (prefers-reduced-motion: reduce) {
    .tag-loading {
      animation: none;
    }

    .tag-pill,
    .link-url {
      transition: none;
    }
  }
</style>
```

**Section 4: Client-Side Enhancement** (~40 lines)
```astro
<script>
  // Type definitions for client-side
  interface LinkRecord {
    id: string
    view_count: number
  }

  interface TagRecord {
    id: string
    name: string
    slug: string
  }

  class ViewCountUpdater {
    private linkId: string
    private viewCountElement: Element | null

    constructor(linkId: string) {
      this.linkId = linkId
      this.viewCountElement = document.querySelector(`[data-link-id="${linkId}"][data-current-count]`)
    }

    async incrementView(): Promise<void> {
      try {
        const controller = new AbortController()
        setTimeout(() => controller.abort(), 5000)

        const response = await fetch(`/api/links/${this.linkId}/view`, {
          method: 'POST',
          signal: controller.signal
        })

        if (response.ok) {
          const updatedLink: LinkRecord = await response.json()
          this.updateDisplay(updatedLink.view_count)
        }
      } catch (error) {
        // Silent failure - page remains functional
        console.debug('View count update failed:', error)
      }
    }

    private updateDisplay(newCount: number): void {
      if (this.viewCountElement) {
        const numberSpan = this.viewCountElement.querySelector('.view-number')
        if (numberSpan) {
          numberSpan.textContent = newCount.toString()
        }
      }
    }
  }

  class TagResolver {
    private linkId: string
    private tagsContainer: Element | null

    constructor(linkId: string) {
      this.linkId = linkId
      this.tagsContainer = document.querySelector(`[data-link-id="${linkId}"] .tags-container`)
    }

    async resolveAllTags(): Promise<void> {
      if (!this.tagsContainer) return

      const tagPills = this.tagsContainer.querySelectorAll('.tag-pill[data-tag-id]')
      const resolutionPromises = Array.from(tagPills).map(pill =>
        this.resolveTagPill(pill as HTMLElement)
      )

      await Promise.allSettled(resolutionPromises)
    }

    private async resolveTagPill(pill: HTMLElement): Promise<void> {
      const tagId = pill.dataset.tagId
      if (!tagId) return

      try {
        const controller = new AbortController()
        setTimeout(() => controller.abort(), 3000)

        const response = await fetch(`/api/collections/tags/records/${tagId}`, {
          signal: controller.signal
        })

        if (response.ok) {
          const tag: TagRecord = await response.json()
          this.updateTagPill(pill, tag)
        } else {
          this.markTagFailed(pill)
        }
      } catch (error) {
        this.markTagFailed(pill)
      }
    }

    private updateTagPill(pill: HTMLElement, tag: TagRecord): void {
      const contentSpan = pill.querySelector('.tag-content')
      const loadingSpan = pill.querySelector('.tag-loading')

      if (contentSpan) {
        contentSpan.textContent = tag.name
        // Convert to link
        const link = document.createElement('a')
        link.href = `/tags/${tag.slug}`
        link.textContent = tag.name
        link.className = 'tag-link'
        contentSpan.replaceWith(link)
      }

      if (loadingSpan) {
        loadingSpan.remove()
      }

      pill.classList.remove('loading')
      pill.classList.add('resolved')
      pill.title = `View all links tagged with "${tag.name}"`
    }

    private markTagFailed(pill: HTMLElement): void {
      const loadingSpan = pill.querySelector('.tag-loading')
      if (loadingSpan) {
        loadingSpan.remove()
      }
      pill.classList.remove('loading')
      // Keep showing tag ID as fallback
    }
  }

  // Initialize on page load
  document.addEventListener('DOMContentLoaded', () => {
    const linkIdElement = document.querySelector('[data-link-id]')
    if (!linkIdElement) return

    const linkId = linkIdElement.getAttribute('data-link-id')
    if (!linkId) return

    // Initialize view count updater
    const viewUpdater = new ViewCountUpdater(linkId)
    viewUpdater.incrementView()

    // Initialize tag resolver
    const tagResolver = new TagResolver(linkId)
    tagResolver.resolveAllTags()
  })
</script>
```

## Component Interaction Flow

### SSR Phase Data Flow

1. **Route Resolution**: Astro extracts `id` from URL parameter `/links/[id]`
2. **Parameter Validation**: Redirect to `/404` if `id` is undefined/null
3. **API Data Fetching**:
   - Call `GET /api/collections/links/records/${id}` with timeout/abort
   - Categorize errors: `notfound`, `network`, `server`, `timeout`
   - Store result in `link: LinkItem | null` and `error: ErrorState | null`
4. **Template Data Preparation**: Pass `link` and `error` to template rendering
5. **Metadata Generation**: Create `pageTitle` based on link data availability

### Client-Side Enhancement Flow

1. **DOM Ready**: Wait for `DOMContentLoaded` event
2. **Element Discovery**: Find `[data-link-id]` element to extract link ID
3. **Parallel Enhancement**:
   - **View Count Update**: `ViewCountUpdater.incrementView()` calls `POST /api/links/:id/view`
   - **Tag Resolution**: `TagResolver.resolveAllTags()` calls `GET /api/collections/tags/records/:tagId`
4. **DOM Updates**:
   - Update view count number in `.view-number` span
   - Replace tag ID text with tag names and convert to links
   - Update CSS classes for visual feedback (`loading` → `resolved`)

### Error Handling Flow

**SSR Errors**:
- Network/timeout → Show error container with retry guidance
- 404 Not Found → Show link-specific not found message
- Server errors → Show generic server error message

**Client-Side Errors**:
- View count increment failure → Silent (page remains functional)
- Tag resolution failure → Keep showing tag IDs as fallback
- Partial tag resolution → Mix of resolved names and fallback IDs

## Integration Points

### External API Dependencies

**PocketBase Auto-Generated Endpoints**:
- `GET /api/collections/links/records/:id` - Link data retrieval
- `GET /api/collections/tags/records/:id` - Individual tag resolution

**Custom Application Endpoints**:
- `POST /api/links/:id/view` - Atomic view count increment (existing)

### Layout System Integration

**BaseLayout.astro Integration**:
- Inherits global CSS custom properties
- Uses established navigation and header structure
- Maintains responsive breakpoint consistency
- Follows accessibility patterns (semantic HTML, focus management)

### Styling System Integration

**CSS Custom Properties Usage**:
- `--content-max-width` - Page container width
- `--spacing-*` values - Consistent spacing scale
- `--text-*` values - Typography scale
- `--surface-color`, `--border-color` - Theme colors
- `--link-color`, `--link-hover-bg` - Link styling

**Animation System Integration**:
- Uses existing `pulse` animation for loading states
- Respects `prefers-reduced-motion` accessibility setting
- Follows established transition timing patterns

## Ordering of Changes

### Implementation Order (Critical Dependencies)

**Phase 1: File Creation**
1. Create `/frontend/src/pages/links/[id].astro` with complete implementation
2. No dependencies - all required APIs and layouts already exist

**Phase 2: Testing & Validation**
1. Test SSR rendering with existing link IDs
2. Verify error handling with non-existent IDs
3. Test client-side enhancements (view count, tag resolution)
4. Validate responsive design and accessibility

**Phase 3: Integration Verification**
1. Confirm navigation integration (links from other pages work)
2. Verify tag links route correctly (future `/tags/[slug]` pages)
3. Test with various link data states (no tags, no description, etc.)

### No Modification Dependencies

**Existing Files Remain Unchanged**:
- Backend routes and APIs are already functional
- Layout and styling systems require no updates
- No migration or configuration changes needed

This structure provides a complete blueprint for implementing the link detail page following established patterns while providing robust error handling and progressive enhancement capabilities.