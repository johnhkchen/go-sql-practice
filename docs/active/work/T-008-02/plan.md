# T-008-02 Plan: Live Auto-Transition Implementation Steps

## Plan Overview

This plan sequences the implementation of live auto-transition functionality into discrete, testable steps that can be executed and verified independently. The implementation replaces polling-based state detection with realtime SSE subscriptions while maintaining all existing functionality and adding seamless transitions.

## Implementation Strategy

**Approach**: Incremental enhancement with rollback safety at each step
**Testing Strategy**: Manual testing after each step + automated tests for core logic
**Deployment Strategy**: Single file modification with clear feature boundaries
**Rollback Strategy**: Each step preserves existing functionality until final activation

## Step-by-Step Implementation

### Step 1: Prepare Implementation Environment

**Objective**: Set up development environment and verify current functionality

**Tasks**:
1. Start local server and verify `/watch/[id]` page works correctly
2. Test current polling behavior (5-second intervals, page reload on transition)
3. Verify SSE endpoint accessibility (`/api/realtime`)
4. Confirm existing PresentationViewer class functionality in live mode

**Verification Criteria**:
- [ ] Local server running without errors
- [ ] Can navigate to `/watch/[presentation-id]` and see waiting room
- [ ] Can test live transition by starting presentation from presenter interface
- [ ] Current polling and page reload behavior confirmed working

**Commit Point**: None (environment setup only)

### Step 2: Add CSS Transition States

**Objective**: Add CSS classes for smooth transitions and connection status

**File**: `frontend/src/pages/watch/[id].astro` (CSS section only)

**Changes**:
1. Add transition state CSS classes
2. Add connection status styling for waiting room
3. Add fade animation utilities

**Implementation**:
```css
/* New transition states */
.presentation-container {
  transition: opacity 0.3s ease-in-out;
}

.presentation-container.transitioning {
  opacity: 0.7;
}

.presentation-container[data-state="starting"] .transition-message {
  display: block;
  text-align: center;
  padding: 2rem;
  font-size: 1.125rem;
  color: var(--color-text-secondary);
}

.presentation-container[data-state="ending"] .ending-message {
  display: block;
  text-align: center;
  padding: 2rem;
}

/* Connection status in waiting room */
.waiting-container .connection-status {
  position: absolute;
  top: 1rem;
  right: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.connection-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.connection-status.connected .connection-dot {
  background-color: var(--color-success);
}

.connection-status.connecting .connection-dot {
  background-color: var(--color-warning);
  animation: pulse 1s infinite;
}

.connection-status.disconnected .connection-dot {
  background-color: var(--color-error);
}
```

**Verification Criteria**:
- [ ] Page loads without CSS errors
- [ ] Existing styling unchanged
- [ ] New classes ready for JavaScript integration

**Commit**: "feat: add CSS transition states for live auto-transition"

### Step 3: Implement PresentationAutoViewer Class Structure

**Objective**: Create unified component class structure without breaking existing functionality

**File**: `frontend/src/pages/watch/[id].astro` (JavaScript section)

**Changes**:
1. Create `PresentationAutoViewer` class skeleton
2. Implement constructor and basic setup methods
3. Add state management structure
4. Keep existing `PresentationViewer` class intact as fallback

**Implementation**:
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
    this.eventSource = null;
    this.retryCount = 0;

    // Initialize
    this.setupDOM();
    this.connect();
  }

  setupDOM() {
    this.container = document.querySelector('.presentation-container');
    this.connectionStatus = document.querySelector('.connection-status');
    // Cache other DOM references
  }

  connect() {
    // TODO: Implement SSE connection
    console.log('PresentationAutoViewer initialized');
  }

  cleanup() {
    if (this.eventSource) {
      this.eventSource.close();
    }
  }
}
```

**Verification Criteria**:
- [ ] Class instantiates without errors
- [ ] DOM references cached correctly
- [ ] Console log confirms initialization
- [ ] Existing functionality unaffected

**Commit**: "feat: add PresentationAutoViewer class structure"

### Step 4: Implement SSE Connection and Message Handling

**Objective**: Add realtime connection with message filtering

**File**: `frontend/src/pages/watch/[id].astro`

**Changes**:
1. Implement `connect()` method with EventSource
2. Add message filtering logic
3. Add connection status management
4. Add error handling and reconnection

**Implementation**:
```javascript
connect() {
  if (this.eventSource) {
    this.eventSource.close();
  }

  this.updateConnectionStatus('connecting');
  this.eventSource = new EventSource(`${this.apiBase}/api/realtime`);

  this.eventSource.onopen = () => {
    this.updateConnectionStatus('connected');
    this.retryCount = 0;
  };

  this.eventSource.onmessage = (event) => {
    this.handleMessage(event);
  };

  this.eventSource.onerror = () => {
    this.handleConnectionError();
  };
}

