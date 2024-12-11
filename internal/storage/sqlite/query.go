package sqlite

import (
	"fmt"
	"learnDB/internal/domain"

	"github.com/jmoiron/sqlx"
)

type QueryStorage struct {
	db *sqlx.DB
}

func NewQueryStorage(db *sqlx.DB) *QueryStorage {
	return &QueryStorage{db: db}
}

func (r *QueryStorage) Insert(q *domain.Query) error {
	_, err := r.db.NamedExec("insert into query(script, info, executed_at, user_id, db_id) values (:script, :info, :executed_at, :user_id, :db_id);", q)
	if err != nil {
		return fmt.Errorf("insert query error: %s", err)
	}
	return nil
}

func (r *QueryStorage) Get(id int) (*domain.Query, error) {
	q := new(domain.Query)
	if err := r.db.Get(q, "select query_id, script, info, executed_at, user_id, db_id from query where query_id = ?;", id); err != nil {
		return nil, fmt.Errorf("select query error: %s", err)
	}
	return q, nil
}

func (r *QueryStorage) GetAll() (queries []domain.Query, err error) {

	if err := r.db.Select(&queries, "select query_id, script, info, executed_at, user_id, db_id from query;"); err != nil {
		return nil, fmt.Errorf("select queries error: %s", err)
	}
	return queries, nil
}

func (r *QueryStorage) Update(q *domain.Query) error {
	_, err := r.db.NamedExec("update query set script=:script, info=:info, executed_at=:executed_at, user_id=:user_id, db_id=:db_id where query_id=:query_id", q)
	if err != nil {
		return fmt.Errorf("update query error: %s", err)
	}
	return nil
}

func (r *QueryStorage) Delete(id int) error {
	if _, err := r.db.Exec("delete from query where query_id = ?", id); err != nil {
		return fmt.Errorf("delete query error: %s", err)
	}
	return nil
}
