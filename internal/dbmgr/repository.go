package dbmgr

import (
	"context"
	"errors"
	"fmt"

	"github.com/pochkachaiki/learndb/internal/domain"

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
		return nil, err
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

func (d *DBRepository) RunSelectContext(ctx context.Context, sql string, schemaName string, limit int) (*domain.QueryResult, error) {
	d.availableConnToken <- true
	defer func() {
		<-d.availableConnToken
	}()

	if schemaName != "" {
		_, err := d.db.ExecContext(ctx, fmt.Sprintf(d.controlStmt, schemaName))
		if err != nil {
			return nil, err
		}
	}
	rows, err := d.db.QueryxContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return d.prepareQueryResult(rows, limit)
}
