# Design: Extract Inline Scripts to TypeScript

## Design Options

### Option 1: Direct Class Export (Selected ✓)
**Approach:** Export classes from `.ts` files, import and instantiate in minimal `<script>` tags.

**Pros:**
- Clean separation of concerns
- Full TypeScript support with proper imports
- Testable with standard testing frameworks
- IDE features work perfectly (autocomplete, refactoring)
- Minimal changes to initialization patterns

**Cons:**
- Need to handle API configuration passing
- Must ensure DOM-ready patterns are preserved

### Option 2: Module with Factory Functions
**Approach:** Export factory functions that create and initialize instances.

**Pros:**
- Encapsulates initialization logic
- Could handle DOM-ready internally

**Cons:**
- More complex API surface
- Harder to test individual components
- Deviates from current patterns unnecessarily

### Option 3: Web Components
**Approach:** Convert classes to custom elements.

**Pros:**
- Native browser API
- Encapsulation benefits

**Cons:**
- Major architectural change
- Would require significant refactoring
- Overkill for this use case

## Selected Design: Direct Class Export

### Core Principles

1. **Minimal Change:** Keep the same class structure and initialization patterns
2. **Type Safety:** Leverage TypeScript imports instead of duplicating types
3. **Configuration:** Pass API base and other config as constructor parameters
4. **DOM Ready:** Keep DOM-ready checks in the Astro script tags

### Configuration Pattern

```typescript
// In the .ts file
export class SyncViewer {
  constructor(
    private sessionId: string,
    private initialProgress: number,
    private apiBase: string
  ) {
    // ...
  }
}

// In the .astro file
<script>
  import { SyncViewer } from '../../scripts/syncViewer';

  const apiBase = import.meta.env.PUBLIC_API_URL || '';
  const sessionId = document.body.dataset.sessionId;
  const initialProgress = parseFloat(document.body.dataset.progress || '0');

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      new SyncViewer(sessionId, initialProgress, apiBase);
    });
  } else {
    new SyncViewer(sessionId, initialProgress, apiBase);
  }
</script>
```

### Import Resolution

1. **Type Imports:** Use `import type` for interfaces to avoid runtime overhead
2. **Utility Imports:** Direct imports from existing modules
3. **Path Convention:** Use relative paths from scripts directory

### SearchInterface Special Case

The search enhancement code needs restructuring:

**Current Issue:** Incomplete DOM update logic falls back to page reload
**Solution:** Simplify to always use server-side navigation

```typescript
export class SearchEnhancer {
  constructor(private apiBase: string) {
    this.init();
  }

  private async performSearch(query: string): Promise<void> {
    // For empty query, navigate to home
    if (!query.trim()) {
      window.location.href = '/';
      return;
    }

    // For all searches, use server-side rendering
    const url = new URL(window.location.href);
    url.searchParams.set('q', query);
    window.location.href = url.toString();
  }
}
```

This removes the half-implemented client-side rendering and ensures consistent behavior.

### API Configuration Standardization

Create a shared helper for API base resolution:

```typescript
// In each script file that needs it
function getApiBase(): string {
  if (typeof import.meta !== 'undefined' && import.meta.env?.PUBLIC_API_URL) {
    return import.meta.env.PUBLIC_API_URL;
  }
  if (typeof window !== 'undefined' && (window as any).PUBLIC_API_URL) {
    return (window as any).PUBLIC_API_URL;
  }
  return '';
}
```

### File Naming Convention

- `syncViewer.ts` - Matches class name in camelCase
- `syncController.ts` - Consistency across all extractions
- `presentationViewer.ts` - Clear purpose indication
- `presenterController.ts` - Distinguishes from viewer
- `statsController.ts` - Domain-specific naming
- `searchEnhancer.ts` - Describes enhancement role

### Type Import Strategy

```typescript
// Good - use type imports for interfaces
import type { StatsData, StatsState } from '../types/api';
import type { LinkItem, SearchResponse } from '../types/api';

// Good - regular imports for utilities
import { progressToStep, stepToProgress } from '../utils/stepConversion';
```

### Error Handling Preservation

Maintain existing error handling patterns:
- Try-catch blocks for async operations
- Loading state management
- User-friendly error messages
- Console logging for debugging

### Testing Considerations

Extraction enables future testing:
- Unit tests for business logic
- Mock DOM elements for initialization tests
- Mock fetch for API interaction tests
- State management verification

## Rejected Approaches

### Why Not Factory Functions?
- Current code already uses `new ClassName()` pattern
- Changing to factories adds complexity without benefit
- Direct instantiation is clearer and more testable

### Why Not Combine Similar Classes?
- SyncViewer and SyncController serve different roles
- Keeping them separate maintains single responsibility
- Easier to understand and maintain

### Why Not Fix SearchInterface Client-Side Rendering?
- Would require implementing complex DOM diffing/templating
- Server-side rendering already works well
- Simpler to maintain one rendering path
- Avoids client/server hydration issues

## Migration Safety

1. **No Breaking Changes:** External API remains the same
2. **Progressive Enhancement:** Scripts still enhance server-rendered HTML
3. **Fallback Behavior:** If scripts fail, server-side still works
4. **Build Verification:** TypeScript compilation catches type errors early