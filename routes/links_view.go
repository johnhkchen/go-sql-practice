package routes

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
)

// registerLinksView registers the links view count increment endpoint
func registerLinksView(e *core.ServeEvent) {
	e.Router.POST("/api/links/:id/view", func(ev *core.RequestEvent) error {
		return handleLinksView(ev, e.App)
	})
}

// handleLinksView handles POST /api/links/:id/view requests
func handleLinksView(e *core.RequestEvent, app core.App) error {
	// Extract link ID from path parameter
	linkId := e.Request.PathParam("id")
	if linkId == "" {
		return e.JSON(400, map[string]string{
			"error": "Link ID is required",
		})
	}

	// Execute atomic SQL update to increment view_count
	sql := "UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?"
	result, err := app.DB().NewQuery(sql).Execute(linkId)
	if err != nil {
		return e.JSON(500, map[string]string{
			"error": "Failed to update view count",
		})
	}

	// Check if any rows were affected (link exists)
	if result.RowsAffected == 0 {
		return e.JSON(404, map[string]string{
			"error": "Link not found",
		})
	}

	// Fetch the updated record using PocketBase DAO
	dao := daos.New(app.DB())
	record, err := dao.FindRecordById("links", linkId)
	if err != nil {
		// This should not happen if the update succeeded
		return e.JSON(500, map[string]string{
			"error": "Failed to fetch updated record",
		})
	}

	// Return the updated record
	return e.JSON(200, record)
}