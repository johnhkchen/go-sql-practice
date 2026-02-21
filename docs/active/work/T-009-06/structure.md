# Structure: Install PocketBase JS SDK (T-009-06)

## File Changes Overview

This implementation requires changes to 4 files: 1 new file creation and 3 existing file modifications.

### Files Created
- `frontend/src/lib/pb.ts` - New shared PocketBase client instance

### Files Modified
- `frontend/package.json` - Add PocketBase SDK dependency
- `frontend/src/pages/sync/[id].astro` - Replace EventSource with SDK subscription
- `frontend/src/pages/watch/[id].astro` - Replace EventSource with SDK subscriptions

### Files Deleted
None.

## Detailed File Structure

### 1. New File: `frontend/src/lib/pb.ts`

**Purpose**: Shared PocketBase client instance and type definitions for realtime subscriptions.

**Module boundary**:
- Exports singleton `pb` client instance
- Exports TypeScript interfaces for collections used in subscriptions
- Encapsulates PocketBase configuration (API_BASE URL)

**Public interface**:
```typescript
export const pb: PocketBase;
export interface SyncSession { id: string; progress: number; status: string; [key: string]: any; }
export interface Presentation { id: string; state: string; [key: string]: any; }
```

**Dependencies**:
- Imports `pocketbase` npm package
- Uses `import.meta.env.PUBLIC_API_URL` (same pattern as existing pages)

**Size estimate**: ~30 lines (small utility module)

### 2. Modified File: `frontend/package.json`

**Changes**: Add single dependency to dependencies section.

**Before**:
```json
{
  "dependencies": {
    "@astrojs/node": "^9.5.4",
    "astro": "^5.17.3"
  }
}
```

**After**:
```json
{
  "dependencies": {
    "@astrojs/node": "^9.5.4",
    "astro": "^5.17.3",
    "pocketbase": "^0.21.0"
  }
}
```

**Impact**: Requires `npm install` to install new dependency.

### 3. Modified File: `frontend/src/pages/sync/[id].astro`

**Modification scope**: Client-side script section only (`<script client:load>` block)

**Changes**:
- **Import addition**: Add `import { pb } from '../../lib/pb.ts'`
- **EventSource removal**: Remove EventSource instantiation and connection logic
- **Subscription addition**: Replace with `pb.collection('sync_sessions').subscribe(sessionId, callback)`
- **Reconnection logic removal**: Remove manual reconnection attempts (5 max attempts, exponential backoff)
- **Cleanup change**: Replace `eventSource.close()` with `pb.collection().unsubscribe()`

**Preserved functionality**:
- Same page structure (Astro SSR + client script)
- Same UI update logic (`updateProgress`, `progressToStep`, `formatStepDisplay`)
- Same error display in UI
- Same session ID extraction from URL params

**Code reduction**: ~30-40 lines of manual connection management removed, ~10-15 lines of SDK usage added.

### 4. Modified File: `frontend/src/pages/watch/[id].astro`

**Modification scope**: Client-side script section only (`<script client:load>` block)

**Changes**:
- **Import addition**: Add `import { pb } from '../../lib/pb.ts'`
- **EventSource removal**: Remove EventSource instantiation and connection logic
- **Dual subscriptions**: Replace with:
  - `pb.collection('presentations').subscribe(presentationId, callback)`
  - `pb.collection('sync_sessions').subscribe('*', callback)` with filtering
- **Reconnection logic removal**: Remove manual reconnection attempts (10 max attempts, exponential backoff)
- **Message parsing simplification**: Remove `JSON.parse(event.data)`, use `e.record` directly
- **Cleanup change**: Replace `eventSource.close()` with appropriate `unsubscribe()` calls

**Preserved functionality**:
- Same page structure (Astro SSR + client script)
- Same state management (waiting → starting → live → ended)
- Same UI update logic for presentation state changes
- Same presentation and session ID extraction from URL params

**Code reduction**: ~40-50 lines of manual connection management removed, ~15-20 lines of SDK usage added.

## Architecture Impact

### Module Dependencies
```
frontend/src/pages/sync/[id].astro --> frontend/src/lib/pb.ts --> pocketbase (npm)
frontend/src/pages/watch/[id].astro --> frontend/src/lib/pb.ts --> pocketbase (npm)
```

### Data Flow Changes

**Before (EventSource)**:
```
PocketBase Server (/api/realtime) --> EventSource --> JSON.parse() --> Message filtering --> UI Update
```

**After (SDK)**:
```
PocketBase Server (/api/realtime) --> PocketBase SDK --> Typed callback --> UI Update
```

### Connection Management Changes

**Before**: Each page manages its own EventSource connection with custom reconnection logic.

**After**: Shared PocketBase client manages all connections automatically. Pages only handle business logic.

## Component Boundaries

### 1. PocketBase Client Layer (`pb.ts`)
- **Responsibility**: PocketBase connection configuration and type definitions
- **Interface**: Exports singleton client and collection interfaces
- **Dependencies**: PocketBase SDK
- **Used by**: Pages that need realtime subscriptions

### 2. Page Components (sync/[id], watch/[id])
- **Responsibility**: UI state management and user interactions
- **Interface**: Astro pages with SSR + client-side scripts
- **Dependencies**: PocketBase client layer
- **Scope**: Business logic for sync progress and presentation viewing

### 3. Existing Components (unchanged)
- All other pages and components remain unchanged
- Continue using manual fetch() for API calls (out of scope)
- No impact on server-side routing or data handling

## Implementation Order

### Phase 1: Setup
1. Install PocketBase SDK (`npm install pocketbase`)
2. Create shared client (`frontend/src/lib/pb.ts`)

### Phase 2: Sync Page Migration
3. Modify `frontend/src/pages/sync/[id].astro`
4. Test sync page functionality

### Phase 3: Watch Page Migration
5. Modify `frontend/src/pages/watch/[id].astro`
6. Test presentation viewer functionality

### Phase 4: Cleanup Verification
7. Verify both pages handle disconnection/reconnection properly
8. Remove any unused imports or variables from original EventSource code

**Rationale for order**: Sync page is simpler (single subscription) so it serves as a good test case for the SDK integration before tackling the more complex watch page with dual subscriptions.

## Risk Mitigation

### Bundle Size Impact
- PocketBase SDK adds ~50KB to frontend bundle
- Acceptable increase given the functionality benefits
- No performance regression expected

### Connection Behavior Changes
- SDK may have different connection timing than manual EventSource
- Both pages should handle initial connection delays gracefully
- Existing error handling UI should still work

### TypeScript Compatibility
- PocketBase SDK has good TypeScript support
- Interface definitions in pb.ts provide type safety for subscription data
- Astro's TypeScript integration should handle SDK imports without issues

### Rollback Plan
If issues arise during implementation:
1. Keep original EventSource code in version control
2. PocketBase SDK can be uninstalled without affecting other functionality
3. `pb.ts` file can be deleted if not working
4. Manual fetch() calls remain unchanged as fallback

## Testing Checkpoints

### After Each File Change
1. **pb.ts creation**: Verify module imports correctly in browser dev tools
2. **sync/[id] migration**: Test sync session progress updates still work
3. **watch/[id] migration**: Test presentation state transitions still work
4. **Final integration**: Test both pages with network disconnection/reconnection

### Acceptance Criteria Verification
- Sync viewer receives realtime progress updates (same UX as before)
- Presentation viewer receives state change notifications (same UX as before)
- Automatic reconnection works without manual intervention
- No manual reconnection logic remains in either page
- PocketBase client instance is shared between pages

The structure maintains clear separation of concerns while eliminating the most problematic aspects of the current EventSource implementation.