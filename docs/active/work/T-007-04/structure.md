# T-007-04 Structure: Presentation Viewer Page Implementation

## Executive Summary

This structure document defines the file-level changes, architecture, and component boundaries for implementing the `/watch/[id]` presentation viewer page. The implementation extends the proven sync viewer patterns with step-awareness capabilities, creating a robust audience-facing presentation viewer.

## 1. File-level Changes

### 1.1 Files to be Created

**Primary Implementation**
- `/frontend/src/pages/watch/[id].astro` (250+ lines)
  - Main presentation viewer page
  - Server-side presentation data fetching
  - Client-side PresentationViewer class
  - Step-aware progress display
  - Real-time SSE integration

**Supporting Components**
- `/frontend/src/components/StepIndicator.astro` (120+ lines)
  - Reusable step indicator component with pills/dots
  - Current step highlighting
  - Mobile-responsive layout
  - Accessibility features (ARIA labels, keyboard navigation)

### 1.2 Files to be Modified

**Existing Utilities (Minor Extensions)**
- `/frontend/src/utils/stepConversion.ts`
  - Add step display formatting utilities (if missing)
  - Enhance existing interfaces for viewer use cases
  - No breaking changes to existing functions

**Navigation Integration (Optional)**
- `/frontend/src/components/Navigation.astro`
  - Add viewer page to site navigation (if required)
  - Maintain existing navigation structure and styling

### 1.3 Files to be Deleted

None. This implementation adds new functionality without removing existing code.

## 2. Module Boundaries and Public Interfaces

### 2.1 Page Level (`/watch/[id].astro`)

**Public Interface:**
- URL Route: `/watch/[presentationId]`
- Query Parameters: None required
- Response: HTML page with embedded JavaScript

**Dependencies:**
- PresentationViewer class (internal)
- StepIndicator component
- stepConversion utilities
- BaseLayout component
- PocketBase realtime API

### 2.2 PresentationViewer Class

**Public Interface:**
```typescript
class PresentationViewer {
  constructor(
    presentationId: string,
    sessionId: string | null,
    presentationData: PresentationData,
    apiBase: string
  )

  // Core lifecycle methods
  connect(): Promise<void>
  disconnect(): void

  // State management
  updatePresentation(data: PresentationData): void

  // Event handling (private methods exposed for testing)
  handleSessionUpdate(progress: number): void
  handleConnectionStateChange(state: ConnectionState): void
}

interface PresentationData {
  id: string;
  name: string;
  step_count: number;
  step_labels?: string[];
  is_live: boolean;
  progress?: number;
  current_step?: number;
  active_session?: { id: string };
}
```

**Internal Boundaries:**
- Session management (SSE connection, reconnection logic)
- UI updates (DOM manipulation, visual transitions)
- Error handling (network errors, session state changes)

### 2.3 StepIndicator Component

**Public Interface:**
```typescript
interface StepIndicatorProps {
  currentStep: number;
  stepCount: number;
  stepLabels?: string[];
  className?: string;
  variant?: 'pills' | 'dots' | 'numbers';
  size?: 'small' | 'medium' | 'large';
}
```

**Data Flow:**
- Input: Current step index, total steps, optional labels
- Output: Visual step progression indicator
- Updates: Receives step changes from parent PresentationViewer

### 2.4 API Integration Boundaries

**Presentation Status API**
- Endpoint: `GET /api/presentations/:id/status`
- Purpose: Initial data fetching and session discovery
- Response: PresentationData interface
- Error Handling: HTTP status codes (404, 500, timeout)

**PocketBase Realtime API**
- Endpoint: `GET /api/realtime` (Server-Sent Events)
- Purpose: Live progress updates
- Message Filtering: `collection: 'sync_sessions'`, `record.id: sessionId`
- Connection Management: Automatic reconnection, error recovery

## 3. Internal Organization

### 3.1 Component Structure and Hierarchy

```
/watch/[id].astro
├── Server-side Frontmatter
│   ├── URL parameter extraction and validation
│   ├── API data fetching (/api/presentations/:id/status)
│   ├── Error state determination
│   └── Initial data preparation
├── HTML Template Structure
│   ├── BaseLayout wrapper
│   ├── Page header with presentation name
│   ├── StepIndicator component
│   ├── Progress section (inherited from sync viewer)
│   ├── Connection status indicators
│   ├── Error state displays
│   └── Loading states
└── Client-side Script
    ├── PresentationViewer class instantiation
    ├── DOM element binding
    ├── SSE connection setup
    └── Real-time update handlers
```

### 3.2 CSS Organization and Class Naming

