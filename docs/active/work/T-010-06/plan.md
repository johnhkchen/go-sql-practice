# Plan: Extract Inline Scripts to TypeScript

## Implementation Steps

### Step 1: Create scripts directory
- Create `frontend/src/scripts/` directory
- Verify creation

### Step 2: Extract SyncViewer class
- Create `frontend/src/scripts/syncViewer.ts`
- Copy class from `sync/[id].astro` (lines 134-146)
- Add proper TypeScript types
- Export the class
- Test: Verify TypeScript compilation

### Step 3: Update sync/[id].astro
- Remove class definition (lines 134-146)
- Add import and instantiation script
- Use `define:vars` to pass variables
- Test: Load sync viewer page, verify SSE connection works

### Step 4: Extract SyncController class
- Create `frontend/src/scripts/syncController.ts`
- Copy class from `sync/[id]/control.astro` (lines 204-550+)
- Add TypeScript types and interfaces
- Export the class
- Test: Verify TypeScript compilation

### Step 5: Update sync/[id]/control.astro
- Remove class definition
- Add import and instantiation script
- Test: Load control page, verify slider and updates work

### Step 6: Extract PresentationAutoViewer class
- Create `frontend/src/scripts/presentationViewer.ts`
- Copy class from `watch/[id].astro` (lines 598-950+)
- Import progressToStep from stepConversion
- Add TypeScript types
- Export the class
- Test: Verify TypeScript compilation

### Step 7: Update watch/[id].astro
- Remove class definition
- Add import and instantiation script
- Test: Load watch page, verify state transitions work

### Step 8: Extract PresenterController class
- Create `frontend/src/scripts/presenterController.ts`
- Copy class from `PresenterController.astro` (lines 804-1100+)
- Import stepConversion utilities
- Add TypeScript types
- Export the class
- Test: Verify TypeScript compilation

### Step 9: Update PresenterController.astro
- Remove class definition
- Add import and instantiation script
- Test: Load presenter page, verify navigation works

### Step 10: Extract StatsController class
- Create `frontend/src/scripts/statsController.ts`
- Copy class from `StatsSummary.astro` (lines 110-380)
- Import types from `types/api.ts`
- Remove duplicate type definitions
- Export the class
- Test: Verify TypeScript compilation

### Step 11: Update StatsSummary.astro
- Remove entire script block with class
- Add import and instantiation script
- Add global refreshStats function
- Test: Load stats, verify refresh works

### Step 12: Extract and simplify SearchEnhancer
- Create `frontend/src/scripts/searchEnhancer.ts`
- Rewrite to use server-side navigation only
- Import types from `types/api.ts`
- Export the class
- Test: Verify TypeScript compilation

### Step 13: Update SearchInterface.astro
- Remove entire script block
- Add simplified import and instantiation
- Test: Verify search still triggers server-side navigation

### Step 14: Final build verification
- Run `cd frontend && npm run build`
- Fix any TypeScript errors
- Verify all pages load correctly

### Step 15: Functional testing
- Test sync viewer real-time updates
- Test sync controller progress updates
- Test presentation viewer state transitions
- Test presenter controller navigation
- Test stats refresh functionality
- Test search functionality

## Testing Strategy

### Unit Test Preparation (Future)
Each extracted class is now testable:
- Mock DOM elements
- Mock fetch calls
- Test state transitions
- Test event handling

### Integration Testing
- Manual testing of each component
- Verify SSE connections
- Verify API calls
- Test keyboard navigation
- Test accessibility features

### Build Verification
- TypeScript compilation must succeed
- No console errors in browser
- All functionality preserved

## Rollback Plan

If issues arise:
1. Git stash changes
2. Revert to inline scripts
3. Debug specific component
4. Re-apply extraction incrementally

## Success Criteria

- ✅ All 6 script blocks extracted to `.ts` files
- ✅ Astro files reduced to <10 lines of script
- ✅ TypeScript compilation succeeds
- ✅ All existing functionality preserved
- ✅ No browser console errors
- ✅ Build completes successfully

## Time Estimates

- Steps 1-2: 5 minutes (SyncViewer)
- Steps 3-4: 5 minutes (SyncController)
- Steps 5-6: 5 minutes (PresentationAutoViewer)
- Steps 7-8: 5 minutes (PresenterController)
- Steps 9-10: 5 minutes (StatsController)
- Steps 11-12: 5 minutes (SearchEnhancer)
- Steps 13-15: 10 minutes (Testing & verification)

Total: ~40 minutes

## Notes

- The sync/[id].astro file has already been modified with API_BASE import
- Use the existing import patterns from the modified file
- Maintain all existing functionality exactly
- Focus on extraction, not optimization (except SearchEnhancer which needs fixing)