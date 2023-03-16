package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thiri-lwin/web_scraper/api"
	"github.com/thiri-lwin/web_scraper/db"
)

const (
	dbstring = "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" // move to config later
)

func main() {
	store, _ := db.NewStore("postgres", dbstring)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	templatePath := "templates/*.html"
	assetsPath := "/assets"
	cssPath := "templates/css"

	server := api.NewServer(store, templatePath, assetsPath, cssPath)
	server.Start("0.0.0.0:8080") // move the server address to config later
}
