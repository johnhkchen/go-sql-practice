# Structure: Clean Up Unused apiFetch and Window Globals

## File Modifications

### 1. SSR Pages - Adopt apiFetch

**frontend/src/pages/index.astro**
- Replace lines 21-65 manual fetch pattern with apiFetch call
- Import ApiError from lib/api
- Simplify error handling to use ApiError.code

**frontend/src/pages/tags/[slug].astro**
- Replace lines 22-48 manual fetch pattern with apiFetch call
- Import ApiError from lib/api
- Map ApiError.code to existing error types

**frontend/src/pages/links/[id].astro**
- Replace lines 20-46 manual fetch pattern with apiFetch call
- Import ApiError from lib/api
- Use ApiError for error categorization

**frontend/src/pages/present/index.astro**
- Replace lines 17-46 manual fetch pattern with apiFetch call
- Import ApiError from lib/api
- Simplify error handling logic

**frontend/src/pages/present/[id].astro**
- Replace lines 62-94 manual fetch pattern with apiFetch call
- Import ApiError from lib/api
- Lines 275-278: Remove window global assignments
- Add data attributes to presenter container element
- Update PresenterController initialization to read from data attributes

**frontend/src/pages/sync/[id]/control.astro**
- Replace lines 40-72 manual fetch pattern with apiFetch call
- Import ApiError from lib/api
- Line 864: Remove window.syncController assignment
- Keep controller as local variable

**frontend/src/pages/watch/[id].astro**
- Replace lines 36-65 manual fetch pattern with apiFetch call
- Import ApiError from lib/api
- Line 950: Remove window.presentationAutoViewer assignment
- Keep viewer as local variable

**frontend/src/pages/sync/[id].astro**
- Line 137: Remove window.syncViewer assignment
- Keep viewer as local variable
- Update import path from '../scripts/syncViewer' to '../lib/syncViewer'

### 2. TypeScript Modules - Remove Window Dependencies

**frontend/src/lib/statsController.ts**
- Lines 13-14: Remove window.PUBLIC_API_URL check
- Update getApiBase() to read from data-api-base attribute
- Pattern:
  ```typescript
  function getApiBase() {
    if (typeof import.meta !== 'undefined' && import.meta.env?.PUBLIC_API_URL) {
      return import.meta.env.PUBLIC_API_URL;
    }
    const container = document.querySelector('[data-api-base]');
    if (container instanceof HTMLElement && container.dataset.apiBase) {
      return container.dataset.apiBase;
    }
    return 'http://localhost:8090';
  }
  ```

**frontend/src/lib/searchEnhancer.ts**
- Line 40: Remove window.PUBLIC_API_URL reference
- Update constructor to read from data-api-base attribute
- Use same pattern as statsController

**frontend/src/lib/presenterController.ts**
- Update to read session data from data attributes
- Remove dependency on window.presenterSessionId, etc.
- Read from container element's data attributes

### 3. Components - Use Data Attributes

**frontend/src/components/PresenterController.astro**
- Line 782: Remove window.PUBLIC_API_URL reference
- Add data-api-base attribute to container
- Update script to read from data attribute

**frontend/src/components/StatsSummary.astro**
- Line 68: Remove window.refreshStats assignment
- Keep refreshStats as local function
- Add data-api-base attribute to stats-container element

### 4. File Relocation

**Move file:**
- FROM: `frontend/src/scripts/syncViewer.ts`
- TO: `frontend/src/lib/syncViewer.ts`
- No internal changes needed, just relocation

## New Patterns Introduced

### Data Attribute Pattern
```astro
<!-- In Astro component -->
<div
  id="controller-container"
  data-api-base={API_BASE}
  data-session-id={sessionId}
  data-admin-token={adminToken}
>

<!-- In script -->
<script>
  const container = document.getElementById('controller-container');
  const apiBase = container?.dataset.apiBase || 'http://localhost:8090';
  const sessionId = container?.dataset.sessionId;
</script>
```

### apiFetch Usage Pattern
```typescript
// Before: 15+ lines of manual fetch
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT);
try {
  const response = await fetch(url, { signal: controller.signal });
  clearTimeout(timeoutId);
  // ... error checking
} catch (err) {
  clearTimeout(timeoutId);
  // ... error handling
}

// After: 3-5 lines with apiFetch
try {
  const data = await apiFetch<ResponseType>(url);
  // ... use data
} catch (err) {
  if (err instanceof ApiError) {
    error = err.code;
  }
}
```

### Local Controller Pattern
```javascript
// Before
window.syncController = new SyncController(sessionId, token, progress, apiBase);

// After
const syncController = new SyncController(sessionId, token, progress, apiBase);
```

## Module Boundaries

### lib/api.ts (Public Exports)
- `API_BASE`: string constant
- `FETCH_TIMEOUT`: number constant
- `FETCH_TIMEOUT_LONG`: number constant
- `ApiError`: error class
- `apiFetch<T>`: generic fetch function

### lib/* Controllers (Public Classes)
Each controller exports its main class:
- `StatsController`
- `SearchEnhancer`
- `PresenterController`
- `SyncController`
- `SyncViewer`
- `PresentationViewer`

## Deleted Code

### Window Assignments (all removed)
- `window.presenterSessionId`
- `window.presenterAdminToken`
- `window.presenterData`
- `window.PUBLIC_API_URL`
- `window.syncController`
- `window.syncViewer`
- `window.presentationAutoViewer`
- `window.refreshStats`

### Manual Fetch Blocks (all removed)
- ~90 lines total of duplicated fetch/abort/timeout patterns
- Replaced with apiFetch calls

## Testing Points

### Build Success
- `cd frontend && npm run build`
- No TypeScript compilation errors
- No missing imports

### Runtime Verification
- All pages load without console errors
- API calls succeed
- Error states display correctly
- Controllers initialize properly

## Import Updates Required

### Pages importing syncViewer
- `frontend/src/pages/sync/[id].astro`
  - FROM: `import { SyncViewer } from '../../scripts/syncViewer';`
  - TO: `import { SyncViewer } from '../../lib/syncViewer';`

### Components/Pages using API utilities
- Add where needed: `import { apiFetch, ApiError } from '../lib/api';`

## No Changes Required

### These files remain unchanged:
- `frontend/src/types/api.ts` - Type definitions unchanged
- `frontend/src/utils/stepConversion.ts` - Not affected
- `frontend/src/lib/presentationViewer.ts` - Already clean
- `frontend/src/layouts/BaseLayout.astro` - Not affected