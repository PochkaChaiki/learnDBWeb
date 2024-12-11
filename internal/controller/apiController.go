package controller

import "learnDB/internal/service"

type APIController struct {
	AnswerController
	QueryController
	QuestionController
}

func New(s *service.APIService) *APIController {
	return &APIController{
		*NewAnswerController(&s.ServiceAnswer),
		*NewQueryController(&s.ServiceQuery),
		*NewQuestionController(&s.ServiceQuestion),
	}
}
