# T-007-04 Implementation Plan: Presentation Viewer Page

## Executive Summary

This plan sequences the implementation of the `/watch/[id]` presentation viewer page into ordered, verifiable tasks that can be executed and committed atomically. Building upon the completed research, design, and structure documents, this implementation extends the proven sync viewer patterns with step-awareness capabilities to create a robust audience-facing presentation viewer.

**Total Estimated Time:** 10-13 hours across 4 phases
**Implementation Approach:** Incremental development with atomic commits
**Risk Level:** Low (leveraging proven infrastructure)
**Key Dependencies:** T-006-04 (sync viewer patterns), T-007-02 (presentation API), existing step conversion utilities

## Implementation Phases Overview

| Phase | Focus | Duration | Deliverables | Key Risks |
|-------|-------|----------|--------------|-----------||
| **Phase 1** | Core Infrastructure | 4-6 hours | Basic page, API integration, error handling | Session discovery reliability |
| **Phase 2** | Step Display | 3-4 hours | StepIndicator component, step-aware UI | Mobile responsiveness complexity |
| **Phase 3** | Real-time Integration | 2-3 hours | SSE connection, live updates | Connection stability under load |
| **Phase 4** | Polish & Testing | 2-3 hours | Accessibility, performance, edge cases | Browser compatibility issues |

## Phase 1: Core Infrastructure Implementation

**Objective:** Establish foundation with presentation data fetching, error handling, and basic UI structure.

### Step 1.1: Create Basic Page Structure (1.5 hours)

**Task:** Create `/frontend/src/pages/watch/[id].astro` with core Astro page architecture.

**Implementation Details:**
- Server-side frontmatter with URL parameter extraction and validation
- Basic HTML structure with BaseLayout integration
- Error handling for invalid IDs (redirect to 404)
- Initial data structure setup for presentation and error states

**Verification Criteria:**
- [ ] Page loads at `/watch/test-id` without errors
- [ ] URL parameter extraction works correctly
- [ ] Invalid IDs redirect to 404
- [ ] Basic HTML structure renders

**Commit Message:** "Add basic page structure for presentation viewer"

### Step 1.2: Implement Presentation Data Fetching (1.5 hours)

**Task:** Add server-side presentation data fetching with comprehensive error handling.

**Implementation Details:**
- API call to `/api/presentations/:id/status` with timeout (5s)
- Error categorization (404, 500, timeout, network)
- Response validation for required fields (id, name, step_count)
- Initial data preparation for client-side hydration

**Verification Criteria:**
- [ ] Successfully fetches valid presentation data
- [ ] Handles 404 responses correctly
- [ ] Handles network timeouts gracefully
- [ ] Validates required presentation fields
- [ ] Server errors display appropriate messages

**Commit Message:** "Add presentation data fetching with error handling"

### Step 1.3: Create Error State Components (1 hour)

**Task:** Implement comprehensive error state displays with user-friendly messaging.

**Implementation Details:**
- Error containers for different error types (not found, timeout, server error)
- "Presentation Not Live" state for inactive presentations
- Recovery actions (retry button, return to home link)
- Mobile-responsive error layouts

**Verification Criteria:**
- [ ] "Presentation Not Found" displays for invalid IDs
- [ ] "Presentation Not Live" displays when `is_live: false`
- [ ] Network error states provide retry options
- [ ] Mobile-responsive error layouts
- [ ] Error actions (retry, return home) function correctly

**Commit Message:** "Add comprehensive error state handling and displays"

### Step 1.4: Create PresentationViewer Class Foundation (1.5 hours)

**Task:** Establish the PresentationViewer class with basic structure and DOM management.

**Implementation Details:**
- TypeScript interfaces for PresentationData and ConnectionState
- PresentationViewer class with constructor and basic methods
- DOM element finding and caching
- Session ID extraction from presentation data
- Basic connection status management

**Verification Criteria:**
- [ ] PresentationViewer class instantiates without errors
- [ ] Session ID extraction works correctly for live presentations
- [ ] DOM elements found and cached properly
- [ ] Connection status updates (placeholder functionality)
- [ ] Console logging provides helpful debugging information

