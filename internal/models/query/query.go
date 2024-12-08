package query

import (
	"fmt"
	"learnDB/internal/storage"
	"time"
)

type Query struct {
	Id         int       `json:"-" db:"query_id"`
	Script     string    `json:"script" db:"script"`
	Info       string    `json:"info" db:"info"`
	ExecutedAt time.Time `json:"executed_at" db:"executed_at"`
	UserId     int       `json:"-" db:"user_id"`
	DBId       int       `json:"-" db:"db_id"`
}

func (q *Query) Insert(st *storage.Storage) error {
	_, err := st.DB.NamedExec("insert into query(script, info, executed_at, user_id, db_id) values (:script, :info, :executed_at, :user_id, :db_id);", q)
	if err != nil {
		return fmt.Errorf("insert query error: %s", err)
	}
	return nil
}

func (q *Query) Get(st *storage.Storage, id int) error {

	if err := st.DB.Get(q, "select query_id, script, info, executed_at, user_id, db_id from query where query_id = ?;", id); err != nil {
		return fmt.Errorf("select query error: %s", err)
	}
	return nil
}

func GetAll(st *storage.Storage) ([]Query, error) {

	querys := []Query{}

	if err := st.DB.Select(&querys, "select query_id, script, info, executed_at, user_id, db_id from query;"); err != nil {
		return nil, fmt.Errorf("select queries error: %s", err)
	}
	return querys, nil
}

func (q *Query) Update(st *storage.Storage) error {
	_, err := st.DB.NamedExec("update query set script=:script, info=:info, executed_at=:executed_at, user_id=:user_id, db_id=:db_id where query_id=:query_id", q)
	if err != nil {
		return fmt.Errorf("update query error: %s", err)
	}
	return nil
}

func (q *Query) Delete(st *storage.Storage) error {
	if _, err := st.DB.Exec("delete from query where query_id = ?", q.Id); err != nil {
		return fmt.Errorf("delete query error: %s", err)
	}
	return nil
}
