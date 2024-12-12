package service

import (
	"learnDB/internal/domain"
	"log"
)

// GET /api/answer
// GET /api/answer/{id}
// POST /api/answer --> /api/query
// DELETE /api/answer/{id}

type AnswerStorage interface {
	Insert(*domain.Answer) error
	Get(int) (*domain.Answer, error)
	GetAll() ([]domain.Answer, error)
	Delete(int) error
}

type ServiceAnswer struct {
	storage AnswerStorage
}

func NewServiceAnswer(s AnswerStorage) *ServiceAnswer {
	return &ServiceAnswer{storage: s}
}

func (srv *ServiceAnswer) Create(ans *domain.Answer) OperationResult {
	if err := srv.storage.Insert(ans); err != nil {
		log.Printf("answer service create error: %s", err)
		return InternalError
	}
	return Ok
}

func (srv *ServiceAnswer) GetAll() ([]domain.Answer, OperationResult) {
	anses, err := srv.storage.GetAll()
	if err != nil {
		log.Printf("answer service get all error: %s", err)
		return nil, InternalError
	}
	return anses, Ok
}

func (srv *ServiceAnswer) Get(id int) (*domain.Answer, OperationResult) {
	ans, err := srv.storage.Get(id)
	if err != nil {
		log.Printf("answer service get error: %s", err)
		return nil, InternalError
	}
	return ans, Ok
}

func (srv *ServiceAnswer) Delete(id int) OperationResult {
	if ans, err := srv.storage.Get(id); err != nil {
		log.Printf("answer service delete error: %s", err)
		return InternalError
	} else if ans == nil {
		return BadRequest
	}

	if err := srv.storage.Delete(id); err != nil {
		log.Printf("answer service delete error: %s", err)
		return InternalError
	}

	return Ok
}
