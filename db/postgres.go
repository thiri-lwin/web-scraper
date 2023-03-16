package db

import (
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register postgres driver
)

type Store struct {
	db *sqlx.DB
}

func NewStore(dbDriver string, dbstring string) (*Store, error) {
	db, err := sqlx.Connect(dbDriver, dbstring)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to postgres '%s'", dbstring)
	}
	// The current max_connections in postgres is 100.
	db.SetMaxOpenConns(5000)
	db.SetMaxIdleConns(10)
	return &Store{db: db}, nil
}
