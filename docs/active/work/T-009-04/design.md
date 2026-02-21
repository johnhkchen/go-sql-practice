# T-009-04 Design: Frontend API URL Configuration

## Problem Statement

The frontend has inconsistent API URL configuration with hardcoded URLs preventing proper deployment flexibility. `StatsSummary.astro` uses a hardcoded `http://127.0.0.1:8094/api/stats` URL that doesn't match server configuration, and there are port mismatches across the development environment.

## Design Options Considered

### Option 1: Relative URLs Only
**Approach**: Convert all API calls to relative URLs (`/api/stats`, `/api/links`, etc.)

**Pros**:
- Simplest solution for same-origin deployment
- No configuration needed in production
- Automatically works with any port/domain
- Eliminates all hardcoded URLs

**Cons**:
- Breaks development workflow with separate Astro dev server
- Requires proxy configuration for `astro dev`
- Less flexible for cross-origin scenarios

**Assessment**: Limited by development constraints. Astro dev server runs on different port than Go server.

### Option 2: Environment Variable Standardization
**Approach**: Fix environment variable configuration and apply existing pattern consistently

**Pros**:
- Builds on established codebase patterns
- Supports both development and production
- Flexible for different deployment scenarios
- Minimal code changes required

**Cons**:
- Still requires environment configuration
- Complex client-side vs server-side variable handling
- Maintenance burden of keeping configs synchronized

**Assessment**: Good fit for existing architecture, but has configuration complexity.

### Option 3: Hybrid Relative + Environment Fallback
**Approach**: Use relative URLs in production, environment variables in development

**Pros**:
- Best of both worlds - simple in production, flexible in development
- Works with existing deployment (same-origin)
- Graceful fallback mechanism
- Supports all use cases

**Cons**:
- More complex logic in components
- Runtime detection needed
- Testing complexity increases

**Assessment**: Most robust but adds complexity.

### Option 4: Build-time Configuration Injection
**Approach**: Inject API base URL at build time based on target environment

**Pros**:
- Single source of truth for API configuration
- No runtime overhead
- Clean separation of concerns
- Eliminates client-side configuration complexity

**Cons**:
- Requires build system changes
- Less flexible for runtime configuration changes
- More complex build pipeline

**Assessment**: Clean but requires significant infrastructure changes.

## Selected Design: Option 2 + Relative URL Optimization

### Rationale

Based on the research findings, the codebase already has a well-established pattern using `import.meta.env.PUBLIC_API_URL` with fallbacks. The best approach is to:

1. **Standardize the existing pattern** across all components
2. **Fix configuration inconsistencies** (port mismatches)
3. **Add relative URL support** where it provides clear benefits
4. **Maintain development workflow** compatibility

This approach leverages existing infrastructure while addressing the specific problems without major architectural changes.

### Core Design Decision

**Primary Strategy**: Environment variable with relative URL fallback
- Server-side: `import.meta.env.PUBLIC_API_URL || '/api'` (relative fallback)
- Client-side: `window.PUBLIC_API_URL || '/api'` (relative fallback)
- Development: Proper proxy configuration to eliminate port conflicts

### Implementation Strategy

#### 1. Fix StatsSummary.astro
Replace hardcoded URL with standard pattern:
```typescript
// Before: fetch('http://127.0.0.1:8094/api/stats')
// After: fetch(`${getApiBase()}/api/stats`)
```

#### 2. Standardize Configuration Pattern
Create reusable API base URL function:
```typescript
function getApiBase(): string {
  // Server-side (build time)
  if (typeof import.meta !== 'undefined' && import.meta.env?.PUBLIC_API_URL) {
    return import.meta.env.PUBLIC_API_URL;
  }
  // Client-side (runtime)
  if (typeof window !== 'undefined' && (window as any).PUBLIC_API_URL) {
    return (window as any).PUBLIC_API_URL;
  }
  // Same-origin fallback (production default)
  return '';
}
```

#### 3. Resolve Port Conflicts
- Update `.env` to match server port: `PUBLIC_API_URL=http://127.0.0.1:8090`
- Add Astro dev server proxy configuration
- Document proper development setup

#### 4. Audit and Update All Components
Apply consistent pattern to all 17 `.astro` files, ensuring:
- Server-side env var usage for SSR/SSG
- Client-side window variable for dynamic content
- Relative URL fallback for same-origin deployment

### Environment Configuration Strategy

#### Development Environment
```bash
# frontend/.env
PUBLIC_API_URL=http://127.0.0.1:8090
```

#### Astro Dev Server Configuration
```javascript
// astro.config.mjs
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://127.0.0.1:8090'
    }
  }
});
```

#### Production Environment
- No environment variables needed (relative URLs work)
- Optional override via environment for different deployments

### Compatibility Considerations

#### Existing Patterns
- Maintain `import.meta.env.PUBLIC_API_URL` pattern
- Keep `window.PUBLIC_API_URL` client-side access
- Preserve fallback mechanisms

#### Astro Framework Integration
- Build-time variable substitution works unchanged
- Client-side hydration patterns remain compatible
- Static build output unchanged

#### PocketBase Integration
- Same-origin serving continues to work
- Embedded frontend patterns unchanged
- Custom route handling unaffected

## Risk Assessment

### Low Risks
- **Configuration errors**: Mitigated by standardized patterns and documentation
- **Development setup**: Proxy configuration solves port conflicts cleanly

### Medium Risks
- **Testing complexity**: Multiple environment configurations need validation
- **Rollback complexity**: Changes affect multiple files but are localized

### Mitigation Strategies
- Comprehensive testing across development and production environments
- Clear documentation of new configuration patterns
- Incremental rollout starting with StatsSummary.astro

## Alternative Approaches Rejected

### Pure Relative URLs
**Why rejected**: Breaks development workflow without significant proxy setup complexity. The environment variable pattern is already established and works well.

### Complete Build-time Injection
**Why rejected**: Requires major build system changes and reduces runtime flexibility. The hybrid approach provides better balance.

### Multiple Configuration Systems
**Why rejected**: Would increase complexity without clear benefits. Standardizing on one proven pattern is more maintainable.

## Success Criteria

### Functional Requirements
1. StatsSummary.astro works with correct API URL
2. All components use consistent API base URL pattern
3. Development and production environments both function correctly
4. No hardcoded host:port combinations remain

### Non-Functional Requirements
1. No performance impact on API calls
2. Development workflow maintains existing convenience
3. Configuration remains simple and clear
4. Future API endpoint additions follow established pattern

### Validation Approach
- Component-level testing with different environment configurations
- Integration testing in both development and production builds
- Verification of all 17 `.astro` files for pattern compliance
- End-to-end testing of StatsSummary functionality