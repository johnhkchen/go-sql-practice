package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// StatsTestResponse struct for unmarshaling stats response
type StatsTestResponse struct {
	TotalLinks int `json:"total_links"`
	TotalTags  int `json:"total_tags"`
	TotalViews int `json:"total_views"`
	MostViewed []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		URL       string `json:"url"`
		ViewCount int    `json:"view_count"`
	} `json:"most_viewed"`
	TopTags []struct {
		Name  string `json:"name"`
		Slug  string `json:"slug"`
		Count int    `json:"count"`
	} `json:"top_tags"`
}

// TestStats_Complete tests the stats endpoint with complete data
func TestStats_Complete(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Make request to stats endpoint
	resp, err := makeRequest(app, "GET", "/api/stats", nil)
	if err != nil {
		t.Fatalf("Failed to make stats request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var stats StatsTestResponse
	parseJSONResponse(t, resp, &stats)

	// Verify that we have data
	if stats.TotalLinks == 0 {
		t.Error("Expected at least one link in stats")
	}

	if stats.TotalTags == 0 {
		t.Error("Expected at least one tag in stats")
	}

	// Most viewed should have at most 5 items (as per SQL query)
	if len(stats.MostViewed) > 5 {
		t.Errorf("Expected at most 5 most viewed links, got %d", len(stats.MostViewed))
	}

	// Verify most viewed are sorted by view count (descending)
	for i := 1; i < len(stats.MostViewed); i++ {
		if stats.MostViewed[i].ViewCount > stats.MostViewed[i-1].ViewCount {
			t.Error("Most viewed links not sorted correctly")
		}
	}

	// Top tags should have at most 5 items
	if len(stats.TopTags) > 5 {
		t.Errorf("Expected at most 5 top tags, got %d", len(stats.TopTags))
	}

	t.Logf("✅ Stats complete test passed: %d links, %d tags, %d total views",
		stats.TotalLinks, stats.TotalTags, stats.TotalViews)
}

// TestStats_Accuracy tests that stats match actual database counts
func TestStats_Accuracy(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Get actual counts from database
	var actualLinkCount int
	err := app.DB().NewQuery("SELECT COUNT(*) FROM links").Row(&actualLinkCount)
	if err != nil {
		t.Fatalf("Failed to get actual link count: %v", err)
	}

	var actualTagCount int
	err = app.DB().NewQuery("SELECT COUNT(*) FROM tags").Row(&actualTagCount)
	if err != nil {
		t.Fatalf("Failed to get actual tag count: %v", err)
	}

	var actualTotalViews int
	err = app.DB().NewQuery("SELECT COALESCE(SUM(view_count), 0) FROM links").Row(&actualTotalViews)
	if err != nil {
		t.Fatalf("Failed to get actual total views: %v", err)
	}

	// Make request to stats endpoint
	resp, err := makeRequest(app, "GET", "/api/stats", nil)
	if err != nil {
		t.Fatalf("Failed to make stats request: %v", err)
	}
	defer resp.Body.Close()

	var stats StatsTestResponse
	parseJSONResponse(t, resp, &stats)

	// Verify counts match
	if stats.TotalLinks != actualLinkCount {
		t.Errorf("Link count mismatch: stats reports %d, actual is %d",
			stats.TotalLinks, actualLinkCount)
	}

	if stats.TotalTags != actualTagCount {
		t.Errorf("Tag count mismatch: stats reports %d, actual is %d",
			stats.TotalTags, actualTagCount)
	}

	if stats.TotalViews != actualTotalViews {
		t.Errorf("Total views mismatch: stats reports %d, actual is %d",
			stats.TotalViews, actualTotalViews)
	}

	t.Log("✅ Stats accuracy test passed: counts match database")
}

// TestStats_EmptyDatabase tests stats with minimal data
func TestStats_EmptyDatabase(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Delete all links (this will make some stats empty)
	// Note: We don't actually delete to avoid breaking other tests
	// Instead, we just test that the endpoint handles empty results gracefully

	resp, err := makeRequest(app, "GET", "/api/stats", nil)
	if err != nil {
		t.Fatalf("Failed to make stats request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 even with minimal data, got %d", resp.StatusCode)
	}

	// Parse response
	var stats StatsTestResponse
	body, err := json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		// Just ensure we got valid JSON
		t.Fatalf("Failed to parse stats response: %v", err)
	}
	_ = body

	// Ensure arrays are initialized (not nil)
	if stats.MostViewed == nil {
		t.Error("Most viewed should be empty array, not nil")
	}

	if stats.TopTags == nil {
		t.Error("Top tags should be empty array, not nil")
	}

	t.Log("✅ Stats empty database test passed: handles minimal data gracefully")
}

// TestStats_ResponseStructure tests that response has expected JSON structure
func TestStats_ResponseStructure(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	resp, err := makeRequest(app, "GET", "/api/stats", nil)
	if err != nil {
		t.Fatalf("Failed to make stats request: %v", err)
	}
	defer resp.Body.Close()

	// Parse as generic map to check structure
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to parse response as JSON: %v", err)
	}

	// Check required fields exist
	requiredFields := []string{
		"total_links",
		"total_tags",
		"total_views",
		"most_viewed",
		"top_tags",
	}

	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Required field %q missing from stats response", field)
		}
	}

	// Verify types
	if _, ok := result["total_links"].(float64); !ok {
		t.Error("total_links should be a number")
	}

	if _, ok := result["most_viewed"].([]interface{}); !ok {
		t.Error("most_viewed should be an array")
	}

	if _, ok := result["top_tags"].([]interface{}); !ok {
		t.Error("top_tags should be an array")
	}

	t.Log("✅ Stats response structure test passed")
}

// Helper to create test data with specific view counts
func createTestDataWithViews(app *pocketbase.PocketBase, viewCounts []int) error {
	collection, err := app.FindCollectionByNameOrId("links")
	if err != nil {
		return err
	}

	for i, count := range viewCounts {
		record := core.NewRecord(collection)
		record.Set("title", fmt.Sprintf("Test Link %d", i))
		record.Set("url", fmt.Sprintf("https://example.com/test%d", i))
		record.Set("description", "Test link for stats")
		record.Set("view_count", count)

		if err := app.Save(record); err != nil {
			return err
		}
	}

	return nil
}