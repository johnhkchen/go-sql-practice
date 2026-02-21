# Plan: Waiting Room Page (T-008-01)

## Implementation Steps

### Step 1: Configure Astro for Hybrid Rendering
**Goal**: Enable SSR capability while keeping existing pages static

**Actions**:
1. Install @astrojs/node adapter dependency
2. Update astro.config.mjs to hybrid mode with Node adapter
3. Test that build still works and existing pages remain static

**Verification**:
- `npm run build` succeeds
- Static pages (index.astro) still generate static HTML
- SSR capability available for dynamic pages

**Commit**: "feat: configure Astro for hybrid rendering to support dynamic routes"

### Step 2: Create Animation Styles
**Goal**: Establish reusable CSS animations for the waiting room

**Actions**:
1. Create `frontend/src/styles/animations.css`
2. Define pulse animation keyframes
3. Add accessibility considerations (prefers-reduced-motion)
4. Export animation utility classes

**Verification**:
- CSS file validates (no syntax errors)
- Animation respects accessibility preferences
- Smooth, subtle animation effect

**Commit**: "feat: add CSS animation utilities for waiting room indicators"

### Step 3: Create Basic Waiting Room Page Structure
**Goal**: Create the dynamic route file with basic structure

**Actions**:
1. Create directory `frontend/src/pages/watch/`
2. Create `frontend/src/pages/watch/[id].astro`
3. Add frontmatter with `prerender = false`
4. Import BaseLayout and basic template structure
5. Add parameter extraction for `id`

**Verification**:
- File structure correct
- Page renders without errors (even before API integration)
- Dynamic route recognized by Astro

**Commit**: "feat: create waiting room page structure with dynamic routing"

### Step 4: Implement Server-Side API Fetching
**Goal**: Fetch presentation data during SSR

**Actions**:
1. Add API_BASE environment variable handling
2. Implement fetch logic with error handling
3. Add timeout and abort controller for robustness
4. Parse and validate API response

**Verification**:
- Successful API calls return presentation data
- 404s are handled gracefully
- Network errors don't crash the page
- Server errors return meaningful messages

**Commit**: "feat: implement server-side presentation data fetching"

### Step 5: Add Waiting Room View
**Goal**: Render waiting room when active_session is null

**Actions**:
1. Add conditional rendering logic for waiting state
2. Create waiting room HTML structure
3. Add presentation name display
4. Include "waiting for presenter" message
5. Add animated pulse indicator
6. Import and apply animation styles

**Verification**:
- Correct display when active_session is null
- Presentation name shows correctly
- Animation plays smoothly
- Mobile responsive layout

**Commit**: "feat: implement waiting room view with animated indicator"

### Step 6: Add Error Handling Views
**Goal**: Handle and display various error states

**Actions**:
1. Create error view templates for different error types
2. Add 404 presentation not found view
3. Add network error view
4. Add server error view
5. Style error messages consistently

**Verification**:
- Invalid IDs show "Presentation not found"
- Network issues show connection error
- Server errors show appropriate message
- All error states are styled correctly

**Commit**: "feat: add error handling views for waiting room page"

### Step 7: Add Live Session Handling
**Goal**: Handle case when presentation is already live

**Actions**:
1. Add conditional logic for active_session existence
2. Create live view or redirect logic
3. Add appropriate messaging for live state
4. Ensure smooth transition between states

**Verification**:
- Live presentations show different view
- State transition logic works correctly
- No waiting room shown for live presentations

**Commit**: "feat: handle live presentation state in waiting room"

### Step 8: Add Progressive Enhancement
**Goal**: Optional JavaScript for live updates

**Actions**:
1. Add client-side script for periodic status checks
2. Implement automatic refresh when session goes live
3. Add cleanup for intervals on page unload
4. Ensure graceful degradation without JavaScript

**Verification**:
- Page works fully without JavaScript
- With JavaScript, automatically detects when presentation goes live
- No memory leaks from polling
- Smooth user experience transitions

**Commit**: "feat: add progressive enhancement with live status updates"

### Step 9: Styling and Mobile Optimization
**Goal**: Polish the visual design and ensure mobile responsiveness