**Class Naming Convention:**
```css
/* Page-level containers */
.presentation-viewer-container { }
.viewer-header { }
.viewer-content { }

/* Step-specific classes */
.step-indicator { }
.step-pill, .step-dot, .step-number { }
.step-pill.active, .step-pill.completed { }
.step-display-text { }

/* Progress section (inherited pattern) */
.progress-section { }
.progress-container { }
.progress-bar { }
.progress-label { }

/* Connection status (inherited) */
.connection-status { }
.status-indicator { }
.status-text { }

/* Error states */
.error-container { }
.error-title { }
.error-message { }
.not-live-message { }

/* Responsive modifiers */
.mobile-layout { }
.desktop-layout { }

/* State classes */
.loading { }
.connected { }
.connecting { }
.disconnected { }
.error { }
```

**CSS Architecture:**
- **Scoped Styles**: Component-level styles using Astro's scoped CSS
- **CSS Variables**: Inherit from BaseLayout design system
- **Mobile-first**: Responsive breakpoints starting from 320px
- **Accessibility**: High contrast, focus indicators, reduced motion support

### 3.3 JavaScript/TypeScript Organization

**File Structure within `/watch/[id].astro`:**
```typescript
// 1. Type definitions and interfaces
interface PresentationData { }
interface ConnectionState { }
enum ViewerState { }

// 2. PresentationViewer class definition
class PresentationViewer {
  // Properties
  private presentationId: string;
  private sessionId: string | null;
  private presentationData: PresentationData;
  private eventSource: EventSource | null;
  private connectionState: ConnectionState;

  // DOM references
  private elements: {
    stepIndicator: HTMLElement;
    progressBar: HTMLProgressElement;
    statusIndicator: HTMLElement;
    stepDisplay: HTMLElement;
  };

  // Core methods
  constructor() { }
  async connect() { }
  disconnect() { }

  // Event handlers
  private handleMessage(event: MessageEvent) { }
  private handleError(error: Event) { }
  private handleOpen() { }
  private handleClose() { }

  // UI update methods
  private updateStepDisplay(step: number, progress: number) { }
  private updateConnectionStatus(state: ConnectionState) { }
  private updateProgressBar(progress: number) { }

  // Utility methods
  private findDOMElements() { }
  private validatePresentationData() { }
  private formatLastUpdate(timestamp: Date) { }
}

// 3. Initialization code
document.addEventListener('DOMContentLoaded', () => {
  // Component instantiation and setup
});
```

### 3.4 Error Handling Patterns

**Error Categories:**
```typescript
enum ViewerErrorType {
  PRESENTATION_NOT_FOUND = 'presentation_not_found',
  PRESENTATION_NOT_LIVE = 'presentation_not_live',
  NETWORK_ERROR = 'network_error',
  SESSION_ERROR = 'session_error',
  TIMEOUT = 'timeout',
  INVALID_DATA = 'invalid_data'
}

const ERROR_MESSAGES = {
  presentation_not_found: {
    title: 'Presentation Not Found',
    message: 'This presentation does not exist or has been removed.',
    recovery: 'Return to presentations list'
  },
  presentation_not_live: {
    title: 'Presentation Not Live',
    message: 'This presentation is not currently being presented.',
    recovery: 'Check back when the presentation starts'
  },
  // ... additional error states
};
```

**Error Handling Strategy:**
- **Server-side**: Comprehensive error detection during initial fetch
- **Client-side**: Real-time connection error recovery
- **User-friendly**: Clear messaging with actionable recovery options
- **Graceful degradation**: Core functionality works without JavaScript
- **Logging**: Console logging for debugging (development only)

## 4. Architecture Integration Details

### 4.1 Routing Integration

**URL Structure:**
- Route Pattern: `/watch/[id]` where `[id]` is presentationId
- Parameter Extraction: `const { id } = Astro.params;`
- Validation: Non-empty string, reasonable length limits
- Error Handling: Invalid IDs redirect to 404 or error page

**Integration Points:**
- **Navigation**: Links from presenter dashboard (`/present/[id]`)
- **URL Sharing**: Direct sharing of viewer URLs for audience access
- **Deep Linking**: Support for bookmarking and direct access

### 4.2 PresentationViewer Architecture

**Extension of SyncViewer Patterns:**
```typescript
// Conceptual inheritance (not literal class inheritance)
class PresentationViewer implements SyncViewerInterface {
  // Inherited patterns:
  - SSE connection management
  - Progress bar updates
  - Connection status indicators
  - Error recovery mechanisms

  // Extensions:
  - Step computation and display
  - Presentation metadata handling
  - Session discovery logic
  - Enhanced error states
}
```

