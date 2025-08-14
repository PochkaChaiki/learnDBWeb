package testreview

import "github.com/pochkachaiki/learndb/internal/domain"

type DBRepository interface {
	RunSelect(sql string, schemaName string, limit int) (*domain.QueryResult, error)
}
