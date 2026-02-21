package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Register registers all migrations with the PocketBase app
func Register(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Run migrations once when server starts
		if err := createCollections(e.App); err != nil {
			return err
		}

		// Seed data after collections are created
		if err := seedData(e.App); err != nil {
			return err
		}

		return e.Next()
	})
}

// createCollections creates the tags and links collections
func createCollections(txApp core.App) error {
	// Check if tags collection already exists
	if _, err := txApp.FindCollectionByNameOrId("tags"); err == nil {
		// Collection already exists, skip creation of tags/links
		// But still need to check sync_sessions and presentations
		if err := createSyncSessionsCollection(txApp); err != nil {
			return err
		}
		return createPresentationsCollection(txApp)
	}

	// Create tags collection first
	tagsCollection := core.NewBaseCollection("tags")
	tagsCollection.Type = core.CollectionTypeBase

	// Add name field
	tagsCollection.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
		Min:      1,
		Max:      100,
	})

	// Add slug field with unique constraint
	tagsCollection.Fields.Add(&core.TextField{
		Name:     "slug",
		Required: true,
		Min:      1,
		Max:      100,
		Pattern:  "^[a-z0-9]+(?:-[a-z0-9]+)*$",
	})

	// Add unique index for slug
	tagsCollection.Indexes = []string{
		"CREATE UNIQUE INDEX idx_tags_slug ON tags (slug)",
	}

	// Set API rules for public read access
	tagsCollection.ListRule = types.Pointer("")
	tagsCollection.ViewRule = types.Pointer("")

	// Save tags collection
	if err := txApp.Save(tagsCollection); err != nil {
		return fmt.Errorf("failed to create tags collection: %w", err)
	}

	// Get the saved tags collection to retrieve its ID
	savedTags, err := txApp.FindCollectionByNameOrId("tags")
	if err != nil {
		return fmt.Errorf("failed to find tags collection: %w", err)
	}

	// Get users collection ID (PocketBase default)
	usersCollection, err := txApp.FindCollectionByNameOrId("users")
	if err != nil {
		// Users collection may not exist yet, that's okay
		usersCollection = nil
	}

	// Create links collection
	linksCollection := core.NewBaseCollection("links")
	linksCollection.Type = core.CollectionTypeBase

	// Add url field
	linksCollection.Fields.Add(&core.URLField{
		Name:     "url",
		Required: true,
	})

	// Add title field
	linksCollection.Fields.Add(&core.TextField{
		Name:     "title",
		Required: true,
		Min:      1,
		Max:      500,
	})

	// Add description field
	linksCollection.Fields.Add(&core.TextField{
		Name:     "description",
		Required: false,
		Max:      2000,
	})

	// Add view_count field
	viewCountMin := 0.0
	linksCollection.Fields.Add(&core.NumberField{
		Name:     "view_count",
		Required: false,
		Min:      &viewCountMin,
		OnlyInt:  true,
	})

	// Add tags relation
	linksCollection.Fields.Add(&core.RelationField{
		Name:         "tags",
		Required:     false,
		CollectionId: savedTags.Id,
		MaxSelect:    100,
	})

	// Add created_by relation (only if users collection exists)
	if usersCollection != nil {
		linksCollection.Fields.Add(&core.RelationField{
			Name:         "created_by",
			Required:     false,
			CollectionId: usersCollection.Id,
			MaxSelect:    1,
		})
	}

	// Set API rules for public read access
	linksCollection.ListRule = types.Pointer("")
	linksCollection.ViewRule = types.Pointer("")

	// Save links collection
	if err := txApp.Save(linksCollection); err != nil {
		return fmt.Errorf("failed to create links collection: %w", err)
	}

	// Create sync_sessions collection
	if err := createSyncSessionsCollection(txApp); err != nil {
		return err
	}

	// Create presentations collection
	return createPresentationsCollection(txApp)
}

