package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// createTestLinks creates multiple test links for search testing
func createTestLinks(app *pocketbase.PocketBase, links []struct {
	title       string
	description string
	tagSlugs    []string
	viewCount   int
}) ([]string, error) {
	var linkIds []string

	for _, linkData := range links {
		collection, err := app.FindCollectionByNameOrId("links")
		if err != nil {
			return nil, fmt.Errorf("failed to find links collection: %v", err)
		}

		record := core.NewRecord(collection)
		record.Set("title", linkData.title)
		record.Set("url", fmt.Sprintf("https://example.com/%s", strings.ReplaceAll(linkData.title, " ", "-")))
		record.Set("description", linkData.description)
		record.Set("view_count", linkData.viewCount)

		// Set tags if provided
		if len(linkData.tagSlugs) > 0 {
			tagIds, err := getTagIdsBySlugs(app, linkData.tagSlugs)
			if err != nil {
				return nil, fmt.Errorf("failed to get tag IDs: %v", err)
			}
			record.Set("tags", tagIds)
		}

		if err := app.Save(record); err != nil {
			return nil, fmt.Errorf("failed to save test link: %v", err)
		}

		linkIds = append(linkIds, record.Id)
	}

	return linkIds, nil
}

// getTagIdsBySlugs returns tag IDs for given slugs
func getTagIdsBySlugs(app *pocketbase.PocketBase, slugs []string) ([]string, error) {
	var tagIds []string

	for _, slug := range slugs {
		record, err := app.FindFirstRecordByData("tags", "slug", slug)
		if err != nil {
			return nil, fmt.Errorf("failed to find tag with slug %s: %v", slug, err)
		}
		tagIds = append(tagIds, record.Id)
	}

	return tagIds, nil
}

