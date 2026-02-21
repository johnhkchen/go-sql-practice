package routes

import (
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// registerLinksView registers the links view count increment endpoint
func registerLinksView(e *core.ServeEvent) {
	e.Router.POST("/api/links/:id/view", func(ev *core.RequestEvent) error {
		return handleLinksView(ev, e.App)
	})
}

// handleLinksView handles POST /api/links/:id/view requests
func handleLinksView(e *core.RequestEvent, app core.App) error {
	// Extract link ID from URL path
	// URL format: /api/links/:id/view
	linkId := extractPathParam(e.Request.URL.Path, "links")

	if linkId == "" || linkId == "view" {
		return e.JSON(400, map[string]string{
			"error": "Link ID is required",
		})
	}

	// Execute atomic SQL update to increment view_count
	sql := "UPDATE links SET view_count = COALESCE(view_count, 0) + 1 WHERE id = {linkId}"
	result, err := app.DB().NewQuery(sql).Bind(dbx.Params{"linkId": linkId}).Execute()
	if err != nil {
		return e.JSON(500, map[string]string{
			"error": "Failed to update view count",
		})
	}

	// Check if any rows were affected (link exists)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return e.JSON(404, map[string]string{
			"error": "Link not found",
		})
	}

	// Fetch the updated record to return complete link data
	record, err := app.FindRecordById("links", linkId)
	if err != nil {
		return e.JSON(500, map[string]string{
			"error": "Failed to retrieve updated record",
		})
	}

	// Build response with full link data
	response := map[string]interface{}{
		"id":          record.Id,
		"url":         record.GetString("url"),
		"title":       record.GetString("title"),
		"description": record.GetString("description"),
		"view_count":  record.GetInt("view_count"),
		"created":     record.GetDateTime("created").Time().Format(time.RFC3339),
		"updated":     record.GetDateTime("updated").Time().Format(time.RFC3339),
		"tags":        []string{}, // Tags handling can be added later if needed
	}

	return e.JSON(200, response)
}