package routes

import (
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// SQL query constants
const (
	// statsTopN defines the number of top items returned by stats queries
	statsTopN = 5

	sqlTotalLinks = "SELECT COUNT(*) as total FROM links"
	sqlTotalTags  = "SELECT COUNT(*) as total FROM tags"
	sqlTotalViews = "SELECT COALESCE(SUM(view_count), 0) as total FROM links"
	sqlMostViewed = `
		SELECT
			id,
			title,
			url,
			COALESCE(view_count, 0) as view_count
		FROM links
		ORDER BY view_count DESC
		LIMIT 5
	`
	sqlTopTags = `
		SELECT
			t.name,
			t.slug,
			(
				SELECT COUNT(*)
				FROM links l
				WHERE json_extract(l.tags, '$') LIKE '%' || t.id || '%'
			) as link_count
		FROM tags t
		ORDER BY link_count DESC
		LIMIT 5
	`
)

// StatsResponse represents the complete stats API response
type StatsResponse struct {
	TotalLinks int64       `json:"total_links"`
	TotalTags  int64       `json:"total_tags"`
	TotalViews int64       `json:"total_views"`
	TopTags    []TagStats  `json:"top_tags"`
	MostViewed []LinkStats `json:"most_viewed"`
}

// TagStats represents tag statistics
type TagStats struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	LinkCount int64  `json:"link_count"`
}

// LinkStats represents link statistics
type LinkStats struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	ViewCount int64  `json:"view_count"`
}

// registerStats registers the stats endpoint
func registerStats(e *core.ServeEvent) {
	e.Router.GET("/api/stats", func(ev *core.RequestEvent) error {
		return handleGetStats(ev, e.App)
	})
}

// handleGetStats handles GET /api/stats requests
func handleGetStats(e *core.RequestEvent, app core.App) error {
	db := app.DB()

	// Get total links
	totalLinks, err := getTotalLinks(db)
	if err != nil {
		app.Logger().Error("getTotalLinks failed", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get total links"})
	}

	// Get total tags
	totalTags, err := getTotalTags(db)
	if err != nil {
		app.Logger().Error("getTotalTags failed", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get total tags"})
	}

	// Get total views
	totalViews, err := getTotalViews(db)
	if err != nil {
		app.Logger().Error("getTotalViews failed", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get total views"})
	}

	// Get most viewed links
	mostViewed, err := getMostViewed(db)
	if err != nil {
		app.Logger().Error("getMostViewed failed", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get most viewed links"})
	}

	// Get top tags
	topTags, err := getTopTags(db)
	if err != nil {
		app.Logger().Error("getTopTags failed", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get top tags"})
	}

	response := StatsResponse{
		TotalLinks: totalLinks,
		TotalTags:  totalTags,
		TotalViews: totalViews,
		TopTags:    topTags,
		MostViewed: mostViewed,
	}

	return e.JSON(http.StatusOK, response)
}

// getTotalLinks returns the total number of links
func getTotalLinks(db dbx.Builder) (int64, error) {
	var total int64
	err := db.NewQuery(sqlTotalLinks).Row(&total)
	return total, err
}

// getTotalTags returns the total number of tags
func getTotalTags(db dbx.Builder) (int64, error) {
	var total int64
	err := db.NewQuery(sqlTotalTags).Row(&total)
	return total, err
}

// getTotalViews returns the sum of all view counts
func getTotalViews(db dbx.Builder) (int64, error) {
	var total int64
	err := db.NewQuery(sqlTotalViews).Row(&total)
	return total, err
}

// getMostViewed returns the top 5 most viewed links
func getMostViewed(db dbx.Builder) ([]LinkStats, error) {
	links := []LinkStats{}
	rows, err := db.NewQuery(sqlMostViewed).Rows()
	if err != nil {
		return links, err
	}
	defer rows.Close()

	for rows.Next() {
		var link LinkStats
		err := rows.Scan(&link.ID, &link.Title, &link.URL, &link.ViewCount)
		if err != nil {
			return links, err
		}
		links = append(links, link)
	}

	return links, rows.Err()
}

// getTopTags returns the top 5 tags by link count
func getTopTags(db dbx.Builder) ([]TagStats, error) {
	tags := []TagStats{}
	rows, err := db.NewQuery(sqlTopTags).Rows()
	if err != nil {
		return tags, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag TagStats
		err := rows.Scan(&tag.Name, &tag.Slug, &tag.LinkCount)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}