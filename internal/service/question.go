package service

import (
	"learnDB/internal/domain"
	"log"
)

// GET /api/question
// GET /api/question/{id}
// POST /api/question
// PUT /api/question
// DELETE /api/question/{id}

type QuestionStorage interface {
	Insert(*domain.Question) error
	Get(int) (*domain.Question, error)
	GetAll() ([]domain.Question, error)
	Update(*domain.Question) error
	Delete(int) error
}

type ServiceQuestion struct {
	storage QuestionStorage
}

func NewServiceQuestion(s QuestionStorage) *ServiceQuestion {
	return &ServiceQuestion{storage: s}
}

func (srv *ServiceQuestion) Create(q *domain.Question) OperationResult {

	if err := srv.storage.Insert(q); err != nil {
		log.Fatalf("question service error: %s", err)
		return InternalError
	}
	return Ok
}

func (srv *ServiceQuestion) Get(id int) (*domain.Question, OperationResult) {

	q, err := srv.storage.Get(id)
	if err != nil {
		log.Fatalf("question service get error: %s", err)
		return nil, InternalError
	}
	return q, Ok
}

func (srv *ServiceQuestion) GetAll() ([]domain.Question, OperationResult) {
	qs, err := srv.storage.GetAll()
	if err != nil {
		log.Fatalf("question service get all error: %s", err)
		return nil, InternalError
	}
	return qs, Ok
}

func (srv *ServiceQuestion) Update(q *domain.Question) OperationResult {

	if question, err := srv.storage.Get(q.Id); err != nil {
		log.Fatalf("question service update error: %s", err)
		return InternalError
	} else if question == nil {
		return BadRequest
	}

	if err := srv.storage.Update(q); err != nil {
		log.Fatalf("question service update error: %s", err)
		return InternalError
	}
	return Ok
}

func (srv *ServiceQuestion) Delete(id int) OperationResult {

	if q, err := srv.storage.Get(id); err != nil {
		log.Fatalf("question service delete error: %s", err)
		return InternalError
	} else if q == nil {
		return BadRequest
	}

	if err := srv.storage.Delete(id); err != nil {
		log.Fatalf("question service delete error: %s", err)
		return InternalError
	}

	return Ok
}
