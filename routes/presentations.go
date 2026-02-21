package routes

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

const (
	TokenLength = 32  // bytes before encoding, matches sync_sessions.go
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

// registerPresentations registers presentation lifecycle routes
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

// generateAdminToken generates a secure random token for session administration
func generateAdminToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// progressToStep converts a progress value (0-1) to a step index
func progressToStep(progress float64, stepCount int) int {
	if stepCount <= 1 {
		return 0
	}
	// Map progress to step boundaries and round to nearest
	stepProgress := progress * float64(stepCount-1)
	return int(stepProgress + 0.5)
}

// stepToProgress converts a step index to progress value (0-1)
func stepToProgress(stepIndex int, stepCount int) float64 {
	if stepCount <= 1 {
		return 0.0
	}
	return float64(stepIndex) / float64(stepCount-1)
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

// validateToken performs constant-time token comparison
func validateToken(provided, stored string) bool {
	if len(provided) != len(stored) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(provided), []byte(stored)) == 1
}

// Placeholder handler functions - to be implemented in subsequent steps

func handleGetStatus(e *core.RequestEvent, app core.App) error {
	// TODO: Implement in Step 2
	return e.JSON(http.StatusNotImplemented, map[string]string{
		"error": "Status endpoint not yet implemented",
	})
}

func handleStopLive(e *core.RequestEvent, app core.App) error {
	// TODO: Implement in Step 3
	return e.JSON(http.StatusNotImplemented, map[string]string{
		"error": "Stop endpoint not yet implemented",
	})
}

func handleStartLive(e *core.RequestEvent, app core.App) error {
	// TODO: Implement in Step 4
	return e.JSON(http.StatusNotImplemented, map[string]string{
		"error": "Start live endpoint not yet implemented",
	})
}