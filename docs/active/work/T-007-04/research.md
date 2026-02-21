# Research: T-007-04 - presentation-viewer-page

## Executive Summary

This research documents the existing infrastructure and patterns for implementing the presentation viewer page at `/watch/[id]`. The codebase already contains substantial infrastructure from dependencies T-006-04 (sync viewer) and T-007-02 (presentation API routes), plus a basic `/watch/[id]` stub that needs enhancement. The implementation requires extending the existing sync viewer pattern with step-aware presentation features.

## Current State Analysis

The codebase contains substantial infrastructure for implementing the presentation viewer page. Key findings:

1. **Existing `/watch/[id]` stub** - Basic implementation that polls for presentation status changes
2. **Proven sync viewer pattern** - Complete realtime sync infrastructure from T-006-04
3. **Step conversion utilities** - Full TypeScript library for progress/step calculations
4. **Database schema** - Presentations linked to sync_sessions via active_session field
5. **API endpoints** - Enhanced presentation status with computed step information

## 1. Existing Sync Infrastructure (from S-006/T-006-04)

### Sync Viewer Implementation
**Location:** `/home/jchen/repos/go-sql-practice/frontend/src/pages/sync/[id].astro`

**Architecture:**
- **Server-side fetching:** Initial session data loaded via fetch to `/api/collections/sync_sessions/records/{id}`
- **Realtime subscription:** Server-Sent Events (SSE) connection to `/api/realtime`
- **Connection states:** connected, connecting, disconnected with visual indicators
- **Progress display:** HTML `<progress>` element with 0-1 values, formatted as percentage
- **Error handling:** Comprehensive error states (not found, network, timeout, server error)

**Key JavaScript Class: SyncViewer**
```javascript
class SyncViewer {
  constructor(sessionId, initialProgress, apiBase)
  // Methods:
  - connect(): Establishes SSE connection
  - handleMessage(event): Filters for sync_sessions updates
  - updateProgress(progress): Updates UI elements
  - updateConnectionStatus(status): Visual connection indicators
}
```

**Realtime Message Filtering:**
- Listens for: `collection === 'sync_sessions'`
- Filters by: `data.record.id === this.sessionId`
- Updates: `data.record.progress` value (0.0-1.0)

**Connection Management:**
- EventSource to `/api/realtime`
- Automatic reconnection logic with attempt limits
- Visual status indicators with CSS classes (.connected, .connecting, .disconnected)

### Backend Sync Sessions
**Location:** `/home/jchen/repos/go-sql-practice/routes/sync_sessions.go`

**API Endpoints:**
- `POST /api/sync/create`: Creates session with admin token
- `POST /api/sync/:id/progress`: Updates progress (requires admin token)

**Data Model (sync_sessions):**
```go
// Fields in PocketBase collection
progress: float64     // 0.0 to 1.0
admin_token: string   // 64-char hex token for admin access
created: datetime
updated: datetime
```

## 2. Presentation Structure (from S-007/T-007-02)

### Data Models
**Location:** `/home/jchen/repos/go-sql-practice/migrations/collections.go`

**Presentations Collection:**
```go
// PocketBase fields
name: string          // Presentation title
step_count: int       // Number of discrete steps
step_labels: JSON     // Optional array of step names
active_session: relation  // Links to sync_sessions
created_by: relation  // Links to users (optional)
```

**Active Session Relationship:**
- `presentations.active_session` -> `sync_sessions.id`
- One-to-one relationship (MaxSelect: 1)
- Nullable (presentations can exist without active sessions)

### Step-Progress Conversion Logic
**Location:** `/home/jchen/repos/go-sql-practice/frontend/src/utils/stepConversion.ts`

**Key Functions:**
```typescript
stepToProgress(stepIndex: number, stepCount: number): number
// Formula: stepIndex / (stepCount - 1)

progressToStep(progress: number, stepCount: number): number
// Formula: Math.round(progress * (stepCount - 1))

formatStepDisplay(stepIndex: number, stepCount: number, stepLabels?: string[]): string
// Returns: "Step 3 of 5 — Label" format

getNavigationState(currentStep: number, stepCount: number): StepNavigationState
// Returns: { currentStep, canGoPrevious, canGoNext, totalSteps }
```

**Backend Implementation:**
**Location:** `/home/jchen/repos/go-sql-practice/routes/presentations.go`
- Same formulas implemented in Go for server-side calculations
- Used in presentation status API responses

