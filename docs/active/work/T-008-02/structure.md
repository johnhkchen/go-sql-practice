# T-008-02 Structure: Live Auto-Transition Implementation Blueprint

## Structure Overview

This document defines the file-level changes, component architecture, and implementation boundaries for implementing live auto-transition functionality. Based on the selected unified realtime component design, the implementation requires modifying one primary file and optionally extracting shared utilities for future reuse.

## File Modification Plan

### Modified Files

#### PRIMARY: `frontend/src/pages/watch/[id].astro`
**Change Type**: Major modification (replace polling with unified realtime component)
**Lines Affected**: ~150 lines (waiting room polling logic + component instantiation)
**Purpose**: Replace current polling + page reload with seamless SSE-based transitions

**Current Structure** (relevant sections):
```
Lines 1-50:    Astro frontmatter (SSR data fetching)
Lines 51-97:   Server-side state detection and error handling
Lines 98-153:  Live view HTML template
Lines 154-162: Waiting room HTML template
Lines 163-457: CSS styles (will add transition states)
Lines 459-727: PresentationViewer class (live state only - TO BE REPLACED)
Lines 729-758: Polling script for waiting room (TO BE REPLACED)
```

**New Structure** (after modification):
```
Lines 1-50:    Astro frontmatter (unchanged)
Lines 51-97:   Server-side state detection (unchanged)
Lines 98-153:  Live view HTML template (minimal updates for connection status)
Lines 154-162: Waiting room HTML template (minor updates for connection status)
Lines 163-500: CSS styles (add transition states and unified container)
Lines 501-900: PresentationAutoViewer class (NEW - replaces both PresentationViewer and polling)
Lines 901-920: Single component initialization (replaces both scripts)
```

### Created Files

#### OPTIONAL: `frontend/src/utils/realtime.js`
**Change Type**: New file (optional refactor for future reuse)
**Purpose**: Extract shared SSE connection management patterns
**Size Estimate**: ~100 lines

**Motivation**: While not required for T-008-02, extracting common SSE patterns would benefit future tickets that need realtime functionality. This is marked as optional since the primary goal is working auto-transition, not code organization.

### No Changes Required

**Backend Files**: No Go code changes needed
- PocketBase SSE endpoint already exists
- All required APIs already implemented
- Database schema unchanged

**Other Frontend Files**: No modifications to other pages
- `frontend/src/pages/sync/[id].astro` - unchanged
- `frontend/src/pages/present/[id].astro` - unchanged
- CSS files unchanged except for new transition states

## Component Architecture

### Core Class Structure

**PresentationAutoViewer Class** (new unified component):
```javascript
class PresentationAutoViewer {
  // Constructor and initialization
  constructor(presentationId, presentationData, apiBase)
  setupDOM()

  // Connection management
  connect()
  disconnect()
  handleConnectionError()
  enablePollingFallback()

  // Message handling
  handleMessage(event)
  handlePresentationUpdate(record)
  handleSessionUpdate(record)

  // State management
  transitionToLive(sessionId)
  transitionToWaiting()
  updateViewState(newState)

  // UI updates
  updateConnectionStatus(status)
  updateProgress(progress)
  showTransitionMessage(message)
  hideTransitionMessage()

  // Utility methods
  fetchCurrentProgress()
  reconcileState(serverState)
  cleanup()
}
```

### State Management Architecture

**State Enum**:
```javascript
const ViewState = {
  WAITING: 'waiting',
  STARTING: 'starting',
  LIVE: 'live',
  ENDING: 'ending'
};
```

**State Data Structure**:
```javascript
this.state = {
  viewState: 'waiting',           // Current UI state
  presentationId: string,         // Presentation being watched
  sessionId: string | null,       // Active session ID (null when waiting)
  connectionStatus: 'disconnected', // SSE connection status
  lastUpdate: Date,              // Timestamp of last state change
  transitionPromise: null        // Pending state transition
};
```

### Message Routing Architecture

