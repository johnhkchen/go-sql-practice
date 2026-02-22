# Plan: Clean Up Unused apiFetch and Window Globals

## Implementation Steps

### Step 1: Move syncViewer.ts to lib directory
**Files:** `frontend/src/scripts/syncViewer.ts`
- Move file from scripts/ to lib/
- Update import in sync/[id].astro
- Verify no other imports exist
- **Verify:** File moved successfully, import updated

### Step 2: Update TypeScript modules to use data attributes
**Files:** `statsController.ts`, `searchEnhancer.ts`
- Update getApiBase() in statsController to check data attributes
- Update searchEnhancer constructor to check data attributes
- Remove window.PUBLIC_API_URL references
- **Verify:** No TypeScript errors, modules compile

### Step 3: Update StatsSummary component
**Files:** `components/StatsSummary.astro`
- Add data-api-base attribute to stats-container
- Remove window.refreshStats assignment
- Keep refreshStats function local
- **Verify:** Stats still load, no console errors

### Step 4: Migrate index.astro to apiFetch
**Files:** `pages/index.astro`
- Import apiFetch and ApiError from lib/api
- Replace manual fetch pattern with apiFetch
- Update error handling to use ApiError.code
- **Verify:** Page loads, search works, errors display correctly

### Step 5: Migrate tags/[slug].astro to apiFetch
**Files:** `pages/tags/[slug].astro`
- Import apiFetch and ApiError
- Replace manual fetch with apiFetch
- Map ApiError.code to error variable
- **Verify:** Tag pages load, error states work

### Step 6: Migrate links/[id].astro to apiFetch
**Files:** `pages/links/[id].astro`
- Import apiFetch and ApiError
- Replace manual fetch with apiFetch
- Simplify error handling
- **Verify:** Link detail pages load correctly

### Step 7: Migrate present/index.astro to apiFetch
**Files:** `pages/present/index.astro`
- Import apiFetch and ApiError
- Replace manual fetch with apiFetch
- Update error handling
- **Verify:** Presenter dashboard loads

### Step 8: Migrate present/[id].astro with data attributes
**Files:** `pages/present/[id].astro`
- Import apiFetch and ApiError
- Replace manual fetch with apiFetch
- Remove window global assignments (lines 275-278)
- Add data attributes to presenter container
- Update presenterController.ts if needed
- **Verify:** Presenter control page works, no window globals

### Step 9: Update PresenterController component
**Files:** `components/PresenterController.astro`
- Remove window.PUBLIC_API_URL reference
- Add data-api-base to container
- Update inline script to read from data attribute
- **Verify:** Presenter controller initializes correctly

### Step 10: Migrate sync/[id]/control.astro
**Files:** `pages/sync/[id]/control.astro`
- Import apiFetch and ApiError
- Replace manual fetch with apiFetch
- Remove window.syncController assignment
- Keep controller as local variable
- **Verify:** Sync control page works

### Step 11: Migrate sync/[id].astro
**Files:** `pages/sync/[id].astro`
- Remove window.syncViewer assignment
- Keep viewer as local variable
- Verify import path is updated to lib/syncViewer
- **Verify:** Sync viewer page works

### Step 12: Migrate watch/[id].astro
**Files:** `pages/watch/[id].astro`
- Import apiFetch and ApiError
- Replace manual fetch with apiFetch
- Remove window.presentationAutoViewer assignment
- Keep viewer as local variable
- **Verify:** Watch page auto-advances correctly

### Step 13: Final build verification
**Command:** `cd frontend && npm run build`
- Run full build
- Check for TypeScript errors
- Check for build errors
- **Verify:** Build succeeds with no errors

## Testing Strategy

### Unit Testing Points
Not applicable - no unit tests in codebase currently

### Integration Testing Points

**Manual testing checklist for each page:**
1. Page loads without console errors
2. Data fetches successfully
3. Error states display when API is down
4. Timeout errors work (can test with dev tools throttling)
5. 404 errors display correctly

**Specific interaction tests:**
1. Stats refresh button works (StatsSummary)
2. Search functionality works (index.astro)
3. Tag filtering works (tags/[slug].astro)
4. Presenter controls work (present/[id].astro)
5. Sync session controls work (sync/[id]/control.astro)
6. Auto-advance works (watch/[id].astro)

### Build Verification
```bash
cd frontend
npm run build
# Should complete with no errors
```

### Browser Console Checks
After each step, verify in browser console:
- No undefined variable errors
- No failed imports
- No network errors (except intentional error testing)
- No TypeScript runtime errors

## Rollback Plan

Each step is atomic and can be reverted independently:

1. **File move rollback**: Move syncViewer.ts back to scripts/
2. **Page rollback**: Revert individual page file to previous state
3. **Component rollback**: Revert component file to previous state
4. **Module rollback**: Revert TypeScript module to previous state

Git commits after each major step allow easy rollback if issues arise.

## Risk Assessment

**Low Risk:**
- Moving syncViewer.ts file (simple relocation)
- Removing window global assignments (unused functionality)
- Adding data attributes (backward compatible)

**Medium Risk:**
- Migrating to apiFetch (core functionality change)
- Updating error handling (user-visible behavior)

**Mitigation:**
- Test each page thoroughly after migration
- Keep changes atomic per file
- Verify build after each step

## Completion Criteria

✅ All manual fetch patterns replaced with apiFetch
✅ No window global assignments remain
✅ No window.PUBLIC_API_URL references
✅ All TypeScript files in lib/ directory
✅ Build succeeds with no errors
✅ All pages load and function correctly
✅ No console errors in browser

## Time Estimate

- Steps 1-3: ~10 minutes (simple file operations)
- Steps 4-7: ~20 minutes (SSR page migrations)
- Steps 8-12: ~25 minutes (complex page migrations)
- Step 13: ~5 minutes (build verification)
- Testing: ~15 minutes (manual verification)

**Total estimate: ~75 minutes**