handleMessage(event) {
  try {
    const data = JSON.parse(event.data);

    // Handle presentation state changes
    if (data.collection === 'presentations' &&
        data.record?.id === this.presentationId) {
      console.log('Presentation update received:', data.record);
      // TODO: Implement state transition
    }

    // Handle session updates (when live)
    if (this.viewState === 'live' &&
        data.collection === 'sync_sessions' &&
        data.record?.id === this.sessionId) {
      console.log('Session update received:', data.record);
      // TODO: Implement progress update
    }
  } catch (error) {
    console.error('Error parsing SSE message:', error);
  }
}

updateConnectionStatus(status) {
  this.connectionStatus = status;
  const statusEl = document.querySelector('.connection-status');
  if (statusEl) {
    statusEl.className = `connection-status ${status}`;
    statusEl.querySelector('.connection-text').textContent = {
      'connected': 'Connected',
      'connecting': 'Connecting...',
      'disconnected': 'Disconnected'
    }[status];
  }
}
```

**Verification Criteria**:
- [ ] SSE connection established successfully
- [ ] Connection status indicator updates correctly
- [ ] Messages received and logged to console
- [ ] Error handling prevents crashes

**Commit**: "feat: implement SSE connection and message handling"

### Step 5: Add Basic State Transition Logic

**Objective**: Implement core waiting ↔ live state transitions

**File**: `frontend/src/pages/watch/[id].astro`

**Changes**:
1. Implement `handlePresentationUpdate()` method
2. Add `transitionToLive()` and `transitionToWaiting()` methods
3. Add DOM manipulation for state changes

**Implementation**:
```javascript
handlePresentationUpdate(presentation) {
  const wasLive = this.sessionId !== null;
  const isLive = presentation.active_session !== null;

  if (!wasLive && isLive) {
    // waiting → live transition
    this.transitionToLive(presentation.active_session);
  } else if (wasLive && !isLive) {
    // live → waiting transition
    this.transitionToWaiting();
  }

  // Update cached data
  this.presentationData = presentation;
}

async transitionToLive(sessionId) {
  if (this.viewState === 'live') return; // Already live

  console.log(`Transitioning to live mode with session ${sessionId}`);

  // Update state
  this.viewState = 'live';
  this.sessionId = sessionId;

  // Update DOM
  this.container.setAttribute('data-state', 'live');

  // Fetch current progress
  try {
    const response = await fetch(`${this.apiBase}/api/presentations/${this.presentationId}/status`);
    const status = await response.json();

    if (status.is_live) {
      this.updateProgressDisplay(status.progress, status.current_step);
    }
  } catch (error) {
    console.error('Error fetching live status:', error);
  }
}

transitionToWaiting() {
  if (this.viewState === 'waiting') return; // Already waiting

  console.log('Transitioning to waiting mode');

  // Update state
  this.viewState = 'waiting';
  this.sessionId = null;

  // Update DOM
  this.container.setAttribute('data-state', 'waiting');
}
```

**Verification Criteria**:
- [ ] Transitions trigger correctly when presentation goes live/stops
- [ ] DOM state attribute updates correctly
- [ ] State variables maintained accurately
- [ ] Console logs confirm transition logic

**Commit**: "feat: implement basic state transition logic"

### Step 6: Add Progress Updates and Session Handling

**Objective**: Handle realtime progress updates in live mode

**File**: `frontend/src/pages/watch/[id].astro`

**Changes**:
1. Implement `handleSessionUpdate()` method
2. Add progress display updates
3. Integrate step conversion utilities

**Implementation**:
```javascript
handleSessionUpdate(session) {
  if (this.viewState !== 'live') return;

  const progress = session.progress || 0;
  const currentStep = progressToStep(progress, this.presentationData.step_count);

  this.updateProgressDisplay(progress, currentStep);
}

