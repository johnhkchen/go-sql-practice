# Plan: T-007-03 Presenter Dashboard

## Overview

This plan sequences the implementation of the presenter dashboard system into atomic, verifiable steps. Each step can be committed independently and tested in isolation. The order minimizes dependencies and allows for incremental verification of functionality.

## Implementation Steps

### Step 1: Foundation Utilities
**Objective**: Create step-progress conversion utilities with comprehensive testing
**Files**: `frontend/src/utils/stepConversion.ts`

**Tasks**:
- Implement `stepToProgress(stepIndex, stepCount)` function
- Implement `progressToStep(progress, stepCount)` function
- Implement `validateStepIndex(stepIndex, stepCount)` function
- Add TypeScript interfaces and documentation
- Handle edge cases (stepCount = 1, boundary values)

**Verification**:
- Unit tests for all conversion functions
- Boundary condition testing (0, 1, max values)
- Mathematical accuracy verification
- TypeScript compilation without errors

**Commit Point**: "Add step-progress conversion utilities with tests"

### Step 2: Navigation Component Update
**Objective**: Add presenter dashboard link to navigation
**Files**: `frontend/src/components/Navigation.astro`

**Tasks**:
- Add "Present" navigation link after "Stats" in desktop menu
- Add "Present" navigation link in mobile hamburger menu
- Verify responsive behavior and accessibility

**Verification**:
- Navigation link renders correctly in both desktop and mobile views
- Link points to correct `/present` route
- No regression in existing navigation functionality
- Mobile menu toggle still works properly

**Commit Point**: "Add presenter dashboard link to navigation"

### Step 3: GoLive Component
**Objective**: Create client island for starting live presentations
**Files**: `frontend/src/components/GoLive.astro`

**Tasks**:
- Create Astro component with TypeScript props interface
- Implement client-side JavaScript for "Go Live" button
- Add loading states and error handling
- Integrate with existing error feedback patterns
- Add accessibility announcements

**Verification**:
- Component renders with provided presentationId prop
- Button click triggers proper API call to `/api/presentations/:id/live`
- Loading state shows during API request
- Success redirects to `/present/[id]` route
- Error states display user-friendly messages
- Screen reader announcements work correctly

**Commit Point**: "Add GoLive component for starting presentations"

### Step 4: Dashboard Listing Page
**Objective**: Create main presenter dashboard at `/present`
**Files**: `frontend/src/pages/present/index.astro`

**Tasks**:
- Create Astro page with BaseLayout wrapper
- Implement server-side presentation fetching from PocketBase
- Render presentations list with status indicators
- Add "New Presentation" form with validation
- Integrate GoLive component for each presentation
- Add error handling for API failures

**Verification**:
- Page renders at `/present` route
- Presentations list displays correctly (name, step count, live status)
- New presentation form creates presentations via PocketBase API
- GoLive buttons work for each presentation
- Error states handled gracefully
- Page works without JavaScript (progressive enhancement)

**Commit Point**: "Add presenter dashboard listing page"

### Step 5: PresenterController Component
**Objective**: Create step navigation control interface
**Files**: `frontend/src/components/PresenterController.astro`

**Tasks**:
- Create Astro component with comprehensive props interface
- Implement PresenterController JavaScript class
- Add previous/next step navigation with boundary checking
- Implement step jump buttons (numbered pills)
- Add fine-control progress slider
- Integrate stop presenting functionality
- Add real-time UI feedback and error handling

**Verification**:
- Previous/next buttons work with proper boundary enforcement
- Step jump buttons navigate to correct steps
- Progress slider provides fine-grained control
- All navigation updates API via `/api/sync/:id/progress`
- Stop button ends session and redirects to dashboard
- Progress updates are throttled appropriately
- Error handling matches existing patterns

**Commit Point**: "Add PresenterController with step navigation"

### Step 6: Control Page
**Objective**: Create presenter control interface at `/present/[id]`
**Files**: `frontend/src/pages/present/[id].astro`

**Tasks**:
- Create dynamic Astro page with `[id]` parameter
- Implement server-side session validation and data fetching
- Add admin token verification and ownership checks
- Render PresenterController with proper props
- Handle error states (invalid session, missing token, etc.)
- Add shareable viewer URL with copy functionality

**Verification**:
- Page renders correctly for valid presentation IDs
- Admin token authentication works properly
- Session data loads and displays correctly
- PresenterController receives proper props
- Error pages for invalid/expired sessions
- Viewer URL generation and copy functionality
- Security boundaries properly enforced

