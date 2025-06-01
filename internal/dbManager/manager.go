package dbManager

import (
	"errors"
	"learnDB/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBManager struct {
	repos map[string]domain.DBRepository
}

func New() *DBManager {
	return &DBManager{
		repos: make(map[string]domain.DBRepository),
	}
}

func (d *DBManager) AddRepository(driver string, connStr string) error {
	var controlStmt string
	switch driver {
	case "postgres":
		controlStmt = "set search_path to %s"
	case "mysql":
		controlStmt = "use %s"
	default:
		return errors.New("database name is uncrecognizable")
	}
	db, err := sqlx.Connect(driver, connStr)
	if err != nil {
		return err
	}
	if driver == "mysql" {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(100)
	}

	d.repos[driver] = &DBRepository{
		db:                 db,
		controlStmt:        controlStmt,
		availableConnToken: make(chan bool, 100),
	}
	return nil
}

func (d *DBManager) GetRepository(driver string) (domain.DBRepository, bool) {
	repo, ok := d.repos[driver]
	return repo, ok
}