// buildSearchURL constructs search URL with query parameters
func buildSearchURL(params map[string]string) string {
	u, _ := url.Parse("/api/links/search")
	q := u.Query()
	for key, value := range params {
		if value != "" {
			q.Set(key, value)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// TestLinksSearch_BasicQuery tests basic text search functionality
func TestLinksSearch_BasicQuery(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Create test links with specific titles for searching
	testLinks := []struct {
		title       string
		description string
		tagSlugs    []string
		viewCount   int
	}{
		{"Go Programming Guide", "Learn Go programming language", []string{"golang"}, 10},
		{"JavaScript Basics", "Introduction to JavaScript", []string{"javascript"}, 20},
		{"Database Design", "How to design databases", []string{"database"}, 15},
		{"Go Testing Tutorial", "Testing in Go applications", []string{"golang", "testing"}, 5},
	}

	linkIds, err := createTestLinks(app, testLinks)
	if err != nil {
		t.Fatalf("Failed to create test links: %v", err)
	}
	_ = linkIds // We have the IDs if we need them

	// Test search for "Go"
	url := buildSearchURL(map[string]string{"q": "Go"})
	resp, err := makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make search request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	parseJSONResponse(t, resp, &searchResp)

	// Should find at least the two "Go" links
	if len(searchResp.Items) < 2 {
		t.Errorf("Expected at least 2 results for 'Go' search, got %d", len(searchResp.Items))
	}

	// Verify results contain "Go" in title or description
	for _, item := range searchResp.Items {
		if !strings.Contains(strings.ToLower(item.Title), "go") &&
			!strings.Contains(strings.ToLower(item.Description), "go") {
			t.Errorf("Search result doesn't contain 'Go': %s", item.Title)
		}
	}

	// Verify pagination fields
	if searchResp.Page != 1 {
		t.Errorf("Expected page 1, got %d", searchResp.Page)
	}

	if searchResp.PerPage != DefaultPerPage {
		t.Errorf("Expected perPage %d, got %d", DefaultPerPage, searchResp.PerPage)
	}

	t.Logf("✅ Basic query search passed: found %d results for 'Go'", len(searchResp.Items))
}

// TestLinksSearch_TagFilter tests filtering by tag slug
func TestLinksSearch_TagFilter(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Test with existing seed data - filter by "golang" tag
	url := buildSearchURL(map[string]string{"tag": "golang"})
	resp, err := makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make tag filter request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	parseJSONResponse(t, resp, &searchResp)

	// Should have at least one result (seed data has golang-tagged links)
	if len(searchResp.Items) == 0 {
		t.Error("Expected at least one result for 'golang' tag filter")
	}

	// Verify all results have the "golang" tag
	for _, item := range searchResp.Items {
		hasGolangTag := false
		for _, tag := range item.Tags {
			if tag == "golang" {
				hasGolangTag = true
				break
			}
		}
		if !hasGolangTag {
			t.Errorf("Search result missing 'golang' tag: %s", item.Title)
		}
	}

	t.Logf("✅ Tag filter search passed: found %d results for 'golang' tag", len(searchResp.Items))
}

// TestLinksSearch_CombinedFilters tests combining text search and tag filter
func TestLinksSearch_CombinedFilters(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Search for "documentation" text with "golang" tag
	url := buildSearchURL(map[string]string{
		"q":   "documentation",
		"tag": "golang",
	})
	resp, err := makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make combined filter request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	parseJSONResponse(t, resp, &searchResp)

	// Verify results match both text and tag criteria
	for _, item := range searchResp.Items {
		// Should contain "documentation" in title or description
		textMatch := strings.Contains(strings.ToLower(item.Title), "documentation") ||
			strings.Contains(strings.ToLower(item.Description), "documentation")

		// Should have "golang" tag
		hasGolangTag := false
		for _, tag := range item.Tags {
			if tag == "golang" {
				hasGolangTag = true
				break
			}
		}

		if !textMatch {
			t.Errorf("Combined filter result missing text match: %s", item.Title)
		}

		if !hasGolangTag {
			t.Errorf("Combined filter result missing golang tag: %s", item.Title)
		}
	}

	t.Logf("✅ Combined filters search passed: found %d results", len(searchResp.Items))
}

// TestLinksSearch_EmptyResults tests queries that return no matches
func TestLinksSearch_EmptyResults(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Search for something that definitely won't exist
	url := buildSearchURL(map[string]string{"q": "zyxwvutsrqponmlkjihgfedcba"})
	resp, err := makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make empty results request: %v", err)
	}
	defer resp.Body.Close()

	// Should return 200, not 404
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for empty results, got %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	parseJSONResponse(t, resp, &searchResp)

	// Should have empty items array, not nil
	if searchResp.Items == nil {
		t.Error("Expected empty items array, got nil")
	}

	if len(searchResp.Items) != 0 {
		t.Errorf("Expected 0 results, got %d", len(searchResp.Items))
	}

	// Should have correct pagination info
	if searchResp.TotalItems != 0 {
		t.Errorf("Expected TotalItems = 0, got %d", searchResp.TotalItems)
	}

	if searchResp.Page != 1 {
		t.Errorf("Expected page 1 even for empty results, got %d", searchResp.Page)
	}

	t.Log("✅ Empty results test passed: returns 200 with empty array")
}

// TestLinksSearch_Pagination tests pagination functionality
func TestLinksSearch_Pagination(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Create many test links to test pagination
	var testLinks []struct {
		title       string
		description string
		tagSlugs    []string
		viewCount   int
	}

	for i := 0; i < 25; i++ {
		testLinks = append(testLinks, struct {
			title       string
			description string
			tagSlugs    []string
			viewCount   int
		}{
			title:       fmt.Sprintf("Test Link %d", i),
			description: "Test description for pagination",
			tagSlugs:    []string{"testing"},
			viewCount:   i,
		})
	}

	_, err := createTestLinks(app, testLinks)
	if err != nil {
		t.Fatalf("Failed to create test links: %v", err)
	}

	// Test first page
	url := buildSearchURL(map[string]string{
		"q":       "Test",
		"page":    "1",
		"perPage": "10",
	})
	resp, err := makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make first page request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	parseJSONResponse(t, resp, &searchResp)

	// Should have 10 items (first page)
	if len(searchResp.Items) != 10 {
		t.Errorf("Expected 10 items on first page, got %d", len(searchResp.Items))
	}

	if searchResp.Page != 1 {
		t.Errorf("Expected page 1, got %d", searchResp.Page)
	}

	if searchResp.PerPage != 10 {
		t.Errorf("Expected perPage 10, got %d", searchResp.PerPage)
	}

	// TotalItems should be at least 25 (our test data)
	if searchResp.TotalItems < 25 {
		t.Errorf("Expected at least 25 total items, got %d", searchResp.TotalItems)
	}

	firstPageTotal := searchResp.TotalItems

	// Test second page
	url = buildSearchURL(map[string]string{
		"q":       "Test",
		"page":    "2",
		"perPage": "10",
	})
	resp, err = makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make second page request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	parseJSONResponse(t, resp, &searchResp)

	if searchResp.Page != 2 {
		t.Errorf("Expected page 2, got %d", searchResp.Page)
	}

	// TotalItems should be consistent across pages
	if searchResp.TotalItems != firstPageTotal {
		t.Errorf("TotalItems inconsistent between pages: page 1: %d, page 2: %d",
			firstPageTotal, searchResp.TotalItems)
	}

	t.Logf("✅ Pagination test passed: page 1: %d items, page 2: %d items, total: %d",
		10, len(searchResp.Items), searchResp.TotalItems)
}

// TestLinksSearch_ParameterValidation tests parameter validation
func TestLinksSearch_ParameterValidation(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Test invalid page number
	url := buildSearchURL(map[string]string{"page": "0"})
	resp, err := makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make invalid page request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid page, got %d", resp.StatusCode)
	}

	// Test invalid perPage number
	url = buildSearchURL(map[string]string{"perPage": "0"})
	resp, err = makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make invalid perPage request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid perPage, got %d", resp.StatusCode)
	}

	// Test perPage too high
	url = buildSearchURL(map[string]string{"perPage": "101"})
	resp, err = makeRequest(app, "GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to make high perPage request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for perPage > 100, got %d", resp.StatusCode)
	}

	t.Log("✅ Parameter validation test passed")
}

