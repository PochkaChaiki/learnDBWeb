package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sqlx.DB
}

func New(filepath string) *Storage {
	db := sqlx.MustConnect("sqlite3", filepath)

	return &Storage{DB: db}
}
