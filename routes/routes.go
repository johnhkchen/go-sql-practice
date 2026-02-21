package routes

import (
	"github.com/pocketbase/pocketbase/core"
)

// Register registers all custom routes with the PocketBase app
func Register(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Register individual routes
		registerHealth(e)
		registerSyncSessions(e)
		registerStats(e)
		registerLinksSearch(e)
		registerLinksSearchSimple(e)

		// Register static file serving (must be last to catch all unmatched routes)
		registerStatic(e)

		// Continue middleware chain
		return e.Next()
	})
}