package dbManager

import (
	"errors"
	"fmt"

	domain "learnDB/internal/domain"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBManager struct {
	db                 *sqlx.DB
	controlStmt        string
	availableConnToken chan bool
}

func New(driver string, connStr string) (*DBManager, error) {
	var controlStmt string
	switch driver {
	case "postgres":
		controlStmt = "set search_path to %s"
	case "mysql":
		controlStmt = "use %s"
	default:
		return nil, errors.New("database name is uncrecognizable")
	}
	db, err := sqlx.Connect(driver, connStr)
	if err != nil {
		return nil, err
	}
	return &DBManager{
		db:                 db,
		controlStmt:        controlStmt,
		availableConnToken: make(chan bool, 100),
	}, nil
}

func (d *DBManager) prepareQueryResult(rows *sqlx.Rows, limit int) (*domain.QueryResult, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.Join(domain.ErrInternalError, err) // fmt.Errorf("columns retrieving error: %w", err)
	}

	data := make([][]any, 0, limit)

	for rows.Next() {
		row, errNew := rows.SliceScan()
		if errNew != nil {
			errors.Join(err, errNew)
		}
		data = append(data, row)
		limit--
		if limit == 0 {
			break
		}
	}

	return &domain.QueryResult{Columns: cols, Data: data}, err

}

func (d *DBManager) RunSelect(sql string, schemaName string, limit int) (*domain.QueryResult, error) {
	d.availableConnToken <- true
	if schemaName != "" {
		_, err := d.db.Exec(fmt.Sprintf(d.controlStmt, schemaName))
		if err != nil {
			return nil, errors.Join(domain.ErrInternalError, err) // fmt.Errorf("RunScript control query error: %w", err)
		}
	}
	rows, err := d.db.Queryx(sql)
	if err != nil {
		return nil, errors.Join(domain.ErrSyntaxError, err) // fmt.Errorf("RunScript query error: %w", err)
	}
	<-d.availableConnToken

	defer rows.Close()

	return d.prepareQueryResult(rows, limit)

}
