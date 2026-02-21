# Design: T-008-02 - live-auto-transition

## Executive Summary

This design implements seamless realtime auto-transition from waiting room to live presentation view in `/watch/[id]` pages. The solution enhances the existing dual-state architecture with unified SSE realtime monitoring, eliminating page reloads while maintaining backward compatibility and progressive enhancement.

## Design Overview

Based on the research findings, this design document evaluates approaches for implementing seamless realtime transitions from waiting room to live presentation view. The core challenge is replacing the current polling + page reload pattern with a unified SSE approach that maintains persistent connections and provides smooth state transitions.

## Design Constraints and Requirements

### Functional Requirements (from Acceptance Criteria)
1. **Astro client island** (`client:load`) subscribes to `presentations` collection via PocketBase realtime
2. **Automatic transition** when `active_session` changes from null to session ID (seamless, no reload)
3. **Live session subscription** begins upon transition (inheriting T-007-04 patterns)
4. **Transition states** with user feedback ("Presentation starting...")
5. **Reverse transition** when presentation ends (`active_session` → null)
6. **Connection status indicator** visible during waiting
7. **Edge case handling** (mid-session joins, network issues)
8. **Visual smoothness** (fade/slide transitions, not jarring swaps)

### Technical Constraints (from Research)
- **Single SSE Connection**: One EventSource per page, client-side message filtering
- **Progressive Enhancement**: Core functionality must work without JavaScript (polling fallback)
- **Mobile Responsive**: Touch-friendly UI across device sizes
- **Astro Architecture**: Server-side rendering with client hydration
- **PocketBase Dependency**: All realtime features depend on PocketBase SSE endpoint

## Design Approaches Evaluated

### Approach 1: Unified Realtime Component (SELECTED)

**Description**: Create a new `PresentationAutoViewer` class that combines waiting room and live viewer functionality within a single realtime-aware component. This class subscribes to presentations collection immediately and manages all state transitions internally.

**Architecture**:
```javascript
class PresentationAutoViewer {
  constructor(presentationId, initialData, apiBase) {
    this.state = initialData.is_live ? 'live' : 'waiting';
    this.connect(); // Start SSE immediately
  }

  connect() {
    // Single EventSource for entire session
    this.eventSource = new EventSource(`${this.apiBase}/api/realtime`);
    this.eventSource.onmessage = (event) => this.handleMessage(event);
  }

  handleMessage(event) {
    const data = JSON.parse(event.data);

    // Handle presentation state changes
    if (data.collection === 'presentations' && data.record.id === this.presentationId) {
      this.handlePresentationUpdate(data.record);
    }

    // Handle session progress updates (when live)
    if (this.state === 'live' && data.collection === 'sync_sessions') {
      this.handleSessionUpdate(data.record);
    }
  }

  handlePresentationUpdate(presentation) {
    const wasLive = this.state === 'live';
    const isLive = presentation.active_session !== null;

    if (!wasLive && isLive) {
      this.transitionToLive(presentation.active_session);
    } else if (wasLive && !isLive) {
      this.transitionToWaiting();
    }
  }
}
```

**Advantages**:
- **Single Connection**: One SSE stream handles all states
- **Smooth Transitions**: DOM manipulation instead of page reloads
- **Reuses Existing Patterns**: Builds on proven T-007-04 SSE implementation
- **Connection Consistency**: Status indicator works across all states
- **Edge Case Friendly**: Handles mid-session joins naturally

**Disadvantages**:
- **Complexity**: Manages multiple states in single component
- **Memory Usage**: Maintains subscription even when not immediately needed
- **Testing Complexity**: Multiple state transitions need comprehensive testing

### Approach 2: Hybrid Polling + SSE

**Description**: Keep polling for presentation state changes but enhance it with faster intervals and better UX. Only initialize SSE connection after detecting live state.

**Architecture**:
```javascript
// Waiting state: Enhanced polling
setInterval(async () => {
  const status = await fetch(`/api/presentations/${id}/status`);
  if (status.is_live && !this.isLive) {
    this.transitionToLive();
  }
}, 1000); // Faster polling (1s instead of 5s)

// Live state: Full SSE (existing pattern)
class PresentationViewer { /* existing implementation */ }
```