**Commit Message:** "Add PresentationViewer class foundation with DOM management"

## Phase 2: Step Display Implementation

**Objective:** Implement step-aware UI components with visual progress indicators and responsive design.

### Step 2.1: Create StepIndicator Component (2 hours)

**Task:** Build a reusable StepIndicator component with pills/dots layout and accessibility features.

**Implementation Details:**
- Create `/frontend/src/components/StepIndicator.astro`
- Support for pills, dots, and numbers variants
- ARIA attributes for accessibility (progressbar role, aria-valuenow)
- Mobile-responsive layout with horizontal scroll
- Size variants (small, medium, large) and CSS transitions

**Verification Criteria:**
- [ ] StepIndicator renders with correct number of steps
- [ ] Active step highlighting works correctly
- [ ] All three variants (pills, dots, numbers) display properly
- [ ] Mobile responsive layout with horizontal scroll
- [ ] ARIA attributes provide proper accessibility
- [ ] High contrast and reduced motion preferences respected

**Commit Message:** "Add StepIndicator component with accessibility and responsive design"

### Step 2.2: Integrate Step Conversion Utilities (1 hour)

**Task:** Import and integrate existing step conversion utilities into the presentation viewer.

**Implementation Details:**
- Import progressToStep and formatStepDisplay from stepConversion.ts
- Add updateStepDisplay method to PresentationViewer class
- Implement step indicator visual state updates
- Add ARIA live regions for step announcements

**Verification Criteria:**
- [ ] Step conversion utilities import successfully
- [ ] Progress-to-step calculation matches expected results
- [ ] Step display text formats correctly with labels
- [ ] Step indicator visual state updates properly
- [ ] ARIA live regions announce step changes
- [ ] Progress bar and percentage display update

**Commit Message:** "Integrate step conversion utilities with real-time display updates"

### Step 2.3: Add Step-Aware HTML Template (1 hour)

**Task:** Enhance the main page template with step indicator and step-aware progress display.

**Implementation Details:**
- Add StepIndicator component to live presentation view
- Implement step display text with proper ARIA labeling
- Enhanced progress section with progress bar and percentage
- Connection status indicators inherited from sync viewer patterns
- Mobile-responsive layout with CSS Grid/Flexbox

**Verification Criteria:**
- [ ] Step indicator displays with correct current step
- [ ] Step display text shows proper formatting
- [ ] Progress bar renders with smooth visual updates
- [ ] Connection status indicator displays
- [ ] Mobile layout adapts properly
- [ ] All ARIA attributes provide accessibility

**Commit Message:** "Add step-aware HTML template with progress display and mobile responsiveness"

## Phase 3: Real-time Integration

**Objective:** Connect SSE real-time updates to step-aware UI components with robust connection management.

### Step 3.1: Implement Session Discovery and Validation (1 hour)

**Task:** Complete session ID discovery logic with validation and error handling.

**Implementation Details:**
- Enhanced constructor with session discovery from presentation data
- Session ID validation (format checking, presence verification)
- Presentation state change handling (live to not-live transitions)
- Not-live state display when presentation ends

**Verification Criteria:**
- [ ] Session ID extracted correctly from live presentations
- [ ] Invalid session IDs handled gracefully
- [ ] Presentation state changes trigger appropriate actions
- [ ] Not-live state displays when presentation ends
- [ ] Session ID validation prevents invalid connections

**Commit Message:** "Add robust session discovery and presentation state management"

### Step 3.2: Implement SSE Connection Management (1.5 hours)

**Task:** Add EventSource connection with message filtering and reconnection logic.

**Implementation Details:**
- EventSource connection to `/api/realtime`
- Message filtering for sync_sessions collection updates
- Exponential backoff reconnection logic (max 5 attempts)
- Connection state management (connected, connecting, reconnecting, disconnected)
- Cleanup on page unload and visibility change handling

