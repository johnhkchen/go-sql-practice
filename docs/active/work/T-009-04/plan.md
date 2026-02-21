# T-009-04 Plan: Implementation Steps

## Overview

This plan implements the standardized API URL configuration across the frontend, fixing hardcoded URLs and port inconsistencies. The work is broken into discrete, testable steps that can be committed atomically.

## Implementation Steps

### Step 1: Fix Environment Configuration
**Objective**: Resolve port conflicts in development environment
**Files**: `frontend/.env`, `frontend/astro.config.mjs`
**Estimated Time**: 10 minutes

**Actions**:
1. Update `frontend/.env` to use correct server port:
   - Change `PUBLIC_API_URL=http://127.0.0.1:8093` to `PUBLIC_API_URL=http://127.0.0.1:8090`
2. Add proxy configuration to `frontend/astro.config.mjs`:
   - Add server.proxy configuration for `/api` routes to `http://127.0.0.1:8090`

**Verification**:
- `npm run dev` starts without errors
- Frontend can access API endpoints through proxy
- Environment variable is correctly loaded in Astro

**Commit Message**: "fix: align frontend dev environment with server port 8090"

### Step 2: Fix StatsSummary.astro Hardcoded URL
**Objective**: Replace hardcoded URL with configurable API base
**Files**: `frontend/src/components/StatsSummary.astro`
**Estimated Time**: 15 minutes

**Actions**:
1. Add `getApiBase()` helper function before `StatsController` class:
   ```typescript
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
2. Replace hardcoded URL in `fetchStats()` method (line 141):
   - Change `'http://127.0.0.1:8094/api/stats'` to `${getApiBase()}/api/stats`

**Verification**:
- Component loads without JavaScript errors
- API call uses correct URL (visible in browser dev tools)
- Statistics data loads successfully
- Fallback behavior works in different environments

**Commit Message**: "fix: replace hardcoded API URL in StatsSummary with configurable base"

### Step 3: Audit Components with API Calls
**Objective**: Verify consistent API URL patterns across all components
**Files**: All `.astro` files using API calls
**Estimated Time**: 20 minutes

**Actions**:
1. Review each component identified in research:
   - `PresenterController.astro`
   - `SearchInterface.astro`
   - `pages/tags/[slug].astro`
   - `pages/index.astro`
   - `pages/links/[id].astro`
   - `pages/watch/[id].astro`
   - `pages/sync/[id].astro`
   - `pages/present/[id].astro`
   - `pages/present/index.astro`
   - `pages/sync/[id]/control.astro`

2. For each component, verify:
   - Server-side uses `import.meta.env.PUBLIC_API_URL || fallback`
   - Client-side uses `(window as any).PUBLIC_API_URL || fallback`
   - Fallback URLs use `http://localhost:8090` (not other ports)

3. Fix any inconsistencies found

**Verification**:
- All components use consistent patterns
- No hardcoded host:port combinations remain
- Fallback URLs are standardized

**Commit Message**: "refactor: standardize API URL patterns across frontend components"

### Step 4: Test Development Environment
**Objective**: Verify all functionality works in development mode
**Files**: None (testing only)
**Estimated Time**: 15 minutes

**Actions**:
1. Start development environment:
   - Run `cd frontend && npm run dev`
   - Run `make dev` in separate terminal
2. Test each major component:
   - Navigate to stats page, verify StatsSummary loads
   - Test search functionality
   - Test link viewing and navigation
   - Test presentation features
3. Verify API calls in browser developer tools:
   - Confirm URLs resolve correctly
   - Check for CORS or proxy issues
   - Validate response data

**Verification**:
- All components load without errors
- API calls use correct URLs
- No console errors related to API communication
- Development workflow is smooth

**Commit Message**: None (testing step)

### Step 5: Test Production Build
**Objective**: Verify functionality in production deployment
**Files**: None (testing only)
**Estimated Time**: 15 minutes

**Actions**:
1. Build and run production version:
   - Run `make build` to create production binary
   - Run `make dev` to start production server
2. Test all functionality:
   - Verify StatsSummary component works
   - Test API-dependent components
   - Check browser developer tools for correct URLs
3. Verify relative URL fallback:
   - Confirm API calls use relative URLs (no host:port)
   - Check that same-origin deployment works correctly

