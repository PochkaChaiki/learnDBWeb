package testReviewer

import "learnDB/internal/dbManager/domain"

type DBRepository interface {
	RunSelect(sql string, limit int) (*domain.QueryResult, error)
}
