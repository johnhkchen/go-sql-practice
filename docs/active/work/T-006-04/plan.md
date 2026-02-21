# Plan: T-006-04 - sync-viewer-page

## Implementation Overview

This plan breaks down the creation of the sync viewer page into discrete, testable steps that can be executed incrementally. Each step builds upon the previous ones and can be validated independently.

## Implementation Steps

### Step 1: Create Basic Page Structure
**Goal**: Establish the basic Astro page file with server-side session fetching
**Files**: `frontend/src/pages/sync/[id].astro` (new)
**Estimated Time**: 30 minutes

#### Tasks:
1. Create the new Astro page file at correct path
2. Add frontmatter with SSR disabled (`export const prerender = false`)
3. Implement session ID extraction and validation from URL params
4. Add basic session data fetching via REST API
5. Include error handling for session not found/network errors
6. Create basic HTML template structure with BaseLayout integration

#### Verification:
- Page loads at `/sync/[session-id]` URL
- Shows "Session not found" for invalid session IDs
- Displays basic session information for valid sessions
- Error states display properly for network failures
- Page integrates properly with site navigation

#### Code Structure:
```astro
---
export const prerender = false;
import BaseLayout from '../../layouts/BaseLayout.astro';

const { id } = Astro.params;
// Session fetching logic
// Error handling
---

<BaseLayout title={pageTitle}>
  {error && (
    <div class="error-container">
      <!-- Error states -->
    </div>
  )}

  {session && (
    <div class="sync-viewer-container">
      <!-- Basic session display -->
    </div>
  )}
</BaseLayout>
```

### Step 2: Implement Static Progress Display
**Goal**: Add progress bar element with initial static display
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 20 minutes

#### Tasks:
1. Add HTML progress element with proper accessibility attributes
2. Display initial progress value from session data
3. Add progress percentage text display
4. Create basic CSS styling for progress bar
5. Ensure mobile-friendly sizing and layout

#### Verification:
- Progress bar displays with correct initial value
- Percentage text matches progress value
- Progress bar is visually appealing and accessible
- Layout works on mobile screens
- Screen readers announce progress properly

#### Code Structure:
```html
<div class="progress-section">
  <h2>Session Progress</h2>
  <div class="progress-container">
    <progress class="progress-bar" value={session.progress} max="1"></progress>
    <div class="progress-label">{Math.round(session.progress * 100)}%</div>
  </div>
</div>
```

### Step 3: Add Connection Status Infrastructure
**Goal**: Create connection status display and basic JavaScript framework
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 45 minutes

#### Tasks:
1. Add connection status header with indicator and text
2. Create JavaScript SyncViewer class structure
3. Implement basic state management for connection status
4. Add methods for updating connection status visually
5. Create CSS for connection status indicators (green/yellow/red dots)

#### Verification:
- Connection status header displays properly
- Status indicators show different colors for different states
- JavaScript class instantiates without errors
- Status can be changed programmatically
- CSS animations work smoothly

#### Code Structure:
```html
<div class="connection-status">
  <span class="status-indicator" id="status-indicator"></span>
  <span class="status-text" id="status-text">Connecting...</span>
  <span class="last-update" id="last-update"></span>
</div>

<script client:load>
class SyncViewer {
  constructor(sessionId, initialProgress, apiBase) {
    // Basic setup
  }

  updateConnectionStatus(status) {
    // Update status display
  }
}
</script>
```

### Step 4: Implement SSE Connection
**Goal**: Establish EventSource connection to PocketBase realtime API
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 60 minutes

#### Tasks:
1. Add EventSource connection setup in SyncViewer class
2. Implement connection event handlers (onopen, onmessage, onerror)
3. Add basic message filtering for sync_sessions collection
4. Update connection status based on SSE events
5. Add console logging for debugging

#### Verification:
- EventSource connection establishes successfully
- Connection status updates to "Connected" when established
- Console logs show incoming SSE messages
- Messages are properly filtered by collection type
- Connection status shows "Connecting" during initial connection

#### Code Structure:
```javascript
connect() {
  this.eventSource = new EventSource('/api/realtime');
  this.eventSource.onopen = () => this.handleConnect();
  this.eventSource.onmessage = (event) => this.handleMessage(event);
  this.eventSource.onerror = () => this.handleError();
}

handleMessage(event) {
  const data = JSON.parse(event.data);
  if (data.collection === 'sync_sessions' && data.record.id === this.sessionId) {
    // Process message
  }
}
```

### Step 5: Add Real-time Progress Updates
**Goal**: Process SSE messages to update progress bar in real-time
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 45 minutes

#### Tasks:
1. Implement progress update method in SyncViewer class
2. Extract progress value from SSE messages
3. Update progress bar value and percentage text
4. Add smooth CSS transitions for progress changes
5. Update timestamp display when progress changes

