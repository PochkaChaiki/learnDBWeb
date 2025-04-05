package testReviewer

import "learnDB/internal/dbRepository/domain"

type DBRepository interface {
	RunScript(sql string, limit int) (*domain.QueryResult, error)
}