**SSE Message Structure**:
```javascript
{
  collection: 'presentations' | 'sync_sessions',
  action: 'create' | 'update' | 'delete',
  record: {...} // Full record data
}
```

**Message Filtering Logic**:
```javascript
handleMessage(event) {
  const data = JSON.parse(event.data);

  // Presentation state changes (always monitor)
  if (data.collection === 'presentations' &&
      data.record?.id === this.presentationId) {
    this.routePresentationMessage(data);
  }

  // Session updates (only when live)
  if (data.collection === 'sync_sessions' &&
      this.state.viewState === 'live' &&
      data.record?.id === this.state.sessionId) {
    this.routeSessionMessage(data);
  }
}
```

## DOM Architecture and Updates

### Container Structure

**Current HTML Structure** (two separate conditional blocks):
```html
{isLive && (
  <div class="live-container">
    <h1 class="live-title">{presentation.name}</h1>
    <div class="connection-status">...</div>
    <div class="step-section">...</div>
    <div class="live-indicator">...</div>
  </div>
)}

{isWaiting && (
  <div class="waiting-container">
    <h1 class="waiting-title">{presentation.name}</h1>
    <p class="waiting-message">Waiting for presenter to start...</p>
    <div class="waiting-indicator animate-pulse"></div>
  </div>
)}
```

**New Unified Container Strategy**:
```html
<div class="presentation-container" data-state="waiting" id="presentation-container">
  <div class="connection-status" id="connection-status">
    <span class="status-indicator" id="status-indicator"></span>
    <span class="status-text" id="status-text">Connecting...</span>
    <span class="last-update" id="last-update"></span>
  </div>

  <div class="content-waiting" id="content-waiting">
    <h1 class="presentation-title" id="presentation-title">{presentation.name}</h1>
    <p class="waiting-message">Waiting for presenter to start...</p>
    <div class="waiting-indicator animate-pulse"></div>
  </div>

  <div class="content-starting" id="content-starting" style="display: none;">
    <h1 class="presentation-title">{presentation.name}</h1>
    <p class="transition-message">Presentation is starting...</p>
    <div class="transition-indicator"></div>
  </div>

  <div class="content-live" id="content-live" style="display: none;">
    <h1 class="presentation-title">{presentation.name}</h1>
    <div class="step-section" id="step-section">
      <h2 class="step-heading" id="step-heading">Step 1 of {presentation.step_count || 1}</h2>
      <div class="step-indicators" id="step-indicators">...</div>
      <div class="progress-container">...</div>
    </div>
    <div class="live-indicator">...</div>
  </div>

  <div class="content-ending" id="content-ending" style="display: none;">
    <h1 class="presentation-title">{presentation.name}</h1>
    <p class="transition-message">Presentation has ended</p>
    <div class="end-indicator"></div>
  </div>
</div>
```

### CSS State Management

**State-based Display Rules**:
```css
.presentation-container[data-state="waiting"] .content-waiting { display: block; }
.presentation-container[data-state="waiting"] .content-live { display: none; }
.presentation-container[data-state="starting"] .content-starting { display: block; }
.presentation-container[data-state="live"] .content-live { display: block; }
```

**Transition Animation Framework**:
```css
.presentation-container {
  transition: opacity 0.3s ease-in-out;
}

.presentation-container.transitioning {
  opacity: 0.7;
}
```

### DOM Update Boundaries

**Minimal DOM Manipulation Scope**:
- **Connection Status**: Single element, updated via textContent + className
- **Progress Display**: Single element, updated via textContent
- **Step Information**: Single element, updated via innerHTML
- **Container State**: Single data attribute change triggers CSS rules

**Preserved Elements**:
- Navigation and header elements (unchanged)
- Base layout and responsive structure (unchanged)
- Error message containers (reused existing)

## Interface Boundaries

### Public Interface (Component API)

