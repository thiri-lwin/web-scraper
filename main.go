package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thiri-lwin/web_scraper/api"
	"github.com/thiri-lwin/web_scraper/db"
	"github.com/thiri-lwin/web_scraper/util"
)

// const (
//
//	dbstring = "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" // move to config later
//
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

	keywordChan := make(chan util.UploadedFile)
	// Define a background job scheduler to scrape keywords
	go scrapKeyword(store, keywordChan)

	templatePath := "templates/*.html"
	assetsPath := "/assets"
	cssPath := "templates/css"

	server := api.NewServer(store, keywordChan, templatePath, assetsPath, cssPath)
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

func scrapKeyword(store *db.Store, keywordChan chan util.UploadedFile) {
	for input := range keywordChan {
		records := input.Records
		userID := input.UserID

		// assume no header and all records are keywords that need to search
		resultCh := make(chan util.SearchInfo)
		var wg sync.WaitGroup
		for _, line := range records {
			for _, keyword := range line {
				keyword = strings.TrimSpace(keyword)
				if keyword == "" {
					continue
				}
				keywordInfo := util.SearchInfo{
					Keyword: keyword,
					Status:  "pending",
				}
				wg.Add(1)
				go util.SearchKeyword(keywordInfo, resultCh, &wg)
			}
		}

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		var allResults []util.SearchInfo

		for result := range resultCh {
			allResults = append(allResults, result)

			// store results in db
			dbRes := util.ResultDB{
				HTMLCode:           result.HTMLCode,
				NumAds:             result.NumAds,
				NumLinks:           result.NumLinks,
				TotalSearchResults: result.TotalSearchResults,
			}
			data, _ := json.Marshal(dbRes)
			_, err := store.InsertSearchResult(context.Background(), db.SearchResult{
				UserID:  userID,
				Keyword: result.Keyword,
				Results: string(data),
				Status:  result.Status,
			})
			if err != nil {
				log.Println("Failed to save result in db:", err)
			}

		}

		fmt.Println("len of all results:", len(allResults))

	}
}
