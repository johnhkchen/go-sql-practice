# Design: T-006-04 - sync-viewer-page

## Problem Statement

Create a viewer page at `/sync/[id]` that displays real-time progress updates for sync sessions. Viewers need to see live progress changes as admins adjust the slider on the control page, with smooth visual updates and robust error handling.

## Design Options Evaluation

### Option 1: Manual EventSource Implementation (Chosen)

**Approach**: Direct SSE connection to PocketBase realtime API using native EventSource

**Pros**:
- No additional dependencies - uses native browser APIs
- Matches existing project's minimal dependency approach
- Direct integration with PocketBase's built-in SSE realtime API
- Full control over connection management and error handling
- Consistent with project's vanilla JavaScript pattern (from control page)
- Lightweight and performant

**Cons**:
- More manual implementation work
- Need to handle SSE protocol specifics
- Manual reconnection logic required
- JSON parsing and filtering needed

**Technical Details**:
```javascript
const eventSource = new EventSource('/api/realtime');
eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  if (data.collection === 'sync_sessions' && data.record.id === sessionId) {
    updateProgress(data.record.progress);
  }
};
```

**Rationale**: Best fits existing architecture patterns and project constraints.

### Option 2: PocketBase JavaScript SDK

**Approach**: Add PocketBase JS SDK dependency and use built-in realtime methods

**Pros**:
- Official SDK with built-in realtime subscriptions
- Automatic reconnection handling
- Type safety and better developer experience
- Built-in authentication support

**Cons**:
- Adds 50KB+ dependency to minimal project
- Would be only external JavaScript library
- Overkill for single-feature usage
- May conflict with project's static generation approach
- Authentication features not needed for public viewing

**Technical Details**:
```javascript
import PocketBase from 'pocketbase';
const pb = new PocketBase('http://localhost:8090');
pb.collection('sync_sessions').subscribe(sessionId, callback);
```

**Rejection Reason**: Violates project's minimal dependency philosophy and adds unnecessary complexity.

### Option 3: WebSocket Implementation

**Approach**: Custom WebSocket connection to PocketBase or separate WebSocket server

**Pros**:
- Full bidirectional communication
- Lower latency than SSE
- More efficient for high-frequency updates

**Cons**:
- PocketBase uses SSE, not WebSockets - would need custom backend
- Increased complexity for simple read-only updates
- Overkill for viewer-only functionality
- Would require significant backend changes

**Rejection Reason**: PocketBase architecture is SSE-based; WebSocket would require backend overhaul.

### Option 4: Polling-Based Updates

**Approach**: Regular fetch requests to session API endpoint

**Pros**:
- Simple implementation using existing patterns
- Works in all network conditions
- Easy error handling and retry logic

**Cons**:
- Not real-time - updates delayed by polling interval
- Higher server load with multiple viewers
- Poor user experience compared to push-based updates
- Inconsistent with "realtime" requirement

**Rejection Reason**: Fails to meet real-time update requirement from acceptance criteria.

## Chosen Design: Manual EventSource Implementation

### Architecture Overview

```
┌─────────────────┐    SSE     ┌─────────────────┐    Record    ┌─────────────────┐
│  Viewer Page    │◄──────────►│ PocketBase      │◄────────────►│ Admin Control   │
│  /sync/[id]     │            │ Realtime API    │              │ Page            │
└─────────────────┘            └─────────────────┘              └─────────────────┘
```

### Data Flow

1. **Page Load**: Fetch initial session data via REST API
2. **SSE Connection**: Establish EventSource to `/api/realtime`
3. **Filter Messages**: Only process sync_sessions updates for current session ID
4. **Update UI**: Animate progress bar when progress changes
5. **Connection Management**: Handle disconnections, errors, and reconnection

### User Interface Design

#### Page Layout
```
┌─────────────────────────────────────────────────────────────────┐
│                        Navigation                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Session Viewer - [Session ID]                                 │
│                                                                 │
│  ┌─ Session Status ────────────────────────────────────────┐   │
│  │ ● Connected     Session: abc123...                      │   │
│  │ Progress: 42.0% (Updated 2 seconds ago)                 │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─ Progress Display ──────────────────────────────────────┐   │
│  │                                                         │   │
│  │  ████████████████████████████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓   │   │
│  │  42.0%                                                  │   │
│  │                                                         │   │
│  │  [Large, smooth progress bar with CSS transitions]     │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─ Session Info ─────────────────────────────────────────┐   │
│  │ Created: Jan 15, 2024 2:30 PM                          │   │
│  │ Last Updated: Jan 15, 2024 2:45 PM                     │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Visual Design Elements

**Connection Status Indicator**:
- Green dot (●): Connected and receiving updates
- Yellow dot (●): Connecting/Reconnecting
- Red dot (●): Disconnected
- Consistent with control page error handling patterns

**Progress Bar**:
- HTML `<progress>` element for accessibility
- CSS transitions for smooth updates (0.3s ease)
- Large size for mobile visibility
- Percentage display with 1 decimal precision

**Responsive Behavior**:
- Mobile-first design matching existing pages
- Touch-friendly elements
- Readable text on small screens

### Technical Implementation Strategy

#### SSE Connection Management

```javascript
class SyncViewer {
  constructor(sessionId) {
    this.sessionId = sessionId;
    this.eventSource = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000; // Start with 1s, exponential backoff
    this.isConnected = false;
  }

