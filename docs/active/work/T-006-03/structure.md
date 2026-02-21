# T-006-03: Sync Admin Page Structure

## Overview

Structure definition for implementing the sync admin page at `/sync/[id]/control` following the progressive enhancement pattern established in the codebase.

## File Changes

### New Files

#### `frontend/src/pages/sync/[id]/control.astro`
**Purpose**: Main admin control page for sync sessions
**Type**: Astro dynamic route with SSR
**Size**: ~300 lines (template + script + styles)

**Public Interface**:
- Route: `/sync/[id]/control?token=<admin_token>`
- Path parameter: `id` (session ID)
- Query parameter: `token` (admin authentication token)

**Internal Structure**:
```astro
---
// Server-side data fetching and validation
export const prerender = false;
import BaseLayout from '../../../layouts/BaseLayout.astro';
// Session data fetching logic
// Token validation
// Error handling
---

<BaseLayout title="Sync Control">
  <!-- Rendered template -->
</BaseLayout>

<script>
  // Client-side SyncController class
</script>

<style>
  /* Component styles */
</style>
```

**Dependencies**:
- `BaseLayout.astro` (existing)
- Sync session API endpoints (existing)

### No Modified Files

All functionality is self-contained in the new route file following the established pattern.

### No Deleted Files

This is a purely additive change.

## Architecture Definition

### Route Structure

```
/sync/[id]/control
├── Server-Side (Astro frontmatter)
│   ├── Parameter extraction (id, token)
│   ├── Session data fetching
│   ├── Token validation
│   ├── Error state determination
│   └── Initial render data preparation
├── Template (Astro template)
│   ├── Error states rendering
│   ├── Session information display
│   ├── Progress control interface
│   ├── Viewer URL section
│   └── Accessibility structure
├── Client Enhancement (Script tag)
│   ├── SyncController class
│   ├── Progressive enhancement logic
│   ├── API communication
│   ├── Slider control and debouncing
│   └── Copy-to-clipboard functionality
└── Styling (Style tag)
    ├── Layout styles
    ├── Control component styles
    ├── Error state styles
    ├── Responsive design
    └── Accessibility styles
```

### Data Flow Architecture

#### Server-to-Client Flow
```
URL Parameters → Server Validation → Initial State → Template Rendering → Client Enhancement
```

#### Client-to-Server Flow
```
User Interaction → Debounced Updates → API Request → Server Update → Client State Update
```

### Component Boundaries

#### Server-Side Boundary
**Scope**: Initial data loading, authentication, error handling
**Interface**:
- Input: URL parameters (id, token)
- Output: Rendered HTML with embedded data
- Error handling: HTTP redirects, error page rendering

#### Client-Side Boundary
**Scope**: Interactive slider control, real-time updates, UI enhancements
**Interface**:
- Input: Embedded session data, DOM events
- Output: DOM updates, API requests
- Error handling: User feedback, retry mechanisms

#### API Boundary
**Scope**: Communication with sync session endpoints
**Interface**:
- Endpoint: `POST /api/sync/:id/progress?token=<token>`
- Request: `{"progress": 0.42}`
- Response: Updated session object or error
- Authentication: Query parameter token validation

## Module Organization

### Server-Side Module (Astro Frontmatter)

```typescript
// Type definitions
interface SessionData {
  id: string;
  progress: number;
  admin_token: string;
  created: string;
  updated: string;
}

interface PageState {
  session: SessionData | null;
  error: ErrorType | null;
  viewerUrl: string;
  isValidToken: boolean;
}

// Core functions
async function fetchSessionData(id: string): Promise<SessionData>
function validateToken(provided: string, session: SessionData): boolean
function generateViewerUrl(sessionId: string): string
function determineErrorState(error: any): ErrorType
```

### Client-Side Module (Script Section)

```typescript
class SyncController {
  // Core properties
  private sessionId: string;
  private adminToken: string;
  private currentProgress: number;
  private elements: {
    slider: HTMLInputElement;
    display: HTMLElement;
    copyButton: HTMLButtonElement;
  };

  // Public interface
  constructor(sessionId: string, token: string, initialProgress: number);
  public async updateProgress(progress: number): Promise<void>;
  public copyViewerUrl(): void;

  // Private methods
  private setupEventListeners(): void;
  private createThrottledUpdate(): Function;
  private handleSliderChange(event: Event): void;
  private updateDisplay(progress: number): void;
  private handleApiError(error: Error): void;
}
```