**Initialization**:
```javascript
// Called from Astro page
window.presentationAutoViewer = new PresentationAutoViewer(
  presentationId,    // string
  presentationData,  // object from server-side rendering
  apiBase           // string
);
```

**Cleanup**:
```javascript
// Called on page unload
window.addEventListener('beforeunload', () => {
  window.presentationAutoViewer?.cleanup();
});
```

### Internal Module Boundaries

**SSE Management Module**:
- Connection establishment and teardown
- Reconnection logic with exponential backoff
- Message parsing and validation
- Connection status updates

**State Transition Module**:
- State machine enforcement
- Transition animations
- DOM updates
- Progress synchronization

**UI Update Module**:
- Connection status indicators
- Progress display
- Step navigation
- Accessibility announcements

### External Dependencies

**Required Utilities**:
```javascript
// Existing utility (frontend/src/utils/stepConversion.ts)
import { progressToStep, stepToProgress } from '../utils/stepConversion.ts';
```

**Browser APIs**:
- `EventSource` (SSE connection)
- `AbortController` (request timeouts)
- `localStorage` (connection state persistence - optional)

**PocketBase APIs**:
- `GET /api/realtime` (SSE stream)
- `GET /api/presentations/:id/status` (fallback data fetching)

## Implementation Roadmap

### Step 1: HTML Template Restructure
**Target**: Transform conditional rendering to unified container
**Changes**:
- Replace `{isLive && (...)}` and `{isWaiting && (...)}` with single `<div class="presentation-container">`
- Add state-based content divs (`content-waiting`, `content-starting`, `content-live`, `content-ending`)
- Move connection status outside conditional blocks
- Add consistent IDs for JavaScript access

### Step 2: CSS State Management
**Target**: Implement state-based display rules
**Changes**:
- Add CSS selectors for `[data-state="waiting|starting|live|ending"]`
- Implement transition animations between states
- Ensure responsive design works across all states
- Add transition indicators and messaging styles

### Step 3: Replace PresentationViewer Class
**Target**: Create PresentationAutoViewer with unified state management
**Key Methods to Implement**:
- `constructor(presentationId, initialState, apiBase)` - Initialize with current server state
- `transitionToState(newState)` - Handle state changes with UI updates
- `handlePresentationMessage(data)` - Process presentation start/end events
- `handleSessionMessage(data)` - Process progress updates during live state
- `setupStateTransition(fromState, toState)` - Manage transition animations

### Step 4: Replace Dual Script Blocks
**Target**: Single component initialization for all states
**Changes**:
- Remove separate `{isLive && <script>...}` and `{isWaiting && <script>...}` blocks
- Create single initialization that works regardless of initial state
- Pass initial state from server-side rendering to component
- Handle both SSE connections and polling fallbacks in one place

### Step 5: State Synchronization
**Target**: Handle mid-session joins and reconnections
**Changes**:
- Add state reconciliation when connecting mid-presentation
- Implement proper error handling for connection failures
- Add polling fallback when SSE is unavailable
- Ensure state consistency between server and client

### Step 6: Testing and Polish
**Target**: Production-ready implementation
**Changes**:
- Add comprehensive error boundaries
- Implement accessibility enhancements
- Add performance optimizations (debouncing, efficient DOM updates)
- Test cross-browser compatibility and mobile responsiveness

## Change Ordering and Dependencies

### Critical Path Dependencies
1. **HTML restructure MUST come first** - CSS and JS depend on new DOM structure
2. **CSS state rules MUST be in place** - Before JavaScript state transitions
3. **PresentationAutoViewer MUST be complete** - Before removing old scripts
4. **State synchronization MUST be tested** - Before deploying to production

### Parallel Development Opportunities
- CSS transition animations can be developed alongside JavaScript class
- Error handling and fallback logic can be implemented in parallel with core functionality
- Accessibility enhancements can be added after core state management works