**Verification Criteria:**
- [ ] SSE connection establishes successfully
- [ ] Message filtering works correctly for sync_sessions
- [ ] Progress updates trigger step display changes
- [ ] Reconnection logic handles network failures
- [ ] Connection cleanup works on page unload
- [ ] Connection status indicators update properly

**Commit Message:** "Implement SSE connection management with automatic reconnection"

### Step 3.3: Complete Real-time Step Updates (30 minutes)

**Task:** Finalize integration between real-time progress updates and step display components.

**Implementation Details:**
- Enhanced updateStepDisplay method with smooth transitions
- Step indicator transitions with CSS animations
- Screen reader announcements for step changes
- Progress bar synchronization with step updates
- Performance optimization with update throttling

**Verification Criteria:**
- [ ] Real-time progress updates trigger step changes
- [ ] Step transitions are smooth and visually appealing
- [ ] Step announcements work with screen readers
- [ ] Progress bar updates synchronize with step changes
- [ ] Reduced motion preferences are respected

**Commit Message:** "Complete real-time step updates with smooth transitions and accessibility"

## Phase 4: Polish and Comprehensive Testing

**Objective:** Finalize accessibility, performance optimization, and comprehensive error scenario testing.

### Step 4.1: Accessibility Enhancement and Validation (1 hour)

**Task:** Comprehensive accessibility audit and improvements.

**Implementation Details:**
- Enhanced ARIA attributes and live regions
- Keyboard navigation support (Space/Enter for reconnect)
- High contrast mode compatibility
- Focus management and screen reader announcements
- Color contrast verification (WCAG 2.1 AA standards)

**Verification Criteria:**
- [ ] Screen reader testing with NVDA/JAWS/VoiceOver
- [ ] Keyboard navigation works without mouse
- [ ] High contrast mode displays properly
- [ ] Focus indicators are visible and logical
- [ ] ARIA live regions announce changes appropriately
- [ ] Color contrast meets WCAG 2.1 AA standards

**Commit Message:** "Enhance accessibility with ARIA improvements and keyboard navigation"

### Step 4.2: Performance Optimization and Large Dataset Handling (1 hour)

**Task:** Optimize performance for presentations with many steps and long viewing sessions.

**Implementation Details:**
- Update throttling to prevent UI performance issues (max 20 updates/sec)
- Compact step indicator for presentations with 20+ steps
- Memory management for long viewing sessions
- Skip redundant updates when step hasn't changed
- Performance monitoring and cleanup routines

**Verification Criteria:**
- [ ] Smooth performance with 100+ step presentations
- [ ] UI remains responsive during rapid progress updates
- [ ] Memory usage stable during long viewing sessions
- [ ] Compact step indicator works for large presentations
- [ ] Update throttling prevents performance issues

**Commit Message:** "Add performance optimizations for large presentations and long sessions"

### Step 4.3: Comprehensive Error Scenario Testing (1 hour)

**Task:** Test and refine error handling for all edge cases and failure scenarios.

**Implementation Details:**
- Network failure recovery with manual retry options
- Session termination handling (presenter ends session)
- Invalid/malformed progress data validation
- Browser compatibility checks (EventSource support)
- Extended network failure handling

**Manual Testing Procedures:**
1. **Network Interruption:** Disconnect network during viewing, verify reconnection
2. **Invalid Presentation ID:** Test with non-existent presentation ID
3. **Presenter Stops Mid-Session:** Stop presentation while viewer is connected
4. **Large Step Count:** Test with presentations having 50+ steps
5. **Rapid Progress Changes:** Simulate very fast progress updates
6. **Mobile Responsiveness:** Test on various mobile devices and orientations
7. **Browser Compatibility:** Test on Chrome, Firefox, Safari, Edge
8. **Screen Reader Testing:** Full navigation with NVDA, JAWS, VoiceOver

**Verification Criteria:**
- [ ] All error scenarios handled gracefully
- [ ] Users always have clear recovery options
- [ ] No JavaScript errors in any failure scenario
- [ ] Performance acceptable under stress conditions
- [ ] Accessibility maintained in error states

