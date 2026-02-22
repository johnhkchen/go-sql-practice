package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

const (
	progressMin = 0.0
	progressMax = 1.0
)

// CreateSessionResponse is the response for creating a new sync session
type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
	AdminURL  string `json:"admin_url"`
	ViewerURL string `json:"viewer_url"`
}

// UpdateProgressRequest is the request body for updating progress
type UpdateProgressRequest struct {
	Progress float64 `json:"progress"`
}

// registerSyncSessions registers the sync session routes
func registerSyncSessions(e *core.ServeEvent) {
	// POST /api/sync/create - Create new session
	e.Router.POST("/api/sync/create", func(ev *core.RequestEvent) error {
		return handleCreateSession(ev, e.App)
	})

	// POST /api/sync/:id/progress - Update session progress
	e.Router.POST("/api/sync/:id/progress", func(ev *core.RequestEvent) error {
		return handleUpdateProgress(ev, e.App)
	})
}

// handleCreateSession creates a new sync session with an admin token
func handleCreateSession(e *core.RequestEvent, app core.App) error {
	// Generate admin token
	token, err := GenerateToken()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate admin token",
		})
	}

	// Get the sync_sessions collection
	collection, err := app.FindCollectionByNameOrId("sync_sessions")
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to find sync_sessions collection",
		})
	}

	// Create new record
	record := core.NewRecord(collection)
	record.Set("progress", 0.0)
	record.Set("admin_token", token)

	// Save the record
	if err := app.Save(record); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create session",
		})
	}

	// Build response
	response := CreateSessionResponse{
		SessionID: record.Id,
		AdminURL:  fmt.Sprintf("/sync/%s/control?token=%s", record.Id, token),
		ViewerURL: fmt.Sprintf("/sync/%s", record.Id),
	}

	return e.JSON(http.StatusCreated, response)
}

// handleUpdateProgress updates the progress of a sync session
func handleUpdateProgress(e *core.RequestEvent, app core.App) error {
	// Get session ID from URL
	sessionID := extractPathParam(e.Request.URL.Path, "sync")

	// Get admin token from query parameter
	adminToken := e.Request.URL.Query().Get("token")
	if adminToken == "" {
		return e.JSON(http.StatusForbidden, map[string]string{
			"error": "Missing admin token",
		})
	}

	// Parse request body
	var req UpdateProgressRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate progress value
	if err := validateProgress(req.Progress); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Find the session record
	record, err := app.FindRecordById("sync_sessions", sessionID)
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]string{
			"error": "Session not found",
		})
	}

	// Verify admin token (constant-time comparison)
	storedToken := record.GetString("admin_token")
	if !ValidateToken(adminToken, storedToken) {
		return e.JSON(http.StatusForbidden, map[string]string{
			"error": "Invalid admin token",
		})
	}

	// Update progress
	record.Set("progress", req.Progress)

	// Save the updated record
	if err := app.Save(record); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update progress",
		})
	}

	// Return the updated session (admin_token excluded for security)
	return e.JSON(http.StatusOK, map[string]interface{}{
		"id":       record.Id,
		"progress": record.GetFloat("progress"),
		"created":  record.GetDateTime("created").Time(),
		"updated":  record.GetDateTime("updated").Time(),
	})
}


// validateProgress checks if the progress value is within valid range
func validateProgress(progress float64) error {
	if progress < progressMin || progress > progressMax {
		return fmt.Errorf("progress must be between %.1f and %.1f", progressMin, progressMax)
	}
	return nil
}

