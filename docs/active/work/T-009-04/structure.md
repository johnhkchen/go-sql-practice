# T-009-04 Structure: File-Level Implementation Plan

## Architecture Overview

The solution implements a standardized API URL configuration pattern across the frontend codebase, eliminating hardcoded URLs and resolving port conflicts. The approach maintains existing patterns while adding consistency and relative URL fallback support.

## File Modifications

### 1. Configuration Files

#### `frontend/.env` (MODIFY)
**Current State**: `PUBLIC_API_URL=http://127.0.0.1:8093`
**Change**: Update to match server port
```bash
PUBLIC_API_URL=http://127.0.0.1:8090
```
**Rationale**: Align development environment with actual server configuration

#### `frontend/astro.config.mjs` (MODIFY)
**Current State**: Basic static build configuration
**Change**: Add development proxy configuration
```javascript
import { defineConfig } from 'astro/config';
import node from '@astrojs/node';

export default defineConfig({
  output: 'static',
  adapter: node({
    mode: 'standalone'
  }),
  build: {
    assets: 'assets',
    format: 'directory'
  },
  server: {
    proxy: {
      '/api': 'http://127.0.0.1:8090'
    }
  }
});
```
**Rationale**: Enable relative URLs in development environment

### 2. Primary Fix - StatsSummary Component

#### `frontend/src/components/StatsSummary.astro` (MODIFY)
**Current Problem**: Line 141 hardcodes `fetch('http://127.0.0.1:8094/api/stats')`
**Change**: Replace with configurable API base URL

**Implementation**:
1. Add API base URL helper function at top of `<script>` section
2. Replace hardcoded URL with function call
3. Maintain existing error handling and response processing

**Code Structure**:
```typescript
// Add before StatsController class
function getApiBase(): string {
  if (typeof import.meta !== 'undefined' && import.meta.env?.PUBLIC_API_URL) {
    return import.meta.env.PUBLIC_API_URL;
  }
  if (typeof window !== 'undefined' && (window as any).PUBLIC_API_URL) {
    return (window as any).PUBLIC_API_URL;
  }
  return ''; // Relative URL fallback
}

// In fetchStats method, line 141:
const response = await fetch(`${getApiBase()}/api/stats`, {
```

### 3. Standardization - Other Components

#### Components Already Using Correct Pattern (AUDIT ONLY)
These files use `import.meta.env.PUBLIC_API_URL || fallback` and need verification only:
- `frontend/src/components/PresenterController.astro`
- `frontend/src/components/SearchInterface.astro`
- `frontend/src/pages/tags/[slug].astro`
- `frontend/src/pages/index.astro`
- `frontend/src/pages/links/[id].astro`
- `frontend/src/pages/watch/[id].astro`
- `frontend/src/pages/sync/[id].astro`
- `frontend/src/pages/present/[id].astro`
- `frontend/src/pages/present/index.astro`
- `frontend/src/pages/sync/[id]/control.astro`

**Verification Requirements**:
1. Confirm fallback URLs use `localhost:8090` (not other ports)
2. Ensure client-side code uses `window.PUBLIC_API_URL` pattern
3. Validate server-side code uses `import.meta.env.PUBLIC_API_URL` pattern

#### Components Without API Calls (NO CHANGE)
These files don't make API calls and require no changes:
- `frontend/src/layouts/BaseLayout.astro`
- `frontend/src/pages/stats.astro` (includes StatsSummary component)
- `frontend/src/components/LinksList.astro`
- `frontend/src/components/LinkCard.astro`
- `frontend/src/components/Navigation.astro`
- `frontend/src/components/GoLive.astro`

## Module Boundaries and Interfaces

### API URL Resolution Module
**Location**: Inline function in each component (existing pattern)
**Interface**:
```typescript
function getApiBase(): string
```
**Behavior**:
- Returns configured API base URL or empty string for relative URLs
- Handles both server-side (build time) and client-side (runtime) contexts
- Provides fallback mechanism for different deployment scenarios

