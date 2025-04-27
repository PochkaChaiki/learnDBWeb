package testReviewer

import "learnDB/internal/dbRepository/domain"

type DBRepository interface {
	RunSelect(sql string, limit int) (*domain.QueryResult, error)
}