### Presentation API Routes
**Location:** `/home/jchen/repos/go-sql-practice/routes/presentations.go`

**Key Endpoints:**
- `GET /api/presentations/:id/status`: Returns enhanced status with live session info
- `POST /api/presentations/:id/live`: Creates sync session and links to presentation
- `POST /api/presentations/:id/stop`: Clears active session

**Response Models:**
```go
type StatusResponse struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    StepCount   int      `json:"step_count"`
    StepLabels  []string `json:"step_labels"`
    IsLive      bool     `json:"is_live"`
    Progress    *float64 `json:"progress,omitempty"`     // When live
    CurrentStep *int     `json:"current_step,omitempty"` // When live
}
```

**Viewer URL Pattern:**
- From `buildStartLiveResponse()`: `/watch/{presentationId}`
- NOT `/sync/{sessionId}` - uses presentation ID, not session ID

## 3. Frontend Architecture Patterns

### Astro Page Structure
**Pattern from existing pages:**
```astro
---
// Server-side data fetching
export const prerender = false;

// API calls with timeout and error handling
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
const FETCH_TIMEOUT = 5000;

// Error handling with comprehensive error types
let data = null;
let error = null;
---

<!-- HTML with error states and loading states -->
{error ? <ErrorComponent /> : data ? <DataComponent /> : <LoadingComponent />}

<!-- Client-side script with global variables -->
<script client:load define:vars={{...}}>
  // Initialize JavaScript classes
</script>
```

### Component Architecture
**PresenterController Pattern:** `/home/jchen/repos/go-sql-practice/frontend/src/components/PresenterController.astro`

**Structure:**
- **Props interface:** TypeScript interface for component props
- **Server-side logic:** Data preparation in frontmatter
- **Template:** Astro template with conditional rendering
- **Styles:** Scoped CSS with responsive design
- **Client script:** TypeScript class for interactivity

**Key Features:**
- Step navigation controls (previous/next buttons)
- Step jump buttons (1, 2, 3, etc.)
- Fine-grained progress slider
- Connection status indicators
- URL sharing functionality
- Loading overlays and status messages

### CSS Design System
**Location:** `/home/jchen/repos/go-sql-practice/frontend/src/layouts/BaseLayout.astro`

**CSS Variables:**
```css
:root {
  /* Colors */
  --color-bg: #ffffff;
  --color-text: #333333;
  --color-primary: #111111;
  --color-border: #e5e5e5;
  --color-footer: #f5f5f5;

  /* Spacing */
  --space-xs: 0.25rem;
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 2rem;
  --space-xl: 3rem;

  /* Layout */
  --max-width: 1200px;
  --header-height: 60px;
}
```

**Responsive Design Patterns:**
- Mobile-first approach with `@media (max-width: 768px)`
- Flexible grid layouts with CSS Grid
- Touch-friendly button sizes
- Accessible focus states

### JavaScript Architecture Patterns

**Class-based Components:**
- Constructor takes initialization data
- Private methods for DOM manipulation
- Event listener setup in `setupDOM()`
- Async API calls with proper error handling
- Throttling for frequent updates (progress slider)

**State Management:**
- Local component state (no global state management)
- Real-time updates via SSE
- Optimistic UI updates with server reconciliation

**Accessibility Features:**
- ARIA attributes (`aria-live`, `aria-labelledby`)
- Screen reader announcements
- Keyboard navigation support
- High contrast and reduced motion support

## 4. Dependencies Analysis

### T-006-04 (Sync Viewer Page) - ✅ COMPLETED
**What it provides:**
- Basic sync viewer at `/sync/[id]`
- SSE connection patterns to PocketBase realtime API
- Progress bar rendering and updates
- Connection status handling
- Error state management
- Mobile-responsive design

**Reusable Components:**
- SyncViewer class structure
- SSE connection logic
- Progress bar styling
- Error message patterns
- Connection status indicators

### T-007-02 (Presentation API Routes) - ✅ COMPLETED
**What it provides:**
- Presentation CRUD via PocketBase collections API
- Live session lifecycle (`/live`, `/stop`, `/status` endpoints)
- Step-to-progress conversion formulas
- Presentation-to-session relationship management
- Enhanced status responses with computed fields

**Available APIs:**
- `GET /api/presentations/:id/status` - Get presentation with live session info
- `GET /api/collections/presentations/records/:id` - Basic presentation data
- PocketBase realtime updates for both presentations and sync_sessions