**Advantages**:
- **Lower Complexity**: Smaller code changes to existing implementation
- **Proven Reliability**: Polling is robust against network issues
- **Resource Efficient**: No persistent connection during waiting
- **Backward Compatible**: Works without SSE support

**Disadvantages**:
- **Delayed Detection**: 1-second delay in best case, longer with network latency
- **Server Load**: Continuous polling from all waiting viewers
- **User Experience**: Visible delays in transition, not truly "seamless"
- **Battery Impact**: Mobile devices drain faster with constant polling

### Approach 3: Separate Component Architecture

**Description**: Create two distinct components (`WaitingRoomViewer`, `LivePresentationViewer`) and swap them dynamically based on state changes. Each component is optimized for its specific role.

**Architecture**:
```javascript
class PresentationPageManager {
  constructor(presentationId, initialData, apiBase) {
    this.currentViewer = null;
    this.initializeViewer(initialData.is_live);
    this.subscribeToStateChanges();
  }

  initializeViewer(isLive) {
    if (isLive) {
      this.currentViewer = new LivePresentationViewer(/*...*/);
    } else {
      this.currentViewer = new WaitingRoomViewer(/*...*/);
    }
  }

  switchToLive(sessionId) {
    this.currentViewer.destroy();
    this.currentViewer = new LivePresentationViewer(/*...*/);
  }
}
```

**Advantages**:
- **Separation of Concerns**: Each component focused on single responsibility
- **Code Reuse**: Can reuse existing PresentationViewer class unchanged
- **Maintainability**: Easier to modify individual states independently
- **Memory Optimization**: Only load what's needed for current state

**Disadvantages**:
- **Connection Overhead**: Need to coordinate SSE connections between components
- **Transition Complexity**: Component swapping may cause visual jumps
- **State Management**: Need additional layer to manage component lifecycle
- **Code Duplication**: Shared functionality (connection status, error handling) duplicated

## Selected Approach: Unified Realtime Component

After evaluating the approaches against the research findings and acceptance criteria, **Approach 1: Unified Realtime Component** is selected for the following reasons:

### Decision Rationale

1. **Acceptance Criteria Alignment**:
   - ✅ "seamlessly swaps in the live view" - unified component enables smooth DOM transitions
   - ✅ "no full reload" - single component maintains persistent state
   - ✅ "connection status indicator visible during waiting" - single connection simplifies status management

2. **Technical Foundation Match**:
   - Builds directly on proven T-007-04 PresentationViewer patterns
   - Leverages existing SSE infrastructure without architectural changes
   - Maintains single EventSource pattern established in codebase

3. **User Experience Priority**:
   - Eliminates all polling delays for truly "seamless" transitions
   - Provides consistent connection status across all states
   - Handles edge cases (mid-session joins, network reconnects) naturally

4. **Performance Benefits**:
   - Single SSE connection reduces server load compared to polling
   - Immediate state change detection (no 1-5 second delays)
   - Better mobile battery life (no polling intervals)

### Rejected Approaches

**Approach 2 (Hybrid)** rejected due to:
- Fails "seamless" requirement with inherent polling delays
- Doesn't improve current server load from polling
- Maintains status quo rather than leveraging realtime capabilities

**Approach 3 (Separate Components)** rejected due to:
- Increased implementation complexity without clear benefits
- Risk of visual jumps during component swapping
- Doesn't align with existing single-component patterns in codebase

## Detailed Design Specification

### Component Architecture

**File**: `frontend/src/pages/watch/[id].astro` (modify existing)

**Core Class**: `PresentationAutoViewer` extends the existing PresentationViewer pattern

```javascript
class PresentationAutoViewer {
  constructor(presentationId, presentationData, apiBase) {
    this.presentationId = presentationId;
    this.presentationData = presentationData;
    this.apiBase = apiBase;

    // State management
    this.viewState = presentationData.is_live ? 'live' : 'waiting';
    this.sessionId = presentationData.active_session;
    this.connectionStatus = 'disconnected';

    // UI references
    this.setupDOM();

    // Start realtime connection immediately
    this.connect();
  }
}
```

### State Machine Design

**States**:
1. **waiting**: `active_session` is null, shows waiting room UI
2. **starting**: Intermediate state during transition (shows "Presentation starting...")
3. **live**: `active_session` exists, shows presentation progress UI
4. **ending**: Intermediate state when presentation stops (shows "Presentation ended")