### Environment Variable Interface
**Build-time**: `import.meta.env.PUBLIC_API_URL`
**Runtime**: `window.PUBLIC_API_URL`
**Fallback**: Empty string (relative URLs)

### Configuration Hierarchy
1. **Build-time environment variable** (`import.meta.env.PUBLIC_API_URL`)
2. **Runtime window variable** (`window.PUBLIC_API_URL`)
3. **Relative URL fallback** (empty string)

## Component Integration Patterns

### Server-Side Rendering (SSR) Components
**Pattern**: Use `import.meta.env.PUBLIC_API_URL` for build-time configuration
**Implementation**: Direct environment variable access in component script sections
**Examples**: Most page-level components

### Client-Side Interactive Components
**Pattern**: Use `window.PUBLIC_API_URL` for runtime configuration
**Implementation**: JavaScript code accessing global window object
**Examples**: PresenterController, SearchInterface

### Hybrid Components (SSR + Client-side)
**Pattern**: Support both build-time and runtime configuration
**Implementation**: Helper function checking both sources
**Examples**: StatsSummary (after fix), complex interactive pages

## Build Process Integration

### Astro Build Integration
**No changes required**: Existing `PUBLIC_API_URL` environment variable processing continues unchanged
**Build output**: Static files with environment variables substituted at build time
**Deployment**: Embedded in Go binary via existing embed.FS mechanism

### Development Workflow
**Change**: Add proxy configuration to Astro dev server
**Benefit**: Enables relative URL testing during development
**Compatibility**: Maintains existing `npm run dev` and `make dev` workflows

## Testing Architecture

### Component-Level Testing
**Target**: Each component's API URL resolution logic
**Approach**: Mock different environment configurations
**Coverage**: Build-time variables, runtime variables, fallback behavior

### Integration Testing
**Target**: End-to-end API communication
**Approach**: Test with development and production configurations
**Coverage**: All API endpoints across different deployment scenarios

### Configuration Testing
**Target**: Environment variable processing
**Approach**: Test with different .env configurations
**Coverage**: Development proxy, production relative URLs, custom deployments

## Deployment Considerations

### Development Deployment
**Configuration**: `PUBLIC_API_URL=http://127.0.0.1:8090` in `.env`
**Proxy**: Astro dev server proxies `/api/*` to Go server
**Workflow**: `npm run dev` for frontend, `make dev` for full stack

### Production Deployment
**Configuration**: No environment variables needed (relative URLs)
**Assets**: Built frontend embedded in Go binary
**Serving**: Same-origin serving from Go server on configured port

### Custom Deployments
**Configuration**: Override `PUBLIC_API_URL` environment variable
**Flexibility**: Support different API servers or cross-origin scenarios
**Compatibility**: Existing patterns continue to work

## Change Ordering

### Phase 1: Configuration Fix
1. Update `frontend/.env` with correct port (8090)
2. Add proxy configuration to `astro.config.mjs`
3. Verify development environment works

### Phase 2: StatsSummary Fix
1. Add `getApiBase()` helper function
2. Replace hardcoded URL in `fetchStats()` method
3. Test component functionality

### Phase 3: Audit and Verification
1. Review all components using API calls
2. Standardize fallback URLs where needed
3. Verify consistent pattern usage

### Phase 4: Testing and Validation
1. Test development and production builds
2. Verify all API endpoints work correctly
3. Confirm no hardcoded URLs remain

## Risk Mitigation Structure

### Rollback Strategy
**Approach**: Changes are localized and reversible
**Critical Path**: StatsSummary.astro is primary change
**Dependencies**: Configuration files can be reverted independently

### Testing Strategy
**Unit Level**: Each component's URL resolution
**Integration Level**: API communication across environments
**System Level**: Full development and production workflows

### Monitoring Points
**Development**: Proxy configuration working
**Build Process**: Environment variable substitution
**Production**: API calls resolving to correct endpoints