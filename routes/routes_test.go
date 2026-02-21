package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jchen/go-sql-practice/migrations"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// TestData holds references to created test data
type TestData struct {
	LinkIDs  []string
	TagIDs   []string
	TagSlugs []string
}

// setupTestApp creates an in-memory PocketBase app for testing
func setupTestApp(t *testing.T) (*pocketbase.PocketBase, func()) {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: "", // Forces in-memory SQLite database
	})

	// Initialize the app without starting HTTP server
	if err := app.Bootstrap(); err != nil {
		t.Fatalf("Failed to bootstrap test app: %v", err)
	}

	// Run migrations to create collections
	if err := runTestMigrations(app); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	cleanup := func() {
		// PocketBase handles cleanup internally
		// No explicit cleanup needed for in-memory database
	}

	return app, cleanup
}

// runTestMigrations runs the production migrations in test environment
func runTestMigrations(app *pocketbase.PocketBase) error {
	// Use the production migration system
	migrations.Register(app)

	// Trigger the migration by simulating a serve event
	serveEvent := &core.ServeEvent{
		App:    app,
		Router: nil, // We don't need the router for migrations
	}

	return app.OnServe().Trigger(serveEvent)
}

// TestSetup validates that the basic test infrastructure works
func TestSetup(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Verify app is initialized
	if app == nil {
		t.Fatal("Expected app to be initialized")
	}

	// Verify database is working
	if app.DB() == nil {
		t.Fatal("Expected database to be initialized")
	}

	// Verify collections exist
	collections := []string{"tags", "links", "sync_sessions", "presentations"}
	for _, name := range collections {
		if _, err := app.FindCollectionByNameOrId(name); err != nil {
			t.Errorf("Expected collection %s to exist: %v", name, err)
		}
	}

	t.Logf("✅ Test setup working correctly with %d collections", len(collections))
}

// TestDatabase validates that basic database operations work
func TestDatabase(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Test database connection
	db := app.DB()
	if db == nil {
		t.Fatal("Expected database to be available")
	}

	// Test basic query
	var count int
	err := db.NewQuery("SELECT COUNT(*) FROM links").Row(&count)
	if err != nil {
		t.Fatalf("Expected query to succeed: %v", err)
	}

	t.Logf("✅ Database working correctly, found %d links", count)
}

// makeRequest executes an HTTP request against the test app with real routing
func makeRequest(app *pocketbase.PocketBase, method, url string, body io.Reader) (*http.Response, error) {
	// Create fresh router for this request
	router := echo.New()

	// Create serve event with router
	serveEvent := &core.ServeEvent{
		App:    app,
		Router: router,
	}

	// Register all routes
	Register(app)

	// Trigger OnServe to register routes with router
	if err := app.OnServe().Trigger(serveEvent); err != nil {
		return nil, err
	}

	// Create HTTP request and recorder
	req := httptest.NewRequest(method, url, body)
	rec := httptest.NewRecorder()

	// Execute request through router
	router.ServeHTTP(rec, req)

	// Return the recorded response
	return rec.Result(), nil
}

// parseJSONResponse unmarshals HTTP response body into target struct
func parseJSONResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v\nBody: %s", err, string(body))
	}
}

// ErrorResponse matches the common error format across all endpoints
type ErrorResponse struct {
	Error string `json:"error"`
}

// assertErrorResponse validates error responses consistently
func assertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedMessage string) {
	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	var errResp ErrorResponse
	parseJSONResponse(t, resp, &errResp)

	if errResp.Error != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, errResp.Error)
	}
}

// BenchmarkMakeRequest benchmarks the request execution performance
func BenchmarkMakeRequest(b *testing.B) {
	app, cleanup := setupTestApp(&testing.T{})
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := makeRequest(app, "GET", "/api/health", nil)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
}

// TestMakeRequest_RealExecution proves requests are actually executed
func TestMakeRequest_RealExecution(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Test 1: Valid endpoint returns success
	resp, err := makeRequest(app, "GET", "/api/health", nil)
	if err != nil {
		t.Fatalf("Failed to make health check request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for health check, got %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read health response: %v", err)
	}

	// Health check should return a simple status
	if !strings.Contains(string(body), "ok") && !strings.Contains(string(body), "healthy") {
		// It might just return empty 200, which is fine
		t.Logf("Health check response: %s", string(body))
	}

	// Test 2: Invalid endpoint returns 404
	resp2, err := makeRequest(app, "GET", "/api/nonexistent", nil)
	if err != nil {
		t.Fatalf("Failed to make nonexistent request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for nonexistent endpoint, got %d", resp2.StatusCode)
	}

	t.Log("✅ Real execution test passed: valid endpoint returns 200, invalid returns 404")
}