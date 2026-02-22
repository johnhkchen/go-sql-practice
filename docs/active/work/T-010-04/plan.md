# Plan: Frontend HTML/CSS/A11y Cleanup (T-010-04)

## Implementation Steps

### Step 1: Fix Astro Configuration
**File**: `frontend/astro.config.mjs`
- Change `output: 'static'` to `output: 'hybrid'`
- **Verification**: Run `cd frontend && npm run build` - should succeed without SSR warnings

### Step 2: Add CSS Tokens to BaseLayout
**File**: `frontend/src/layouts/BaseLayout.astro`
- Add semantic color tokens and border-radius variables to :root
- Update copyright year to dynamic `{new Date().getFullYear()}`
- **Verification**: Inspect CSS variables in browser DevTools

### Step 3: Create ErrorState Component
**File**: `frontend/src/components/ErrorState.astro`
- Create new component with title, message, showRetry, backLink, backLabel props
- Include all error state styles within component
- **Verification**: Component file exists and has valid Astro syntax

### Step 4: Fix Navigation Accessibility
**File**: `frontend/src/components/Navigation.astro`
- Add `aria-label="Main navigation"` to nav element
- Replace checkbox hack with proper button element
- Add aria-expanded, aria-controls attributes
- Implement active link detection with aria-current="page"
- Add client-side script for menu toggle
- **Verification**:
  - Mobile menu toggles with button click
  - Active page has aria-current attribute
  - Accessibility tree shows proper ARIA labels

### Step 5: Remove Console.logs - Batch 1 (Pages)
**Files**:
- `frontend/src/pages/index.astro` - Remove 2 console.logs
- `frontend/src/pages/tags/[slug].astro` - Remove 2 console.logs
- `frontend/src/pages/links/[id].astro` - Remove console.logs
- `frontend/src/pages/sync/[id].astro` - Remove console.logs
- `frontend/src/pages/watch/[id].astro` - Remove console.logs
- **Verification**: `grep -r "console.log" frontend/src/pages/` returns no results

### Step 6: Fix Nested Main Elements
**Files**:
- `frontend/src/pages/index.astro` - Change `<main>` to `<div>` (line 97)
- `frontend/src/pages/tags/[slug].astro` - Change `<main>` to `<div>` (line 79)
- `frontend/src/pages/present/[id].astro` - Change `<main>` to `<div>` (line 201)
- `frontend/src/components/PresenterController.astro` - Change `<main>` to `<div>` (line 67)
- **Verification**: HTML validator shows no nested main elements

### Step 7: Fix Tag Navigation URLs
**File**: `frontend/src/pages/links/[id].astro`
- Line 100: Change `/search?tag=${tag}` to `/tags/${encodeURIComponent(tag)}`
- Line 363: Update script section URL similarly
- **Verification**: Click tag links, should navigate to /tags/[tagname]

### Step 8: Apply ErrorState Component to Index Page
**File**: `frontend/src/pages/index.astro`
- Import ErrorState component
- Replace error HTML (lines 108-117) with ErrorState component
- Remove error CSS (lines 162-224)
- **Verification**: Trigger error state, should display consistently

### Step 9: Apply ErrorState to Tags Page
**File**: `frontend/src/pages/tags/[slug].astro`
- Import ErrorState component
- Replace error HTML with component
- Remove duplicate error CSS
- Remove `.loading-skeleton` CSS (lines 329-341)
- **Verification**: Navigate to invalid tag, see error state

### Step 10: Remove Console.logs - Batch 2 (Components)
**Files**:
- `frontend/src/components/SearchInterface.astro`
- `frontend/src/components/StatsSummary.astro`
- `frontend/src/components/PresenterController.astro`
- **Verification**: `grep -r "console.log" frontend/src/components/` returns no results

### Step 11: Fix Component Issues
**Files**:
- `frontend/src/components/LinksList.astro` - Remove `role="main"` from div (line 98)
- `frontend/src/components/GoLive.astro` - Import API_BASE, use in fetch
- `frontend/src/components/PresenterController.astro` - Import API_BASE instead of redefine
- **Verification**: Components use consistent API configuration

### Step 12: Clean Sync Control Page
**File**: `frontend/src/pages/sync/[id]/control.astro`
- Remove `.temp-debug` CSS (lines 1426-1448)
- Remove any console.logs
- **Verification**: No temp-debug styles in rendered page

### Step 13: Apply ErrorState to Remaining Pages
**Files**:
- `frontend/src/pages/present/index.astro`
- `frontend/src/pages/present/[id].astro`
- `frontend/src/pages/links/[id].astro`
- Replace error HTML with ErrorState component where applicable
- Remove duplicate error CSS
- **Verification**: All error states render consistently

### Step 14: Final Build Verification
- Run `cd frontend && npm run build`
- Build should complete successfully
- **Verification**: Build succeeds, dist folder created

### Step 15: Runtime Verification
- Start dev server: `cd frontend && npm run dev`
- Check browser console - no console.logs
- Test navigation menu on mobile
- Click tag links - should route correctly
- Check active nav link highlighting
- Trigger error states - should display consistently
- **Verification**: All acceptance criteria met

## Testing Strategy

### Unit Testing
Not applicable - no new business logic, only UI cleanup

### Integration Testing
Manual testing required for:
1. Navigation menu functionality (mobile and desktop)
2. Tag link routing
3. Error state display
4. Active navigation highlighting

### Visual Regression
Check that pages still render correctly:
- Home page with links
- Tag page with filtered results
- Individual link page
- Presentation control pages
- Error states

### Accessibility Testing
1. Run axe DevTools on each page
2. Check ARIA attributes in navigation
3. Verify no nested main elements
4. Test keyboard navigation

## Rollback Plan

If issues discovered:
1. All changes are in frontend/ directory only
2. Can revert individual files as needed
3. No database or backend changes
4. CSS class names preserved for compatibility

## Success Criteria

- [ ] `npm run build` succeeds
- [ ] No console.log statements in browser console
- [ ] HTML validator shows no nested main elements
- [ ] Navigation has aria-current on active page
- [ ] Tag links navigate to /tags/ routes
- [ ] Error states use consistent component
- [ ] Mobile nav menu works with button (not checkbox)
- [ ] Copyright year is dynamic