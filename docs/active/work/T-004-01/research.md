# Research: Astro Layout and Navigation

## Current State

### Frontend Structure
The frontend exists at `frontend/` with a minimal Astro 5.17.3 setup. Key files:
- `src/pages/index.astro` - Basic HTML page, no component structure
- `astro.config.mjs` - Static site configuration with directory format
- `package.json` - Minimal deps (only Astro), Node >=24 requirement

The index page contains:
- Inline HTML with `<!doctype html>` structure
- Embedded styles in `<style>` tag
- No component imports or layout inheritance
- Basic system font stack and simple styling

### Missing Components
No existing layout or component infrastructure:
- No `src/layouts/` directory
- No `src/components/` directory
- No shared navigation component
- No CSS framework or utility system

### Page Requirements
From the acceptance criteria, we need:
- Home page at `/` (exists as index.astro)
- Stats page at `/stats` (does not exist)

### Backend Context
The backend uses PocketBase (main.go), a self-contained backend with:
- Built-in database (SQLite by default)
- REST API endpoints
- Admin UI at `/_/`
- API typically at `/api/collections/{collection}/records`

### Astro Capabilities
Astro 5.x provides:
- Component-based architecture with `.astro` files
- Scoped CSS via `<style>` tags in components
- Layout inheritance via slots
- Static site generation (configured)
- TypeScript support (tsconfig.json exists)

### Styling Options
Current setup uses:
- Inline `<style>` tags with plain CSS
- System font stack
- Basic rem-based spacing
- No CSS preprocessing or utilities

Astro supports:
- Scoped styles (default behavior)
- Global styles via `is:global`
- CSS modules
- Sass/PostCSS (requires deps)
- Tailwind (requires deps)

### Build Configuration
- Output: static HTML
- Assets directory: `assets`
- Format: directory (clean URLs)
- Dev server: `npm run dev`
- Build: `npm run build`
- Preview: `npm run preview`

### Responsive Requirements
The acceptance criteria specify mobile support. Current index.astro has:
- `viewport` meta tag for mobile
- Rem-based spacing (scales)
- No responsive breakpoints
- No mobile navigation pattern

### Related Work
Previous tickets (T-001-*) set up:
- Backend with PocketBase
- Frontend initialization
- Basic project structure

This ticket (T-004-01) is the first frontend component work, establishing patterns for all subsequent UI development.

## Constraints and Assumptions

### Technical Constraints
- Must use existing Astro 5.17.3 (no framework changes)
- Node >=24 requirement
- Static site output (no SSR)
- Must work with PocketBase backend

### Design Constraints
- All pages must use the common layout
- Navigation must include Home and Stats links
- Must be responsive/mobile-friendly
- CSS approach should be maintainable

### Assumptions
- Navigation will be in header (typical pattern)
- Footer mentioned but content not specified
- "Basic CSS" suggests avoiding heavy frameworks
- Layout slot pattern for page content
- Stats page will be created in subsequent ticket

## Key Decisions Needed

1. **Component Organization**: Where to place layouts vs components
2. **CSS Strategy**: Scoped styles vs utility approach vs hybrid
3. **Navigation Pattern**: Desktop nav bar vs mobile menu handling
4. **Layout Composition**: How to structure BaseLayout for flexibility
5. **Path Handling**: Absolute vs relative paths for navigation