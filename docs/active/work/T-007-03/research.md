# Research: T-007-03 Presenter Dashboard

## Overview

This ticket requires implementing a comprehensive presenter dashboard at `/present` for managing and controlling live presentations. The dashboard needs to integrate with existing PocketBase APIs and provide real-time presentation controls with step-by-step navigation.

## Current Codebase Architecture

### Technology Stack
- **Backend**: Go with PocketBase framework (`main.go:12`)
- **Frontend**: Astro 5.17.3 with Node.js 24+ (`frontend/package.json:13`)
- **Database**: PocketBase embedded SQLite with collections
- **Styling**: CSS custom properties with component-scoped styles
- **Client Interactivity**: Astro client islands with vanilla JavaScript

### Database Schema

#### Presentations Collection (`migrations/collections.go:228`)
- `id`: Auto-generated PocketBase ID
- `name`: Text field, required, 1-255 chars
- `step_count`: Number field, required, minimum 1, integer only
- `step_labels`: JSON array field, optional (stores step names)
- `active_session`: Relation to sync_sessions, max 1, optional
- `created_by`: Relation to users, max 1, optional (for ownership)

**API Rules:**
- List/View: Public access (`""`)
- Create: Authenticated users only (`"@request.auth.id != ''`)
- Update/Delete: Owner only (`"@request.auth.id = created_by`)

#### Sync Sessions Collection (`migrations/collections.go:157`)
- `id`: Auto-generated PocketBase ID
- `progress`: Float field, 0.0-1.0 range, optional
- `admin_token`: Text field, required, exactly 64 chars (hex-encoded 32 bytes)

**API Rules:**
- List/View: Public access (`""`)
- Create/Update/Delete: Restricted (handled by custom routes)

### Existing API Endpoints

#### Presentations API (`routes/presentations.go`)

**GET /api/presentations/:id/status** (`presentations.go:47`)
- Returns enhanced status with live session info
- Includes: id, name, step_count, step_labels, is_live, progress, current_step
- Uses `buildStatusResponse()` to compute current step from progress

**POST /api/presentations/:id/live** (`presentations.go:57`)
- Starts live presentation session
- Requires authentication and ownership validation
- Creates new sync_session with admin token
- Returns: session_id, admin_url, viewer_url, step_count, step_labels
- Admin URL format: `/sync/{session_id}/control?token={admin_token}`
- Viewer URL format: `/watch/{presentation_id}`

**POST /api/presentations/:id/stop** (`presentations.go:53`)
- Stops live presentation session
- Requires authentication and ownership validation
- Clears `active_session` field from presentation

#### Sync Sessions API (`routes/sync_sessions.go`)

**POST /api/sync/:id/progress** (`sync_sessions.go:40`)
- Updates session progress (0.0-1.0)
- Requires admin token authentication
- Token passed as query parameter or in request body
- Validates progress range and performs constant-time token comparison

### Step-Progress Conversion Logic (`routes/presentations.go:95`)

**Progress to Step Formula:**
```javascript
step_index = round(progress * (step_count - 1)) for step_count > 1, else 0
```

**Step to Progress Formula:**
```javascript
progress = step_index / (step_count - 1) for step_count > 1, else 0.0
```

This creates discrete step boundaries: 0, 0.5, 1.0 for 3 steps (0, 1, 2).

### Frontend Architecture

#### Layout System (`frontend/src/layouts/BaseLayout.astro`)
- Global CSS custom properties for theming
- Responsive design with mobile breakpoints
- Navigation integration and footer
- Semantic HTML structure with accessibility considerations

#### Navigation Component (`frontend/src/components/Navigation.astro`)
- Sticky header with mobile hamburger menu
- Current links: Home (`/`), Stats (`/stats`)
- CSS-only mobile toggle functionality
- Missing presenter dashboard link