**Transitions**:
```javascript
waiting → starting → live → ending → waiting
```

**State Actions**:
- **waiting**: Subscribe to presentations, show connection status, display waiting UI
- **starting**: Show transition message, prepare live UI, fetch session data
- **live**: Subscribe to sync_sessions, update progress, handle presenter actions
- **ending**: Show ending message, cleanup session subscription, offer re-wait option

### SSE Message Handling Strategy

**Dual Collection Subscription**:
```javascript
handleMessage(event) {
  const data = JSON.parse(event.data);

  // Always listen for presentation changes
  if (data.collection === 'presentations' &&
      data.record.id === this.presentationId) {
    this.handlePresentationUpdate(data.record);
  }

  // Only listen for session changes when live
  if (this.viewState === 'live' &&
      data.collection === 'sync_sessions' &&
      data.record.id === this.sessionId) {
    this.handleSessionUpdate(data.record);
  }
}
```

**Presentation Update Logic**:
```javascript
handlePresentationUpdate(presentation) {
  const wasActive = this.sessionId !== null;
  const isActive = presentation.active_session !== null;

  if (!wasActive && isActive) {
    // waiting → starting → live
    this.initiateTransitionToLive(presentation.active_session);
  } else if (wasActive && !isActive) {
    // live → ending → waiting
    this.initiateTransitionToWaiting();
  }

  // Update stored presentation data
  this.presentationData = presentation;
}
```

### DOM Transition Implementation

**Container Strategy**: Single container with state-based content swapping
```html
<div class="presentation-container" data-state="waiting">
  <!-- Content swapped based on state -->
</div>
```

**Transition Sequence**:
1. **Fade out current content** (300ms CSS transition)
2. **Update DOM structure** for new state
3. **Update data attributes** for styling
4. **Fade in new content** (300ms CSS transition)

**CSS State Management**:
```css
.presentation-container[data-state="waiting"] .live-content { display: none; }
.presentation-container[data-state="live"] .waiting-content { display: none; }
.presentation-container[data-state="starting"] .transition-message { display: block; }
```

### Connection Status Integration

**Unified Status Display**: Same connection status component visible in all states

```javascript
updateConnectionStatus(status) {
  this.connectionStatus = status;
  const indicator = document.querySelector('.connection-status');
  indicator.className = `connection-status ${status}`;
  indicator.textContent = {
    'connected': 'Connected',
    'connecting': 'Connecting...',
    'disconnected': 'Disconnected'
  }[status];
}
```

**Status Positioning**:
- **Waiting State**: Top-right corner, non-intrusive
- **Live State**: Integrated with existing UI (same location as current implementation)
- **Transition States**: Maintain visibility during transitions

### Error Handling and Resilience

**Fallback Strategy**: Progressive degradation to polling if SSE fails
```javascript
handleConnectionError() {
  this.retryCount++;
  if (this.retryCount > MAX_RETRIES) {
    console.warn('SSE failed, falling back to polling');
    this.enablePollingFallback();
  } else {
    this.scheduleReconnect();
  }
}

enablePollingFallback() {
  this.pollingInterval = setInterval(() => {
    this.checkPresentationStatus();
  }, 5000); // Graceful degradation to original polling
}
```

**Network Resilience**:
- **Connection Loss**: Show "reconnecting" status, attempt reconnection
- **Rapid Changes**: Debounce state transitions to prevent UI flicker
- **Race Conditions**: Use promise-based transitions to ensure sequential state changes

### Edge Case Handling

**Mid-Session Join**:
```javascript
// If user opens page while presentation is live
if (this.viewState === 'live') {
  this.fetchCurrentProgress().then(progress => {
    this.initializeLiveView(progress);
  });
}
```

**Network Disconnect/Reconnect**:
```javascript
handleReconnect() {
  // Re-sync current state from server
  this.fetchPresentationStatus().then(status => {
    this.reconcileState(status);
  });
}
```

**Presenter Actions During Transition**:
- **Queue state changes** if transition in progress
- **Debounce rapid changes** (presenter starts/stops quickly)
- **Maintain user feedback** during any delays

## Visual Design Considerations

### Transition Animation Strategy

**Fade Transition** (selected over slide for reliability):
```css
.presentation-container {
  transition: opacity 0.3s ease-in-out;
}

.presentation-container.transitioning {
  opacity: 0.7;
}
```