**Key Architectural Differences:**
- **Data Source**: Presentation API → Session ID discovery → SSE subscription
- **Display Logic**: Progress + Steps (dual display modes)
- **Error Scenarios**: Additional "not live" and "presentation not found" states
- **URL Pattern**: Uses presentation ID, not session ID

### 4.3 Step-aware Display Components

**Component Hierarchy:**
```
StepIndicator (Reusable Component)
├── Step Pills/Dots (Visual indicators)
├── Step Labels (Optional text labels)
├── Progress Context (Current/Total display)
└── Accessibility Features (ARIA attributes)

Step Display Text (Integrated into main page)
├── Current Step Calculation
├── Label Resolution
├── Progress Percentage
└── Live Update Handling
```

**Data Flow Architecture:**
```
Real-time Progress Update
│
├── progressToStep() calculation
├── formatStepDisplay() formatting
├── updateStepIndicator() visual update
└── updateProgressBar() inherited behavior
```

### 4.4 Mobile Responsiveness Architecture

**Breakpoint Strategy:**
```css
/* Mobile First - Base styles for 320px+ */
.step-indicator {
  flex-direction: column;
  align-items: center;
}

/* Tablet - 768px+ */
@media (min-width: 768px) {
  .step-indicator {
    flex-direction: row;
    justify-content: center;
  }
}

/* Desktop - 1024px+ */
@media (min-width: 1024px) {
  .step-indicator {
    gap: var(--space-lg);
  }
}
```

**Mobile Specific Features:**
- **Touch Optimization**: Larger tap targets, smooth scrolling
- **Performance**: Reduced animations, efficient DOM updates
- **Layout Adaptation**: Stacking on narrow screens, horizontal scroll for many steps
- **Accessibility**: Voice-over compatibility, high contrast support

## 5. Implementation Ordering and Dependencies

### 5.1 Phase 1: Core Infrastructure (4-6 hours)

**Dependencies:** None (uses existing infrastructure)

**Implementation Order:**
1. **Create page stub** (`/frontend/src/pages/watch/[id].astro`)
   - Basic Astro page structure
   - URL parameter extraction
   - BaseLayout integration
   - Basic HTML structure

2. **Implement presentation fetching**
   - API call to `/api/presentations/:id/status`
   - Error handling (timeout, not found, server error)
   - Data validation and sanitization
   - Server-side error state determination

3. **Add error state displays**
   - Error message components
   - "Not found" and "Not live" specific messaging
   - Navigation back to home/presentations list
   - Error state styling

4. **Create basic PresentationViewer class shell**
   - Constructor and basic properties
   - DOM element finding and binding
   - Basic connection management setup
   - Logging and debugging infrastructure

**Testing:** Basic page loads, error states display correctly, no JavaScript errors

### 5.2 Phase 2: Step Display Implementation (3-4 hours)

**Dependencies:** Phase 1 completion, stepConversion utilities

**Implementation Order:**
1. **Create StepIndicator component**
   - Basic pill/dot layout
   - Props interface and validation
   - CSS styling with design system variables
   - Responsive layout implementation

2. **Integrate step computation**
   - Import and use stepConversion utilities
   - Progress-to-step calculation
   - Step display formatting
   - Current step highlighting

3. **Add step display to main page**
   - StepIndicator component integration
   - Step text display ("Step 3 of 5 — Label")
   - Progress percentage alongside step information
   - ARIA live regions for screen readers

4. **Implement visual transitions**
   - CSS transitions for step changes
   - Smooth progress bar updates
   - Loading state animations
   - Reduced motion support

**Testing:** Step indicators display correctly, transitions work smoothly, mobile layout functions

### 5.3 Phase 3: Real-time Integration (2-3 hours)

**Dependencies:** Phase 2 completion, PocketBase realtime API

**Implementation Order:**
1. **Implement session discovery**
   - Extract active_session.id from presentation status
   - Validate session ID presence and format
   - Handle missing session (not live) scenarios
   - Session state change detection

2. **Add SSE connection management**
   - EventSource setup to `/api/realtime`
   - Message filtering for sync_sessions updates
   - Session ID matching logic
   - Connection state management (connecting, connected, disconnected)

3. **Connect real-time updates to UI**
   - Handle progress update messages
   - Trigger step recalculation
   - Update both progress bar and step indicators
   - Maintain smooth visual transitions

4. **Implement connection status display**
   - Inherit status indicators from sync viewer
   - Connection state visual feedback
   - Last update timestamp display
   - Reconnection attempt messaging

**Testing:** End-to-end presentation viewing, real-time updates, connection recovery

### 5.4 Phase 4: Polish and Optimization (2-3 hours)

**Dependencies:** Phase 3 completion