#### Page Structure
- `frontend/src/pages/index.astro`: Basic welcome page
- `frontend/src/pages/stats.astro`: Statistics page
- `frontend/src/pages/sync/[id]/control.astro`: Existing sync control interface
- Missing: `/present` dashboard and `/present/[id]` control views

### Existing Control Interface Analysis

The existing sync control at `frontend/src/pages/sync/[id]/control.astro` provides:

**Server-Side Features:**
- Session data fetching from PocketBase API
- Admin token validation
- Error handling for missing/invalid tokens
- Session information display

**Client-Side Features (`control.astro:182`):**
- Real-time progress slider (0-1 range)
- Throttled API updates (30 updates/sec max)
- Progress validation and error handling
- Viewer URL with copy-to-clipboard functionality
- Visual feedback for API calls and copy operations
- Responsive design with mobile optimizations

**JavaScript Architecture:**
- Single `SyncController` class with DOM management
- Event-driven progress updates
- Robust error handling with user feedback
- Accessibility features (screen reader announcements)

### API Integration Patterns

**Authentication Flow:**
1. PocketBase built-in auth for presentation management
2. Admin tokens for sync session control
3. Constant-time token comparison for security

**Data Fetching:**
- Server-side API calls in Astro components
- Client-side updates via fetch API
- Error handling with user-friendly messages
- Timeout handling (5s default)

**URL Structure:**
- Admin control: `/sync/{session_id}/control?token={admin_token}`
- Public viewer: `/watch/{presentation_id}`
- API base: `http://localhost:8090` (configurable via env)

## Key Patterns and Constraints

### Astro Conventions
- TypeScript interfaces for props
- Server-side rendering with `export const prerender = false`
- Client islands for interactive components
- Component-scoped CSS with CSS custom properties
- Semantic HTML with ARIA labels

### Security Considerations
- Admin tokens are 64-character hex strings (32 bytes)
- Constant-time token comparison to prevent timing attacks
- Presentation ownership validation before control access
- Public read access to presentations for viewer functionality

### Error Handling
- Structured error objects with user-friendly messages
- Network timeout handling (5s default)
- Graceful degradation for missing data
- Screen reader compatibility for error states

### Performance Optimizations
- Throttled API updates (33ms minimum interval)
- Immediate UI updates with async API calls
- Timeout handling to prevent hanging requests
- Component-level CSS to minimize bundle size

## Dependencies and Integration Points

### Required Dependencies
- Presentations must exist before going live
- Sync sessions created only through presentation workflow
- Admin tokens required for control access
- Authentication required for presentation management

### Missing Components
1. **Presenter Dashboard Page** (`/present`): Main listing and management interface
2. **Presenter Control Page** (`/present/[id]`): Step-by-step control interface
3. **Navigation Updates**: Link to presenter dashboard
4. **New Presentation Form**: Create presentations with step configuration

### Integration Requirements
- PocketBase built-in CRUD API for presentations listing/creation
- Custom presentation lifecycle API for live session management
- Sync progress API for real-time control
- Existing control interface patterns for consistency

## Technical Considerations

### Step Navigation Logic
The presenter control needs to convert between:
- **Step Index** (0-based integer): User-facing step numbers
- **Progress Value** (0.0-1.0 float): API storage format
- **Step Labels** (string array): Human-readable step names

### Real-time Updates
- Progress updates flow: UI → API → Database
- No real-time sync between multiple clients (by design)
- Progress persistence allows session recovery

### Mobile Responsiveness
- Existing patterns use CSS Grid and Flexbox
- Mobile-first breakpoint at 767px
- Touch-friendly controls for sliders and buttons

### Accessibility
- ARIA labels and semantic HTML throughout
- Screen reader announcements for dynamic updates
- Keyboard navigation support
- High contrast mode considerations

This research provides the foundation for designing a presenter dashboard that integrates seamlessly with the existing PocketBase architecture while following established frontend patterns and maintaining security best practices.