// TestLinksSearch_SQLInjectionProtection tests SQL injection protection
func TestLinksSearch_SQLInjectionProtection(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Test various SQL injection attempts
	injectionAttempts := []string{
		"'; DROP TABLE links; --",
		"' UNION SELECT * FROM users --",
		"' OR '1'='1",
		"<script>alert('xss')</script>",
		"' OR 1=1 --",
		"'; DELETE FROM links WHERE 1=1; --",
	}

	for _, attempt := range injectionAttempts {
		url := buildSearchURL(map[string]string{"q": attempt})
		resp, err := makeRequest(app, "GET", url, nil)
		if err != nil {
			t.Fatalf("Failed to make SQL injection test request: %v", err)
		}
		defer resp.Body.Close()

		// Should return 200 (treated as normal search query)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("SQL injection attempt should return 200, got %d for: %s", resp.StatusCode, attempt)
		}

		var searchResp SearchResponse
		parseJSONResponse(t, resp, &searchResp)

		// Response should be well-formed (no SQL errors)
		if searchResp.Items == nil {
			t.Errorf("SQL injection broke response format for: %s", attempt)
		}
	}

	// Test tag parameter injection
	for _, attempt := range injectionAttempts {
		url := buildSearchURL(map[string]string{"tag": attempt})
		resp, err := makeRequest(app, "GET", url, nil)
		if err != nil {
			t.Fatalf("Failed to make tag injection test request: %v", err)
		}
		defer resp.Body.Close()

		// Should return 200 (no matching tag is fine)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Tag injection attempt should return 200, got %d for: %s", resp.StatusCode, attempt)
		}
	}

	t.Log("✅ SQL injection protection test passed")
}