**Implementation Order:**
1. **Mobile responsiveness refinement**
   - Test on various screen sizes
   - Optimize touch interactions
   - Ensure readable typography
   - Performance optimization for mobile devices

2. **Accessibility enhancement**
   - Screen reader testing and optimization
   - Keyboard navigation support
   - ARIA label verification
   - High contrast mode compatibility

3. **Performance optimization**
   - Large step count handling (50+ steps)
   - Memory leak prevention
   - DOM update optimization
   - Connection cleanup on page unload

4. **Error scenario testing**
   - Network failure simulation
   - Invalid presentation ID handling
   - Session termination during viewing
   - API timeout scenarios

**Testing:** Full accessibility audit, performance testing, comprehensive error scenarios

### 5.5 Critical Integration Points

**Between Phases:**
- **Phase 1 → Phase 2**: Presentation data structure must be established
- **Phase 2 → Phase 3**: Step computation logic must be proven working
- **Phase 3 → Phase 4**: Real-time updates must be stable

**External Dependencies:**
- **PocketBase API**: Realtime endpoint availability and message format
- **Presentation API**: Status endpoint data structure consistency
- **Step Conversion**: Utility function compatibility and correctness
- **Design System**: CSS variables and component styling consistency

**Testing Integration Points:**
- **Manual Testing**: Real presenter → viewer workflow
- **Automated Testing**: Unit tests for step conversion logic
- **Performance Testing**: Large presentations, slow connections
- **Accessibility Testing**: Screen reader compatibility, keyboard navigation

## 6. Risk Mitigation and Quality Assurance

### 6.1 Implementation Risks

**High Risk:**
- Session discovery reliability from presentation API
- Real-time connection stability under network conditions
- Performance with presentations containing many steps

**Medium Risk:**
- Mobile layout complexity for step indicators
- Error state message clarity and user guidance
- Browser compatibility for SSE connections

**Low Risk:**
- CSS styling consistency (well-established design system)
- Step computation accuracy (proven utilities)
- Basic page routing and navigation

### 6.2 Quality Gates

**Each Phase:**
- Code review for architectural consistency
- Manual testing of new functionality
- No regression in existing features

**Final Implementation:**
- Full accessibility audit (WCAG 2.1 AA compliance)
- Performance benchmarking (load times, update latency)
- Cross-browser testing (Chrome, Firefox, Safari, Edge)
- Mobile device testing (iOS, Android)

### 6.3 Rollback Strategy

**File-level Rollback:**
- New files can be removed without affecting existing functionality
- No modifications to critical existing files
- Database schema unchanged

**Feature Flags:**
- Implementation can be behind feature flag if needed
- Progressive rollout to subset of presentations
- Easy disabling if issues discovered

## 7. Success Metrics and Acceptance Criteria

### 7.1 Functional Requirements

**Core Functionality:**
- ✅ Presentation viewer accessible at `/watch/[presentationId]`
- ✅ Displays presentation name, current step, and progress
- ✅ Real-time updates reflect presenter changes within 500ms
- ✅ Step indicators show current position in presentation
- ✅ "Not live" state for inactive presentations
- ✅ Error handling for invalid/missing presentations

**User Experience:**
- ✅ Professional, clean interface suitable for audience use
- ✅ Mobile-responsive design works on phones and tablets
- ✅ Loading states provide clear feedback
- ✅ Error messages include recovery actions

### 7.2 Technical Requirements

**Performance:**
- ✅ Initial page load under 2 seconds on 3G connection
- ✅ Real-time updates processed under 100ms
- ✅ Memory usage stable during long viewing sessions
- ✅ Graceful handling of presentations with 100+ steps

**Accessibility:**
- ✅ Screen reader compatible with proper ARIA attributes
- ✅ Keyboard navigation support
- ✅ High contrast mode compatibility
- ✅ Reduced motion preference support

### 7.3 Integration Requirements

**System Integration:**
- ✅ Seamless integration with existing PocketBase realtime API
- ✅ Consistent styling with existing design system
- ✅ Compatible with BaseLayout and Navigation components
- ✅ No breaking changes to existing sync viewer functionality

**API Compatibility:**
- ✅ Works with current presentation status API format
- ✅ Handles session lifecycle changes gracefully
- ✅ Maintains connection through network interruptions
- ✅ Supports all existing step conversion utilities

## Conclusion

This structure provides a comprehensive blueprint for implementing the presentation viewer page while maintaining architectural consistency with existing patterns. The phased approach allows for incremental development and testing, while the clear module boundaries ensure maintainable code that integrates seamlessly with the existing codebase.

The implementation leverages proven infrastructure extensively (SSE patterns, step conversion utilities, design system) while adding the necessary enhancements for a polished audience-facing presentation viewer experience.