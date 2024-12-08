package db

import (
	"fmt"
	"learnDB/internal/storage"
)

type DB struct {
	Id   int    `json:"-" db:"db_id"`
	Name string `json:"db_name" db:"db_name"`
}

func (db *DB) Insert(st *storage.Storage) error {
	_, err := st.DB.NamedExec("insert into db(db_name) values (:db_name);", db)
	if err != nil {
		return fmt.Errorf("insert db error: %s", err)
	}
	return nil
}

func (db *DB) Get(st *storage.Storage, id int) error {

	if err := st.DB.Get(db, "select db_id, db_name from db where db_id = ?;", id); err != nil {
		return fmt.Errorf("select db error: %s", err)
	}
	return nil
}

func GetAll(st *storage.Storage) ([]DB, error) {

	dbs := []DB{}

	if err := st.DB.Select(&dbs, "select db_id, db_name from db;"); err != nil {
		return nil, fmt.Errorf("select dbs error: %s", err)
	}
	return dbs, nil
}

func (db *DB) Update(st *storage.Storage) error {
	_, err := st.DB.NamedExec("update db set db_name=:db_name where db_id=:db_id", db)
	if err != nil {
		return fmt.Errorf("update db error: %s", err)
	}
	return nil
}

func (db *DB) Delete(st *storage.Storage) error {
	if _, err := st.DB.Exec("delete from db where db_id = ?", db.Id); err != nil {
		return fmt.Errorf("delete db error: %s", err)
	}
	return nil
}
