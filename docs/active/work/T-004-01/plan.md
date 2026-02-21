# Plan: Astro Layout and Navigation

## Implementation Steps

### Step 1: Create Directory Structure
**Actions:**
- Create `frontend/src/layouts/` directory
- Create `frontend/src/components/` directory

**Verification:**
- Directories exist at correct paths
- Proper permissions set

**Commit:** "feat: create layout and component directories"

### Step 2: Implement BaseLayout Component
**Actions:**
- Create `frontend/src/layouts/BaseLayout.astro`
- Add HTML document structure
- Define Props interface for title and description
- Add meta tags for SEO and viewport
- Implement CSS variables in `:root`
- Add global reset styles
- Create main content slot
- Add simple footer

**Verification:**
- File compiles without errors
- Props properly typed
- CSS variables accessible

**Commit:** "feat: implement BaseLayout component with global styles"

### Step 3: Implement Navigation Component
**Actions:**
- Create `frontend/src/components/Navigation.astro`
- Add header and nav structure
- Implement nav-brand link
- Add navigation links (Home, Stats)
- Implement mobile toggle with checkbox hack
- Add hamburger icon in label
- Style for desktop (horizontal)
- Add media query for mobile (vertical + toggle)

**Verification:**
- Component renders correctly
- Links have correct hrefs
- No JavaScript errors (pure CSS)

**Commit:** "feat: create responsive Navigation component"

### Step 4: Update Index Page
**Actions:**
- Modify `frontend/src/pages/index.astro`
- Import BaseLayout component
- Remove full HTML structure
- Wrap content in BaseLayout
- Pass title prop

**Verification:**
- Page renders with layout
- Navigation visible
- Footer present
- Existing content preserved

**Commit:** "refactor: update index page to use BaseLayout"

### Step 5: Test Mobile Responsiveness
**Actions:**
- Start dev server with `npm run dev`
- Test at various viewport widths
- Verify mobile menu toggle works
- Check navigation links function
- Ensure layout doesn't break

**Verification:**
- Mobile menu hidden by default
- Toggle button visible < 768px
- Menu shows/hides on toggle
- Desktop nav always visible >= 768px
- No horizontal scroll

**Commit:** "test: verify responsive behavior"

### Step 6: Build and Verify Production Output
**Actions:**
- Run `npm run build`
- Check `dist/` output
- Verify static HTML generated
- Ensure CSS properly scoped
- Check file sizes reasonable

**Verification:**
- Build completes without errors
- HTML files in dist/
- CSS inlined or linked properly
- No missing assets

**Commit:** "build: verify production build"

## Testing Strategy

### Manual Testing Checklist
- [ ] Desktop navigation displays horizontally
- [ ] Mobile toggle button appears < 768px
- [ ] Mobile menu opens/closes on toggle
- [ ] All navigation links work
- [ ] Footer stays at bottom
- [ ] Page content renders in main slot
- [ ] No horizontal scroll on mobile
- [ ] CSS variables apply correctly

### Browser Testing
- Chrome (latest)
- Firefox (latest)
- Safari (if available)
- Mobile Chrome/Safari

### Accessibility Checks
- Keyboard navigation works
- Links have proper focus states
- Mobile menu accessible via keyboard
- Semantic HTML structure

## Rollback Plan

If issues arise:
1. Git revert the problematic commit
2. Identify the issue in dev environment
3. Fix and re-apply changes
4. Re-test thoroughly

## Dependencies

### Required Before Start
- Astro dev server working (`npm run dev`)
- Write access to frontend/src/

### No External Dependencies
- No new npm packages
- No API integrations
- No database setup

## Success Metrics

### Acceptance Criteria Validation
- ✅ BaseLayout.astro provides HTML shell
- ✅ Navigation has Home and Stats links
- ✅ CSS scoped or minimal utility
- ✅ Responsive on mobile widths
- ✅ Reusable by future pages

### Performance Targets
- Page loads in < 1 second
- No layout shift on load
- CSS < 10KB (minified)
- No JavaScript for basic nav

## Risk Mitigation

### Potential Issues
1. **CSS Specificity Conflicts**
   - Mitigation: Use scoped styles
   - Fallback: Add more specific selectors

2. **Mobile Menu Not Working**
   - Mitigation: Test checkbox hack early
   - Fallback: Add minimal JS if needed

3. **Build Errors**
   - Mitigation: Test after each step
   - Fallback: Check Astro docs/errors

## Notes

- Keep initial implementation simple
- Avoid premature optimization
- Document any deviations in progress.md
- Future tickets will enhance styling