**Verification**:
- Production build completes successfully
- All components work with embedded frontend
- API calls use relative URLs as expected
- No hardcoded URLs in network tab

**Commit Message**: None (testing step)

### Step 6: Documentation Update
**Objective**: Document configuration patterns for future development
**Files**: Project documentation (if needed)
**Estimated Time**: 10 minutes

**Actions**:
1. Update development setup documentation if it exists
2. Document environment variable configuration
3. Add notes about API URL patterns for future components

**Verification**:
- Documentation is clear and accurate
- Setup instructions work for new developers

**Commit Message**: "docs: update API URL configuration patterns"

## Testing Strategy

### Unit Testing
**Scope**: Individual component API URL resolution
**Approach**:
- Mock different environment configurations
- Test `getApiBase()` function behavior
- Verify fallback mechanisms

**Test Cases**:
- Environment variable present (build time)
- Window variable present (runtime)
- No configuration (relative URL fallback)
- Invalid configuration handling

### Integration Testing
**Scope**: End-to-end API communication
**Approach**:
- Test with real API endpoints
- Verify data loading and error handling
- Check different deployment configurations

**Test Cases**:
- Development environment with proxy
- Production environment with relative URLs
- Custom environment with different API base
- Network error handling

### Configuration Testing
**Scope**: Environment variable processing
**Approach**:
- Test different .env configurations
- Verify Astro build-time substitution
- Check runtime variable access

**Test Cases**:
- Default development configuration
- Production build with no env vars
- Custom API base URL override
- Missing or malformed configuration

## Rollback Plan

### Quick Rollback (< 5 minutes)
**Scope**: Revert configuration changes only
**Actions**:
1. Revert `frontend/.env` to original port (8093)
2. Remove proxy configuration from `astro.config.mjs`
3. Revert StatsSummary.astro to hardcoded URL

### Full Rollback (< 15 minutes)
**Scope**: Revert all changes if major issues found
**Actions**:
1. Use git to revert all commits made during implementation
2. Verify original functionality is restored
3. Document issues encountered for future resolution

## Risk Mitigation

### High-Priority Risks
1. **StatsSummary breaks**: Primary component failure
   - **Mitigation**: Test thoroughly before committing
   - **Detection**: Component fails to load statistics
   - **Response**: Immediate rollback of StatsSummary changes

2. **Development workflow breaks**: Proxy configuration issues
   - **Mitigation**: Test development environment after config changes
   - **Detection**: API calls fail in development mode
   - **Response**: Fix proxy configuration or revert

### Medium-Priority Risks
1. **Other components affected**: Unintended side effects
   - **Mitigation**: Comprehensive audit in Step 3
   - **Detection**: Components fail after environment changes
   - **Response**: Review and fix affected components

2. **Production deployment issues**: Build or runtime problems
   - **Mitigation**: Test production build before finalizing
   - **Detection**: Build fails or components don't work in production
   - **Response**: Debug and fix or rollback

### Low-Priority Risks
1. **Performance impact**: Additional function calls
   - **Mitigation**: Function is simple and fast
   - **Detection**: Slower API response times
   - **Response**: Profile and optimize if needed

## Success Criteria Verification

### Functional Requirements
- [ ] StatsSummary.astro loads statistics correctly
- [ ] All API calls use configurable URLs (no hardcoded host:port)
- [ ] Development environment works with proxy configuration
- [ ] Production environment works with relative URLs
- [ ] Custom API base URL can be configured via environment variable

### Non-Functional Requirements
- [ ] No performance degradation in API calls
- [ ] Development workflow remains convenient
- [ ] Configuration is simple and well-documented
- [ ] Code follows established patterns and conventions

### Quality Assurance
- [ ] No console errors in browser developer tools
- [ ] All API endpoints tested in both environments
- [ ] Code review passes (if applicable)
- [ ] Documentation is updated and accurate

## Deployment Notes

### Development Setup
- Developers need to restart Astro dev server after config changes
- New proxy configuration may require clearing browser cache
- Environment variables are loaded at Astro startup

### Production Deployment
- No special deployment steps required
- Existing build and embed process works unchanged
- Environment variables optional (relative URLs work automatically)

### Troubleshooting
- Check browser developer tools Network tab for actual URLs used
- Verify environment variable loading with `console.log()` if needed
- Test API endpoints directly to isolate configuration vs server issues