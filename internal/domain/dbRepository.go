package domain

type DBRepository interface {
	RunSelect(sql string, schemaName string, limit int) (*QueryResult, error)
}
