# Plan: Frontend Shared Types and API Client
## Ticket T-010-03

### Implementation Steps

#### Step 1: Create Type Definitions
**File**: `frontend/src/types/api.ts`
**Action**: Create new file with all interface definitions
**Validation**:
- Run `cd frontend && npm run build` - should succeed
- No runtime impact (types only)
**Commit**: "feat: create shared TypeScript interfaces for frontend"

#### Step 2: Create API Utilities
**File**: `frontend/src/lib/api.ts`
**Action**: Create new file with API_BASE, FETCH_TIMEOUT, ApiError, and apiFetch
**Dependencies**: Import types from Step 1
**Validation**:
- TypeScript compilation should succeed
- Exports should be accessible
**Commit**: "feat: create shared API utilities and fetch helper"

#### Step 3: Fix control.astro Variable Bug
**File**: `frontend/src/pages/sync/[id]/control.astro`
**Action**: Move `let error = null;` declaration before line 19
**Validation**:
- No JavaScript errors in browser console
- Page should load without hoisting issues
**Commit**: "fix: correct variable declaration order in control.astro"

#### Step 4: Refactor index.astro
**File**: `frontend/src/pages/index.astro`
**Actions**:
1. Add imports for types and api utilities
2. Remove inline interface definitions
3. Remove API_BASE and FETCH_TIMEOUT constants
4. Update LinkCard import (remove type import)
**Validation**:
- Homepage loads correctly
- Search functionality works
- Links display properly
**Commit**: "refactor: use shared types and API utilities in index page"

#### Step 5: Refactor tags/[slug].astro
**File**: `frontend/src/pages/tags/[slug].astro`
**Actions**:
1. Add imports
2. Remove duplicate interfaces
3. Remove constants
**Validation**:
- Tag pages load
- Links display correctly
- Error states work
**Commit**: "refactor: use shared types in tag pages"

#### Step 6: Refactor links/[id].astro
**File**: `frontend/src/pages/links/[id].astro`
**Actions**:
1. Add imports
2. Remove interfaces and constants
3. Update field references if needed (created → created_at)
4. Remove duplicate API_BASE in client script
**Validation**:
- Individual link pages load
- View count increment works
- Tag resolution works
**Commit**: "refactor: use shared types in link detail pages"

#### Step 7: Refactor Presentation Pages
**Files**:
- `frontend/src/pages/present/index.astro`
- `frontend/src/pages/present/[id].astro`
**Actions**:
1. Add imports to both files
2. Remove all duplicate interfaces
3. Use FETCH_TIMEOUT_LONG where appropriate
**Validation**:
- Presenter dashboard loads
- Live presentations work
- Control panel functions
**Commit**: "refactor: use shared types in presentation pages"

#### Step 8: Refactor Sync Pages
**Files**:
- `frontend/src/pages/sync/[id].astro`
- `frontend/src/pages/sync/[id]/control.astro` (already partially done)
**Actions**:
1. Add imports
2. Remove constants
3. Complete control.astro refactor
**Validation**:
- Sync viewer works
- Sync control panel works
- Token validation functions
**Commit**: "refactor: use shared types in sync pages"

#### Step 9: Refactor Remaining Pages
**Files**:
- `frontend/src/pages/watch/[id].astro`
- `frontend/src/pages/stats.astro`
**Actions**:
1. Add imports where needed
2. Remove any duplicate constants
**Validation**:
- Watch page loads
- Stats page displays data
**Commit**: "refactor: complete page refactoring for shared types"

#### Step 10: Update SearchInterface Component
**File**: `frontend/src/components/SearchInterface.astro`
**Actions**:
1. Remove interface definitions from script block (lines 433-457)
2. Add type imports in script block
**Validation**:
- Search interface works
- Client-side search functions
- No TypeScript errors in browser
**Commit**: "refactor: use shared types in SearchInterface component"