**Commit Message:** "Add comprehensive error handling and browser compatibility checks"

## Testing Strategy

### Unit Testing Requirements

**Step Conversion Logic Testing:**
```typescript
// Tests for step conversion utilities integration
describe('Step Conversion Integration', () => {
  test('progressToStep matches backend calculation', () => {
    expect(progressToStep(0.0, 5)).toBe(0);
    expect(progressToStep(0.25, 5)).toBe(1);
    expect(progressToStep(0.5, 5)).toBe(2);
    expect(progressToStep(1.0, 5)).toBe(4);
  });

  test('formatStepDisplay with labels', () => {
    const result = formatStepDisplay(2, 5, ['Intro', 'Demo', 'Conclusion', 'Q&A', 'End']);
    expect(result).toContain('Step 3 of 5');
    expect(result).toContain('Conclusion');
  });
});
```

**PresentationViewer Class Testing:**
- Session ID extraction logic
- Connection state management
- Progress update throttling
- Error handling scenarios

### Integration Testing Requirements

**API Integration Tests:**
- `/api/presentations/:id/status` endpoint responses
- PocketBase realtime message filtering
- Session discovery and validation
- Error response handling

**Real-time Connection Tests:**
- EventSource connection establishment
- Message parsing and validation
- Reconnection logic under various network conditions
- Connection cleanup on page unload

**UI Integration Tests:**
- Step indicator updates with various step counts
- Progress bar synchronization
- Mobile responsive layout
- Accessibility features

### Manual Testing Procedures

**Functional Testing Checklist:**
- [ ] Load valid presentation ID - displays correctly
- [ ] Load invalid presentation ID - shows error
- [ ] Live presentation connects and updates in real-time
- [ ] Not-live presentation shows appropriate message
- [ ] Network interruption triggers reconnection
- [ ] Presenter stopping session updates viewer state
- [ ] Step changes display smoothly with transitions
- [ ] Progress bar updates match step progression

**Accessibility Testing Checklist:**
- [ ] Screen reader announces step changes
- [ ] Keyboard navigation works without mouse
- [ ] High contrast mode displays properly
- [ ] Focus indicators visible and logical
- [ ] ARIA attributes provide proper context
- [ ] Color contrast meets WCAG 2.1 AA

**Performance Testing Checklist:**
- [ ] Load time under 2 seconds on 3G
- [ ] Smooth updates with rapid progress changes
- [ ] Memory stable during 30+ minute sessions
- [ ] UI responsive with 100+ step presentations
- [ ] Mobile performance acceptable

**Cross-browser Testing:**
- [ ] Chrome (latest 2 versions)
- [ ] Firefox (latest 2 versions)
- [ ] Safari (latest 2 versions)
- [ ] Edge (latest 2 versions)
- [ ] Mobile Safari (iOS)
- [ ] Mobile Chrome (Android)

## Implementation Dependencies and Integration Points

**Existing Infrastructure Leveraged:**
- SSE connection patterns from `/sync/[id].astro`
- Step conversion utilities from `/utils/stepConversion.ts`
- Presentation API endpoints from T-007-02
- Design system variables from BaseLayout
- Error handling patterns from sync viewer

**New Components Created:**
- `PresentationViewer` class (step-aware version of SyncViewer)
- `StepIndicator.astro` component (reusable step progression display)
- Enhanced `/watch/[id].astro` page (extends existing stub)

### Key Code Changes

**Server-Side Data Flow** (Step 1):
```astro
// Replace basic API call with status endpoint
const response = await fetch(`${API_BASE}/api/presentations/${id}/status`);
const { is_live, progress, current_step, session_id, ...presentation } = await response.json();
```

**Client-Side SSE Integration** (Steps 3-4):
```typescript
class PresentationViewer {
  constructor(presentationId, sessionId, initialData) {
    this.presentationId = presentationId;
    this.sessionId = sessionId;
    this.currentStep = initialData.current_step;
    this.stepCount = initialData.step_count;
  }

  handleProgressUpdate(progress) {
    const newStep = progressToStep(progress, this.stepCount);
    if (newStep !== this.currentStep) {
      this.updateStepDisplay(newStep);
      this.announceStepChange(newStep);
    }
    this.updateProgressBar(progress);
  }
}
```

