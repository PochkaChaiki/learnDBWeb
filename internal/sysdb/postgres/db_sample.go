package postgres

import (
	"fmt"

	"github.com/pochkachaiki/learndb/internal/domain"
)

func (r *Repository) CreateDBSample(sample *domain.DBSample) error {
	const query = `INSERT INTO testing_info.db_sample 
                  (name, path, description, dbms) 
                  VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query,
		sample.Name, sample.Path, sample.Description)
	if err != nil {
		return fmt.Errorf("failed to create db sample: %w", err)
	}

	return nil
}

func (r *Repository) ReadDBSamplesByTaskID(taskID int) ([]domain.DBSample, error) {
	const query = `SELECT ds.id, ds.name, ds.path, ds.description, ds.dbms
                  FROM testing_info.db_sample ds
                  JOIN testing_info.task_db_sample tds ON ds.id = tds.db_sample_id
                  WHERE tds.task_id = $1`

	var samples []domain.DBSample
	err := r.db.Select(&samples, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get db samples: %w", err)
	}

	return samples, nil
}

func (r *Repository) UpdateDBSample(sample *domain.DBSample) error {
	const query = `UPDATE testing_info.db_sample 
                  SET path = $1, description = $2, dbms = $3
                  WHERE name = $4`

	_, err := r.db.Exec(query,
		sample.Path, sample.Description, sample.DBMS, sample.Name)
	if err != nil {
		return fmt.Errorf("failed to update db sample: %w", err)
	}

	return nil
}

func (r *Repository) DeleteDBSample(name string) error {
	const query = `DELETE FROM testing_info.db_sample WHERE name = $1`

	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete db sample: %w", err)
	}

	return nil
}
