# Design: Install PocketBase JS SDK (T-009-06)

## Problem Statement

The frontend currently uses raw `EventSource` for PocketBase realtime subscriptions and manual `fetch()` for API calls. This creates several maintenance and reliability issues:

1. **Manual connection management**: Each page reimplements EventSource connection logic with inconsistent reconnection strategies (5 vs 10 max attempts)
2. **Code duplication**: Similar PocketBase API fetch patterns and error handling repeated across multiple pages
3. **Type safety gaps**: No centralized typing for PocketBase collections, manual JSON parsing of realtime messages
4. **Reliability concerns**: Custom reconnection logic without the robustness of a tested client library

The PocketBase JS SDK provides a typed client with built-in realtime subscription management, auto-reconnect, and auth helpers that would address these issues.

## Design Approaches Evaluated

### Approach 1: Full SDK Migration with Shared Client Instance

**Description**: Install PocketBase JS SDK, create a shared client instance, and migrate all PocketBase API usage to use the SDK instead of manual fetch() and EventSource.

**Scope**:
- Install `pocketbase` npm package
- Create shared client at `frontend/src/lib/pb.ts`
- Convert EventSource usage in `/sync/[id]` and `/watch/[id]`
- Convert fetch() calls across all pages to use SDK
- Remove manual reconnection logic

**Pros**:
- Eliminates all manual API handling
- Maximum type safety across the application
- Consistent error handling and retry logic
- Single source of truth for PocketBase configuration
- Built-in authentication support for future needs

**Cons**:
- Larger scope than ticket requirements
- Risk of introducing bugs in currently stable API calls
- More changes to test and verify

### Approach 2: EventSource-Only Migration (Ticket Scope)

**Description**: Install PocketBase JS SDK but only migrate EventSource usage to SDK subscriptions. Leave existing fetch() calls unchanged.

**Scope**:
- Install `pocketbase` npm package
- Create shared client at `frontend/src/lib/pb.ts`
- Convert EventSource usage in `/sync/[id]` and `/watch/[id]` to `pb.collection().subscribe()`
- Keep existing fetch() API calls as-is
- Remove only EventSource-related reconnection logic

**Pros**:
- Focused scope matching ticket requirements exactly
- Minimizes risk by keeping stable fetch() logic unchanged
- Easier to test and verify
- Immediate improvement to the most problematic area (realtime subscriptions)

**Cons**:
- Hybrid approach means some API calls still manual
- Doesn't address fetch() code duplication
- Type safety improvements only for subscription data

### Approach 3: Gradual Migration with SDK Wrapper

**Description**: Create a wrapper around both SDK and manual fetch() that provides a consistent interface, then migrate one API endpoint at a time.

**Scope**:
- Install `pocketbase` npm package
- Create abstraction layer in `frontend/src/lib/api.ts`
- Migrate EventSource first, then individual fetch() endpoints over time
- Maintain backward compatibility during transition

**Pros**:
- Allows incremental migration
- Maintains consistency during transition
- Lower risk of breaking existing functionality

**Cons**:
- More complex initial setup
- Adds abstraction layer overhead
- May create confusion about which approach to use

## Decision: Approach 2 - EventSource-Only Migration

### Rationale

**Alignment with ticket requirements**: The ticket specifically mentions refactoring `/sync/[id]` and `/watch/[id]` to use SDK subscriptions, not a complete API migration. The acceptance criteria focus on EventSource replacement:
- "Refactor `/sync/[id]` to use `pb.collection('sync_sessions').subscribe()`"
- "Refactor `/watch/[id]` to use SDK subscriptions"
- "SDK handles reconnection automatically (remove manual reconnect logic)"

**Risk mitigation**: The research shows that manual fetch() calls are currently stable and working. The problematic area is EventSource management with inconsistent reconnection logic. By focusing only on EventSource → SDK subscriptions, we address the core problem without introducing unnecessary risk.

**Clear improvement path**: The most significant benefits come from replacing manual EventSource management:
- Automatic reconnection handling
- Better connection state management
- Simplified code in the most complex pages
- Type safety for realtime message parsing

**Future compatibility**: Using a shared PocketBase client instance sets up the foundation for future API migrations without requiring them now.

## Technical Design

### 1. Package Installation

Add PocketBase JS SDK to frontend dependencies:
```json
{
  "dependencies": {
    "pocketbase": "^0.21.0"  // Compatible with PocketBase v0.36.5
  }
}
```

### 2. Shared Client Instance

Create `frontend/src/lib/pb.ts`:
```typescript
import PocketBase from 'pocketbase';

// Use the same API_BASE pattern as existing code
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';

// Create singleton client instance
export const pb = new PocketBase(API_BASE);

// Type definitions for collections used in realtime subscriptions
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

### 3. EventSource Replacement Strategy

**For `/sync/[id]` (Sync viewer)**:
- Replace EventSource with `pb.collection('sync_sessions').subscribe(id, callback)`
- Remove manual reconnection logic (5 max attempts, exponential backoff)
- Simplify error handling since SDK manages connection state
- Keep same UI update logic (`updateProgress`)

**For `/watch/[id]` (Presentation viewer)**:
- Replace EventSource with both:
  - `pb.collection('presentations').subscribe(presentationId, callback)`
  - `pb.collection('sync_sessions').subscribe('*', callback)` with filtering
- Remove manual reconnection logic (10 max attempts, exponential backoff)
- Keep complex state management (waiting → starting → live → ended)

### 4. Message Handling Migration

Current manual parsing:
```typescript
const data = JSON.parse(event.data);
if (data.collection === 'sync_sessions' && data.action === 'update') {
  // Handle update
}
```

SDK approach:
```typescript
pb.collection('sync_sessions').subscribe(sessionId, (e) => {
  if (e.action === 'update') {
    // e.record is already parsed and typed
    this.updateProgress(e.record.progress);
  }
});
```

### 5. Error Handling Simplification

Remove custom reconnection logic:
```typescript
// Remove: Manual reconnect attempts with counters
// Remove: Exponential backoff calculations
// Remove: Connection state tracking (connected/connecting/disconnected)
// Keep: UI error display for user feedback
```

SDK handles connection management automatically, so pages only need to handle business logic errors.

### 6. Cleanup Management

Replace manual EventSource cleanup:
```typescript
// Before: eventSource.close() on page unload
// After: pb.collection().unsubscribe() on page unload
```

## Implementation Considerations

### Backward Compatibility
- No changes to PocketBase server required
- No changes to API endpoints or data structures
- Same realtime message format (PocketBase standard)
- UI behavior remains identical to users

### Testing Strategy
- Verify both pages still receive realtime updates
- Test reconnection behavior (disconnect network, reconnect)
- Confirm UI state transitions work the same way
- Validate that subscription cleanup prevents memory leaks

### Performance Impact
- SDK adds ~50KB to bundle (reasonable for the benefits)
- Automatic connection management should improve reliability
- No performance regression expected for realtime updates

### Error Scenarios
- Network disconnection: SDK handles reconnection automatically
- Invalid session/presentation IDs: Same error handling as before
- Server restart: SDK auto-reconnects, no manual logic needed

## Future Migration Path

This design enables future migration of fetch() calls to SDK without breaking changes:
1. Keep shared client instance from this ticket
2. Gradually replace individual fetch() calls with `pb.collection().getOne()`, etc.
3. Add more type definitions to `pb.ts` as needed
4. Eventually achieve full SDK usage across the application

The EventSource migration provides immediate benefits while setting up infrastructure for broader API improvements later.