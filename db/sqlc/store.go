package db

import (
	"cnfs/domain"
	"database/sql"
)

type Store interface {
	Querier
}

type PostgresqlStore struct {
	db  *sql.DB
	cfg domain.IConfig
	*Queries
}

func NewStore(db *sql.DB, cfg domain.IConfig) Store {
	return &PostgresqlStore{
		db:      db,
		cfg:     cfg,
		Queries: New(db),
	}
}
