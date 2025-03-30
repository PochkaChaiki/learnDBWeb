package service

import "learnDB/internal/web/storage"

type APIService struct {
	ServiceAnswer
	ServiceQuery
	ServiceQuestion
}

func New(s *storage.Storage) *APIService {
	return &APIService{
		*NewServiceAnswer(s.AnswerStorage),
		*NewServiceQuery(s.QueryStorage),
		*NewServiceQuestion(s.QuestionStorage),
	}
}
