# Structure: T-007-03 Presenter Dashboard

## Overview

This document defines the specific file-level changes, architectural boundaries, and implementation structure for the presenter dashboard system. Based on the design decision to extend existing sync control patterns, this structure maximizes code reuse while maintaining clear separation of concerns.

## File Modifications

### New Files to Create

#### Frontend Pages

**`frontend/src/pages/present/index.astro`** - Dashboard listing page
- **Purpose**: Main presenter dashboard with presentations list and creation form
- **Exports**: Default Astro page component
- **Dependencies**: BaseLayout, PocketBase API integration
- **Client Islands**: GoLive component for live session management
- **Size Estimate**: ~150 lines (similar to existing pages)

**`frontend/src/pages/present/[id].astro`** - Presenter control page
- **Purpose**: Step-by-step presentation control interface
- **Exports**: Default Astro page component with dynamic route handling
- **Dependencies**: BaseLayout, PresenterController client component
- **Server Logic**: Session validation, admin token verification, presentation data fetching
- **Size Estimate**: ~180 lines (based on existing control.astro pattern)

#### Frontend Components

**`frontend/src/components/GoLive.astro`** - Go Live interaction component
- **Purpose**: Client island for starting live presentations from dashboard
- **Exports**: Astro component with embedded client-side JavaScript
- **Public Interface**: `{ presentationId: string }` props
- **Internal Methods**: `handleGoLive()`, `handleError()`, `showLoading()`
- **Dependencies**: Fetch API, error handling utilities
- **Size Estimate**: ~100 lines

**`frontend/src/components/PresenterController.astro`** - Step navigation controls
- **Purpose**: Main control interface for step-by-step navigation
- **Exports**: Astro component with PresenterController JavaScript class
- **Public Interface**: `{ sessionId: string, adminToken: string, presentationData: object }` props
- **Internal Classes**: `PresenterController` (extends patterns from existing SyncController)
- **Methods**:
  - `updateStep(stepIndex)` - Convert step to progress and update API
  - `handlePrevious()` / `handleNext()` - Boundary-aware navigation
  - `handleStepJump(index)` - Direct step navigation
  - `handleSliderChange(progress)` - Fine-grained progress control
  - `stopPresenting()` - End session and return to dashboard
- **Dependencies**: Step-progress conversion utilities, throttled API updates
- **Size Estimate**: ~200 lines

#### Utility Modules

**`frontend/src/utils/stepConversion.ts`** - Step-progress conversion utilities
- **Purpose**: Centralized logic for step-to-progress mapping
- **Exports**:
  - `stepToProgress(stepIndex: number, stepCount: number): number`
  - `progressToStep(progress: number, stepCount: number): number`
  - `validateStepIndex(stepIndex: number, stepCount: number): boolean`
- **Implementation**: Uses formulas from backend (`presentations.go:95`)
- **Size Estimate**: ~40 lines

### Files to Modify

#### Navigation Component

**`frontend/src/components/Navigation.astro`** - Add presenter dashboard link
- **Modification**: Add "Present" navigation item after "Stats"
- **Changes**:
  - Line ~15: Add navigation link `<a href="/present">Present</a>`
  - Line ~25: Add mobile menu item with same link
- **Impact**: Minimal, follows existing pattern

#### Layout Integration

**`frontend/src/layouts/BaseLayout.astro`** - No changes required
- **Assessment**: Existing layout supports new pages without modification
- **Verification**: CSS custom properties and responsive breakpoints compatible

## Component Boundaries and Interfaces

### Public API Contracts

#### GoLive Component Interface
```typescript
interface GoLiveProps {
  presentationId: string;
}

interface GoLiveResponse {
  session_id: string;
  admin_url: string;
  viewer_url: string;
  step_count: number;
  step_labels: string[];
}
```

#### PresenterController Component Interface
```typescript
interface PresenterControllerProps {
  sessionId: string;
  adminToken: string;
  presentationData: {
    id: string;
    name: string;
    step_count: number;
    step_labels: string[];
    current_step: number;
    progress: number;
  };
}

interface StepNavigationState {
  currentStep: number;
  canGoPrevious: boolean;
  canGoNext: boolean;
  isUpdating: boolean;
}
```