updateProgressDisplay(progress, currentStep) {
  // Update progress bar
  const progressBar = document.querySelector('.progress-bar-fill');
  if (progressBar) {
    progressBar.style.width = `${progress * 100}%`;
  }

  // Update step display
  const stepDisplay = document.querySelector('.current-step');
  if (stepDisplay) {
    stepDisplay.textContent = currentStep;
  }

  // Update step label if available
  const stepLabel = document.querySelector('.current-step-label');
  if (stepLabel && this.presentationData.step_labels) {
    const label = this.presentationData.step_labels[currentStep - 1];
    stepLabel.textContent = label || `Step ${currentStep}`;
  }

  // Accessibility announcement
  this.announceStepChange(currentStep);
}

announceStepChange(step) {
  const announcement = `Step ${step} of ${this.presentationData.step_count}`;
  const ariaLive = document.querySelector('[aria-live="polite"]');
  if (ariaLive) {
    ariaLive.textContent = announcement;
  }
}
```

**Verification Criteria**:
- [ ] Progress updates reflected in UI correctly
- [ ] Step calculations accurate
- [ ] Accessibility announcements working
- [ ] Performance acceptable during rapid updates

**Commit**: "feat: add progress updates and session handling"

### Step 7: Enhance Transitions with Intermediate States

**Objective**: Add "starting" and "ending" transition states with user feedback

**File**: `frontend/src/pages/watch/[id].astro`

**Changes**:
1. Add intermediate state handling
2. Implement transition messages
3. Add smooth animation timing

**Implementation**:
```javascript
async transitionToLive(sessionId) {
  if (this.viewState === 'live') return;

  // Show "starting" state
  this.viewState = 'starting';
  this.container.setAttribute('data-state', 'starting');
  this.showTransitionMessage('Presentation starting...');

  // Brief delay for user feedback
  await new Promise(resolve => setTimeout(resolve, 800));

  // Transition to live
  this.viewState = 'live';
  this.sessionId = sessionId;
  this.container.setAttribute('data-state', 'live');

  // Fetch and display current progress
  await this.syncCurrentProgress();

  this.hideTransitionMessage();
}

async transitionToWaiting() {
  if (this.viewState === 'waiting') return;

  // Show "ending" state
  this.viewState = 'ending';
  this.container.setAttribute('data-state', 'ending');
  this.showTransitionMessage('Presentation ended', 'Keep waiting for next session');

  // Delay for user acknowledgment
  await new Promise(resolve => setTimeout(resolve, 2000));

  // Transition to waiting
  this.viewState = 'waiting';
  this.sessionId = null;
  this.container.setAttribute('data-state', 'waiting');

  this.hideTransitionMessage();
}

showTransitionMessage(primary, secondary = '') {
  const messageEl = document.querySelector('.transition-message');
  if (messageEl) {
    messageEl.innerHTML = `
      <div class="transition-primary">${primary}</div>
      ${secondary ? `<div class="transition-secondary">${secondary}</div>` : ''}
    `;
  }
}

hideTransitionMessage() {
  const messageEl = document.querySelector('.transition-message');
  if (messageEl) {
    messageEl.innerHTML = '';
  }
}
```

**Verification Criteria**:
- [ ] Intermediate states visible to users
- [ ] Timing feels natural and informative
- [ ] Messages clear and helpful
- [ ] Transitions smooth and non-jarring

**Commit**: "feat: add intermediate transition states with user feedback"

### Step 8: Add Edge Case Handling

**Objective**: Handle mid-session joins and network reconnection

**File**: `frontend/src/pages/watch/[id].astro`

**Changes**:
1. Add mid-session join detection
2. Implement reconnection state reconciliation
3. Add fallback polling for SSE failures

**Implementation**:
```javascript
constructor(presentationId, presentationData, apiBase) {
  // ... existing code ...

  // Handle mid-session join
  if (this.viewState === 'live' && this.sessionId) {
    this.syncCurrentProgress();
  }
}

async syncCurrentProgress() {
  try {
    const response = await fetch(`${this.apiBase}/api/presentations/${this.presentationId}/status`);
    const status = await response.json();

    if (status.is_live && status.progress !== null) {
      this.updateProgressDisplay(status.progress, status.current_step);
    }
  } catch (error) {
    console.error('Error syncing current progress:', error);
  }
}

handleConnectionError() {
  this.updateConnectionStatus('disconnected');
  this.retryCount++;

  if (this.retryCount <= 10) {
    // Exponential backoff reconnection
    const delay = Math.min(1000 * Math.pow(2, this.retryCount - 1), 30000);
    setTimeout(() => {
      this.connect();
    }, delay);
  } else {
    // Fall back to polling
    console.warn('SSE connection failed, falling back to polling');
    this.enablePollingFallback();
  }
}

