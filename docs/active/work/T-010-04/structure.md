# Structure: Frontend HTML/CSS/A11y Cleanup (T-010-04)

## File Operations Overview

### Create (1 file)
- `frontend/src/components/ErrorState.astro`

### Modify (16 files)
- `frontend/astro.config.mjs`
- `frontend/src/layouts/BaseLayout.astro`
- `frontend/src/components/Navigation.astro`
- `frontend/src/components/LinksList.astro`
- `frontend/src/components/GoLive.astro`
- `frontend/src/components/PresenterController.astro`
- `frontend/src/components/SearchInterface.astro`
- `frontend/src/components/StatsSummary.astro`
- `frontend/src/pages/index.astro`
- `frontend/src/pages/tags/[slug].astro`
- `frontend/src/pages/links/[id].astro`
- `frontend/src/pages/present/[id].astro`
- `frontend/src/pages/present/index.astro`
- `frontend/src/pages/sync/[id].astro`
- `frontend/src/pages/sync/[id]/control.astro`
- `frontend/src/pages/watch/[id].astro`

### Delete (0 files)
- None

## Detailed File Changes

### 1. Create ErrorState Component

**File**: `frontend/src/components/ErrorState.astro`

```astro
---
export interface Props {
  title: string;
  message: string;
  showRetry?: boolean;
  backLink?: string;
  backLabel?: string;
}

const {
  title,
  message,
  showRetry = true,
  backLink = '/',
  backLabel = '← Back to Home'
} = Astro.props;
---

<div class="error-state">
  <div class="error-icon">⚠️</div>
  <h1 class="error-title">{title}</h1>
  <p class="error-message">{message}</p>
  <div class="error-actions">
    <a href={backLink} class="back-link">{backLabel}</a>
    {showRetry && (
      <button onclick="window.location.reload()" class="retry-button">Try Again</button>
    )}
  </div>
</div>

<style>
  /* Error state styles */
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    min-height: 50vh;
  }

  .error-icon {
    font-size: 3rem;
    margin-bottom: var(--space-lg);
    opacity: 0.8;
  }

  .error-title {
    font-size: 2rem;
    margin-bottom: var(--space-md);
    color: var(--color-primary);
    font-weight: 600;
  }

  .error-message {
    font-size: 1.125rem;
    color: var(--color-text);
    opacity: 0.8;
    margin-bottom: var(--space-xl);
    max-width: 500px;
    line-height: 1.6;
  }

  .error-actions {
    display: flex;
    gap: var(--space-md);
    flex-wrap: wrap;
    justify-content: center;
  }

  .back-link,
  .retry-button {
    color: var(--color-primary);
    text-decoration: none;
    font-weight: 500;
    padding: var(--space-sm) var(--space-lg);
    border: 1px solid var(--color-border);
    border-radius: var(--border-radius);
    background-color: transparent;
    font-family: var(--font-body);
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .back-link:hover,
  .retry-button:hover {
    background-color: var(--color-footer);
    text-decoration: none;
    border-color: var(--color-primary);
  }

  .retry-button:active {
    transform: translateY(1px);
  }
</style>
```

### 2. Astro Configuration

**File**: `frontend/astro.config.mjs`
- Line 5: Change `output: 'static'` to `output: 'hybrid'`

### 3. BaseLayout Updates

**File**: `frontend/src/layouts/BaseLayout.astro`
- Lines 21-28: Add new CSS custom properties after existing ones:
  ```css
  --color-success: #10b981;
  --color-error: #ef4444;
  --color-secondary: #6b7280;
  --border-radius: 4px;
  --border-radius-lg: 8px;
  ```
- Line 105: Change `&copy; 2024` to `&copy; {new Date().getFullYear()}`

### 4. Navigation Component

**File**: `frontend/src/components/Navigation.astro`
- Lines 2-3: Add script imports and current path logic
- Line 5: Add `aria-label="Main navigation"` to nav element
- Lines 9-14: Replace checkbox hack with button element:
  ```astro
  <button
    class="nav-toggle-button"
    aria-expanded="false"
    aria-controls="nav-menu"
    aria-label="Toggle navigation menu"
  >
    <span></span>
    <span></span>
    <span></span>
  </button>
  ```
- Line 16: Add id="nav-menu" to ul element
- Lines 17-19: Add aria-current logic to each nav link
- Add script section for button functionality

### 5. Page Main Element Fixes

**File**: `frontend/src/pages/index.astro`
- Line 97: Change `<main class="home-page-container">` to `<div class="home-page-container">`
- Line 131: Change closing `</main>` to `</div>`
- Lines 316, 321: Remove console.log statements

**File**: `frontend/src/pages/tags/[slug].astro`
- Line 79: Change `<main class="tag-page-container">` to `<div class="tag-page-container">`
- Corresponding closing tag change
- Lines 329-341: Remove `.loading-skeleton` class and `@keyframes loading`
- Lines 373, 378: Remove console.log statements
- Lines 81-90: Replace error state with ErrorState component usage

**File**: `frontend/src/pages/present/[id].astro`
- Line 201: Change `<main class="presenter-control-page">` to `<div class="presenter-control-page">`
- Corresponding closing tag change

**File**: `frontend/src/components/PresenterController.astro`
- Line 67: Change `<main class="controller-main">` to `<div class="controller-main">`
- Corresponding closing tag change
- Line 34: Remove API_BASE definition, import from lib/api instead
- Remove console.log statements

### 6. LinksList Component

**File**: `frontend/src/components/LinksList.astro`
- Line 98: Remove `role="main"` attribute from div

### 7. Links Page Tag URL Fix

**File**: `frontend/src/pages/links/[id].astro`
- Line 100: Change `/search?tag=${tag}` to `/tags/${encodeURIComponent(tag)}`
- Line 363: Change `/search?tag=${tagSlug}` to `/tags/${encodeURIComponent(tagSlug)}`
- Remove console.log statements

### 8. API Standardization

**File**: `frontend/src/components/GoLive.astro`
- Add import: `import { API_BASE } from '../lib/api';`
- Line 38: Change `/api/presentations/` to `${API_BASE}/api/presentations/`

### 9. Console.log Removals

**File**: `frontend/src/components/SearchInterface.astro`
- Remove all console.log statements (keep console.error if any)

**File**: `frontend/src/components/StatsSummary.astro`
- Remove all console.log statements

**File**: `frontend/src/pages/sync/[id].astro`
- Remove all console.log statements

**File**: `frontend/src/pages/watch/[id].astro`
- Remove all console.log statements

### 10. Dead CSS Removal

**File**: `frontend/src/pages/sync/[id]/control.astro`
- Lines 1426-1448: Remove entire `.temp-debug` CSS block and related rules
- Remove any console.log statements

### 11. Error State Replacements

**File**: `frontend/src/pages/index.astro`
- Import ErrorState component
- Lines 108-117: Replace error state HTML with:
  ```astro
  <ErrorState
    title={errorMessages[error].title}
    message={errorMessages[error].message}
  />
  ```
- Lines 162-224: Remove duplicate error state CSS

**File**: `frontend/src/pages/present/index.astro`
- Apply same ErrorState component pattern
- Remove duplicate error CSS

## Component Interfaces

No changes to existing component props or public interfaces. The ErrorState component is new and self-contained.

## CSS Architecture

- Preserve all existing class names for backward compatibility
- Add new CSS custom properties to root for theming
- ErrorState component contains its own scoped styles
- Removed CSS is confirmed unused through grep searches

## Build Impact

All changes are compatible with Astro's hybrid output mode. No breaking changes to build process.