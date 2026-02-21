package main

import (
	"log"

	"github.com/jchen/go-sql-practice/migrations"
	"github.com/jchen/go-sql-practice/routes"
	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()

	// Register migrations
	migrations.Register(app)

	// Register custom routes
	routes.Register(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}