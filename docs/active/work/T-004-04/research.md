# Research: Stats Page Implementation (T-004-04)

## Overview

This ticket requires building a stats page at `/stats` that displays aggregate data from the custom stats endpoint `GET /api/stats`. The implementation must exercise data fetching and presentational components in Astro with client-side loading capabilities.

## Codebase Architecture

### Backend API Layer

**Endpoint**: `/api/stats`
**File**: `routes/stats.go`
**Registration**: Already integrated into the routing system via `routes/routes.go`

The stats endpoint provides a complete JSON response with the following structure:
```go
type StatsResponse struct {
    TotalLinks int64       `json:"total_links"`
    TotalTags  int64       `json:"total_tags"`
    TotalViews int64       `json:"total_views"`
    TopTags    []TagStats  `json:"top_tags"`
    MostViewed []LinkStats `json:"most_viewed"`
}
```

**Data Source**: SQLite database via PocketBase's dbx.Builder
**Query Performance**: Optimized SQL with aggregates, JOINs, and LIMIT clauses
**Error Handling**: Returns HTTP 500 with error messages on database failures
**Empty Database Handling**: Uses COALESCE to return zeros and empty arrays

### Frontend Architecture

**Framework**: Astro 5.17.3 with Node.js adapter
**Configuration**: Static output with standalone mode in `astro.config.mjs`
**Build Tool**: Native Astro build system

**Layout System**:
- `BaseLayout.astro`: Provides HTML shell, navigation integration, global CSS variables
- `Navigation.astro`: Responsive header with brand link and nav items (Home, Stats)
- CSS Architecture: CSS custom properties with responsive design patterns
- Mobile Support: Pure CSS hamburger menu using checkbox technique

**Existing Pages**:
- `/` (index.astro): Simple home page using BaseLayout
- `/watch/[id]` (dynamic route): Exists but not relevant to stats functionality

**Styling Patterns**:
- CSS custom properties for theming (`--color-*`, `--space-*`, `--font-*`)
- Global reset and typography in BaseLayout
- Scoped component styles using Astro's `<style>` blocks
- Mobile-first responsive design with `@media (max-width: 767px)`

### Data Flow Requirements

**Server-Side Rendering**: Current Astro configuration uses static output
**Client-Side Requirements**: Acceptance criteria specifies "data loads client-side (Astro island)"
**Target Architecture**: Mixed approach - SSR initial page, CSR for data fetching

## Dependency Analysis

### T-004-01: Astro Layout and Navigation ✅
**Status**: Complete (phase: done)
**Artifacts**: BaseLayout.astro and Navigation.astro fully implemented
**Relevance**: Direct dependency - stats page will use BaseLayout and navigation includes `/stats` link

**Key Implementation Details**:
- BaseLayout provides consistent HTML structure and styling
- Navigation component already includes "Stats" link pointing to `/stats`
- Responsive design patterns established
- CSS variable system ready for extension

### T-003-03: Stats Endpoint ✅
**Status**: Complete (phase: done)
**Artifacts**: Complete API implementation with testing
**Relevance**: Critical dependency - provides data source for stats page

**Key Implementation Details**:
- Endpoint tested and validated with production-like data
- Response time: ~8ms (well within performance requirements)
- Full error handling and edge case coverage
- Database seeding ensures non-empty responses in development

## Current State Assessment

**Missing Components**:
- No `/stats` page exists in `frontend/src/pages/`
- No client-side data fetching patterns established in codebase
- No presentational components for stats display

**Available Foundation**:
- Backend API ready and tested
- Frontend layout system complete
- Navigation already includes stats link
- CSS design system established

**Technical Constraints**:
- Astro 5.17.3 with static output configuration
- Node.js adapter in standalone mode
- Must support client-side data refresh without full page reload

## Client-Side Data Fetching Patterns

**Astro Islands**: Need to create interactive client-side components
**Available Options**:
1. Vanilla JavaScript with fetch API (no external dependencies)
2. Framework integration (would require adding React/Vue/Svelte)
3. Astro's built-in client directives (`client:load`, `client:visible`)

**Current Dependencies**: Minimal - only Astro and Node adapter
**Philosophy**: Project appears to favor minimal dependencies

## Component Architecture Patterns

**Existing Patterns**:
- Props-based component interfaces (BaseLayout accepts title, description)
- TypeScript interfaces for prop definitions
- Scoped styling with CSS custom properties
- Semantic HTML structure

**Layout Hierarchy**:
```
BaseLayout (HTML shell + global styles)
├─ Navigation (header + responsive menu)
├─ main (content slot)
└─ footer (copyright)
```

## Styling System Analysis

**Design Tokens** (from BaseLayout):
- Colors: `--color-bg`, `--color-text`, `--color-primary`, `--color-border`, `--color-footer`
- Spacing: `--space-xs` through `--space-xl` (0.25rem to 3rem)
- Layout: `--max-width` (1200px), `--header-height` (60px)
- Typography: `--font-body` (system font stack), `--line-height` (1.6)

**Component Patterns**:
- Cards/containers: likely need new patterns for summary cards
- Lists: need ranked list styling for "Top Tags" and "Most Viewed"
- Loading states: need to handle client-side loading UX
- Responsive layout: must work within established mobile-first approach

## Technical Boundaries

**File Structure Requirements**:
- Page component: `frontend/src/pages/stats.astro`
- Client-side components: likely in `frontend/src/components/`
- TypeScript interfaces: follow existing patterns for props

**Integration Points**:
- BaseLayout integration (consistent with existing pages)
- API endpoint consumption (new pattern for this codebase)
- Client-side hydration (new requirement for this codebase)

**Performance Considerations**:
- API response time: ~8ms (excellent)
- Client-side bundle size: minimize JavaScript for data fetching
- Loading UX: need loading states for client-side requests

## Assumptions and Constraints

**Assumptions**:
- Stats page should follow existing visual design patterns
- Client-side refresh means updating data without navigation
- Summary cards and ranked lists need responsive design
- API endpoint is stable and matches documented schema

**Constraints**:
- Must use existing BaseLayout and Navigation components
- Must maintain current CSS custom property system
- Static build output with client-side interactivity islands
- No external UI framework dependencies (inferred from current setup)

**Questions to Resolve in Design Phase**:
- JavaScript framework choice for client-side functionality
- Loading state and error handling UX patterns
- Refresh mechanism (manual button vs automatic vs both)
- Visual hierarchy for stats display (cards, lists, emphasis)