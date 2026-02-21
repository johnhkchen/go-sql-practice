# Structure: Waiting Room Page (T-008-01)

## File Operations Overview

### Files to Modify
1. `frontend/astro.config.mjs` - Switch to hybrid rendering
2. `frontend/package.json` - Add Node adapter for SSR

### Files to Create
1. `frontend/src/pages/watch/[id].astro` - Main waiting room page
2. `frontend/src/styles/animations.css` - Reusable animation styles

### Files Unchanged
- `frontend/src/layouts/BaseLayout.astro` - Will be imported and used
- `frontend/src/components/Navigation.astro` - No changes needed
- All backend files - No backend changes required

## Detailed File Structures

### 1. `frontend/astro.config.mjs` (Modify)
```javascript
import { defineConfig } from 'astro/config';
import node from '@astrojs/node';  // New import

export default defineConfig({
  output: 'hybrid',  // Changed from 'static'
  adapter: node({    // New adapter config
    mode: 'standalone'
  }),
  build: {
    assets: 'assets',
    format: 'directory'
  }
});
```

**Changes**:
- Output mode from 'static' to 'hybrid'
- Add Node.js adapter for SSR capability
- Keeps existing build configuration

### 2. `frontend/package.json` (Modify)
```json
{
  "dependencies": {
    "astro": "^5.17.3",
    "@astrojs/node": "^8.3.4"  // New dependency
  }
  // ... rest unchanged
}
```

**Changes**:
- Add @astrojs/node adapter for SSR support

### 3. `frontend/src/pages/watch/[id].astro` (Create)

**Structure**:
```
---
// Frontmatter (Server-side code)
- Import statements
- SSR configuration
- Parameter extraction
- API fetch logic
- State determination
- Error handling
---

<!-- HTML Template -->
- BaseLayout wrapper
- Conditional rendering:
  - Waiting room view
  - Live view (future)
  - Error view
- CSS animations
- Optional client-side script
```

**Key Sections**:

```astro
---
// Disable prerendering for SSR
export const prerender = false;

import BaseLayout from '../../layouts/BaseLayout.astro';

// Get presentation ID from URL
const { id } = Astro.params;

// Fetch presentation data
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
let presentation = null;
let error = null;

try {
  const response = await fetch(`${API_BASE}/api/collections/presentations/records/${id}`);
  if (response.ok) {
    presentation = await response.json();
  } else if (response.status === 404) {
    error = 'notfound';
  } else {
    error = 'server';
  }
} catch (e) {
  error = 'network';
}

// Determine view state
const isLive = presentation?.active_session !== null;
---

<BaseLayout title={presentation?.name || 'Waiting Room'}>
  <!-- Main content -->
</BaseLayout>
```

### 4. `frontend/src/styles/animations.css` (Create)

**Purpose**: Centralized animation definitions for reuse

```css
/* Pulse animation for waiting indicator */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(0.98);
  }
}

/* Respect motion preferences */
@media (prefers-reduced-motion: reduce) {
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.8; }
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
```

## Component Architecture

### Page Component Hierarchy
```
[id].astro
├── BaseLayout (imported)
│   ├── Navigation (automatic via BaseLayout)
│   └── Footer (automatic via BaseLayout)
├── Waiting Room View (conditional)
│   ├── Presentation Title
│   ├── Status Message
│   └── Animated Indicator
├── Error View (conditional)
│   ├── Error Title
│   └── Error Message
└── Client Script (optional enhancement)
```

### State Machine
```
INITIAL → FETCH_PRESENTATION
         ↓
    [Success] → CHECK_SESSION
    [404]     → RENDER_404
    [Error]   → RENDER_ERROR
         ↓
    [active_session = null] → RENDER_WAITING
    [active_session exists] → RENDER_LIVE
```

## CSS Architecture

### Scoped Styles (in [id].astro)
```css
<style>
  .waiting-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 50vh;
    text-align: center;
    padding: var(--space-lg);
  }

  .waiting-title {
    font-size: 2rem;
    margin-bottom: var(--space-md);
    color: var(--color-primary);
  }

  .waiting-message {
    font-size: 1.125rem;
    color: var(--color-text);
    opacity: 0.8;
    margin-bottom: var(--space-lg);
  }

  .waiting-indicator {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--color-primary), var(--color-border));
    opacity: 0.2;
  }

  /* Mobile responsive */
  @media (max-width: 767px) {
    .waiting-title {
      font-size: 1.5rem;
    }
    .waiting-message {
      font-size: 1rem;
    }
  }
</style>
```

## API Integration Pattern

### Fetch Pattern
```typescript
// Server-side fetch configuration
const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
const FETCH_TIMEOUT = 5000;

// Fetch with timeout
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT);

try {
  const response = await fetch(url, {
    signal: controller.signal,
    headers: {
      'Accept': 'application/json'
    }
  });
  clearTimeout(timeoutId);
  // ... handle response
} catch (error) {
  // ... handle error
}
```

## Progressive Enhancement Structure

### Base Layer (No JavaScript)
- Full waiting room rendered server-side
- CSS animations run automatically
- Page is fully functional

### Enhancement Layer (With JavaScript)
```javascript
<script>
  // Only enhance if JavaScript is available
  if (typeof window !== 'undefined') {
    // 1. Set up periodic status checks
    const checkInterval = setInterval(async () => {
      // Fetch presentation status
      // If active_session changes, update view
    }, 5000);

    // 2. Clean up on page unload
    window.addEventListener('beforeunload', () => {
      clearInterval(checkInterval);
    });
  }
</script>
```

## Error Handling Structure

### Error States
1. **404 Not Found**: Presentation ID doesn't exist
2. **Network Error**: Cannot reach API
3. **Server Error**: API returned 5xx
4. **Invalid Response**: API returned unexpected format

### Error Messages
```javascript
const errorMessages = {
  notfound: {
    title: 'Presentation Not Found',
    message: 'This presentation does not exist or has been removed.'
  },
  network: {
    title: 'Connection Error',
    message: 'Unable to connect to the server. Please check your connection and try again.'
  },
  server: {
    title: 'Server Error',
    message: 'Something went wrong on our end. Please try again later.'
  }
};
```

## Build Configuration

### Environment Variables
```bash
# .env (local development)
PUBLIC_API_URL=http://localhost:8090

# Production would use
PUBLIC_API_URL=https://api.production.com
```

### Package Scripts (no changes needed)
- `npm run dev` - Runs in SSR mode automatically with hybrid
- `npm run build` - Builds both static and SSR pages
- `npm run preview` - Preview with Node adapter

## Module Boundaries

### Public Interface (what other pages can use)
- None - this is a standalone page

### Internal Organization
- Frontmatter: All server-side logic
- Template: All markup
- Style block: All scoped styles
- Script block: Progressive enhancements only

### Dependencies
- BaseLayout from T-004-01
- Presentations API structure from T-007-01
- No new shared components created

## Summary

This structure implements the waiting room as a server-rendered Astro page using hybrid rendering. The page fetches presentation data server-side, renders appropriate content based on session state, includes CSS animations for visual feedback, and supports optional JavaScript enhancements. The implementation maintains progressive enhancement principles while providing a polished user experience.