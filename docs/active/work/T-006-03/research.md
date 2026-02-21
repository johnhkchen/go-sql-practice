# Research: sync-admin-page (T-006-03)

## Context

This ticket requires creating an Astro admin page for controlling sync sessions at `/sync/[id]/control` with token-based authentication. The page displays a progress slider (0-1 range) that updates the server via API, shows current progress, displays viewer URL for sharing, and handles error states.

## Current State

### Project Architecture

**Frontend Structure**:
- Astro v5.17.3 with Node adapter in static output mode
- Directory structure: `frontend/src/pages/`, `frontend/src/components/`, `frontend/src/layouts/`
- Existing BaseLayout provides HTML shell with navigation, global CSS variables, mobile responsive design
- CSS custom properties system: colors, spacing, typography, layout constraints
- Animation utilities in `src/styles/animations.css` with reduced motion support

**Existing Page Patterns**:

1. **Dynamic Routes**: `/watch/[id].astro` demonstrates:
   - Server-side rendering with `prerender = false`
   - URL parameter extraction via `Astro.params`
   - API integration with error handling and timeouts
   - Progressive enhancement with client-side JavaScript
   - State-based conditional rendering (error, live, waiting states)

2. **BaseLayout Integration**: All pages use consistent:
   - HTML head with meta tags, title prop
   - Global CSS variables and reset styles
   - Navigation component with responsive mobile menu
   - Footer with consistent styling

3. **Error Handling Patterns**:
   - Comprehensive error states: notfound, network, server, timeout
   - User-friendly error messages with suggested actions
   - Fallback UI for failed API requests

### API Infrastructure

**Backend Architecture**:
- Go-based PocketBase application with custom routes
- Routes registered in `routes/routes.go` with standardized pattern
- API endpoints served from Go server, frontend assets served statically

**Sync Session API** (implemented by T-006-02):
- `POST /api/sync/create`: Creates session with admin token, returns URLs
- `POST /api/sync/:id/progress`: Updates progress with token validation
- Admin tokens are 64-character hex strings (32 bytes encoded)
- Progress validation: float between 0.0 and 1.0 inclusive
- Token validation uses constant-time comparison for security

**API Response Format**:
```json
// Create session response
{
  "session_id": "abc123",
  "admin_url": "/sync/abc123/control?token=<admin_token>",
  "viewer_url": "/sync/abc123"
}

// Update progress request body
{"progress": 0.42}

// Update progress response
{
  "id": "abc123",
  "progress": 0.42,
  "admin_token": "...",
  "created": "2024-...",
  "updated": "2024-..."
}
```

**Error Responses**:
- 403: Invalid/missing admin token
- 404: Session not found
- 400: Invalid progress value (outside 0-1 range)
- 500: Server errors

### Database Schema

**sync_sessions Collection** (from migrations/collections.go):
- `progress`: NumberField (0.0-1.0 range, float precision)
- `admin_token`: TextField (required, exactly 64 chars)
- Standard PocketBase fields: `id`, `created`, `updated`
- Public read rules, no standard create/update rules (handled by custom routes)

## Technical Constraints

### Astro Framework Limitations

**Static vs SSR**:
- Current config uses `output: 'static'` but existing dynamic route uses `prerender = false`
- Individual pages can override to SSR for dynamic content
- Admin page needs SSR for token validation and error handling

**Client-Side Interactivity**:
- Astro islands required for interactive components (`client:load` directive)
- Vanilla JavaScript recommended for lightweight interactivity
- No additional JS frameworks in current setup

**API Integration Pattern**:
- Server-side fetch in Astro component for initial data
- Client-side fetch for real-time updates
- Environment variable pattern: `import.meta.env.PUBLIC_API_URL || 'http://localhost:8090'`

### Security Considerations

**Token Exposure**:
- Admin token passed via URL query parameter (visible in browser history)
- Token required for all update requests
- No additional authentication layers specified

**Client-Side Validation**:
- Progress range validation should occur client-side before API calls
- Server validates again, client validation is UX enhancement only

**CORS and Request Headers**:
- API expects JSON content-type for progress updates
- No authentication headers beyond token query parameter

### Performance Requirements

**Debouncing/Throttling**:
- Ticket specifies ~30 updates/sec maximum to prevent API flooding
- Need debouncing or throttling mechanism for slider input events
- Balance between responsiveness and server load

**Real-time Updates**:
- No requirement for receiving updates from other sessions
- Admin page is write-only interface (except for initial state)

## Implementation Considerations

### URL Structure and Routing

