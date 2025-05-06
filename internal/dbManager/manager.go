package dbManager

import (
	"errors"
	"fmt"

	domain "learnDB/internal/dbManager/domain"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBManager struct {
	db          *sqlx.DB
	selectStmt  string
	controlStmt string
}

func New(driver string, connStr string) (*DBManager, error) {
	var selectStmt string
	var controlStmt string
	switch driver {
	case "postgres":
		controlStmt = "set search_path to %s"
		selectStmt = "select * from (%s) s limit %d;"
	case "mysql":
		controlStmt = "use %s"
		selectStmt = "select * from (%s) s limit %d;"
	case "sqlite":
		selectStmt = "select * from (%s) s limit %d;"
	default:
		return nil, errors.New("database name is uncrecognizable")
	}
	db, err := sqlx.Connect(driver, connStr)
	if err != nil {
		return nil, err
	}
	return &DBManager{db, selectStmt, controlStmt}, nil
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

func (d *DBManager) RunSelect(sql string, schemaName string, limit int) (*domain.QueryResult, error) {
	if schemaName != "" {
		_, err := d.db.Exec(fmt.Sprintf(d.controlStmt, schemaName))
		if err != nil {
			return nil, fmt.Errorf("RunScript control query error: %w", err)
		}
	}
	rows, err := d.db.Queryx(fmt.Sprintf(d.selectStmt, sql, limit))
	if err != nil {
		return nil, fmt.Errorf("RunScript query error: %w", err)
	}

	defer rows.Close()

	return d.prepareQueryResult(rows, limit)

}