### Style Module (Style Section)

```css
/* Layout containers */
.sync-control-container { /* Main container */ }
.session-header { /* Session info section */ }
.progress-section { /* Slider controls */ }
.viewer-section { /* URL sharing */ }

/* Interactive components */
.progress-slider { /* Range input styling */ }
.progress-display { /* Value display */ }
.copy-button { /* Copy to clipboard */ }
.copy-feedback { /* Success/error feedback */ }

/* State styles */
.error-state { /* Error display */ }
.loading-state { /* Loading indicators */ }
.disabled-state { /* Disabled controls */ }

/* Responsive design */
@media (max-width: 767px) { /* Mobile styles */ }
@media (min-width: 768px) { /* Desktop styles */ }
```

## Component Interface Specifications

### Server-Side Interface

**Input Parameters**:
- `Astro.params.id`: Session ID from URL path
- `Astro.url.searchParams.get('token')`: Admin token from query string

**Output Data**:
- `session`: Session object with progress and metadata
- `viewerUrl`: Formatted viewer URL for sharing
- `error`: Error state information for rendering
- `isValidAccess`: Boolean for access control

**Error States**:
- `invalid_token`: 403-equivalent, show access denied
- `session_not_found`: 404-equivalent, show not found
- `network_error`: Connection issues, show retry options
- `server_error`: 500-equivalent, show server error message

### Client-Side Interface

**Initialization**:
```javascript
// Embedded in page via script tag data
window.syncData = {
  sessionId: "abc123",
  adminToken: "64-char-hex-string",
  initialProgress: 0.42,
  viewerUrl: "/sync/abc123"
};
```

**Event Handling**:
- `input` event on slider → throttled progress update
- `click` event on copy button → clipboard operation
- `focus`/`blur` events → accessibility enhancements

**API Communication**:
- Endpoint: `POST /api/sync/:id/progress?token=:token`
- Throttling: Maximum 30 requests per second
- Error handling: Network failures, token expiration, validation errors
- Response processing: Update UI state based on server response

### External Dependencies

#### Existing System Dependencies
- `BaseLayout.astro`: Layout shell and global styles
- `/api/sync/:id/progress`: Progress update endpoint
- CSS custom properties from BaseLayout (--color-*, --space-*)
- Navigation component integration

#### Browser API Dependencies
- `fetch()`: API communication
- `navigator.clipboard.writeText()`: Preferred copy method
- `document.execCommand('copy')`: Fallback copy method
- `URLSearchParams`: Query string parsing
- Event listeners: DOM interaction handling

## Implementation Ordering

### Phase 1: Server-Side Foundation
1. Create route file with basic Astro structure
2. Implement session data fetching logic
3. Add token validation and error handling
4. Create basic HTML template structure

### Phase 2: Static UI Implementation
1. Add session information display
2. Implement viewer URL section with static copy button
3. Create progress slider with initial value
4. Add error state templates

### Phase 3: Client-Side Enhancement
1. Implement SyncController class skeleton
2. Add slider event handling and display updates
3. Implement throttled API communication
4. Add copy-to-clipboard functionality

### Phase 4: Styling and Accessibility
1. Add component-specific styles
2. Implement responsive design
3. Add accessibility attributes and announcements
4. Test and refine visual design

### Phase 5: Error Handling and Edge Cases
1. Comprehensive client-side error handling
2. Network failure recovery
3. Token expiration handling
4. Input validation and user feedback

## Validation Criteria

### Functional Requirements
- [ ] Admin page accessible at `/sync/[id]/control?token=<token>`
- [ ] Progress slider updates server state via API
- [ ] Current progress value displayed numerically
- [ ] Viewer URL shown with functional copy button
- [ ] Invalid token shows 403-equivalent error
- [ ] Updates throttled to ~30/second maximum

### Technical Requirements
- [ ] Uses BaseLayout for consistent shell
- [ ] Follows existing SSR + progressive enhancement pattern
- [ ] No client islands or framework dependencies
- [ ] Responsive design works on mobile and desktop
- [ ] Accessible via keyboard navigation and screen readers
- [ ] Graceful degradation when JavaScript disabled

### Integration Requirements
- [ ] Integrates with existing sync API endpoints
- [ ] Consistent with established routing patterns
- [ ] Uses existing CSS custom properties and design system
- [ ] Compatible with current build and deployment process

This structure provides a clear blueprint for implementation while maintaining consistency with the existing codebase architecture and patterns.