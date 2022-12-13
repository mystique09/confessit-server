package db

import (
	"cnfs/config"
	"database/sql"
)

type Store interface {
	Querier
}

type PostgresqlStore struct {
	db  *sql.DB
	cfg *config.Config
	*Queries
}

func NewStore(db *sql.DB, cfg *config.Config) Store {
	return &PostgresqlStore{
		db:      db,
		cfg:     cfg,
		Queries: New(db),
	}
}