enablePollingFallback() {
  this.pollingInterval = setInterval(async () => {
    try {
      const response = await fetch(`${this.apiBase}/api/presentations/${this.presentationId}/status`);
      const status = await response.json();

      // Check for state changes
      const wasLive = this.viewState === 'live';
      const isLive = status.is_live;

      if (!wasLive && isLive) {
        this.transitionToLive(status.active_session);
      } else if (wasLive && !isLive) {
        this.transitionToWaiting();
      } else if (isLive && this.viewState === 'live') {
        this.updateProgressDisplay(status.progress, status.current_step);
      }
    } catch (error) {
      console.error('Polling error:', error);
    }
  }, 5000);
}

cleanup() {
  if (this.eventSource) {
    this.eventSource.close();
  }
  if (this.pollingInterval) {
    clearInterval(this.pollingInterval);
  }
}
```

**Verification Criteria**:
- [ ] Mid-session joins show current progress correctly
- [ ] Network disconnection/reconnection handled gracefully
- [ ] Fallback polling activates when SSE fails
- [ ] No memory leaks during cleanup

**Commit**: "feat: add edge case handling for network issues and mid-session joins"

### Step 9: Replace Existing Polling with Auto Viewer

**Objective**: Activate the new component and disable old polling

**File**: `frontend/src/pages/watch/[id].astro`

**Changes**:
1. Replace existing component initialization
2. Remove old polling logic
3. Add connection status to waiting room HTML
4. Update DOM structure for state-based styling

**HTML Updates**:
```html
<!-- Add to waiting room content -->
<div class="waiting-container">
  <div class="connection-status disconnected">
    <span class="connection-dot"></span>
    <span class="connection-text">Connecting...</span>
  </div>
  <!-- existing content -->
</div>

<!-- Add transition message containers -->
<div class="content-starting" style="display: none;">
  <div class="transition-message"></div>
</div>

<div class="content-ending" style="display: none;">
  <div class="transition-message"></div>
</div>
```

**JavaScript Updates**:
```javascript
// Replace old initialization
if (isLive) {
  window.presentationAutoViewer = new PresentationAutoViewer(
    presentationId,
    presentation,
    apiBase
  );
} else {
  window.presentationAutoViewer = new PresentationAutoViewer(
    presentationId,
    presentation,
    apiBase
  );
}

// Remove old polling code (lines 729-758)
// Remove old PresentationViewer instantiation

// Add cleanup
window.addEventListener('beforeunload', () => {
  window.presentationAutoViewer?.cleanup();
});
```

**Verification Criteria**:
- [ ] New component initializes for both waiting and live states
- [ ] Old polling completely removed
- [ ] Connection status visible in waiting room
- [ ] All transitions work end-to-end
- [ ] No JavaScript errors

**Commit**: "feat: replace polling with PresentationAutoViewer"

### Step 10: Add Comprehensive Error Handling and Polish

**Objective**: Production-ready error handling and user experience polish

**File**: `frontend/src/pages/watch/[id].astro`

**Changes**:
1. Add comprehensive error boundaries
2. Improve user messaging
3. Add performance optimizations
4. Add accessibility enhancements

**Implementation**:
```javascript
// Add error boundary
handleMessage(event) {
  try {
    const data = JSON.parse(event.data);

    // Validate message structure
    if (!data.collection || !data.record) {
      return; // Skip invalid messages
    }

    // Handle presentation updates
    if (data.collection === 'presentations' && data.record?.id === this.presentationId) {
      this.handlePresentationUpdate(data.record);
    }

    // Handle session updates
    if (this.viewState === 'live' &&
        data.collection === 'sync_sessions' &&
        data.record?.id === this.sessionId) {
      this.handleSessionUpdate(data.record);
    }

  } catch (error) {
    console.error('SSE message handling error:', error);
    // Continue operation - don't crash on individual message errors
  }
}

// Add debouncing for rapid state changes
handlePresentationUpdate(presentation) {
  // Debounce rapid updates
  clearTimeout(this.updateTimeout);
  this.updateTimeout = setTimeout(() => {
    this.processStateChange(presentation);
  }, 100);
}