**Commit Point**: "Add presenter control page with session management"

### Step 7: Integration Testing
**Objective**: End-to-end workflow testing and bug fixes
**Files**: Various (bug fixes as needed)

**Tasks**:
- Test complete workflow: Dashboard → Go Live → Control → Stop
- Verify API integration with all PocketBase endpoints
- Test error scenarios (network failures, invalid data)
- Verify responsive design on mobile devices
- Test accessibility with screen readers
- Performance testing (page load times, API response times)

**Verification**:
- Full presenter workflow works without issues
- All API endpoints respond correctly
- Error states provide clear user guidance
- Mobile experience is fully functional
- Accessibility requirements met
- Performance within acceptable ranges

**Commit Point**: "Fix integration issues and complete presenter dashboard"

## Verification Strategy

### Unit Testing Approach
- **stepConversion.ts**: Mathematical accuracy, boundary conditions, type safety
- **Component Logic**: Isolated testing of event handlers and state management
- **API Integration**: Mock responses for all PocketBase endpoints

### Integration Testing Approach
- **Page Navigation**: Dashboard → Control → Dashboard flow
- **API Workflow**: Create → Go Live → Control → Stop sequence
- **Error Scenarios**: Network failures, authentication errors, invalid data

### Manual Testing Checklist
- [ ] Dashboard displays presentations correctly
- [ ] New presentation creation works
- [ ] Go Live button starts sessions properly
- [ ] Control interface navigates steps correctly
- [ ] Stop presenting returns to dashboard
- [ ] Error states display helpful messages
- [ ] Mobile experience fully functional
- [ ] Accessibility features work with screen readers
- [ ] Performance acceptable on slow connections

### Testing Data Requirements
- **Test Presentations**: Various step counts (1, 3, 10 steps)
- **Test Users**: Authenticated and unauthenticated scenarios
- **Test Sessions**: Active and expired sessions for error testing

## Risk Mitigation

### Risk 1: Step-Progress Conversion Edge Cases
**Mitigation**: Comprehensive unit testing of mathematical functions first
**Verification**: Boundary value testing (0, 1, max steps)
**Fallback**: Revert to continuous progress if conversion fails

### Risk 2: API Integration Complexity
**Mitigation**: Incremental integration starting with simple endpoints
**Verification**: Mock testing before real API integration
**Fallback**: Graceful degradation with error messages

### Risk 3: Mobile Usability Issues
**Mitigation**: Progressive enhancement from desktop to mobile
**Verification**: Testing on actual mobile devices throughout development
**Fallback**: Desktop-only release if mobile issues persist

### Risk 4: Performance Degradation
**Mitigation**: Follow existing patterns for API throttling and caching
**Verification**: Performance monitoring during development
**Fallback**: Code splitting or lazy loading if bundle size grows

## Dependencies and Prerequisites

### External Dependencies
- PocketBase server running with existing presentation/sync APIs
- Astro development server for frontend
- Authentication system for presentation ownership

### Code Dependencies
- BaseLayout component (existing)
- Navigation component (existing)
- Existing API error handling patterns
- CSS custom properties system

### Data Dependencies
- Presentations collection schema (implemented in T-004-01)
- Sync sessions collection schema (implemented in T-007-02)
- Authentication system with user management

## Success Criteria

### Functional Requirements
- [ ] All acceptance criteria from ticket met
- [ ] Dashboard lists presentations with correct status
- [ ] New presentation creation works
- [ ] Go Live functionality starts sessions
- [ ] Step navigation works in all modes (buttons, jumps, slider)
- [ ] Stop presenting returns to dashboard
- [ ] Shareable viewer URLs generated correctly

### Technical Requirements
- [ ] Code follows existing Astro patterns and conventions
- [ ] TypeScript compilation without errors
- [ ] No accessibility regressions
- [ ] Mobile responsive design maintained
- [ ] API integration secure and robust
- [ ] Performance within acceptable ranges

### Quality Requirements
- [ ] All unit tests pass
- [ ] Integration tests verify complete workflow
- [ ] Manual testing checklist completed
- [ ] Code review standards met
- [ ] Documentation updated as needed

This plan provides a structured path to implementation with clear verification points and risk mitigation strategies. Each step builds incrementally toward the complete presenter dashboard system.