## 5. Integration Architecture for /watch/[id]

### Data Flow Requirements
1. **URL Parameter:** Extract presentation ID from `/watch/[id]`
2. **Initial Data:** Fetch presentation data via `/api/presentations/:id/status`
3. **Session Discovery:** If `is_live: true`, extract session ID from status response
4. **Realtime Subscription:** Subscribe to sync_sessions updates for the active session
5. **Step Computation:** Convert progress updates to current step using stepConversion utilities
6. **UI Updates:** Display both progress bar and step information

### Component Reuse Strategy
**From SyncViewer:**
- SSE connection management
- Progress bar component
- Connection status indicators
- Error handling patterns

**From PresenterController:**
- Step display formatting
- Responsive design patterns
- Loading overlays
- Status messaging

**New Components Needed:**
- PresentationViewer class (extends SyncViewer patterns)
- Step-aware progress display
- Presentation metadata display
- "Not live" state handling

### API Integration Points
**Initial Load:**
```javascript
// Get presentation status (includes live session info)
const response = await fetch(`${API_BASE}/api/presentations/${presentationId}/status`);
const presentation = await response.json();

if (presentation.is_live) {
  // Extract active session ID for SSE subscription
  const sessionId = presentation.active_session_id; // Need to verify field name
}
```

**Realtime Updates:**
```javascript
// Subscribe to sync_sessions updates (same as existing sync viewer)
eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  if (data.collection === 'sync_sessions' && data.record.id === sessionId) {
    const progress = data.record.progress;
    const currentStep = progressToStep(progress, presentation.step_count);
    updateUI(progress, currentStep);
  }
};
```

## 6. Implementation Constraints and Assumptions

### Discovered Constraints
1. **URL Pattern:** Backend expects `/watch/{presentationId}`, not `/sync/{sessionId}`
2. **Session Discovery:** Must derive session ID from presentation status, not URL
3. **Realtime Scope:** Still subscribe to sync_sessions updates, not presentation updates
4. **Step Formula:** Must match backend implementation exactly for consistency
5. **Error States:** Need "presentation not found" and "presentation not live" states

### Technical Assumptions
1. **PocketBase Realtime:** SSE connection patterns remain the same
2. **Session Lifetime:** Active sessions persist until explicitly stopped
3. **Progress Mapping:** Linear mapping between progress and steps
4. **Mobile Support:** Must work on touch devices
5. **Accessibility:** Must support screen readers and keyboard navigation

### Integration Assumptions
1. **API Consistency:** Presentation status API returns consistent data structure
2. **Session Relationship:** `active_session` field contains valid session ID when live
3. **Step Labels:** Optional step_labels array may be empty or null
4. **Connection Stability:** SSE connections can be interrupted and reconnected

## 7. Key File Locations

### Frontend Files
- **Existing sync viewer:** `/frontend/src/pages/sync/[id].astro`
- **Step utilities:** `/frontend/src/utils/stepConversion.ts`
- **Presenter controller:** `/frontend/src/components/PresenterController.astro`
- **Base layout:** `/frontend/src/layouts/BaseLayout.astro`
- **Navigation:** `/frontend/src/components/Navigation.astro`

### Backend Files
- **Sync session routes:** `/routes/sync_sessions.go`
- **Presentation routes:** `/routes/presentations.go`
- **Database migrations:** `/migrations/collections.go`

### Documentation Files
- **T-006-04 ticket:** `/docs/active/tickets/T-006-04.md`
- **T-007-02 ticket:** `/docs/active/tickets/T-007-02.md`
- **Implementation progress:** Various `/docs/active/work/T-*/` directories

## 8. Next Steps for Implementation

### Phase 1: Basic Structure
1. Create `/frontend/src/pages/watch/[id].astro`
2. Implement presentation data fetching with error handling
3. Create basic HTML structure with presentation info display

### Phase 2: Realtime Integration
1. Adapt SyncViewer class to PresentationViewer
2. Implement session ID discovery from presentation status
3. Add step-aware progress updates using stepConversion utilities

### Phase 3: UI Enhancement
1. Add step indicator component
2. Implement "not live" state display
3. Add responsive design and accessibility features

### Phase 4: Testing & Polish
1. Test error scenarios (presentation not found, not live, network issues)
2. Verify mobile responsiveness
3. Ensure accessibility compliance
4. Performance optimization for large step counts