#### Step 11: Update StatsSummary Component
**File**: `frontend/src/components/StatsSummary.astro`
**Actions**:
1. Remove interface definitions (lines 3-24)
2. Add re-export from types/api
**Validation**:
- Stats component renders
- Stats page still works
- Type exports available
**Commit**: "refactor: re-export shared types from StatsSummary"

#### Step 12: Update LinkCard Component
**File**: `frontend/src/components/LinkCard.astro`
**Actions**:
1. Check if LinkItem is exported
2. Remove export if present
3. Import from types/api if used internally
**Validation**:
- Link cards render correctly
- No broken imports in pages
**Commit**: "refactor: update LinkCard to use shared types"

#### Step 13: Final Validation
**Actions**:
1. Run `cd frontend && npm run build`
2. Check for TypeScript errors
3. Test each page type manually
4. Verify no duplicate definitions remain
**Validation Checklist**:
- [ ] Build succeeds
- [ ] No TypeScript errors
- [ ] Homepage loads and search works
- [ ] Tag pages function
- [ ] Link detail pages work
- [ ] Presentation features work
- [ ] Stats display correctly
- [ ] Sync features function
**Commit**: "chore: validate complete refactoring"

### Testing Strategy

#### Unit Testing
Not applicable (no test framework detected)

#### Manual Testing Checklist
1. **Homepage** (`/`)
   - [ ] Initial load shows links
   - [ ] Search functionality works
   - [ ] Error states display correctly

2. **Tag Pages** (`/tags/[slug]`)
   - [ ] Valid tag shows links
   - [ ] Invalid tag shows error
   - [ ] Link cards display

3. **Link Details** (`/links/[id]`)
   - [ ] Link information displays
   - [ ] View count increments
   - [ ] Tags are clickable

4. **Presentations** (`/present/*`)
   - [ ] Dashboard loads
   - [ ] Can create presentation
   - [ ] Control panel works

5. **Stats** (`/stats`)
   - [ ] Statistics load
   - [ ] Data displays correctly

6. **Sync** (`/sync/*`)
   - [ ] Viewer page works
   - [ ] Control requires token
   - [ ] Progress updates work

### Rollback Plan

If issues occur at any step:
1. Git revert the specific commit
2. Investigate the issue
3. Fix and retry the step

Since each step is atomic and committed separately, rollback is straightforward.

### Success Criteria

1. **No Duplicate Definitions**
   ```bash
   grep -r "interface LinkItem" frontend/src/ | wc -l  # Should be 1
   grep -r "interface SearchResponse" frontend/src/ | wc -l  # Should be 1
   ```

2. **Single Source for Constants**
   ```bash
   grep -r "const API_BASE = " frontend/src/ | wc -l  # Should be 1
   grep -r "const FETCH_TIMEOUT = " frontend/src/ | wc -l  # Should be 1
   ```

3. **Build Success**
   ```bash
   cd frontend && npm run build  # Should succeed with no errors
   ```

4. **Variable Bug Fixed**
   - control.astro should not have variable used before declaration

5. **All Pages Functional**
   - Manual testing confirms all features work

### Risk Mitigation

**Risk**: Type mismatches after refactoring
**Mitigation**: Test each page after refactoring, use optional fields

**Risk**: Build failures
**Mitigation**: Commit after each successful step, easy rollback

**Risk**: Client-side script issues
**Mitigation**: Test progressive enhancement features specifically

**Risk**: Import path problems
**Mitigation**: Use consistent relative paths, test builds frequently

### Time Estimate
- Steps 1-2: 15 minutes (new file creation)
- Step 3: 5 minutes (bug fix)
- Steps 4-12: 45 minutes (refactoring)
- Step 13: 15 minutes (validation)
- Total: ~80 minutes

### Notes
- Prioritize fixing the control.astro bug early (Step 3)
- Keep commits atomic for easy rollback
- Test incrementally to catch issues early
- The apiFetch helper can be adopted gradually (not required in first pass)