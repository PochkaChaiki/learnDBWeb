package testReviewer

import "learnDB/internal/testReviewer/domain"

type TestReader interface {
	Read(bookname string) (domain.StudentWork, error)
	Write(rw domain.ReviewedStudentWork, bookname string) error
}
