# Research: Custom Health Route Implementation

## Current Codebase State

### Project Structure
```
go-sql-practice/
├── main.go                  # Entry point, creates PocketBase app
├── migrations/
│   └── collections.go       # Defines collections and registers via OnServe
├── go.mod                   # Module: github.com/jchen/go-sql-practice
└── pb_data/                 # PocketBase data files (SQLite DBs)
```

### Core Dependencies
- PocketBase v0.36.5 - Web framework with built-in admin UI
- Go 1.26 runtime

### Entry Point Analysis
`main.go:10-19` creates a PocketBase app instance and starts it:
- `app := pocketbase.New()` - Creates default PocketBase instance
- `migrations.Register(app)` - Registers collection migrations
- `app.Start()` - Starts HTTP server with all PocketBase routes

### Migration Pattern
`migrations/collections.go:10-17` demonstrates the hook registration pattern:
- Uses `app.OnServe()` to attach initialization logic
- `BindFunc` registers a function to run when server starts
- Calls `e.Next()` to continue the chain

## PocketBase Routing Architecture

### Built-in Routes
PocketBase provides several built-in route groups:
- `/_/` - Admin UI dashboard
- `/api/` - REST API endpoints for collections
- `/auth/` - Authentication endpoints

### Hook System
PocketBase uses an event-driven architecture with hooks:
- `OnServe()` - Fired when the HTTP server starts
- `OnRequest()` - Global request interceptor
- `OnModelAfter*()` - Model lifecycle hooks

### Custom Route Registration Methods

#### Method 1: OnServe Hook with Router Access
The `OnServe` event provides access to the underlying echo router through `e.Router`.
This is the recommended approach for custom routes as shown in migrations.

#### Method 2: Echo Router Groups
PocketBase uses Echo v5 internally. Routes can be registered via:
- `e.Router.GET()`, `e.Router.POST()`, etc.
- `e.Router.Group()` for route grouping

#### Method 3: Middleware Chain
Routes can be added with middleware for auth, CORS, etc.

## Route Registration Context

### ServeEvent Structure
The `core.ServeEvent` passed to OnServe contains:
- `App` - The PocketBase app instance
- `Router` - Echo router for registering routes
- `Next()` - Continues the middleware chain

### Authentication Consideration
PocketBase has built-in auth middleware, but health endpoints typically bypass auth.
Routes registered directly on the router don't inherit auth by default.

## File Organization Patterns

### Existing Pattern
The codebase already demonstrates package organization:
- `migrations/` package for database-related code
- Single registration function `Register()` called from main

### Common Go Patterns
- `routes/` or `handlers/` package for HTTP handlers
- `api/` package for API-specific routes
- `internal/` for non-exported packages

## Constraints and Assumptions

### From Acceptance Criteria
- Must use `GET /api/health` path
- Must return JSON `{"status": "ok"}` with 200 status
- Must register via OnServe hook (not router hijacking)
- Must be accessible without authentication

### Technical Constraints
- PocketBase already claims `/api/` namespace for collections
- Custom routes should not conflict with built-in routes
- The Echo router is accessible via `ServeEvent.Router`

### Code Style Observations
- Uses standard Go package structure
- Error handling with wrapped errors (`fmt.Errorf`)
- Comments explain purpose of functions
- Registration pattern: package exports single `Register()` function

## Missing Components

Currently there is no:
- Custom route registration code
- HTTP handler functions
- Routes package or similar organization
- Health check endpoint implementation

## Summary

The codebase has a clean foundation with PocketBase initialized and a migration system using the OnServe hook.
The same hook pattern can be extended for custom route registration. The existing `migrations/` package
demonstrates the registration pattern we need to follow. The Echo router is available through the ServeEvent
and supports standard HTTP method registration.