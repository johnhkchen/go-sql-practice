package routes

import (
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)


// StartLiveResponse is the response for starting a live presentation
type StartLiveResponse struct {
	SessionID  string   `json:"session_id"`
	AdminURL   string   `json:"admin_url"`
	ViewerURL  string   `json:"viewer_url"`
	StepCount  int      `json:"step_count"`
	StepLabels []string `json:"step_labels"`
}

// StopLiveResponse is the response for stopping a live presentation
type StopLiveResponse struct {
	Message string `json:"message"`
}

// StatusResponse is the enhanced presentation status with computed fields
type StatusResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	StepCount   int      `json:"step_count"`
	StepLabels  []string `json:"step_labels"`
	IsLive      bool     `json:"is_live"`
	Progress    *float64 `json:"progress,omitempty"`
	CurrentStep *int     `json:"current_step,omitempty"`
}

// registerPresentations registers presentation lifecycle routes with the PocketBase router.
// These routes handle the live session workflow that PocketBase's built-in CRUD API cannot manage.
func registerPresentations(e *core.ServeEvent) {
	// GET /api/presentations/:id/status - Get presentation status with live info
	e.Router.GET("/api/presentations/:id/status", func(ev *core.RequestEvent) error {
		return handleGetStatus(ev, e.App)
	})

	// POST /api/presentations/:id/stop - Stop live session
	e.Router.POST("/api/presentations/:id/stop", func(ev *core.RequestEvent) error {
		return handleStopLive(ev, e.App)
	})

	// POST /api/presentations/:id/live - Start live session
	e.Router.POST("/api/presentations/:id/live", func(ev *core.RequestEvent) error {
		return handleStartLive(ev, e.App)
	})
}

// getAuthenticatedUser extracts the authenticated user from the request context
func getAuthenticatedUser(e *core.RequestEvent) (*core.Record, error) {
	authRecord := e.Auth
	if authRecord == nil {
		return nil, fmt.Errorf("authentication required")
	}
	return authRecord, nil
}

// checkPresentationOwnership verifies the user owns the presentation
func checkPresentationOwnership(presentation, user *core.Record) error {
	createdBy := presentation.GetString("created_by")
	if createdBy == "" {
		// Presentation has no owner, allow access (for backwards compatibility)
		return nil
	}
	if createdBy != user.Id {
		return fmt.Errorf("not authorized")
	}
	return nil
}


// progressToStep converts a progress value (0-1) to a step index using the formula:
// step_index = round(progress * (step_count - 1)) for step_count > 1, else 0
func progressToStep(progress float64, stepCount int) int {
	if stepCount <= 1 {
		return 0
	}
	// Map progress to step boundaries and round to nearest
	stepProgress := progress * float64(stepCount-1)
	return int(stepProgress + 0.5)
}


// buildStartLiveResponse constructs the response for starting a live session
func buildStartLiveResponse(session, presentation *core.Record, token string) StartLiveResponse {
	stepCount := presentation.GetInt("step_count")
	stepLabels := make([]string, 0)

	// Extract step_labels if they exist
	if labelsInterface := presentation.Get("step_labels"); labelsInterface != nil {
		if labelsSlice, ok := labelsInterface.([]interface{}); ok {
			for _, label := range labelsSlice {
				if labelStr, ok := label.(string); ok {
					stepLabels = append(stepLabels, labelStr)
				}
			}
		}
	}

	return StartLiveResponse{
		SessionID:  session.Id,
		AdminURL:   fmt.Sprintf("/sync/%s/control?token=%s", session.Id, token),
		ViewerURL:  fmt.Sprintf("/watch/%s", presentation.Id),
		StepCount:  stepCount,
		StepLabels: stepLabels,
	}
}

// buildStatusResponse constructs the enhanced status response
func buildStatusResponse(presentation, session *core.Record) StatusResponse {
	stepCount := presentation.GetInt("step_count")
	stepLabels := make([]string, 0)

	// Extract step_labels if they exist
	if labelsInterface := presentation.Get("step_labels"); labelsInterface != nil {
		if labelsSlice, ok := labelsInterface.([]interface{}); ok {
			for _, label := range labelsSlice {
				if labelStr, ok := label.(string); ok {
					stepLabels = append(stepLabels, labelStr)
				}
			}
		}
	}

	response := StatusResponse{
		ID:         presentation.Id,
		Name:       presentation.GetString("name"),
		StepCount:  stepCount,
		StepLabels: stepLabels,
		IsLive:     session != nil,
	}

	// Add live session details if active
	if session != nil {
		progress := session.GetFloat("progress")
		currentStep := progressToStep(progress, stepCount)
		response.Progress = &progress
		response.CurrentStep = &currentStep
	}

	return response
}

