package db

import (
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *Store

const (
	dbDriver = "postgres"
	dbstring = "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" // move to config later
)

func TestMain(m *testing.M) {
	// config, err := util.LoadConfig("../..")
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }

	// testDB, err = sql.Open(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal("cannot connect to db:", err)
	// }
	var err error
	testDB, err = NewStore(dbDriver, dbstring)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	os.Exit(m.Run())
}
