package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// SearchParams represents the search query parameters
type SearchParams struct {
	Q       string
	Tag     string
	Page    int
	PerPage int
}

// SearchResponse represents the search API response
type SearchResponse struct {
	Items      []LinkItem `json:"items"`
	Page       int        `json:"page"`
	PerPage    int        `json:"perPage"`
	TotalItems int        `json:"totalItems"`
}

// LinkItem represents a single link in the search results
type LinkItem struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ViewCount   int      `json:"view_count"`
	Tags        []string `json:"tags"`
	Created     string   `json:"created"`
	Updated     string   `json:"updated"`
}

// registerLinksSearch registers the search endpoint
func registerLinksSearch(e *core.ServeEvent) {
	e.Router.GET("/api/links/search", func(ev *core.RequestEvent) error {
		return handleSearch(ev, e.App)
	})
}

// handleSearch handles the search endpoint
func handleSearch(e *core.RequestEvent, app core.App) error {
	// Parse parameters
	params := parseSearchParams(e.Request)

	// Validate parameters
	if err := validateSearchParams(params); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Execute search query
	links, err := executeSearchQuery(app, params)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to execute search",
		})
	}

	// Get total count
	totalItems, err := executeCountQuery(app, params)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get total count",
		})
	}

	// Fetch tags for links
	if len(links) > 0 {
		linkIDs := make([]string, len(links))
		for i, link := range links {
			linkIDs[i] = link.ID
		}

		tagMap, err := fetchTagsForLinks(app, linkIDs)
		if err != nil {
			// Log error but don't fail the request
			// Links will be returned without tags
			tagMap = make(map[string][]string)
		}

		// Assign tags to links
		for i := range links {
			if tags, ok := tagMap[links[i].ID]; ok {
				links[i].Tags = tags
			} else {
				links[i].Tags = []string{}
			}
		}
	}

	// Build response
	response := SearchResponse{
		Items:      links,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
	}

	return e.JSON(http.StatusOK, response)
}

