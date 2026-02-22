# Research: Frontend HTML/CSS/A11y Cleanup (T-010-04)

## Overview

The frontend Astro application has accumulated various code quality issues: HTML semantics problems, dead CSS, accessibility gaps, hardcoded values, and ~46 production console.log statements. The task involves cleaning these up without breaking functionality.

## Current Frontend Architecture

### Build Configuration
- **astro.config.mjs**: Currently set to `output: 'static'` but uses Node.js adapter
- 8 pages explicitly set `export const prerender = false` for SSR
- Contradiction: static output mode incompatible with SSR pages requiring hybrid/server mode

### Layout Structure
- **BaseLayout.astro**: Main layout wrapper at lines 100-102, wraps slot with `<main>` element
- Provides CSS custom properties but missing semantic tokens (--color-success, --color-error)
- Footer has hardcoded copyright year "2024" at line 105

### Component Issues Identified

#### Nested `<main>` Elements
Four pages create invalid HTML by adding their own `<main>` inside BaseLayout's `<main>`:
1. **index.astro** line 97: `<main class="home-page-container">`
2. **tags/[slug].astro** line 79: `<main class="tag-page-container">`
3. **present/[id].astro** line 201: `<main class="presenter-control-page">`
4. **PresenterController.astro** line 67: `<main class="controller-main">`

#### Invalid ARIA Roles
- **LinksList.astro** line 98: `<div role="main">` - incorrect use of main role on div element

#### Broken Tag Navigation
- **links/[id].astro** lines 100, 363: Links use `/search?tag=${tag}` instead of `/tags/${tag}`
- Inconsistent with actual tag route structure

#### Navigation Accessibility
**Navigation.astro** missing:
- No `aria-current="page"` for active navigation link
- Uses checkbox hack for mobile menu (lines 9-14) instead of proper button with ARIA
- No `aria-label` on nav element

#### Dead CSS
- **tags/[slug].astro** lines 329-341: `.loading-skeleton` class and `@keyframes loading` animation never used
- **sync/[id]/control.astro** lines 1426-1448: `.temp-debug` styles marked as temporary, never used

## Console.log Analysis

46 console.log statements found across 9 files:

### By File:
1. **index.astro**: 2 logs (lines 316, 321)
2. **tags/[slug].astro**: 2 logs (lines 373, 378)
3. **SearchInterface.astro**: Multiple logs for search operations
4. **StatsSummary.astro**: Logs for stats loading
5. **PresenterController.astro**: Controller initialization logs
6. **links/[id].astro**: Link view tracking logs
7. **sync/[id].astro**: Sync session logs
8. **sync/[id]/control.astro**: Control panel logs
9. **watch/[id].astro**: Watch mode logs

These are informational/debug logs, not error handling.

## API Configuration

### Existing Infrastructure
- **lib/api.ts**: Provides `API_BASE`, `FETCH_TIMEOUT`, `apiFetch` utility
- Most pages use `API_BASE` correctly
- **GoLive.astro** line 38: Uses relative `/api/` instead of `API_BASE`
- **PresenterController.astro** line 34: Duplicates API_BASE instead of importing

## Error State Patterns

Multiple pages duplicate identical error state HTML/CSS (~50 lines markup, ~80 lines CSS):
- **index.astro**: Lines 108-130 (markup), 162-224 (styles)
- **tags/[slug].astro**: Lines 81-90 (markup), similar styles
- **links/[id].astro**: Error container pattern
- **present/index.astro**: Error states
- **present/[id].astro**: Error handling

Common pattern includes:
- Error icon, title, message
- Back link and retry button
- Identical styling across all instances

## CSS Token System

**BaseLayout.astro** defines base tokens:
- Color tokens: --color-bg, --color-text, --color-primary, --color-border, --color-footer
- Missing semantic tokens: --color-success, --color-error, --color-secondary
- Missing: --border-radius token (hardcoded as 4px/8px throughout)

## Data Passing Patterns

**PresenterController.astro**:
- Takes props for sessionId, adminToken, presentationData
- May use window globals for client-side script communication
- Could benefit from data-* attributes or define:vars

## Dependency Check

The ticket depends on T-010-03 which should be complete based on git status showing modified files.

## File Impact Summary

### Files to Modify:
1. **astro.config.mjs** - Change output mode
2. **BaseLayout.astro** - Add CSS tokens, fix copyright
3. **Navigation.astro** - Add accessibility attributes
4. **index.astro** - Remove nested main, console.logs
5. **tags/[slug].astro** - Remove nested main, dead CSS, console.logs
6. **links/[id].astro** - Fix tag URLs, remove console.logs
7. **present/[id].astro** - Remove nested main
8. **present/index.astro** - Review for API URLs
9. **sync/[id]/control.astro** - Remove temp-debug CSS
10. **PresenterController.astro** - Remove nested main, use lib/api
11. **GoLive.astro** - Use API_BASE from lib/api
12. **SearchInterface.astro** - Remove console.logs
13. **StatsSummary.astro** - Remove console.logs
14. **sync/[id].astro** - Remove console.logs
15. **watch/[id].astro** - Remove console.logs
16. **LinksList.astro** - Remove role="main"

### New File to Create:
1. **components/ErrorState.astro** - Shared error state component

## Constraints

- Must maintain all existing functionality
- Changes are cleanup/quality improvements only
- No breaking changes to component props or page behavior
- Build must succeed: `cd frontend && npm run build`