# Design: Astro Layout and Navigation

## Design Options

### Option 1: Minimal Scoped Components
Use Astro's built-in scoped styles with simple component composition.

**Approach:**
- `BaseLayout.astro` with scoped styles
- `Navigation.astro` component with inline responsive CSS
- CSS variables for theme consistency
- Media queries for mobile breakpoints

**Pros:**
- No additional dependencies
- Fast build times
- Simple to understand
- Matches current inline style approach

**Cons:**
- Styles scattered across components
- No utility classes for rapid development
- Mobile menu requires custom JavaScript

### Option 2: Utility-First with Tailwind
Add Tailwind CSS for utility-based styling.

**Approach:**
- Install and configure Tailwind
- Use utility classes for all styling
- Tailwind's responsive prefixes for mobile
- Component classes via `@apply`

**Pros:**
- Rapid development with utilities
- Consistent spacing/colors
- Built-in responsive system
- Large ecosystem

**Cons:**
- Additional dependency
- Build complexity increase
- Learning curve if unfamiliar
- Overkill for simple layouts

### Option 3: CSS Modules + Design Tokens
Use CSS modules with centralized design tokens.

**Approach:**
- `styles/tokens.css` with CSS variables
- Component-specific `.module.css` files
- Shared layout styles in `global.css`
- PostCSS for nesting/features

**Pros:**
- Explicit style dependencies
- No global namespace pollution
- Reusable token system
- Modern CSS features

**Cons:**
- More boilerplate
- Additional build config
- Split between .astro and .css files

## Selected Approach: Minimal Scoped Components

**Rationale:**
1. **Simplicity First**: The project needs basic layout/nav. Starting minimal avoids premature optimization.
2. **No Dependencies**: Aligns with current setup (only Astro installed).
3. **Progressive Enhancement**: Can add utilities later if needed.
4. **Fast Iteration**: Inline styles allow rapid prototyping.
5. **Astro-Native**: Uses framework's intended patterns.

## Component Architecture

### BaseLayout.astro
```
Props: { title, description? }
Structure:
  <html>
    <head> - meta, title, global styles
    <body>
      <Navigation /> - shared header/nav
      <main>
        <slot /> - page content
      </main>
      <footer> - minimal footer
```

### Navigation.astro
```
Props: none (hardcoded links for now)
Structure:
  <header>
    <nav>
      <div class="nav-brand"> - site name/logo
      <button class="nav-toggle"> - mobile menu button
      <ul class="nav-links"> - Home, Stats
```

## Styling Strategy

### Design Tokens (CSS Variables)
```css
:root {
  /* Colors */
  --color-bg: white;
  --color-text: #333;
  --color-primary: #111;

  /* Spacing */
  --space-xs: 0.25rem;
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 2rem;

  /* Layout */
  --max-width: 1200px;
  --mobile-breakpoint: 768px;
}
```

### Mobile Navigation Pattern
- Desktop: Horizontal nav bar
- Mobile: Hidden by default
- Toggle button shows/hides menu
- Use checkbox hack or minimal JS
- Slide or dropdown animation

### Layout Constraints
```css
main {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: var(--space-lg);
}
```

## Responsive Breakpoints

### Mobile First Approach
- Base styles: Mobile (< 768px)
- Tablet: @media (min-width: 768px)
- Desktop: @media (min-width: 1024px)

### Navigation Behavior
- Mobile: Vertical stack, toggle button
- Tablet+: Horizontal bar, no toggle

## Implementation Path

1. Create `src/layouts/BaseLayout.astro`
2. Create `src/components/Navigation.astro`
3. Add CSS variables to BaseLayout
4. Implement mobile toggle (pure CSS first)
5. Update index.astro to use BaseLayout
6. Test responsive behavior

## Trade-offs Accepted

- **No utility classes**: Slower styling but simpler setup
- **Inline styles**: Less reusable but faster to modify
- **Basic mobile menu**: Not as polished but functional
- **No CSS preprocessing**: Modern CSS is sufficient

## Success Criteria Validation

✓ HTML shell with slot pattern
✓ Navigation with required links
✓ Scoped styles (Astro default)
✓ Responsive via media queries
✓ Reusable by all pages