  connect() {
    this.eventSource = new EventSource('/api/realtime');
    this.setupEventHandlers();
  }

  setupEventHandlers() {
    this.eventSource.onopen = () => this.handleConnect();
    this.eventSource.onmessage = (event) => this.handleMessage(event);
    this.eventSource.onerror = () => this.handleError();
  }
}
```

#### Message Processing Pattern

```javascript
handleMessage(event) {
  try {
    const data = JSON.parse(event.data);

    // Filter for sync_sessions updates to our session
    if (data.collection === 'sync_sessions' &&
        data.action === 'update' &&
        data.record &&
        data.record.id === this.sessionId) {

      this.updateProgress(data.record.progress);
      this.updateTimestamp(data.record.updated);
    }
  } catch (error) {
    console.warn('Invalid SSE message:', error);
  }
}
```

#### Error Handling Strategy

**Connection Errors**:
- Exponential backoff reconnection (1s, 2s, 4s, 8s, 16s max)
- Visual connection status updates
- User notification after 3 failed attempts
- Manual retry button after 5 failures

**Message Processing Errors**:
- Log malformed messages but continue processing
- Validate data structure before UI updates
- Graceful handling of unexpected message formats

**Network Changes**:
- Detect online/offline events
- Automatic reconnection when network returns
- Pause reconnection attempts while offline

### Mobile Optimization

#### Performance Considerations
- Minimize DOM updates during progress changes
- Use CSS transforms instead of layout changes
- Debounce rapid progress updates (max 30fps)
- Efficient memory usage for long viewing sessions

#### Battery Life
- Reduce update frequency when page not visible
- Use Page Visibility API to pause non-essential updates
- Efficient event processing to minimize CPU usage

### Accessibility Design

#### Screen Reader Support
- ARIA live regions for progress announcements
- Descriptive labels for all interactive elements
- Progress updates announced at reasonable intervals (not every small change)

#### Keyboard Navigation
- Tab navigation through all interactive elements
- Focus indicators matching project standards
- Keyboard shortcuts for common actions

#### Visual Accessibility
- High contrast color scheme support
- Large touch targets (minimum 44px)
- Readable fonts and appropriate text sizing
- Reduced motion support for animations

### Error State Design

#### Session Not Found
```
┌─────────────────────────────────────────────────────────────────┐
│                        Navigation                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ⚠ Session Not Found                                            │
│                                                                 │
│  The sync session you're looking for doesn't exist or          │
│  has been removed.                                              │
│                                                                 │
│  [← Back to Home]                                               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Connection Problems
```
┌─ Connection Status ─────────────────────────────────────────┐
│ ⚠ Connection Lost                                          │
│ Trying to reconnect... (Attempt 3/5)                      │
│                                                            │
│ [Retry Now] [View Offline]                                 │
└────────────────────────────────────────────────────────────┘
```

### CSS Architecture

#### Design System Integration
- Use existing CSS custom properties from BaseLayout
- Follow established color scheme and spacing patterns
- Consistent typography with other pages
- Responsive breakpoints matching site standards

#### Component-Specific Styles
```css
.sync-viewer {
  --progress-height: 24px;
  --progress-color: var(--color-primary);
  --status-indicator-size: 12px;
}

.progress-container {
  position: relative;
  margin: var(--space-lg) 0;
}

.progress-bar {
  width: 100%;
  height: var(--progress-height);
  transition: width 0.3s ease;
}

@media (prefers-reduced-motion: reduce) {
  .progress-bar {
    transition: none;
  }
}
```

## Integration Points

### Admin Control Page Integration
- Share viewer URL generation pattern
- Consistent error handling approach
- Similar responsive design principles
- Unified CSS design system

### PocketBase Integration
- Use existing session collection structure
- Leverage built-in realtime API
- Follow established API patterns
- Consistent error response handling

### Astro Framework Integration
- Server-side rendering for initial page load
- Client-side hydration for realtime functionality
- Static asset optimization
- SEO-friendly metadata

## Rejected Design Alternatives

### Complex State Management
**Why rejected**: Overkill for simple progress display; project favors vanilla JavaScript approaches

### Real-time Bidirectional Communication
**Why rejected**: Viewer page is read-only; two-way communication not needed

### Advanced Animation Library
**Why rejected**: CSS transitions sufficient; avoid external dependencies

### Iframe Embedding
**Why rejected**: Would complicate styling and interaction; not necessary for this use case

## Security Considerations

### Public Access Model
- No authentication required - intentionally public
- Session ID acts as access token (share-based security)
- Read-only access prevents malicious modifications
- No sensitive data exposure beyond progress value

### Input Validation
- Validate session ID format before API calls
- Sanitize all data received from SSE messages
- Prevent XSS through proper data handling

## Success Metrics

### Functional Requirements
- ✅ Real-time progress updates (< 1s latency)
- ✅ Mobile browser compatibility
- ✅ Connection status visibility
- ✅ Graceful error handling
- ✅ Accessibility compliance

### Performance Requirements
- Page load time < 2s on 3G
- Memory usage stable over 30+ minute sessions
- Battery impact minimal on mobile devices
- Smooth progress animations (no jank)

### User Experience Requirements
- Intuitive interface requiring no explanation
- Reliable connection with clear status feedback
- Professional appearance matching site design
- Accessible to users with disabilities

This design provides a robust, user-friendly sync viewer page that integrates seamlessly with the existing architecture while meeting all acceptance criteria through proven technical approaches.