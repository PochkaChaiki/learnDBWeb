package domain

import "learnDB/internal/domain/answer"

type CheckResult struct {
	Points int
	Error  error
}

// Struct to check it in "checkAnswer" function
type CheckUnit struct {
	Script         string
	CorrectAnswers []answer.CorrectAnswer
	Review         CheckResult
}