**Rationale**: Fade transitions are:
- Less prone to layout issues across different content heights
- More accessible (respects reduced motion preferences)
- Simpler to implement reliably across mobile/desktop

### Loading and Transition States

**"Presentation starting..." State**:
- Gentle pulsing animation
- Brief duration (should resolve within 1-2 seconds)
- Clear messaging about what's happening

**"Presentation ended" State**:
- Option to "Keep waiting for next session"
- Brief positive confirmation of end
- Smooth return to waiting state

### Mobile Considerations

**Touch-Friendly**: All interactive elements maintain existing mobile-friendly sizing
**Performance**: Minimize DOM manipulation scope to prevent jank on slower devices
**Battery**: Single persistent connection is more efficient than polling

## Integration with Existing Code

### Minimal Modification Approach

**File Changes Required**:
- **Primary**: `frontend/src/pages/watch/[id].astro` - replace polling with unified component
- **Optional**: Extract shared SSE utilities to `frontend/src/utils/realtime.js` for future reuse

**Backward Compatibility**:
- Server-side rendering preserved (progressive enhancement)
- Existing API endpoints unchanged
- CSS and styling patterns maintained
- Error handling patterns preserved

### Reuse Strategy

**Leverage Existing Code**:
- Connection management patterns from PresentationViewer
- Progress/step conversion from stepConversion.ts
- CSS animations and responsive design unchanged
- Error handling and user messaging patterns

**Code Extraction Opportunities**:
```javascript
// Shared utilities (future refactor)
class SSEConnection { /* Connection management */ }
class StateManager { /* State transition helpers */ }
```

## Testing Strategy

### Unit Testing Focus
1. **State Transition Logic**: Verify correct state changes for various scenarios
2. **Message Filtering**: Ensure proper collection/ID filtering
3. **Error Recovery**: Test fallback to polling, reconnection behavior
4. **Edge Cases**: Mid-session joins, rapid state changes

### Integration Testing
1. **Full User Flows**: Waiting → Live → End sequences
2. **Connection Resilience**: Network disconnect/reconnect scenarios
3. **Multi-viewer**: Multiple viewers watching same presentation
4. **Performance**: Memory leaks, connection cleanup

### Manual Testing Requirements
1. **Cross-device**: Mobile/desktop/tablet responsiveness
2. **Accessibility**: Screen reader announcements, keyboard navigation
3. **Network Conditions**: Slow/fast/intermittent connections
4. **Browser Compatibility**: EventSource support across target browsers

## Risk Assessment and Mitigation

### Technical Risks

**Risk**: SSE connection failures break auto-transition
**Mitigation**: Graceful degradation to polling fallback, comprehensive error handling

**Risk**: State transition race conditions
**Mitigation**: Promise-based sequential transitions, debouncing rapid changes

**Risk**: Memory leaks from persistent connections
**Mitigation**: Proper cleanup on page unload, connection monitoring

### User Experience Risks

**Risk**: Jarring visual transitions
**Mitigation**: Smooth fade animations, transition state feedback

**Risk**: Connection status confusion
**Mitigation**: Clear, consistent status indicators across all states

**Risk**: Edge case confusion (mid-session joins)
**Mitigation**: Intelligent state detection and user guidance

## Performance Implications

### Positive Impacts
- **Reduced Server Load**: Eliminates polling requests from waiting viewers
- **Faster Response**: Immediate state change detection vs. polling delays
- **Better Battery Life**: No periodic polling on mobile devices

### Potential Concerns
- **Persistent Connections**: Slightly higher memory usage per viewer
- **Message Processing**: Client-side filtering of SSE messages

### Monitoring Strategy
- **Connection Metrics**: Track SSE connection success/failure rates
- **Transition Performance**: Monitor state change response times
- **Error Rates**: Track fallback to polling frequency

## Implementation Priority

This design aligns with the acceptance criteria and leverages existing infrastructure effectively. The unified realtime component approach provides the best balance of user experience, technical feasibility, and maintainability while building on proven patterns from T-007-04.

The design preserves progressive enhancement and maintains all existing functionality while adding the seamless realtime transitions required by the ticket. Implementation can proceed confidently with this architecture as it addresses all identified constraints and requirements from the research phase.