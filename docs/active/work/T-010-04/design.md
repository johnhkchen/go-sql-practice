# Design: Frontend HTML/CSS/A11y Cleanup (T-010-04)

## Design Objectives

1. Fix HTML semantic violations (nested mains, invalid roles)
2. Improve accessibility (navigation ARIA attributes)
3. Remove dead code (unused CSS, console.logs)
4. Centralize duplicated patterns (error states)
5. Fix configuration inconsistencies (Astro output mode)
6. Standardize API usage patterns

## Approach Options

### Option 1: Minimal Touch
- Fix only breaking issues (nested mains, Astro config)
- Leave console.logs for potential debugging value
- Keep duplicated error states as-is
- **Rejected**: Doesn't address code quality debt, console.logs hurt production

### Option 2: Full Rewrite
- Extract all duplicated code into components
- Create component library for all UI patterns
- Full accessibility audit and fixes
- **Rejected**: Scope creep beyond ticket requirements, high risk of regressions

### Option 3: Targeted Cleanup (Selected)
- Fix all HTML/accessibility issues
- Create ErrorState component for duplicated patterns
- Remove all console.logs (keep console.error)
- Fix configuration and API usage
- Add missing CSS tokens
- **Selected**: Balances improvement with stability, matches ticket scope exactly

## Design Decisions

### 1. Astro Configuration Fix
Change `output: 'static'` to `output: 'hybrid'` in astro.config.mjs. This resolves the contradiction where 8 pages use SSR (`prerender = false`) but config specifies static mode.

### 2. ErrorState Component Strategy

Create shared component with props:
```typescript
interface Props {
  title: string;
  message: string;
  showRetry?: boolean;
  backLink?: string;
  backLabel?: string;
}
```

Benefits:
- Reduces ~130 lines per page to ~10 line component usage
- Centralizes error UI updates
- Maintains consistent UX

### 3. Nested Main Resolution

Convert page-specific `<main>` elements to `<div>` since BaseLayout provides the semantic main wrapper:
- Keep existing classes for styling continuity
- No CSS changes needed (classes remain functional)
- HTML becomes semantically correct

### 4. Navigation Accessibility

Enhance Navigation.astro with:
- Compare `Astro.url.pathname` to link hrefs for active state
- Add `aria-current="page"` to active link
- Replace checkbox hack with button + JavaScript
- Add proper ARIA attributes (aria-expanded, aria-controls, aria-label)

Implementation approach:
- Use progressive enhancement (works without JS)
- Button controls menu visibility via aria-expanded
- Mobile menu references via aria-controls

### 5. Console.log Removal Strategy

- Remove all `console.log` statements entirely
- Keep `console.error` for actual error conditions
- No replacement needed - these are debug artifacts
- Search/remove pattern: `grep -r "console.log" frontend/src/`

### 6. Tag Navigation Fix

Change `/search?tag=${tag}` to `/tags/${encodeURIComponent(tag)}`:
- Aligns with existing route structure
- Maintains URL encoding for safety
- Update in both template and script sections

### 7. API Configuration Standardization

- Import `API_BASE` from `lib/api.ts` in all components
- Remove duplicate definitions
- GoLive.astro: Change `/api/presentations/` to `${API_BASE}/api/presentations/`
- PresenterController.astro: Import instead of redefine

### 8. CSS Token Additions

Add to BaseLayout.astro `:root`:
```css
--color-success: #10b981;
--color-error: #ef4444;
--color-secondary: #6b7280;
--border-radius: 4px;
--border-radius-lg: 8px;
```

Then replace hardcoded values in components.

### 9. Dead CSS Removal

Simply delete:
- `.loading-skeleton` and `@keyframes loading` from tags/[slug].astro
- `.temp-debug` rules from sync/[id]/control.astro

No functional impact since unused.

### 10. Dynamic Copyright Year

Replace `&copy; 2024` with `&copy; {new Date().getFullYear()}` in BaseLayout footer.

## Implementation Order

1. **Config & Infrastructure** (astro.config.mjs, BaseLayout CSS tokens)
2. **Shared Component** (Create ErrorState.astro)
3. **Page Updates** (Fix mains, apply ErrorState, remove console.logs)
4. **Component Updates** (Navigation a11y, API standardization)
5. **Cleanup** (Dead CSS removal, copyright year)

## Risk Mitigation

- Each change is isolated and testable
- No changes to component interfaces or props
- CSS classes preserved for styling continuity
- Build verification after each phase
- Changes are removal/fixes, not new functionality

## Verification Strategy

After implementation:
- `npm run build` must succeed
- HTML validator shows no nesting errors
- No console.logs in browser console
- Navigation shows current page indicator
- Tag links navigate correctly
- Error states render consistently