### Risk Mitigation Order
1. **Preserve existing functionality** - Each step should maintain current behavior
2. **Add new capabilities incrementally** - Test state transitions before removing old code
3. **Implement fallbacks early** - Polling backup should be ready before removing old polling
4. **Comprehensive testing last** - Full integration testing after core implementation

## Data Flow Architecture

### Server-Side Data Flow (Unchanged)
```
Astro Page Request → API calls → Server-side rendering → HTML with initial data
```

### Client-Side Data Flow (New)
```
Page Load → PresentationAutoVivier init → SSE connection → Message handling → State updates → DOM updates
```

### State Synchronization Flow
```
PocketBase change → SSE message → Message filter → State transition → UI update → Accessibility announcement
```

### Error Handling Flow
```
Connection error → Retry logic → Fallback to polling → User notification → Recovery on reconnection
```

## Memory Management

### Object Lifecycle
- **PresentationAutoViewer**: Created on page load, cleaned up on page unload
- **EventSource**: Created during connect(), closed during disconnect()
- **DOM references**: Cached during setupDOM(), released during cleanup()

### Event Listener Management
```javascript
// Setup (during initialization)
this.eventSource = new EventSource(url);
this.eventSource.onmessage = this.boundMessageHandler;
window.addEventListener('beforeunload', this.boundCleanup);

// Cleanup (on page unload or component destruction)
this.eventSource.close();
window.removeEventListener('beforeunload', this.boundCleanup);
```

### Memory Leak Prevention
- Proper EventSource closure
- Event listener removal
- Clear interval/timeout cleanup
- Null reference assignments

## Testing Boundaries

### Unit Test Targets
- **PresentationAutoViewer class**: State transitions, message handling
- **State management**: Transition validation, edge case handling
- **DOM updates**: Verify correct elements updated
- **Error handling**: Connection failures, malformed messages

### Integration Test Boundaries
- **SSE connection**: End-to-end message flow
- **State synchronization**: Server changes → UI updates
- **Fallback behavior**: Polling activation on SSE failure
- **Multi-viewer**: Multiple tabs watching same presentation

### Manual Test Scenarios
- **Cross-device**: Mobile, tablet, desktop rendering
- **Network conditions**: Slow, fast, intermittent connections
- **Accessibility**: Screen reader compatibility, keyboard navigation
- **Edge cases**: Mid-session joins, rapid presenter actions

## Performance Considerations

### Optimization Boundaries
- **Message filtering**: Client-side filtering scope limited to essential checks
- **DOM updates**: Minimal scope, batch operations where possible
- **State changes**: Debounced to prevent excessive transitions
- **Memory usage**: Connection cleanup, reference management

### Performance Monitoring Points
- **Connection establishment time**: Track SSE connection speed
- **Message processing time**: Monitor filtering and update performance
- **Transition smoothness**: Measure animation frame rates
- **Memory usage**: Monitor for leaks over extended sessions

## Implementation Risk Mitigation

### High-Risk Areas
1. **State transition race conditions** → Promise-based sequential transitions
2. **EventSource connection failures** → Comprehensive error handling + polling fallback
3. **DOM manipulation performance** → Minimal scope updates, efficient selectors
4. **Memory leaks from persistent connections** → Thorough cleanup procedures

### Low-Risk Areas
- CSS styling (reusing existing patterns)
- Server-side rendering (unchanged)
- Backend APIs (no modifications needed)
- Mobile responsive design (existing framework)

## Rollback Strategy

### Implementation Safety
- **Minimal changes**: Single file modification with clear boundaries
- **Feature flag potential**: Can disable SSE and revert to polling
- **Backward compatibility**: All existing functionality preserved
- **Gradual deployment**: Can test with subset of presentations

### Rollback Triggers
- **SSE connection reliability** below acceptable threshold
- **Performance regressions** in page load or transition smoothness
- **Accessibility issues** not resolved within iteration
- **Cross-browser compatibility** problems

This structure provides a clear blueprint for implementing the live auto-transition functionality while maintaining system stability and providing clear boundaries for testing and rollback if needed.