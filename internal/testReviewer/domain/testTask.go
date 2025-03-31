package domain

import "learnDB/internal/domain/answer"

type TestTask struct {
	QuestionText   string
	CorrectAnswers []answer.CorrectAnswer
	Answer         string
}
