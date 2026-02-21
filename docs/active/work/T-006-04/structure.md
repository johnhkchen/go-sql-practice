# Structure: T-006-04 - sync-viewer-page

## File Changes Overview

This implementation requires **1 new file** and **0 modifications** to existing files. The design leverages existing infrastructure without changes to backend code or shared components.

## Files to Create

### 1. `frontend/src/pages/sync/[id].astro` (NEW)
**Purpose**: Main viewer page for real-time sync session progress display
**Size**: ~800-1000 lines (following control page patterns)
**Type**: Astro page with SSR disabled for dynamic routing

#### File Structure Breakdown:

```astro
---
// Disable prerendering for SSR (lines 1-10)
export const prerender = false;
import BaseLayout from '../../layouts/BaseLayout.astro';

// URL parameter extraction (lines 11-20)
const { id } = Astro.params;
if (!id) return Astro.redirect('/404');

// API configuration (lines 21-30)
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
const FETCH_TIMEOUT = 5000;

// Initial session data fetch (lines 31-80)
let session = null;
let error = null;
// [Error handling logic similar to control page pattern]

// Page metadata (lines 81-90)
const pageTitle = session ? `Sync Viewer - Session ${id}` : 'Error';
---

<!-- HTML Template (lines 91-200) -->
<BaseLayout title={pageTitle}>
  <!-- Error states -->
  <!-- Session viewer UI -->
  <!-- Progress display -->
  <!-- Session info -->
</BaseLayout>

<!-- Client-side JavaScript (lines 201-700) -->
{session && (
  <script client:load define:vars={{ sessionId: id, initialProgress: session.progress, apiBase: API_BASE }}>
    // SyncViewer class implementation
    // SSE connection management
    // Progress updates
    // Error handling
  </script>
)}

<!-- CSS Styles (lines 701-1000) -->
<style>
  /* Component styles */
  /* Responsive design */
  /* Accessibility enhancements */
</style>
```

#### Module Boundaries:

**Server-Side Section (Astro frontmatter)**:
- Session ID validation and routing
- Initial session data fetching via REST API
- Error state determination
- Page metadata generation

**Template Section (Astro HTML)**:
- Static HTML structure
- Conditional error/success state rendering
- Semantic markup for accessibility
- BaseLayout integration

**Client-Side Section (JavaScript)**:
- `SyncViewer` class encapsulating all realtime functionality
- SSE connection management with EventSource
- Progress bar updates with smooth animations
- Connection status indicator management
- Error handling and retry logic

**Styling Section (Scoped CSS)**:
- Component-specific styles
- Responsive design rules
- Accessibility enhancements
- Dark mode support

## Component Architecture

### SyncViewer JavaScript Class Structure

```javascript
class SyncViewer {
  // Core Properties (lines 220-240)
  constructor(sessionId, initialProgress, apiBase)

  // Connection Management (lines 241-300)
  connect()
  setupEventHandlers()
  handleConnect()
  handleError()
  handleMessage(event)

  // UI Updates (lines 301-360)
  updateProgress(progress)
  updateConnectionStatus(status)
  updateTimestamp(timestamp)
  showError(message)

  // Reconnection Logic (lines 361-420)
  scheduleReconnect()
  attemptReconnect()
  resetReconnectionState()

  // Cleanup (lines 421-440)
  destroy()
}
```

### HTML Structure Organization

```html
<!-- Main Container -->
<div class="sync-viewer-container">

  <!-- Connection Status Header -->
  <div class="connection-status">
    <span class="status-indicator"></span>
    <span class="status-text">Connected</span>
    <span class="last-update">Updated 2 seconds ago</span>
  </div>

  <!-- Progress Display -->
  <div class="progress-section">
    <h2>Session Progress</h2>
    <div class="progress-container">
      <progress class="progress-bar" value="0" max="1"></progress>
      <div class="progress-label">0.0%</div>
    </div>
  </div>

  <!-- Session Information -->
  <div class="session-info">
    <h3>Session Details</h3>
    <div class="info-grid">
      <div class="info-item">
        <label>Session ID:</label>
        <code>{session.id}</code>
      </div>
      <!-- Additional metadata -->
    </div>
  </div>

</div>
```

### CSS Architecture Layers

```css
/* Layer 1: CSS Custom Properties (lines 701-720) */
.sync-viewer-container {
  --progress-height: 24px;
  --status-indicator-size: 12px;
  --connection-green: #10b981;
  --connection-yellow: #f59e0b;
  --connection-red: #ef4444;
}

/* Layer 2: Layout Components (lines 721-800) */
.sync-viewer-container { /* Main layout */ }
.connection-status { /* Status header layout */ }
.progress-section { /* Progress area layout */ }

/* Layer 3: Interactive Elements (lines 801-880) */
.progress-bar { /* Progress styling */ }
.status-indicator { /* Connection dots */ }
.progress-label { /* Percentage display */ }

/* Layer 4: State Variants (lines 881-920) */
.status-indicator.connected { /* Green state */ }
.status-indicator.connecting { /* Yellow state */ }
.status-indicator.disconnected { /* Red state */ }

/* Layer 5: Responsive Design (lines 921-980) */
@media (max-width: 767px) { /* Mobile adaptations */ }

/* Layer 6: Accessibility (lines 981-1000) */
@media (prefers-reduced-motion: reduce) { /* Animation control */ }
@media (prefers-contrast: high) { /* High contrast support */ }
```

