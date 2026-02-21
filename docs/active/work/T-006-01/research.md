# Research: sync-sessions-collection (T-006-01)

## Context

This ticket requires creating a PocketBase collection for sync sessions. Each session will hold:
- A progress float (0 to 1) to track synchronization progress
- An admin token for secure admin access control
- Built-in realtime subscription support for viewers

The collection needs API rules that allow public viewing but restrict updates to requests with matching admin tokens.

## Current State

### Project Structure

The project is a Go-based PocketBase application with the following structure:
- `main.go`: Entry point that initializes PocketBase, registers migrations and routes
- `migrations/collections.go`: Migration system for creating collections programmatically
- `routes/`: Custom routes directory (referenced in main.go, needs to be created)
- `pb_data/`: PocketBase data directory containing SQLite databases
- `frontend/`: Astro frontend application (separate from backend)

### PocketBase Version

Using PocketBase v0.36.5, which provides:
- Built-in realtime subscription support via Server-Sent Events (SSE)
- Collection-based data modeling
- Programmatic migration system
- API rules and custom route capabilities
- Echo v5 framework is also included in dependencies

### Existing Collections

Two collections already exist (created by T-002-01):

1. **tags**: Base collection with fields:
   - `name` (text, required, 1-100 chars)
   - `slug` (text, required, 1-100 chars, pattern validated)
   - Unique index on slug field

2. **links**: Base collection with fields:
   - `url` (URL, required)
   - `title` (text, required, 1-500 chars)
   - `description` (text, optional, max 2000 chars)
   - `view_count` (number, integer only, min 0)
   - `tags` (relation to tags collection, max 100)
   - `created_by` (relation to users if exists, max 1)

### Migration Pattern

The existing migration system follows this pattern:
1. Register migrations in `main.go` via `migrations.Register(app)`
2. Hook into `OnServe()` event to run migrations when server starts
3. Check if collections exist before creating (idempotent)
4. Use PocketBase's core API to define collection schema
5. Create collections using `txApp.Save(collection)`

Key observations:
- Migrations run automatically on server startup
- Collection existence is checked to prevent duplicates
- Field validation is defined programmatically (min/max, patterns, etc.)
- Relations are handled via collection IDs
- Indexes can be created for unique constraints

### Routes Registration

The main.go file shows a routes.Register(app) call, indicating custom routes are expected. This will be needed for admin token validation on updates. The routes package doesn't exist yet and needs to be created.

## Technical Constraints

### Field Types Available

For the sync_sessions collection, relevant field types:
- `NumberField`: For progress (supports float, min/max validation)
- `TextField`: For admin_token (supports pattern, min/max length)
- No native "token" field type, will use text with appropriate generation

### API Rules System

PocketBase API rules:
- Defined per collection for list/view/create/update/delete operations
- Rules are SQL-like expressions evaluated against request context
- Can reference `@request` object for headers, auth, data
- Custom routes can bypass rules for specialized logic

For admin_token validation:
- Standard API rules can't directly validate against field values in request
- Need custom route to handle admin token verification
- Can still use rules for basic view permissions

### Realtime Subscriptions

PocketBase realtime features:
- Automatically available via `/api/realtime` SSE endpoint
- Clients subscribe to collection changes
- No additional configuration needed for basic functionality
- Changes broadcast to all subscribed clients automatically

## Implementation Considerations

### Token Generation

Admin tokens need to be:
- Generated server-side (not client-provided)
- Sufficiently random/unguessable
- Stored securely in the collection
- Options: UUID v4, random hex string, or crypto-random base64

### Access Control Strategy

Two-tier access model:
1. **Viewer access**: Anyone can read session records via standard API
2. **Admin access**: Updates require matching admin_token

Implementation approach:
- Use permissive view rules in collection settings
- Create custom route for admin updates that validates token
- Keep token validation server-side only

### Migration Integration

New migration should:
- Follow existing pattern in `migrations/collections.go`
- Add to existing `createCollections` function
- Check for sync_sessions existence before creating
- Define all fields with appropriate validation

### Custom Routes Structure

Need to create routes package to handle:
- Admin token validation for updates
- Possibly session creation with token generation
- Structure should match migration pattern for consistency

### Default Values

Progress field needs:
- Default value of 0
- Min constraint of 0
- Max constraint of 1
- Should accept float values

Admin token needs:
- No default (generated on session creation)
- Required field
- Sufficient length for security (suggest 32+ chars)

## Dependencies

- Depends on T-002-01 (completed) which established migration pattern
- Routes package referenced but not yet created (will be part of this work)
- No other blocking dependencies identified
- Frontend integration will be separate concern

## Risks and Unknowns

1. **Token validation in API rules**: PocketBase API rules don't support complex token matching logic, requiring custom route implementation

2. **Routes package**: Referenced in main.go but doesn't exist - needs to be created as part of this work

3. **Progress precision**: NumberField should support sufficient decimal precision for progress tracking (0.0 to 1.0)

4. **Session lifecycle**: No requirements specified for session cleanup/expiration

5. **Concurrent updates**: Multiple admin updates to same session could conflict - may need optimistic locking

6. **Token security**: Tokens will be visible to anyone who can view sessions - this is acceptable per requirements

## File System State

Current files:
```
main.go                     # Entry point, registers migrations and routes
migrations/collections.go   # Contains createCollections function
routes/                     # Referenced but doesn't exist yet (needs creation)
pb_data/                    # Database files (data.db, auxiliary.db)
go.mod                      # Using PocketBase v0.36.5, Echo v5
```

## Next Steps

Design phase should address:
1. Specific token generation strategy (crypto/rand vs UUID)
2. Custom route implementation for admin updates
3. Routes package structure and organization
4. Whether to add session metadata fields (name/identifier)
5. Error handling for invalid tokens
6. API endpoint design for session management