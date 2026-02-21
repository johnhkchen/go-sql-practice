package routes

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

// registerHealth registers middleware to intercept the built-in health endpoint
func registerHealth(e *core.ServeEvent) {
	// Use middleware to intercept and override the built-in /api/health response
	e.Router.Bind(&hook.Handler[*core.RequestEvent]{
		Func: func(ev *core.RequestEvent) error {
			// Check if this is the health endpoint
			if ev.Request.URL.Path == "/api/health" && ev.Request.Method == "GET" {
				// Return our custom response instead of the built-in one
				return ev.JSON(200, map[string]string{
					"status": "ok",
				})
			}
			// Continue to next middleware/handler for other routes
			return ev.Next()
		},
		Priority: -1000, // Run before other handlers
	})
}