# Design: T-007-03 Presenter Dashboard

## Overview

Based on the research, we need to implement a two-page presenter dashboard system that integrates with the existing PocketBase API and Astro frontend architecture. This design evaluates different approaches for the dashboard listing page (`/present`) and control interface (`/present/[id]`), making decisions grounded in the codebase reality.

## Design Options

### Option 1: Extend Existing Sync Control Pattern
**Approach**: Build new pages that mirror the existing `/sync/[id]/control.astro` architecture with server-side rendering and client islands.

**Advantages:**
- Leverages proven patterns from `control.astro:182` (SyncController class, throttled updates)
- Maintains consistency with existing error handling and accessibility features
- Reuses established API integration patterns (5s timeouts, user-friendly error messages)
- Aligns with current Astro conventions (TypeScript props, component-scoped CSS)

**Disadvantages:**
- Requires adapting sync-focused code to step-based navigation
- Step-to-progress conversion adds complexity to UI logic

**Assessment**: Strong candidate due to proven architecture and consistency with codebase.

### Option 2: Client-Heavy SPA Architecture
**Approach**: Build dashboard as client-side rendered components with extensive JavaScript state management.

**Advantages:**
- More responsive interactions for step navigation
- Easier real-time updates across multiple views

**Disadvantages:**
- Contradicts Astro's server-first philosophy
- Increases bundle size and complexity
- Breaks established patterns in the codebase
- Complicates SEO and initial load performance

**Assessment**: Rejected due to architectural inconsistency.

### Option 3: Hybrid Server/Client with New Component Patterns
**Approach**: Create new base components specifically for presenter functionality, separate from sync patterns.

**Advantages:**
- Clean separation between sync and presentation concerns
- Opportunity to optimize for step-based workflow from scratch

**Disadvantages:**
- Duplicates proven functionality from existing control interface
- Increases maintenance burden with parallel patterns
- Higher risk of introducing bugs in new patterns

**Assessment**: Rejected due to unnecessary complexity and code duplication.

## Chosen Approach: Extended Sync Control Pattern (Option 1)

**Rationale**: The existing sync control interface (`control.astro`) already demonstrates successful integration with the PocketBase API, admin token authentication, and real-time progress updates. The presenter dashboard requirements are essentially an enhanced version of this pattern with step-discrete navigation instead of continuous progress.

**Key Design Decisions:**

### 1. Page Architecture

**Dashboard Page** (`/present`):
- Server-side rendering with PocketBase presentations API integration
- Static form for creating new presentations
- Astro client island for "Go Live" interactions
- Follows existing page patterns with BaseLayout wrapper

**Control Page** (`/present/[id]`):
- Server-side session validation and data fetching
- Client island for step navigation controls
- Inherits security patterns from existing control interface

### 2. Step Navigation Interface

**Discrete Step Controls:**
- Previous/Next buttons with disabled states at boundaries
- Step jump buttons (numbered pills) for direct navigation
- Fine-control slider inherited from existing progress control
- Current step indicator with optional labels

**Progress Mapping Strategy:**
Use existing step-progress conversion formulas from `presentations.go:95`:
- Step → Progress: `progress = step_index / (step_count - 1)`
- Progress → Step: `step_index = round(progress * (step_count - 1))`

This maintains API compatibility while providing discrete step semantics.

### 3. Component Reuse Strategy

**Extend SyncController Pattern:**
- Create `PresenterController` class based on existing `SyncController`
- Inherit throttled API updates and error handling
- Add step-specific UI management methods
- Maintain accessibility and responsive design patterns

**Shared Utilities:**
- Reuse clipboard functionality for shareable URLs
- Inherit timeout handling and network error patterns
- Maintain consistent visual feedback systems

## Technical Design Specifications

### Dashboard Page Structure

```
/present
├── Server-side: Fetch presentations via PocketBase API
├── Render: Presentation list with status indicators
├── Form: New presentation creation (static)
└── Client Island: Go Live button interactions
```

**Data Flow:**
1. Server fetches presentations with `enhanced status` (research.md:44)
2. Component renders list with live/offline indicators
3. "Go Live" triggers `POST /api/presentations/:id/live`
4. Success redirects to `/present/[id]` with admin token

### Control Page Structure

```
/present/[id]
├── Server-side: Session validation and presentation data
├── Render: Step controls with current state
└── Client Island: Navigation controls
    ├── Previous/Next buttons
    ├── Step jump buttons
    ├── Fine-control slider
    └── Stop presenting functionality
```

**Navigation Logic:**
- Previous/Next: Calculate new step index, convert to progress, update API
- Step Jump: Direct step index to progress conversion
- Slider: Maintain existing continuous progress control
- All updates use existing `POST /api/sync/:id/progress` endpoint

### Security and Authentication

**Admin Token Flow** (inherited from research.md:157):
- Admin tokens passed via URL query parameter
- Constant-time comparison for security
- Session ownership validation before control access

**Presentation Ownership** (research.md:29):
- Creation requires authentication (`@request.auth.id != ''`)
- Control requires ownership match (`@request.auth.id = created_by`)

### Error Handling Strategy

**Network Errors** (following research.md:163):
- 5-second timeouts for API calls
- User-friendly error messages
- Graceful degradation for missing data
- Screen reader announcements for accessibility

**Validation Errors:**
- Step boundary enforcement (0 to step_count-1)
- Progress range validation (0.0 to 1.0)
- Admin token format verification

## Integration Points

### Navigation Updates
Add presenter dashboard link to `Navigation.astro` following existing patterns:
- Home, Stats, **Present** navigation structure
- Mobile hamburger menu compatibility

### API Endpoint Usage
- **List Presentations**: PocketBase built-in CRUD API (`"" public access`)
- **Create Presentation**: PocketBase CRUD with auth (`@request.auth.id != ''`)
- **Start Session**: `POST /api/presentations/:id/live` (custom route)
- **Update Progress**: `POST /api/sync/:id/progress` (existing endpoint)
- **Stop Session**: `POST /api/presentations/:id/stop` (custom route)

### URL Structure
- **Dashboard**: `/present`
- **Control**: `/present/[id]` (redirects from admin_url after go-live)
- **Admin Control**: `/sync/[session_id]/control?token=[admin_token]` (existing fallback)

## Risk Assessment and Mitigation

### Risk 1: Step-Progress Conversion Complexity
**Mitigation**: Extensive testing of boundary conditions and rounding behavior. Use existing formulas proven in backend logic.

### Risk 2: UI State Synchronization
**Mitigation**: Follow throttled update pattern from existing control interface. Single source of truth in API progress value.

### Risk 3: Mobile Usability
**Mitigation**: Inherit responsive patterns from existing control interface. Touch-friendly button sizing and layout.

## Success Criteria

1. **Functional Requirements Met**: All acceptance criteria implemented using established patterns
2. **Code Consistency**: New components follow existing Astro conventions and styling approaches
3. **Security Maintained**: Admin token handling and presentation ownership validation preserved
4. **Performance**: No regressions in page load time or interaction responsiveness
5. **Accessibility**: Screen reader compatibility and keyboard navigation maintained

This design leverages proven patterns while extending functionality in a maintainable way. The approach minimizes risk by building on successful existing code rather than introducing new architectural patterns.