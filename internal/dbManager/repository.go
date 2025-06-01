package dbManager

import (
	"errors"
	"fmt"

	"learnDB/internal/domain"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBRepository struct {
	db                 *sqlx.DB
	controlStmt        string
	availableConnToken chan bool
}

func (d *DBRepository) prepareQueryResult(rows *sqlx.Rows, limit int) (*domain.QueryResult, error) {
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

func (d *DBRepository) RunSelect(sql string, schemaName string, limit int) (*domain.QueryResult, error) {
	d.availableConnToken <- true
	if schemaName != "" {
		_, err := d.db.Exec(fmt.Sprintf(d.controlStmt, schemaName))
		if err != nil {
			<-d.availableConnToken
			return nil, errors.Join(domain.ErrInternalError, err) // fmt.Errorf("RunScript control query error: %w", err)
		}
	}
	rows, err := d.db.Queryx(sql)
	if err != nil {
		<-d.availableConnToken
		return nil, errors.Join(domain.ErrSyntaxError, err) // fmt.Errorf("RunScript query error: %w", err)
	}
	<-d.availableConnToken

	defer rows.Close()

	return d.prepareQueryResult(rows, limit)

}