**Step Indicator Component** (Step 2):
```astro
<div class="step-indicators">
  {Array.from({ length: presentation.step_count }, (_, i) => (
    <div
      class={`step-dot ${i === current_step ? 'active' : ''}`}
      aria-label={`Step ${i + 1}`}
    />
  ))}
</div>
```

### File Modification Summary

**Single File Enhancement**: `/frontend/src/pages/watch/[id].astro`
- Server-side: Enhanced API integration (~30 lines changed)
- Client-side: SSE class implementation (~150 lines added)
- Templates: Step display components (~50 lines added)
- Styles: Step indicators and animations (~100 lines added)

**Dependencies**:
- Import from `/utils/stepConversion.ts` (already exists)
- Use existing SSE patterns from `/sync/[id].astro`
- Leverage existing BaseLayout and styling system

### Error Handling Strategy

**Network Errors**:
- Maintain existing timeout handling (5s)
- Add exponential backoff for SSE reconnection
- Graceful degradation when realtime unavailable

**Data Validation**:
- Validate step indices before display
- Handle missing or malformed session data
- Fallback to progress percentage if step calculation fails

**User Feedback**:
- Connection status prominently displayed
- Clear error messages for presentation issues
- Loading states during reconnection attempts

## Verification Plan

### Step-by-Step Verification

**After Step 1** (Enhanced API):
```bash
# Verify API integration
curl http://localhost:8090/api/presentations/{id}/status
# Check page renders with session data
curl http://localhost:3000/watch/{id}
```

**After Step 2** (Step Display):
- Visual inspection of step indicators
- Test with presentations of different step counts
- Verify mobile layout with browser dev tools

**After Step 3** (SSE Connection):
- Monitor Network tab for SSE connection
- Check console for connection events
- Verify status indicator updates

**After Step 4** (Live Updates):
- Use presenter dashboard to change steps
- Time delay between change and viewer update
- Verify smooth animations

**After Steps 5-6** (Enhancement & Polish):
- Lighthouse performance audit
- Accessibility testing with screen reader
- Cross-browser compatibility check

### Performance Acceptance Criteria

**Loading Performance**:
- Initial page load < 2 seconds
- Time to interactive < 3 seconds
- SSE connection established < 1 second

**Runtime Performance**:
- Step updates reflected < 2 seconds
- Smooth animations (60fps where supported)
- Memory usage stable during long sessions

**Accessibility Standards**:
- WCAG AA compliance for color contrast
- Screen reader compatibility
- Keyboard navigation support

## Risk Mitigation Strategies

### High-Risk Scenarios

**Session Discovery Failure:**
- **Risk:** Presentation API doesn't return active_session data
- **Mitigation:** Comprehensive validation and fallback to "not live" state
- **Detection:** Server-side validation in Astro frontmatter
- **Recovery:** Clear error messaging with manual refresh option

**SSE Connection Instability:**
- **Risk:** Real-time connection drops frequently
- **Mitigation:** Exponential backoff reconnection with manual retry
- **Detection:** Connection state monitoring and timeout detection
- **Recovery:** Visual indicators and user-initiated reconnection

**Performance Issues with Large Presentations:**
- **Risk:** UI becomes sluggish with 50+ steps
- **Mitigation:** Update throttling and compact indicator mode
- **Detection:** Performance monitoring and step count thresholds
- **Recovery:** Automatic fallback to simplified UI

### Rollback Strategy

**File-Level Rollback:**
- All new files can be removed without affecting existing functionality
- No modifications to critical existing files
- Database schema unchanged

**Feature Rollback Plan:**
1. **Quick Disable:** Comment out route in Astro config
2. **Partial Rollback:** Disable real-time features, keep static display
3. **Full Rollback:** Remove all new files, restore any modified files

**Data Rollback:**
- No database changes required
- No user data affected
- No migration rollback needed

