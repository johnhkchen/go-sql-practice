package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

const (
	DefaultPage2    = 1
	DefaultPerPage2 = 20
	MaxPerPage2     = 100
)

// SearchParams2 represents the search query parameters
type SearchParams2 struct {
	Q       string
	Tag     string
	Page    int
	PerPage int
}

// SearchResponse2 represents the search API response
type SearchResponse2 struct {
	Items      []LinkItem2 `json:"items"`
	Page       int         `json:"page"`
	PerPage    int         `json:"perPage"`
	TotalItems int         `json:"totalItems"`
}

// LinkItem2 represents a single link in the search results
type LinkItem2 struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ViewCount   int      `json:"view_count"`
	Tags        []string `json:"tags"`
}

// registerLinksSearchSimple registers a simple search endpoint for testing
func registerLinksSearchSimple(e *core.ServeEvent) {
	e.Router.GET("/api/links/search-simple", func(ev *core.RequestEvent) error {
		return handleSearchSimple(ev, e.App)
	})
}

// handleSearchSimple handles the simple search endpoint
func handleSearchSimple(e *core.RequestEvent, app core.App) error {
	// Parse parameters
	params := parseSearchParams2(e.Request)

	// Validate parameters
	if err := validateSearchParams2(params); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Simple query without tag filtering for now
	db := app.DB()

	// Build simple query
	query := `
		SELECT
			id,
			url,
			title,
			COALESCE(description, '') as description,
			COALESCE(view_count, 0) as view_count
		FROM links
	`

	var whereClauses []string
	queryParams := make(map[string]interface{})

	// Add text search condition if provided
	if params.Q != "" {
		searchPattern := escapeLikePattern2(params.Q)
		whereClauses = append(whereClauses, "(title LIKE {searchPattern} OR description LIKE {searchPattern})")
		queryParams["searchPattern"] = searchPattern
	}

	// Add WHERE clause if conditions exist
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Add ordering and pagination
	query += " ORDER BY id DESC LIMIT {limit} OFFSET {offset}"
	queryParams["limit"] = params.PerPage
	queryParams["offset"] = (params.Page - 1) * params.PerPage

	// Execute query with parameter binding to prevent SQL injection
	rows, err := db.NewQuery(query).Bind(dbx.Params(queryParams)).Rows()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Query failed: %v", err),
		})
	}
	defer rows.Close()

	// Scan results
	var links []LinkItem2
	for rows.Next() {
		var link LinkItem2
		err := rows.Scan(
			&link.ID,
			&link.URL,
			&link.Title,
			&link.Description,
			&link.ViewCount,
		)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Row scan failed: %v", err),
			})
		}
		link.Tags = []string{} // Empty for now
		links = append(links, link)
	}

	if links == nil {
		links = []LinkItem2{}
	}

	// Simple count - just use length for now
	totalItems := len(links)

	// Build response
	response := SearchResponse2{
		Items:      links,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
	}

	return e.JSON(http.StatusOK, response)
}

// parseSearchParams2 extracts and parses query parameters
func parseSearchParams2(r *http.Request) SearchParams2 {
	params := SearchParams2{
		Q:       r.URL.Query().Get("q"),
		Tag:     r.URL.Query().Get("tag"),
		Page:    DefaultPage2,
		PerPage: DefaultPerPage2,
	}

	// Parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			params.Page = page
		}
	}

	// Parse perPage
	if perPageStr := r.URL.Query().Get("perPage"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil {
			params.PerPage = perPage
		}
	}

	return params
}

// validateSearchParams2 validates the search parameters
func validateSearchParams2(params SearchParams2) error {
	if params.Page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	if params.PerPage < 1 || params.PerPage > MaxPerPage2 {
		return fmt.Errorf("perPage must be between 1 and %d", MaxPerPage2)
	}
	return nil
}

// escapeLikePattern2 escapes special characters for SQL LIKE queries
func escapeLikePattern2(pattern string) string {
	// Escape LIKE wildcards (single quotes handled by parameter binding)
	pattern = strings.ReplaceAll(pattern, "\\", "\\\\")
	pattern = strings.ReplaceAll(pattern, "%", "\\%")
	pattern = strings.ReplaceAll(pattern, "_", "\\_")

	// Add wildcards for partial matching
	return "%" + pattern + "%"
}