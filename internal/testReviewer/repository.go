package testReviewer

import "learnDB/internal/dbRepository/domain"

type Repository interface {
	RunScript(sql string, limit int) (*domain.QueryResult, error)
}
