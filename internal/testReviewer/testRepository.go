package testReviewer

import (
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/domain"
)

type TestRepository interface {
	Read(xlsx *config.ExcelConfig) ([]domain.StudentWork, error)
	Write(rw domain.ReviewedWorks) error
}
