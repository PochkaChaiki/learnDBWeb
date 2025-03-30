package excelReader

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// type QueryResult struct {
// 	columns []string
// 	data    [][]any
// }

type DBManager struct {
	db   *sqlx.DB
	stmt *sqlx.Stmt
}

func NewDB(db *sqlx.DB) (*DBManager, error) {
	stmt, err := db.Preparex("select * from (?) s limit ?;")
	if err != nil {
		return nil, err
	}
	return &DBManager{db, stmt}, nil
}

func (d *DBManager) RunScript(sql string, limit int) (*QueryResult, error) {

	rows, err := d.stmt.Queryx(sql, limit)
	if err != nil {
		return nil, fmt.Errorf("RunScript query error: %w", err)
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("RunScript columns retrieving error: %w", err)
	}

	data := make([][]any, 0, limit)
	defer rows.Close()

	err = nil
	for rows.Next() {
		row, errNew := rows.SliceScan()
		if errNew != nil {
			errors.Join(err, errNew)
		}
		data = append(data, row)
	}

	return &QueryResult{Columns: cols, Data: data}, err

}
