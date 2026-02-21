# Research: Install PocketBase JS SDK (T-009-06)

## Project Architecture Overview

The project is a Go application built with PocketBase (v0.36.5) as the backend that serves an Astro frontend. The backend runs on port 8090 and provides PocketBase's standard collections API along with custom routes.

### Backend Structure
- **Main application**: Go app using PocketBase framework (`main.go`)
- **Routes**: Custom route handlers in `/routes/` directory
- **Frontend serving**: Static files served from embedded frontend build
- **Collections**: PocketBase collections for links, tags, presentations, sync_sessions
- **Realtime**: PocketBase's built-in `/api/realtime` SSE endpoint for live updates

### Frontend Structure (Astro + TypeScript)
- **Framework**: Astro v5.17.3 with TypeScript
- **Location**: `/frontend/` directory with standard Astro structure
- **Pages**: SSR-enabled pages in `src/pages/`
- **Components**: Reusable components in `src/components/`
- **No external npm packages**: Currently only has Astro and Node adapter

## Current PocketBase API Usage Patterns

### 1. Collections API Usage
Multiple pages use PocketBase's standard collections endpoints:

**Data Fetching:**
```typescript
// Pattern used across sync/[id], present/index, links/[id], etc.
const response = await fetch(`${API_BASE}/api/collections/sync_sessions/records/${id}`)
const response = await fetch(`${API_BASE}/api/collections/presentations/records?page=1&perPage=50&sort=-created`)
const response = await fetch(`${API_BASE}/api/collections/links/records/${id}`)
```

**Common endpoints used:**
- `/api/collections/sync_sessions/records/*` - Session data
- `/api/collections/presentations/records/*` - Presentation data
- `/api/collections/links/records/*` - Link data
- `/api/collections/tags/records/*` - Tag data

### 2. Custom API Routes
Backend also provides custom routes handled by Go:
- `/api/presentations/{id}/status` - Enhanced presentation status
- `/api/sync/{sessionId}/progress` - Progress updates
- `/api/links/search` - Search functionality
- `/api/stats` - Statistics aggregation

### 3. Error Handling Patterns
Standard pattern across pages:
```typescript
const FETCH_TIMEOUT = 5000;
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT);

try {
  const response = await fetch(url, { signal: controller.signal });
  if (!response.ok) throw new Error(`HTTP ${response.status}`);
  const data = await response.json();
} catch (err) {
  // Handle AbortError for timeouts, network errors, etc.
}
```

## Current EventSource (SSE) Implementation

### 1. Real-time Connection Pattern
Two main pages implement EventSource for real-time updates:

**Sync Viewer (`/sync/[id]`):**
- Connects to `/api/realtime` via EventSource
- Listens for `sync_sessions` collection updates
- Filters messages by session ID to update progress bar
- Manual reconnection logic with max attempts (5)

**Presentation Viewer (`/watch/[id]`):**
- Also connects to `/api/realtime` via EventSource
- Listens for both `presentations` and `sync_sessions` collection updates
- More complex state management (waiting → starting → live → ended)
- Exponential backoff reconnection (up to 10 attempts)

### 2. Message Filtering Logic
```typescript
// Both pages filter PocketBase realtime messages
const data = JSON.parse(event.data);

// For sync progress updates
if (data.collection === 'sync_sessions' &&
    data.action === 'update' &&
    data.record.id === this.sessionId) {
  this.updateProgress(data.record.progress);
}

// For presentation state changes
if (data.collection === 'presentations' &&
    data.action === 'update' &&
    data.record.id === this.presentationId) {
  this.handlePresentationStateChange(data.record);
}
```

### 3. Connection Management Issues
Current manual EventSource management includes:
- **Reconnection logic**: Custom exponential backoff implemented manually
- **Connection state tracking**: Manual status indicators (connected/connecting/disconnected)
- **Error handling**: Basic error detection without detailed connection status
- **Cleanup**: Manual `eventSource.close()` on page unload

## Frontend Development Patterns

### 1. Client-Side Scripts
Pages use Astro's `<script client:load>` for client-side functionality:
- Complex classes for managing EventSource connections
- Manual DOM manipulation for UI updates
- Inline utility functions (progressToStep, formatStepDisplay)

### 2. Configuration Management
Environment variables handled manually:
```typescript
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
```

### 3. TypeScript Interfaces
Each page defines its own interfaces for PocketBase responses:
```typescript
interface PocketBaseResponse {
  page: number;
  perPage: number;
  totalItems: number;
  totalPages: number;
  items: Presentation[];
}
```

## Analysis of Current Limitations

### 1. Manual Connection Management
- Each page reimplements EventSource connection logic
- No automatic retry or reconnection built into PocketBase's browser EventSource
- Custom reconnection logic varies between pages (5 vs 10 max attempts)
- No standard connection state management

### 2. Repetitive Code Patterns
- Same PocketBase collections API fetch patterns repeated across pages
- Similar error handling logic duplicated
- Manual timeout management on every fetch
- Inconsistent response interface definitions

### 3. Type Safety Gaps
- No centralized typing for PocketBase collections
- Manual JSON parsing of realtime messages
- Potential runtime errors from schema mismatches

### 4. Authentication Blind Spot
- No authentication handling visible in current code
- Could become issue if auth is needed later
- Manual header management for API calls

## Key Files for SDK Integration

### High-Impact Files (require EventSource → SDK conversion):
- `frontend/src/pages/sync/[id].astro` - Sync session viewer
- `frontend/src/pages/watch/[id].astro` - Presentation viewer

### Medium-Impact Files (can benefit from typed client):
- `frontend/src/pages/present/index.astro` - Presentations list
- `frontend/src/pages/links/[id].astro` - Link detail pages
- `frontend/src/components/StatsSummary.astro` - Statistics component

### Configuration Files:
- `frontend/package.json` - Package dependencies
- `frontend/src/env.d.ts` - Environment type definitions (if exists)

## PocketBase Server Capabilities

Based on `go.mod` and backend structure:
- **Version**: PocketBase v0.36.5 (current/recent version)
- **Realtime**: Built-in SSE at `/api/realtime`
- **Collections**: Standard CRUD operations via `/api/collections/*`
- **Authentication**: PocketBase auth system available but not currently used
- **File uploads**: PocketBase file handling available
- **Hooks**: Server-side hooks for custom logic (routes package handles custom endpoints)

## Dependencies and Constraints

### Current Frontend Dependencies:
```json
{
  "dependencies": {
    "@astrojs/node": "^9.5.4",
    "astro": "^5.17.3"
  }
}
```

### Constraints:
- **No build changes needed**: Astro already handles TypeScript and modern JS
- **Server compatibility**: PocketBase v0.36.5 should work with current PocketBase JS SDK
- **API stability**: Using standard PocketBase endpoints, so SDK should be compatible
- **No breaking changes**: Existing functionality should remain unchanged during migration

## Success Criteria Mapping

The ticket requirements map to current code as follows:

1. **"pocketbase npm package installed"** → Add to frontend/package.json
2. **"Create shared PocketBase client instance"** → New file `frontend/src/lib/pb.ts`
3. **"Refactor /sync/[id] to use pb.collection().subscribe()"** → Replace EventSource in sync/[id].astro
4. **"Refactor /watch/[id] to use SDK subscriptions"** → Replace EventSource in watch/[id].astro
5. **"Existing functionality unchanged"** → Same UI behavior, just different underlying transport
6. **"SDK handles reconnection automatically"** → Remove manual reconnection logic from both files

The research shows a clear, well-structured codebase with consistent patterns that will translate cleanly to PocketBase SDK usage.