#### Verification:
- Progress bar updates when admin changes progress on control page
- Percentage text updates to match progress value
- Transitions are smooth and visually appealing
- Timestamp updates with each change
- Multiple rapid updates don't cause UI issues

#### Code Structure:
```javascript
handleMessage(event) {
  const data = JSON.parse(event.data);
  if (data.collection === 'sync_sessions' &&
      data.action === 'update' &&
      data.record.id === this.sessionId) {
    this.updateProgress(data.record.progress);
    this.updateTimestamp(data.record.updated);
  }
}

updateProgress(progress) {
  const progressBar = document.getElementById('progress-bar');
  const progressLabel = document.getElementById('progress-label');

  progressBar.value = progress;
  progressLabel.textContent = Math.round(progress * 100) + '%';
}
```

### Step 6: Implement Error Handling and Reconnection
**Goal**: Add robust error handling with automatic reconnection
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 90 minutes

#### Tasks:
1. Implement exponential backoff reconnection logic
2. Add connection timeout handling
3. Create user-friendly error messages
4. Add manual retry functionality
5. Handle browser online/offline events
6. Implement connection cleanup on page unload

#### Verification:
- Automatic reconnection works when connection is lost
- Error messages are user-friendly and informative
- Manual retry button appears after several failed attempts
- Page handles network online/offline transitions
- No memory leaks when leaving the page
- Connection attempts stop after reasonable number of retries

#### Code Structure:
```javascript
handleError() {
  this.isConnected = false;
  this.updateConnectionStatus('disconnected');
  this.scheduleReconnect();
}

scheduleReconnect() {
  if (this.reconnectAttempts < this.maxReconnectAttempts) {
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 16000);
    setTimeout(() => this.attemptReconnect(), delay);
    this.reconnectAttempts++;
  } else {
    this.showManualRetry();
  }
}
```

### Step 7: Add Session Information Display
**Goal**: Display comprehensive session metadata and information
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 30 minutes

#### Tasks:
1. Add session information section to HTML template
2. Display session ID, creation time, and last update time
3. Format dates and times appropriately
4. Style session information consistently with site design
5. Ensure information is accessible and mobile-friendly

#### Verification:
- Session information displays correctly
- Dates and times are properly formatted
- Information updates when session data changes
- Layout is responsive on mobile devices
- Screen readers can access all information

#### Code Structure:
```html
<div class="session-info">
  <h3>Session Details</h3>
  <div class="info-grid">
    <div class="info-item">
      <label>Session ID:</label>
      <code>{session.id}</code>
    </div>
    <div class="info-item">
      <label>Created:</label>
      <span>{new Date(session.created).toLocaleString()}</span>
    </div>
    <!-- Additional metadata -->
  </div>
</div>
```

### Step 8: Enhance CSS Styling and Responsive Design
**Goal**: Polish visual design and ensure mobile compatibility
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 45 minutes

#### Tasks:
1. Refine progress bar styling with consistent design system
2. Add responsive breakpoints for mobile layout
3. Implement dark mode support using CSS media queries
4. Add CSS transitions and animations
5. Ensure high contrast mode compatibility
6. Optimize for touch interfaces

#### Verification:
- Design matches existing site aesthetic
- Layout adapts properly to different screen sizes
- Dark mode works correctly if system preference is set
- Animations are smooth and don't cause performance issues
- Touch targets are appropriately sized for mobile
- High contrast mode displays clearly

#### Code Structure:
```css
.sync-viewer-container {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--space-lg);
}

.progress-bar {
  width: 100%;
  height: 24px;
  transition: width 0.3s ease;
}

@media (max-width: 767px) {
  .sync-viewer-container {
    padding: var(--space-md);
  }
}

@media (prefers-reduced-motion: reduce) {
  .progress-bar {
    transition: none;
  }
}
```

### Step 9: Add Accessibility Enhancements
**Goal**: Ensure comprehensive accessibility compliance
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 30 minutes

#### Tasks:
1. Add ARIA labels and descriptions to all interactive elements
2. Implement live regions for progress announcements
3. Ensure keyboard navigation works properly
4. Add screen reader announcements for important state changes
5. Verify color contrast meets WCAG guidelines

#### Verification:
- Screen readers announce progress changes appropriately
- All elements have proper ARIA labels
- Keyboard navigation works for all interactive elements
- Color contrast meets accessibility standards
- Focus indicators are visible and appropriate

#### Code Structure:
```html
<progress
  class="progress-bar"
  value={session.progress}
  max="1"
  aria-labelledby="progress-heading"
  aria-describedby="progress-description">
</progress>

<div id="progress-description" class="sr-only">
  Progress updates automatically as the session advances
</div>

<div aria-live="polite" id="progress-announcements"></div>
```

### Step 10: Final Testing and Performance Optimization
**Goal**: Comprehensive testing and performance validation
**Files**: `frontend/src/pages/sync/[id].astro` (modify)
**Estimated Time**: 60 minutes