// parseSearchParams extracts and parses query parameters
func parseSearchParams(r *http.Request) SearchParams {
	params := SearchParams{
		Q:       r.URL.Query().Get("q"),
		Tag:     r.URL.Query().Get("tag"),
		Page:    DefaultPage,
		PerPage: DefaultPerPage,
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

// validateSearchParams validates the search parameters
func validateSearchParams(params SearchParams) error {
	if params.Page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	if params.PerPage < 1 || params.PerPage > MaxPerPage {
		return fmt.Errorf("perPage must be between 1 and %d", MaxPerPage)
	}
	return nil
}

// escapeLikePattern escapes special characters for SQL LIKE queries
func escapeLikePattern(pattern string) string {
	// Escape special characters
	pattern = strings.ReplaceAll(pattern, "\\", "\\\\")
	pattern = strings.ReplaceAll(pattern, "%", "\\%")
	pattern = strings.ReplaceAll(pattern, "_", "\\_")

	// Add wildcards for partial matching
	return "%" + pattern + "%"
}

// executeSearchQuery executes the main search query
func executeSearchQuery(app core.App, params SearchParams) ([]LinkItem, error) {
	db := app.DB()

	// Build query
	query := `
		SELECT DISTINCT
			l.id,
			l.url,
			l.title,
			COALESCE(l.description, '') as description,
			COALESCE(l.view_count, 0) as view_count,
			l.created,
			l.updated
		FROM links l
	`

	var args []interface{}
	var whereClauses []string

	// Add JOINs if tag filter is present
	if params.Tag != "" {
		query += `
			JOIN json_each(l.tags) AS jt ON 1=1
			JOIN tags t ON t.id = jt.value
		`
		whereClauses = append(whereClauses, "t.slug = ?")
		args = append(args, params.Tag)
	}

	// Add text search condition
	if params.Q != "" {
		searchPattern := escapeLikePattern(params.Q)
		whereClauses = append(whereClauses, "(l.title LIKE ? OR l.description LIKE ?)")
		args = append(args, searchPattern, searchPattern)
	}

	// Add WHERE clause if conditions exist
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Add ordering and pagination
	query += " ORDER BY l.created DESC LIMIT ? OFFSET ?"
	args = append(args, params.PerPage, (params.Page-1)*params.PerPage)

	// Execute query using string interpolation for parameters since dbx doesn't support varargs
	finalQuery := query
	for _, arg := range args {
		argStr := fmt.Sprintf("%v", arg)
		// For numeric values (int), don't add quotes
		switch arg.(type) {
		case int:
			finalQuery = strings.Replace(finalQuery, "?", argStr, 1)
		default:
			// For strings, add quotes and escape single quotes
			argStr = strings.ReplaceAll(argStr, "'", "''")
			finalQuery = strings.Replace(finalQuery, "?", "'"+argStr+"'", 1)
		}
	}
	rows, err := db.NewQuery(finalQuery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	var links []LinkItem
	for rows.Next() {
		var link LinkItem
		err := rows.Scan(
			&link.ID,
			&link.URL,
			&link.Title,
			&link.Description,
			&link.ViewCount,
			&link.Created,
			&link.Updated,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if links == nil {
		links = []LinkItem{}
	}

	return links, rows.Err()
}

// executeCountQuery executes the count query for pagination
func executeCountQuery(app core.App, params SearchParams) (int, error) {
	db := app.DB()

	// Build count query
	query := `
		SELECT COUNT(DISTINCT l.id)
		FROM links l
	`

	var args []interface{}
	var whereClauses []string

	// Add JOINs if tag filter is present
	if params.Tag != "" {
		query += `
			JOIN json_each(l.tags) AS jt ON 1=1
			JOIN tags t ON t.id = jt.value
		`
		whereClauses = append(whereClauses, "t.slug = ?")
		args = append(args, params.Tag)
	}

	// Add text search condition
	if params.Q != "" {
		searchPattern := escapeLikePattern(params.Q)
		whereClauses = append(whereClauses, "(l.title LIKE ? OR l.description LIKE ?)")
		args = append(args, searchPattern, searchPattern)
	}

	// Add WHERE clause if conditions exist
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Execute count query using string interpolation for parameters
	var count int
	finalQuery := query
	for _, arg := range args {
		argStr := fmt.Sprintf("%v", arg)
		// For numeric values (int), don't add quotes
		switch arg.(type) {
		case int:
			finalQuery = strings.Replace(finalQuery, "?", argStr, 1)
		default:
			// For strings, add quotes and escape single quotes
			argStr = strings.ReplaceAll(argStr, "'", "''")
			finalQuery = strings.Replace(finalQuery, "?", "'"+argStr+"'", 1)
		}
	}
	err := db.NewQuery(finalQuery).Row(&count)

	return count, err
}

// fetchTagsForLinks fetches tags for the given link IDs
func fetchTagsForLinks(app core.App, linkIDs []string) (map[string][]string, error) {
	if len(linkIDs) == 0 {
		return make(map[string][]string), nil
	}

	db := app.DB()

	// Build placeholders for IN clause
	placeholders := make([]string, len(linkIDs))
	args := make([]interface{}, len(linkIDs))
	for i, id := range linkIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT l.id, t.slug
		FROM links l
		JOIN json_each(l.tags) AS jt ON 1=1
		JOIN tags t ON t.id = jt.value
		WHERE l.id IN (%s)
		ORDER BY t.slug
	`, strings.Join(placeholders, ","))

	// Execute tag query using string interpolation for parameters
	finalQuery := query
	for _, arg := range args {
		argStr := fmt.Sprintf("%v", arg)
		// For numeric values (int), don't add quotes
		switch arg.(type) {
		case int:
			finalQuery = strings.Replace(finalQuery, "?", argStr, 1)
		default:
			// For strings, add quotes and escape single quotes
			argStr = strings.ReplaceAll(argStr, "'", "''")
			finalQuery = strings.Replace(finalQuery, "?", "'"+argStr+"'", 1)
		}
	}
	rows, err := db.NewQuery(finalQuery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build map of link_id to tags
	tagMap := make(map[string][]string)
	for rows.Next() {
		var linkID, tagSlug string
		if err := rows.Scan(&linkID, &tagSlug); err != nil {
			return nil, err
		}
		tagMap[linkID] = append(tagMap[linkID], tagSlug)
	}

	return tagMap, rows.Err()
}