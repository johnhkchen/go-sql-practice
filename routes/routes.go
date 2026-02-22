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
		registerPresentations(e)
		registerStats(e)
		registerLinksSearch(e)
		registerLinksView(e)

		// Register static file serving
		registerStatic(e)

		// Continue middleware chain
		return e.Next()
	})
}