**Actions**:
1. Apply consistent styling using design system variables
2. Add responsive breakpoints for mobile
3. Optimize animation performance
4. Test accessibility (color contrast, focus management)
5. Ensure touch-friendly design

**Verification**:
- Matches existing design system
- Readable on small screens (≥375px width)
- Animations perform well on mobile devices
- Accessible to screen readers
- Good color contrast ratios

**Commit**: "style: polish waiting room design and mobile responsiveness"

### Step 10: Integration Testing
**Goal**: Ensure all components work together correctly

**Actions**:
1. Test with real presentation data from PocketBase
2. Test all error scenarios (404, network, server errors)
3. Test both waiting and live states
4. Test progressive enhancement behavior
5. Cross-browser testing (Chrome, Firefox, Safari, mobile browsers)

**Verification**:
- All acceptance criteria met
- No console errors or warnings
- Consistent behavior across browsers
- Performance acceptable on mobile devices

**Commit**: "test: verify waiting room functionality and cross-browser compatibility"

## Testing Strategy

### Unit Testing
- **CSS Animations**: Visual testing of animation smoothness and accessibility
- **Error Handling**: Mock various API failure scenarios
- **State Logic**: Test waiting vs live state determination

### Integration Testing
- **API Integration**: Test with running PocketBase instance
- **Route Handling**: Test URL parameter extraction and validation
- **Progressive Enhancement**: Test with and without JavaScript enabled

### Manual Testing Scenarios

1. **Happy Path**:
   - Visit `/watch/valid-id` with waiting presentation
   - Verify waiting room displays correctly
   - Start presentation, verify it transitions

2. **Error Scenarios**:
   - Visit `/watch/invalid-id` → should show 404
   - Disconnect network, visit page → should show connection error
   - Test with API server down → should show server error

3. **Mobile Testing**:
   - Test on phone screens (375px, 414px widths)
   - Verify touch interactions work
   - Check animation performance

4. **Accessibility Testing**:
   - Test with screen reader
   - Verify keyboard navigation works
   - Check color contrast meets WCAG standards

### Performance Considerations

- **Server-Side Rendering**: Keep API fetch time under 2s with 5s timeout
- **Animation Performance**: Use transform/opacity only, avoid layout thrashing
- **Bundle Size**: Minimal JavaScript, leverage CSS animations
- **Loading Time**: Critical path CSS inlined, non-critical deferred

## Risk Mitigation

### Technical Risks

1. **API Timeout**: 5-second timeout with clear error message
2. **Network Issues**: Graceful degradation with retry suggestions
3. **Invalid Data**: Validate API responses before rendering
4. **Browser Compatibility**: Use standard CSS features, progressive enhancement

### User Experience Risks

1. **Confusing States**: Clear messaging for each state (waiting, live, error)
2. **Loading Delays**: Show immediate HTML render, then enhance
3. **Mobile Experience**: Test on real devices, not just dev tools
4. **Accessibility**: Include proper ARIA labels, respect motion preferences

## Dependencies and Coordination

### Internal Dependencies
- **T-004-01**: BaseLayout component (completed)
- **T-007-01**: Presentations collection structure (completed)

### External Dependencies
- PocketBase API running on localhost:8090
- Node.js runtime for SSR (development and production)

### No Coordination Required
- This is a new isolated page with no impact on existing functionality
- Uses established patterns from BaseLayout and Navigation

## Definition of Done

All acceptance criteria verified:
- [ ] `/watch/[id]` fetches presentation record on load
- [ ] When `active_session` is null, renders waiting room with:
  - [ ] Presentation name as heading
  - [ ] "Waiting for presenter to start..." status message
  - [ ] Subtle CSS animation (pulse/breathing effect)
  - [ ] Clean, centered layout - minimal and professional
- [ ] When `active_session` is set, skips waiting room and renders live view
- [ ] Shows "Presentation not found" for invalid IDs
- [ ] Works on mobile - responsive, readable on phone screens
- [ ] No JavaScript required for initial waiting room render
- [ ] Uses BaseLayout from T-004-01

Technical requirements:
- [ ] Code follows existing patterns and conventions
- [ ] No console errors or warnings
- [ ] Performance acceptable on mobile devices
- [ ] Accessible to screen readers
- [ ] Cross-browser compatible
- [ ] Commits are atomic and well-described