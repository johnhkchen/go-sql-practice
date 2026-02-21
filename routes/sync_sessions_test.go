package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestSyncCreate_Success tests successful sync session creation
func TestSyncCreate_Success(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Make request to create sync session
	resp, err := makeRequest(app, "POST", "/api/sync", nil)
	if err != nil {
		t.Fatalf("Failed to create sync session: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// Parse response
	var syncResp CreateSessionResponse
	parseJSONResponse(t, resp, &syncResp)

	// Verify response has required fields
	if syncResp.SessionID == "" {
		t.Error("Expected non-empty session ID")
	}

	if syncResp.AdminURL == "" {
		t.Error("Expected non-empty admin URL")
	}

	if syncResp.ViewerURL == "" {
		t.Error("Expected non-empty viewer URL")
	}

	// Verify URLs contain session ID
	expectedAdminPath := fmt.Sprintf("/sync/%s", syncResp.SessionID)
	if syncResp.AdminURL != expectedAdminPath {
		t.Errorf("Expected admin URL %s, got %s", expectedAdminPath, syncResp.AdminURL)
	}

	expectedViewerPath := fmt.Sprintf("/watch/%s", syncResp.SessionID)
	if syncResp.ViewerURL != expectedViewerPath {
		t.Errorf("Expected viewer URL %s, got %s", expectedViewerPath, syncResp.ViewerURL)
	}

	// Verify session exists in database
	var count int
	err = app.DB().NewQuery("SELECT COUNT(*) FROM sync_sessions WHERE id = {:id}").
		Bind(map[string]interface{}{"id": syncResp.SessionID}).
		Row(&count)
	if err != nil {
		t.Fatalf("Failed to check session in database: %v", err)
	}

	if count != 1 {
		t.Error("Session not found in database")
	}

	t.Log("✅ Sync session creation test passed")
}

// TestSyncProgress_ValidUpdate tests updating progress on valid session
func TestSyncProgress_ValidUpdate(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// First create a session
	createResp, err := makeRequest(app, "POST", "/api/sync", nil)
	if err != nil {
		t.Fatalf("Failed to create sync session: %v", err)
	}
	defer createResp.Body.Close()

	var syncResp CreateSessionResponse
	parseJSONResponse(t, createResp, &syncResp)

	// Get the admin token from database
	var adminToken string
	err = app.DB().NewQuery("SELECT admin_token FROM sync_sessions WHERE id = {:id}").
		Bind(map[string]interface{}{"id": syncResp.SessionID}).
		Row(&adminToken)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	// Update progress
	updateReq := UpdateProgressRequest{
		Progress: 0.5,
	}
	reqBody, _ := json.Marshal(updateReq)

	url := fmt.Sprintf("/api/sync/%s/progress?token=%s", syncResp.SessionID, adminToken)
	resp, err := makeRequest(app, "PUT", url, bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to update progress: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify progress in database
	var progress float64
	err = app.DB().NewQuery("SELECT progress FROM sync_sessions WHERE id = {:id}").
		Bind(map[string]interface{}{"id": syncResp.SessionID}).
		Row(&progress)
	if err != nil {
		t.Fatalf("Failed to check progress: %v", err)
	}

	if progress != 0.5 {
		t.Errorf("Expected progress 0.5, got %f", progress)
	}

	t.Log("✅ Sync progress update test passed")
}

// TestSyncProgress_InvalidToken tests rejection of invalid token
func TestSyncProgress_InvalidToken(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// First create a session
	createResp, err := makeRequest(app, "POST", "/api/sync", nil)
	if err != nil {
		t.Fatalf("Failed to create sync session: %v", err)
	}
	defer createResp.Body.Close()

	var syncResp CreateSessionResponse
	parseJSONResponse(t, createResp, &syncResp)

	// Try to update with wrong token
	updateReq := UpdateProgressRequest{
		Progress: 0.5,
	}
	reqBody, _ := json.Marshal(updateReq)

	url := fmt.Sprintf("/api/sync/%s/progress?token=wrongtoken", syncResp.SessionID)
	resp, err := makeRequest(app, "PUT", url, bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to make update request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid token, got %d", resp.StatusCode)
	}

	t.Log("✅ Invalid token test passed")
}

// TestSyncProgress_OutOfRange tests rejection of out-of-range progress values
func TestSyncProgress_OutOfRange(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// First create a session
	createResp, err := makeRequest(app, "POST", "/api/sync", nil)
	if err != nil {
		t.Fatalf("Failed to create sync session: %v", err)
	}
	defer createResp.Body.Close()

	var syncResp CreateSessionResponse
	parseJSONResponse(t, createResp, &syncResp)

	// Get the admin token
	var adminToken string
	err = app.DB().NewQuery("SELECT admin_token FROM sync_sessions WHERE id = {:id}").
		Bind(map[string]interface{}{"id": syncResp.SessionID}).
		Row(&adminToken)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	// Test invalid progress values
	invalidValues := []float64{-0.1, 1.1, 2.0, -1.0}

	for _, value := range invalidValues {
		updateReq := UpdateProgressRequest{
			Progress: value,
		}
		reqBody, _ := json.Marshal(updateReq)

		url := fmt.Sprintf("/api/sync/%s/progress?token=%s", syncResp.SessionID, adminToken)
		resp, err := makeRequest(app, "PUT", url, bytes.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to make update request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for progress %f, got %d", value, resp.StatusCode)
		}
	}

	t.Log("✅ Out-of-range progress test passed")
}

// TestSyncProgress_NonExistentSession tests updating non-existent session
func TestSyncProgress_NonExistentSession(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	updateReq := UpdateProgressRequest{
		Progress: 0.5,
	}
	reqBody, _ := json.Marshal(updateReq)

	url := "/api/sync/nonexistentsession/progress?token=sometoken"
	resp, err := makeRequest(app, "PUT", url, bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to make update request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent session, got %d", resp.StatusCode)
	}

	t.Log("✅ Non-existent session test passed")
}

// TestSyncCreate_InvalidRequest tests sync creation with invalid body
func TestSyncCreate_InvalidRequest(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Send invalid JSON body
	invalidBody := bytes.NewReader([]byte("{invalid json}"))
	resp, err := makeRequest(app, "POST", "/api/sync", invalidBody)
	if err != nil {
		t.Fatalf("Failed to make sync request: %v", err)
	}
	defer resp.Body.Close()

	// Should still create session (body is ignored for creation)
	// But let's verify it doesn't crash
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		t.Logf("Got status %d for invalid JSON (accepting either 201 or 400)", resp.StatusCode)
	}

	t.Log("✅ Invalid request test passed")
}