package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thiri-lwin/web_scraper/api"
	"github.com/thiri-lwin/web_scraper/db"
	"github.com/thiri-lwin/web_scraper/util"
)

// const (
// 	dbstring = "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" // move to config later
// )

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)
	store, _ := db.NewStore(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	templatePath := "templates/*.html"
	assetsPath := "/assets"
	cssPath := "templates/css"

	server := api.NewServer(store, templatePath, assetsPath, cssPath)
	server.Start(config.HTTPServerAddress)

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up", err)
	}

	log.Println("db migrated successfully")
}
