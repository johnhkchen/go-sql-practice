# Design: Clean Up Unused apiFetch and Window Globals

## Decision Summary

1. **Use `apiFetch<T>()`** - Refactor pages to use the centralized function
2. **Replace window globals with data attributes** - Pass SSR values via DOM
3. **Remove window controller assignments** - Keep controllers local
4. **Consolidate to `lib/` directory** - Move syncViewer.ts from scripts/

## Option Analysis

### 1. apiFetch Usage Decision

**Option A: Use apiFetch (SELECTED)**
- Eliminates 80+ lines of duplicated fetch boilerplate across 7 pages
- Provides consistent error handling with typed errors
- Already implemented and tested in T-010-03
- Reduces maintenance burden with single source of truth
- Easy migration path - mostly drop-in replacement

**Option B: Delete apiFetch (REJECTED)**
- Would preserve status quo but perpetuate duplication
- Manual patterns work but violate DRY principle
- Error handling inconsistencies likely to grow over time
- Missed opportunity to use already-built abstraction

**Rationale**: The function exists, provides value, and the manual duplication is technical debt. Using it is the clear choice.

### 2. API Base Configuration Strategy

**Option A: Data attributes on DOM (SELECTED)**
```html
<div id="stats-container" data-api-base={API_BASE}>
```
- Clean separation of SSR and client concerns
- No global namespace pollution
- Works reliably in all Astro script contexts
- Standard web platform approach

**Option B: Keep window globals (REJECTED)**
- Current approach works but is messy
- TypeScript requires casting to `any`
- Global namespace pollution
- Not idiomatic for modern web development

**Option C: Environment variable only (REJECTED)**
- import.meta.env not available in client scripts
- Would require build-time injection
- Less flexible than data attributes

**Rationale**: Data attributes are the standard way to pass server-rendered values to client scripts. They're cleaner than globals and more reliable than environment variables in this context.

### 3. Controller Instance Management

**Option A: Remove global assignments (SELECTED)**
```javascript
// Instead of: window.syncController = new SyncController(...)
const syncController = new SyncController(...);
```
- Cleaner scope management
- No global pollution
- Controllers already work without global access
- Prevents accidental cross-component interference

**Option B: Keep globals for debugging (REJECTED)**
- Debugging can be done with browser devtools
- Global access rarely needed in practice
- Creates maintenance confusion

**Rationale**: The global assignments appear to be legacy patterns. The controllers work fine as local variables, and debugging doesn't require global access with modern devtools.

### 4. TypeScript File Organization

**Option A: Consolidate in lib/ (SELECTED)**
- Move `scripts/syncViewer.ts` → `lib/syncViewer.ts`
- All 7 controller/enhancement files in one directory
- Clear convention: lib/ for all extracted business logic
- Single import path pattern

**Option B: Consolidate in scripts/ (REJECTED)**
- Would require moving 6 files instead of 1
- "lib" better describes shared libraries
- More disruptive to existing imports

**Option C: Keep split (REJECTED)**
- Inconsistent and confusing
- No clear reasoning for the split
- Makes maintenance harder

**Rationale**: Moving one file is less disruptive than six. The lib/ directory name better represents the purpose of these shared modules.

### 5. Error Handling Approach

**Approach: Preserve existing error UI patterns**
- apiFetch throws ApiError with same codes pages expect
- Pages catch and map to existing error messages
- No changes to error display components
- Smooth transition with no user-visible changes

**Error mapping:**
```typescript
catch (err) {
  if (err instanceof ApiError) {
    error = err.code; // 'timeout' | 'network' | 'server' | 'notfound'
  } else {
    error = 'network'; // fallback
  }
}
```

### 6. Migration Strategy

**Incremental file-by-file approach:**
1. Each page migrated independently
2. Test each page after migration
3. No big-bang changes
4. Can rollback individual pages if needed

**Data attribute migration pattern:**
```astro
<!-- Component with data attribute -->
<div id="presenter-container" data-session-id={sessionId} data-api-base={API_BASE}>

<!-- Script reads from data attribute -->
<script>
  const container = document.getElementById('presenter-container');
  const sessionId = container?.dataset.sessionId;
  const apiBase = container?.dataset.apiBase || 'http://localhost:8090';
</script>
```

### 7. Import.meta.env Access Pattern

**For TypeScript modules:**
```typescript
// In lib/*.ts files, keep the dual check pattern
function getApiBase() {
  // First try import.meta.env (works in module context)
  if (typeof import.meta !== 'undefined' && import.meta.env?.PUBLIC_API_URL) {
    return import.meta.env.PUBLIC_API_URL;
  }
  // Then try data attribute from container
  const container = document.querySelector('[data-api-base]');
  if (container?.dataset.apiBase) {
    return container.dataset.apiBase;
  }
  // Fallback
  return 'http://localhost:8090';
}
```

### 8. Testing Considerations

**What needs testing:**
- Each migrated page loads without errors
- API calls succeed with apiFetch
- Error states display correctly
- Controllers initialize properly
- No console errors about undefined globals

**Build verification:**
- `npm run build` must pass
- No TypeScript errors
- No runtime errors

## Rejected Alternatives

### Using define:vars Without Window Assignment
- Technically possible but complex
- Would need different patterns per component
- Data attributes are cleaner

### Creating a Global Config Object
- Just moves the problem
- Still pollutes global namespace
- Data attributes are more idiomatic

### Build-Time Environment Injection
- Would require build tool changes
- Less flexible than runtime configuration
- Complicates deployment

## Implementation Priority

1. **High Priority**: Remove window globals for data (security/clarity)
2. **High Priority**: Adopt apiFetch (eliminate duplication)
3. **Medium Priority**: Consolidate file locations (consistency)
4. **Low Priority**: Remove controller globals (nice-to-have)

This design maintains backward compatibility while cleaning up technical debt systematically.