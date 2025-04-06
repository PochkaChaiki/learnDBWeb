package testReviewer

import (
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/domain"
)

type TestService interface {
	Read(bookname string, sheetname string, xlsx *config.ExcelConfig) ([]domain.StudentWork, error)
	Write(rw domain.ReviewedStudentWork, bookname string) error
}
