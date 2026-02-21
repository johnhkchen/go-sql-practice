package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// tagSeed represents seed data for a tag
type tagSeed struct {
	name string
	slug string
}

// linkSeed represents seed data for a link
type linkSeed struct {
	url         string
	title       string
	description string
	viewCount   int
	tagSlugs    []string
}

// tagSeeds defines the tags to create for development
var tagSeeds = []tagSeed{
	{"Go", "golang"},
	{"JavaScript", "javascript"},
	{"Database", "database"},
	{"DevOps", "devops"},
	{"Frontend", "frontend"},
	{"Backend", "backend"},
	{"Testing", "testing"},
	{"Architecture", "architecture"},
}

// linkSeeds defines the links to create for development
var linkSeeds = []linkSeed{
	{
		url:         "https://go.dev/doc/",
		title:       "Go Documentation",
		description: "Official Go programming language documentation",
		viewCount:   42,
		tagSlugs:    []string{"golang", "backend"},
	},
	{
		url:         "https://react.dev/",
		title:       "React",
		description: "The library for web and native user interfaces",
		viewCount:   38,
		tagSlugs:    []string{"javascript", "frontend"},
	},
	{
		url:         "https://www.postgresql.org/docs/",
		title:       "PostgreSQL Documentation",
		description: "The world's most advanced open source database",
		viewCount:   25,
		tagSlugs:    []string{"database", "backend"},
	},
	{
		url:         "https://hub.docker.com/",
		title:       "Docker Hub",
		description: "Container image library and community",
		viewCount:   31,
		tagSlugs:    []string{"devops"},
	},
	{
		url:         "https://go.dev/doc/tutorial/add-a-test",
		title:       "Testing in Go",
		description: "Learn how to write unit tests in Go",
		viewCount:   18,
		tagSlugs:    []string{"golang", "testing", "backend"},
	},
	{
		url:         "https://developer.mozilla.org/",
		title:       "MDN Web Docs",
		description: "Resources for developers, by developers",
		viewCount:   45,
		tagSlugs:    []string{"javascript", "frontend"},
	},
	{
		url:         "https://kubernetes.io/docs/",
		title:       "Kubernetes Documentation",
		description: "Production-grade container orchestration",
		viewCount:   22,
		tagSlugs:    []string{"devops", "architecture"},
	},
	{
		url:         "https://astro.build/",
		title:       "Astro",
		description: "The web framework for content-driven websites",
		viewCount:   15,
		tagSlugs:    []string{"frontend", "javascript"},
	},
	{
		url:         "https://pocketbase.io/docs/",
		title:       "PocketBase",
		description: "Open source backend in 1 file",
		viewCount:   8,
		tagSlugs:    []string{"database", "backend", "golang"},
	},
	{
		url:         "https://docs.github.com/actions",
		title:       "GitHub Actions Documentation",
		description: "Automate your workflow from idea to production",
		viewCount:   12,
		tagSlugs:    []string{"devops", "testing"},
	},
}

// seedData creates initial seed data for development
func seedData(txApp core.App) error {
	// Check if seed data already exists
	if seedDataExists(txApp) {
		return nil
	}

	// Create tags first
	tagMap, err := createSeedTags(txApp)
	if err != nil {
		return err
	}

	// Create links with tag relations
	return createSeedLinks(txApp, tagMap)
}

// seedDataExists checks if seed data has already been created
func seedDataExists(txApp core.App) bool {
	// Check for a known seed tag as indicator
	_, err := txApp.FindFirstRecordByData("tags", "slug", "golang")
	return err == nil
}

// createSeedTags creates all seed tags and returns a map of slug to ID
func createSeedTags(txApp core.App) (map[string]string, error) {
	tagsCollection, err := txApp.FindCollectionByNameOrId("tags")
	if err != nil {
		return nil, fmt.Errorf("failed to find tags collection: %w", err)
	}

	tagMap := make(map[string]string)
	for _, t := range tagSeeds {
		record := core.NewRecord(tagsCollection)
		record.Set("name", t.name)
		record.Set("slug", t.slug)

		if err := txApp.Save(record); err != nil {
			return nil, fmt.Errorf("failed to create tag %s: %w", t.slug, err)
		}
		tagMap[t.slug] = record.Id
	}

	return tagMap, nil
}

// createSeedLinks creates all seed links with their tag relations
func createSeedLinks(txApp core.App, tagMap map[string]string) error {
	linksCollection, err := txApp.FindCollectionByNameOrId("links")
	if err != nil {
		return fmt.Errorf("failed to find links collection: %w", err)
	}

	for _, l := range linkSeeds {
		record := core.NewRecord(linksCollection)
		record.Set("url", l.url)
		record.Set("title", l.title)
		record.Set("description", l.description)
		record.Set("view_count", l.viewCount)

		// Convert tag slugs to IDs
		tagIds := []string{}
		for _, slug := range l.tagSlugs {
			if id, ok := tagMap[slug]; ok {
				tagIds = append(tagIds, id)
			}
		}
		if len(tagIds) > 0 {
			record.Set("tags", tagIds)
		}

		if err := txApp.Save(record); err != nil {
			return fmt.Errorf("failed to create link %s: %w", l.title, err)
		}
	}

	return nil
}