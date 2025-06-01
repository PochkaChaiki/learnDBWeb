package scriptManager

import "learnDB/internal/domain"

type DBRepository interface {
	RunSelect(sql string, schemaName string, limit int) (*domain.QueryResult, error)
}
