package storage

import "learnDB/internal/domain"

type QuestionStorage interface {
	Insert(*domain.Question) error
	Get(int) (*domain.Question, error)
	GetAll() ([]domain.Question, error)
	Update(*domain.Question) error
	Delete(int) error
}

type AnswerStorage interface {
	Insert(*domain.Answer) error
	Get(int) (*domain.Answer, error)
	GetAll() ([]domain.Answer, error)
	Update(*domain.Answer) error
	Delete(int) error
}

type DBStorage interface {
	Insert(*domain.DB) error
	Get(int) (*domain.DB, error)
	GetAll() ([]domain.DB, error)
	Update(*domain.DB) error
	Delete(int) error
}

type DBSampleStorage interface {
	Insert(*domain.DBSample) error
	Get(int) (*domain.DBSample, error)
	GetAll() ([]domain.DBSample, error)
	Update(*domain.DBSample) error
	Delete(int) error
}

type QueryStorage interface {
	Insert(*domain.Query) error
	Get(int) (*domain.Query, error)
	GetAll() ([]domain.Query, error)
	Update(*domain.Query) error
	Delete(int) error
}

type UserStorage interface {
	Get(int) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetUserByUsername(string) (*domain.User, error)
}

type Storage struct {
	AnswerStorage
	DBStorage
	DBSampleStorage
	QueryStorage
	QuestionStorage
	UserStorage
}
