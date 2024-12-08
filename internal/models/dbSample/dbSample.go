package dbSample

import (
	"fmt"
	"learnDB/internal/storage"
)

type DBSample struct {
	Id          int    `json:"-" db:"db_sample_id"`
	Description string `json:"description" db:"description"`
	Filepath    string `json:"-" db:"filepath"`
	DBId        int    `json:"-" db:"db_id"`
}

func (dbs *DBSample) Insert(st *storage.Storage) error {
	_, err := st.DB.NamedExec("insert into db_sample(description, filepath, db_id) values (:description, :filepath, :db_id);", dbs)
	if err != nil {
		return fmt.Errorf("insert dbsample error: %s", err)
	}
	return nil
}

func (dbs *DBSample) Get(st *storage.Storage, id int) error {

	if err := st.DB.Get(dbs, "select db_sample_id, description, filepath, db_id from db_sample where db_sample_id = ?;", id); err != nil {
		return fmt.Errorf("select dbsample error: %s", err)
	}
	return nil
}

func GetAll(st *storage.Storage) ([]DBSample, error) {

	dbsamples := []DBSample{}

	if err := st.DB.Select(&dbsamples, "select db_sample_id, description, filepath, db_id from db_sample;"); err != nil {
		return nil, fmt.Errorf("select dbsamples error: %s", err)
	}
	return dbsamples, nil
}

func (dbs *DBSample) Update(st *storage.Storage) error {
	_, err := st.DB.NamedExec("update db_sample set description=:description, filepath=:filepath, db_id=:db_id where db_sample_id=:db_sample_id", dbs)
	if err != nil {
		return fmt.Errorf("update dbsample error: %s", err)
	}
	return nil
}

func (dbs *DBSample) Delete(st *storage.Storage) error {
	if _, err := st.DB.Exec("delete from db_sample where db_sample_id = ?", dbs.Id); err != nil {
		return fmt.Errorf("delete dbsample error: %s", err)
	}
	return nil
}