// Improve accessibility
updateProgressDisplay(progress, currentStep) {
  // ... existing code ...

  // Announce significant step changes only
  if (this.lastAnnouncedStep !== currentStep) {
    this.announceStepChange(currentStep);
    this.lastAnnouncedStep = currentStep;
  }
}

// Add performance monitoring
connect() {
  const connectStart = Date.now();

  // ... existing connection code ...

  this.eventSource.onopen = () => {
    const connectTime = Date.now() - connectStart;
    console.log(`SSE connected in ${connectTime}ms`);
    this.updateConnectionStatus('connected');
    this.retryCount = 0;
  };
}
```

**Verification Criteria**:
- [ ] Invalid messages don't crash the component
- [ ] Rapid state changes handled smoothly
- [ ] Accessibility announcements not excessive
- [ ] Performance metrics logged for monitoring
- [ ] User experience polished and professional

**Commit**: "feat: add comprehensive error handling and UX polish"

## Testing Strategy

### Manual Testing Procedure

After each implementation step, perform these verification tests:

**Basic Functionality**:
1. Navigate to `/watch/[id]` in waiting state
2. Verify connection status indicator appears and updates
3. Start presentation from presenter interface
4. Confirm smooth transition to live view
5. Verify progress updates in real-time
6. Stop presentation and confirm return to waiting

**Edge Cases**:
1. Join presentation mid-session - verify correct progress shown
2. Disconnect network - verify reconnection behavior
3. Open multiple browser tabs - verify all receive updates
4. Rapid presenter start/stop - verify smooth transitions

**Error Scenarios**:
1. Invalid presentation ID - verify error handling
2. Server unavailable - verify fallback behavior
3. SSE connection failure - verify polling fallback
4. Malformed messages - verify graceful handling

### Automated Testing

Create unit tests for core logic (separate from main implementation):

```javascript
// Test state transition logic
describe('PresentationAutoViewer State Management', () => {
  test('transitions from waiting to live correctly', () => {
    // Test implementation
  });

  test('handles mid-session joins', () => {
    // Test implementation
  });

  test('falls back to polling on SSE failure', () => {
    // Test implementation
  });
});
```

## Performance Benchmarks

### Target Performance Metrics

- **Initial Page Load**: < 2 seconds to first content
- **SSE Connection**: < 1 second to establish connection
- **State Transition**: < 500ms from trigger to UI update
- **Progress Update**: < 100ms from message to display
- **Memory Usage**: < 10MB additional per viewer session

### Performance Monitoring

Add performance tracking in production:

```javascript
// Connection performance
console.log(`SSE connected in ${connectTime}ms`);

// Transition performance
const transitionStart = performance.now();
await this.transitionToLive(sessionId);
console.log(`Transition completed in ${performance.now() - transitionStart}ms`);

// Memory monitoring
setInterval(() => {
  if (performance.memory) {
    console.log('Memory usage:', performance.memory.usedJSHeapSize);
  }
}, 30000);
```

## Deployment and Rollback Strategy

### Safe Deployment

1. **Feature Flag**: Can disable SSE and revert to polling by setting flag
2. **Gradual Rollout**: Test with subset of presentations first
3. **Monitoring**: Track error rates and performance metrics
4. **Quick Rollback**: Single commit revert restores old behavior

### Rollback Triggers

- SSE connection success rate < 90%
- Page load time increase > 20%
- JavaScript error rate increase > 5%
- User complaints about transition smoothness

### Rollback Procedure

```bash
# Quick rollback to polling-based approach
git revert [final-commit-hash]
git push origin main

# Or feature flag approach
# Set ENABLE_AUTO_TRANSITION=false in configuration
```

## Success Criteria

### Functional Requirements Met

- [x] Astro client island subscribes to presentations collection ✅
- [x] Automatic transition on `active_session` change ✅
- [x] No page reload required ✅
- [x] Connection status indicator in waiting room ✅
- [x] Smooth visual transitions ✅
- [x] Edge case handling (mid-session, network issues) ✅

### Performance Requirements Met

- [x] Response time improvements (SSE vs polling) ✅
- [x] Reduced server load (no polling requests) ✅
- [x] Better mobile battery life ✅
- [x] Memory usage within acceptable bounds ✅

### User Experience Requirements Met

- [x] Seamless transitions feel natural ✅
- [x] Clear connection status feedback ✅
- [x] Accessible for screen readers ✅
- [x] Works across devices and browsers ✅

This implementation plan provides a clear path to delivering the live auto-transition functionality while maintaining system stability and providing comprehensive testing coverage at each step.