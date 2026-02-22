# Research: Extract Inline Scripts to TypeScript

## Current State Overview

The frontend contains six Astro components with substantial inline `<script>` blocks (100-500+ lines each) that define complete classes for handling various interactive features. These scripts are embedded directly in the Astro files, making them difficult to test, lint, and reuse.

## Affected Files and Their Scripts

### 1. SyncViewer Class (`sync/[id].astro`)
- **Location:** Lines 136-291 (~155 lines)
- **Class:** `SyncViewer`
- **Purpose:** Real-time sync viewer for progress tracking
- **Dependencies:**
  - Uses EventSource API for SSE connections
  - Manages connection states (connected/connecting/disconnected)
  - Updates progress bars and status indicators
- **Key Features:**
  - SSE connection management with reconnection logic
  - Real-time progress updates from server
  - DOM manipulation for status indicators

### 2. SyncController Class (`sync/[id]/control.astro`)
- **Location:** Lines 204-550+ (~350 lines)
- **Class:** `SyncController`
- **Purpose:** Admin control interface for sync sessions
- **Dependencies:**
  - Requires admin token for API calls
  - Manages progress slider with keyboard controls
  - Clipboard API for URL copying
- **Key Features:**
  - Throttled progress updates (33ms)
  - Enhanced keyboard navigation
  - Accessibility features (ARIA announcements)

### 3. PresentationAutoViewer Class (`watch/[id].astro`)
- **Location:** Lines 598-950+ (~350 lines)
- **Class:** `PresentationAutoViewer`
- **Purpose:** Auto-following presentation viewer
- **Dependencies:**
  - Uses stepConversion utilities (progressToStep function)
  - EventSource for real-time updates
  - Multiple presentation states (waiting/starting/live/ended)
- **Key Features:**
  - State machine for presentation lifecycle
  - Mid-session join handling
  - Server state synchronization

### 4. PresenterController Class (`PresenterController.astro`)
- **Location:** Lines 804-1100+ (~200 lines)
- **Class:** `PresenterController`
- **Purpose:** Presenter's control panel for presentations
- **Dependencies:**
  - Imports stepConversion utilities from '../utils/stepConversion.js'
  - Uses presentation data structure with step_count, step_labels
  - Admin token for API calls
- **Key Features:**
  - Step navigation (prev/next/jump)
  - Progress slider with keyboard support
  - Stop presentation functionality

### 5. StatsController Class (`StatsSummary.astro`)
- **Location:** Lines 110-380 (~270 lines)
- **Class:** `StatsController`
- **Purpose:** Statistics dashboard controller
- **Dependencies:**
  - Re-exports types from '../types/api' at component level
  - Duplicates StatsData and StatsState interfaces inline (lines 66-87)
  - Uses getApiBase() helper function defined inline
- **Key Features:**
  - Async stats fetching with error handling
  - DOM updates for cards and lists
  - Screen reader announcements
  - Refresh functionality

### 6. Search Enhancement (`SearchInterface.astro`)
- **Location:** Lines 431-700+ (~270 lines)
- **Not a class:** Direct event-driven code
- **Purpose:** Client-side search enhancement
- **Dependencies:**
  - Imports LinkItem, SearchResponse from '../types/api'
  - Falls back to window.location.reload() for non-empty results (line 657)
  - State management object for search state
- **Key Issues:**
  - Incomplete client-side DOM update implementation
  - Falls back to page reload instead of updating DOM
  - Mixed client/server rendering approach

## Existing Infrastructure

### Type Definitions (`frontend/src/types/api.ts`)
Already contains:
- `LinkItem`, `SearchResponse` interfaces
- `StatsData`, `StatsState` interfaces
- `Presentation`, `PresentationStatus` interfaces
- `SearchState` interface

### Utility Functions (`frontend/src/utils/stepConversion.ts`)
Already contains:
- `stepToProgress`, `progressToStep` functions
- `getNavigationState`, `getNextStep`, `getPreviousStep`
- `formatStepDisplay` function
- TypeScript types for navigation state

## Common Patterns Observed

### API Configuration
- Most files use: `import.meta.env.PUBLIC_API_URL || 'http://localhost:8090'`
- StatsSummary has a `getApiBase()` helper that checks both import.meta and window

### DOM Initialization
- All classes check `document.readyState` before initialization
- Use DOMContentLoaded event listener if still loading
- Query and cache DOM elements in setupDOM() methods

### Event Handling
- Enhanced keyboard support (Arrow keys, Home/End, PageUp/PageDown)
- Throttling/debouncing for performance-sensitive operations
- ARIA attributes for accessibility

### Error Handling
- Try-catch blocks around async operations
- Loading states during operations
- User-friendly error messages

## Technical Constraints

### Import/Export in Astro Scripts
- Astro `<script>` tags run on the client
- Can use ES module imports for TypeScript files
- Need to maintain current initialization patterns

### Type Safety
- Current inline scripts duplicate type definitions
- Extraction will allow proper type imports from `types/api.ts`
- Better IDE support and type checking

### Build System
- Frontend uses npm build process
- TypeScript compilation already configured
- Must maintain compatibility with Astro's client-side script handling

## Key Observations

1. **Duplication:** StatsController duplicates types that already exist in `types/api.ts`
2. **Inconsistency:** Some components import utilities (PresenterController) while others inline everything
3. **Testing Gap:** Inline scripts cannot be unit tested
4. **Search Issue:** SearchInterface has incomplete client-side implementation that needs fixing
5. **Instantiation Pattern:** All classes use similar DOM-ready patterns and element querying
6. **API Base:** Inconsistent API base URL resolution across components