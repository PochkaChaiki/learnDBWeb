package storage

import (
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
}

func New(filepath string) *Storage {
	db := sqlx.MustConnect("sqlite3", filepath)

	return &Storage{DB: db}
}
