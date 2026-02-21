# Research: Waiting Room Page (T-008-01)

## Context
This ticket implements a waiting room page at `/watch/:id` that displays when a presentation is not live (when `active_session` is null). The page must fetch presentation data, show a waiting state with animations, and automatically transition to the live view when the presentation goes live.

## Codebase Landscape

### Frontend Architecture
The frontend is an Astro application in static generation mode (`output: 'static'`):
- **Location**: `/frontend/`
- **Framework**: Astro 5.17.3, configured for static site generation
- **Build Output**: Static HTML/CSS/JS files in directory format
- **Current Pages**: Only `index.astro` exists
- **No Dynamic Routes**: Currently no `/watch/[id]` or similar patterns exist

### Existing Components
1. **BaseLayout.astro** (`frontend/src/layouts/BaseLayout.astro`):
   - Provides HTML shell with head, body structure
   - Accepts `title` and `description` props
   - Global CSS variables for colors, spacing, typography
   - Responsive design with mobile-first approach
   - Footer with copyright notice
   - Clean, minimal design aesthetic

2. **Navigation.astro** (`frontend/src/components/Navigation.astro`):
   - Sticky header with brand name "Link Bookmarks"
   - Links to Home (`/`) and Stats (`/stats`)
   - Mobile-responsive hamburger menu
   - Uses CSS transitions for hover effects

### Backend Architecture
PocketBase application with custom routes and migrations:
- **Entry Point**: `main.go` - initializes PocketBase with migrations and routes
- **Custom Routes**: `routes/` directory with health checks and sync session handlers
- **API Base**: PocketBase provides REST API at `/api/collections/`

### Data Collections

1. **presentations** collection (from T-007-01):
   - `name` (text, required): Human-readable presentation name
   - `step_count` (number, required, min 1): Total number of steps
   - `step_labels` (json, optional): Array of labels for each step
   - `active_session` (relation to sync_sessions, optional): Currently live session (null when not presenting)
   - `created_by` (relation to users, optional): Owner
   - API Rules: Anyone can view, authenticated users can create, owners can update/delete

2. **sync_sessions** collection (from T-006-01):
   - `progress` (number, min 0, max 1, default 0): Current progress value
   - `admin_token` (text, required): Token for admin access
   - API Rules: Anyone can view, update requires matching admin_token
   - Supports PocketBase realtime subscriptions

### Key Observations

1. **No Dynamic Routing**: Astro is configured for static generation, but `/watch/[id]` needs dynamic parameters. This requires either:
   - Server-side rendering (SSR) mode
   - Static paths with getStaticPaths
   - Client-side routing
   - Hybrid rendering mode

2. **API Integration Pattern**: No existing Astro pages fetch from PocketBase API yet. Need to establish:
   - API client configuration
   - Fetch patterns (during build vs. client-side)
   - Error handling conventions

3. **Realtime Requirements**: The ticket mentions "When `active_session` is set (already live), skips waiting room". This implies:
   - Initial server-side check of presentation state
   - Potential client-side polling or SSE/WebSocket for live updates
   - PocketBase supports realtime subscriptions natively

4. **Progressive Enhancement**: Requirement states "No JavaScript required for initial waiting room render" which aligns with Astro's philosophy but needs careful consideration for:
   - Initial static/SSR render
   - Optional JavaScript for live updates
   - Graceful degradation

5. **CSS Animation Requirement**: "Subtle CSS animation" needed. Current codebase uses:
   - CSS transitions (Navigation component)
   - CSS custom properties for theming
   - No animation libraries or complex keyframes yet

### Dependencies and Constraints

1. **From T-004-01**: Must use the BaseLayout component
2. **From T-007-01**: Must interact with presentations collection structure
3. **PocketBase API**: Available at localhost:8090/api/collections/presentations/records/:id
4. **Mobile Responsiveness**: Existing patterns use media queries at 767px breakpoint
5. **Design System**: Uses CSS variables for consistent theming

### Technical Gaps to Address

1. **Dynamic Routing Solution**: Need to implement dynamic paths in Astro
2. **API Client**: No fetch utilities or PocketBase SDK integration exists
3. **Error Pages**: No 404 or error handling pages exist
4. **Loading States**: No existing loading or skeleton components
5. **Animation Patterns**: No established keyframe animations yet

### File System Impact

Current structure:
```
frontend/
├── src/
│   ├── layouts/
│   │   └── BaseLayout.astro
│   ├── components/
│   │   └── Navigation.astro
│   └── pages/
│       └── index.astro
```

Will need to add:
- Dynamic route handling for `/watch/[id]`
- Potential API utility functions
- Error boundary or error pages
- Waiting room specific styles/animations

## Summary

The codebase has a clean Astro frontend with established layout and navigation components, and a PocketBase backend with presentations and sync_sessions collections. The main challenge is implementing dynamic routing in a statically configured Astro app while maintaining progressive enhancement. The waiting room must fetch presentation data, show loading states, and handle both the waiting and live scenarios gracefully.