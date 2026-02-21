package routes

import (
	"github.com/jchen/go-sql-practice/internal/frontend"
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
		registerLinksSearchSimple(e)
		registerLinksView(e)

		// Register static file serving with availability check
		if frontend.FrontendExists() {
			registerStatic(e)
		} else {
			e.App.Logger().Warn("Frontend assets not found, static serving disabled")
		}

		// Continue middleware chain
		return e.Next()
	})
}