## Verification and Acceptance Criteria

### Functional Requirements Verification

**Core Functionality:**
- [ ] ✅ Presentation viewer accessible at `/watch/[presentationId]`
- [ ] ✅ Displays presentation name, current step, and progress
- [ ] ✅ Real-time updates reflect presenter changes <500ms
- [ ] ✅ Step indicators show current position accurately
- [ ] ✅ "Not live" state for inactive presentations
- [ ] ✅ Error handling for invalid/missing presentations

**User Experience Requirements:**
- [ ] ✅ Professional, clean interface suitable for audience use
- [ ] ✅ Mobile-responsive design works on phones and tablets
- [ ] ✅ Loading states provide clear feedback
- [ ] ✅ Error messages include actionable recovery options
- [ ] ✅ Smooth visual transitions between steps

### Technical Requirements Verification

**Performance Metrics:**
- [ ] ✅ Initial page load <2 seconds on 3G connection
- [ ] ✅ Real-time updates processed <100ms
- [ ] ✅ Memory usage stable during long viewing sessions
- [ ] ✅ Graceful handling of presentations with 100+ steps

**Accessibility Compliance:**
- [ ] ✅ Screen reader compatible with proper ARIA attributes
- [ ] ✅ Keyboard navigation support
- [ ] ✅ High contrast mode compatibility
- [ ] ✅ Reduced motion preference support
- [ ] ✅ WCAG 2.1 AA compliance verified

### Integration Requirements Verification

**System Integration:**
- [ ] ✅ Seamless integration with existing PocketBase realtime API
- [ ] ✅ Consistent styling with existing design system
- [ ] ✅ Compatible with BaseLayout and Navigation components
- [ ] ✅ No breaking changes to existing sync viewer functionality

**API Compatibility:**
- [ ] ✅ Works with current presentation status API format
- [ ] ✅ Handles session lifecycle changes gracefully
- [ ] ✅ Maintains connection through network interruptions
- [ ] ✅ Supports all existing step conversion utilities

## Implementation Timeline

### Week 1: Foundation and Core Features
- **Days 1-2:** Phase 1 - Core Infrastructure (4-6 hours)
- **Days 3-4:** Phase 2 - Step Display Implementation (3-4 hours)
- **Day 5:** Phase 3 - Real-time Integration (2-3 hours)

### Week 2: Polish and Testing
- **Days 1-2:** Phase 4 - Accessibility and Performance (2-3 hours)
- **Days 3-4:** Comprehensive testing and bug fixes
- **Day 5:** Documentation and final verification

### Milestone Checkpoints

**End of Phase 1:**
- Basic page structure loads without errors
- Presentation data fetching works reliably
- All error states display appropriately
- PresentationViewer class foundation established

**End of Phase 2:**
- StepIndicator component displays correctly
- Step conversion utilities integrated
- Step-aware HTML template renders properly
- Mobile responsiveness verified

**End of Phase 3:**
- Real-time SSE connection established
- Progress updates trigger step display changes
- Connection management handles failures gracefully
- End-to-end presentation viewing works

**Final Completion:**
- All acceptance criteria verified
- Accessibility audit passed
- Performance requirements met
- Cross-browser compatibility confirmed

## Conclusion

This implementation plan provides a comprehensive roadmap for building the presentation viewer page while leveraging existing infrastructure and maintaining high code quality standards. The phased approach ensures each component can be developed, tested, and verified independently while building toward a robust, user-friendly presentation viewing experience.

The plan's emphasis on proven patterns, comprehensive testing, and accessibility ensures the final implementation will integrate seamlessly with the existing codebase while providing a polished, professional experience for presentation audiences.

**Key Success Factors:**
- Building on proven sync viewer and step conversion infrastructure
- Incremental development with atomic commits
- Comprehensive error handling and graceful degradation
- Strong focus on accessibility and mobile experience
- Thorough testing at each phase
- Clear rollback and recovery strategies

This structured approach minimizes implementation risk while ensuring all requirements are met with high quality and maintainability.