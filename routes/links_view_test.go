package routes

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// TestLinksView_Success tests successful view count increment
func TestLinksView_Success(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Get a link ID from the database
	var linkID string
	err := app.DB().NewQuery("SELECT id FROM links LIMIT 1").Row(&linkID)
	if err != nil {
		t.Fatalf("Failed to get test link ID: %v", err)
	}

	// Get initial view count
	var initialCount int
	err = app.DB().NewQuery("SELECT COALESCE(view_count, 0) FROM links WHERE id = {:id}").
		Bind(map[string]interface{}{"id": linkID}).
		Row(&initialCount)
	if err != nil {
		t.Fatalf("Failed to get initial view count: %v", err)
	}

	// Make request to increment view count
	url := fmt.Sprintf("/api/links/%s/view", linkID)
	resp, err := makeRequest(app, "POST", url, nil)
	if err != nil {
		t.Fatalf("Failed to make view request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var viewResp struct {
		ID        string `json:"id"`
		ViewCount int    `json:"view_count"`
	}
	parseJSONResponse(t, resp, &viewResp)

	// Verify the response
	if viewResp.ID != linkID {
		t.Errorf("Expected link ID %s, got %s", linkID, viewResp.ID)
	}

	if viewResp.ViewCount != initialCount+1 {
		t.Errorf("Expected view count %d, got %d", initialCount+1, viewResp.ViewCount)
	}

	// Verify in database
	var finalCount int
	err = app.DB().NewQuery("SELECT view_count FROM links WHERE id = {:id}").
		Bind(map[string]interface{}{"id": linkID}).
		Row(&finalCount)
	if err != nil {
		t.Fatalf("Failed to get final view count: %v", err)
	}

	if finalCount != initialCount+1 {
		t.Errorf("Database view count not updated: expected %d, got %d", initialCount+1, finalCount)
	}

	t.Logf("✅ View count increment test passed: %d -> %d", initialCount, finalCount)
}

// TestLinksView_NotFound tests view count on non-existent link
func TestLinksView_NotFound(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Use a non-existent ID
	url := "/api/links/nonexistentid12345/view"
	resp, err := makeRequest(app, "POST", url, nil)
	if err != nil {
		t.Fatalf("Failed to make view request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent link, got %d", resp.StatusCode)
	}

	t.Log("✅ View count not found test passed")
}

// TestLinksView_InvalidID tests view count with invalid ID format
func TestLinksView_InvalidID(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Test with various invalid IDs
	invalidIDs := []string{
		"",
		"../etc/passwd",
		"'; DROP TABLE links; --",
		"<script>alert('xss')</script>",
	}

	for _, id := range invalidIDs {
		url := fmt.Sprintf("/api/links/%s/view", id)
		resp, err := makeRequest(app, "POST", url, nil)
		if err != nil {
			t.Fatalf("Failed to make view request with ID %q: %v", id, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Errorf("Expected error status for invalid ID %q, got 200", id)
		}
	}

	t.Log("✅ Invalid ID test passed")
}

// TestLinksView_Concurrent tests concurrent view count increments
func TestLinksView_Concurrent(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Get a link ID from the database
	var linkID string
	err := app.DB().NewQuery("SELECT id FROM links LIMIT 1").Row(&linkID)
	if err != nil {
		t.Fatalf("Failed to get test link ID: %v", err)
	}

	// Get initial view count
	var initialCount int
	err = app.DB().NewQuery("SELECT COALESCE(view_count, 0) FROM links WHERE id = {:id}").
		Bind(map[string]interface{}{"id": linkID}).
		Row(&initialCount)
	if err != nil {
		t.Fatalf("Failed to get initial view count: %v", err)
	}

	// Number of concurrent increments
	concurrentRequests := 10
	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	// Track success count
	successCount := 0
	var mu sync.Mutex

	// Launch concurrent requests
	for i := 0; i < concurrentRequests; i++ {
		go func() {
			defer wg.Done()

			url := fmt.Sprintf("/api/links/%s/view", linkID)
			resp, err := makeRequest(app, "POST", url, nil)
			if err != nil {
				t.Errorf("Concurrent request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	// Wait for all requests to complete
	wg.Wait()

	// Verify final count in database
	var finalCount int
	err = app.DB().NewQuery("SELECT view_count FROM links WHERE id = {:id}").
		Bind(map[string]interface{}{"id": linkID}).
		Row(&finalCount)
	if err != nil {
		t.Fatalf("Failed to get final view count: %v", err)
	}

	// The final count should be initial + number of successful requests
	expectedCount := initialCount + successCount
	if finalCount != expectedCount {
		t.Errorf("Concurrent increments failed: expected %d, got %d (initial: %d, successful requests: %d)",
			expectedCount, finalCount, initialCount, successCount)
	}

	t.Logf("✅ Concurrent view count test passed: %d concurrent requests, %d -> %d",
		concurrentRequests, initialCount, finalCount)
}

// Helper to create a test link with known ID
func createTestLinkWithID(app *pocketbase.PocketBase) (string, error) {
	collection, err := app.FindCollectionByNameOrId("links")
	if err != nil {
		return "", err
	}

	record := core.NewRecord(collection)
	record.Set("title", "Test Link for View Count")
	record.Set("url", "https://example.com/test")
	record.Set("description", "Test link for view count testing")
	record.Set("view_count", 0)

	if err := app.Save(record); err != nil {
		return "", err
	}

	return record.Id, nil
}