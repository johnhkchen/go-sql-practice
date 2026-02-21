# Plan: Install PocketBase JS SDK (T-009-06)

## Implementation Strategy

This plan follows the structure defined in structure.md, implementing PocketBase SDK integration in 4 ordered steps with independent verification at each stage. The approach prioritizes safety by testing simpler changes first before moving to more complex modifications.

## Step 1: Install PocketBase SDK and Setup Shared Client

### Implementation Tasks
1. **Install PocketBase SDK**
   - Navigate to frontend directory: `cd frontend/`
   - Install package: `npm install pocketbase@^0.21.0`
   - Verify package.json updated correctly

2. **Create shared PocketBase client module**
   - Create file: `frontend/src/lib/pb.ts`
   - Implement singleton PocketBase client instance
   - Add TypeScript interfaces for SyncSession and Presentation
   - Use same API_BASE pattern as existing pages

### Code Implementation
```typescript
// frontend/src/lib/pb.ts
import PocketBase from 'pocketbase';

const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';

export const pb = new PocketBase(API_BASE);

export interface SyncSession {
  id: string;
  progress: number;
  status: string;
  [key: string]: any;
}

export interface Presentation {
  id: string;
  state: string;
  [key: string]: any;
}
```

### Verification Criteria
- [ ] `pocketbase` package appears in frontend/package.json dependencies
- [ ] `frontend/src/lib/pb.ts` file exists and TypeScript compiles without errors
- [ ] Can import `pb` in browser dev tools console without errors
- [ ] API_BASE resolves to correct URL (http://localhost:8090 in dev)

### Commit Point
**Commit message**: "feat: add PocketBase JS SDK and shared client instance"

## Step 2: Migrate Sync Page EventSource to SDK

### Implementation Tasks
1. **Modify sync/[id].astro client script**
   - Add import: `import { pb } from '../../lib/pb.ts'`
   - Replace EventSource instantiation with SDK subscription
   - Remove manual reconnection logic (5 attempts, exponential backoff)
   - Update cleanup from `eventSource.close()` to `pb.collection().unsubscribe()`
   - Keep all existing UI update logic unchanged

### Detailed Changes

**Remove EventSource code**:
```typescript
// Remove: EventSource connection setup
// Remove: Manual reconnection logic with maxAttempts = 5
// Remove: Exponential backoff calculation
// Remove: Connection state tracking
```

**Add SDK subscription**:
```typescript
import { pb } from '../../lib/pb.ts';

// Replace EventSource with SDK subscription
const unsubscribe = pb.collection('sync_sessions').subscribe(this.sessionId, (e) => {
  if (e.action === 'update') {
    this.updateProgress(e.record.progress);
  }
});

// Update cleanup
window.addEventListener('beforeunload', () => unsubscribe());
```

### Verification Criteria
- [ ] Sync page loads without JavaScript errors
- [ ] Progress bar updates when sync progress changes
- [ ] Connection automatically reconnects after network interruption
- [ ] No manual reconnection logic remains in client script
- [ ] Page cleanup properly unsubscribes on navigation away

### Testing Procedure
1. Start a sync session from another page
2. Navigate to `/sync/[session-id]`
3. Verify progress bar updates as sync progresses
4. Disconnect network, reconnect, verify updates resume
5. Navigate away and back to verify no memory leaks

### Commit Point
**Commit message**: "feat: migrate sync page to use PocketBase SDK subscriptions"

## Step 3: Migrate Presentation Viewer EventSource to SDK

### Implementation Tasks
1. **Modify watch/[id].astro client script**
   - Add import: `import { pb } from '../../lib/pb.ts'`
   - Replace EventSource with dual SDK subscriptions
   - Remove manual reconnection logic (10 attempts, exponential backoff)
   - Simplify message parsing (use `e.record` instead of `JSON.parse`)
   - Update cleanup to unsubscribe from both collections

### Detailed Changes

**Remove EventSource code**:
```typescript
// Remove: EventSource connection setup
// Remove: Manual reconnection logic with maxAttempts = 10
// Remove: Exponential backoff calculation
// Remove: JSON.parse(event.data) parsing
// Remove: Connection state management
```

**Add dual SDK subscriptions**:
```typescript
import { pb } from '../../lib/pb.ts';

// Subscribe to presentation state changes
const unsubscribePresentation = pb.collection('presentations').subscribe(this.presentationId, (e) => {
  if (e.action === 'update') {
    this.handlePresentationStateChange(e.record);
  }
});

// Subscribe to sync session updates (filtered)
const unsubscribeSync = pb.collection('sync_sessions').subscribe('*', (e) => {
  if (e.action === 'update' &&
      e.record.presentation_id === this.presentationId) {
    this.handleSyncUpdate(e.record);
  }
});

// Update cleanup
window.addEventListener('beforeunload', () => {
  unsubscribePresentation();
  unsubscribeSync();
});
```

### Verification Criteria
- [ ] Presentation viewer loads without JavaScript errors
- [ ] State transitions work: waiting → starting → live → ended
- [ ] Receives both presentation state and sync progress updates
- [ ] Connection automatically reconnects after network interruption
- [ ] No manual reconnection logic remains in client script
- [ ] Page cleanup properly unsubscribes from both collections

### Testing Procedure
1. Create a presentation and navigate to `/watch/[presentation-id]`
2. Start the presentation from presenter interface
3. Verify state changes from "waiting" to "starting" to "live"
4. Verify sync progress updates appear during presentation
5. Disconnect network, reconnect, verify updates resume
6. End presentation and verify state changes to "ended"

### Commit Point
**Commit message**: "feat: migrate presentation viewer to use PocketBase SDK subscriptions"

## Step 4: Integration Testing and Cleanup

### Implementation Tasks
1. **Cross-page testing**
   - Test both sync and presentation viewers simultaneously
   - Verify no subscription conflicts or memory leaks
   - Test rapid navigation between pages
   - Verify cleanup works properly

2. **Code cleanup**
   - Remove any unused imports from EventSource migration
   - Remove unused variables or functions
   - Verify no EventSource-related code remains
   - Check TypeScript compilation passes

3. **Documentation update**
   - Update any inline comments referencing EventSource
   - Verify code follows existing patterns and conventions

### Verification Criteria
- [ ] Both pages work correctly when used simultaneously
- [ ] No memory leaks when rapidly switching between pages
- [ ] TypeScript compilation passes without warnings
- [ ] No unused imports or dead code remains
- [ ] All manual reconnection logic has been removed
- [ ] Both pages handle network disconnection/reconnection gracefully

### Testing Procedure
1. **Simultaneous usage test**:
   - Open sync viewer in one tab
   - Open presentation viewer in another tab
   - Start a presentation that includes sync operations
   - Verify both pages receive appropriate updates

2. **Memory leak test**:
   - Navigate between sync and presentation viewers rapidly
   - Check browser dev tools for increasing memory usage
   - Verify WebSocket connections are properly closed

3. **Network reliability test**:
   - Use browser dev tools to simulate network conditions
   - Test offline → online transitions
   - Verify automatic reconnection without user intervention

### Commit Point
**Commit message**: "test: verify PocketBase SDK integration and remove EventSource dependencies"

## Testing Strategy

### Unit Testing Approach
While this project doesn't appear to have a formal test suite, verification will be done through:
- Browser dev tools console verification
- Manual functional testing of realtime features
- Network condition simulation for reliability testing

### Integration Testing Focus
- **Cross-browser compatibility**: Test in Chrome and Firefox minimum
- **Network resilience**: Test reconnection behavior with poor connections
- **Performance**: Verify no regression in page load times or memory usage
- **Concurrent usage**: Multiple tabs/pages using realtime features simultaneously

### Acceptance Criteria Verification Plan
Each acceptance criterion maps to specific verification steps:

1. **"pocketbase npm package installed"**
   - ✅ Verify package appears in package.json
   - ✅ Verify npm install succeeds

2. **"Create shared PocketBase client instance"**
   - ✅ Verify pb.ts file exists and exports client
   - ✅ Verify TypeScript compilation passes

3. **"Refactor /sync/[id] to use pb.collection().subscribe()"**
   - ✅ Verify sync page uses SDK subscription
   - ✅ Verify progress updates still work
   - ✅ Verify no EventSource code remains

4. **"Refactor /watch/[id] to use SDK subscriptions"**
   - ✅ Verify presentation viewer uses SDK subscriptions
   - ✅ Verify state transitions still work
   - ✅ Verify dual subscription handling

5. **"Existing functionality unchanged"**
   - ✅ Verify sync progress bar behavior identical
   - ✅ Verify presentation state transitions identical
   - ✅ Verify UI/UX unchanged from user perspective

6. **"SDK handles reconnection automatically"**
   - ✅ Verify manual reconnection code removed
   - ✅ Verify automatic reconnection works
   - ✅ Test network interruption scenarios

## Risk Management

### Rollback Strategy
If any step fails verification:
1. **Step 1 failure**: Remove pb.ts file, uninstall package
2. **Step 2 failure**: Revert sync/[id].astro to EventSource version
3. **Step 3 failure**: Revert watch/[id].astro to EventSource version
4. **Step 4 failure**: Keep working steps, investigate specific issues

### Common Issues and Solutions
- **Import path errors**: Verify TypeScript module resolution
- **Connection timing differences**: SDK may connect faster/slower than EventSource
- **Message format changes**: Verify e.record matches expected data structure
- **Cleanup timing**: Ensure unsubscribe happens before page unload

### Performance Monitoring
- Monitor bundle size increase (~50KB expected)
- Watch for memory leaks during development
- Verify no performance regression in realtime updates
- Test with multiple concurrent connections

## Success Metrics

### Functional Success
- Both pages receive realtime updates reliably
- Automatic reconnection works without user intervention
- No user-visible changes in functionality or performance

### Technical Success
- No EventSource-related code remains in codebase
- TypeScript compilation passes without warnings
- Clean separation of concerns with shared client instance
- Proper subscription lifecycle management

### Quality Metrics
- Zero JavaScript console errors during normal usage
- Graceful handling of network interruptions
- Memory usage stable during extended sessions
- Fast initial connection establishment

The plan ensures a safe, step-by-step migration from EventSource to PocketBase SDK while maintaining all existing functionality and improving connection reliability.