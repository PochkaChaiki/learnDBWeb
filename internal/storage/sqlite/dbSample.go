package sqlite

import (
	"fmt"
	"learnDB/internal/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DBSampleStorage struct {
	db *sqlx.DB
}

func NewDBSampleStorage(db *sqlx.DB) *DBSampleStorage {
	return &DBSampleStorage{db: db}
}

func (r *DBSampleStorage) Insert(dbs *domain.DBSample) error {
	_, err := r.db.NamedExec("insert into db_sample(description, filepath, db_id) values (:description, :filepath, :db_id);", dbs)
	if err != nil {
		return fmt.Errorf("insert dbsample error: %s", err)
	}
	return nil
}

func (r *DBSampleStorage) Get(id int) (*domain.DBSample, error) {
	dbs := new(domain.DBSample)
	if err := r.db.Get(dbs, "select db_sample_id, description, filepath, db_id from db_sample where db_sample_id = ?;", id); err != nil {
		return nil, fmt.Errorf("select dbsample error: %s", err)
	}
	return dbs, nil
}

func (r *DBSampleStorage) GetAll() (dbSamples []domain.DBSample, err error) {

	if err := r.db.Select(&dbSamples, "select db_sample_id, description, filepath, db_id from db_sample;"); err != nil {
		return nil, fmt.Errorf("select dbsamples error: %s", err)
	}
	return dbSamples, nil
}

func (r *DBSampleStorage) Update(dbs *domain.DBSample) error {
	_, err := r.db.NamedExec("update db_sample set description=:description, filepath=:filepath, db_id=:db_id where db_sample_id=:db_sample_id", dbs)
	if err != nil {
		return fmt.Errorf("update dbsample error: %s", err)
	}
	return nil
}

func (r *DBSampleStorage) Delete(id int) error {
	if _, err := r.db.Exec("delete from db_sample where db_sample_id = ?", id); err != nil {
		return fmt.Errorf("delete dbsample error: %s", err)
	}
	return nil
}