## Data Flow Architecture

### Initial Page Load Sequence
1. **Server-Side Rendering**: Astro processes frontmatter
2. **Session Fetch**: REST API call to get initial session data
3. **HTML Generation**: Template rendered with session data
4. **Client Hydration**: JavaScript class instantiated
5. **SSE Connection**: EventSource established to realtime API

### Real-time Update Flow
```
Admin Control Page → Updates progress via API
           ↓
PocketBase Record Update → Triggers realtime broadcast
           ↓
SSE Message → Received by EventSource
           ↓
Message Filter → Check collection/session ID match
           ↓
Progress Update → Animate progress bar
           ↓
UI Refresh → Update timestamp and status
```

### Error Recovery Flow
```
Connection Error → Increment retry counter
        ↓
Exponential Backoff → Wait (1s, 2s, 4s, 8s, 16s)
        ↓
Reconnect Attempt → Create new EventSource
        ↓
Success/Failure → Update status indicator
        ↓
Max Retries → Show manual retry button
```

## Integration Interfaces

### Astro Framework Integration
- **Routing**: Uses existing `[id]` dynamic route pattern
- **Layout**: Inherits from `BaseLayout.astro` without modifications
- **SSR**: Disabled with `export const prerender = false`
- **Client Scripts**: Uses `client:load` directive for hydration

### PocketBase API Integration
- **REST Endpoint**: `/api/collections/sync_sessions/records/{id}` for initial data
- **Realtime Endpoint**: `/api/realtime` for SSE connection
- **Message Format**: Standard PocketBase realtime message structure
- **Error Responses**: Standard HTTP status codes and JSON error format

### CSS Design System Integration
- **Variables**: Uses existing CSS custom properties from BaseLayout
- **Typography**: Inherits font families and sizing scales
- **Colors**: Uses established color palette
- **Spacing**: Uses standard spacing scale (--space-xs through --space-xl)

## Internal Component Organization

### Event Handling Structure
```javascript
// Connection Events
eventSource.onopen = () => handleConnect();
eventSource.onmessage = (event) => handleMessage(event);
eventSource.onerror = () => handleError();

// Browser Events
window.addEventListener('online', () => handleOnline());
window.addEventListener('offline', () => handleOffline());
window.addEventListener('visibilitychange', () => handleVisibilityChange());
```

### State Management Pattern
```javascript
// Centralized state object
const state = {
  sessionId: '',
  progress: 0,
  connectionStatus: 'disconnected', // 'connected', 'connecting', 'disconnected'
  lastUpdate: null,
  reconnectAttempts: 0,
  isVisible: true
};

// Single state update method
updateState(newState) {
  Object.assign(this.state, newState);
  this.render();
}
```

### Memory Management
- **EventSource Cleanup**: Proper connection closure on page unload
- **Event Listener Removal**: All browser event listeners cleaned up
- **Timer Management**: Clear timeouts and intervals on destruction
- **Reference Cleanup**: Null out DOM references to prevent leaks

## Error Boundary Definitions

### Server-Side Error Handling
- **Invalid Session ID**: Redirect to 404 page
- **Network Timeout**: Show connection error state
- **API Error Response**: Display appropriate error message
- **Session Not Found**: Show "session not found" page

### Client-Side Error Boundaries
- **SSE Connection Failure**: Automatic retry with exponential backoff
- **Malformed Messages**: Log warning, continue processing other messages
- **DOM Update Failures**: Fallback to basic progress display
- **JavaScript Errors**: Graceful degradation to static display

## Performance Optimization Structure

### Bundle Size Optimization
- **No External Dependencies**: Pure vanilla JavaScript implementation
- **Minimal DOM Queries**: Cache DOM references in constructor
- **Efficient Event Handling**: Single message handler with filtering
- **CSS Optimization**: Scoped styles with no global impact

### Runtime Performance
- **Update Throttling**: Max 30fps progress bar updates
- **Visibility API**: Reduce activity when page not visible
- **Memory Efficient**: Minimal object allocation during updates
- **Battery Conscious**: Pause non-essential updates when offline

## Testing Interfaces

### Manual Testing Hooks
- **Connection Status Display**: Visual feedback for all connection states
- **Console Logging**: Detailed logs for debugging SSE messages
- **Error Simulation**: Network tab can simulate connection failures
- **State Inspection**: Window-level debugging object exposure

### Integration Test Points
- **Session Data Fetch**: Initial API call testable via network mocking
- **SSE Message Processing**: Message handler testable with synthetic events
- **Progress Updates**: UI changes observable via DOM queries
- **Error States**: Error conditions reproducible via API/network manipulation

## Deployment Considerations

### Static Asset Generation
- **Astro Build Process**: Page included in static build output
- **CSS Inlining**: Styles embedded in page for optimal loading
- **JavaScript Bundling**: Client script included inline for performance
- **Asset Optimization**: Astro handles image and asset optimization

### Server Configuration
- **No Server Changes**: Implementation uses existing PocketBase APIs
- **Route Registration**: Automatic via Astro file-based routing
- **Cache Headers**: Astro default caching for static assets
- **CORS Configuration**: Uses existing PocketBase CORS settings

This structure provides a clean, maintainable implementation that integrates seamlessly with the existing codebase while delivering all required functionality through a single, well-organized file.