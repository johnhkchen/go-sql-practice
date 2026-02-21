# Research: T-006-04 - sync-viewer-page

## Project Architecture

The project is structured as a full-stack application with a PocketBase backend and Astro frontend:

### Backend (Go + PocketBase)
- **Main Entry**: `main.go` - PocketBase app initialization with custom migrations and routes
- **Database**: PocketBase SQLite with realtime capabilities built-in
- **Custom Routes**: `routes/` directory contains API endpoints
- **Migrations**: `migrations/collections.go` creates database collections including `sync_sessions`

### Frontend (Astro)
- **Location**: `frontend/` directory with Astro static site generation
- **Configuration**: `frontend/astro.config.mjs` - Node.js adapter in standalone mode
- **Layout**: `frontend/src/layouts/BaseLayout.astro` - shared shell with navigation
- **Pages**: `frontend/src/pages/` - file-based routing

## Existing Sync Infrastructure

### Database Schema (sync_sessions collection)
Located in `migrations/collections.go:157-205`:
- `id` (auto-generated): Session identifier
- `progress` (float): 0.0 to 1.0 range representing sync progress
- `admin_token` (text): 64-character hex string for admin authentication
- `created` / `updated` (timestamps): Automatic PocketBase fields

API rules allow public read access but no direct CRUD operations - controlled via custom routes.

### Backend API Endpoints
Implemented in `routes/sync_sessions.go`:

1. **POST /api/sync/create**
   - Creates new session with random admin token
   - Returns session_id, admin_url, and viewer_url
   - No authentication required

2. **POST /api/sync/:id/progress**
   - Updates session progress (admin-only)
   - Requires admin token via query parameter
   - Validates progress range (0.0-1.0)
   - Uses constant-time token comparison for security

### Existing Sync Control Page
Location: `frontend/src/pages/sync/[id]/control.astro`

This comprehensive control page (1,482 lines) provides:
- Admin token validation and error handling
- Interactive progress slider with real-time updates
- Throttled API calls (max 30/sec)
- Viewer URL generation and clipboard functionality
- Accessibility features (ARIA labels, screen reader support)
- Robust error handling with retry dialogs
- Mobile-responsive design

Key JavaScript capabilities:
- Real-time progress updates via throttled fetch requests
- Advanced keyboard navigation support
- Copy-to-clipboard with fallbacks
- Network error recovery with connection monitoring
- Graceful degradation without JavaScript

## PocketBase Realtime Architecture

PocketBase provides built-in SSE (Server-Sent Events) realtime subscriptions:
- Endpoint: `/api/realtime` (PocketBase standard)
- Protocol: SSE (Server-Sent Events)
- Authentication: Not required for public read collections
- Subscription format: Collection-based with record-level filtering

### Realtime API Usage Pattern
Based on PocketBase documentation and control page patterns:
```javascript
// PocketBase JS SDK approach (not currently used in project)
pb.collection('sync_sessions').subscribe(id, callback)

// Manual SSE approach (would need implementation)
const eventSource = new EventSource(`/api/realtime?collections=sync_sessions`)
```

## Current Directory Structure

### Pages Routing
- `/` - Home page (`frontend/src/pages/index.astro`)
- `/stats` - Statistics page (`frontend/src/pages/stats.astro`)
- `/sync/[id]/control` - Admin control page (✅ exists)
- `/sync/[id]` - **TARGET: Viewer page (missing)**
- `/links/[id]` - Link detail pages
- `/tags/[slug]` - Tag pages
- `/watch/[id]` - Watch pages

### Component Architecture
- `frontend/src/layouts/BaseLayout.astro` - HTML shell, global CSS, navigation
- `frontend/src/components/Navigation.astro` - Site navigation
- `frontend/src/components/LinkCard.astro` - Reusable link display
- `frontend/src/components/StatsSummary.astro` - Statistics display

## Dependencies and Constraints

### Frontend Dependencies (package.json)
- `astro: ^5.17.3` - Static site generator
- `@astrojs/node: ^9.5.4` - Node.js adapter for SSR capability
- Node.js >= 24 required

### Missing Dependencies for Realtime
The project does NOT currently include:
- PocketBase JavaScript SDK
- Manual SSE/EventSource implementation
- WebSocket libraries
- Realtime state management

### API Configuration
- Base API URL: `import.meta.env.PUBLIC_API_URL || 'http://localhost:8090'`
- Default timeout: 5000ms for API requests
- PocketBase collection API: `/api/collections/sync_sessions/records/[id]`
- Custom progress API: `/api/sync/[id]/progress`

## Technical Constraints

### Astro Framework Limitations
- Static site generation by default
- Limited client-side reactivity without framework integration
- SSR available but configured for static output
- JavaScript interactions via `<script>` tags with `client:load` directive

### PocketBase Realtime
- Uses SSE, not WebSockets
- Collection-level subscriptions
- Automatic JSON serialization
- Requires manual connection management
- No built-in retry/reconnection logic

### Browser Compatibility
- EventSource (SSE) supported in all modern browsers
- IE/older browsers would need polyfills
- Mobile browsers fully supported
- Network connectivity changes need manual handling

## Security Considerations

### Viewer Page Security Model
- **No authentication required** - public viewing
- **Read-only access** - viewers cannot modify progress
- **Session ID exposure** - URLs contain session ID (acceptable for sharing)
- **No admin token exposure** - viewer page should not have admin capabilities

### Data Flow Security
- Admin control page: Requires admin token for updates
- Viewer page: Public read access to session data
- Realtime subscriptions: Public for sync_sessions collection
- API endpoints: Proper validation and error handling

## Integration Points

### Admin Control → Viewer Page
- Admin control page generates viewer URL: `/sync/${id}`
- Progress updates from admin are broadcasted via PocketBase realtime
- Viewer page subscribes to same session record for live updates

### Existing Error Patterns
From control page implementation:
- Network timeouts with retry functionality
- Session not found handling
- Connection status indicators
- Graceful degradation patterns

### CSS Design System
Established in BaseLayout.astro:
- CSS custom properties for colors, spacing, typography
- Responsive breakpoints (767px mobile)
- Consistent component styling patterns
- Accessibility-focused design

## Missing Implementation Pieces

### Core Functionality
1. **Viewer page file**: `frontend/src/pages/sync/[id].astro` (needs creation)
2. **Realtime subscription**: SSE connection to PocketBase
3. **Progress display**: HTML progress element with live updates
4. **Connection status**: Online/offline/reconnecting indicators

### Technical Gaps
1. **PocketBase JS SDK**: Not included in dependencies
2. **Manual SSE implementation**: EventSource connection management
3. **Error recovery**: Network failure and reconnection logic
4. **Mobile optimization**: Touch-friendly progress display

## Success Dependencies

### Completed Dependencies
- ✅ T-004-01: BaseLayout and Navigation components exist
- ✅ T-006-02: Sync API routes and database schema implemented

### Required for Implementation
- PocketBase realtime API understanding
- SSE/EventSource implementation pattern
- Progress bar styling consistent with control page
- Error handling patterns from existing control page
- Mobile-responsive design matching project standards