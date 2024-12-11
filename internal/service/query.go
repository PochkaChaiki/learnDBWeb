package service

import (
	"learnDB/internal/domain"
	"log"
)

// GET /api/query
// GET /api/query/{id}
// POST /api/query

type QueryStorage interface {
	Get(int) (*domain.Query, error)
	GetAll() ([]domain.Query, error)
	Insert(*domain.Query) error
}

type ServiceQuery struct {
	storage QueryStorage
}

func NewServiceQuery(s QueryStorage) *ServiceQuery {
	return &ServiceQuery{storage: s}
}

func (srv *ServiceQuery) GetAll() ([]domain.Query, OperationResult) {
	qs, err := srv.storage.GetAll()
	if err != nil {
		log.Fatalf("query service get all error: %s", err)
		return nil, InternalError
	}
	return qs, Ok
}

func (srv *ServiceQuery) Get(id int) (*domain.Query, OperationResult) {
	q, err := srv.storage.Get(id)
	if err != nil {
		log.Fatalf("query service get error: %s", err)
		return nil, InternalError
	}
	return q, Ok
}
func (srv *ServiceQuery) Create(q *domain.Query) OperationResult {
	if err := srv.storage.Insert(q); err != nil {
		log.Fatalf("query service error: %s", err)
		return InternalError
	}
	return Ok
}
