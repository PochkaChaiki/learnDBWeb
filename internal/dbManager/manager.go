package dbManager

import (
	"errors"
	"fmt"

	domain "learnDB/internal/dbManager/domain"

	"github.com/jmoiron/sqlx"
)

type DBManager struct {
	db         *sqlx.DB
	selectStmt *sqlx.Stmt
}

func NewMySQL(db *sqlx.DB) (*DBManager, error) {

	selectStmt, err := db.Preparex("select * from (?) s limit ?;")
	if err != nil {
		return nil, err
	}
	return &DBManager{db, selectStmt}, nil
}

func NewPostgreSQL(db *sqlx.DB) (*DBManager, error) {

	selectStmt, err := db.Preparex("select * from (?) s limit ?;")
	if err != nil {
		return nil, err
	}
	return &DBManager{db, selectStmt}, nil
}

func (d *DBManager) prepareQueryResult(rows *sqlx.Rows, limit int) (*domain.QueryResult, error) {

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("columns retrieving error: %w", err)
	}

	data := make([][]any, 0, limit)

	err = nil
	for rows.Next() {
		row, errNew := rows.SliceScan()
		if errNew != nil {
			errors.Join(err, errNew)
		}
		data = append(data, row)
	}

	return &domain.QueryResult{Columns: cols, Data: data}, err

}

func (d *DBManager) RunSelect(sql string, limit int) (*domain.QueryResult, error) {

	rows, err := d.selectStmt.Queryx(sql, limit)
	if err != nil {
		return nil, fmt.Errorf("RunScript query error: %w", err)
	}

	defer rows.Close()

	return d.prepareQueryResult(rows, limit)

}
