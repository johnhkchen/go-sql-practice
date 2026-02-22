# Research: Clean Up Unused apiFetch and Window Globals

## Current State Analysis

### 1. The `apiFetch<T>()` Function

**Location**: `frontend/src/lib/api.ts` (lines 22-83)

**Status**: Created in T-010-03 but currently unused

The function provides:
- Generic type support for response typing
- Automatic timeout handling via AbortController
- Error categorization (timeout/network/server/notfound)
- Custom ApiError class with error codes
- Default headers for JSON accept
- Configurable timeout (default 5000ms)

### 2. Manual Fetch Patterns in SSR Pages

All SSR pages currently use manual fetch + AbortController + setTimeout pattern:

**Pages with manual fetch (lines of boilerplate each):**
1. `index.astro` (lines 21-65): ~15 lines of fetch boilerplate
2. `tags/[slug].astro` (lines 22-48): ~12 lines of fetch boilerplate
3. `links/[id].astro` (lines 20-46): ~12 lines of fetch boilerplate
4. `present/index.astro` (lines 17-46): ~13 lines of fetch boilerplate
5. `present/[id].astro` (lines 62-94): ~14 lines of fetch boilerplate
6. `sync/[id]/control.astro` (lines 40-72): ~13 lines of fetch boilerplate
7. `watch/[id].astro` (lines 36-65): ~12 lines of fetch boilerplate

**Common pattern observed:**
```typescript
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT);
const response = await fetch(url, {
  signal: controller.signal,
  headers: { 'Accept': 'application/json' }
});
clearTimeout(timeoutId);
// Error handling for ok/404/server/timeout/network
```

### 3. Window Global Usage for API_BASE

**Files reading from `(window as any).PUBLIC_API_URL`:**

1. `frontend/src/lib/statsController.ts` (lines 13-14)
   - Falls back to window after checking import.meta.env
   - Used in getApiBase() helper function

2. `frontend/src/lib/searchEnhancer.ts` (line 40)
   - Sets `this.API_BASE` from window global in constructor
   - No import.meta.env check

3. `frontend/src/components/PresenterController.astro` (line 782)
   - Inline script reads from window global
   - Falls back to localhost default

### 4. Window Global Assignments for Data Passing

**Global assignments found:**

1. **`present/[id].astro` (lines 275-278)**
   - Sets: `window.presenterSessionId`, `window.presenterAdminToken`
   - Sets: `window.presenterData`, `window.PUBLIC_API_URL`
   - Uses `define:vars` directive but still assigns to window

2. **`sync/[id]/control.astro` (line 864)**
   - Sets: `window.syncController = new SyncController(...)`
   - Creates global controller instance

3. **`sync/[id].astro` (line 137)**
   - Sets: `window.syncViewer = new SyncViewer(...)`
   - Creates global viewer instance

4. **`watch/[id].astro` (line 950)**
   - Sets: `window.presentationAutoViewer = new PresentationAutoViewer(...)`
   - Creates global auto-viewer instance

5. **`components/StatsSummary.astro` (line 68)**
   - Sets: `(window as any).refreshStats = () => { ... }`
   - Exposes refresh function globally

### 5. TypeScript File Organization

**Current distribution:**
- `frontend/src/lib/`: 6 files
  - api.ts
  - presenterController.ts
  - presentationViewer.ts
  - searchEnhancer.ts
  - statsController.ts
  - syncController.ts

- `frontend/src/scripts/`: 1 file (the outlier)
  - syncViewer.ts

- `frontend/src/types/`: 1 file
  - api.ts (type definitions)

- `frontend/src/utils/`: 1 file
  - stepConversion.ts

### 6. API Configuration Export

`frontend/src/lib/api.ts` exports:
- `API_BASE`: Uses import.meta.env.PUBLIC_API_URL with fallback
- `FETCH_TIMEOUT`: 5000ms default
- `FETCH_TIMEOUT_LONG`: 10000ms for complex ops
- `ApiError`: Custom error class
- `apiFetch<T>()`: Generic fetch wrapper

### 7. Astro-Specific Constraints

**Client-side script limitations:**
- Inline `<script>` tags in Astro don't have direct access to import.meta.env at runtime
- The `define:vars` directive is available for passing SSR values to client scripts
- Data attributes on DOM elements are an alternative to window globals

**Current workarounds:**
- Pages use `define:vars` but still assign to window object
- TypeScript modules check both import.meta.env and window fallbacks
- Some components hard-code localhost fallback

### 8. Error Handling Patterns

**Manual fetch pages have consistent error categories:**
- `notfound`: 404 responses
- `server`: Non-404 error responses
- `timeout`: AbortError from timeout
- `network`: Other fetch errors

**ApiError class matches these categories:**
- Provides typed error codes
- Includes optional HTTP status
- Extends native Error class

### 9. Dependencies and Build Process

**Build command**: `cd frontend && npm run build`
- Must succeed after all changes
- TypeScript compilation included
- Astro SSR/SSG build process

### 10. Related Ticket Dependencies

**T-010-03**: Created the `apiFetch<T>()` function
- Function exists but unused
- Part of API utilities centralization

**T-010-06**: Extracted scripts to .ts files
- Created the current lib/scripts split
- Some files still use window globals post-extraction