func handleGetStatus(e *core.RequestEvent, app core.App) error {
	// Get presentation ID from URL
	presentationID := extractPathParam(e.Request.URL.Path, "presentations")
	if presentationID == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing presentation ID",
		})
	}

	// Find the presentation record
	presentation, err := app.FindRecordById("presentations", presentationID)
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]string{
			"error": "Presentation not found",
		})
	}

	// Check if presentation has an active session
	var session *core.Record
	activeSessionID := presentation.GetString("active_session")
	if activeSessionID != "" {
		session, err = app.FindRecordById("sync_sessions", activeSessionID)
		if err != nil {
			// If session not found, treat as not live (orphaned reference)
			session = nil
		}
	}

	// Build enhanced status response
	response := buildStatusResponse(presentation, session)

	return e.JSON(http.StatusOK, response)
}

func handleStopLive(e *core.RequestEvent, app core.App) error {
	// Get presentation ID from URL
	presentationID := extractPathParam(e.Request.URL.Path, "presentations")
	if presentationID == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing presentation ID",
		})
	}

	// Authenticate user
	user, err := getAuthenticatedUser(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Authentication required",
		})
	}

	// Find the presentation record
	presentation, err := app.FindRecordById("presentations", presentationID)
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]string{
			"error": "Presentation not found",
		})
	}

	// Check ownership
	if err := checkPresentationOwnership(presentation, user); err != nil {
		return e.JSON(http.StatusForbidden, map[string]string{
			"error": "Not authorized to control this presentation",
		})
	}

	// Check if presentation is currently live
	activeSessionID := presentation.GetString("active_session")
	if activeSessionID == "" {
		return e.JSON(http.StatusConflict, map[string]string{
			"error": "Presentation is not currently live",
		})
	}

	// Clear the active session
	presentation.Set("active_session", "")

	// Save the updated presentation
	if err := app.Save(presentation); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to stop presentation",
		})
	}

	// Return success response
	response := StopLiveResponse{
		Message: "Presentation stopped successfully",
	}

	return e.JSON(http.StatusOK, response)
}

func handleStartLive(e *core.RequestEvent, app core.App) error {
	// Get presentation ID from URL
	presentationID := extractPathParam(e.Request.URL.Path, "presentations")
	if presentationID == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing presentation ID",
		})
	}

	// Authenticate user
	user, err := getAuthenticatedUser(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Authentication required",
		})
	}

	// Find the presentation record
	presentation, err := app.FindRecordById("presentations", presentationID)
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]string{
			"error": "Presentation not found",
		})
	}

	// Check ownership
	if err := checkPresentationOwnership(presentation, user); err != nil {
		return e.JSON(http.StatusForbidden, map[string]string{
			"error": "Not authorized to control this presentation",
		})
	}

	// Check if presentation is already live
	activeSessionID := presentation.GetString("active_session")
	if activeSessionID != "" {
		return e.JSON(http.StatusConflict, map[string]string{
			"error": "Presentation is already live",
		})
	}

	// Generate admin token
	adminToken, err := GenerateToken()
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

	// Create new sync session record
	session := core.NewRecord(collection)
	session.Set("progress", 0.0)
	session.Set("admin_token", adminToken)

	// Save the session record
	if err := app.Save(session); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create sync session",
		})
	}

	// Update presentation with active session
	presentation.Set("active_session", session.Id)

	// Save the updated presentation
	if err := app.Save(presentation); err != nil {
		// If presentation update fails, we have an orphaned session
		// For now, accept this as acceptable risk per design decision
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to link session to presentation",
		})
	}

	// Build response with URLs and metadata
	response := buildStartLiveResponse(session, presentation, adminToken)

	return e.JSON(http.StatusCreated, response)
}