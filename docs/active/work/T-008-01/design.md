# Design: Waiting Room Page (T-008-01)

## Requirements Recap
- Create `/watch/[id]` page that fetches presentation data
- Show waiting room when `active_session` is null
- Show live view when `active_session` is set
- Include CSS animations for visual feedback
- Ensure mobile responsiveness
- Support progressive enhancement (no JS required for initial render)

## Design Options

### Option 1: Server-Side Rendering (SSR) with Astro
**Approach**: Switch Astro to SSR mode, render waiting room server-side
- **Pros**:
  - True progressive enhancement, HTML delivered complete
  - Direct API calls from server, no CORS issues
  - SEO-friendly, accessible
  - Natural handling of dynamic routes
- **Cons**:
  - Requires changing entire app from static to SSR
  - Needs Node.js runtime in production
  - More complex deployment
  - Impacts other static pages unnecessarily

### Option 2: Hybrid Rendering (SSR for Watch, Static for Others)
**Approach**: Use Astro's hybrid mode - SSR for `/watch/*`, static for rest
- **Pros**:
  - Best of both worlds - dynamic where needed, static elsewhere
  - Progressive enhancement maintained
  - Minimal impact on existing pages
- **Cons**:
  - Slightly more complex configuration
  - Mixed deployment requirements
  - Need to manage both static and dynamic routes

### Option 3: Client-Side Rendering with Static Shell
**Approach**: Create static `/watch/[id].astro` page with client-side fetching
- **Pros**:
  - Keeps static deployment simple
  - Can use PocketBase realtime subscriptions easily
  - Familiar SPA-like pattern
- **Cons**:
  - Violates "no JS required" requirement
  - Flash of loading state
  - Poor SEO and accessibility
  - Against Astro philosophy

### Option 4: Static Generation with Redirect Pattern
**Approach**: Static `/watch/index.astro` that reads ID from URL params client-side
- **Pros**:
  - Maintains fully static build
  - Single page handles all IDs
  - Simple deployment
- **Cons**:
  - Requires JavaScript for basic functionality
  - No server-side validation
  - Poor user experience without JS

## Decision: Option 2 - Hybrid Rendering

**Rationale**: Hybrid rendering best satisfies all requirements while maintaining architectural coherence.

1. **Progressive Enhancement**: Server-renders the waiting room HTML, ensuring it works without JavaScript
2. **Minimal Disruption**: Other pages remain static, no unnecessary complexity added
3. **Future-Proof**: When the presentation goes live, we can enhance with client-side updates via PocketBase subscriptions
4. **Performance**: Initial page load is fast (SSR), subsequent updates can be handled client-side
5. **Correct Semantics**: Dynamic content (presentations) gets dynamic handling, static content stays static

## Implementation Design

### Routing Strategy
```
frontend/src/pages/
├── index.astro (static)
├── stats.astro (static - future)
└── watch/
    └── [id].astro (SSR, prerender=false)
```

### Data Flow
1. **Initial Request**: Browser requests `/watch/abc123`
2. **Server Fetch**: Astro server fetches from `http://localhost:8090/api/collections/presentations/records/abc123`
3. **State Check**:
   - If `active_session` is null → render waiting room
   - If `active_session` exists → render live view (or redirect)
4. **Progressive Enhancement**: Optional JavaScript adds:
   - Polling for session status changes
   - Smooth transitions when going live
   - Enhanced animations

### API Integration
```typescript
// Pseudo-code for server-side fetch
const response = await fetch(`${API_BASE}/api/collections/presentations/records/${id}`);
if (!response.ok) {
  return render404();
}
const presentation = await response.json();

if (presentation.active_session) {
  // Render live view or redirect
} else {
  // Render waiting room
}
```

### CSS Animation Design
Three animation options for the waiting room:

1. **Pulse Effect** (Chosen):
   ```css
   @keyframes pulse {
     0%, 100% { opacity: 1; }
     50% { opacity: 0.5; }
   }
   ```
   - Simple, universal, low motion
   - Communicates "active waiting"
   - Accessible, respects prefers-reduced-motion

2. **Loading Dots**:
   - More complex, requires multiple elements
   - Common pattern, might feel generic

3. **Gradient Shimmer**:
   - Smooth but potentially distracting
   - Higher GPU usage on mobile

### Error Handling
- **404**: Presentation ID not found → Custom error message
- **500**: API unreachable → Generic error with retry suggestion
- **Network**: Connection issues → Graceful degradation message

### Mobile Considerations
- Use existing breakpoint (767px) for consistency
- Stack content vertically on small screens
- Larger touch targets for any interactive elements
- Ensure text remains readable (min 16px)

## Rejected Approaches

### Why Not Full SSR?
Unnecessary overhead for static pages like the homepage. Hybrid gives us targeted dynamism.

### Why Not Client-Side Only?
Violates the "no JavaScript required" constraint and goes against Astro's progressive enhancement philosophy.

### Why Not Static with Build-Time Fetching?
Presentations are dynamic - their state changes. Build-time fetching would show stale data.

### Why Not WebSockets Immediately?
Progressive enhancement means starting with a working baseline (SSR HTML), then adding realtime updates as an enhancement.

## Configuration Changes Required

1. Update `astro.config.mjs`:
   ```javascript
   export default defineConfig({
     output: 'hybrid', // Changed from 'static'
     // ... rest of config
   });
   ```

2. Mark `/watch/[id]` as SSR:
   ```astro
   ---
   export const prerender = false; // This page needs SSR
   ---
   ```

## Next Steps
The Structure phase will detail exact file modifications, component boundaries, and integration points. The chosen hybrid approach balances requirements with architectural simplicity.