### Internal Module Boundaries

#### Server-Side Data Flow
```
[Astro Page] → [PocketBase API] → [Component Props] → [Client Island]
```

- **Pages**: Handle authentication, data fetching, error states
- **Components**: Receive props, manage client-side interactions
- **API Integration**: Isolated in server-side page logic

#### Client-Side Architecture
```
[User Interaction] → [Component Event] → [API Update] → [UI Feedback]
```

- **Event Handling**: Component-level event listeners
- **API Communication**: Fetch-based with error handling
- **State Management**: Component-local state with DOM updates

## File Organization and Dependencies

### Directory Structure
```
frontend/src/
├── components/
│   ├── GoLive.astro              (new)
│   ├── Navigation.astro          (modified)
│   └── PresenterController.astro (new)
├── pages/
│   └── present/
│       ├── index.astro           (new)
│       └── [id].astro            (new)
└── utils/
    └── stepConversion.ts         (new)
```

### Dependency Graph
```
present/index.astro
├── BaseLayout.astro (existing)
└── GoLive.astro (new)
    └── stepConversion.ts (new)

present/[id].astro
├── BaseLayout.astro (existing)
└── PresenterController.astro (new)
    └── stepConversion.ts (new)
```

### External Dependencies
- **PocketBase API**: Built-in CRUD and custom presentation routes
- **Existing Patterns**: Error handling, timeout management, accessibility features
- **CSS Framework**: Component-scoped styles with existing custom properties

## Implementation Ordering

### Phase 1: Foundation Components
1. **stepConversion.ts** - Core utilities (no dependencies)
2. **GoLive.astro** - Simple client island (depends on utilities)
3. **Navigation.astro** - Add link (independent modification)

### Phase 2: Dashboard Page
4. **present/index.astro** - Dashboard listing (depends on GoLive component)

### Phase 3: Control Interface
5. **PresenterController.astro** - Step navigation (depends on utilities)
6. **present/[id].astro** - Control page (depends on PresenterController)

### Rationale for Ordering
- **Bottom-up dependency resolution**: Utilities before consumers
- **Independent testing**: Each component can be tested in isolation
- **Incremental functionality**: Dashboard works before control interface
- **Risk mitigation**: Complex control logic implemented last

## Data Flow Architecture

### Server-Side Rendering Flow
```
1. [Route Request] → [Astro Page Component]
2. [API Fetch] → [Data Transformation]
3. [Props Generation] → [Component Rendering]
4. [HTML + Islands] → [Client Hydration]
```

### Client-Side Interaction Flow
```
1. [User Action] → [Event Handler]
2. [State Validation] → [API Request]
3. [Response Handling] → [DOM Update]
4. [Error/Success Feedback] → [Accessibility Announcements]
```

### Error Propagation
```
[API Error] → [Component Error State] → [User Feedback] → [Recovery Options]
```

## Security Boundaries

### Authentication Layers
- **Page Level**: Server-side auth validation for presentation access
- **Component Level**: Admin token verification for control operations
- **API Level**: Backend validation with constant-time comparison

### Data Isolation
- **Admin Tokens**: Never exposed to client-side storage
- **Session Data**: Scoped to authenticated presentation owners
- **Public Data**: Presentations list available to all users

## Testing Integration Points

### Component Testing
- **GoLive**: API interaction, error handling, loading states
- **PresenterController**: Step navigation logic, boundary conditions
- **stepConversion**: Mathematical accuracy, edge cases

### Integration Testing
- **Navigation Flow**: Dashboard → Go Live → Control → Stop → Dashboard
- **API Integration**: Mock PocketBase responses for all endpoints
- **Error Scenarios**: Network failures, invalid tokens, missing data

### Accessibility Testing
- **Screen Readers**: All interactive controls properly announced
- **Keyboard Navigation**: Tab order and keyboard shortcuts
- **High Contrast**: Visual feedback in accessibility modes

This structure provides a clear blueprint for implementation while maintaining architectural consistency with the existing codebase. The modular approach allows for incremental development and testing, reducing integration risk.