**Dynamic Route Pattern**:
- File: `frontend/src/pages/sync/[id]/control.astro`
- Nested directory structure for cleaner URLs
- Token validation via query parameter extraction

**Parameter Extraction**:
```javascript
const { id } = Astro.params;
const token = Astro.request.url.searchParams.get('token');
```

### State Management Strategy

**Initial Page Load**:
1. Extract session ID and token from URL
2. Validate token exists (redirect to error if missing)
3. Optionally fetch current session state for display
4. Render page with initial progress value

**Error State Handling**:
- Invalid/missing token: 403 error state
- Session not found: 404 error state
- Network errors: Connection error state
- Server errors: Generic server error state

### Interactive Component Design

**Slider Implementation**:
- HTML5 range input: `<input type="range" min="0" max="1" step="0.001">`
- Value display: Real-time number display next to slider
- Debounced API updates to prevent excessive requests

**Copy-to-Clipboard Feature**:
- Viewer URL generation: `/sync/${sessionId}` (no token needed)
- Modern Clipboard API with fallback for older browsers
- Success/error feedback for copy operation

**Client Island Structure**:
```astro
<SyncControl
  client:load
  sessionId={id}
  adminToken={token}
  initialProgress={progress}
/>
```

### API Integration Pattern

**Debouncing Strategy**:
- Use setTimeout to delay API calls
- Cancel previous timeout on new input
- Immediate update for UI feedback, delayed server sync

**Error Recovery**:
- Retry failed requests with exponential backoff
- Show error states without blocking UI
- Allow manual retry actions

**Request Format**:
```javascript
fetch(`${API_BASE}/api/sync/${sessionId}/progress?token=${token}`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ progress: parseFloat(value) })
});
```

### CSS and Styling Integration

**Component Styling**:
- Follow existing design system in BaseLayout
- Use CSS custom properties for consistent theming
- Mobile-first responsive design approach

**Animation Considerations**:
- Respect `prefers-reduced-motion` media query
- Smooth transitions for slider and progress display
- Loading states for API operations

## Dependencies Analysis

### T-004-01 (astro-layout-and-nav) - Completed

Provides foundation for admin page:
- BaseLayout component with HTML shell
- Navigation component (though admin page may not need nav)
- CSS design system with custom properties
- Mobile responsive patterns

### T-006-02 (sync-api-routes) - Completed

Provides API endpoints for admin functionality:
- Session creation with token generation
- Progress update endpoint with validation
- Error response patterns
- Token security implementation

### Environment Configuration

**API Base URL**:
- Development: `http://localhost:8090` (default)
- Production: Via `PUBLIC_API_URL` environment variable
- Same pattern as existing pages

## File System Impact

### New Files Required

**Primary Page**:
- `frontend/src/pages/sync/[id]/control.astro`

**Potential Component** (if complex):
- `frontend/src/components/SyncControl.astro` (for client island)

### Directory Structure

```
frontend/src/
├── pages/
│   ├── index.astro               # Existing
│   ├── stats.astro               # Existing
│   ├── watch/[id].astro          # Existing (similar pattern)
│   └── sync/
│       └── [id]/
│           └── control.astro     # New admin page
├── components/
│   ├── Navigation.astro          # Existing
│   └── SyncControl.astro         # Potential new component
└── layouts/
    └── BaseLayout.astro          # Existing
```

## Risk Analysis

### Authentication Model

**Token Security**:
- Admin token exposed in URL (browser history, server logs)
- No token expiration or refresh mechanism
- Single token grants full admin access to session

**Access Control**:
- No user authentication beyond token possession
- Token sharing gives admin access to anyone
- No audit trail for admin actions

### Client-Side Implementation

**Browser Compatibility**:
- Range input requires modern browser support
- Clipboard API may need fallbacks
- Fetch API requires polyfills for very old browsers

**JavaScript Reliability**:
- Page should function without JavaScript (graceful degradation)
- Server-side rendering should show current state
- Progressive enhancement for interactivity

### Performance and UX

**Network Reliability**:
- Admin operations depend on stable connection
- Need offline handling for poor connectivity
- Retry mechanisms for failed updates

**Concurrent Usage**:
- Multiple admin sessions could conflict
- Last-write-wins for progress updates
- No conflict resolution mechanism

## Next Steps

Design phase should address:
1. Detailed component architecture (single page vs islands)
2. Specific debouncing implementation (timeout values, strategies)
3. Error state UI design and user flows
4. Clipboard functionality implementation details
5. Progressive enhancement strategy
6. Mobile UX considerations for range slider
7. Loading and feedback state designs
8. Integration with existing design system