#### Tasks:
1. Test with multiple concurrent viewers
2. Verify memory usage remains stable over long sessions
3. Test network interruption scenarios
4. Validate performance on low-powered devices
5. Cross-browser compatibility testing
6. Mobile device testing

#### Verification:
- Multiple viewers can watch same session simultaneously
- Memory usage doesn't increase over time
- Graceful handling of network interruptions
- Smooth performance on older mobile devices
- Consistent behavior across modern browsers
- Touch interactions work properly on mobile

## Testing Strategy

### Manual Testing Checklist

#### Basic Functionality:
- [ ] Page loads with valid session ID
- [ ] Displays "Session not found" for invalid session ID
- [ ] Shows progress bar with current session progress
- [ ] Real-time updates when admin changes progress
- [ ] Connection status indicator works correctly
- [ ] Session information displays properly

#### Error Scenarios:
- [ ] Network interruption handling
- [ ] Server unavailable scenarios
- [ ] Invalid session data handling
- [ ] SSE connection failures
- [ ] Rapid connection loss/recovery cycles

#### Accessibility:
- [ ] Screen reader compatibility
- [ ] Keyboard navigation
- [ ] High contrast mode
- [ ] Reduced motion preference
- [ ] ARIA labels and live regions

#### Performance:
- [ ] Initial page load speed
- [ ] Memory usage over 30+ minutes
- [ ] Mobile device performance
- [ ] Multiple tab scenarios
- [ ] Background tab behavior

### Integration Testing

#### With Admin Control Page:
1. Open admin control page in one browser tab
2. Open viewer page in another tab (same session)
3. Move progress slider on admin page
4. Verify viewer page updates in real-time
5. Test connection recovery after network interruption

#### With PocketBase Backend:
1. Verify SSE messages are properly formatted
2. Test session CRUD operations don't break viewer
3. Validate error responses are handled correctly
4. Check realtime API performance under load

### Browser Compatibility Testing

**Primary Browsers** (must work perfectly):
- Chrome/Chromium latest
- Firefox latest
- Safari latest
- Mobile Chrome (Android)
- Mobile Safari (iOS)

**Secondary Browsers** (should work with graceful degradation):
- Edge latest
- Firefox ESR
- Older mobile browsers

### Performance Benchmarks

**Page Load Performance**:
- Initial load time: < 2 seconds on 3G
- Time to first contentful paint: < 1 second
- JavaScript bundle size: < 50KB (inline)

**Runtime Performance**:
- Memory usage: Stable over 60+ minutes
- CPU usage: Minimal when page not visible
- Battery impact: Low on mobile devices

**Network Performance**:
- SSE reconnection time: < 5 seconds
- Failed connection recovery: < 30 seconds
- Offline detection: < 10 seconds

## Risk Mitigation

### High-Risk Areas:

**SSE Connection Stability**:
- *Risk*: Connection drops frequently on mobile networks
- *Mitigation*: Robust reconnection logic with exponential backoff
- *Testing*: Extensive mobile network testing

**Performance on Low-End Devices**:
- *Risk*: Animations cause frame drops on older phones
- *Mitigation*: Reduced motion preference support, efficient CSS
- *Testing*: Test on low-end Android devices

**Cross-Browser SSE Support**:
- *Risk*: Inconsistent EventSource behavior across browsers
- *Mitigation*: Thorough browser testing, fallback strategies
- *Testing*: Comprehensive cross-browser test suite

### Rollback Strategy:

If critical issues are discovered after deployment:
1. Temporarily redirect `/sync/[id]` to a static "Coming Soon" page
2. Fix issues in development environment
3. Re-deploy with fixes
4. Enable viewer page with proper testing

### Monitoring and Debugging:

**Console Logging**:
- Connection state changes
- Message processing events
- Error conditions
- Performance metrics

**User-Facing Indicators**:
- Clear connection status display
- Helpful error messages
- Manual retry options

## Success Criteria

### Functional Requirements:
- ✅ Page accessible at `/sync/[id]` URL
- ✅ Real-time progress updates with < 1 second latency
- ✅ Connection status visibility
- ✅ Graceful error handling and recovery
- ✅ Mobile browser compatibility
- ✅ Accessibility compliance

### Performance Requirements:
- ✅ Page load time < 2 seconds on 3G
- ✅ Stable memory usage over 60+ minutes
- ✅ Smooth progress animations (60fps target)
- ✅ Low CPU usage when page not active

### User Experience Requirements:
- ✅ Intuitive interface requiring no explanation
- ✅ Professional appearance matching site design
- ✅ Reliable connection with clear status feedback
- ✅ Accessible to users with disabilities

This plan provides a systematic approach to implementing the sync viewer page with incremental validation at each step, ensuring a robust and user-friendly final product.