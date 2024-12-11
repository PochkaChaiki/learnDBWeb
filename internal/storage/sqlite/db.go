package sqlite

import (
	"fmt"
	"learnDB/internal/domain"

	"github.com/jmoiron/sqlx"
)

type DBStorage struct {
	db *sqlx.DB
}

func NewDBStorage(db *sqlx.DB) *DBStorage {
	return &DBStorage{db: db}
}
func (r *DBStorage) Insert(db *domain.DB) error {
	_, err := r.db.NamedExec("insert into db(db_name) values (:db_name);", db)
	if err != nil {
		return fmt.Errorf("insert db error: %s", err)
	}
	return nil
}

func (r *DBStorage) Get(id int) (db *domain.DB, err error) {

	if err = r.db.Get(db, "select db_id, db_name from db where db_id = ?;", id); err != nil {
		return nil, fmt.Errorf("select db error: %s", err)
	}
	return db, nil
}

func (r *DBStorage) GetAll() (dbs []domain.DB, err error) {

	if err = r.db.Select(&dbs, "select db_id, db_name from db;"); err != nil {
		return nil, fmt.Errorf("select dbs error: %s", err)
	}
	return dbs, nil
}

func (r *DBStorage) Update(db *domain.DB) error {
	_, err := r.db.NamedExec("update db set db_name=:db_name where db_id=:db_id", db)
	if err != nil {
		return fmt.Errorf("update db error: %s", err)
	}
	return nil
}

func (r *DBStorage) Delete(id int) error {
	if _, err := r.db.Exec("delete from db where db_id = ?", id); err != nil {
		return fmt.Errorf("delete db error: %s", err)
	}
	return nil
}