// createSyncSessionsCollection creates the sync_sessions collection
func createSyncSessionsCollection(txApp core.App) error {
	// Check if sync_sessions collection already exists
	if _, err := txApp.FindCollectionByNameOrId("sync_sessions"); err == nil {
		// Collection already exists, skip
		return nil
	}

	// Create sync_sessions collection
	syncSessionsCollection := core.NewBaseCollection("sync_sessions")
	syncSessionsCollection.Type = core.CollectionTypeBase

	// Add progress field (float between 0 and 1)
	progressMin := 0.0
	progressMax := 1.0
	syncSessionsCollection.Fields.Add(&core.NumberField{
		Name:     "progress",
		Required: false,
		Min:      &progressMin,
		Max:      &progressMax,
		OnlyInt:  false,
		System:   false,
	})
	// Default value would be set here if the field supported it
	// For now, we'll handle defaults in the application layer

	// Add admin_token field (required text field)
	syncSessionsCollection.Fields.Add(&core.TextField{
		Name:     "admin_token",
		Required: true,
		Min:      64, // Hex-encoded 32 bytes = 64 chars
		Max:      64,
	})

	// Set API rules for public read access
	syncSessionsCollection.ListRule = types.Pointer("") // Allow all to list
	syncSessionsCollection.ViewRule = types.Pointer("") // Allow all to view
	// No create, update, or delete rules - these will be handled by custom routes
	syncSessionsCollection.CreateRule = nil
	syncSessionsCollection.UpdateRule = nil
	syncSessionsCollection.DeleteRule = nil

	// Save sync_sessions collection
	if err := txApp.Save(syncSessionsCollection); err != nil {
		return fmt.Errorf("failed to create sync_sessions collection: %w", err)
	}

	return nil
}

// createPresentationsCollection creates the presentations collection
func createPresentationsCollection(txApp core.App) error {
	// Check if presentations collection already exists
	if _, err := txApp.FindCollectionByNameOrId("presentations"); err == nil {
		// Collection already exists, skip
		return nil
	}

	// Find sync_sessions collection (required dependency)
	syncSessionsCollection, err := txApp.FindCollectionByNameOrId("sync_sessions")
	if err != nil {
		return fmt.Errorf("failed to find sync_sessions collection: %w", err)
	}

	// Find users collection (optional dependency)
	usersCollection, err := txApp.FindCollectionByNameOrId("users")
	if err != nil {
		// Users collection may not exist, that's okay
		usersCollection = nil
	}

	// Create presentations collection
	presentationsCollection := core.NewBaseCollection("presentations")
	presentationsCollection.Type = core.CollectionTypeBase

	// Add name field
	presentationsCollection.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
		Min:      1,
		Max:      255,
	})

	// Add step_count field
	stepCountMin := 1.0
	presentationsCollection.Fields.Add(&core.NumberField{
		Name:     "step_count",
		Required: true,
		Min:      &stepCountMin,
		OnlyInt:  true,
	})

	// Add step_labels field (JSON array)
	presentationsCollection.Fields.Add(&core.JSONField{
		Name:     "step_labels",
		Required: false,
	})

	// Add active_session relation
	presentationsCollection.Fields.Add(&core.RelationField{
		Name:         "active_session",
		Required:     false,
		CollectionId: syncSessionsCollection.Id,
		MaxSelect:    1,
	})

	// Add created_by relation (only if users collection exists)
	if usersCollection != nil {
		presentationsCollection.Fields.Add(&core.RelationField{
			Name:         "created_by",
			Required:     false,
			CollectionId: usersCollection.Id,
			MaxSelect:    1,
		})
	}

	// Set API rules
	presentationsCollection.ListRule = types.Pointer("")  // Anyone can list
	presentationsCollection.ViewRule = types.Pointer("")  // Anyone can view
	presentationsCollection.CreateRule = types.Pointer("@request.auth.id != ''")  // Authenticated users only
	presentationsCollection.UpdateRule = types.Pointer("@request.auth.id = created_by")  // Owner only
	presentationsCollection.DeleteRule = types.Pointer("@request.auth.id = created_by")  // Owner only

	// Save presentations collection
	if err := txApp.Save(presentationsCollection); err != nil {
		return fmt.Errorf("failed to create presentations collection: %w", err)
	}

	return nil
}