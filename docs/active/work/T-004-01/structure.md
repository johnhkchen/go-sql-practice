# Structure: Astro Layout and Navigation

## File Operations

### Create Files
```
frontend/src/layouts/BaseLayout.astro
frontend/src/components/Navigation.astro
```

### Modify Files
```
frontend/src/pages/index.astro
```

### Directory Structure After Changes
```
frontend/
  src/
    components/
      Navigation.astro       [NEW]
    layouts/
      BaseLayout.astro       [NEW]
    pages/
      index.astro            [MODIFIED]
```

## Component Specifications

### BaseLayout.astro

**Location:** `frontend/src/layouts/BaseLayout.astro`

**Props Interface:**
```typescript
export interface Props {
  title: string;
  description?: string;
}
```

**Structure:**
```astro
---
import Navigation from '../components/Navigation.astro';
const { title, description = "Link Bookmarks Application" } = Astro.props;
---
<!DOCTYPE html>
<html lang="en">
<head>
  <!-- Meta tags -->
  <!-- Global styles with CSS variables -->
</head>
<body>
  <Navigation />
  <main>
    <slot />
  </main>
  <footer>
    <!-- Simple footer content -->
  </footer>
</body>
</html>
```

**Responsibilities:**
- HTML document structure
- Meta tags and SEO
- Global CSS variables and reset styles
- Page layout grid
- Slot for page content

**CSS Scope:**
- `:root` CSS variables
- Global reset/normalize
- Body and main layout
- Footer positioning

### Navigation.astro

**Location:** `frontend/src/components/Navigation.astro`

**Props Interface:**
```typescript
// No props - hardcoded navigation for now
```

**Structure:**
```astro
---
// No frontmatter needed initially
---
<header>
  <nav>
    <div class="nav-container">
      <a href="/" class="nav-brand">Link Bookmarks</a>
      <input type="checkbox" id="nav-toggle" class="nav-toggle">
      <label for="nav-toggle" class="nav-toggle-label">
        <!-- Hamburger icon -->
      </label>
      <ul class="nav-links">
        <li><a href="/">Home</a></li>
        <li><a href="/stats">Stats</a></li>
      </ul>
    </div>
  </nav>
</header>
<style>
  /* Component-scoped navigation styles */
</style>
```

**Responsibilities:**
- Site branding/logo link
- Navigation link list
- Mobile menu toggle
- Responsive behavior

**CSS Scope:**
- Navigation layout
- Link styling
- Mobile toggle button
- Responsive breakpoints

### Updated index.astro

**Location:** `frontend/src/pages/index.astro`

**Structure:**
```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
---
<BaseLayout title="Link Bookmarks">
  <h1>Link Bookmarks</h1>
  <p>Welcome to the Link Bookmarks application.</p>
</BaseLayout>
```

**Changes:**
- Remove `<!DOCTYPE html>` and full HTML structure
- Import and use BaseLayout
- Pass title prop
- Keep only page-specific content

## CSS Architecture

### Global Variables (in BaseLayout)
```css
:root {
  /* Colors */
  --color-bg: #ffffff;
  --color-text: #333333;
  --color-primary: #111111;
  --color-border: #e5e5e5;

  /* Spacing */
  --space-xs: 0.25rem;
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 2rem;
  --space-xl: 3rem;

  /* Layout */
  --max-width: 1200px;
  --header-height: 60px;

  /* Typography */
  --font-body: system-ui, -apple-system, sans-serif;
  --line-height: 1.6;
}
```

### Component Styles
- Each component uses `<style>` tags (scoped by default)
- Reference CSS variables for consistency
- Mobile-first approach with min-width media queries

## Mobile Navigation Strategy

### Checkbox Hack Implementation
- Hidden checkbox input for state
- Label as hamburger button
- CSS sibling selectors for menu visibility
- No JavaScript required initially

### Breakpoints
- Mobile: 0-767px (vertical menu, toggle visible)
- Desktop: 768px+ (horizontal menu, toggle hidden)

## Module Boundaries

### BaseLayout Exports
- Props interface (TypeScript types)
- Default component export

### Navigation Exports
- Default component export only
- No props (self-contained)

### Page Integration
- Pages import BaseLayout
- Pass required title prop
- Optional description prop
- Content goes in default slot

## Dependencies

### Internal
- No circular dependencies
- Clear hierarchy: Pages → Layout → Components

### External
- Only Astro framework (already installed)
- No additional npm packages needed

## Build Impact

### File Size Estimates
- BaseLayout.astro: ~150 lines
- Navigation.astro: ~100 lines
- Updated index.astro: ~10 lines

### Performance
- No JavaScript by default (pure CSS mobile menu)
- Scoped styles prevent conflicts
- CSS variables reduce duplication